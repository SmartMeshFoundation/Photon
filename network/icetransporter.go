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

type iceStatus int

const (
	iceTransporterStateInit = iceStatus(iota)
	iceTranspoterStateInitComplete
	iceTransporterStateNegotiateComplete
	iceTransporterStateClosed
)

const (
	statusCanReceive  = 0
	statusStopReceive = 1
	statusCanSend     = 0
	statusStopSend    = 1
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

//InitIceTransporter init ice configuration
func InitIceTransporter(turnServer, turnUser, turnPassowrd, signalServerURL string) {
	once.Do(func() {
		signalServer = signalServerURL
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
	ic   *iceCallback
}
type iceFail struct {
	receiver common.Address
	err      error
}

/*
IceTransport is a transport use ICE to setup P2P connection
it implements Transporter
there should only one connection between each pair nodes
*/
type IceTransport struct {
	key                   *ecdsa.PrivateKey
	addr                  common.Address
	address2IceStreamMaps map[common.Address]*iceCallback
	icestream2AddressMaps map[*iceCallback]common.Address
	lock                  sync.Mutex
	receiveStatus         int32
	sendStatus            int32
	signal                SignalInterface.SignalProxy
	protocol              ProtocolReceiver
	name                  string //for test
	connLastReceiveMap    map[common.Address]time.Time
	checkInterval         time.Duration
	stopChan              chan struct{}
	sendChan              chan *iceSend
	receiveChan           chan *iceReceive
	iceFailChan           chan *iceFail
	log                   log.Logger
}

//NewIceTransporter create IceTransport
func NewIceTransporter(key *ecdsa.PrivateKey, name string) (it *IceTransport, err error) {
	it = &IceTransport{
		key:                   key,
		receiveStatus:         statusStopReceive,
		address2IceStreamMaps: make(map[common.Address]*iceCallback),
		icestream2AddressMaps: make(map[*iceCallback]common.Address),
		connLastReceiveMap:    make(map[common.Address]time.Time),
		stopChan:              make(chan struct{}),
		sendChan:              make(chan *iceSend, 100),
		receiveChan:           make(chan *iceReceive, 100),
		iceFailChan:           make(chan *iceFail, 10),
		checkInterval:         time.Second * 180,
		addr:                  crypto.PubkeyToAddress(key.PublicKey),
		name:                  name,
		log:                   log.New("name", fmt.Sprintf("%s-IceTransport", name)),
	}
	sp, err := xmpp.NewSignalConnection(signalServer, it.addr, newpassword(it.key), func(from common.Address, sdp string) (mysdp string, err error) {
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

//RegisterProtocol register a receiver
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
			ic, ok := it.address2IceStreamMaps[r]
			if ok {
				it.log.Trace(fmt.Sprintf("%s connection has been removed", utils.APex(r)))
				delete(it.address2IceStreamMaps, r)
				delete(it.icestream2AddressMaps, ic)
				ic.ist.Stop()
			}
			it.lock.Unlock()
		}
	}
}

/*
Send one message:
1. check if has a connection,
2. if connection is invalid (maybe on setup), just return fail
3. if connection is valid ,just send
4. if no connection,  try to setup a p2p connection use ice.
*/
func (it *IceTransport) Send(receiver common.Address, host string, port int, data []byte) error {
	it.log.Trace(fmt.Sprintf("send to %s , message=%s,hash=%s\n", utils.APex2(receiver), encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:]))))
	if it.sendStatus != statusCanSend {
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
	ic, ok := it.address2IceStreamMaps[receiver]
	if ok {
		if ic.Status != iceTransporterStateNegotiateComplete {
			return errIceStreamTransporterNotReady
		}
		it.log.Trace(fmt.Sprintf("send to %s,msg=%s,data=\n%s", utils.APex2(receiver), encoding.MessageType(data[0]), hex.Dump(data)))
		err = ic.ist.SendData(data)
		return err
	}
	//start new p2p
	ic = &iceCallback{
		it:         it,
		partner:    receiver,
		datatosend: data,
		Status:     iceTransporterStateInit,
	}
	it.address2IceStreamMaps[receiver] = ic
	it.icestream2AddressMaps[ic] = receiver
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

type iceCallback struct {
	it         *IceTransport
	partner    common.Address
	datatosend []byte
	ist        *ice.StreamTransport
	Status     iceStatus
}

func (ic *iceCallback) OnReceiveData(data []byte, from net.Addr) {
	ic.it.log.Trace(fmt.Sprintf("icecallback receive data from %s, l=%d", from.String(), len(data)))
	if ic.it.receiveStatus == statusStopReceive {
		ic.it.log.Debug(fmt.Sprintf("receivie data from %s, but ice transport has stopped", from))
		return
	}
	ic.it.receiveChan <- &iceReceive{from, data, ic}
}
func (ic *iceCallback) OnIceComplete(result error) {
	ic.it.log.Trace(fmt.Sprintf("icecallback complete result=%v,partner=%s", result, utils.APex(ic.partner)))
	if result != nil {
		ic.it.log.Error(fmt.Sprintf("ice complete callback error err=%s", result))
		ic.it.removeIceStreamTransport((ic.partner))
		ic.Status = iceTransporterStateClosed
	} else {
		ic.Status = iceTransporterStateNegotiateComplete
		if len(ic.datatosend) > 0 {
			ic.it.sendChan <- &iceSend{ic.partner, ic.datatosend}
		}
	}
}
func (it *IceTransport) handleSdpArrived(partner common.Address, sdp string) (mysdp string, err error) {
	it.lock.Lock()
	defer it.lock.Unlock()
	if it.receiveStatus != statusCanReceive {
		err = errStoppedReceive
		return
	}
	it.log.Trace(fmt.Sprintf("handleSdpArrived from %s, sdp=%s", utils.APex2(partner), sdp))
	ic, ok := it.address2IceStreamMaps[partner]
	if ok {
		/*
			already have a connection, why remote create new connection, reasons:
			1. they are trying to send data each  other at the same time. The probability of this is very low
			2. partner thinks the connection is invalid,and drops this connection. but I think this connection is valid.
		*/
		err = fmt.Errorf("%s trying to send each other at the same time?", utils.APex(partner))
		it.log.Error(fmt.Sprintf("handleSdpArrived from %s,but my ice connection is ok", utils.APex2(partner)))
		it.lock.Unlock()
		//if partner think this connection is valid,this is invalid.
		it.removeIceStreamTransport(partner)
		it.lock.Lock()
	}
	ic = &iceCallback{
		partner: partner,
		it:      it,
		Status:  iceTransporterStateInit,
	}
	ic.ist, err = ice.NewIceStreamTransport(cfg, utils.APex2(partner))
	if err != nil {
		return
	}
	ic.ist.SetCallBack(ic)
	it.address2IceStreamMaps[partner] = ic
	it.icestream2AddressMaps[ic] = partner
	it.addCheck(partner)
	sdpresult, err := it.startIceWithSdp(ic, sdp)
	if err != nil {
		it.log.Error(fmt.Sprintf("startIceWithSdp:%s err:%s", utils.APex(partner), err))
	}
	return sdpresult, err

}
func (it *IceTransport) startIceWithSdp(ic *iceCallback, rsdp string) (sdpresult string, err error) {
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
	ic, ok := it.address2IceStreamMaps[receiver]
	if !ok {
		return
	}
	it.log.Info(fmt.Sprintf("removeIceStreamTransport %s", utils.APex2(receiver)))
	delete(it.address2IceStreamMaps, receiver)
	delete(it.icestream2AddressMaps, ic)
	if ic.ist != nil {
		ic.ist.Stop()
	}
}
func (it *IceTransport) startIce(ic *iceCallback, receiver common.Address) (err error) {
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
func (it *IceTransport) receiveData(ic *iceCallback, data []byte, host string, port int) error {

	it.lock.Lock()
	defer it.lock.Unlock()
	addr, ok := it.icestream2AddressMaps[ic]
	if !ok {
		it.log.Error("recevie data from unkown icestream, it must be a error")
		return nil
	}
	it.connLastReceiveMap[addr] = time.Now()
	return it.receive(data, host, port)
}
func (it *IceTransport) addCheck(addr common.Address) {
	it.connLastReceiveMap[addr] = time.Now()
}

//receive receive data and notify upper layer
func (it *IceTransport) receive(data []byte, host string, port int) error {
	it.log.Trace(fmt.Sprintf("receive message,message=%s,hash=%s\n", encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data))))
	if it.protocol != nil {
		it.log.Trace(fmt.Sprintf("message for protocol"))
		it.protocol.receive(data, host, port)
		it.log.Trace(fmt.Sprintf("message for protocol complete..."))

	}
	return nil
}

//Start transport,ready for send and receive
func (it *IceTransport) Start() {
	it.receiveStatus = statusCanReceive
}

//Stop send and receive
func (it *IceTransport) Stop() {
	it.StopAccepting()
	atomic.SwapInt32(&it.sendStatus, statusStopSend)
	close(it.stopChan)
	it.log.Trace("stopped")
	it.signal.Close()
	close(it.sendChan)
	//close(it.iceFailChan) //avoid crash, sendChan will make loop  quit.
	close(it.receiveChan)
	it.lock.Lock()
	for a, i := range it.address2IceStreamMaps {
		delete(it.address2IceStreamMaps, a)
		delete(it.icestream2AddressMaps, i)
		if i.ist != nil {
			i.ist.Stop()
		}
	}
	it.lock.Unlock()
}

//StopAccepting stops receive
func (it *IceTransport) StopAccepting() {
	atomic.SwapInt32(&it.receiveStatus, statusStopReceive)
}

//IceHelperDicovery a mock discovery for ICE,because ICE discovery nodes by signal server
type IceHelperDicovery struct {
}

//NewIceHelperDiscovery create IceHelperDicovery
func NewIceHelperDiscovery() *IceHelperDicovery {
	return new(IceHelperDicovery)
}

//Register just to implement Discover
func (id *IceHelperDicovery) Register(address common.Address, host string, port int) error {
	return nil
}

//Get just to implement Discover
func (id *IceHelperDicovery) Get(address common.Address) (host string, port int, err error) {
	return address.String(), 0, nil
}

//NodeIDByHostPort just to implement Discover
func (id *IceHelperDicovery) NodeIDByHostPort(host string, port int) (node common.Address, err error) {
	return common.HexToAddress(host), nil
}
