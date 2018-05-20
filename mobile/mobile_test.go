package mobile

import (
	"testing"

	"os"
	"path"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestMobile(t *testing.T) {
	go StartUp("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa", "../testdata/keystore", "ws://127.0.0.1:18546" /*rpc.TestRpcEndpoint,*/, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "127.0.0.1:5009")
	time.Sleep(time.Second * 5)
}

func TestFormat(t *testing.T) {
	a := utils.NewRandomAddress()
	t.Logf("a=%q,a=%v,a=%s", a, a, a)
}
