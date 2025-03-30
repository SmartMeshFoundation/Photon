package network

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/Photon/codefortest"

	"github.com/SmartMeshFoundation/Photon/network/gomatrix"

	"github.com/ethereum/go-ethereum/common"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
)

var testPrivKey *ecdsa.PrivateKey
var testAddress common.Address
var testTrustedServers = []string{
	//"transport01.smartraiden.network",
	"transport01.smartmesh.cn",
	"transport13.smartmesh.cn",
	//"transport03.smartmesh.cn",
}

func init() {
	var err error
	//ALIASFRAGMENT = fmt.Sprintf("testdiscovery-%s", utils.RandomString(10))
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
func getMatrixEnvConfig() (cfg1, cfg2 map[string]string) {
	cfg1 = make(map[string]string)
	cfg2 = make(map[string]string)
	//cfg3 = make(map[string]string)
	cfg1["transport01.smartmesh.cn"] = "http://transport01.smartmesh.cn:8008"
	cfg2["transport13.smartmesh.cn"] = "http://transport13.smartmesh.cn:8008"
	//cfg3["transport03.smartmesh.cn"] = "http://transport03.smartmesh.cn:8008"
	return
}
func newTestMatrixTransport(name string, cfg map[string]string) (m1 *MatrixTransport) {
	key, _ := utils.MakePrivateKeyAddress()
	m1 = NewMatrixTransport(name, key, "other", cfg, &codefortest.MockDb{})
	m1.setTrustServers(testTrustedServers)
	return m1
}

func newFourTestMatrixTransport() (m0, m1, m2 *MatrixTransport) {
	cfg0 := params.TrustMatrixServers
	cfg1, cfg2 := getMatrixEnvConfig()
	m0 = newTestMatrixTransport("m0", cfg0)
	m1 = newTestMatrixTransport("m1", cfg1)
	m2 = newTestMatrixTransport("m2", cfg2)
	//m3 = newTestMatrixTransport("m3", cfg3)
	log.Info(fmt.Sprintf("node0=%s,node1=%s,node2=%s",
		utils.APex2(m0.NodeAddress),
		utils.APex2(m1.NodeAddress),
		utils.APex2(m2.NodeAddress)))
	//utils.APex2(m3.NodeAddress)
	return
}

func TestCreateMatrixTransport(t *testing.T) {
	if testing.Short() {
		return
	}
	m1 := newTestMatrixTransport("mrand", params.TrustMatrixServers)
	m1.Stop()
}
func TestLoginAndJoinDiscoveryRoom(t *testing.T) {
	if testing.Short() {
		return
	}
	cfg1, _ := getMatrixEnvConfig()
	m1 := NewMatrixTransport("test", testPrivKey, "other", cfg1, &codefortest.MockDb{})
	m1.setTrustServers(testTrustedServers)
	log.Trace(fmt.Sprintf("privkey=%s", hex.EncodeToString(crypto.FromECDSA(m1.key))))
	defer m1.Stop()
	m1.Start()
	time.Sleep(time.Second * 1)
}

func TestGetJoinedRoomAlias(t *testing.T) {
	if testing.Short() {
		return
	}
	m1 := NewMatrixTransport("test", testPrivKey, "other", params.TrustMatrixServers, &codefortest.MockDb{})
	m1.setTrustServers(testTrustedServers)
	defer m1.Stop()
	m1.Start()
	//err := m1.loginOrRegister()
	//if err != nil {
	//	panic(err)
	//}
	for m1.matrixcli == nil {
		time.Sleep(time.Second)
	}
	rooms, err := m1.matrixcli.JoinedRooms()
	log.Trace(fmt.Sprintf("rooms=%s", rooms))
	if err != nil {
		panic(err)
	}
	//for _, r := range rooms.JoinedRooms {
	//	var dat map[string]interface{}
	//	m1.matrixcli.StateEvent(r, "m.room.aliases", &dat)
	//	log.Info(fmt.Sprintf("room %s name=%s", r, dat))
	//}
}
func TestInvite(t *testing.T) {
	if testing.Short() {
		return
	}
	_, m1, m2 := newFourTestMatrixTransport()
	m1.Start()
	m2.Start()
	//m3.Start()
	log.Info(fmt.Sprintf(" m1=%s,m2=%s", m1.UserID, m2.UserID))
	for m1.matrixcli == nil {
		time.Sleep(time.Second)
	}
	r, err := m1.matrixcli.CreateRoom(&gomatrix.ReqCreateRoom{
		Invite:        []string{m2.UserID},
		Visibility:    "public",
		Preset:        "public_chat",
		RoomAliasName: fmt.Sprintf("photon_%s_%s_y", utils.APex2(m1.NodeAddress), utils.APex2(m2.NodeAddress)),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("new room=%s", r.RoomID)
	_, err = m1.matrixcli.InviteUser(r.RoomID, &gomatrix.ReqInviteUser{
		UserID: m2.UserID,
	})
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second * 3)
}

func TestSendMessage(t *testing.T) {
	if testing.Short() {
		return
	}
	_, m1, m2 := newFourTestMatrixTransport()
	m1.db.(*codefortest.MockDb).AddPartner(m2.NodeAddress)
	m2.db.(*codefortest.MockDb).AddPartner(m1.NodeAddress)
	m1.Start()
	m2.Start()
	for m1.matrixcli == nil || m1.matrixcli.Syncer == nil || m2.matrixcli == nil || m2.matrixcli.Syncer == nil {
		time.Sleep(time.Second)
	}
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
	////let server sync
	time.Sleep(time.Second * 6)
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
	m1.Stop()
	m2.Stop()
	time.Sleep(time.Second)
}

func TestSendMessageReLoginOnAnotherServer(t *testing.T) {
	if testing.Short() {
		return
	}
	_, m1, m2 := newFourTestMatrixTransport()
	m1.db.(*codefortest.MockDb).AddPartner(m2.NodeAddress)
	m2.db.(*codefortest.MockDb).AddPartner(m1.NodeAddress)
	m1.Start()
	//let server sync
	//time.Sleep(time.Second * 6)
	m2.Start()
	for m1.matrixcli == nil || m1.matrixcli.Syncer == nil || m2.matrixcli == nil || m2.matrixcli.Syncer == nil || !m1.hasDoneStartCheck {
		time.Sleep(time.Second)
	}
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
	_, cfg3 := getMatrixEnvConfig()
	//m2 relogin on transport03
	m2Again := NewMatrixTransport("m2", m2.key, "other", cfg3, &codefortest.MockDb{})
	if err != nil {
		t.Error(err)
	}
	m2Again.setTrustServers(testTrustedServers)
	m2Again.db.(*codefortest.MockDb).AddPartner(m1.NodeAddress)
	m2Again.Start()
	for m2Again.matrixcli == nil || m2Again.matrixcli.Syncer == nil {
		time.Sleep(time.Second)
	}
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
	time.Sleep(time.Second * 10)
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
	if testing.Short() {
		return
	}
	_, m1, m2 := newFourTestMatrixTransport()
	m1.Start()
	m2.Start()
	time.Sleep(time.Second * 6)
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

	//m1 send use
	err := m1.Send(m2.NodeAddress, []byte("aaa"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second * 300):
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
		if txt != "aaa" {
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
	if testing.Short() {
		return
	}
	_, m1, m2 := newFourTestMatrixTransport()
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
		if txt != "aaa" {
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
	_, cfg2 := getMatrixEnvConfig()
	//m2 relogin on transport03
	m2Again := NewMatrixTransport("m2", m2.key, "other", cfg2, &codefortest.MockDb{})
	m2Again.setTrustServers(testTrustedServers)
	m2Again.db.(*codefortest.MockDb).AddPartner(m1.NodeAddress)
	m2Again.Start()

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

func TestVerifySignature(t *testing.T) {
	users := []*gomatrix.UserInfo{
		&gomatrix.UserInfo{
			DisplayName: "10b2-4cdaa7fc3722665d4a9b80ec321b77b95b690a0f0cc91268d4ba5e8dc9a9088b6b69188dba2a3bed5a94ae42204595638eef18d2357653e15418096003e713dc1c",
			UserID:      "@0x10b256b3c83904d524210958fa4e7f9caffb76c6:transport01.smartmesh.cn",
		},
		{
			DisplayName: "10b2-15117d6774e2c98fcc2d09bc1489e74a613d657683f44d0a6b8df8f7fb6d9ee22bf75c4e1e1c09f93e076ecf70050055d0f91c8b13245d42671519ba27a7f13f1c",
			UserID:      "@0x10b256b3c83904d524210958fa4e7f9caffb76c6:transport02.smartmesh.cn",
		},
	}
	for _, u := range users {
		_, err := validateUseridSignature(u)
		if err != nil {
			t.Error(err)
			return
		}
	}

}

func TestRoomTimeLineEvents(t *testing.T) {
	var err error
	if testing.Short() {
		return
	}
	_, m1, m2 := newFourTestMatrixTransport()
	m1.db.(*codefortest.MockDb).AddPartner(m2.NodeAddress)
	m2.db.(*codefortest.MockDb).AddPartner(m1.NodeAddress)
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
	time.Sleep(time.Second * 10)
	for i := 0; i < 10000; i++ {
		//m1 send use
		err = m1.Send(m2.NodeAddress, []byte("aaa"))
		if err != nil {
			t.Error(err)
			return
		}
	}
	time.Sleep(time.Second * 30)
	m1.Stop()
	m2.Stop()

	//重新登录,看看事件有没有问题
	cfg1, cfg2 := getMatrixEnvConfig()
	m1Again := NewMatrixTransport("m1", m1.key, "other", cfg1, &codefortest.MockDb{})
	if err != nil {
		t.Error(err)
	}
	m1Again.setTrustServers(testTrustedServers)

	m2Again := NewMatrixTransport("m2", m2.key, "other", cfg2, &codefortest.MockDb{})
	if err != nil {
		t.Error(err)
	}
	m2Again.setTrustServers(testTrustedServers)
	time.Sleep(time.Second * 20)
	//看看下次获取的事件信息
	m1Again.Start()
	m2Again.Start()
	time.Sleep(time.Second * 10)
}

func TestMatrixTransport_splitAlias(t *testing.T) {
	prefix, isChannel, addr1, addr2, err := splitRoomAlias("#photon_y_37bd76c0187ebc21e3fd3d474d83810bb495a518_4533775cfd13a2b07bf910c04d2038fd028ff73c:transport13.smartmesh.cn")
	if err != nil {
		t.Error(err)
		return
	}
	if prefix != "photon" {
		t.Error("prefix")
		return
	}
	if isChannel != "y" {
		t.Error("channel")
		return
	}
	if addr1 != common.HexToAddress("0x37bd76c0187ebc21e3fd3d474d83810bb495a518") {
		t.Error("addr1")
		return
	}
	if addr2 != common.HexToAddress("0x4533775cfd13a2b07bf910c04d2038fd028ff73c") {
		t.Error("addr2")
		return
	}
	_, _, _, _, err = splitRoomAlias("")
	if err == nil {
		t.Error("must fail")
		return
	}
	_, _, _, _, err = splitRoomAlias("#photon_ropsten_discovery:transport01.smartmesh.cn")
	if err == nil {
		t.Error("must fail")
		return
	}
}

func TestLeaveUselessRoom(t *testing.T) {
	if testing.Short() {
		return
	}
	cfg1, _ := getMatrixEnvConfig()
	m1 := NewMatrixTransport("test", testPrivKey, "other", cfg1, &codefortest.MockDb{})
	m1.setTrustServers(testTrustedServers)
	log.Trace(fmt.Sprintf("privkey=%s", hex.EncodeToString(crypto.FromECDSA(m1.key))))
	defer m1.Stop()
	m1.Start()
	m1.leaveUselessRoom()
}

func TestPresenceListFunction(t *testing.T) { //注册次测试过程需要手动修改sync函数里面的presence参数为空
	if testing.Short() {
		return
	}
	//same server
	delete(params.TrustMatrixServers, "transport13.smartmesh.cn")
	cfg := params.TrustMatrixServers
	trans01 := newTestMatrixTransport("t1", cfg)
	trans02 := newTestMatrixTransport("t2", cfg)
	trans03 := newTestMatrixTransport("t3", cfg)
	trans01.Start()
	trans02.Start()
	trans03.Start()

	time.Sleep(5 * time.Second)

	//设置presence list
	trans01.matrixcli.PostPresenceList(&gomatrix.ReqPresenceList{
		Invite: []string{trans02.UserID, trans03.UserID},
	})
	//trans01读取列表
	reps01, err := trans01.matrixcli.GetPresenceList(trans01.UserID)
	if err != nil {
		t.Error("t1 GetPresenceList(t2/t3) failed", err)
		return
	}
	for _, content1 := range reps01 {
		log.Trace(fmt.Sprintf("UserID:%s ,presence=%s", content1.UserID, content1.Presence))
		if content1.Presence != "online" {
			t.Error("get t2/t3 presence err,Expected=online,but =offline")
			return
		}
	}
	log.Trace(fmt.Sprintf(" t1 GetPresenceList=%s", utils.StringInterface(reps01, 8)))
	log.Trace("====================================================\n")
	//trans02下线/上线交互
	var setPresence = ""
	var errTimes = 0
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		if i%2 == 0 {
			setPresence = "offline"
		} else {
			setPresence = "online"
		}

		err2 := trans02.matrixcli.SetPresenceState(&gomatrix.ReqPresenceUser{
			Presence:  setPresence,
			StatusMsg: "other",
		})
		if err2 != nil {
			t.Error("test", i, ": SetPresenceState(offline) failed", err2)
			err = err2
			return
		}

		reps02, err2 := trans01.matrixcli.GetPresenceList(trans01.UserID)
		if err2 != nil {
			t.Error("test", i, ": GetPresenceList(t2/t3) failed", err2)
			err = err2
			return
		}
		for _, content2 := range reps02 {
			log.Trace(fmt.Sprintf("test %d : UserID:%s ,presence=%s", i, content2.UserID, content2.Presence))
			if content2.UserID == trans02.UserID && content2.Presence != setPresence {
				errTimes++
			}
		}
	}
	if errTimes > 3 {
		t.Error("test for changing presence on one server failed")
	}

	log.Trace("====================================================\n")

	time.Sleep(time.Second)
	//节点处于不同的服务器
	trans02.Stop()
	trans03.Stop()
	//移除trans01对trans02 trans03的监控列表
	trans01.matrixcli.PostPresenceList(&gomatrix.ReqPresenceList{
		Drop: []string{trans02.UserID, trans03.UserID},
	})

	params.TrustMatrixServers["transport13.smartmesh.cn"] = "http://transport13.smartmesh.cn:8008"
	delete(params.TrustMatrixServers, "transport01.smartmesh.cn")
	//params.TrustMatrixServers["transport01.smartmesh.cn"]="http://transport01.smartmesh.cn:8008"
	cfg1 := params.TrustMatrixServers
	trans04 := newTestMatrixTransport("t4", cfg1)
	trans05 := newTestMatrixTransport("t5", cfg1)
	trans04.Start()
	trans05.Start()
	time.Sleep(time.Second * 5)
	//设置presence list
	trans01.matrixcli.PostPresenceList(&gomatrix.ReqPresenceList{
		Invite: []string{trans04.UserID, trans05.UserID},
	})
	var setPresence1 = ""
	var errTimes1 = 0
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		if i%2 == 0 {
			setPresence1 = "offline"
		} else {
			setPresence1 = "online"
		}

		err = trans04.matrixcli.SetPresenceState(&gomatrix.ReqPresenceUser{
			Presence:  setPresence1,
			StatusMsg: "other",
		})
		if err != nil {
			t.Error("test", i, ": SetPresenceState(offline) failed", err)
			return
		}

		reps02, err := trans01.matrixcli.GetPresenceList(trans01.UserID)
		if err != nil {
			t.Error("test", i, ": GetPresenceList(t2/t3) failed", err)
			return
		}
		for _, content2 := range reps02 {
			log.Trace(fmt.Sprintf("test %d : UserID:%s ,presence=%s", i, content2.UserID, content2.Presence))
			if content2.UserID == trans04.UserID && content2.Presence != setPresence1 {
				errTimes1++
			}
		}
	}
	if errTimes1 > 3 {
		t.Error("test for changing presence on different server failed")
	}

}
