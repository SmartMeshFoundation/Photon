package network

import (
	"crypto/ecdsa"
	"sync/atomic"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

type MixTransporter struct {
	udp *UDPTransport
	ice *IceTransport
	t   atomic.Value
}
type MixDiscovery struct {
	udp *Discovery
	ice *IceHelperDicovery
	d   atomic.Value
}

func NewMixTranspoter(key *ecdsa.PrivateKey, name string, host string, port int, conn *SafeUdpConnection, protocol ProtocolReceiver, policy Policier) (t *MixTransporter, d *MixDiscovery) {
	var err error
	var h Transporter
	t = &MixTransporter{}
	t.udp = NewUDPTransport(host, port, conn, protocol, policy)
	t.ice, err = NewIceTransporter(key, name)
	if err != nil {
		log.Error(fmt.Sprintf("new ice transport error %s,default will be udp transport", err))
		h = t.udp
		d = NewMixDiscovery(false)
	} else {
		h = t.ice
		d = NewMixDiscovery(true)
	}
	t.t.Store(&h)
	return
}
func (t *MixTransporter) getTranspoter() Transporter {
	return *t.t.Load().(*Transporter)
}
func (t *MixTransporter) Send(receiver common.Address, host string, port int, data []byte) error {
	return t.getTranspoter().Send(receiver, host, port, data)
}
func (t *MixTransporter) Receive(data []byte, host string, port int) error {
	panic("useless")
}
func (t *MixTransporter) Start() {
	t.getTranspoter().Start()
}
func (t *MixTransporter) Stop() {
	if t.ice != nil {
		t.ice.Stop()
	}
	if t.udp != nil {
		t.udp.Stop()
	}
}
func (t *MixTransporter) StopAccepting() {
	t.getTranspoter().StopAccepting()
} //stop receiving data
func (t *MixTransporter) RegisterProtocol(protcol ProtocolReceiver) {
	if t.ice != nil {
		t.ice.RegisterProtocol(protcol)
	}
	if t.udp != nil {
		t.udp.RegisterProtocol(protcol)
	}

} //register transporter to protocol
func (t *MixTransporter) switchToUdp() bool {
	u := t.getTranspoter()
	_, ok := u.(*UDPTransport)
	if ok {
		log.Error(fmt.Sprintf("mixTransporter already uses udp ."))
		return false
	}
	u = t.udp
	t.t.Store(&u)
	t.ice.StopAccepting()
	t.udp.Start()
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
	t.udp.StopAccepting()
	t.ice.Start()
	return true
}
func NewMixDiscovery(useIce bool) *MixDiscovery {
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
func (d *MixDiscovery) Register(address common.Address, host string, port int) error {
	return d.getDefault().Register(address, host, port)
}
func (d *MixDiscovery) Get(address common.Address) (host string, port int, err error) {
	return d.getDefault().Get(address)
}
func (d *MixDiscovery) NodeIdByHostPort(host string, port int) (node common.Address, err error) {
	return d.getDefault().NodeIdByHostPort(host, port)
}
func (d *MixDiscovery) switchToUdp() bool {
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
		for k, v := range u.NodeIdHostPortMap {
			log.Trace(fmt.Sprintf("%s:%s", utils.APex(k), v))
		}
	} else {
		log.Trace(fmt.Sprintf("it's ice helper discovery ,cannot get nodes."))
	}
}
