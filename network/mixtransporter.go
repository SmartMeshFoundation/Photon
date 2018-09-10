package network

import (
	"crypto/ecdsa"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/network/netshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport"
	"github.com/ethereum/go-ethereum/common"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

/*
MixTransporter is a wrapper for two Transporter(UDP and XMPP)
if I can reach the node by UDP,then UDP,
if I cannot reach the node, try XMPP
*/
type MixTransporter struct {
	udp      *UDPTransport
	xmpp     *XMPPTransport
	matirx	 *MatrixTransport
	name     string
	protocol ProtocolReceiver
}

//NewMixTranspoter create a MixTransporter and discover
func NewMixTranspoter(name, xmppServer, host string, port int, key *ecdsa.PrivateKey, protocol ProtocolReceiver, policy Policier, deviceType string) (t *MixTransporter, err error) {
	t = &MixTransporter{
		name:     name,
		protocol: protocol,
	}
	t.udp, err = NewUDPTransport(name, host, port, protocol, policy)
	if err != nil {
		return
	}
	/*t.xmpp = NewXMPPTransport(name, xmppServer, key, deviceType)
	t.RegisterProtocol(protocol)*/

	t.matirx, err = InitMatrixTransport(name, params.DefaultMatrixServer, key, deviceType)
	t.RegisterProtocol(protocol)
	return
}

/*
Send message
优先选择局域网,在局域网走不通的情况下,才会考虑 xmpp
*/
func (t *MixTransporter) Send(receiver common.Address, data []byte) error {
/*	_, isOnline := t.udp.NodeStatus(receiver)
	if isOnline {
		return t.udp.Send(receiver, data)
	} else if t.xmpp != nil {
		return t.xmpp.Send(receiver, data)
	} else {
		err := fmt.Errorf("no valid %s send to %s , message=%s,response hash=%s", t.name, utils.APex2(receiver), encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:])))
		log.Error(err.Error())
		return err
	}*/
		_, isOnline := t.udp.NodeStatus(receiver)
	if isOnline {
		t.udp.Send(receiver, data)
	} else if t.matirx != nil {
		//t.xmpp.Send(receiver, data)
		t.matirx.Send(receiver, data)
	}
	return nil
}

//Start the two transporter
func (t *MixTransporter) Start() {
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
func (t *MixTransporter) Stop() {
/*	if t.xmpp != nil {
		t.xmpp.Stop()
	}*/
	if t.udp != nil {
		t.udp.Stop()
	}
	if t.matirx != nil {
		t.matirx.Stop()
	}
}

//StopAccepting stops receiving for the two transporter
func (t *MixTransporter) StopAccepting() {
/*	if t.xmpp != nil {
		t.xmpp.StopAccepting()
	}*/
	if t.udp != nil {
		t.udp.StopAccepting()
	}
	if t.matirx!=nil{
		t.matirx.StopAccepting()
	}
}

//RegisterProtocol register receiver for the two transporter
func (t *MixTransporter) RegisterProtocol(protcol ProtocolReceiver) {
/*	if t.xmpp != nil {
		t.xmpp.RegisterProtocol(protcol)
	}*/
	if t.udp != nil {
		t.udp.RegisterProtocol(protcol)
	}
	if t.matirx != nil {
		t.matirx.RegisterProtocol(protcol)
	}
}

//NodeStatus get node's status and is online right now
func (t *MixTransporter) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	deviceType, isOnline = t.udp.NodeStatus(addr)
	if isOnline {
		return
	}
	/*return t.xmpp.NodeStatus(addr)*/
	return t.matirx.NodeStatus(addr)
}

//GetNotify notification of connection status change
func (t *MixTransporter) GetNotify() (notify <-chan netshare.Status, err error) {
	if t.xmpp.conn != nil {
		return t.xmpp.statusChan, nil
	}
	return nil, errors.New("connection not established")
}

//SubscribeNeighbor get the status change notification of partner node
func (t *MixTransporter) SubscribeNeighbor(db xmpptransport.XMPPDb) error {
	/*if t.xmpp.conn == nil {
		return fmt.Errorf("try to subscribe neighbor,but xmpp connection is disconnected")
	}
	return t.xmpp.conn.CollectNeighbors(db)*/
	return nil
}
