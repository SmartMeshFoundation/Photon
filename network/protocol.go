package network

import (
	"crypto/ecdsa"

	"encoding/hex"

	"reflect"

	"fmt"
	"time"

	"sync"

	"errors"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

const NODE_NETWORK_UNKNOWN = "unknown"
const NODE_NETWORK_UNREACHABLE = "unreachable"
const NODE_NETWORK_REACHABLE = "reachable"

type NodesStatusMap map[common.Address]*NetworkStatus

//const NetworkStatusUnkwn = 0
//const NetworkStatusUnReachable = 1
//const NetworkStatusReachable = 2

type NetworkStatus struct {
	LastTime time.Time //last update time
	Status   string
}
type MessageToRaiden struct {
	Msg      encoding.SignedMessager
	EchoHash common.Hash
}
type AsyncResult struct {
	Result chan error
	Sub    ethereum.Subscription
}

type SentMessageState struct {
	AsyncResult     *AsyncResult
	AckChannel      chan error
	ReceiverAddress common.Address
}
type NodesStatusGeter interface {
	GetNetworkStatus(addr common.Address) string
}

func NewAsyncResult() *AsyncResult {
	return &AsyncResult{Result: make(chan error, 1)}
}

type RaidenProtocol struct {
	Transport                    Transporter
	discovery                    DiscoveryInterface
	privKey                      *ecdsa.PrivateKey
	nodeAddr                     common.Address
	ReceivedHashesToAck          map[common.Hash]*encoding.Ack //also need a timeout to clear this map
	SentHashesToChannel          map[common.Hash]*SentMessageState
	retryTimes                   int
	retryInterval                time.Duration
	mapLock                      sync.Mutex
	address2NetworkStatus        map[common.Address]*NetworkStatus //todo need a lock .or to a new struct keep status mananger
	ReceivedMessageChannel       chan *MessageToRaiden
	ReceivedMessageResultChannel chan error
}

func NewRaidenProtocol(transport Transporter, discovery DiscoveryInterface, privKey *ecdsa.PrivateKey) *RaidenProtocol {
	rp := &RaidenProtocol{
		Transport:                    transport,
		discovery:                    discovery,
		privKey:                      privKey,
		retryTimes:                   10,
		retryInterval:                time.Millisecond * 6000,
		ReceivedHashesToAck:          make(map[common.Hash]*encoding.Ack),
		SentHashesToChannel:          make(map[common.Hash]*SentMessageState),
		address2NetworkStatus:        make(map[common.Address]*NetworkStatus),
		ReceivedMessageChannel:       make(chan *MessageToRaiden),
		ReceivedMessageResultChannel: make(chan error),
	}
	rp.nodeAddr = crypto.PubkeyToAddress(privKey.PublicKey)
	tr, ok := transport.(*UDPTransport)
	if ok {
		tr.Register(rp)
	}
	return rp
}

func New(sample interface{}) interface{} {
	t := reflect.ValueOf(sample)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	tt := t.Type()
	v := reflect.New(tt).Interface()
	return v
}
func (this *RaidenProtocol) _sendAck(host string, port int, data []byte) {
	reciver, err := this.discovery.NodeIdByHostPort(host, port)
	if err != nil {
		log.Error(fmt.Sprintf("unkonw %s:%d ,no such address", host, port))
	}
	this.Transport.Send(reciver, host, port, data)
}
func (this *RaidenProtocol) sendAck(receiver common.Address, ack *encoding.Ack) {
	this.sendRawWitNoAck(receiver, ack.Pack())
}
func (this *RaidenProtocol) sendRawWitNoAck(receiver common.Address, data []byte) error {
	host, port, err := this.discovery.Get(receiver)
	if err != nil {
		return err
	}
	return this.Transport.Send(receiver, host, port, data)
}

/*
msg should be signed
*/
func (this *RaidenProtocol) sendWithResult(receiver common.Address,
	msg encoding.Messager) (result *AsyncResult, err error) {
	/*
		1. if this packet is on sending, retry send immediately
		2. retry send many times, until receive a ack
		3. after many tries, the receiver maybe offline,so tell the caller

		caller can read from chan reusltChannel to get if this packet is successfully sent to receiver
		the caller can unsubscribe this event to cancel to wait
	*/
	if true {
		signed, ok := msg.(encoding.SignedMessager)
		if ok && len(signed.GetSignature()) <= 0 {
			log.Error("send unsigned message")
			return nil, errors.New("send unsigned message")
		}
	}
	data := msg.Pack()
	echohash := utils.Sha3(data, receiver[:])
	this.mapLock.Lock()
	msgState, ok := this.SentHashesToChannel[echohash]
	this.mapLock.Unlock()
	if ok {
		result = msgState.AsyncResult
		return
	}
	msgState = &SentMessageState{
		AsyncResult:     NewAsyncResult(),
		ReceiverAddress: receiver,
		AckChannel:      make(chan error, 1),
	}
	stopSendChannel := make(chan struct{})
	sub := event.NewSubscription(func(quit <-chan struct{}) error {
		select {
		case _, ok := <-msgState.AckChannel: //tell caller ack has arrived
			if ok {
				msgState.AsyncResult.Result <- nil
			} else {
				msgState.AsyncResult.Result <- errors.New("channel closed")
			}
			stopSendChannel <- struct{}{}
			close(stopSendChannel)
			this.mapLock.Lock()
			close(msgState.AsyncResult.Result)
			this.mapLock.Unlock()
		case <-quit: //the caller cancel to wait this result
			stopSendChannel <- struct{}{}
			close(stopSendChannel)
			this.mapLock.Lock()
			delete(this.SentHashesToChannel, echohash)
			close(msgState.AckChannel)
			close(msgState.AsyncResult.Result)
			this.mapLock.Unlock()
			return nil
		}
		return nil
	})
	msgState.AsyncResult.Sub = sub
	this.mapLock.Lock()
	this.SentHashesToChannel[echohash] = msgState
	this.mapLock.Unlock()
	result = msgState.AsyncResult
	go func() { //try to send many times
		for i := 0; i < this.retryTimes; i++ {
			//log.Trace(fmt.Sprintf("before send:\n%s\n", hex.Dump(data)))
			err := this.sendRawWitNoAck(receiver, data)
			if err != nil {
				log.Info(err.Error())
			}
			timeout := time.After(this.retryInterval)
			select {
			case <-stopSendChannel:
				//receive ack or user canceld
				return
			case <-timeout: //retry

			}
		}
	}()
	return
}

// send this packet and wait ack until timeout
func (this *RaidenProtocol) SendAndWait(receiver common.Address, msg encoding.Messager, timeout time.Duration) error {
	result, err := this.sendWithResult(receiver, msg)
	if err != nil {
		return err
	}
	timeoutCh := time.After(timeout)
	select {
	case err = <-result.Result:
		if err == nil {
			this.updateNetworkStatus(receiver, NODE_NETWORK_REACHABLE)
		}
	case <-timeoutCh:
		err = errors.New("time out of sendWithResult")
		this.updateNetworkStatus(receiver, NODE_NETWORK_UNREACHABLE)
	}
	return err
}
func (this *RaidenProtocol) SendAsync(receiver common.Address, msg encoding.Messager) (*AsyncResult, error) {
	return this.sendWithResult(receiver, msg)
}
func (this *RaidenProtocol) createAck(echohash common.Hash) *encoding.Ack {
	return encoding.NewAck(this.nodeAddr, echohash)
}
func (this *RaidenProtocol) updateNetworkStatus(addr common.Address, status string) {
	s, ok := this.address2NetworkStatus[addr]
	if !ok {
		s = &NetworkStatus{
			time.Now(), NODE_NETWORK_UNKNOWN,
		}
		this.address2NetworkStatus[addr] = s
	}
	s.Status = status
	s.LastTime = time.Now()
}
func (this *RaidenProtocol) GetNetworkStatus(addr common.Address) string {
	s, ok := this.address2NetworkStatus[addr]
	if !ok {
		return NODE_NETWORK_UNKNOWN
	}
	return s.Status
}
func (this *RaidenProtocol) Receive(data []byte, host string, port int) {
	if len(data) > params.UDP_MAX_MESSAGE_SIZE {
		log.Error("receive packet larger than maximum size :", len(data))
		return
	}
	cmdid := int(data[0])
	echohash := utils.Sha3(data, this.nodeAddr[:])
	this.mapLock.Lock()
	ack, ok := this.ReceivedHashesToAck[echohash]
	this.mapLock.Unlock()
	if ok { //avoid to notify new message.
		this._sendAck(host, port, ack.Pack())
		return
	}
	messager, ok := encoding.MessageMap[cmdid]
	if !ok {
		log.Warn("receive unknown message:", hex.Dump(data))
		return
	}
	msg := New(messager)
	messager = msg.(encoding.Messager)
	err := messager.UnPack(data)
	if err != nil {
		log.Warn("message unpack error : ", err)
		return
	}
	if messager.Cmd() == encoding.ACK_CMDID { //some one may be waiting this ack
		ackMsg := messager.(*encoding.Ack)
		this.updateNetworkStatus(ackMsg.Sender, NODE_NETWORK_REACHABLE)
		this.mapLock.Lock()
		msgState, ok := this.SentHashesToChannel[ackMsg.Echo]
		if ok {
			msgState.AckChannel <- nil
			close(msgState.AckChannel)
			delete(this.SentHashesToChannel, ackMsg.Echo)
		} else {
			log.Debug(fmt.Sprintf("receive duplicate ack  from %s:%d ", host, port))
		}
		this.mapLock.Unlock()
	} else {
		signedMessager, ok := messager.(encoding.SignedMessager)
		if !ok {
			log.Warn("message should be signed except for ack")
		}
		err := signedMessager.VerifySignature(data)
		if err != nil {
			log.Warn(fmt.Sprint("verify message  signature error,length:%d, from %s:%d ", len(data), host, port))
			return
		}
		this.updateNetworkStatus(signedMessager.GetSender(), NODE_NETWORK_REACHABLE)
		if messager.Cmd() == encoding.PING_CMDID { //send ack
			this.sendAck(signedMessager.GetSender(), this.createAck(echohash))
		} else {
			//send message to raiden ,and wait result
			this.ReceivedMessageChannel <- &MessageToRaiden{signedMessager, echohash}
			err = <-this.ReceivedMessageResultChannel
			//only send the Ack if the message was handled without exceptions
			if err == nil {
				ack := this.createAck(echohash)
				this.sendAck(signedMessager.GetSender(), ack)
				this.mapLock.Lock()
				this.ReceivedHashesToAck[echohash] = ack
				this.mapLock.Unlock()
			}
		}
	}

}

func (this *RaidenProtocol) StopAndWait() {
	this.Transport.StopAccepting()
	this.mapLock.Lock()
	for k, c := range this.SentHashesToChannel {
		delete(this.SentHashesToChannel, k)
		close(c.AckChannel)
		close(c.AsyncResult.Result)
	}
	this.mapLock.Unlock()
	//what about the outgoing packets, maybe lost
	this.Transport.Stop()
	close(this.ReceivedMessageResultChannel)
	close(this.ReceivedMessageChannel)

}

func (this *RaidenProtocol) Start() {
	this.Transport.Start()
}
