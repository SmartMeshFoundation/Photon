package xmpptransport

import (
	"fmt"
	"os"
	"testing"

	"crypto/ecdsa"

	"encoding/hex"

	"time"

	"bytes"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
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
func TestNewXmpp(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	key2, _ := crypto.GenerateKey()
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	log.Trace(fmt.Sprintf("addr1=%s,addr=%s\n", addr1.String(), addr2.String()))
	x1handler := newTestDataHandler("x1")
	x2handler := newTestDataHandler("x2")
	x1, err := NewConnection(params.DefaultXMPPServer, addr1, &testPasswordGeter{key1}, x1handler, "client1", TypeMobile, make(chan Status, 10))
	if err != nil {
		t.Error(err)
		return
	}
	deviceType, isOnline, err := x1.IsNodeOnline(addr2)
	if err != nil {
		t.Error(err)
		return
	}
	if isOnline {
		t.Error("should offline")
		return
	}
	defer x1.Close()
	x2, err := NewConnection(params.DefaultXMPPServer, addr2, &testPasswordGeter{key2}, x2handler, "client2", TypeOtherDevice, make(chan Status, 10))
	if err != nil {
		t.Error(err)
		return
	}
	defer x2.Close()
	deviceType, isOnline, err = x1.IsNodeOnline(addr2)
	if err != nil {
		t.Error(err)
		return
	}
	if !isOnline {
		t.Error("should online")
		return
	}
	if deviceType != TypeOtherDevice {
		t.Error("type error")
		return
	}
	deviceType, isOnline, err = x2.IsNodeOnline(addr1)
	if deviceType != TypeMobile {
		t.Error("type error")
		return
	}
	err = x1.SendData(addr2, []byte("abc"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second):
		t.Error("recevie timeout")
	case data := <-x2handler.data:
		if !bytes.Equal(data, []byte("abc")) {
			t.Error("not equal")
		}
	}
}

func BenchmarkNewXmpp(b *testing.B) {
	b.N = 10
	for i := 0; i < b.N; i++ {
		key1, _ := crypto.GenerateKey()
		addr1 := crypto.PubkeyToAddress(key1.PublicKey)
		x1, err := NewConnection("139.199.6.114:5222", addr1, &testPasswordGeter{key1}, newTestDataHandler("x1"), "client1", TypeOtherDevice, make(chan Status, 10))
		if err != nil {
			return
		}
		x1.Close()
	}
}
