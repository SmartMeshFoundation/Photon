package accounts

import (
	"encoding/hex"
	"testing"

	"runtime"

	"github.com/SmartMeshFoundation/Photon/utils"
)

func TestDefaultKeyStoreDir(t *testing.T) {
	t.Log(utils.DefaultKeyStoreDir())
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
