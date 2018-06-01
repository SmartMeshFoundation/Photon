package network

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
MixTransporter is a wrapper for two Transporter(UDP and XMPP)
if I can reach the node by UDP,then UDP,
if I cannot reach the node, try XMPP
*/
type MixTransporter struct {
	udp  *UDPTransport
	xmpp *XMPPTransporter
	d    *MixDiscovery
	name string
}

/*
MixDiscovery is a wrapper for two Discover ,so it can switch between these discover
*/
type MixDiscovery struct {
	udp  *Discovery
	mock *MockDicovery
}

//NewMixTranspoter create a MixTransporter and discover
func NewMixTranspoter(name, host string, port int, protocol ProtocolReceiver, policy Policier) (t *MixTransporter, d *MixDiscovery, err error) {
	t = &MixTransporter{}
	t.udp, err = NewUDPTransport(host, port, protocol, policy)
	if err != nil {
		return
	}
	t.xmpp, err = NewXMPPTransporter()
	if err != nil {
		log.Error(fmt.Sprintf("new xmpp transport error %s,default will be udp transport", err))
		return
	}
	d = newMixDiscovery()
	t.d = d
	t.name = name
	return
}

/*
Send message
优先选择局域网,在局域网走不通的情况下,才会考虑 xmpp
*/
func (t *MixTransporter) Send(receiver common.Address, host string, port int, data []byte) error {
	_, ok := t.d.udp.NodeIDHostPortMap[receiver]
	if ok {
		return t.udp.Send(receiver, host, port, data)
	} else if t.xmpp != nil {
		return t.xmpp.Send(receiver, host, port, data)
	} else {
		err := fmt.Errorf("no valid %s send to %s %s:%d, message=%s,response hash=%s", t.name, utils.APex2(receiver), host, port, encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:])))
		log.Error(err.Error())
		return err
	}
}

//receive  just for Transporter interface
func (t *MixTransporter) receive(data []byte, host string, port int) error {
	panic("useless")
}

//Start the two transporter
func (t *MixTransporter) Start() {
	if t.udp != nil {
		t.udp.Start()
	}
	if t.xmpp != nil {
		t.xmpp.Start()
	}
}

//Stop the two transporter
func (t *MixTransporter) Stop() {
	if t.xmpp != nil {
		t.xmpp.Stop()
	}
	if t.udp != nil {
		t.udp.Stop()
	}
}

//StopAccepting stops receiving for the two transporter
func (t *MixTransporter) StopAccepting() {
	if t.xmpp != nil {
		t.xmpp.StopAccepting()
	}
	if t.udp != nil {
		t.udp.StopAccepting()
	}
}

//RegisterProtocol register receiver for the two transporter
func (t *MixTransporter) RegisterProtocol(protcol ProtocolReceiver) {
	if t.xmpp != nil {
		t.xmpp.RegisterProtocol(protcol)
	}
	if t.udp != nil {
		t.udp.RegisterProtocol(protcol)
	}

}

func newMixDiscovery() *MixDiscovery {
	m := &MixDiscovery{
		udp:  NewDiscovery(),
		mock: NewMockDicovery(),
	}
	return m
}

//Register a node
func (d *MixDiscovery) Register(address common.Address, host string, port int) error {
	return d.udp.Register(address, host, port)
}

//Get node's ip and port
func (d *MixDiscovery) Get(address common.Address) (host string, port int, err error) {
	host, port, err = d.udp.Get(address)
	if err != nil {
		host, port, err = d.mock.Get(address)
	}
	return
}

//NodeIDByHostPort find a node by host and port
func (d *MixDiscovery) NodeIDByHostPort(host string, port int) (node common.Address, err error) {
	node, err = d.udp.NodeIDByHostPort(host, port)
	if err != nil {
		node, err = d.mock.NodeIDByHostPort(host, port)
	}
	return
}
func (d *MixDiscovery) printNodes() {
	u := d.udp
	log.Trace(fmt.Sprintf("nodes are:\n"))
	for k, v := range u.NodeIDHostPortMap {
		log.Trace(fmt.Sprintf("%s:%s", utils.APex(k), v))
	}
}
