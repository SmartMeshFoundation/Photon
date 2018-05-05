package ice

import (
	"flag"
	"fmt"
	"net"

	_ "net/http/pprof"

	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/turn"
)

var (
	network = flag.String("net", "udp", "network to listen")
	address = flag.String("addr", "0.0.0.0:3478", "address to listen")
	profile = flag.Bool("profile", false, "profile")
)
var (
	ErrTimeout            = errors.New("timed out")
	ErrInvalidMessage     = errors.New("invalid message")
	ErrDuplicateWaiter    = errors.New("waiter with uid already exists")
	ErrWaiterClosed       = errors.New("waiter closed")
	ErrClientDisconnected = errors.New("client disconnected")
)

type serverSockMode int

const (
	/*
		服务器启动以后进入的是等待ice 协商阶段,这时收到的数据全部都是 stun.Message
	*/
	StageNegotiation serverSockMode = iota
	/*
		ICE 协商完毕,建立了通道,我这里没有经过turn Server 中转来接收数据. 所以这里面是不包含 channel data 的,如果不是 stun.message, 那就是直接交付给用户的数据
	*/
	StunModeData
	/*
		我发送接收数据都要经过 TurnServer 中转,所有的 data 都是 channel 通道,这种情况下数据全都解析为 stun message 或者 channel data
	*/
	TurnModeData
)
const (
	MinChannelNumber = 0x4000
	MaxChannelNUmber = 0x7fff
)

type sendRequest struct {
	data   []byte
	toaddr net.Addr
}

/*
StunServerSock 是用来 ICE 协商以及协商成功以后节点之间直接发送数据需要的.
ICE 协商时需要从指定的 ip 地址上发送stun message.
ICE 协商完毕以后,节点之间互相发送数据也需要 Server 保持在线,因为需要接收来自对方的 SendIndication/BindIndication 来保持连接有效性.
如果是 turn server 中转,还需要 ChannelNumber 信息.


Server 可能收到以下消息
1. ICE 协商过程中的 BindRequest, 这个消息是需要短期凭证的.
2. 来自 Stun/turn server 的 refresh reponse.
3. 来自 turn server 的 DataIndication 这是对 peer 的 BindResponse 的封装
4. 连接建立以后,通信的数据,可能是 channel data 封装的数据,也可能是直接的数据.
5. 来自对方的SendIndication/BindIndication,用来保持连接有效性的. 比如较长时间没有通信,仍然需要保持连接有效.
同时也要通过 Server 的 Connection 发送消息:
主要发送如下消息:
除了上面的1,4,5,还有就是
用 SendIndication 封装的由 turn server relay的 BindRequest.
*/
type cachedResponse struct {
	cacheTime time.Time
	msg       *stun.Message //todo need store a full copy  or a pointer?
}
type sendreq struct {
	data []byte
	to   net.Addr
}
type StunServerSock struct {
	Addr                  string //address listening on
	mode                  serverSockMode
	LogAllErrors          bool
	cb                    ServerSockCallbacker
	c                     net.PacketConn
	channelNumber2Address map[int]string // channel number-> address
	address2ChannelNumber map[string]int
	waiters               map[stun.TransactionID]chan *serverSockResponse
	lock                  sync.RWMutex
	syncMessageTimeout    time.Duration //default 10 seconds?
	Name                  string
	cachedResponse        map[stun.TransactionID]*cachedResponse //重复的 bindingrequest, 就不要提交给上层了.
	sendchan              chan *sendreq
	stoped                bool
	log                   log.Logger
}
type serverSockResponse struct {
	res  *stun.Message
	from string
}
type ServerSockCallbacker interface {
	/*
	 收到一个 stun.Message, 可能是 Bind Request/Bind Response 等等.
	*/
	RecieveStunMessage(localAddr, remoteAddr string, msg *stun.Message)
	/*
		ICE 协商建立连接以后,收到了对方发过来的数据,可能是经过 turn server 中转的 channel data( 不接受 sendData data request),也可能直接是数据.
		如果是经过 turn server 中转的, channelNumber 一定介于0x4000-0x7fff 之间.否则一定为0
	*/
	ReceiveData(localAddr, peerAddr string, data []byte)
}

var (
	software          = stun.NewSoftware("nkbai@163.com/ice")
	errNotSTUNMessage = errors.New("not stun message")
)

func (s *StunServerSock) serveConn(c net.PacketConn, req *stun.Message) error {
	if c == nil {
		return nil
	}
	buf := make([]byte, 1024)
	n, addr, err := c.ReadFrom(buf)
	if err != nil {
		s.log.Trace(fmt.Sprintf("ReadFrom: %v", err))
		return err
	}
	raw := buf[:n]
	if _, err = req.Write(raw); err != nil {
		s.dataReceived(udpAddrToAddr(addr), raw)
		return nil
	}
	if req.Type == stun.BindingIndication || req.Type == turn.SendIndication {
		return nil //ignore indication ,只是为了保持心跳而已.
	}
	s.stunMessageReceived(s.Addr, addr.String(), req)
	return nil
}

/*
from: address sendData this message directly.
peerAddr: address who really sendData this message.
在 stun 模式下,两者完全一致,只有在 turn 中转情况下,两者才不一致,
turn 模式下: from 是 turnserver 的地址
peerAddr 才是真正的通信节点地址
*/
func (s *StunServerSock) dataReceived(peerAddr string, data []byte) {
	s.log.Trace(fmt.Sprintf("---- recevied data from %s,len=%d -----", peerAddr, len(data)))
	if s.cb != nil {
		s.cb.ReceiveData(s.Addr, peerAddr, data)
	}
}
func (s *StunServerSock) stunMessageReceived(localaddr, from string, msg *stun.Message) {
	s.log.Trace(fmt.Sprintf("--receive stun message %s<----%s  --\n%s\n", localaddr, from, msg))
	if msg.Type.Method == stun.MethodChannelData {
		s.log.Trace(fmt.Sprintf("\n%s", hex.Dump(msg.Raw)))
		//debug.PrintStack()
	}
	var err error
	/*
		收到 channeldata 要特殊处理,如果是 turn server 模式下,
		如果是在 negiotiation 阶段,说明出错了.
		如果是 stunmode, 说明解析错了,把普通的 data 解析成了 channeldata 了
	*/
	if msg.Type.Method == stun.MethodChannelData {
		if s.mode == StageNegotiation {
			s.log.Error(fmt.Sprintf("receive data error when negiotiation"))
			/*
				在 channel binding success 和 changemode 之间接收到了数据怎么办?直接丢弃,反正对方会重传.
			*/
			//s.dataReceived(from, msg.Raw)
			return
		} else if s.mode == StunModeData {
			/*
				收到了普通的数据,但是被误判为 channelData, 直接纠正即可.
			*/
			s.dataReceived(from, msg.Raw)
			return
		} else if s.mode == TurnModeData {
			var data turn.ChannelData
			err = data.GetFrom(msg)
			if err != nil {
				s.log.Error(fmt.Sprintf("received channel data,but Channel Data err:%s", err))
				return
			}
			peer, ok := s.channelNumber2Address[int(data.ChannelNumber)]
			if !ok {
				s.log.Info(fmt.Sprintf("received data ,but wrong channel number got %d  ", data.ChannelNumber))
				return
			}
			s.dataReceived(peer, data.Data)
		}
	}
	ch, ok := s.getAndRemoveWaiter(msg.TransactionID)
	if ok {
		ch <- &serverSockResponse{msg, from} //对一个消息的 response.提供来自于什么地方,有可能是第三方伪造的消息?
		close(ch)
		return
	}
	if s.checkCachedResponse(msg, from) {
		return
	}
	//需要报告给上层的其他消息
	if s.cb != nil {
		s.cb.RecieveStunMessage(localaddr, from, msg)
	}
}

//如果对应的消息应答,已经缓存了,直接发送即可.
func (s *StunServerSock) checkCachedResponse(req *stun.Message, from string) bool {
	if len(s.cachedResponse) <= 0 {
		return false
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now()
	for id, c := range s.cachedResponse {
		if c.cacheTime.Add(StunResponseCacheDuration).Before(now) {
			delete(s.cachedResponse, id)
		}
	}
	for _, c := range s.cachedResponse {
		if c.msg.Type.Method == req.Type.Method && c.msg.TransactionID == req.TransactionID {
			s.log.Trace(fmt.Sprintf("id %s duplicated", hex.EncodeToString(req.TransactionID[:])))
			s.sendData(c.msg.Raw, s.Addr, from)
			return true
		}
	}
	return false
}

//sendData packet to peer
func (s *StunServerSock) sendData(data []byte, fromaddr, toaddr string) (err error) {
	if s.Addr != fromaddr {
		panic(fmt.Sprintf("each binding..., me=%s,got=%s", s.Addr, fromaddr))
	}
	if s.stoped {
		return
	}
	s.sendchan <- &sendreq{data, addrToUdpAddr(toaddr)}
	return
}

func (s *StunServerSock) sendStunMessageAsync(msg *stun.Message, fromaddr, toaddr string) error {
	s.log.Trace(fmt.Sprintf("---sendData stun message %s-->%s ---\n%s\n", s.Addr, toaddr, msg))
	if msg.Type.Class == stun.ClassSuccessResponse || msg.Type.Class == stun.ClassErrorResponse {
		s.lock.Lock()
		s.cachedResponse[msg.TransactionID] = &cachedResponse{time.Now(), msg}
		s.lock.Unlock()
	}
	return s.sendData(msg.Raw, fromaddr, toaddr)
}

/*
create channel etc...
*/
func (s *StunServerSock) sendStunMessageWithResult(msg *stun.Message, fromaddr, toaddr string) (key stun.TransactionID, ch chan *serverSockResponse, err error) {
	wait := make(chan *serverSockResponse)
	err = s.addWaiter(msg.TransactionID, wait)
	if err != nil {
		return
	}
	err = s.sendStunMessageAsync(msg, fromaddr, toaddr)
	if err != nil {
		return
	}
	ch = s.waiters[msg.TransactionID]
	return
}
func (s *StunServerSock) sendStunMessageSync(msg *stun.Message, fromaddr, toaddr string) (res *stun.Message, err error) {
	wait := make(chan *serverSockResponse)
	err = s.addWaiter(msg.TransactionID, wait)
	if err != nil {
		return
	}
	//defer s.getAndRemoveWaiter(msg.TransactionID)
	err = s.sendStunMessageAsync(msg, fromaddr, toaddr)
	if err != nil {
		return
	}
	return s.wait(wait)
}
func (s *StunServerSock) addWaiter(key stun.TransactionID, ch chan *serverSockResponse) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.waiters[key]; ok {
		return ErrDuplicateWaiter
	}
	s.waiters[key] = ch
	return nil
}
func (s *StunServerSock) getAndRemoveWaiter(key stun.TransactionID) (ch chan *serverSockResponse, exists bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	ch, exists = s.waiters[key]
	if exists {
		delete(s.waiters, key)
	}
	return
}
func (s *StunServerSock) wait(ch chan *serverSockResponse) (res *stun.Message, err error) {
	select {
	case res, ok := <-ch:
		if !ok {
			return nil, ErrWaiterClosed
		}
		return res.res, nil
	case <-time.After(s.syncMessageTimeout):
		return nil, ErrTimeout
	}
}

/*
根据需要发生了 channel binding 以后,需要指定 channel number, 这样才知道收到了来自哪里的消息.
*/
func (s *StunServerSock) SetChannelNumber(channelNumber int, addr string) {
	//todo fixit ,need a lock?
	s.channelNumber2Address[channelNumber] = addr
	s.address2ChannelNumber[addr] = channelNumber
}
func (s *StunServerSock) FinishNegotiation(mode serverSockMode) {
	s.log.Trace(fmt.Sprintf("change mode from %d to %d", s.mode, mode))
	s.mode = mode
}

// Serve reads packets from connections and responds to BINDING requests.
func (s *StunServerSock) Serve(c net.PacketConn) error {
	go func() {
		//writeto 是阻塞函数,不要阻塞 sendasync
		for {
			select {
			case r, ok := <-s.sendchan:
				if !ok {
					return
				}
				s.c.WriteTo(r.data, r.to)
			}
		}
	}()
	for {
		req := new(stun.Message)
		if err := s.serveConn(c, req); err != nil {
			s.log.Info(fmt.Sprintf("serve: %v", err))
			return err
		}
	}
}
func (s *StunServerSock) Close() {
	s.log.Trace(fmt.Sprintf("%s closed", s.Addr))
	s.stoped = true
	s.c.Close()
	for key, ch := range s.waiters {
		s.getAndRemoveWaiter(key)
		close(ch)
	}
	return
}

/*
监听指定的地址 bindAddr,
同时指定相关的用户密码密码等信息.
*/
func NewStunServerSock(bindAddr string, cb ServerSockCallbacker, name string) (s *StunServerSock, err error) {
	c, err := net.ListenPacket("udp", bindAddr)
	if err != nil {
		return
	}
	s = &StunServerSock{
		Addr:               bindAddr,
		mode:               StageNegotiation,
		c:                  c,
		waiters:            make(map[stun.TransactionID]chan *serverSockResponse),
		syncMessageTimeout: time.Second * 30,
		cb:                 cb,
		Name:               name,
		channelNumber2Address: make(map[int]string),
		address2ChannelNumber: make(map[string]int),
		cachedResponse:        make(map[stun.TransactionID]*cachedResponse),
		sendchan:              make(chan *sendreq, 10),
		log:                   log.New("name", fmt.Sprintf("%s-StunServerSock", name)),
	}
	go func() {
		s.Serve(s.c)
	}()
	return
}
