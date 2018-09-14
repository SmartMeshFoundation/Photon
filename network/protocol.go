package network

import (
	"crypto/ecdsa"

	"encoding/hex"

	"reflect"

	"fmt"
	"time"

	"sync"

	"errors"

	"net"
	"strconv"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var errTimeout = errors.New("wait timeout")
var errExpired = errors.New("message expired")

/*
MessageToRaiden message and it's echo hash
*/
type MessageToRaiden struct {
	Msg      encoding.SignedMessager
	EchoHash common.Hash
}

// SentMessageState is the state of message on sending
type SentMessageState struct {
	AsyncResult     *utils.AsyncResult
	AckChannel      chan error
	ReceiverAddress common.Address
	Success         bool

	Message  encoding.Messager //message to send
	EchoHash common.Hash       //message echo hash
	Data     []byte            //packed message
}

// PingSender do send ping task
type PingSender interface {
	//SendPing send a ping to receiver,and not block
	SendPing(receiver common.Address) error
}

/*
ChannelStatusGetter get the status of channel address, so sender can remove msg based on channel status
	for example :
		A send B a mediated transfer, but B is offline
		when B is online ,this transfer is invalid, so A will never receive ack
		if A  remove this msg, this channel can not be used only more.
		but if A does't remove, when A settle/withdraw/reopen channel with B,this msg will make the new channel unusable too.
		So A need to remove channel when channel status change.
*/
type ChannelStatusGetter interface {
	GetChannelStatus(channelIdentifier common.Hash) int
}

/*
BlockNumberGetter get the lastest block number,so sender can remove expired mediated transfer.
	for example :
		A send B a mediated transfer, but B is offline
		when B is online ,this transfer is invalid, so A will never receive  ack ,so A will try forever.
		message secret,secretRequest,revealSecret won't allow error
*/
//type BlockNumberGetter interface {
//	// GetBlockNumber return latest block number
//	GetBlockNumber() int64
//}

type timeoutGenerator func() time.Duration

/*
	Timeouts generator with an exponential backoff strategy.

	Timeouts start spaced by `timeout`, after `retries` exponentially increase
	the retry delays until `maximum`, then maximum is returned indefinitely.
*/
func timeoutExponentialBackoff(retries int, timeout, maximumTimeout time.Duration) timeoutGenerator {
	tries := 1
	return func() time.Duration {
		tries++
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

/*
RaidenProtocol is a UDP protocol,
every message needs a ack to make sure sent success.
*/
type RaidenProtocol struct {
	Transport           Transporter
	privKey             *ecdsa.PrivateKey
	nodeAddr            common.Address
	SentHashesToChannel map[common.Hash]*SentMessageState
	retryTimes          int
	retryInterval       time.Duration
	mapLock             sync.Mutex
	statusLock          sync.RWMutex
	/*
		message from other nodes
	*/
	ReceivedMessageChan chan *MessageToRaiden
	/*
		this is a synchronized chan,reading  process message result from raiden
	*/
	ReceivedMessageResultChan chan error
	sendingQueueMap           map[string]chan *SentMessageState //write to this channel to send a message
	receivedMessageSaver      ReceivedMessageSaver
	ChannelStatusGetter       ChannelStatusGetter
	onStop                    bool //flag for stop
	//notify quit
	quitChan chan struct{}
	//receive data
	receiveChan chan []byte
	log         log.Logger
}

// NewRaidenProtocol create RaidenProtocol
func NewRaidenProtocol(transport Transporter, privKey *ecdsa.PrivateKey, channelStatusGetter ChannelStatusGetter) *RaidenProtocol {
	rp := &RaidenProtocol{
		Transport:                 transport,
		privKey:                   privKey,
		retryTimes:                10,
		retryInterval:             time.Millisecond * 6000,
		SentHashesToChannel:       make(map[common.Hash]*SentMessageState),
		ReceivedMessageChan:       make(chan *MessageToRaiden),
		ReceivedMessageResultChan: make(chan error),
		sendingQueueMap:           make(map[string]chan *SentMessageState),
		ChannelStatusGetter:       channelStatusGetter,
		quitChan:                  make(chan struct{}),
		receiveChan:               make(chan []byte, 20),
	}
	rp.nodeAddr = crypto.PubkeyToAddress(privKey.PublicKey)
	transport.RegisterProtocol(rp)
	rp.log = log.New("name", utils.APex2(rp.nodeAddr))
	go rp.loop()
	return rp
}

// New create new object from sample.
func New(sample interface{}) interface{} {
	t := reflect.ValueOf(sample)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	tt := t.Type()
	v := reflect.New(tt).Interface()
	return v
}

// SetReceivedMessageSaver set db saver
func (p *RaidenProtocol) SetReceivedMessageSaver(saver ReceivedMessageSaver) {
	p.receivedMessageSaver = saver
}

func (p *RaidenProtocol) sendAck(receiver common.Address, ack *encoding.Ack) {
	p.log.Trace(fmt.Sprintf("send to %s, ack=%s", utils.APex2(receiver), ack))
	err := p.sendRawWitNoAck(receiver, ack.Pack())
	if err != nil {
		log.Warn(fmt.Sprintf("sesendRawWitNoAck err %s ", err))
	}
}
func (p *RaidenProtocol) sendRawAck(receiver common.Address, data []byte) {
	p.log.Trace(fmt.Sprintf("send to %s raw ack", utils.APex2(receiver)))
	err := p.sendRawWitNoAck(receiver, data)
	if err != nil {
		log.Warn(fmt.Sprintf("sesendRawWitNoAck err %s ", err))
	}
}
func (p *RaidenProtocol) sendRawWitNoAck(receiver common.Address, data []byte) error {
	return p.Transport.Send(receiver, data)
}

// SendPing PingSender
func (p *RaidenProtocol) SendPing(receiver common.Address) error {
	ping := encoding.NewPing(utils.NewRandomInt64())
	err := ping.Sign(p.privKey, ping)
	if err != nil {
		return err
	}
	data := ping.Pack()
	return p.sendRawWitNoAck(receiver, data)
}

/*
	message mediatedTransfer  can safely be discarded when channel not open only more
	当channel被移除后,可以安全的移除待发送的消息,否则会导致新channel无法使用
	(之前的实现是交易中的锁过期后移除,但这可能会导致通道双方状态不同步)
*/
func (p *RaidenProtocol) messageCanBeSent(msg encoding.Messager, channelIdentifier common.Hash) bool {
	if channelIdentifier != utils.EmptyHash {
		status := channeltype.StateOpened
		switch msg.(type) {
		case *encoding.DirectTransfer, *encoding.MediatedTransfer,
			*encoding.RemoveExpiredHashlockTransfer, *encoding.UnLock, *encoding.AnnounceDisposedResponse,
			*encoding.SettleRequest, *encoding.WithdrawRequest:
			status = p.ChannelStatusGetter.GetChannelStatus(channelIdentifier)
		}
		if status == channeltype.StateInValid {
			p.log.Info(fmt.Sprintf("message cannot be send because of channel status =%d", status))
			return false
		}
	}
	return true
}

func (p *RaidenProtocol) getChannelQueue(receiver common.Address, channelAddr common.Hash) chan<- *SentMessageState {

	p.mapLock.Lock()
	defer p.mapLock.Unlock()
	key := fmt.Sprintf("%s-%s", receiver.String(), channelAddr.String())
	var sendingChan chan *SentMessageState
	var ok bool
	/*
		no channelAddr means that p message doesn't need ordered.
		if  channel address is not nil,it must contain a new balance proof.
		balance proof must be sent ordered
	*/
	if channelAddr == utils.EmptyHash {
		sendingChan = make(chan *SentMessageState, 1) //should not block sender
	} else {
		sendingChan, ok = p.sendingQueueMap[key]
		if ok {
			return sendingChan
		}
		sendingChan = make(chan *SentMessageState, 1000) //should not block sender
		p.sendingQueueMap[key] = sendingChan
	}
	go func() {
		defer rpanic.PanicRecover(fmt.Sprintf("protocol ChannelQueue %s", key))
		/*
			1. if p packet is on sending, retry send immediately
			2. retry infinite, until receive a ack
			3. p message should be sent by caller after restart.

			caller can read from chan reusltChannel to get if p packet is successfully sent to receiver
		*/
	labelNextMessage:
		for {
			p.log.Trace(fmt.Sprintf("queue %s try send next message", key))
			var msgState *SentMessageState
			var ok bool
			select {
			case msgState, ok = <-sendingChan:
			case <-p.quitChan:
				return
			}
			if !ok {
				p.log.Info(fmt.Sprintf("queue %s quit, because of chan closed", key))
				return
			}
			p.log.Trace(fmt.Sprintf("send to %s,msg=%s, echoash=%s",
				utils.APex2(msgState.ReceiverAddress), msgState.Message,
				utils.HPex(msgState.EchoHash)))
			for {
				if !p.messageCanBeSent(msgState.Message, channelAddr) {
					msgState.AsyncResult.Result <- errExpired
					break
				}
				nextTimeout := timeoutExponentialBackoff(p.retryTimes, p.retryInterval, p.retryInterval*10)
				err := p.sendRawWitNoAck(receiver, msgState.Data)
				if err != nil {
					p.log.Info(fmt.Sprintf("sendRawWitNoAck %s msg error %s", key, err.Error()))
				}
				timeout := time.After(nextTimeout())
				select {
				case _, ok = <-msgState.AckChannel:
					if ok {
						p.log.Trace(fmt.Sprintf("msg=%s, sent success :%s", encoding.MessageType(msgState.Message.Cmd()), utils.HPex(msgState.EchoHash)))
						msgState.AsyncResult.Result <- nil
						goto labelNextMessage
					} else {
						//message must send success, otherwise keep trying...
						p.log.Info(fmt.Sprintf("queue %s quit, because of chan closed", key))
						return //user call stop
					}
				case <-timeout: //retry
				case <-p.quitChan:
					return
				}
			}
		}

	}()
	return sendingChan
}

func getMessageChannelAddress(msg encoding.Messager) common.Hash {
	var channelAddress common.Hash
	switch msg2 := msg.(type) {
	case *encoding.DirectTransfer:
		channelAddress = msg2.ChannelIdentifier
	case *encoding.MediatedTransfer:
		channelAddress = msg2.ChannelIdentifier
	case *encoding.AnnounceDisposedResponse:
		channelAddress = msg2.ChannelIdentifier
	case *encoding.UnLock:
		channelAddress = msg2.ChannelIdentifier
	case *encoding.RemoveExpiredHashlockTransfer:
		channelAddress = msg2.ChannelIdentifier
	}
	return channelAddress
}

/*
	msg should be signed.
	msg must be sent success.
*/
func (p *RaidenProtocol) sendWithResult(receiver common.Address,
	msg encoding.Messager) (result *utils.AsyncResult) {
	//no more message...
	if p.onStop {
		return utils.NewAsyncResult()
	}
	if true {
		signed, ok := msg.(encoding.SignedMessager)
		if ok && signed.GetSender() == utils.EmptyAddress {
			p.log.Error("send unsigned message")
			panic("send unsigned message")
		}
	}
	data := msg.Pack()
	echohash := utils.Sha3(data, receiver[:])
	p.mapLock.Lock()
	msgState, ok := p.SentHashesToChannel[echohash]
	if ok {
		p.mapLock.Unlock()
		result = msgState.AsyncResult
		return
	}
	p.log.Debug(fmt.Sprintf("send msg=%s to=%s,expected hash=%s", encoding.MessageType(msg.Cmd()), utils.APex2(receiver), utils.HPex(echohash)))
	msgState = &SentMessageState{
		AsyncResult:     utils.NewAsyncResult(),
		ReceiverAddress: receiver,
		AckChannel:      make(chan error, 1),
		Message:         msg,
		Data:            data,
		EchoHash:        echohash,
	}
	p.SentHashesToChannel[echohash] = msgState
	p.mapLock.Unlock()
	result = msgState.AsyncResult
	channelAddress := getMessageChannelAddress(msg)
	//make sure not block
	p.getChannelQueue(receiver, channelAddress) <- msgState
	return
}

// SendAndWait send this packet and wait ack until timeout
func (p *RaidenProtocol) SendAndWait(receiver common.Address, msg encoding.Messager, timeout time.Duration) error {
	var err error
	result := p.sendWithResult(receiver, msg)
	timeoutCh := time.After(timeout)
	select {
	case err = <-result.Result:
	case <-timeoutCh:
		err = errTimeout
	case <-p.quitChan:
		err = errTimeout
	}
	return err
}

// SendAsync send a message asynchronize ,notify by `AsyncResult`
func (p *RaidenProtocol) SendAsync(receiver common.Address, msg encoding.Messager) *utils.AsyncResult {
	return p.sendWithResult(receiver, msg)
}

// CreateAck creat a ack message,
func (p *RaidenProtocol) CreateAck(echohash common.Hash) *encoding.Ack {
	return encoding.NewAck(p.nodeAddr, echohash)
}

// GetNetworkStatus return `addr` node's network status
func (p *RaidenProtocol) GetNetworkStatus(addr common.Address) (deviceType string, isOnline bool) {
	return p.Transport.NodeStatus(addr)
}

func (p *RaidenProtocol) receive(data []byte) {
	//todo fix ,remove copy and fix deadlock of send and receive
	cdata := make([]byte, len(data))
	copy(cdata, data)
	p.receiveChan <- cdata
}

func (p *RaidenProtocol) loop() {
	for {
		select {
		case <-p.quitChan:
			return
		case data := <-p.receiveChan:
			p.receiveInternal(data)
		}
	}
}

func (p *RaidenProtocol) receiveInternal(data []byte) {
	if len(data) > params.UDPMaxMessageSize {
		p.log.Error("receive packet larger than maximum size :", len(data))
		return
	}
	//ignore incomming message when stop
	if p.onStop {
		return
	}
	cmdid := int(data[0])
	messager, ok := encoding.MessageMap[cmdid]
	if !ok {
		p.log.Warn("receive unknown message:", hex.Dump(data))
		return
	}
	messager = New(messager).(encoding.Messager)
	err := messager.UnPack(data)
	if err != nil {
		p.log.Warn(fmt.Sprintf("message unpack error : %s", err))
		return
	}
	echohash := utils.Sha3(data, p.nodeAddr[:])
	if p.receivedMessageSaver != nil && messager.Cmd() != encoding.AckCmdID {
		ackdata := p.receivedMessageSaver.GetAck(echohash)
		if len(ackdata) > 0 {
			sm, ok := messager.(encoding.SignedMessager)
			if !ok {
				p.log.Error(fmt.Sprintf("received a message %s, not ack ,and don't signed", messager))
				return
			}
			p.sendRawAck(sm.GetSender(), ackdata)
			return
		}
	}
	if messager.Cmd() == encoding.AckCmdID { //some one may be waiting p ack
		ackMsg := messager.(*encoding.Ack)
		p.log.Debug(fmt.Sprintf("receive ack ,hash=%s", utils.HPex(ackMsg.Echo)))
		p.mapLock.Lock()
		msgState, ok := p.SentHashesToChannel[ackMsg.Echo]
		if ok && msgState.Success == false {
			msgState.AckChannel <- nil
			close(msgState.AckChannel)
			msgState.Success = true
		} else {
			p.log.Debug(fmt.Sprintf("receive duplicate ack  from %s", utils.APex(ackMsg.Sender)))
		}
		p.mapLock.Unlock()
	} else {
		signedMessager, ok := messager.(encoding.SignedMessager)
		p.log.Trace(fmt.Sprintf("received msg=%s from=%s,expect ack=%s", messager, utils.APex2(signedMessager.GetSender()), utils.HPex(echohash)))
		if !ok {
			p.log.Warn("message should be signed except for ack")
			return
		}
		if messager.Cmd() == encoding.PingCmdID { //send ack
			p.sendAck(signedMessager.GetSender(), p.CreateAck(echohash))
		} else {
			//send message to raiden ,and wait result
			p.log.Trace(fmt.Sprintf("protocol send message to raiden... %s", signedMessager))
			p.ReceivedMessageChan <- &MessageToRaiden{signedMessager, echohash}
			select {
			case err, ok = <-p.ReceivedMessageResultChan:
			case <-p.quitChan:
				ok = false
				err = errors.New("protocol stoped")
			}
			p.log.Trace(fmt.Sprintf("protocol receive message response from raiden ok=%v,err=%v", ok, err))
			//only send the Ack if the message was handled without exceptions
			if err == nil && ok {
				ack := p.CreateAck(echohash)
				p.sendAck(signedMessager.GetSender(), ack)
				if p.receivedMessageSaver != nil {
					p.receivedMessageSaver.SaveAck(echohash, messager, ack.Pack())
				}
			} else {
				p.log.Info(fmt.Sprintf("and raiden report error %s, for Received Message %s", err, utils.StringInterface(signedMessager, 3)))
			}
		}
	}

}

// StopAndWait stop andf wait for clean.
func (p *RaidenProtocol) StopAndWait() {
	p.log.Info("RaidenProtocol stop...")
	p.onStop = true
	close(p.quitChan)
	p.Transport.StopAccepting()
	//what about the outgoing packets, maybe lost
	p.Transport.Stop()

	p.log.Info("raiden protocol stop ok...")
}

// Start raiden protocol
func (p *RaidenProtocol) Start() {
	p.Transport.Start()
}

// NodeInfo get from user
type NodeInfo struct {
	Address    string `json:"address"`
	IPPort     string `json:"ip_port"`
	DeviceType string `json:"device_type"` // must be mobile?
}

// UpdateMeshNetworkNodes update nodes in this intranet
func (p *RaidenProtocol) UpdateMeshNetworkNodes(nodes []*NodeInfo) error {
	p.log.Trace(fmt.Sprintf("nodes=%s", utils.StringInterface(nodes, 3)))
	nodesmap := make(map[common.Address]*net.UDPAddr)
	for _, n := range nodes {
		addr := common.HexToAddress(n.Address)
		host, port, err := net.SplitHostPort(n.IPPort)
		if err != nil {
			return err
		}
		porti, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		ua := &net.UDPAddr{
			IP:   net.ParseIP(host),
			Port: porti,
		}
		nodesmap[addr] = ua
	}
	if transport, ok := p.Transport.(*MixTransporter); ok {
		transport.udp.setHostPort(nodesmap)
	} else if transport, ok := p.Transport.(*UDPTransport); ok {
		transport.setHostPort(nodesmap)
	} else {
		return errors.New("no need to register nodes while udp doesn't work")
	}
	return nil
}
