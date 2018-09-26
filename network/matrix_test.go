package network

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/network/gomatrix"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/models/cb"

	"github.com/ethereum/go-ethereum/common"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

var testPrivKey *ecdsa.PrivateKey
var testAddress common.Address

type MockDb struct {
	channels []*channeltype.Serialization
}

func (db *MockDb) addPartner(address common.Address) {
	db.channels = append(db.channels, &channeltype.Serialization{
		PartnerAddressBytes: address[:],
	})
}
func (db *MockDb) XMPPIsAddrSubed(addr common.Address) bool {
	return true
}
func (db *MockDb) XMPPMarkAddrSubed(addr common.Address) {
	return
}
func (db *MockDb) GetChannelList(token, partner common.Address) (cs []*channeltype.Serialization, err error) {
	return db.channels, nil
}
func (db *MockDb) RegisterNewChannellCallback(f cb.ChannelCb) {

}
func (db *MockDb) RegisterChannelStateCallback(f cb.ChannelCb) {

}
func (db *MockDb) XMPPUnMarkAddr(addr common.Address) {

}
func init() {
	var err error
	ALIASFRAGMENT = fmt.Sprintf("testdiscovery-%s", utils.RandomString(10))
	//ALIASFRAGMENT = "testdiscovery"
	bin, err := hex.DecodeString("67fcbd8f1ed7b411813ab744c7327c79e860ac556d26a87046af3d2e0676d0c9")
	if err != nil {
		panic(err)
	}
	testPrivKey, err = crypto.ToECDSA(bin)
	if err != nil {
		panic(err)
	}
	testAddress = crypto.PubkeyToAddress(testPrivKey.PublicKey)
}
func getMatrixEnvConfig() (cfg1, cfg2, cfg3 map[string]string) {
	cfg1 = make(map[string]string)
	cfg2 = make(map[string]string)
	cfg3 = make(map[string]string)
	cfg1["transport01.smartmesh.cn"] = "http://transport01.smartmesh.cn:8008"
	cfg2["transport02.smartmesh.cn"] = "http://transport02.smartmesh.cn:8008"
	cfg3["transport03.smartmesh.cn"] = "http://transport03.smartmesh.cn:8008"
	return
}
func newTestMatrixTransport(name string, cfg map[string]string) (m1 *MatrixTransport) {
	key, _ := utils.MakePrivateKeyAddress()
	m1 = NewMatrixTransport(name, key, "other", cfg)
	m1.setDB(&MockDb{})
	return m1
}

func newFourTestMatrixTransport() (m0, m1, m2, m3 *MatrixTransport) {
	cfg0 := params.MatrixServerConfig
	cfg1, cfg2, cfg3 := getMatrixEnvConfig()
	m0 = newTestMatrixTransport("m0", cfg0)
	m1 = newTestMatrixTransport("m1", cfg1)
	m2 = newTestMatrixTransport("m2", cfg2)
	m3 = newTestMatrixTransport("m3", cfg3)
	m0.setDB(&MockDb{})
	m1.setDB(&MockDb{})
	m2.setDB(&MockDb{})
	m3.setDB(&MockDb{})
	log.Info(fmt.Sprintf("node0=%s,node1=%s,node2=%s,node3=%s",
		utils.APex2(m0.NodeAddress),
		utils.APex2(m1.NodeAddress),
		utils.APex2(m2.NodeAddress),
		utils.APex2(m3.NodeAddress)))
	return
}

func TestCreateMatrixTransport(t *testing.T) {
	m1 := newTestMatrixTransport("mrand", params.MatrixServerConfig)
	m1.Stop()
}
func TestRegisterAndJoinDiscoveryRoom(t *testing.T) {
	cfg1, _, _ := getMatrixEnvConfig()
	m1 := newTestMatrixTransport("mrand", cfg1)
	m1.Start()
	/*
		观察初次注册加入 discovery room 是否返回其他人在线信息.
	*/
	time.Sleep(time.Minute * 10)
}
func TestLoginAndJoinDiscoveryRoom(t *testing.T) {
	m1 := NewMatrixTransport("test", testPrivKey, "other", params.MatrixServerConfig)
	m1.setDB(&MockDb{})
	log.Trace(fmt.Sprintf("privkey=%s", hex.EncodeToString(crypto.FromECDSA(m1.key))))
	defer m1.Stop()
	m1.Start()
	time.Sleep(time.Minute * 10)
}

func TestGetJoinedRoomAlias(t *testing.T) {
	m1 := NewMatrixTransport("test", testPrivKey, "other", params.MatrixServerConfig)
	m1.setDB(&MockDb{})
	err := m1.loginOrRegister()
	if err != nil {
		panic(err)
	}
	rooms, err := m1.matrixcli.JoinedRooms()
	log.Trace(fmt.Sprintf("rooms=%s", rooms))
	if err != nil {
		panic(err)
	}
	for _, r := range rooms.JoinedRooms {
		var dat map[string]interface{}
		m1.matrixcli.StateEvent(r, "m.room.aliases", &dat)
		log.Info(fmt.Sprintf("room %s name=%s", r, dat))
	}
}
func TestInvite(t *testing.T) {
	_, m1, m2, m3 := newFourTestMatrixTransport()
	m1.Start()
	m2.Start()
	m3.Start()
	log.Info(fmt.Sprintf(" m1=%s,m2=%s,m3=%s", m1.UserID, m2.UserID, m3.UserID))
	r, err := m1.matrixcli.CreateRoom(&gomatrix.ReqCreateRoom{
		Invite:        []string{m2.UserID},
		Visibility:    "public",
		Preset:        "public_chat",
		RoomAliasName: fmt.Sprintf("smartraiden_%s_%s", utils.APex2(m1.NodeAddress), utils.APex2(m2.NodeAddress)),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("new room=%s", r.RoomID)
	_, err = m1.matrixcli.InviteUser(r.RoomID, &gomatrix.ReqInviteUser{
		UserID: m3.UserID,
	})
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second * 10)
}
func TestSearchNode(t *testing.T) {
	m0, m1, m2, m3 := newFourTestMatrixTransport()
	m0.Stop()
	//defer m1.Stop()
	//defer m2.Stop()
	//defer m3.Stop()
	err := m1.loginOrRegister()
	if err != nil {
		t.Error(err)
		return
	}
	err = m2.loginOrRegister()
	if err != nil {
		t.Error(err)
		return
	}
	err = m3.loginOrRegister()
	if err != nil {
		t.Error(err)
		return
	}
	users, err := m3.searchNode(m1.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	m3.log.Trace(fmt.Sprintf("users=%s", utils.StringInterface(users, 3)))
	if len(users) != 0 {
		t.Errorf("can not search cross homeserver when they dont join the same room")
		return
	}
	err = m1.joinDiscoveryRoom()
	if err != nil {
		t.Error(err)
		return
	}
	err = m3.joinDiscoveryRoom()
	if err != nil {
		t.Error(err)
		return
	}
	//must wait some time for home server to sync users to db?
	time.Sleep(time.Second * 6)
	users, err = m1.searchNode(m3.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	if len(users) == 0 {
		t.Errorf("must find node3")
		return
	}
	users, err = m3.searchNode(m1.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	if len(users) == 0 {
		t.Error("must find node1")
		return
	}
	m3.log.Trace(fmt.Sprintf("users_2=%s", utils.StringInterface(users, 3)))
}

func TestSendMessage(t *testing.T) {
	_, m1, m2, _ := newFourTestMatrixTransport()
	m1.db.(*MockDb).addPartner(m2.NodeAddress)
	m2.db.(*MockDb).addPartner(m1.NodeAddress)
	m1.Start()
	////let server sync
	//time.Sleep(time.Second * 6)
	m2.Start()
	m1Chan := make(chan string)
	m2Chan := make(chan string)
	m1.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m1Chan <- string(data)
	})
	m2.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m2Chan <- string(data)
	})
	//wait m2 to join m1's discovery room
	time.Sleep(time.Second * 10)
	_, isOnline := m1.NodeStatus(m2.NodeAddress)
	if !isOnline {
		t.Error("m2 should online")
		return
	}
	_, isOnline = m2.NodeStatus(m1.NodeAddress)
	if !isOnline {
		t.Error("m1 should online")
		return
	}
	//m1 send use
	err := m1.Send(m2.NodeAddress, []byte("aaa"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "aaa" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "aaa" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	err = m2.Send(m1.NodeAddress, []byte("bbb"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "bbb" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "bbb" {
			t.Errorf("m2 content error %s", txt)
		}
	}
}

func TestSendMessageReLoginOnAnotherServer(t *testing.T) {
	_, m1, m2, _ := newFourTestMatrixTransport()
	m1.db.(*MockDb).addPartner(m2.NodeAddress)
	m2.db.(*MockDb).addPartner(m1.NodeAddress)
	m1.Start()
	//let server sync
	//time.Sleep(time.Second * 6)
	m2.Start()
	m1Chan := make(chan string)
	m2Chan := make(chan string)
	m1.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m1Chan <- string(data)
	})
	m2.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m2Chan <- string(data)
	})
	//wait m2 to join m1's discovery room
	time.Sleep(time.Second * 10)
	//m1 send use
	err := m1.Send(m2.NodeAddress, []byte("aaa"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "aaa" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "aaa" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	m2.Stop()
	time.Sleep(time.Second)
	_, _, cfg3 := getMatrixEnvConfig()
	//m2 relogin on transport03
	m2Again := NewMatrixTransport("m2", m2.key, "other", cfg3)
	if err != nil {
		t.Error(err)
	}
	m2Again.setDB(new(MockDb))
	m2Again.db.(*MockDb).addPartner(m1.NodeAddress)
	m2Again.Start()
	//wait new m2 to join m1's chat room
	time.Sleep(time.Second * 10)
	err = m1.Send(m2Again.NodeAddress, []byte("ccc"))
	if err != nil {
		t.Error(err)
		return
	}
	m2Again.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m2Chan <- string(data)
	})
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "ccc" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "ccc" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	err = m2Again.Send(m1.NodeAddress, []byte("ddd"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "ddd" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "ddd" {
			t.Errorf("m2 content error %s", txt)
		}
	}
}

func TestSendMessageWithoutChannel(t *testing.T) {
	_, m1, m2, _ := newFourTestMatrixTransport()
	m1.Start()
	m2.Start()
	m1Chan := make(chan string)
	m2Chan := make(chan string)
	m1.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m1Chan <- string(data)
	})
	m2.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m2Chan <- string(data)
	})
	//wait m2 to join m1's discovery room
	time.Sleep(time.Second * 10)
	_, isOnline := m1.NodeStatus(m2.NodeAddress)
	if isOnline {
		t.Error("m2 should offline")
		return
	}
	_, isOnline = m2.NodeStatus(m1.NodeAddress)
	if isOnline {
		t.Error("m1 should offline")
		return
	}
	//m1 send use
	err := m1.Send(m2.NodeAddress, []byte("aaa"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "aaa" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 3):
		//should fail
	case txt := <-m2Chan:
		if len(txt) > 0 || txt == "aaa" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	err = m1.Send(m2.NodeAddress, []byte("ccc"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "ccc" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out ")
	case txt := <-m2Chan:
		if txt != "ccc" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	err = m2.Send(m1.NodeAddress, []byte("bbb"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "bbb" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "bbb" {
			t.Errorf("m2 content error %s", txt)
		}
	}
}

func TestSendMessageWithoutChannelAndOfflineOnline(t *testing.T) {
	_, m1, m2, _ := newFourTestMatrixTransport()
	m1.Start()
	m2.Start()
	m1Chan := make(chan string)
	m2Chan := make(chan string)
	m1.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m1Chan <- string(data)
	})
	m2.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m2Chan <- string(data)
	})
	//wait m2 to join m1's discovery room
	time.Sleep(time.Second * 10)
	_, isOnline := m1.NodeStatus(m2.NodeAddress)
	if isOnline {
		t.Error("m2 should offline")
		return
	}
	_, isOnline = m2.NodeStatus(m1.NodeAddress)
	if isOnline {
		t.Error("m1 should offline")
		return
	}
	//m1 send use
	err := m1.Send(m2.NodeAddress, []byte("aaa"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "aaa" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 3):
		//should fail
	case txt := <-m2Chan:
		if len(txt) > 0 || txt == "aaa" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	err = m1.Send(m2.NodeAddress, []byte("ccc"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "ccc" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out ")
	case txt := <-m2Chan:
		if txt != "ccc" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	err = m2.Send(m1.NodeAddress, []byte("bbb"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "bbb" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "bbb" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	m2.Stop()
	time.Sleep(time.Second)
	_, cfg2, _ := getMatrixEnvConfig()
	//m2 relogin on transport03
	m2Again := NewMatrixTransport("m2", m2.key, "other", cfg2)
	m2Again.setDB(new(MockDb))
	m2Again.db.(*MockDb).addPartner(m1.NodeAddress)
	m2Again.Start()
	//wait new m2 to join m1's chat room
	time.Sleep(time.Second * 10)
	err = m1.Send(m2Again.NodeAddress, []byte("ddd"))
	if err != nil {
		t.Error(err)
		return
	}
	m2Again.matrixcli.Syncer.(*gomatrix.DefaultSyncer).OnEventType("m.room.message", func(msg *gomatrix.Event) {
		txt, _ := msg.Body()
		data, _ := base64.StdEncoding.DecodeString(txt)
		m2Chan <- string(data)
	})
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "ddd" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "ddd" {
			t.Errorf("m2 content error %s", txt)
		}
	}
	err = m2Again.Send(m1.NodeAddress, []byte("eee"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m1 receive timeout ")
		return
	case txt := <-m1Chan:
		if txt != "eee" {
			t.Errorf("m1 content error %s", txt)
		}
	}
	select {
	case <-time.After(time.Second * 10):
		t.Error("m2 receive time out")
	case txt := <-m2Chan:
		if txt != "eee" {
			t.Errorf("m2 content error %s", txt)
		}
	}
}
