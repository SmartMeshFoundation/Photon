package network

import (
	"crypto/ecdsa"

	"errors"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/network/netshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport"
	"github.com/ethereum/go-ethereum/common"
)

/*
MatrixMixTransporter is a wrapper for two Transporter(UDP and Matrix)
if I can reach the node by UDP,then UDP,
if I cannot reach the node, try Matrix
*/
type MatrixMixTransporter struct {
	udp      *UDPTransport
	matirx   *MatrixTransport
	name     string
	protocol ProtocolReceiver
}

//NewMatrixMixTransporter create a MixTransporter and discover
func NewMatrixMixTransporter(name, host string, port int, key *ecdsa.PrivateKey, protocol ProtocolReceiver, policy Policier, deviceType string) (t *MatrixMixTransporter, err error) {
	t = &MatrixMixTransporter{
		name:     name,
		protocol: protocol,
	}
	t.udp, err = NewUDPTransport(name, host, port, protocol, policy)
	if err != nil {
		return
	}
	t.matirx, err = InitMatrixTransport(name, key, deviceType)
	t.RegisterProtocol(protocol)
	return
}

/*
Send message
优先选择局域网,在局域网走不通的情况下,才会考虑 matrix
*/
/*
 *	Send message prefers to choose LAN,
 *	after LAN does not work, then try matrix.
 */
func (t *MatrixMixTransporter) Send(receiver common.Address, data []byte) error {
	_, isOnline := t.udp.NodeStatus(receiver)
	if isOnline {
		err := t.udp.Send(receiver, data)
		if err != nil {
			return err
		}
	} else if t.matirx != nil {
		err := t.matirx.Send(receiver, data)
		if err != nil {
			return err
		}
	}
	return nil
}

//Start the two transporter
func (t *MatrixMixTransporter) Start() {
	if t.udp != nil {
		t.udp.Start()
	}
	/*	if t.xmpp != nil {
		t.xmpp.Start()
	}*/
	if t.matirx != nil {
		t.matirx.Start()
	}
}

//Stop the two transporter
func (t *MatrixMixTransporter) Stop() {

	if t.udp != nil {
		t.udp.Stop()
	}
	if t.matirx != nil {
		t.matirx.Stop()
	}
}

//StopAccepting stops receiving for the two transporter
func (t *MatrixMixTransporter) StopAccepting() {
	if t.udp != nil {
		t.udp.StopAccepting()
	}
	if t.matirx != nil {
		t.matirx.StopAccepting()
	}
}

//RegisterProtocol register receiver for the two transporter
func (t *MatrixMixTransporter) RegisterProtocol(protcol ProtocolReceiver) {
	if t.udp != nil {
		t.udp.RegisterProtocol(protcol)
	}
	if t.matirx != nil {
		t.matirx.RegisterProtocol(protcol)
	}
}

//NodeStatus get node's status and is online right now
func (t *MatrixMixTransporter) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	deviceType, isOnline = t.udp.NodeStatus(addr)
	if isOnline {
		return
	}
	/*return t.xmpp.NodeStatus(addr)*/
	return t.matirx.NodeStatus(addr)
}

//GetNotify notification of connection status change
func (t *MatrixMixTransporter) GetNotify() (notify <-chan netshare.Status, err error) {
	return nil, errors.New("connection not established")
}

//SubscribeNeighbor get the status change notification of partner node
//func (t *MatrixMixTransporter) SubscribeNeighbor(db xmpptransport.XMPPDb) error {
func (t *MatrixMixTransporter) SubscribeNeighbor(db xmpptransport.XMPPDb) error {
	if t.matirx == nil {
		return fmt.Errorf("Try to subscribe neighbor,but matrix connection is disconnected")
	}
	return t.matirx.CollectNeighbors(db)
}
