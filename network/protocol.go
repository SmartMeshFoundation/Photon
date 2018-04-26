package network

import (
	"crypto/ecdsa"

	"encoding/hex"

	"reflect"

	"fmt"
	"time"

	"sync"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const NODE_NETWORK_UNKNOWN = "unknown"
const NODE_NETWORK_UNREACHABLE = "unreachable"
const NODE_NETWORK_REACHABLE = "reachable"

var errTimeout = errors.New("wait timeout")
var errExpired = errors.New("message expired")

type NodesStatusMap map[common.Address]*NetworkStatus

type NetworkStatus struct {
	LastAckTime time.Time //time of last received ack message
	Status      string
}
type MessageToRaiden struct {
	Msg      encoding.SignedMessager
	EchoHash common.Hash
}
type AsyncResult struct {
	Result chan error
	Tag    interface{}
}

type SentMessageState struct {
	AsyncResult     *AsyncResult
	AckChannel      chan error
	ReceiverAddress common.Address
	Success         bool

	Message  encoding.Messager //message to send
	EchoHash common.Hash       //message echo hash
	Data     []byte            //packed message
}
type NodesStatusGetter interface {
	GetNetworkStatus(addr common.Address) string
	GetNetworkStatusAndLastAckTime(addr common.Address) (status string, lastAckTime time.Time)
}
type PingSender interface {
	SendPing(receiver common.Address) error
}

/*
get the lastest block number,so sender can remove expired mediated transfer.
for example :
A send B a mediated transfer, but B is offline
when B is online ,this transfer is invalid, so A will never receive  ack ,so A will try forever.
message secret,secretRequest,revealSecret won't allow error
*/
type BlockNumberGetter interface {
	GetBlockNumber() int64
}

func NewAsyncResult() *AsyncResult {
	return &AsyncResult{Result: make(chan error, 1)}
}

type timeoutGenerator func() time.Duration

/*
Timeouts generator with an exponential backoff strategy.

    Timeouts start spaced by `timeout`, after `retries` exponentially increase
    the retry delays until `maximum`, then maximum is returned indefinitely.
*/
func timeoutExponentialBackoff(retries int, timeout, maximumTimeout time.Duration) timeoutGenerator {
	tries := 1
	return func() time.Duration {
		tries += 1
		if tries < retries {
			return timeout
		}
		if timeout < maximumTimeout {
			timeout = timeout * 2
			if timeout < maximumTimeout {
				return timeout
			}
		}
		return maximumTimeout
	}
}

type RaidenProtocol struct {
	Transport                    Transporter
	discovery                    DiscoveryInterface
	privKey                      *ecdsa.PrivateKey
	nodeAddr                     common.Address
	SentHashesToChannel          map[common.Hash]*SentMessageState
	retryTimes                   int
	retryInterval                time.Duration
	mapLock                      sync.Mutex
	address2NetworkStatus        map[common.Address]*NetworkStatus
	statusLock                   sync.RWMutex
	ReceivedMessageChannel       chan *MessageToRaiden
	ReceivedMessageResultChannel chan error
	sendingQueueMap              map[string]chan *SentMessageState //write to this channel to send a message
	quitWaitGroup                sync.WaitGroup                    //wait before quit
	receivedMessageSaver         ReceivedMessageSaver
	BlockNumberGetter            BlockNumberGetter
	onStop                       bool //flag for stop
}

func NewRaidenProtocol(transport Transporter, discovery DiscoveryInterface, privKey *ecdsa.PrivateKey, blockNumberGetter BlockNumberGetter) *RaidenProtocol {
	rp := &RaidenProtocol{
		Transport:                    transport,
		discovery:                    discovery,
		privKey:                      privKey,
		retryTimes:                   10,
		retryInterval:                time.Millisecond * 6000,
		SentHashesToChannel:          make(map[common.Hash]*SentMessageState),
		address2NetworkStatus:        make(map[common.Address]*NetworkStatus),
		ReceivedMessageChannel:       make(chan *MessageToRaiden),
		ReceivedMessageResultChannel: make(chan error),
		sendingQueueMap:              make(map[string]chan *SentMessageState),
		BlockNumberGetter:            blockNumberGetter,
	}
	rp.nodeAddr = crypto.PubkeyToAddress(privKey.PublicKey)
	tr, ok := transport.(*UDPTransport)
	if ok {
		tr.Register(rp)
	}
	tr2, ok := transport.(*IceTransport)
	if ok {
		tr2.Register(rp)
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
func (this *RaidenProtocol) SetReceivedMessageSaver(saver ReceivedMessageSaver) {
	this.receivedMessageSaver = saver
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
func (this *RaidenProtocol) SendPing(receiver common.Address) error {
	ping := encoding.NewPing(utils.NewRandomInt64())
	ping.Sign(this.privKey, ping)
	data := ping.Pack()
	return this.sendRawWitNoAck(receiver, data)
}

/*
message mediatedTransfer and refundTransfer can safely be discarded when expired.
*/
func (this *RaidenProtocol) messageCanBeSent(msg encoding.Messager) bool {
	var expired int64 = 0
	switch msg2 := msg.(type) {
	case *encoding.MediatedTransfer:
		expired = msg2.Expiration
	case *encoding.RefundTransfer:
		expired = msg2.Expiration
	}
	if expired > 0 && expired <= this.BlockNumberGetter.GetBlockNumber() {
		return false
	}
	return true
}
func (this *RaidenProtocol) getChannelQueue(receiver, token common.Address) chan<- *SentMessageState {

	this.mapLock.Lock()
	defer this.mapLock.Unlock()
	key := fmt.Sprintf("%s-%s", receiver.String(), token.String())
	var sendingChan chan *SentMessageState
	var ok bool
	if token == utils.EmptyAddress { //no token means that this message doesn't need ordered.
		sendingChan = make(chan *SentMessageState, 1) //should not block sender
	} else {
		sendingChan, ok = this.sendingQueueMap[key]
		if ok {
			return sendingChan
		}
		sendingChan = make(chan *SentMessageState, 1000) //should not block sender
		this.sendingQueueMap[key] = sendingChan
	}
	go func() {
		/*
			1. if this packet is on sending, retry send immediately
			2. retry infinite, until receive a ack
			3. this message should be sent by caller after restart.

			caller can read from chan reusltChannel to get if this packet is successfully sent to receiver
		*/
	labelNextMessage:
		for {
			log.Trace(fmt.Sprintf("queue %s try send next message", key))
			this.quitWaitGroup.Add(1)
			msgState, ok := <-sendingChan
			if !ok {
				log.Info(fmt.Sprintf("queue %s quit, because of chan closed", key))
				this.quitWaitGroup.Done() //user stop
				return
			}
			for {
				if !this.messageCanBeSent(msgState.Message) {
					log.Info(fmt.Sprintf("message cannot be send because of expired msg=%s", msgState.Message))
					msgState.AsyncResult.Result <- errExpired
					close(msgState.AsyncResult.Result)
					this.quitWaitGroup.Done()
					break
				}
				nextTimeout := timeoutExponentialBackoff(this.retryTimes, this.retryInterval, this.retryInterval*10)
				err := this.sendRawWitNoAck(receiver, msgState.Data)
				if err != nil {
					log.Info(fmt.Sprintf("sendRawWitNoAck %s msg error %s", key, err.Error()))
				}
				timeout := time.After(nextTimeout())
				select {
				case _, ok = <-msgState.AckChannel:
					if ok {
						log.Trace(fmt.Sprintf("msg=%s, sent success :%s", encoding.MessageType(msgState.Message.Cmd()), utils.HPex(msgState.EchoHash)))
						msgState.AsyncResult.Result <- nil
						close(msgState.AsyncResult.Result)
						this.quitWaitGroup.Done()
						goto labelNextMessage
					} else {
						//message must send success, otherwise keep trying...
						log.Info(fmt.Sprintf("queue %s quit, because of chan closed", key))
						this.quitWaitGroup.Done()
						return //user call stop
					}
				case <-timeout: //retry

				}
			}
		}

	}()
	return sendingChan
}
func getMessageTokenAddress(msg encoding.Messager) common.Address {
	var tokenAddress common.Address
	switch msg2 := msg.(type) {
	case *encoding.DirectTransfer:
		tokenAddress = msg2.Token
	case *encoding.MediatedTransfer:
		tokenAddress = msg2.Token
	case *encoding.RefundTransfer:
		tokenAddress = msg2.Token
	}
	return tokenAddress
}
func getMessageChannelAddress(msg encoding.Messager) common.Address {
	var channelAddress common.Address
	switch msg2 := msg.(type) {
	case *encoding.DirectTransfer:
		channelAddress = msg2.Channel
	case *encoding.MediatedTransfer:
		channelAddress = msg2.Channel
	case *encoding.RefundTransfer:
		channelAddress = msg2.Channel
	case *encoding.Secret:
		channelAddress = msg2.Channel
	}
	return channelAddress
}

/*
msg should be signed.
msg must be sent success.
*/
func (this *RaidenProtocol) sendWithResult(receiver common.Address,
	msg encoding.Messager) (result *AsyncResult) {
	//no more message...
	if this.onStop {
		return NewAsyncResult()
	}
	this.quitWaitGroup.Add(1)
	defer this.quitWaitGroup.Done()
	if true {
		signed, ok := msg.(encoding.SignedMessager)
		if ok && len(signed.GetSignature()) <= 0 {
			log.Error("send unsigned message")
			panic("send unsigned message")
			return
		}
	}
	data := msg.Pack()
	echohash := utils.Sha3(data, receiver[:])
	this.mapLock.Lock()
	msgState, ok := this.SentHashesToChannel[echohash]
	if ok {
		this.mapLock.Unlock()
		result = msgState.AsyncResult
		return
	}
	log.Debug(fmt.Sprintf("send msg=%s to=%s,expected hash=%s", encoding.MessageType(msg.Cmd()), utils.APex2(receiver), utils.HPex(echohash)))
	msgState = &SentMessageState{
		AsyncResult:     NewAsyncResult(),
		ReceiverAddress: receiver,
		AckChannel:      make(chan error, 1),
		Message:         msg,
		Data:            data,
		EchoHash:        echohash,
	}
	this.SentHashesToChannel[echohash] = msgState
	this.mapLock.Unlock()
	result = msgState.AsyncResult
	tokenAddress := getMessageTokenAddress(msg)
	//make sure not block
	this.getChannelQueue(receiver, tokenAddress) <- msgState
	return
}

// send this packet and wait ack until timeout
func (this *RaidenProtocol) SendAndWait(receiver common.Address, msg encoding.Messager, timeout time.Duration) error {
	var err error
	result := this.sendWithResult(receiver, msg)
	timeoutCh := time.After(timeout)
	select {
	case err = <-result.Result:
		if err == nil {
			this.updateNetworkStatus(receiver, NODE_NETWORK_REACHABLE)
		}
	case <-timeoutCh:
		err = errTimeout
		this.updateNetworkStatus(receiver, NODE_NETWORK_UNREACHABLE)
	}
	return err
}
func (this *RaidenProtocol) SendAsync(receiver common.Address, msg encoding.Messager) *AsyncResult {
	return this.sendWithResult(receiver, msg)
}
func (this *RaidenProtocol) CreateAck(echohash common.Hash) *encoding.Ack {
	return encoding.NewAck(this.nodeAddr, echohash)
}
func (this *RaidenProtocol) updateNetworkStatus(addr common.Address, status string) {
	this.statusLock.Lock()
	defer this.statusLock.Unlock()
	s, ok := this.address2NetworkStatus[addr]
	if !ok {
		s = &NetworkStatus{
			time.Now(), NODE_NETWORK_UNKNOWN,
		}
		this.address2NetworkStatus[addr] = s
	}
	s.Status = status
	s.LastAckTime = time.Now()
}
func (this *RaidenProtocol) GetNetworkStatus(addr common.Address) string {
	this.statusLock.Lock()
	defer this.statusLock.Unlock()
	s, ok := this.address2NetworkStatus[addr]
	if !ok {
		return NODE_NETWORK_UNKNOWN
	}
	return s.Status
}
func (this *RaidenProtocol) GetNetworkStatusAndLastAckTime(addr common.Address) (status string, lastAckTime time.Time) {
	this.statusLock.Lock()
	defer this.statusLock.Unlock()
	s, ok := this.address2NetworkStatus[addr]
	if !ok {
		return NODE_NETWORK_UNKNOWN, time.Now()
	}
	return s.Status, s.LastAckTime
}
func (this *RaidenProtocol) Receive(data []byte, host string, port int) {
	if len(data) > params.UDP_MAX_MESSAGE_SIZE {
		log.Error("receive packet larger than maximum size :", len(data))
		return
	}
	//ignore incomming message when stop
	if this.onStop {
		return
	}
	//wait finish this packet when stop
	this.quitWaitGroup.Add(1)
	defer this.quitWaitGroup.Done()
	cmdid := int(data[0])
	echohash := utils.Sha3(data, this.nodeAddr[:])
	if this.receivedMessageSaver != nil {
		ackdata := this.receivedMessageSaver.GetAck(echohash)
		if len(ackdata) > 0 {
			this._sendAck(host, port, ackdata)
			return
		}
	}
	messager, ok := encoding.MessageMap[cmdid]
	if !ok {
		log.Warn("receive unknown message:", hex.Dump(data))
		return
	}
	messager = New(messager).(encoding.Messager)
	err := messager.UnPack(data)
	if err != nil {
		log.Warn(fmt.Sprintf("message unpack error : %s", err))
		return
	}
	if messager.Cmd() == encoding.ACK_CMDID { //some one may be waiting this ack
		ackMsg := messager.(*encoding.Ack)
		log.Debug(fmt.Sprintf("receive ack ,hash=%s", utils.HPex(ackMsg.Echo)))
		this.updateNetworkStatus(ackMsg.Sender, NODE_NETWORK_REACHABLE)
		this.mapLock.Lock()
		msgState, ok := this.SentHashesToChannel[ackMsg.Echo]
		if ok && msgState.Success == false {
			msgState.AckChannel <- nil
			close(msgState.AckChannel)
			msgState.Success = true
			//delete(this.SentHashesToChannel, ackMsg.Echo)
		} else {
			log.Debug(fmt.Sprintf("receive duplicate ack  from %s:%d ", host, port))
		}
		this.mapLock.Unlock()
	} else {
		signedMessager, ok := messager.(encoding.SignedMessager)
		log.Trace(fmt.Sprintf("received msg=%s from=%s,expect ack=%s", encoding.MessageType(messager.Cmd()), utils.APex2(signedMessager.GetSender()), utils.HPex(echohash)))
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
			this.sendAck(signedMessager.GetSender(), this.CreateAck(echohash))
		} else {
			//send message to raiden ,and wait result
			log.Trace(fmt.Sprintf("protocol send message to raiden... %s", signedMessager))
			this.ReceivedMessageChannel <- &MessageToRaiden{signedMessager, echohash}
			err, ok = <-this.ReceivedMessageResultChannel
			log.Trace(fmt.Sprintf("protocol receive message response from raiden ok=%v,err=%s", ok, err))
			//only send the Ack if the message was handled without exceptions
			if err == nil && ok {
				ack := this.CreateAck(echohash)
				this.sendAck(signedMessager.GetSender(), ack)
				if this.receivedMessageSaver != nil {
					this.receivedMessageSaver.SaveAck(echohash, messager, ack.Pack())
				}
			} else {
				log.Info(fmt.Sprintf("and raiden report error %s, for Received Message %s", err, utils.StringInterface(signedMessager, 3)))
			}
		}
	}

}

func (this *RaidenProtocol) StopAndWait() {
	log.Info("RaidenProtocol stop...")
	this.onStop = true
	this.Transport.StopAccepting()
	this.mapLock.Lock()
	for k, c := range this.SentHashesToChannel {
		delete(this.SentHashesToChannel, k)
		if !c.Success {
			close(c.AckChannel)
			//close(c.AsyncResult.Result) //caller waiting for result, it must be a successful result.
		}
	}
	//stop sending..
	for _, c := range this.sendingQueueMap {
		close(c)
	}
	this.mapLock.Unlock()
	//what about the outgoing packets, maybe lost
	this.Transport.Stop()
	close(this.ReceivedMessageResultChannel)
	close(this.ReceivedMessageChannel)
	this.quitWaitGroup.Wait()
	log.Info("raiden protocol stop ok...")
}

func (this *RaidenProtocol) Start() {
	this.Transport.Start()
}
