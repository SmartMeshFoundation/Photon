package mobile

import (
	"testing"

	"os"
	"path"

	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/utils"
)

func TestMobile(t *testing.T) {
	api, err := MobileStartUp("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa", "d:\\privnet\\keystore\\", rpc.TestRpcEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "c:\\pass")
	if err != nil {
		t.Error(err)
		return
	}
	//time.Sleep(2 * time.Second)
	chs, err := api.GetChannelList()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(chs)
}
