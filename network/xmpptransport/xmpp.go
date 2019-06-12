package xmpptransport

import (
	"errors"
	"time"

	"github.com/SmartMeshFoundation/Photon/network/wakeuphandler"

	"sync"

	"encoding/base64"
	"fmt"

	"strings"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mattn/go-xmpp"
)

var (
	errTimeout            = errors.New("timed out")
	errInvalidMessage     = errors.New("invalid message")
	errDuplicateWaiter    = errors.New("waiter with uid already exists")
	errWaiterClosed       = errors.New("waiter closed")
	errClientDisconnected = errors.New("client disconnected")
)

const (
	/*
			after offer send sdp,answer do the following jobs:
			1. receive sdp
			2. createicesteamtransport, contact with stun server,tunserver
			3. get it's own sdp
			4. send it's own sdp back to offer
		so how long should be better?
	*/
	defaultTimeout   = 15 * time.Second
	defaultReconnect = true
	nameSuffix       = "@mobileraiden"
	//TypeMobile photon run on a mobile device
	TypeMobile = "mobile"
	//TypeMeshBox photon run on a meshbox
	TypeMeshBox = "meshbox"
	//TypeOtherDevice photon run on a other device
	TypeOtherDevice = "other"
)

// Config contains various client options.
type Config struct {
	Timeout time.Duration
}

// DefaultConfig with standard private channel prefix and 1 second timeout.
var DefaultConfig = &Config{
	Timeout: defaultTimeout,
}

/*
PasswordGetter generate login password
*/
type PasswordGetter interface {
	//get current password
	GetPassWord() string
}

//DataHandler handels new data from peer node
type DataHandler interface {
	//DataHandler handle recevied data
	DataHandler(from common.Address, data []byte)
}

//NodeStatus is status of a photon node
type NodeStatus struct {
	IsOnline   bool
	DeviceType string
}

// XMPPConnection describes client connection to xmpp server.
type XMPPConnection struct {
	mutex          sync.RWMutex
	config         *Config
	options        xmpp.Options
	client         *xmpp.Client
	waitersMutex   sync.RWMutex
	waiters        map[string]chan interface{} //message waiting for response
	closed         chan struct{}
	reconnect      bool
	status         netshare.Status
	statusChan     chan<- netshare.Status
	NextPasswordFn PasswordGetter
	dataHandler    DataHandler
	name           string
	nodesStatus    map[string]*NodeStatus
	db             XMPPDb
	addrMap        map[common.Address]int //addr neighbor count
	*wakeuphandler.WakeUpHandler
}

/*
NewConnection create Xmpp connection to signal sever
*/
func NewConnection(ServerURL string, User common.Address, passwordFn PasswordGetter, dataHandler DataHandler, name, deviceType string, statusChan chan<- netshare.Status, dao XMPPDb) (x2 *XMPPConnection, err error) {
	x := &XMPPConnection{
		mutex:  sync.RWMutex{},
		config: DefaultConfig,
		options: xmpp.Options{
			Host:                         ServerURL,
			User:                         fmt.Sprintf("%s%s", strings.ToLower(User.String()), nameSuffix),
			Password:                     passwordFn.GetPassWord(),
			NoTLS:                        true,
			InsecureAllowUnencryptedAuth: true,
			Debug:                        false,
			Session:                      false,
			Status:                       "xa",
			StatusMessage:                name,
			Resource:                     deviceType,
		},
		client:         nil,
		waitersMutex:   sync.RWMutex{},
		waiters:        make(map[string]chan interface{}),
		nodesStatus:    make(map[string]*NodeStatus),
		closed:         make(chan struct{}),
		addrMap:        make(map[common.Address]int),
		reconnect:      true,
		status:         netshare.Disconnected,
		statusChan:     statusChan,
		NextPasswordFn: passwordFn,
		dataHandler:    dataHandler,
		name:           name,
		WakeUpHandler:  wakeuphandler.NewWakeupHandler("xmpp"),
		db:             dao,
	}
	log.Trace(fmt.Sprintf("%s new xmpp user %s password %s", name, User.String(), x.options.Password))
	x.client, err = x.options.NewClient()
	if err != nil {
		log.Error(fmt.Sprintf("%s new xmpp client err %s", name, err))
		return
	}
	x.changeStatus(netshare.Connected)
	x.registerChannelStateCallback() // 注册回调
	go x.loop()
	x2 = x
	return
}
func (x *XMPPConnection) loop() {
	defer rpanic.PanicRecover("xmpp")
	// 第一次启动时订阅通道partner的在线状态
	x.CollectNeighbors()
	for {
		chat, err := x.client.Recv()
		if x.status == netshare.Closed {
			return
		}
		if err != nil {
			log.Error(fmt.Sprintf("%s receive error %s ,try to reconnect ", x.name, err))
			err = x.client.Close()
			if err != nil {
				log.Error(fmt.Sprintf("xmpp close err %s", err))
			}
			x.reConnect()
			continue
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			if v.Type != "chat" {
				continue //error
			}
			remoteUser := strings.Split(v.Remote, "/")[0]
			remoteUser = strings.Split(remoteUser, "@")[0]
			raddr := common.HexToAddress(remoteUser)
			log.Trace(fmt.Sprintf("receive from %s message %s", utils.APex2(raddr), v.Text))
			data, err := base64.StdEncoding.DecodeString(v.Text)
			if err != nil {
				log.Error(fmt.Sprintf("receive unkown message :\n%s", utils.StringInterface(v, 3)))
			} else {
				x.dataHandler.DataHandler(raddr, data)
			}
		case xmpp.IQ:
			log.Info(fmt.Sprintf("receive unkonwn iq message :\n%s", utils.StringInterface(v, 3)))
		case xmpp.Presence:
			if len(v.ID) > 0 {
				// 订阅/取消订阅请求的应答消息,不用做处理
			} else {
				// 节点在线状态变更
				var id, device string
				ss := strings.Split(v.From, "/")
				if len(ss) >= 2 {
					device = ss[1]
				}
				id = ss[0]
				oldBs := x.nodesStatus[id]
				bs := &NodeStatus{
					DeviceType: device,
					IsOnline:   len(v.Type) == 0,
				}
				if bs.IsOnline && len(bs.DeviceType) == 0 {
					log.Error(fmt.Sprintf("receive unexpected presence message :\n%s", utils.StringInterface(v, 3)))
				}
				if bs.IsOnline && (oldBs == nil || !oldBs.IsOnline) {
					/*
						上线唤醒
					*/
					ss := strings.Split(id, "@")
					x.WakeUp(common.HexToAddress(ss[0]))
				}
				x.nodesStatus[id] = bs
				//log.Trace(fmt.Sprintf("node status change %s, deviceType=%s,isonline=%v", id, bs.DeviceType, bs.IsOnline))
			}
		default:
			log.Trace(fmt.Sprintf("recv unknown type message :\n %s", utils.StringInterface(v, 3)))
		}
	}
}
func (x *XMPPConnection) changeStatus(newStatus netshare.Status) {
	log.Info(fmt.Sprintf("xmpp changeStatus from %d to %d", x.status, newStatus))
	x.status = newStatus
	select {
	case x.statusChan <- newStatus:
	default:
		//never block
	}
}

//Reconnect :
func (x *XMPPConnection) Reconnect() {
	err := x.client.Close()
	if err != nil {
		log.Warn(fmt.Sprintf("xmpp client close err %s", err))
	}
	return
}

func (x *XMPPConnection) reConnect() {
	x.changeStatus(netshare.Reconnecting)
	o := x.options
	for {
		if x.status == netshare.Closed {
			return
		}
		o.Password = x.NextPasswordFn.GetPassWord()
		client, err := o.NewClient()
		if err != nil {
			log.Error(fmt.Sprintf("%s xmpp reconnect error %s", x.name, err))
			time.Sleep(time.Second)
			continue
		}
		x.mutex.Lock()
		x.client = client
		x.mutex.Unlock()
		break
	}
	// 重连成功过后,订阅所有partner的在线状态
	x.CollectNeighbors()
	x.changeStatus(netshare.Connected)
}
func (x *XMPPConnection) sendSyncIQ(msg *xmpp.IQ) (response *xmpp.IQ, err error) {
	uid := msg.ID
	wait := make(chan interface{})
	err = x.addWaiter(uid, wait)
	if err != nil {
		return nil, err
	}
	defer x.removeWaiter(uid)
	err = x.sendIQ(msg)
	if err != nil {
		return nil, err
	}
	r, err := x.wait(wait)
	response, ok := r.(*xmpp.IQ)
	if !ok {
		log.Error(fmt.Sprintf("recevie response %s,but type error ", utils.StringInterface(r, 3)))
		err = errors.New("type error")
	}
	return
}
func (x *XMPPConnection) send(msg *xmpp.Chat) error {
	select {
	case <-x.closed:
		return errClientDisconnected
	default:
		x.mutex.Lock()
		cli := x.client
		x.mutex.Unlock()
		//log.Trace(fmt.Sprintf("%s send msg %s:%s %s", x.name, msg.Remote, msg.Subject, msg.Text))
		_, err := cli.Send(*msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (x *XMPPConnection) sendIQ(msg *xmpp.IQ) error {
	select {
	case <-x.closed:
		return errClientDisconnected
	default:
		x.mutex.Lock()
		cli := x.client
		x.mutex.Unlock()
		log.Trace(fmt.Sprintf("%s send msg %s:%s %s", x.name, msg.From, msg.To, msg.ID))
		_, err := cli.SendIQ(*msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (x *XMPPConnection) addWaiter(uid string, ch chan interface{}) error {
	x.waitersMutex.Lock()
	defer x.waitersMutex.Unlock()
	if _, ok := x.waiters[uid]; ok {
		return errDuplicateWaiter
	}
	x.waiters[uid] = ch
	return nil
}

func (x *XMPPConnection) removeWaiter(uid string) error {
	x.waitersMutex.Lock()
	defer x.waitersMutex.Unlock()
	delete(x.waiters, uid)
	return nil
}

func (x *XMPPConnection) wait(ch chan interface{}) (response interface{}, err error) {
	select {
	case data, ok := <-ch:
		if !ok {
			return nil, errWaiterClosed
		}
		return data, nil
	case <-time.After(x.config.Timeout):
		return nil, errTimeout
	case <-x.closed:
		return nil, errClientDisconnected
	}
}

//Close this connection
func (x *XMPPConnection) Close() {
	x.changeStatus(netshare.Closed)
	close(x.closed)
	err := x.client.Close()
	if err != nil {
		log.Error(fmt.Sprintf("close err %s", err))
	}
}

//Connected returns true when this connection is ready for sent
func (x *XMPPConnection) Connected() bool {
	return x.status == netshare.Connected
}

//SendData to peer
func (x *XMPPConnection) SendData(addr common.Address, data []byte) error {
	chat := &xmpp.Chat{
		Remote: fmt.Sprintf("%s%s", strings.ToLower(addr.String()), nameSuffix),
		Type:   "chat",
		Stamp:  time.Now(),
	}
	chat.Text = base64.StdEncoding.EncodeToString(data)
	return x.send(chat)
}

const (
	resultOnline  = "pong"
	resultOffline = "pang"
)

type iqResult struct {
	Result   string
	Resource string
}

//IsNodeOnline test node is online
func (x *XMPPConnection) IsNodeOnline(addr common.Address) (deviceType string, isOnline bool, err error) {
	id := fmt.Sprintf("%s%s", strings.ToLower(addr.String()), nameSuffix)
	defer log.Trace(fmt.Sprintf("query nodeonline on xmpp %s %v", strings.ToLower(addr.String()), isOnline))
	ns, ok := x.nodesStatus[id]
	if ok {
		return ns.DeviceType, ns.IsOnline, nil
	}
	//log.Info(fmt.Sprintf("try to get status of %s, but no record", utils.APex2(addr)))
	return "", false, nil

}
func (x *XMPPConnection) sendPresence(msg *xmpp.Presence) error {
	select {
	case <-x.closed:
		return errClientDisconnected
	default:
		x.mutex.Lock()
		cli := x.client
		x.mutex.Unlock()
		//log.Trace(fmt.Sprintf("%s send msg %s:%s %s", x.name, msg.From, msg.To, msg.ID))
		_, err := cli.SendPresence(*msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (x *XMPPConnection) sendSyncPresence(msg *xmpp.Presence) (response *xmpp.Presence, err error) {
	uid := msg.ID
	wait := make(chan interface{})
	err = x.addWaiter(uid, wait)
	if err != nil {
		return nil, err
	}
	defer x.removeWaiter(uid)
	err = x.sendPresence(msg)
	if err != nil {
		return nil, err
	}
	r, err := x.wait(wait)
	if err != nil {
		return
	}
	response, ok := r.(*xmpp.Presence)
	if !ok {
		log.Error(fmt.Sprintf("recevie response %s,but type error ", utils.StringInterface(r, 3)))
		err = errors.New("type error")
	}
	return
}

//SubscribeNeighbour the status change of `addr`
// 这里修改为发送订阅消息即返回,不关心订阅成功与否,因为连接正常但订阅失败的概率较小,如果需要解决,考虑做有限次数的重试
func (x *XMPPConnection) SubscribeNeighbour(addr common.Address) error {
	addrName := fmt.Sprintf("%s%s", strings.ToLower(addr.String()), nameSuffix)
	p := xmpp.Presence{
		From: x.options.User,
		To:   addrName,
		Type: "subscribe",
		ID:   utils.RandomString(10), // 这里带上ID只是为了在loop()中区分该消息的应答消息与节点在线状态通知消息
	}
	return x.sendPresence(&p)
}

//Unsubscribe the status change of `addr`
/*
```xml
<presence id='xk3h1v69' to='leon@mobilephoton' type='unsubscribe'/>
```
*/
func (x *XMPPConnection) Unsubscribe(addr common.Address) error {
	addrName := fmt.Sprintf("%s%s", strings.ToLower(addr.String()), nameSuffix)
	p := xmpp.Presence{
		From: x.options.User,
		To:   addrName,
		Type: "unsubscribe",
		ID:   utils.RandomString(10), // 这里带上ID只是为了在loop()中区分该消息的应答消息与节点在线状态通知消息
	}
	return x.sendPresence(&p)
}

//SubscribeNeighbors I want to know these `addrs` status change
func (x *XMPPConnection) SubscribeNeighbors(addrs []common.Address) error {
	for _, addr := range addrs {
		err := x.SubscribeNeighbour(addr)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
XMPPDb 解耦 db 依赖
*/
type XMPPDb interface {
	XMPPIsAddrSubed(addr common.Address) bool
	XMPPMarkAddrSubed(addr common.Address)
	GetChannelList(token, partner common.Address) (cs []*channeltype.Serialization, err error)
	RegisterNewChannelCallback(f cb.ChannelCb)
	RegisterChannelStateCallback(f cb.ChannelCb)
	RegisterChannelSettleCallback(f cb.ChannelCb)
	XMPPUnMarkAddr(addr common.Address)
}

//CollectNeighbors subscribe status change from database
func (x *XMPPConnection) CollectNeighbors() {
	db := x.db
	cs, err := db.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		log.Error(fmt.Sprintf("db GetChannelList err %s", err.Error()))
	}
	for _, c := range cs {
		if c.State == channeltype.StateOpened {
			x.addrMap[c.PartnerAddress()]++
		}
	}
	/*
	   异步,一次性订阅所有参数通道列表里的所有地址,这里没有并发问题,只是可能重复,由于不记录订阅状态,重复订阅的成本很低
	*/
	go func(cs []*channeltype.Serialization) {
		for _, c := range cs {
			err = x.SubscribeNeighbour(c.PartnerAddress())
			if err != nil {
				log.Error(fmt.Sprintf("Subscribe partner %s err : %s", c.PartnerAddress().String(), err.Error()))
			}
		}
		log.Info(fmt.Sprintf("SubscribeNeighbour of %d channels", len(cs)))
	}(cs)
	return
}

/*
	向数据库注册回调,用于在通道创建/关闭时订阅/取消订阅用户在线状态
*/
func (x *XMPPConnection) registerChannelStateCallback() {
	db := x.db
	db.RegisterNewChannelCallback(func(c *channeltype.Serialization) (remove bool) {
		if x.status == netshare.Closed {
			return true
		}
		err := x.SubscribeNeighbour(c.PartnerAddress())
		if err != nil {
			log.Error(fmt.Sprintf("sub %s err %s", c.PartnerAddress().String(), err))
		} else {
			x.db.XMPPMarkAddrSubed(c.PartnerAddress())
		}
		return false
	})
	db.RegisterChannelStateCallback(func(c *channeltype.Serialization) (remove bool) {
		if x.status == netshare.Closed {
			return true
		}
		return false
	})
	db.RegisterChannelSettleCallback(func(c *channeltype.Serialization) (remove bool) {
		if x.status == netshare.Closed {
			return true
		}
		x.addrMap[c.PartnerAddress()]--
		if x.addrMap[c.PartnerAddress()] <= 0 {
			err := x.Unsubscribe(c.PartnerAddress())
			if err != nil {
				log.Error(fmt.Sprintf("unsub %s err %s", c.PartnerAddress().String(), err))
				return false
			}
			db.XMPPUnMarkAddr(c.PartnerAddress())
		}
		return false
	})
}
