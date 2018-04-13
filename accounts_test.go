package smartraiden

import (
	"encoding/hex"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

func TestDefaultKeyStoreDir(t *testing.T) {
	t.Log(params.DefaultKeyStoreDir())
}

func TestAccountManager(t *testing.T) {
	am := NewAccountManager("testdata/keystore")
	privkey, err := am.GetPrivateKey(am.Accounts[0].Address, "123")
	if err != nil {
		t.Error(err)
	}
	t.Logf("privkey=0x%s", hex.EncodeToString(privkey))
}
