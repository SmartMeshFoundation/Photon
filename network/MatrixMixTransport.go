package network

import (
	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/Photon/models"

	"github.com/SmartMeshFoundation/Photon/network/wakeuphandler"

	"github.com/SmartMeshFoundation/Photon/params"

	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/ethereum/go-ethereum/common"
)

/*
MatrixMixTransport is a wrapper for two Transporter(UDP and Matrix)
if I can reach the node by UDP,then UDP,
if I cannot reach the node, try Matrix
*/
type MatrixMixTransport struct {
	udp      *UDPTransport
	matirx   *MatrixTransport
	name     string
	protocol ProtocolReceiver
	*wakeuphandler.MixWakeUpHandler
}

//NewMatrixMixTransporter create a MixTransport and discover
func NewMatrixMixTransporter(name, host string, port int, key *ecdsa.PrivateKey, protocol ProtocolReceiver, policy Policier, deviceType string, dao models.Dao) (t *MatrixMixTransport, err error) {
	t = &MatrixMixTransport{
		name:     name,
		protocol: protocol,
	}
	t.udp, err = NewUDPTransport(name, host, port, protocol, policy)
	if err != nil {
		return
	}
	t.matirx = NewMatrixTransport(name, key, deviceType, params.TrustMatrixServers, dao)
	t.RegisterProtocol(protocol)
	t.MixWakeUpHandler = wakeuphandler.NewMixWakeUpHandler(t.udp.WakeUpHandler, t.matirx.WakeUpHandler)
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
func (t *MatrixMixTransport) Send(receiver common.Address, data []byte) error {
	_, isOnline := t.udp.NodeStatus(receiver)
	if isOnline {
		err := t.udp.Send(receiver, data)
		if err != nil {
			return err
		}
	}
	if t.matirx != nil {
		err := t.matirx.Send(receiver, data)
		if err != nil {
			return err
		}
	}
	return nil
}

//Start the two transporter
func (t *MatrixMixTransport) Start() {
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
func (t *MatrixMixTransport) Stop() {

	if t.udp != nil {
		t.udp.Stop()
	}
	if t.matirx != nil {
		t.matirx.Stop()
	}
}

//StopAccepting stops receiving for the two transporter
func (t *MatrixMixTransport) StopAccepting() {
	if t.udp != nil {
		t.udp.StopAccepting()
	}
	if t.matirx != nil {
		t.matirx.StopAccepting()
	}
}

//RegisterProtocol register receiver for the two transporter
func (t *MatrixMixTransport) RegisterProtocol(protcol ProtocolReceiver) {
	if t.udp != nil {
		t.udp.RegisterProtocol(protcol)
	}
	if t.matirx != nil {
		t.matirx.RegisterProtocol(protcol)
	}
}

//NodeStatus get node's status and is online right now
func (t *MatrixMixTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	deviceType, isOnline = t.udp.NodeStatus(addr)
	if isOnline {
		return
	}
	return t.matirx.NodeStatus(addr)
}

//UDPNodeStatus get node's status of UDPTransport
func (t *MatrixMixTransport) UDPNodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	deviceType, isOnline = t.udp.NodeStatus(addr)
	return
}

//GetNotify notification of connection status change
func (t *MatrixMixTransport) GetNotify() (notify <-chan netshare.Status, err error) {
	//if t.matirx != nil {
	return t.matirx.statusChan, nil
	//}
	//return nil, errors.New("connection not established")
}

////SetMatrixDB get the status change notification of partner node
////func (t *MatrixMixTransport) SetMatrixDB(db xmpptransport.XMPPDb)  {
//func (t *MatrixMixTransport) SetMatrixDB(db xmpptransport.XMPPDb) {
//	t.matirx.setDB(db)
//	return
//}
