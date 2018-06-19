package xmpptransport

import (
	"fmt"
	"os"
	"testing"

	"crypto/ecdsa"

	"encoding/hex"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/netshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport/xmpppass"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
		data: make(chan []byte),
	}
}

//DataHandler handles received data
func (t *testDataHandler) DataHandler(from common.Address, data []byte) {
	log.Trace(fmt.Sprintf("%s receive sdp request from %s,data=\n%s", t.name, utils.APex(from), hex.Dump(data)))
	t.data <- data
}
func TestSubscribe(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	key2, _ := crypto.GenerateKey()
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	key3, addr3 := utils.MakePrivateKeyAddress()
	log.Trace(fmt.Sprintf("addr1=%s,addr2=%s,addr3=%s\n", addr1.String(), addr2.String(), addr3.String()))
	x1handler := newTestDataHandler("x1")
	x2handler := newTestDataHandler("x2")
	x1, err := NewConnection(params.DefaultXMPPServer, addr1, &testPasswordGeter{key1}, x1handler, "client1", TypeMobile, make(chan netshare.Status, 10))
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
	x2, err := NewConnection(params.DefaultXMPPServer, addr2, &testPasswordGeter{key2}, x2handler, "client2", TypeOtherDevice, make(chan netshare.Status, 10))
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
	x3, err := NewConnection(params.DefaultXMPPServer, addr3, &testPasswordGeter{key3}, nil, "client3", TypeOtherDevice, make(chan netshare.Status, 10))
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
	x2, err = NewConnection(params.DefaultXMPPServer, addr2, &testPasswordGeter{key2}, x2handler, "client2", TypeOtherDevice, make(chan netshare.Status, 10))
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace("client2 will logout")
	x2.Close()

}
func BenchmarkNewXmpp(b *testing.B) {
	b.N = 10
	for i := 0; i < b.N; i++ {
		key1, _ := crypto.GenerateKey()
		addr1 := crypto.PubkeyToAddress(key1.PublicKey)
		x1, err := NewConnection("139.199.6.114:5222", addr1, &testPasswordGeter{key1}, newTestDataHandler("x1"), "client1", TypeOtherDevice, make(chan netshare.Status, 10))
		if err != nil {
			return
		}
		x1.Close()
	}
}
