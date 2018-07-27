package accounts

import (
	"encoding/hex"
	"testing"

	"runtime"

	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

func TestDefaultKeyStoreDir(t *testing.T) {
	t.Log(params.DefaultKeyStoreDir())
}

func TestAccountManager(t *testing.T) {
	am := NewAccountManager("../testdata/keystore")
	privkey, err := am.GetPrivateKey(am.Accounts[0].Address, "123")
	if err != nil {
		t.Error(err)
	}
	t.Logf("privkey=0x%s", hex.EncodeToString(privkey))
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
	t.Logf("privkey=0x%s", hex.EncodeToString(privkey))
}
