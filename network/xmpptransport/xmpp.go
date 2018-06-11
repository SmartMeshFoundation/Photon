package xmpptransport

import (
	"errors"
	"time"

	"sync"

	"encoding/base64"
	"fmt"

	"strings"

	"encoding/json"

	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
	//TypeMobile raiden run on a mobile device
	TypeMobile = "mobile"
	//TypeMeshBox raiden run on a meshbox
	TypeMeshBox = "meshbox"
	//TypeOtherDevice raiden run on a other device
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

// Status shows actual connection status.
type Status int

const (
	//Disconnected init status
	Disconnected = Status(iota)
	//Connected connection status
	Connected
	//Closed user closed
	Closed
	//Reconnecting connection error
	Reconnecting
)

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

// XMPPConnection describes client connection to xmpp server.
type XMPPConnection struct {
	mutex          sync.RWMutex
	config         *Config
	options        xmpp.Options
	client         *xmpp.Client
	waitersMutex   sync.RWMutex
	waiters        map[string]chan *xmpp.IQ
	closed         chan struct{}
	reconnect      bool
	status         Status
	statusChan     chan<- Status
	NextPasswordFn PasswordGetter
	dataHandler    DataHandler
	name           string
}

/*
NewConnection create Xmpp connection to signal sever
*/
func NewConnection(ServerURL string, User common.Address, passwordFn PasswordGetter, dataHandler DataHandler, name, deviceType string, statusChan chan<- Status) (x2 *XMPPConnection, err error) {
	x := &XMPPConnection{
		mutex:  sync.RWMutex{},
		config: DefaultConfig,
		options: xmpp.Options{
			Host:                         ServerURL,
			User:                         fmt.Sprintf("%s%s", User.String(), nameSuffix),
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
		waiters:        make(map[string]chan *xmpp.IQ),
		closed:         make(chan struct{}),
		reconnect:      true,
		status:         Disconnected,
		statusChan:     statusChan,
		NextPasswordFn: passwordFn,
		dataHandler:    dataHandler,
		name:           name,
	}
	log.Trace(fmt.Sprintf("%s new xmpp user %s password %s", name, User.String(), x.options.Password))
	x.client, err = x.options.NewClient()
	if err != nil {
		log.Trace(fmt.Sprintf("%s new xmpp client err %s", name, err))
		return
	}
	x.changeStatus(Connected)
	go x.loop()
	x2 = x
	return
}
func (x *XMPPConnection) loop() {
	defer rpanic.PanicRecover("xmpp")
	for {
		chat, err := x.client.Recv()
		if x.status == Closed {
			return
		}
		if err != nil {
			//todo how to detect network error ,disconnect
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
			data, err := base64.StdEncoding.DecodeString(v.Text)
			if err != nil {
				log.Error(fmt.Sprintf("receive unkown message %s", utils.StringInterface(v, 3)))
			} else {
				x.dataHandler.DataHandler(raddr, data)
			}
		case xmpp.IQ:

			uid := v.ID
			x.waitersMutex.Lock()
			ch, ok := x.waiters[uid]
			x.waitersMutex.Unlock()
			if ok {
				log.Trace(fmt.Sprintf("%s %s received response", x.name, uid))
				ch <- &v
			} else {
				log.Info(fmt.Sprintf("receive unkonwn iq message %s", utils.StringInterface(v, 3)))
			}
		default:
			//log.Trace(fmt.Sprintf("recv %s", utils.StringInterface(v, 3)))
		}
	}
}
func (x *XMPPConnection) changeStatus(newStatus Status) {
	log.Info(fmt.Sprintf("changeStatus from %d to %d", x.status, newStatus))
	x.status = newStatus
	select {
	case x.statusChan <- newStatus:
	default:
		//never block
	}
}
func (x *XMPPConnection) reConnect() {
	x.changeStatus(Reconnecting)
	o := x.options
	for {
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
	x.changeStatus(Connected)
}
func (x *XMPPConnection) sendSyncIQ(msg *xmpp.IQ) (response *xmpp.IQ, err error) {
	uid := msg.ID
	wait := make(chan *xmpp.IQ)
	err = x.addWaiter(uid, wait)
	if err != nil {
		return nil, err
	}
	defer x.removeWaiter(uid)
	err = x.sendIQ(msg)
	if err != nil {
		return nil, err
	}
	response, err = x.wait(wait)
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
		log.Trace(fmt.Sprintf("%s send msg %s:%s %s", x.name, msg.Remote, msg.Subject, msg.Text))
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
func (x *XMPPConnection) addWaiter(uid string, ch chan *xmpp.IQ) error {
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

func (x *XMPPConnection) wait(ch chan *xmpp.IQ) (response *xmpp.IQ, err error) {
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
	x.changeStatus(Closed)
	close(x.closed)
	err := x.client.Close()
	if err != nil {
		log.Error(fmt.Sprintf("close err %s", err))
	}
}

//Connected returns true when this connection is ready for sent
func (x *XMPPConnection) Connected() bool {
	return x.status == Connected
}

//SendData to peer
func (x *XMPPConnection) SendData(addr common.Address, data []byte) error {
	chat := &xmpp.Chat{
		Remote: fmt.Sprintf("%s%s", addr.String(), nameSuffix),
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
	iq := &xmpp.IQ{
		ID:    utils.RandomString(10),
		From:  x.options.User,
		To:    fmt.Sprintf("%s%s", addr.String(), nameSuffix),
		Type:  "get",
		Query: []byte("<ping xmlns='urn:xmpp:ping'/>"),
	}
	r, err := x.sendSyncIQ(iq)
	if err != nil {
		err = fmt.Errorf("sendsynciq to %s err %s", utils.APex(addr), err)
		return
	}
	if len(r.Query) < 13 {
		err = fmt.Errorf("iq body too short")
		return
	}

	//log.Trace(fmt.Sprintf("query=%s", string(r.Query)))
	//log.Trace(fmt.Sprintf("body=%s", r.Query))
	var ir iqResult
	err = json.Unmarshal([]byte(r.Query), &ir)
	if err != nil {
		return
	}
	//log.Trace(fmt.Sprintf("ir=%s", utils.StringInterface(ir, 3)))
	isOnline = ir.Result == resultOnline
	deviceType = ir.Resource
	return
}
