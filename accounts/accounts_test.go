package accounts

import (
	"encoding/hex"
	"testing"

	"runtime"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDefaultKeyStoreDir(t *testing.T) {
	t.Log(utils.DefaultKeyStoreDir())
}

func TestAccountManager(t *testing.T) {
	am := NewAccountManager("/home/chuck/code/run/smc/testnet/keystore")
	privkey, err := am.GetPrivateKey(am.Accounts[2].Address, "123")
	if err != nil {
		t.Error(err)
	}
	t.Logf("privkey=0x%s", hex.EncodeToString(privkey))
	p, _ := crypto.ToECDSA(privkey)
	t.Logf("pubKey=0x%s", crypto.PubkeyToAddress(p.PublicKey).String())
}

func TestFininalize(t *testing.T) {
	for i := 0; i < 10; i++ {
		testNewAccount(t)
		runtime.GC()
	}
}
func testNewAccount(t *testing.T) {
	am := NewAccountManager("../testdata/keystore")
	privkey, err := am.GetPrivateKey(am.Accounts[0].Address, "123")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s privkey=0x%s", am.Accounts[0].Address.String(), hex.EncodeToString(privkey))
}
