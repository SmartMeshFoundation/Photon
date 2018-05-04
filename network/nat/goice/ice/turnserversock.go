package ice

import (
	"net"
	"strconv"

	"fmt"

	"time"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/turn"
)

const StunKeepAliveInterval = time.Second * 20

type TurnServerSockConfig struct {
	user         string //turn server user
	password     string //turn server password
	nonce        string
	realm        string
	credentials  stun.MessageIntegrity //long term
	lifetime     turn.Lifetime         //create permission life time.
	relayAddress string
	serverAddr   string
}
type TurnServerSock struct {
	s           *StunServerSock
	cfg         *TurnServerSockConfig
	cb          ServerSockCallbacker
	allowedPeer []turn.PeerAddress
	Name        string
	stopchan    chan struct{} //for stop refresh.
	log         log.Logger
}

func NewTurnServerSockWrapper(bindAddr, name string, cb ServerSockCallbacker, cfg *TurnServerSockConfig) (ts *TurnServerSock, err error) {
	ts = &TurnServerSock{
		cfg:      cfg,
		cb:       cb,
		Name:     name,
		stopchan: make(chan struct{}),
		log:      log.New("name", fmt.Sprintf("%s-TurnServerSock", name)),
	}
	s, err := NewStunServerSock(bindAddr, ts, name)
	if err != nil {
		return
	}
	ts.s = s
	return
}

/*
 收到一个 stun.Message, 可能是 Bind Request/Bind Response 等等.
*/
func (ts *TurnServerSock) RecieveStunMessage(localAddr, remoteAddr string, msg *stun.Message) {
	/*
		需要在协商阶段处理 turn server 中转来的 Data Indication.将其解码,然后把其中的 binding response 交给调用者.
	*/
	if msg.Type == turn.DataIndication {
		var data turn.Data
		var peer turn.PeerAddress
		if remoteAddr != ts.cfg.serverAddr {
			panic("data indication from unkown address")
		}
		err := data.GetFrom(msg)
		if err != nil {
			//todo fix all panic shoulde be removed ,attacker...
			panic(fmt.Sprintf("unexpected message.. %s", msg))
		}
		if len(data) <= 0 {
			panic(fmt.Sprintf("unexpected message.. %s", msg))
		}
		err = peer.GetFrom(msg)
		if err != nil {
			panic(fmt.Sprintf("unexpected message.. %s", msg))
		}
		res := new(stun.Message)
		_, err = res.Write([]byte(data))
		if err != nil || res.Type.Method == stun.MethodChannelData {
			//有可能我认为协商没完成,但是对方认为已经完成了,所以直接发送了数据过来.但是我还没有进行 channel binding. 所以还是要处理数据的.
			if ts.cb != nil {
				ts.cb.ReceiveData(localAddr, peer.String(), []byte(data))
			}
		} else {
			ts.log.Trace(fmt.Sprintf("actual message:%s", res))
			if res.Type == stun.BindingSuccess || res.Type != stun.BindingError || res.Type != stun.BindingRequest {
				ts.s.stunMessageReceived(ts.cfg.relayAddress, peer.String(), res)
			} else {
				panic("data indication must carry bind response")
			}
		}
		return
	}
	if ts.cb != nil {
		ts.cb.RecieveStunMessage(localAddr, remoteAddr, msg)
	}
}

/*
	ICE 协商建立连接以后,收到了对方发过来的数据,可能是经过 turn server 中转的 channel data( 不接受 sendData data request),也可能直接是数据.
	如果是经过 turn server 中转的, channelNumber 一定介于0x4000-0x7fff 之间.否则一定为0
*/
func (ts *TurnServerSock) ReceiveData(localAddr, peerAddr string, data []byte) {
	msg2 := new(stun.Message)
	_, err := msg2.Write(data)
	if err == nil && msg2.Type.Method != stun.MethodChannelData {
		//收到了发到中转地址的一个 stun message
		ts.s.stunMessageReceived(ts.cfg.relayAddress, peerAddr, msg2)
		return
	}
	if ts.cb != nil {
		ts.cb.ReceiveData(localAddr, peerAddr, data)
	}
}

func (ts *TurnServerSock) createPermission(remoteCandidates []*Candidate) (res *stun.Message, err error) {
	var req *stun.Message
	var peers []stun.Setter
	for _, c := range remoteCandidates {
		host, port, err := net.SplitHostPort(c.addr)
		if err != nil {
			//panic?
			ts.log.Error(fmt.Sprintf("split error for %s,err:%s", c.addr, err))
			continue
		}
		peer := turn.PeerAddress{
			IP: net.ParseIP(host),
		}
		peer.Port, _ = strconv.Atoi(port)
		peers = append(peers, peer)
	}
	req = new(stun.Message)
	err = req.Build(stun.TransactionIDSetter, turn.CreatePermissionRequest,
		stun.Realm(ts.cfg.realm), stun.Nonce(ts.cfg.nonce),
		stun.Username(ts.cfg.user),
	)
	if err != nil {
		ts.log.Error(fmt.Sprintf("build err %s", err))
	}
	for _, p := range peers {
		err = p.AddTo(req)
		if err != nil {
			ts.log.Error(fmt.Sprintf("build err %s", err))
		}
	}
	err = ts.cfg.credentials.AddTo(req)
	if err != nil {
		ts.log.Error(fmt.Sprintf("build err %s", err))
	}
	err = stun.Fingerprint.AddTo(req)
	if err != nil {
		ts.log.Error(fmt.Sprintf("build err %s", err))
	}
	res, err = ts.s.sendStunMessageSync(req, ts.s.Addr, ts.cfg.serverAddr)
	return
}

func (ts *TurnServerSock) wrapperData(fromaddr string, toaddr string, data []byte) (outdata []byte, toaddr2 string) {
	return nil, ""
}
func (ts *TurnServerSock) wrapperStunMessage(fromaddr string, toaddr string, msg *stun.Message) (msg2 *stun.Message, fromaddr2, toaddr2 string) {
	if fromaddr == ts.s.Addr {
		return msg, fromaddr, toaddr
	}
	if fromaddr != ts.cfg.relayAddress {
		panic(fmt.Sprintf("sendData from unkonw address.. ts.s.Addr=%s,fromaddr=%s,relay=%s", ts.s.Addr, fromaddr, ts.cfg.relayAddress))
	}
	msg2 = new(stun.Message)
	to := addrToUdpAddr(toaddr)
	peer := &turn.PeerAddress{
		IP:   to.IP,
		Port: to.Port,
	}
	msg2.Build(stun.TransactionIDSetter,
		turn.SendIndication,
		peer, turn.Data(msg.Raw), stun.Fingerprint,
	)
	return msg2, ts.s.Addr, ts.cfg.serverAddr
}
func (ts *TurnServerSock) sendStunMessageAsync(msg *stun.Message, fromaddr, toaddr string) error {
	ts.log.Trace(fmt.Sprintf("---sendData stun message %s-->%s ---\n%s\n", fromaddr, toaddr, msg))
	msg2, fromaddr2, toaddr2 := ts.wrapperStunMessage(fromaddr, toaddr, msg)
	if fromaddr2 != fromaddr {
		ts.log.Trace(fmt.Sprintf("message actually from %s to %s", fromaddr2, toaddr2))
	}
	return ts.s.sendStunMessageAsync(msg2, fromaddr2, toaddr2) // sendData(msg2.Raw, fromaddr2, toaddr2)
}

/*
create channel etc...
*/
func (ts *TurnServerSock) sendStunMessageWithResult(msg *stun.Message, fromaddr, toaddr string) (key stun.TransactionID, ch chan *serverSockResponse, err error) {
	wait := make(chan *serverSockResponse)
	err = ts.s.addWaiter(msg.TransactionID, wait)
	if err != nil {
		return
	}
	err = ts.sendStunMessageAsync(msg, fromaddr, toaddr)
	if err != nil {
		return
	}
	ch = ts.s.waiters[msg.TransactionID]
	return
}
func (ts *TurnServerSock) sendStunMessageSync(msg *stun.Message, fromaddr, toaddr string) (res *stun.Message, err error) {
	wait := make(chan *serverSockResponse)
	err = ts.s.addWaiter(msg.TransactionID, wait)
	if err != nil {
		return
	}
	//defer ts.s.getAndRemoveWaiter(msg.TransactionID)
	msg2, fromaddr2, toaddr2 := ts.wrapperStunMessage(fromaddr, toaddr, msg)
	err = ts.s.sendStunMessageAsync(msg2, fromaddr2, toaddr2)
	if err != nil {
		return
	}
	return ts.s.wait(wait)
}
func (ts *TurnServerSock) Close() {
	close(ts.stopchan)
	ts.s.Close()
}

/*
这个连接上有三种情况
1.直接通信
2.通过 stun server 反馈到的 地址通信
3.通过 turn 中转.
*/
func (ts *TurnServerSock) StartRefresh() {
	go func() {
		for {
			ts.keepAlive()
			select {
			case <-time.After(StunKeepAliveInterval):
				continue
			case <-ts.stopchan:
				return
			}
		}
	}()
	if ts.s.mode == TurnModeData {
		go func() {
			for {
				ts.refreshRequest(ts.cfg.lifetime)
				select {
				case <-time.After(ts.cfg.lifetime.Duration / 2):
					continue
				case <-ts.stopchan:
					return
				}
			}
		}()
	} else {
		//stop turn's allocate right now
		ts.log.Debug(fmt.Sprintf("release turn allocated ."))
		ts.refreshRequest(turn.Lifetime{})
	}

}
func (ts *TurnServerSock) sendData(data []byte, fromaddr, toaddr string) error {
	if fromaddr == ts.cfg.relayAddress {
		/*
			分成两个阶段,第一阶段协商完毕可以发送数据,但是 check 仍在继续,发送链接随时可能变化.
			第二阶段: 协商完毕,我这边的已经稳定下来了,那么这时候就应该通过 channel 来发送数据.
		*/
		number := ts.s.address2ChannelNumber[toaddr]
		if number >= turn.MinChannelNumber && number <= turn.MaxChannelNumber {
			wdata := &turn.ChannelData{
				ChannelNumber: uint16(number),
				Data:          data,
			}
			r, _ := stun.Build(turn.ChannelDataRequest)
			err := wdata.AddTo(r)
			if err != nil {
				panic("turn.channeldata error")
			}
			ts.log.Trace(fmt.Sprintf("send  channel data %d, %s---->%s", len(r.Raw), ts.s.Addr, ts.cfg.serverAddr))
			ts.s.sendData(r.Raw, ts.s.Addr, ts.cfg.serverAddr)
		} else {
			if ts.s.mode == TurnModeData {
				ts.log.Warn(fmt.Sprintf("should not happen only if channel binding fail"))
			}
			to := addrToUdpAddr(toaddr)
			peer := turn.PeerAddress{
				IP:   to.IP,
				Port: to.Port,
			}
			r, err := stun.Build(stun.TransactionIDSetter, turn.SendIndication, turn.Data(data), peer, stun.Fingerprint)
			if err != nil {
				panic("build error")
			}
			ts.log.Trace(fmt.Sprintf("send data use send indication %s--->%s  message:%s\n", ts.s.Addr, ts.cfg.serverAddr, r))
			ts.s.sendStunMessageAsync(r, ts.s.Addr, ts.cfg.serverAddr)
		}
	} else {
		ts.log.Trace(fmt.Sprintf("send directly data %d   %s----->%s", len(data), fromaddr, toaddr))
		return ts.s.sendData(data, fromaddr, toaddr)
	}
	return nil
}

/*
绑定到 channel, 节省流量.
*/
func (ts *TurnServerSock) channelBind(addr string) error {
	uaddr := addrToUdpAddr(addr)
	peerAddr := &turn.PeerAddress{
		IP:   uaddr.IP,
		Port: uaddr.Port,
	}
	req, err := stun.Build(stun.TransactionIDSetter,
		turn.ChannelBindRequest,
		turn.ChannelNumber(turn.MinChannelNumber),
		peerAddr,
		stun.Username(ts.cfg.user),
		stun.Realm(ts.cfg.realm),
		stun.Nonce(ts.cfg.nonce),
		ts.cfg.credentials,
	)
	if err != nil {
		panic("....")
	}
	res, err := ts.s.sendStunMessageSync(req, ts.s.Addr, ts.cfg.serverAddr)
	if err != nil {
		return err
	}
	if res.Type.Method != stun.MethodChannelBind && res.Type.Class != stun.ClassSuccessResponse {
		ts.log.Error(fmt.Sprintf("channel bind response :%s", res))
		return errors.New("channel bind error")
	}
	ts.s.SetChannelNumber(turn.MinChannelNumber, addr)
	return nil
}

/*
我这边认为协商成功了,但是对方可能还灭与偶成功,所以仍然可能收到 stun message 消息,也就是通过 channel data 收到的还有可能是 stun 消息而不是真实的数据
*/
func (ts *TurnServerSock) FinishNegotiation(mode serverSockMode) {
	ts.log.Trace(fmt.Sprintf("change mode from %d to %d", ts.s.mode, mode))
	ts.s.mode = mode
	ts.StartRefresh()
}
func (ts *TurnServerSock) refreshRequest(lifetime turn.Lifetime) {
	req, err := stun.Build(stun.TransactionIDSetter,
		turn.RefreshRequest,
		stun.Username(ts.cfg.user),
		stun.Realm(ts.cfg.realm),
		stun.Nonce(ts.cfg.nonce),
		lifetime,
		ts.cfg.credentials,
	)
	if err != nil {
		panic("....")
	}
	res, err := ts.s.sendStunMessageSync(req, ts.s.Addr, ts.cfg.serverAddr)
	if err != nil {
		ts.log.Error(fmt.Sprintf("refresh request error %s", err))
		return
	}
	if res.Type != turn.RefreshResponse {
		//must refresh error response
		var code stun.ErrorCodeAttribute
		err = code.GetFrom(res)
		if err != nil {
			ts.log.Error("i don't know why?..")
		}
		ts.log.Error(fmt.Sprintf("%s channel refresh response  err:%s", ts.Name, code))
	} else {
		err = lifetime.GetFrom(res)
		if err != nil {
			ts.log.Error(fmt.Sprintf("unexpected err :%s", err))
		} else {
			ts.cfg.lifetime = lifetime
		}
	}
}

/*
only keep server reflexive port is valid.
keep the allocate address valid ,should call refersh request.
*/
func (ts *TurnServerSock) keepAlive() {
	req, _ := stun.Build(stun.TransactionIDSetter, stun.BindingIndication)
	ts.s.sendStunMessageAsync(req, ts.s.Addr, ts.cfg.serverAddr)
}

type ServerSocker interface {
	sendStunMessageSync(msg *stun.Message, fromaddr, toaddr string) (res *stun.Message, err error)
	sendStunMessageWithResult(msg *stun.Message, fromaddr, toaddr string) (key stun.TransactionID, ch chan *serverSockResponse, err error)
	sendStunMessageAsync(msg *stun.Message, fromaddr, toaddr string) error
	sendData(data []byte, fromaddr, toaddr string) error
	Close()
	FinishNegotiation(mode serverSockMode)
	//StartRefresh()
}
