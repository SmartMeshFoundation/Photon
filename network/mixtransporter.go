package network

import (
	"fmt"

	"crypto/ecdsa"

	"errors"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/network/xmpptransport"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
MixTransport is a wrapper for two Transporter(UDP and XMPP)
if I can reach the node by UDP,then UDP,
if I cannot reach the node, try XMPP
*/
type MixTransport struct {
	udp      *UDPTransport
	xmpp     *XMPPTransport
	name     string
	protocol ProtocolReceiver
}

//NewMixTranspoter create a MixTransport and discover
func NewMixTranspoter(name, xmppServer, host string, port int, key *ecdsa.PrivateKey, protocol ProtocolReceiver, policy Policier, deviceType string) (t *MixTransport, err error) {
	t = &MixTransport{
		name:     name,
		protocol: protocol,
	}
	t.udp, err = NewUDPTransport(name, host, port, protocol, policy)
	if err != nil {
		return
	}
	t.xmpp = NewXMPPTransport(name, xmppServer, key, deviceType)
	t.RegisterProtocol(protocol)
	return
}

/*
Send message
优先选择局域网,在局域网走不通的情况下,才会考虑 xmpp
*/
/*
 *	Send : function to send out messages.
 *
 *	Note that this function prefers to choose LAN, ifor new c local network does not work,
 * 	then it chooses xmpp.
 */
func (t *MixTransport) Send(receiver common.Address, data []byte) error {
	_, isOnline := t.udp.NodeStatus(receiver)
	if isOnline {
		return t.udp.Send(receiver, data)
	} else if t.xmpp != nil {
		return t.xmpp.Send(receiver, data)
	} else {
		err := fmt.Errorf("no valid %s send to %s , message=%s,response hash=%s", t.name, utils.APex2(receiver), encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:])))
		log.Error(err.Error())
		return err
	}
}

//Start the two transporter
func (t *MixTransport) Start() {
	if t.udp != nil {
		t.udp.Start()
	}
	if t.xmpp != nil {
		t.xmpp.Start()
	}
}

//Stop the two transporter
func (t *MixTransport) Stop() {
	if t.xmpp != nil {
		t.xmpp.Stop()
	}
	if t.udp != nil {
		t.udp.Stop()
	}
}

//StopAccepting stops receiving for the two transporter
func (t *MixTransport) StopAccepting() {
	if t.xmpp != nil {
		t.xmpp.StopAccepting()
	}
	if t.udp != nil {
		t.udp.StopAccepting()
	}
}

//RegisterProtocol register receiver for the two transporter
func (t *MixTransport) RegisterProtocol(protcol ProtocolReceiver) {
	if t.xmpp != nil {
		t.xmpp.RegisterProtocol(protcol)
	}
	if t.udp != nil {
		t.udp.RegisterProtocol(protcol)
	}

}

//NodeStatus get node's status and is online right now
func (t *MixTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	deviceType, isOnline = t.udp.NodeStatus(addr)
	if isOnline {
		return
	}
	return t.xmpp.NodeStatus(addr)
}

//GetNotify notification of connection status change
func (t *MixTransport) GetNotify() (notify <-chan netshare.Status, err error) {
	if t.xmpp.conn != nil {
		return t.xmpp.statusChan, nil
	}
	return nil, errors.New("connection not established")
}

//SubscribeNeighbor get the status change notification of partner node
func (t *MixTransport) SubscribeNeighbor(db xmpptransport.XMPPDb) error {
	if t.xmpp.conn == nil {
		return fmt.Errorf("try to subscribe neighbor,but xmpp connection is disconnected")
	}
	return t.xmpp.conn.CollectNeighbors(db)
}

// Reconnect :
func (t *MixTransport) Reconnect() {
	t.xmpp.conn.Reconnect()
}
