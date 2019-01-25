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

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var errTimeout = errors.New("wait timeout")
var errExpired = errors.New("message expired")

/*
MessageToPhoton message and it's echo hash
*/
type MessageToPhoton struct {
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
	GetChannelStatus(channelIdentifier common.Hash) (int, int64)
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

type queueMessagesAndLock struct {
	messages   []*SentMessageState //todo fixme use channel to avoid lock
	lock       sync.Mutex
	wakeUpChan chan int
}

/*
PhotonProtocol is a UDP protocol,
every message needs a ack to make sure sent success.
*/
type PhotonProtocol struct {
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
	ReceivedMessageChan chan *MessageToPhoton
	/*
		this is a synchronized chan,reading  process message result from photon
	*/
	ReceivedMessageResultChan chan error
	sendingChanMap            map[string]chan *SentMessageState //write to this channel to send a message
	sendingQueueMap           map[string]*queueMessagesAndLock
	receivedMessageSaver      ReceivedMessageSaver
	ChannelStatusGetter       ChannelStatusGetter
	onStop                    bool //flag for stop
	//notify quit
	quitChan chan struct{}
	//receive data
	receiveChan chan []byte
	log         log.Logger
	isReceiving bool
}

// NewPhotonProtocol create PhotonProtocol
func NewPhotonProtocol(transport Transporter, privKey *ecdsa.PrivateKey, channelStatusGetter ChannelStatusGetter) *PhotonProtocol {
	rp := &PhotonProtocol{
		Transport:                 transport,
		privKey:                   privKey,
		retryTimes:                10,
		retryInterval:             time.Millisecond * 6000,
		SentHashesToChannel:       make(map[common.Hash]*SentMessageState),
		ReceivedMessageChan:       make(chan *MessageToPhoton),
		ReceivedMessageResultChan: make(chan error),
		sendingChanMap:            make(map[string]chan *SentMessageState),
		sendingQueueMap:           make(map[string]*queueMessagesAndLock),
		ChannelStatusGetter:       channelStatusGetter,
		quitChan:                  make(chan struct{}),
		receiveChan:               make(chan []byte, 200),
		mapLock:                   sync.Mutex{},
	}
	rp.nodeAddr = crypto.PubkeyToAddress(privKey.PublicKey)
	transport.RegisterProtocol(rp)
	rp.log = log.New("name", utils.APex2(rp.nodeAddr))
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
func (p *PhotonProtocol) SetReceivedMessageSaver(saver ReceivedMessageSaver) {
	p.receivedMessageSaver = saver
}

func (p *PhotonProtocol) sendAck(receiver common.Address, ack *encoding.Ack) {
	p.log.Trace(fmt.Sprintf("send ack EchoHash=%s to %s, ", utils.HPex(ack.Echo), utils.APex2(receiver)))
	err := p.sendRawWitNoAck(receiver, ack.Pack())
	if err != nil {
		log.Warn(fmt.Sprintf("sesendRawWitNoAck err %s ", err))
	}
}
func (p *PhotonProtocol) sendRawAck(receiver common.Address, data []byte) {
	p.log.Trace(fmt.Sprintf("send to %s raw ack", utils.APex2(receiver)))
	err := p.sendRawWitNoAck(receiver, data)
	if err != nil {
		log.Warn(fmt.Sprintf("sesendRawWitNoAck err %s ", err))
	}
}
func (p *PhotonProtocol) sendRawWitNoAck(receiver common.Address, data []byte) error {
	return p.Transport.Send(receiver, data)
}

// SendPing PingSender
func (p *PhotonProtocol) SendPing(receiver common.Address) error {
	ping := encoding.NewPing(utils.NewRandomInt64())
	err := ping.Sign(p.privKey, ping)
	if err != nil {
		return err
	}
	data := ping.Pack()
	return p.sendRawWitNoAck(receiver, data)
}

/*
	message mediatedTransfer  can safely be discarded when channel not exist only more
	当channel被移除后,可以安全的移除待发送的消息,否则会导致新channel无法使用
	(之前的实现是交易中的锁过期后移除,但这可能会导致通道双方状态不同步)
*/
/*
 *	messageCanBeSent : function to check MediatedTransfer can be discarded securely when channel no long exists.
 *
 *	Note that once this channel gets removed, those pending message should also be securely removed,
 *	otherwise new channel can't be created.
 */
func (p *PhotonProtocol) messageCanBeSent(msg encoding.Messager) bool {
	var channelIdentifier common.Hash
	var openBlockNumber int64
	channelIdentifier, openBlockNumber = getMessageChannelIdentifier(msg)
	if channelIdentifier != utils.EmptyHash {
		status, localOpenBlockNumber := p.ChannelStatusGetter.GetChannelStatus(channelIdentifier)
		if status == channeltype.StateInValid {
			p.log.Info(fmt.Sprintf("message cannot be send because of channel status =%d", status))
			return false
		}
		if openBlockNumber != localOpenBlockNumber {
			p.log.Info(fmt.Sprintf("message cannot be send because of channel open block number does't match,now=%d,msg.OpenBlockNumber=%d", localOpenBlockNumber, openBlockNumber))
			return false
		}
	}
	return true
}

/*
此函数会修改Protocol关键Map,所以必须有p.mapLock保护
此函数保证不阻塞
此函数启动的goroutine会自动处理退出问题
*/
func (p *PhotonProtocol) processSentMessageState(receiver common.Address, channelIdentifier common.Hash, msgState *SentMessageState) {
	if channelIdentifier == utils.EmptyHash {
		// 不带balance proof的消息,单独发送,不使用队列
		log.Trace(fmt.Sprintf("send message EchoHash=%s without SendingQueue", utils.HPex(msgState.EchoHash)))
		go p.sendMessage(receiver, msgState)
		return
	}
	key := fmt.Sprintf("%s-%s", receiver.String(), channelIdentifier.String())
	p.mapLock.Lock()
	defer p.mapLock.Unlock()
	ql, ok := p.sendingQueueMap[key]
	if ok {
		//存在,顺序投递并唤醒发送goroutine,然后退出
		log.Trace(fmt.Sprintf("send message EchoHash=%s with exist SendingQueue", utils.HPex(msgState.EchoHash)))
		ql.messages = append(ql.messages, msgState)
		select {
		case ql.wakeUpChan <- 0:
		default:
			// never block
		}
		return
	}
	// 创建ql并启动goroutine发送
	ql = &queueMessagesAndLock{
		wakeUpChan: make(chan int),
	}
	ql.messages = append(ql.messages, msgState)
	p.sendingQueueMap[key] = ql
	// 启动该ql的发送routing
	go func() {
		defer rpanic.PanicRecover(fmt.Sprintf("protocol ChannelQueue %s", key))
		log.Trace(fmt.Sprintf("send message EchoHash=%s with New SendingQueue", utils.HPex(msgState.EchoHash)))
		for {
			p.mapLock.Lock()
			if len(ql.messages) == 0 {
				p.mapLock.Unlock()
				// goroutine保留一段时间,防止频繁创建
				select {
				case <-p.quitChan: //其他地方要求退出了
					return
				case <-ql.wakeUpChan: // 阻塞等待新message唤醒,goroutine 一直保留
					continue
					//case <-time.After(60 * time.Second):
					//	// 清理sendingQueueMap, 退出goroutine
					//	p.mapLock.Lock()
					//	delete(p.sendingQueueMap, key)
					//	p.mapLock.Unlock()
					//	log.Trace(fmt.Sprintf("sendingQueueMap %s released", key))
					//	return
				}
			}
			msg := ql.messages[0]
			ql.messages = ql.messages[1:]
			p.mapLock.Unlock()
			p.sendMessage(receiver, msg)
		}
	}()
}

func (p *PhotonProtocol) sendMessage(receiver common.Address, msgState *SentMessageState) {
	p.log.Trace(fmt.Sprintf("send to %s,msg=%s, echohash=%s",
		utils.APex2(msgState.ReceiverAddress), msgState.Message,
		utils.HPex(msgState.EchoHash)))
	for {
		if !p.messageCanBeSent(msgState.Message) {
			msgState.AsyncResult.Result <- errExpired
			p.mapLock.Lock()
			delete(p.SentHashesToChannel, msgState.EchoHash)
			p.mapLock.Unlock()
			return
		}
		nextTimeout := timeoutExponentialBackoff(p.retryTimes, p.retryInterval, p.retryInterval*10)
		err := p.sendRawWitNoAck(receiver, msgState.Data)
		if err != nil {
			p.log.Info(fmt.Sprintf("sendRawWitNoAck msg echoHash=%s error %s", utils.HPex(msgState.EchoHash), err.Error()))
		}
		timeout := time.After(nextTimeout())
		var ok bool
		select {
		case _, ok = <-msgState.AckChannel:
			if ok {
				p.log.Trace(fmt.Sprintf("msg=%s EchoHash=%s, sent success", encoding.MessageType(msgState.Message.Cmd()), utils.HPex(msgState.EchoHash)))
				msgState.AsyncResult.Result <- nil
				p.mapLock.Lock()
				delete(p.SentHashesToChannel, msgState.EchoHash)
				p.mapLock.Unlock()

			} else {
				p.log.Info(fmt.Sprintf("sendMessage EchoHash=%s stop retry, because of chan closed", utils.HPex(msgState.EchoHash)))
			}
			return
		case <-timeout: //retry
			// 如果是matrix且对方不在线,挂起并等待唤醒
			_, isOnline := p.Transport.NodeStatus(receiver)
			transport, ok1 := p.Transport.(*MatrixMixTransport)
			if ok1 && !isOnline && transport != nil {
				log.Warn(fmt.Sprintf("receiver %s is not online,sleep until when he back online", receiver.String()))
				wakeUpChan := make(chan int)
				// 向transport注册wakeUpChan
				transport.RegisterWakeUpChan(receiver, wakeUpChan)
				// 挂起并等待对方上线
				<-wakeUpChan
				// 继续发送并注销wakeUpChan
				transport.UnRegisterWakeUpChan(receiver)
			}
		case <-p.quitChan:
			return
		}
	}
}

func getMessageChannelIdentifier(msg encoding.Messager) (common.Hash, int64) {
	var channelIdentifier common.Hash
	var openBlockNumber int64
	switch msg2 := msg.(type) {
	case *encoding.DirectTransfer:
		channelIdentifier = msg2.ChannelIdentifier
		openBlockNumber = msg2.OpenBlockNumber
	case *encoding.MediatedTransfer:
		channelIdentifier = msg2.ChannelIdentifier
		openBlockNumber = msg2.OpenBlockNumber
	case *encoding.AnnounceDisposedResponse:
		channelIdentifier = msg2.ChannelIdentifier
		openBlockNumber = msg2.OpenBlockNumber
	case *encoding.UnLock:
		channelIdentifier = msg2.ChannelIdentifier
		openBlockNumber = msg2.OpenBlockNumber
	case *encoding.RemoveExpiredHashlockTransfer:
		channelIdentifier = msg2.ChannelIdentifier
		openBlockNumber = msg2.OpenBlockNumber
	case *encoding.SettleRequest:
		channelIdentifier = msg2.ChannelIdentifier
		openBlockNumber = msg2.OpenBlockNumber
	case *encoding.WithdrawRequest:
		channelIdentifier = msg2.ChannelIdentifier
		openBlockNumber = msg2.OpenBlockNumber
	}
	return channelIdentifier, openBlockNumber
}

/*
	msg should be signed.
	msg must be sent success.
*/
func (p *PhotonProtocol) sendWithResult(receiver common.Address,
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
		log.Info(fmt.Sprintf(`trying to send same message, this should only occur when token swap,
because they sending same SecretRequest and RevealSecret,msg=%s`, msg))
		p.mapLock.Unlock()
		result = msgState.AsyncResult
		return
	}
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
	channelIdentifier, _ := getMessageChannelIdentifier(msg)
	p.processSentMessageState(receiver, channelIdentifier, msgState)
	return
}

// SendAndWait send this packet and wait ack until timeout
func (p *PhotonProtocol) SendAndWait(receiver common.Address, msg encoding.Messager, timeout time.Duration) error {
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
func (p *PhotonProtocol) SendAsync(receiver common.Address, msg encoding.Messager) *utils.AsyncResult {
	return p.sendWithResult(receiver, msg)
}

// CreateAck creat a ack message,
func (p *PhotonProtocol) CreateAck(echohash common.Hash) *encoding.Ack {
	return encoding.NewAck(p.nodeAddr, echohash)
}

// GetNetworkStatus return `addr` node's network status
func (p *PhotonProtocol) GetNetworkStatus(addr common.Address) (deviceType string, isOnline bool) {
	return p.Transport.NodeStatus(addr)
}

func (p *PhotonProtocol) receive(data []byte) {
	//todo fix 使用可以反复使用的缓冲区,而不是每次都分配.
	cdata := make([]byte, len(data))
	copy(cdata, data)

	//p.log.Trace(fmt.Sprintf("try to send receive data l=%d,message=%s", len(cdata), encoding.MessageType(cdata[0])))
	p.receiveChan <- cdata
	//p.log.Trace(fmt.Sprintf("receive complete l=%d", len(cdata)))
}

func (p *PhotonProtocol) loop() {
	p.isReceiving = true
	for {
		select {
		case <-p.quitChan:
			return
		case data := <-p.receiveChan:
			p.receiveInternal(data)
		}
	}
}

func (p *PhotonProtocol) receiveInternal(data []byte) {
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
		p.log.Debug(fmt.Sprintf("receive ack ,EchoHash=%s", utils.HPex(ackMsg.Echo)))
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
		p.log.Trace(fmt.Sprintf("received msg=%s from=%s,expect ack EchoHash=%s", messager, utils.APex2(signedMessager.GetSender()), utils.HPex(echohash)))
		if !ok {
			p.log.Warn("message should be signed except for ack")
			return
		}
		if messager.Cmd() == encoding.PingCmdID { //send ack
			p.sendAck(signedMessager.GetSender(), p.CreateAck(echohash))
		} else {
			//send message to photon ,and wait result
			p.log.Trace(fmt.Sprintf("protocol send message to photon... %s", signedMessager))
			p.ReceivedMessageChan <- &MessageToPhoton{signedMessager, echohash}
			select {
			case err, ok = <-p.ReceivedMessageResultChan:
			case <-p.quitChan:
				ok = false
				err = errors.New("protocol stoped")
			}
			p.log.Trace(fmt.Sprintf("protocol receive message response from photon ok=%v,err=%v", ok, err))
			//only send the Ack if the message was handled without exceptions
			if err == nil && ok {
				ack := p.CreateAck(echohash)
				p.sendAck(signedMessager.GetSender(), ack)
				if p.receivedMessageSaver != nil {
					p.receivedMessageSaver.SaveAck(echohash, messager, ack.Pack())
				}
			} else {
				p.log.Info(fmt.Sprintf("and photon report error %s, for Received Message %s", err, utils.StringInterface(signedMessager, 3)))
			}
		}
	}

}

// StopAndWait stop andf wait for clean.
func (p *PhotonProtocol) StopAndWait() {
	p.log.Info("PhotonProtocol stop...")
	p.onStop = true
	close(p.quitChan)
	p.Transport.StopAccepting()
	//what about the outgoing packets, maybe lost
	p.Transport.Stop()

	p.log.Info("photon protocol stop ok...")
}

// Start photon protocol
func (p *PhotonProtocol) Start(receive bool) {
	if receive {
		go p.loop()
	}
	p.Transport.Start()
}

//StartReceive start event loop if not start,otherwise crash
func (p *PhotonProtocol) StartReceive() {
	if p.isReceiving {
		panic("can not receive twice")
	}
	go p.loop()
}

// NodeInfo get from user
type NodeInfo struct {
	Address    string `json:"address"`
	IPPort     string `json:"ip_port"`
	DeviceType string `json:"device_type"` // must be mobile?
}

// UpdateMeshNetworkNodes update nodes in this intranet
func (p *PhotonProtocol) UpdateMeshNetworkNodes(nodes []*NodeInfo) error {
	//p.log.Trace(fmt.Sprintf("nodes=%s", utils.StringInterface(nodes, 3)))
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
	if transport, ok := p.Transport.(*MixTransport); ok {
		transport.udp.setHostPort(nodesmap)
	} else if transport, ok := p.Transport.(*MatrixMixTransport); ok {
		transport.udp.setHostPort(nodesmap)
	} else if transport, ok := p.Transport.(*UDPTransport); ok {
		transport.setHostPort(nodesmap)
	} else {
		return errors.New("no need to register nodes while udp doesn't work")
	}
	return nil
}
