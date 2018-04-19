package network

import (
	"crypto/ecdsa"

	"fmt"

	"sync/atomic"

	"sync"

	"errors"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/interface"
	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/signalshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/xmpp"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/ice"
	"net"
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
	turnServer=""
	turnUser=""
	turnPassword=""
	cfg *ice.TransportConfig
)
var (
	errHasStopped                   = errors.New("ice transporter has stopped")
	errStoppedReceive               = errors.New("ice transporter has stopped receiving")
	errIceStreamTransporterNotReady = errors.New("icestreamtransport not ready to send")
)

func InitIceTransporter(turnServer, turnUser, turnPassowrd, signalServerUrl string) {
	signalServer = signalServerUrl
	cfg=ice.NewTransportConfigWithTurn(turnServer,turnUser,turnPassowrd)
}
func newpassword(key *ecdsa.PrivateKey) xmpp.GetCurrentPasswordFunc {
	f1 := func() string {
		pass, _ := signalshare.CreatePassword(key)
		return pass
	}
	return f1
}

type IceStream struct {
	ist            *ice.IceStreamTransport
	Status         IceStatus
	Role ice.SessionRole
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
this function should be protected by lock
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
				is.ist.Stop()
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
	log.Trace(fmt.Sprintf("%s send to %s , message=%s,hash=%s\n", it.name, utils.APex2(receiver), encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:]))))
	var err error
	if it.sendStatus != StatusCanSend {
		return errHasStopped
	}
	is, ok := it.IceStreamMap[receiver]
	if ok {
		if is.Status != IceTransporterStateNegotiateComplete {
			return errIceStreamTransporterNotReady
		}
		err = is.ist.SendData(data)
		return err
	} else { //start new p2p
		err := it.signal.TryReach(receiver)
		if err != nil {
			return err
		}
		is := &IceStream{
			Status: IceTransporterStateInit,
			Role:ice.SessionRoleControlling,
		}
		is.ist, err =  ice.NewIceStreamTransport(cfg,it.name)
		if err != nil {
			log.Trace(fmt.Sprintf("%s NewIceStreamTransport err %s", it.name, err))
			return err
		}
		ic:=&IceCallback{
			it:it,
			is:is,
			partner:receiver,
			datatosend:data,
		}
		is.ist.SetCallBack(ic)
		it.IceStreamMap[receiver] = is
		is.CanExchangeSdp = true
		it.startIce(is,receiver)
	}
	return nil
}

type sdpresult struct {
	sdp string
	err error
}
type IceCallback struct{
	it *IceTransport
	partner common.Address
	is *IceStream
	datatosend []byte
}
func (ic*IceCallback) OnReceiveData(data []byte,from net.Addr) {
	addr:=from.(*net.UDPAddr)
	ic.it.Receive(data,addr.IP.String(),addr.Port)
}
func (ic*IceCallback) OnIceComplete(result error) {
	if result != nil {
		log.Error(fmt.Sprintf("%s ice complete callback error err=%s", ic.it.name,result))
		ic.it.removeIceStreamTransport((ic.partner))
		ic.is.Status=IceTransporterStateClosed
	} else{
		ic.is.Status=IceTransporterStateNegotiateComplete
		if len(ic.datatosend)>0{
			ic.is.ist.SendData(ic.datatosend)
		}
	}
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
		Role:ice.SessionRoleControlled,
		CanExchangeSdp: true,
	}

	is.ist, err =ice.NewIceStreamTransport(cfg,it.name)
	if err != nil {
		return
	}
	ic:=&IceCallback{
		partner:partner,
		it:it,
		is:is,
	}
	is.ist.SetCallBack(ic)
	it.IceStreamMap[partner] = is
	sdpresult,err:=it.startIceWithSdp(is,sdp)
	log.Debug(fmt.Sprintf("%s get sdp:%s err:%s", it.name, sdpresult,err))
	return sdpresult,err

}
func (it *IceTransport) startIceWithSdp(is *IceStream,rsdp string) (sdpresult string, err error) {
	 err=is.ist.InitIce(ice.SessionRoleControlled)
	 if err!=nil{
	 	return
	 }
	 sdpresult,err=is.ist.EncodeSession()
	 if err!=nil{
	 	return
	 }
	 err=is.ist.StartNegotiation(rsdp)
	 return
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
	is.ist.Stop()
}
func (it *IceTransport) startIce(is *IceStream, receiver common.Address) {
	var err error
	defer func(){
		if err!=nil{
			it.removeIceStreamTransport(receiver)
		}
	}()
	err = is.ist.InitIce(ice.SessionRoleControlling)
		if err != nil {
			log.Error(fmt.Sprintf("%s %s InitIceSession err ", it.name, utils.APex(receiver), err))

			return
		}
		sdp, err := is.ist.EncodeSession()
		if err != nil {
			log.Error(fmt.Sprintf("%s get sdp error %s", it.name, err))
			return
		}

			partnersdp, err := it.signal.ExchangeSdp(receiver, sdp)
			if err != nil {
				log.Error(fmt.Sprintf("%s exchange sdp error %s", it.name, err))
				return
			}
			err = is.ist.StartNegotiation(partnersdp)
			if err != nil {
				log.Error(fmt.Sprintf("%s %s StartIce error %s", it.name, utils.APex(receiver), err))
			return
			}
}
func (it *IceTransport) Receive(data []byte, host string, port int) error {
	if it.receiveStatus == StatusStopReceive {
		return errStoppedReceive
	}
	log.Trace(fmt.Sprintf("%s receive message,message=%s,hash=%s\n", it.name, encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data))))
	if it.protocol != nil {
		log.Trace(fmt.Sprintf("%s message for protocol", it.name))
		go func() {
			// icestream  seems that the same thread is used for sending and receiving, so the reception must not be blocked. Otherwise, it will cause no transmission.
			it.protocol.Receive(data, host, port)
			log.Trace(fmt.Sprintf("%s message for protocol complete...", it.name))
		}()

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
		i.ist.Stop()
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
