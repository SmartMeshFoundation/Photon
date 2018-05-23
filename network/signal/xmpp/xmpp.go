package xmpp

import (
	"encoding/json"
	"errors"
	"time"

	"sync"

	"fmt"
	"strings"

	"github.com/SmartMeshFoundation/SmartRaiden/network/signal/interface"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mattn/go-xmpp"
	"github.com/nkbai/log"
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
	disconnected = Status(iota)
	connected
	closed
	reconnecting
)

/*
GetCurrentPasswordFunc generate login password
*/
type GetCurrentPasswordFunc func() string

// SignalConnection describes client connection to Centrifugo server.
type SignalConnection struct {
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

/*
NewSignalConnection create Xmpp connection to signal sever
*/
func NewSignalConnection(ServerURL string, User common.Address, passwordFn GetCurrentPasswordFunc, sdphandler SignalInterface.SdpHandler, name string) (sp SignalInterface.SignalProxy, err error) {
	x := &SignalConnection{
		options: xmpp.Options{
			Host:     ServerURL,
			User:     fmt.Sprintf("%s%s", User.String(), nameSuffix),
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
		status:         disconnected,
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
	x.status = connected
	sp = x
	go func() {
		for {
			chat, err := x.client.Recv()
			if x.status == closed {
				return
			}
			if err != nil {
				//todo how to detect network error ,disconnect
				log.Error(fmt.Sprintf("%s receive error %s ,try to reconnect ", name, err))
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
					var cmd xmppCommand
					err := json.Unmarshal([]byte(v.Text), &cmd)
					if err != nil {
						log.Debug(fmt.Sprintf("%s recieve unkown message from:%s, subject:%s,text:%s", name, v.Remote, v.Subject, v.Text))
						continue //
					} else {
						x.handleNewCommand(v.Remote, v.Subject, &cmd)
					}
				}
			}
		}
	}()
	return
}
func (x *SignalConnection) reConnect() {
	x.status = reconnecting
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
	x.status = connected
}
func (x *SignalConnection) sendSync(msg *xmpp.Chat) (response *xmpp.Chat, err error) {
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

func (x *SignalConnection) send(msg *xmpp.Chat) error {
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

func (x *SignalConnection) addWaiter(uid string, ch chan *xmpp.Chat) error {
	x.waitersMutex.Lock()
	defer x.waitersMutex.Unlock()
	if _, ok := x.waiters[uid]; ok {
		return errDuplicateWaiter
	}
	x.waiters[uid] = ch
	return nil
}

func (x *SignalConnection) removeWaiter(uid string) error {
	x.waitersMutex.Lock()
	defer x.waitersMutex.Unlock()
	delete(x.waiters, uid)
	return nil
}

func (x *SignalConnection) wait(ch chan *xmpp.Chat) (response *xmpp.Chat, err error) {
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

var reachok string

func init() {
	var cmd = xmppCommand{
		Command: commandTryReach,
	}
	data, err := json.Marshal(&cmd)
	if err != nil {
		log.Error(fmt.Sprintf("json marshal cmd err %s", err))
	}
	reachok = string(data)
}

func (x *SignalConnection) handleNewCommand(from, subject string, cmd *xmppCommand) {
	log.Trace(fmt.Sprintf("%s new command from:%s,subect:%s,cmd=%s", x.name, from, subject, utils.StringInterface1(cmd)))
	switch cmd.Command {
	case commandTryReach:
		err := x.send(&xmpp.Chat{
			Type:    "chat",
			Remote:  from,
			Subject: subject,
			Text:    reachok,
		})
		if err != nil {
			log.Error(fmt.Sprintf("handleNewCommand send err %s", err))
		}
	case commandExChangeSdp:
		fromaddr := strings.Split(from, "@")[0]
		r, err := x.SdpHandler(common.HexToAddress(fromaddr), cmd.OtherInfo)
		cmd2 := xmppCommand{
			Command: commandExChangeSdp,
		}
		if err != nil {
			cmd2.Error = err.Error()
			log.Error(fmt.Sprintf("%s cannot handle sdp request cmd=%s,err=%s", x.name, spew.Sdump(cmd), err))
		} else {
			cmd2.OtherInfo = r
		}
		data, err := json.Marshal(cmd2)
		if err != nil {
			log.Error(fmt.Sprintf("json marshal cmd error %s", err))
		}
		err = x.send(&xmpp.Chat{
			Type:    "chat",
			Remote:  from,
			Subject: subject,
			Text:    string(data),
		})
		if err != nil {
			log.Error(fmt.Sprintf("handleNewCommand send err %s", err))
		}
	default:
		log.Error(fmt.Sprintf("%s receive unkown from:%s,subject:%s cmd:%s", x.name, from, subject, spew.Sdump(cmd)))
	}
}

//Close this connection
func (x *SignalConnection) Close() {
	x.status = closed
	close(x.closed)
	err := x.client.Close()
	if err != nil {
		log.Error(fmt.Sprintf("handleNewCommand send err %s", err))
	}
}

//Connected returns true when this connection is ready for sent
func (x *SignalConnection) Connected() bool {
	return x.status == connected
}

type command int

const (
	commandTryReach = command(iota)
	commandExChangeSdp
)

type xmppCommand struct {
	Command   command
	Error     string
	OtherInfo string
}

func (x *SignalConnection) sendCommand(addr common.Address, cmd *xmppCommand) (cmdResponse *xmppCommand, err error) {
	chat := &xmpp.Chat{
		Remote:  fmt.Sprintf("%s%s", addr.String(), nameSuffix),
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
	var cmd2 xmppCommand
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

//TryReach test if `addr` is online or not
func (x *SignalConnection) TryReach(addr common.Address) error {
	_, err := x.sendCommand(addr, &xmppCommand{
		Command: commandTryReach,
	})
	return err
}

//ExchangeSdp exchange sdp info with `addr`
func (x *SignalConnection) ExchangeSdp(addr common.Address, sdp string) (partnerSdp string, err error) {
	r, err := x.sendCommand(addr, &xmppCommand{
		Command:   commandExChangeSdp,
		OtherInfo: sdp,
	})
	if err != nil {
		return "", err
	}
	return r.OtherInfo, nil
}
