package network

import (
	"crypto/ecdsa"

	"fmt"

	"sync/atomic"

	"sync"

	"errors"

	"time"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/network/nat/gopjnath"
	"github.com/SmartMeshFoundation/raiden-network/network/signal/interface"
	"github.com/SmartMeshFoundation/raiden-network/network/signal/signalshare"
	"github.com/SmartMeshFoundation/raiden-network/network/signal/xmpp"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type IceStatus int

const (
	IceTransporterStateInit = IceStatus(iota)
	IceTranspoterStateInitComplete
	IceTransporterStateNegotiateComplete
	IceTransporterStateClosed
)

const (
	StatusCanReceive  = 0
	StatusStopReceive = 1
	StatusCanSend     = 0
	StatusStopSend    = 1
)

var (
	signalServer = ""
)
var (
	errHasStopped                   = errors.New("ice transporter has stopped")
	errStoppedReceive               = errors.New("ice transporter has stopped receiving")
	errIceStreamTransporterNotReady = errors.New("gopjnath.icestreamtransport not ready to send")
)

func InitIceTransporter(turnServer, turnUser, turnPassowrd, signalServerUrl string) {
	err := gopjnath.IceInit(turnServer, turnServer, turnUser, turnPassowrd)
	if err != nil {
		panic(fmt.Sprintf("init ice error %s", err))
	}
	signalServer = signalServerUrl
}
func newpassword(key *ecdsa.PrivateKey) xmpp.GetCurrentPasswordFunc {
	f1 := func() string {
		pass, _ := signalshare.CreatePassword(key)
		return pass
	}
	return f1
}

type IceStream struct {
	ist            *gopjnath.IceStreamTransport
	Status         IceStatus
	CanExchangeSdp bool
}
type IceTransport struct {
	key                *ecdsa.PrivateKey
	Addr               common.Address
	IceStreamMap       map[common.Address]*IceStream
	lock               sync.Mutex
	receiveStatus      int32
	sendStatus         int32
	signal             SignalInterface.SignalProxy
	protocol           ProtocolReceiver
	name               string //for test
	connLastReceiveMap map[common.Address]time.Time
	checkInterval      time.Duration
	lastCheck          time.Time
}

func NewIceTransporter(key *ecdsa.PrivateKey, name string) *IceTransport {
	it := &IceTransport{
		key:                key,
		receiveStatus:      StatusStopReceive,
		IceStreamMap:       make(map[common.Address]*IceStream),
		connLastReceiveMap: make(map[common.Address]time.Time),
		lastCheck:          time.Time{},
		checkInterval:      time.Minute,
		Addr:               crypto.PubkeyToAddress(key.PublicKey),
		name:               name,
	}
	sp, err := xmpp.NewXmpp(signalServer, it.Addr, newpassword(it.key), func(from common.Address, sdp string) (mysdp string, err error) {
		return it.handleSdpArrived(from, sdp)
	}, name)
	if err != nil {
		panic(fmt.Sprintf("create ice transpoter error %s", err))
	}
	it.signal = sp
	return it
}
func (it *IceTransport) Register(protcol ProtocolReceiver) {
	it.protocol = protcol
}

/*
for connections that don't use for a long time, just to remove.
for connections in use but may be invalid because of network, remove too.
*/
func (it *IceTransport) removeExpiredConnection() {
	if time.Now().Sub(it.lastCheck) < it.checkInterval {
		return
	}
	now := time.Now()
	for r, t := range it.connLastReceiveMap {
		if now.Sub(t) > 2*it.checkInterval {
			is, ok := it.IceStreamMap[r]
			if ok {
				delete(it.IceStreamMap, r)
				is.ist.Destroy()
			}
		}
	}
}

/*
for send one message:
1. check if has a connection,
2. if connection is invalid (maybe on setup), just return fail
3. if connection is valid ,just send
4. if no connection,  try to setup a p2p connection use ice.
*/
func (it *IceTransport) Send(receiver common.Address, host string, port int, data []byte) error {
	it.lock.Lock()
	defer it.lock.Unlock()
	it.removeExpiredConnection()
	log.Trace(fmt.Sprintf("%s send to %s , message=%s,hash=%s\n", it.name, utils.APex(receiver), encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:]))))
	var err error
	if it.sendStatus != StatusCanSend {
		return errHasStopped
	}
	is, ok := it.IceStreamMap[receiver]
	if ok {
		if is.Status != IceTransporterStateNegotiateComplete {
			return errIceStreamTransporterNotReady
		}
		err = is.ist.Send(data)
		return err
	} else { //start new p2p
		is := &IceStream{
			Status: IceTransporterStateInit,
		}
		log.Trace("aaaaa")
		is.ist, err = gopjnath.NewIceStreamTransport(it.name, func(u uint, bytes []byte, addr gopjnath.SockAddr) {
			it.Receive(bytes, "", 0)
		}, func(op gopjnath.IceTransportOp, e error) {
			it.handelIceCompleteForControlling(is, receiver, op, e, data)
		})
		log.Trace("xxxxx")
		if err != nil {
			log.Trace(fmt.Sprintf("%s NewIceStreamTransport err %s", it.name, err))
			return err
		}
		log.Trace("bbbbb")
		it.IceStreamMap[receiver] = is
		err := it.signal.TryReach(receiver)
		if err != nil {
			return err
		}
		log.Trace("cccc")
		is.CanExchangeSdp = true
	}
	return nil
}

type sdpresult struct {
	sdp string
	err error
}

func (it *IceTransport) handleSdpArrived(partner common.Address, sdp string) (mysdp string, err error) {
	it.lock.Lock()
	defer it.lock.Unlock()
	if it.receiveStatus != StatusCanReceive {
		err = errStoppedReceive
		return
	}
	log.Trace(fmt.Sprintf("%s handleSdpArrived from %s", it.name, utils.APex2(partner)))
	is, ok := it.IceStreamMap[partner]
	if ok { //already have a connection, why remote create new connection,  maybe they are trying to send data each  other at the same time.
		err = errors.New(fmt.Sprintf("%s trying to send each other at the same time?", it.name))
		return
	}
	is = &IceStream{
		Status:         IceTransporterStateInit,
		CanExchangeSdp: true,
	}
	sdpchan := make(chan *sdpresult, 1)
	is.ist, err = gopjnath.NewIceStreamTransport(it.name, func(u uint, bytes []byte, addr gopjnath.SockAddr) {
		//log.Trace(fmt.Sprintf("receive:%s", string(bytes)))
		it.Receive(bytes, partner.String(), 0)
	}, func(op gopjnath.IceTransportOp, e error) {
		it.handelIceCompleteForControlled(is, partner, op, e, sdp, sdpchan)
	})
	if err != nil {
		return
	}
	it.IceStreamMap[partner] = is
	res := <-sdpchan
	if res == nil {
		err = errors.New(fmt.Sprintf("%s sdp chan return nil", it.name))
		return
	}
	log.Debug(fmt.Sprintf("%s get sdp:%s err:%s", it.name, res.sdp, res.err))
	return res.sdp, res.err

}
func (it *IceTransport) handelIceCompleteForControlled(is *IceStream, receiver common.Address, op gopjnath.IceTransportOp, err error, partnerSdp string, sdp chan *sdpresult) {
	log.Trace(fmt.Sprintf("%s ice complete callback op=%d,err=%s\n", it.name, op, err))
	if err != nil {
		log.Error(fmt.Sprintf("%s ice complete callback error op=%d, err=%s", it.name, op, err))
		it.removeIceStreamTransport(receiver)
		sdp <- &sdpresult{
			err: errors.New(fmt.Sprintf("%s %s ice complete callback err %s", it.name, utils.APex(receiver), err)),
		}
		return
	}
	if op == gopjnath.IceTransportOpStateInit { //could start exchange sdp
		is.Status = IceTranspoterStateInitComplete
		err := is.ist.InitIceSession(gopjnath.IceSessRoleControlled)
		if err != nil {
			log.Error(fmt.Sprintf("%s %s InitIceSession err ", it.name, utils.APex(receiver), err))
			it.removeIceStreamTransport(receiver)
			sdp <- &sdpresult{
				err: err,
			}
			return
		}
		mysdp, err := is.ist.GetLocalSdp()
		if err != nil {
			log.Error(fmt.Sprintf("%s get sdp error %s", it.name, err))
			it.removeIceStreamTransport(receiver)
			sdp <- &sdpresult{
				err: err,
			}
			return
		}
		err = is.ist.StartIce(partnerSdp)
		if err != nil {
			log.Error(fmt.Sprintf("%s controlled start ice err %s", it.name, err))
			it.removeIceStreamTransport(receiver)
			sdp <- &sdpresult{
				err: err,
			}
			return
		}
		sdp <- &sdpresult{
			sdp: mysdp,
			err: err,
		}
		return
	} else if op == gopjnath.IceTransportOpStateNegotiation {
		is.Status = IceTransporterStateNegotiateComplete
	}
}
func (it *IceTransport) removeIceStreamTransport(receiver common.Address) {
	it.lock.Lock()
	defer it.lock.Unlock()
	is, ok := it.IceStreamMap[receiver]
	if !ok {
		return
	}
	log.Info(fmt.Sprintf("%s removeIceStreamTransport %s", it.name, utils.APex2(receiver)))
	delete(it.IceStreamMap, receiver)
	is.ist.Destroy()
}
func (it *IceTransport) handelIceCompleteForControlling(is *IceStream, receiver common.Address, op gopjnath.IceTransportOp, err error, data []byte) {
	log.Trace(fmt.Sprintf("%s ice complete callback op=%d,err=%s\n", it.name, op, err))
	if err != nil {
		log.Error(fmt.Sprintf("%s ice complete callback error op=%d, err=%s", it.name, op, err))
		it.removeIceStreamTransport(receiver)
	}
	if op == gopjnath.IceTransportOpStateInit { //could start exchange sdp
		is.Status = IceTranspoterStateInitComplete
		if !is.CanExchangeSdp {
			log.Error(fmt.Sprintf("%s ready to negotiate ,but cannot find %s", it.name, utils.APex(receiver)))
			it.removeIceStreamTransport(receiver)
			return
		}
		err := is.ist.InitIceSession(gopjnath.IceSessRoleControlling)
		if err != nil {
			log.Error(fmt.Sprintf("%s %s InitIceSession err ", it.name, utils.APex(receiver), err))
			it.removeIceStreamTransport(receiver)
			return
		}
		sdp, err := is.ist.GetLocalSdp()
		if err != nil {
			log.Error(fmt.Sprintf("%s get sdp error %s", it.name, err))
			it.removeIceStreamTransport(receiver)
			return
		}
		//must not block for pjsip's callback
		go func() {
			partnersdp, err := it.signal.ExchangeSdp(receiver, sdp)
			if err != nil {
				log.Error(fmt.Sprintf("%s exchange sdp error %s", it.name, err))
				it.removeIceStreamTransport(receiver)
				return
			}
			err = is.ist.StartIce(partnersdp)
			if err != nil {
				log.Error(fmt.Sprintf("%s %s StartIce error %s", it.name, utils.APex(receiver), err))
				it.removeIceStreamTransport(receiver)
			}
		}()
	} else if op == gopjnath.IceTransportOpStateNegotiation {
		is.Status = IceTransporterStateNegotiateComplete
		err := is.ist.Send(data)
		if err != nil {
			log.Error(fmt.Sprintf("%s send data err", it.name))
		}
	}
}
func (it *IceTransport) Receive(data []byte, host string, port int) error {
	if it.receiveStatus == StatusStopReceive {
		return errStoppedReceive
	}
	log.Trace(fmt.Sprintf("%s receive message,message=%s,hash=%s\n", it.name, encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data))))
	if it.protocol != nil {
		it.protocol.Receive(data, host, port)
	}
	return nil
}
func (it *IceTransport) Start() {
	it.receiveStatus = StatusCanReceive
}
func (it *IceTransport) Stop() {
	it.StopAccepting()
	atomic.SwapInt32(&it.sendStatus, StatusStopSend)
	it.signal.Close()
	it.lock.Lock()
	for _, i := range it.IceStreamMap {
		i.ist.Destroy()
	}
	it.lock.Unlock()
}
func (it *IceTransport) StopAccepting() {
	atomic.SwapInt32(&it.receiveStatus, StatusStopReceive)
}

type IceHelperDicovery struct {
}

func NewIceHelperDiscovery() *IceHelperDicovery {
	return new(IceHelperDicovery)
}
func (this *IceHelperDicovery) Register(address common.Address, host string, port int) error {
	return nil
}
func (this *IceHelperDicovery) Get(address common.Address) (host string, port int, err error) {
	return address.String(), 0, nil
}
func (this *IceHelperDicovery) NodeIdByHostPort(host string, port int) (node common.Address, err error) {
	return common.HexToAddress(host), nil
}
