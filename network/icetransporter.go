package network

import (
	"crypto/ecdsa"

	"fmt"

	"sync/atomic"

	"sync"

	"errors"

	"time"

	"net"

	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/ice"
	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/interface"
	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/signalshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/xmpp"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	turnServer   = ""
	turnUser     = ""
	turnPassword = ""
	cfg          *ice.TransportConfig
)
var (
	errHasStopped                   = errors.New("ice transporter has stopped")
	errStoppedReceive               = errors.New("ice transporter has stopped receiving")
	errIceStreamTransporterNotReady = errors.New("icestreamtransport not ready to send")
)
var once sync.Once

func InitIceTransporter(turnServer, turnUser, turnPassowrd, signalServerUrl string) {
	once.Do(func() {
		signalServer = signalServerUrl
		cfg = ice.NewTransportConfigWithTurn(turnServer, turnUser, turnPassowrd)
	})
}
func newpassword(key *ecdsa.PrivateKey) xmpp.GetCurrentPasswordFunc {
	f1 := func() string {
		pass, _ := signalshare.CreatePassword(key)
		return pass
	}
	return f1
}

/*
send data to receiver.
*/
type iceSend struct {
	receiver common.Address
	data     []byte
}

/*
receive data from `from`, and related iceStreamTranspoter is stored in `ic`
*/
type iceReceive struct {
	from net.Addr
	data []byte
	ic   *IceCallback
}
type iceFail struct {
	receiver common.Address
	err      error
}
type IceTransport struct {
	key                  *ecdsa.PrivateKey
	Addr                 common.Address
	Address2IceStreamMap map[common.Address]*IceCallback
	Icestream2AddressMap map[*IceCallback]common.Address
	lock                 sync.Mutex
	receiveStatus        int32
	sendStatus           int32
	signal               SignalInterface.SignalProxy
	protocol             ProtocolReceiver
	name                 string //for test
	connLastReceiveMap   map[common.Address]time.Time
	checkInterval        time.Duration
	stopChan             chan struct{}
	sendChan             chan *iceSend
	receiveChan          chan *iceReceive
	iceFailChan          chan *iceFail
	log                  log.Logger
}

func NewIceTransporter(key *ecdsa.PrivateKey, name string) (it *IceTransport, err error) {
	it = &IceTransport{
		key:                  key,
		receiveStatus:        StatusStopReceive,
		Address2IceStreamMap: make(map[common.Address]*IceCallback),
		Icestream2AddressMap: make(map[*IceCallback]common.Address),
		connLastReceiveMap:   make(map[common.Address]time.Time),
		stopChan:             make(chan struct{}),
		sendChan:             make(chan *iceSend, 100),
		receiveChan:          make(chan *iceReceive, 100),
		iceFailChan:          make(chan *iceFail, 10),
		checkInterval:        time.Second * 180,
		Addr:                 crypto.PubkeyToAddress(key.PublicKey),
		name:                 name,
		log:                  log.New("name", fmt.Sprintf("%s-IceTransport", name)),
	}
	sp, err := xmpp.NewXmpp(signalServer, it.Addr, newpassword(it.key), func(from common.Address, sdp string) (mysdp string, err error) {
		return it.handleSdpArrived(from, sdp)
	}, name)
	if err != nil {
		err = fmt.Errorf("create ice transpoter error %s", err)
		return
	}
	it.signal = sp
	go it.loop()
	return
}
func (it *IceTransport) RegisterProtocol(protcol ProtocolReceiver) {
	it.protocol = protcol
}
func (it *IceTransport) loop() {
	var ok bool
	var s *iceSend
	var r *iceReceive
	var f *iceFail
	var err error
	defer func() {
		it.log.Info(fmt.Sprintf("IceTransport quit loop"))
	}()
	for {
		id := utils.RandomString(10)
		it.log.Trace(fmt.Sprintf("IceTransport loop %s start", id))
		select {
		case s, ok = <-it.sendChan:
			if !ok {
				return
			}
			it.log.Trace(fmt.Sprintf("start send to %s, l=%d", utils.APex(s.receiver), len(s.data)))
			err = it.sendInternal(s.receiver, s.data)
			if err != nil {
				it.log.Info(fmt.Sprintf("send to %s, error:%s", utils.APex(s.receiver), err))
			}
		case r, ok = <-it.receiveChan:
			if !ok {
				return
			}
			it.log.Trace(fmt.Sprintf("received data from %s,addr=%s", r.from, utils.APex(r.ic.partner)))
			addr := r.from.(*net.UDPAddr)
			it.receiveData(r.ic, r.data, addr.IP.String(), addr.Port)
		case f, ok = <-it.iceFailChan:
			if !ok {
				return
			}
			it.log.Trace(fmt.Sprintf("ice %s failed,because of %v", utils.APex(f.receiver), f.err))
			it.removeIceStreamTransport(f.receiver)
		case <-time.After(it.checkInterval):
			if len(it.connLastReceiveMap) > 0 {
				it.removeExpiredConnection()
			}
		case <-it.stopChan:
			return
		}
		it.log.Trace(fmt.Sprintf("IceTransport  loop %s end", id))

	}
}

/*
for connections that don't use for a long time, just to remove.
for connections in use but may be invalid because of network, remove too.
this function should be protected by lock
*/
func (it *IceTransport) removeExpiredConnection() {
	now := time.Now()
	for r, t := range it.connLastReceiveMap {
		if now.Sub(t) > 2*it.checkInterval {
			it.lock.Lock()
			ic, ok := it.Address2IceStreamMap[r]
			if ok {
				it.log.Trace(fmt.Sprintf("%s connection has been removed", utils.APex(r)))
				delete(it.Address2IceStreamMap, r)
				delete(it.Icestream2AddressMap, ic)
				ic.ist.Stop()
			}
			it.lock.Unlock()
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
	it.log.Trace(fmt.Sprintf("send to %s , message=%s,hash=%s\n", utils.APex2(receiver), encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:]))))
	if it.sendStatus != StatusCanSend {
		it.log.Info(fmt.Sprintf("send data to %s, but icetransport has been stopped", utils.APex(receiver)))
		return errHasStopped
	}
	it.sendChan <- &iceSend{receiver, data}
	return nil
}
func (it *IceTransport) sendInternal(receiver common.Address, data []byte) error {
	var err error
	it.lock.Lock()
	defer it.lock.Unlock()
	ic, ok := it.Address2IceStreamMap[receiver]
	if ok {
		if ic.Status != IceTransporterStateNegotiateComplete {
			return errIceStreamTransporterNotReady
		}
		it.log.Trace(fmt.Sprintf("send to %s,msg=%s,data=\n%s", utils.APex2(receiver), encoding.MessageType(data[0]), hex.Dump(data)))
		err = ic.ist.SendData(data)
		return err
	}
	//start new p2p
	ic = &IceCallback{
		it:         it,
		partner:    receiver,
		datatosend: data,
		Status:     IceTransporterStateInit,
	}
	it.Address2IceStreamMap[receiver] = ic
	it.Icestream2AddressMap[ic] = receiver
	go func() {
		/*
			其他节点之间的 ice, 不能影响已经协商完毕的连接.
		*/
		err := it.signal.TryReach(receiver)
		if err != nil {
			it.iceFailChan <- &iceFail{receiver, err}
			return
		}
		ic.ist, err = ice.NewIceStreamTransport(cfg, utils.APex2(receiver))
		if err != nil {
			it.log.Trace(fmt.Sprintf("NewIceStreamTransport err %s", err))
			it.iceFailChan <- &iceFail{receiver, err}
			return
		}
		ic.ist.SetCallBack(ic)
		it.addCheck(receiver)
		err = it.startIce(ic, receiver)
		if err != nil {
			it.iceFailChan <- &iceFail{receiver, err}
			return
		}
	}()
	return nil
}

type IceCallback struct {
	it         *IceTransport
	partner    common.Address
	datatosend []byte
	ist        *ice.IceStreamTransport
	Status     IceStatus
}

func (ic *IceCallback) OnReceiveData(data []byte, from net.Addr) {
	ic.it.log.Trace(fmt.Sprintf("icecallback receive data from %s, l=%d", from.String(), len(data)))
	if ic.it.receiveStatus == StatusStopReceive {
		ic.it.log.Debug(fmt.Sprintf("receivie data from %s, but ice transport has stopped", from))
		return
	}
	ic.it.receiveChan <- &iceReceive{from, data, ic}
}
func (ic *IceCallback) OnIceComplete(result error) {
	ic.it.log.Trace(fmt.Sprintf("icecallback complete result=%v,partner=%s", result, utils.APex(ic.partner)))
	if result != nil {
		ic.it.log.Error(fmt.Sprintf("ice complete callback error err=%s", result))
		ic.it.removeIceStreamTransport((ic.partner))
		ic.Status = IceTransporterStateClosed
	} else {
		ic.Status = IceTransporterStateNegotiateComplete
		if len(ic.datatosend) > 0 {
			ic.it.sendChan <- &iceSend{ic.partner, ic.datatosend}
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
	it.log.Trace(fmt.Sprintf("handleSdpArrived from %s, sdp=%s", utils.APex2(partner), sdp))
	ic, ok := it.Address2IceStreamMap[partner]
	if ok {
		/*
			already have a connection, why remote create new connection, reasons:
			1. they are trying to send data each  other at the same time. The probability of this is very low
			2. partner thinks the connection is invalid,and drops this connection. but I think this connection is valid.
		*/
		err = errors.New(fmt.Sprintf("%s trying to send each other at the same time?", utils.APex(partner)))
		it.log.Error(fmt.Sprintf("handleSdpArrived from %s,but my ice connection is ok", utils.APex2(partner)))
		it.lock.Unlock()
		//if partner think this connection is valid,this is invalid.
		it.removeIceStreamTransport(partner)
		it.lock.Lock()
	}
	ic = &IceCallback{
		partner: partner,
		it:      it,
		Status:  IceTransporterStateInit,
	}
	ic.ist, err = ice.NewIceStreamTransport(cfg, utils.APex2(partner))
	if err != nil {
		return
	}
	ic.ist.SetCallBack(ic)
	it.Address2IceStreamMap[partner] = ic
	it.Icestream2AddressMap[ic] = partner
	it.addCheck(partner)
	sdpresult, err := it.startIceWithSdp(ic, sdp)
	if err != nil {
		it.log.Error(fmt.Sprintf("startIceWithSdp:%s err:%s", utils.APex(partner), err))
	}
	return sdpresult, err

}
func (it *IceTransport) startIceWithSdp(ic *IceCallback, rsdp string) (sdpresult string, err error) {
	err = ic.ist.InitIce(ice.SessionRoleControlled)
	if err != nil {
		it.log.Trace(fmt.Sprintf("startIceWithSdp init ice err %s with %s", err, ic.ist.Name))
		return
	}
	sdpresult, err = ic.ist.EncodeSession()
	if err != nil {
		it.log.Trace(fmt.Sprintf("%s EncodeSession err %s", ic.ist.Name, err))
		return
	}
	go ic.ist.StartNegotiation(rsdp)
	return
}
func (it *IceTransport) removeIceStreamTransport(receiver common.Address) {
	it.lock.Lock()
	defer it.lock.Unlock()
	ic, ok := it.Address2IceStreamMap[receiver]
	if !ok {
		return
	}
	it.log.Info(fmt.Sprintf("removeIceStreamTransport %s", utils.APex2(receiver)))
	delete(it.Address2IceStreamMap, receiver)
	delete(it.Icestream2AddressMap, ic)
	if ic.ist != nil {
		ic.ist.Stop()
	}
}
func (it *IceTransport) startIce(ic *IceCallback, receiver common.Address) (err error) {
	err = ic.ist.InitIce(ice.SessionRoleControlling)
	if err != nil {
		it.log.Error(fmt.Sprintf("%s %s InitIceSession err ", utils.APex(receiver), err))
		return
	}
	sdp, err := ic.ist.EncodeSession()
	if err != nil {
		it.log.Error(fmt.Sprintf("get sdp error %s for %s", err, utils.APex(receiver)))
		return
	}
	partnersdp, err := it.signal.ExchangeSdp(receiver, sdp)
	if err != nil {
		it.log.Error(fmt.Sprintf("exchange sdp error %s for %s", err, utils.APex(receiver)))
		return
	}
	err = ic.ist.StartNegotiation(partnersdp)
	if err != nil {
		it.log.Error(fmt.Sprintf("%s StartIce error %s", utils.APex(receiver), err))
		return
	}
	return nil
}
func (it *IceTransport) receiveData(ic *IceCallback, data []byte, host string, port int) error {

	it.lock.Lock()
	defer it.lock.Unlock()
	addr, ok := it.Icestream2AddressMap[ic]
	if !ok {
		it.log.Error("recevie data from unkown icestream, it must be a error")
		return nil
	}
	it.connLastReceiveMap[addr] = time.Now()
	return it.Receive(data, host, port)
}
func (it *IceTransport) addCheck(addr common.Address) {
	it.connLastReceiveMap[addr] = time.Now()
}
func (it *IceTransport) Receive(data []byte, host string, port int) error {
	it.log.Trace(fmt.Sprintf("receive message,message=%s,hash=%s\n", encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data))))
	if it.protocol != nil {
		it.log.Trace(fmt.Sprintf("message for protocol"))
		it.protocol.Receive(data, host, port)
		it.log.Trace(fmt.Sprintf("message for protocol complete..."))

	}
	return nil
}
func (it *IceTransport) Start() {
	it.receiveStatus = StatusCanReceive
}
func (it *IceTransport) Stop() {
	it.StopAccepting()
	atomic.SwapInt32(&it.sendStatus, StatusStopSend)
	close(it.stopChan)
	it.log.Trace("stopped")
	it.signal.Close()
	close(it.sendChan)
	//close(it.iceFailChan) //avoid crash, sendChan will make loop  quit.
	close(it.receiveChan)
	it.lock.Lock()
	for a, i := range it.Address2IceStreamMap {
		delete(it.Address2IceStreamMap, a)
		delete(it.Icestream2AddressMap, i)
		if i.ist != nil {
			i.ist.Stop()
		}
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
