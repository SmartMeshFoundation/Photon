package network

import (
	"sync/atomic"

	"crypto/ecdsa"
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
MixTransporter is a wrapper for two Transporter(UDP and ICE)
*/
type MixTransporter struct {
	udp *UDPTransport
	ice *IceTransport
	t   atomic.Value
}

/*
MixDiscovery is a wrapper for two Discover ,so it can switch between these discover
*/
type MixDiscovery struct {
	udp *Discovery
	ice *IceHelperDicovery
	d   atomic.Value
}

//NewMixTranspoter create a MixTransporter and discover
func NewMixTranspoter(key *ecdsa.PrivateKey, name string, host string, port int, conn *SafeUDPConnection, protocol ProtocolReceiver, policy Policier) (t *MixTransporter, d *MixDiscovery) {
	var err error
	var h Transporter
	t = &MixTransporter{}
	t.udp = NewUDPTransport(host, port, conn, protocol, policy)
	t.ice, err = NewIceTransporter(key, name)
	if err != nil {
		log.Error(fmt.Sprintf("new ice transport error %s,default will be udp transport", err))
		h = t.udp
		d = newMixDiscovery(false)
	} else {
		h = t.ice
		d = newMixDiscovery(true)
	}
	t.t.Store(&h)
	return
}
func (t *MixTransporter) getTranspoter() Transporter {
	return *t.t.Load().(*Transporter)
}

//Send message
func (t *MixTransporter) Send(receiver common.Address, host string, port int, data []byte) error {
	return t.getTranspoter().Send(receiver, host, port, data)
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
	if t.ice != nil {
		t.ice.Start()
	}
}

//Stop the two transporter
func (t *MixTransporter) Stop() {
	if t.ice != nil {
		t.ice.Stop()
	}
	if t.udp != nil {
		t.udp.Stop()
	}
}

//StopAccepting stops receiving for the two transporter
func (t *MixTransporter) StopAccepting() {
	if t.ice != nil {
		t.ice.StopAccepting()
	}
	if t.udp != nil {
		t.udp.StopAccepting()
	}
}

//RegisterProtocol register receiver for the two transporter
func (t *MixTransporter) RegisterProtocol(protcol ProtocolReceiver) {
	if t.ice != nil {
		t.ice.RegisterProtocol(protcol)
	}
	if t.udp != nil {
		t.udp.RegisterProtocol(protcol)
	}

}
func (t *MixTransporter) switchToUDP() bool {
	u := t.getTranspoter()
	_, ok := u.(*UDPTransport)
	if ok {
		log.Error(fmt.Sprintf("mixTransporter already uses udp ."))
		return false
	}
	u = t.udp
	t.t.Store(&u)
	return true
}
func (t *MixTransporter) switchToIce() bool {
	if t.ice == nil {
		log.Error(fmt.Sprintf("switch to ice,but  there is no ice transporter.use previous transpoter."))
		return false
	}
	i := t.getTranspoter()
	_, ok := i.(*IceTransport)
	if ok {
		log.Error(fmt.Sprintf("mixTransporter already uses ice ."))
		return false
	}
	i = t.ice
	t.t.Store(&i)
	return true
}
func newMixDiscovery(useIce bool) *MixDiscovery {
	m := &MixDiscovery{
		udp: NewDiscovery(),
		ice: NewIceHelperDiscovery(),
	}
	var h DiscoveryInterface
	if useIce {
		h = m.ice
	} else {
		h = m.udp
	}
	m.d.Store(&h)
	return m
}
func (d *MixDiscovery) getDefault() DiscoveryInterface {
	return *d.d.Load().(*DiscoveryInterface)
}

//Register a node
func (d *MixDiscovery) Register(address common.Address, host string, port int) error {
	return d.getDefault().Register(address, host, port)
}

//Get node's ip and port
func (d *MixDiscovery) Get(address common.Address) (host string, port int, err error) {
	return d.getDefault().Get(address)
}

//NodeIDByHostPort find a node by host and port
func (d *MixDiscovery) NodeIDByHostPort(host string, port int) (node common.Address, err error) {
	return d.getDefault().NodeIDByHostPort(host, port)
}
func (d *MixDiscovery) switchToUDP() bool {
	u := d.getDefault()
	_, ok := u.(*Discovery)
	if ok {
		log.Error(fmt.Sprintf("MixDiscovery already uses udp discovery"))
		return false
	}
	u = d.udp
	d.d.Store(&u)
	return true
}
func (d *MixDiscovery) switchToIce() bool {
	i := d.getDefault()
	_, ok := i.(*IceHelperDicovery)
	if ok {
		log.Error(fmt.Sprintf("MixDiscovery already uses ice helper discovery"))
		return false
	}
	i = d.ice
	d.d.Store(&i)
	return true
}

func (d *MixDiscovery) printNodes() {
	u, ok := d.getDefault().(*Discovery)
	if ok {
		log.Trace(fmt.Sprintf("nodes are:\n"))
		for k, v := range u.NodeIDHostPortMap {
			log.Trace(fmt.Sprintf("%s:%s", utils.APex(k), v))
		}
	} else {
		log.Trace(fmt.Sprintf("it's ice helper discovery ,cannot get nodes."))
	}
}
