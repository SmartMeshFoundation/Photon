package network

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport/xmpppass"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-errors/errors"
)

var errXMPPConnectionNotReady = errors.New("xmpp connection not ready")

//XMPPTransporter use XMPP to comminucate with other raiden nodes
type XMPPTransporter struct {
	conn          *xmpptransport.XMPPConnection
	quitChan      chan struct{}
	stopped       bool
	stopReceiving bool
	log           log.Logger
	protocol      ProtocolReceiver
}

/*
NewXMPPTransporter create xmpp transporter,
if not success ,for example cannot connect to xmpp server, will try background
*/
func NewXMPPTransporter(name, ServerURL string, key *ecdsa.PrivateKey, deviceType string) (x *XMPPTransporter) {
	x = &XMPPTransporter{
		quitChan: make(chan struct{}),
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	passwordFn := func() xmpptransport.GetCurrentPasswordFunc {
		f1 := func() string {
			pass, _ := xmpppass.CreatePassword(key)
			return pass
		}
		return f1
	}
	x.log = log.New("name", name)
	go func() {
		wait := time.Millisecond
		var err error
		for {
			select {
			case <-time.After(wait):
				x.conn, err = xmpptransport.NewConnection(ServerURL, addr, passwordFn(), x.ReceiveData, name, deviceType)
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
	return
}

//Send a message
func (x *XMPPTransporter) Send(receiver common.Address, data []byte) error {
	if x.stopped || x.conn == nil {
		return errXMPPConnectionNotReady
	}
	x.log.Trace(fmt.Sprintf("send to %s, message=%s", utils.APex2(receiver), encoding.MessageType(data[0])))
	return x.conn.SendData(receiver, data)
}

//ReceiveData call back of xmpp connection
func (x *XMPPTransporter) ReceiveData(from common.Address, data []byte) {
	x.log.Trace(fmt.Sprintf("received from %s, message=%s", utils.APex2(from), encoding.MessageType(data[0])))
	if x.stopped || x.stopReceiving {
		return
	}
	if x.protocol != nil {
		x.protocol.receive(data)
	}
}

//Start ,ready for send and receive
func (x *XMPPTransporter) Start() {

}

//Stop send and receive
func (x *XMPPTransporter) Stop() {
	x.stopped = true
	close(x.quitChan)
}

//StopAccepting stops receiving
func (x *XMPPTransporter) StopAccepting() {
	x.stopReceiving = true
}

//RegisterProtocol a receiver
func (x *XMPPTransporter) RegisterProtocol(protcol ProtocolReceiver) {
	x.protocol = protcol
}

//NodeStatus get node's status and is online right now
func (x *XMPPTransporter) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
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
