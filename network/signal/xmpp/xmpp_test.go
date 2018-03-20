package xmpp

import (
	"fmt"
	"os"
	"testing"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/raiden-network/network/signal/signalshare"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func newpassword(key *ecdsa.PrivateKey) GetCurrentPasswordFunc {
	f1 := func() string {
		pass, _ := signalshare.CreatePassword(key)
		return pass
	}
	return f1
}
func testSdpHandler(from common.Address, sdp string) (mysdp string, err error) {
	log.Trace(fmt.Sprintf("receive sdp request from %s,sdp=%s", utils.APex(from), sdp))
	return sdp, nil
}
func TestNewXmpp(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	key2, _ := crypto.GenerateKey()
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	log.Trace(fmt.Sprintf("addr1=%s,addr=%s\n", addr1.String(), addr2.String()))
	sdp := "test test test"
	x1, err := NewXmpp("119.28.43.121:5222", addr1, newpassword(key1), testSdpHandler, "client1")
	if err != nil {
		t.Error(err)
		return
	}
	x2, err := NewXmpp("119.28.43.121:5222", addr2, newpassword(key2), testSdpHandler, "client2")
	if err != nil {
		t.Error(err)
		return
	}
	err = x1.TryReach(addr2)
	if err != nil {
		t.Error(err)
		return
	}
	sdp2, err := x1.ExchangeSdp(addr2, sdp)
	if err != nil {
		t.Error(err)
		return
	}
	if sdp != sdp2 {
		t.Error(fmt.Sprintf("sdp not equal sdp:%s,sdp2:%s", sdp, sdp2))
	} else {
		t.Log("sdp exchange ok")
	}
	x1.Close()
	x2.Close()
}
func TestNewXmppError(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	sdp := "test test test"
	log.Trace(fmt.Sprintf("addr1 is  %s", addr1.String()))
	x1, err := NewXmpp("119.28.43.121:5222", addr1, newpassword(key1), testSdpHandler, "client1")
	if err != nil {
		t.Error(err)
		return
	}
	err = x1.TryReach(utils.NewRandomAddress())
	if err == nil {
		t.Error(fmt.Sprintf("should not reach"))
		return
	}
	_, err = x1.ExchangeSdp(utils.NewRandomAddress(), sdp)
	if err == nil {
		t.Error(fmt.Sprintf("should fail"))
		return
	}
	x1.Close()
}
