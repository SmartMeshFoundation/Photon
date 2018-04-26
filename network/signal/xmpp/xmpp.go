package xmpp

import (
	"errors"
	"sync"
	"time"

	"fmt"

	"encoding/json"

	"strings"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/interface"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mattn/go-xmpp"
)

var (
	ErrTimeout            = errors.New("timed out")
	ErrInvalidMessage     = errors.New("invalid message")
	ErrDuplicateWaiter    = errors.New("waiter with uid already exists")
	ErrWaiterClosed       = errors.New("waiter closed")
	ErrClientDisconnected = errors.New("client disconnected")
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
	DefaultTimeout   = 15 * time.Second
	DefaultReconnect = true
	NameSuffix       = "@mobileraiden"
)

// Config contains various client options.
type Config struct {
	Timeout time.Duration
}

// DefaultConfig with standard private channel prefix and 1 second timeout.
var DefaultConfig = &Config{
	Timeout: DefaultTimeout,
}

// Status shows actual connection status.
type Status int

const (
	DISCONNECTED = Status(iota)
	CONNECTED
	CLOSED
	RECONNECTING
)

type GetCurrentPasswordFunc func() string

// XmppWrapper describes client connection to Centrifugo server.
type XmppWrapper struct {
	mutex          sync.RWMutex
	config         *Config
	options        xmpp.Options
	client         *xmpp.Client
	waitersMutex   sync.RWMutex
	waiters        map[string]chan *xmpp.Chat
	closed         chan struct{}
	reconnect      bool
	status         Status
	NextPasswordFn GetCurrentPasswordFunc
	SdpHandler     SignalInterface.SdpHandler
	name           string
}

func NewXmpp(ServerUrl string, User common.Address, passwordFn GetCurrentPasswordFunc, sdphandler SignalInterface.SdpHandler, name string) (sp SignalInterface.SignalProxy, err error) {
	x := &XmppWrapper{
		options: xmpp.Options{
			Host:     ServerUrl,
			User:     fmt.Sprintf("%s%s", User.String(), NameSuffix),
			Password: passwordFn(),
			NoTLS:    true,
			InsecureAllowUnencryptedAuth: true,
			Debug:         false,
			Session:       false,
			Status:        "xa",
			StatusMessage: "welcome",
		},
		config:         DefaultConfig,
		reconnect:      true,
		NextPasswordFn: passwordFn,
		status:         DISCONNECTED,
		SdpHandler:     sdphandler,
		waiters:        make(map[string]chan *xmpp.Chat),
		closed:         make(chan struct{}),
		name:           name,
	}
	log.Trace(fmt.Sprintf("%s new xmpp user %s password %s", name, User.String(), x.options.Password))
	x.client, err = x.options.NewClient()
	if err != nil {
		log.Trace(fmt.Sprintf("%s new xmpp client err %s", name, err))
		return
	}
	x.status = CONNECTED
	sp = x
	go func() {
		for {
			chat, err := x.client.Recv()
			if x.status == CLOSED {
				return
			}
			if err != nil {
				//todo how to detect network error ,disconnect
				log.Error(fmt.Sprintf("%s receive error %s ,try to reconnect ", name, err))
				x.client.Close()
				x.ReConnect()
				continue
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				if v.Type != "chat" {
					continue //error
				}
				log.Trace(fmt.Sprintf("%s receive message %s,%s", name, v.Remote, v.Text))
				remoteUser := strings.Split(v.Remote, "/")[0]
				v.Remote = remoteUser
				uid := fmt.Sprintf("%s-%s", v.Remote, v.Subject)
				x.waitersMutex.Lock()
				ch, ok := x.waiters[uid]
				x.waitersMutex.Unlock()
				if ok {
					log.Trace(fmt.Sprintf("%s %s received response", name, uid))
					ch <- &v
				} else {
					var cmd XmppCommand
					err := json.Unmarshal([]byte(v.Text), &cmd)
					if err != nil {
						log.Debug(fmt.Sprintf("%s recieve unkown message from:%s, subject:%s,text:%s", name, v.Remote, v.Subject, v.Text))
						continue //
					} else {
						x.HandleNewCommand(v.Remote, v.Subject, &cmd)
					}
				}
			}
		}
	}()
	return
}
func (x *XmppWrapper) ReConnect() {
	x.status = RECONNECTING
	o := x.options
	for {
		o.Password = x.NextPasswordFn()
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
	x.status = CONNECTED
}
func (x *XmppWrapper) sendSync(msg *xmpp.Chat) (response *xmpp.Chat, err error) {
	uid := fmt.Sprintf("%s-%s", msg.Remote, msg.Subject)
	wait := make(chan *xmpp.Chat)
	err = x.addWaiter(uid, wait)
	if err != nil {
		return nil, err
	}
	defer x.removeWaiter(uid)
	err = x.send(msg)
	if err != nil {
		return nil, err
	}
	return x.wait(wait)
}

func (x *XmppWrapper) send(msg *xmpp.Chat) error {
	select {
	case <-x.closed:
		return ErrClientDisconnected
	default:
		x.mutex.Lock()
		cli := x.client
		x.mutex.Unlock()
		log.Trace(fmt.Sprintf("%s send msg %s:%s %s", x.name, msg.Remote, msg.Subject, msg.Text))
		cli.Send(*msg)
	}
	return nil
}

func (x *XmppWrapper) addWaiter(uid string, ch chan *xmpp.Chat) error {
	x.waitersMutex.Lock()
	defer x.waitersMutex.Unlock()
	if _, ok := x.waiters[uid]; ok {
		return ErrDuplicateWaiter
	}
	x.waiters[uid] = ch
	return nil
}

func (x *XmppWrapper) removeWaiter(uid string) error {
	x.waitersMutex.Lock()
	defer x.waitersMutex.Unlock()
	delete(x.waiters, uid)
	return nil
}

func (x *XmppWrapper) wait(ch chan *xmpp.Chat) (response *xmpp.Chat, err error) {
	select {
	case data, ok := <-ch:
		if !ok {
			return nil, ErrWaiterClosed
		}
		return data, nil
	case <-time.After(x.config.Timeout):
		return nil, ErrTimeout
	case <-x.closed:
		return nil, ErrClientDisconnected
	}
}

var reachok string

func init() {
	var cmd = XmppCommand{
		Command: CommandTryReach,
	}
	data, _ := json.Marshal(&cmd)
	reachok = string(data)
}

func (x *XmppWrapper) HandleNewCommand(from, subject string, cmd *XmppCommand) {
	log.Trace(fmt.Sprintf("%s new command from:%s,subect:%s,cmd=%s", x.name, from, subject, utils.StringInterface1(cmd)))
	switch cmd.Command {
	case CommandTryReach:
		x.send(&xmpp.Chat{
			Type:    "chat",
			Remote:  from,
			Subject: subject,
			Text:    reachok,
		})
	case CommandExChangeSdp:
		fromaddr := strings.Split(from, "@")[0]
		r, err := x.SdpHandler(common.HexToAddress(fromaddr), cmd.OtherInfo)
		cmd2 := XmppCommand{
			Command: CommandExChangeSdp,
		}
		if err != nil {
			cmd2.Error = err.Error()
			log.Error(fmt.Sprintf("%s cannot handle sdp request cmd=%s,err=%s", x.name, spew.Sdump(cmd), err))
		} else {
			cmd2.OtherInfo = r
		}
		data, _ := json.Marshal(cmd2)
		x.send(&xmpp.Chat{
			Type:    "chat",
			Remote:  from,
			Subject: subject,
			Text:    string(data),
		})
	default:
		log.Error(fmt.Sprintf("%s receive unkown from:%s,subject:%s cmd:%s", x.name, from, subject, spew.Sdump(cmd)))
	}
}
func (x *XmppWrapper) Close() {
	x.status = CLOSED
	close(x.closed)
	x.client.Close()
}

func (x *XmppWrapper) Connected() bool {
	return x.status == CONNECTED
}

type Command int

const (
	CommandTryReach = Command(iota)
	CommandExChangeSdp
)

type XmppCommand struct {
	Command   Command
	Error     string
	OtherInfo string
}

func (x *XmppWrapper) sendCommand(addr common.Address, cmd *XmppCommand) (cmdResponse *XmppCommand, err error) {
	chat := &xmpp.Chat{
		Remote:  fmt.Sprintf("%s%s", addr.String(), NameSuffix),
		Type:    "chat",
		Subject: utils.RandomString(10),
		Stamp:   time.Now(),
	}
	data, err := json.Marshal(cmd)
	if err != nil {
		return
	}
	chat.Text = string(data)
	response, err := x.sendSync(chat)
	if err != nil {
		return
	}
	log.Trace(fmt.Sprintf("%s receive  response %s,%s", x.name, response.Remote, response.Text))
	var cmd2 XmppCommand
	err = json.Unmarshal([]byte(response.Text), &cmd2)
	if err != nil {
		return
	}
	if len(cmd2.Error) != 0 {
		err = errors.New(cmd2.Error)
		return
	}
	cmdResponse = &cmd2
	return
}
func (x *XmppWrapper) TryReach(addr common.Address) error {
	_, err := x.sendCommand(addr, &XmppCommand{
		Command: CommandTryReach,
	})
	return err
}

func (x *XmppWrapper) ExchangeSdp(addr common.Address, sdp string) (partnerSdp string, err error) {
	r, err := x.sendCommand(addr, &XmppCommand{
		Command:   CommandExChangeSdp,
		OtherInfo: sdp,
	})
	if err != nil {
		return "", err
	}
	return r.OtherInfo, nil
}
