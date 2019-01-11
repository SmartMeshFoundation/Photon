package network

import (
	"testing"

	"net"

	"time"

	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestNewMixTransport(t *testing.T) {
	if testing.Short() {
		return
	}
	key1, _ := utils.MakePrivateKeyAddress()
	key2, _ := utils.MakePrivateKeyAddress()
	key3, _ := utils.MakePrivateKeyAddress()
	m1, err := NewMixTranspoter("m1", params.DefaultTestXMPPServer, "127.0.0.1", 50001, key1, newDummyProtocol("m1"), &dummyPolicy{}, DeviceTypeMobile)
	if err != nil {
		t.Error(err)
		return
	}
	m2, err := NewMixTranspoter("m1", params.DefaultTestXMPPServer, "127.0.0.1", 50002, key2, newDummyProtocol("m2"), &dummyPolicy{}, DeviceTypeOther)
	if err != nil {
		t.Error(err)
		return
	}
	m3, err := NewMixTranspoter("m1", params.DefaultTestXMPPServer, "127.0.0.1", 50003, key3, newDummyProtocol("m3"), &dummyPolicy{}, DeviceTypeMobile)
	if err != nil {
		t.Error(err)
		return
	}
	m1.Start()
	m2.Start()
	m3.Start()
	defer m1.Stop()
	defer m2.Stop()
	defer m3.Stop()
	nodes := map[common.Address]*net.UDPAddr{
		m1.xmpp.NodeAddress: m1.udp.UAddr,
		m3.xmpp.NodeAddress: m3.udp.UAddr,
	}
	//m3.udp.setHostPort(nodes)
	m1.udp.setHostPort(nodes)
	datam12 := []byte("m1->m2")
	datam13 := []byte("m1->m3")
	err = m1.Send(m2.xmpp.NodeAddress, datam12)
	if err != nil {
		t.Error(err)
		return
	}
	err = m1.Send(m3.xmpp.NodeAddress, datam13)
	if err != nil {
		t.Error(err)
		return
	}
	err = m1.xmpp.conn.SubscribeNeighbour(m2.xmpp.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	err = m1.xmpp.conn.SubscribeNeighbour(m3.xmpp.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	err = m3.xmpp.conn.SubscribeNeighbour(m1.xmpp.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	deviceType, isOnline := m1.NodeStatus(m2.xmpp.NodeAddress)
	if !isOnline || deviceType != DeviceTypeOther {
		t.Error("type error")
		return
	}
	deviceType, isOnline = m1.NodeStatus(m3.xmpp.NodeAddress)
	if !isOnline || deviceType != DeviceTypeOther {
		t.Error("type error")
		return
	}
	deviceType, isOnline = m3.NodeStatus(m1.xmpp.NodeAddress)
	if !isOnline || deviceType != DeviceTypeMobile {
		t.Error("type error")
		return
	}
	err = m3.Send(m1.xmpp.NodeAddress, []byte("m3->m1"))
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second * 3)
}
