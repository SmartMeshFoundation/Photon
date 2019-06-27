package xmpptransport

import (
	"fmt"
	"os"
	"testing"

	"github.com/SmartMeshFoundation/Photon/codefortest"

	"github.com/mattn/go-xmpp"

	"crypto/ecdsa"

	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/network/xmpptransport/xmpppass"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

type testPasswordGeter struct {
	key *ecdsa.PrivateKey
}

func (t *testPasswordGeter) GetPassWord() string {
	pass, _ := xmpppass.CreatePassword(t.key)
	return pass
}

type testDataHandler struct {
	name string
	data chan []byte
}

func newTestDataHandler(name string) *testDataHandler {
	return &testDataHandler{
		name: name,
		data: make(chan []byte, 1000),
	}
}

//DataHandler handles received data
func (t *testDataHandler) DataHandler(from common.Address, data []byte) {
	log.Trace(fmt.Sprintf("%s receive sdp request from %s,data=\n%s", t.name, utils.APex(from), string(data)))
	t.data <- data
}
func TestSubscribe(t *testing.T) {
	if testing.Short() {
		return
	}
	db := &codefortest.MockDb{}
	key1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	key2, _ := crypto.GenerateKey()
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	key3, addr3 := utils.MakePrivateKeyAddress()
	log.Trace(fmt.Sprintf("addr1=%s,addr2=%s,addr3=%s\n", addr1.String(), addr2.String(), addr3.String()))
	x1handler := newTestDataHandler("x1")
	x2handler := newTestDataHandler("x2")
	x1, err := NewConnection(params.DefaultXMPPServer, addr1, &testPasswordGeter{key1}, x1handler, "client1", TypeMobile, make(chan netshare.Status, 10), db)
	if err != nil {
		t.Error(err)
		return
	}
	err = x1.SubscribeNeighbour(addr2)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace(fmt.Sprintf("subscribe %s", addr2.String()))

	defer x1.Close()
	_, isOnline, err := x1.IsNodeOnline(addr2)
	if isOnline {
		t.Error("should not online")
		return
	}
	log.Trace("client2 will login")
	x2, err := NewConnection(params.DefaultXMPPServer, addr2, &testPasswordGeter{key2}, x2handler, "client2", TypeOtherDevice, make(chan netshare.Status, 10), db)
	if err != nil {
		t.Error(err)
		return
	}
	//wait notification from server
	time.Sleep(time.Millisecond * 1000)
	deviceType, isOnline, err := x1.IsNodeOnline(addr2)
	if err != nil || !isOnline || deviceType != TypeOtherDevice {
		t.Errorf("should  online,err=%v,isonline=%v,devicetype=%s", err, isOnline, deviceType)
		return
	}
	log.Trace("client2 will logout")
	x2.Close()
	time.Sleep(time.Millisecond * 1000)
	deviceType, isOnline, err = x1.IsNodeOnline(addr2)
	if err != nil || isOnline || deviceType != TypeOtherDevice {
		t.Error("should  offline")
		return
	}
	log.Trace("client3 will login")
	x3, err := NewConnection(params.DefaultXMPPServer, addr3, &testPasswordGeter{key3}, nil, "client3", TypeOtherDevice, make(chan netshare.Status, 10), db)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace("client3 will logout")
	x3.Close()
	err = x1.Unsubscribe(addr2)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 100)
	log.Trace("client2 will relogin")
	x2, err = NewConnection(params.DefaultXMPPServer, addr2, &testPasswordGeter{key2}, x2handler, "client2", TypeOtherDevice, make(chan netshare.Status, 10), db)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace("client2 will logout")
	x2.Close()

}
func BenchmarkNewXmpp(b *testing.B) {
	if testing.Short() {
		return
	}
	db := &codefortest.MockDb{}
	b.N = 10
	for i := 0; i < b.N; i++ {
		key1, _ := crypto.GenerateKey()
		addr1 := crypto.PubkeyToAddress(key1.PublicKey)
		x1, err := NewConnection("139.199.6.114:5222", addr1, &testPasswordGeter{key1}, newTestDataHandler("x1"), "client1", TypeOtherDevice, make(chan netshare.Status, 10), db)
		if err != nil {
			return
		}
		chat := &xmpp.Chat{
			Remote: fmt.Sprintf("%s%s", "leon", nameSuffix),
			Type:   "chat",
			Stamp:  time.Now(),
			Text:   "aaa",
		}
		err = x1.send(chat)
		if err != nil {
			b.Error(err)
		}
		x1.Close()
	}
}
func TestSend(t *testing.T) {
	if testing.Short() {
		return
	}
	db := &codefortest.MockDb{}
	key1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	x1, err := NewConnection(params.DefaultTestXMPPServer, addr1, &testPasswordGeter{key1}, newTestDataHandler("x1"), "client1", TypeOtherDevice, make(chan netshare.Status, 10), db)
	if err != nil {

		return
	}
	for i := 0; i < 10; i++ {

		chat := &xmpp.Chat{
			Remote: fmt.Sprintf("%s%s/Miranda", "leon", nameSuffix),
			Type:   "chat",
			Stamp:  time.Now(),
			Text:   fmt.Sprintf("%d", i),
		}
		err = x1.send(chat)
		if err != nil {
			t.Error(err)
		}

	}
	time.Sleep(time.Second)
	x1.Close()
}
func TestXMPPConnection_SendData(t *testing.T) {
	//key1, _ := crypto.GenerateKey()
	//addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	//key2, _ := crypto.GenerateKey()
	//addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	//log.Trace(fmt.Sprintf("addr1=%s,addr2=%s\n", addr1.String(), addr2.String()))
	//x1handler := newTestDataHandler("x1")
	//x2handler := newTestDataHandler("x2")
	//x1, err := NewConnection(params.DefaultXMPPServer, addr1, &testPasswordGeter{key1}, x1handler, "client1", TypeMobile, make(chan netshare.Status, 10))
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//log.Trace("client2 will login")
	//x2, err := NewConnection(params.DefaultXMPPServer, addr2, &testPasswordGeter{key2}, x2handler, "client2", TypeOtherDevice, make(chan netshare.Status, 10))
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//type sendInfo struct {
	//	start time.Time
	//	end   time.Time
	//	data  string
	//}
	//type receiveInfo struct {
	//	start time.Time
	//	end   time.Time
	//	data  string
	//}
	//sm := make(map[string]*sendInfo)
	//rm := make(map[string]*receiveInfo)
	//log.Trace(fmt.Sprintf("status=%d", x2.status))
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	for {
	//		ri := &receiveInfo{
	//			start: time.Now(),
	//		}
	//		select {
	//		case <-time.After(time.Second * 20):
	//			return
	//		case data := <-x2handler.data:
	//			ri.data = string(data)
	//			ri.end = time.Now()
	//			rm[ri.data] = ri
	//		}
	//	}
	//}()
	//totalTime := time.Now()
	//number := 1000
	//lock := sync.Mutex{}
	//for i := 0; i < number; i++ {
	//	data := fmt.Sprintf("%d", i)
	//	data2 := data
	//	si := &sendInfo{
	//		start: time.Now(),
	//		data:  data,
	//	}
	//	err = x1.SendData(addr2, []byte(data2))
	//	if err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	si.end = time.Now()
	//	lock.Lock()
	//	sm[data2] = si
	//	lock.Unlock()
	//}
	//wg.Wait()
	////t.Logf("sm=%s,rm=%s", utils.StringInterface(sm, 3), utils.StringInterface(rm, 3))
	//for i := 0; i < number; i++ {
	//	data := fmt.Sprintf("%d", i)
	//	si := sm[data]
	//	ri := rm[data]
	//	if si == nil || ri == nil {
	//		t.Errorf("send or receive error for %s,ri=%s,si=%s", data, utils.StringInterface(ri, 2), utils.StringInterface(si, 2))
	//		continue
	//	}
	//	sendTakeTime := si.end.Sub(si.start)
	//	receiveTakeTime := ri.end.Sub(si.end)
	//	receiveWait := ri.end.Sub(ri.start)
	//	t.Logf("message %s send=%s,receive=%s,receiveWait=%s", data, sendTakeTime, receiveTakeTime, receiveWait)
	//}
	//t.Logf("message number=%d,total time=%s", number, time.Now().Sub(totalTime.Add(time.Second*20)))
}
