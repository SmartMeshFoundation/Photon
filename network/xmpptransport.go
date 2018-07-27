package network

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/netshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport/xmpppass"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-errors/errors"
)

var errXMPPConnectionNotReady = errors.New("xmpp connection not ready")

//XMPPTransport use XMPP to comminucate with other raiden nodes
type XMPPTransport struct {
	conn          *xmpptransport.XMPPConnection
	quitChan      chan struct{}
	stopped       bool
	stopReceiving bool
	log           log.Logger
	protocol      ProtocolReceiver
	NodeAddress   common.Address
	key           *ecdsa.PrivateKey
	statusChan    chan netshare.Status
}

/*
NewXMPPTransport create xmpp transporter,
if not success ,for example cannot connect to xmpp server, will try background
*/
func NewXMPPTransport(name, ServerURL string, key *ecdsa.PrivateKey, deviceType string) (x *XMPPTransport) {
	x = &XMPPTransport{
		quitChan:    make(chan struct{}),
		NodeAddress: crypto.PubkeyToAddress(key.PublicKey),
		key:         key,
		statusChan:  make(chan netshare.Status, 10),
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	x.log = log.New("name", name)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wait := time.Millisecond
		var err error
		//only wait one time
		var first bool
		for {
			select {
			case <-time.After(wait):
				x.conn, err = xmpptransport.NewConnection(ServerURL, addr, x, x, name, deviceType, x.statusChan)
				if !first {
					first = true
					wg.Done()
				}
				if err != nil {
					x.log.Error(fmt.Sprintf("cannot connect to xmpp server %s", ServerURL))
					time.Sleep(time.Second * 5)
				} else {
					return
				}
			case <-x.quitChan:
				return
			}

		}

	}()
	//only wait for one try
	wg.Wait()
	return
}

//GetPassWord returns current login password
func (x *XMPPTransport) GetPassWord() string {
	pass, err := xmpppass.CreatePassword(x.key)
	if err != nil {
		log.Error(fmt.Sprintf("GetPassWord for %s err %s", utils.APex2(x.NodeAddress), err))
	}
	err = xmpppass.VerifySignature(x.NodeAddress.String(), pass)
	if err != nil {
		panic(err)
	}
	return pass
}

//Send a message
func (x *XMPPTransport) Send(receiver common.Address, data []byte) error {
	x.log.Trace(fmt.Sprintf("send to %s, message=%s", utils.APex2(receiver), encoding.MessageType(data[0])))
	if x.stopped || x.conn == nil {
		return errXMPPConnectionNotReady
	}
	return x.conn.SendData(receiver, data)
}

//DataHandler call back of xmpp connection
func (x *XMPPTransport) DataHandler(from common.Address, data []byte) {
	x.log.Trace(fmt.Sprintf("received from %s, message=%s", utils.APex2(from), encoding.MessageType(data[0])))
	if x.stopped || x.stopReceiving {
		return
	}
	if x.protocol != nil {
		x.protocol.receive(data)
	}
}

//Start ,ready for send and receive
func (x *XMPPTransport) Start() {

}

//Stop send and receive
func (x *XMPPTransport) Stop() {
	x.stopped = true
	close(x.quitChan)
	if x.conn != nil {
		x.conn.Close()
	}
}

//StopAccepting stops receiving
func (x *XMPPTransport) StopAccepting() {
	x.stopReceiving = true
}

//RegisterProtocol a receiver
func (x *XMPPTransport) RegisterProtocol(protcol ProtocolReceiver) {
	x.protocol = protcol
}

//NodeStatus get node's status and is online right now
func (x *XMPPTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	if x.conn == nil {
		return "", false
	}
	var err error
	deviceType, isOnline, err = x.conn.IsNodeOnline(addr)
	if err != nil {
		x.log.Error(fmt.Sprintf("IsNodeOnline query %s err %s", utils.APex2(addr), err))
	}
	return
}
