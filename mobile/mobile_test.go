package mobile

import (
	"testing"

	"os"
	"path"

	"encoding/json"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/restful/v1"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestMobile(t *testing.T) {
	nodeAddr := common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa")
	api, err := StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if api.Address() != common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa").String() {
		t.Error("address error")
	}
	var tokens []common.Address
	tokensstr := api.Tokens()
	err = json.Unmarshal([]byte(tokensstr), &tokens)
	if err != nil {
		t.Error(err)
		return
	}
	if len(tokensstr) <= 0 {
		t.Errorf("tokens length err")
		return
	}
	var channels []*v1.ChannelData
	channelsstr, err := api.GetChannelList()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(channelsstr)
	err = json.Unmarshal([]byte(channelsstr), &channels)
	if err != nil {
		t.Error(err)
		return
	}
	if len(channels) <= 0 {
		t.Error("channels length error")
		return
	}

	partnerAddr := utils.NewRandomAddress()
	channelstr, err := api.OpenChannel(partnerAddr.String(), tokens[0].String(), 30, "3")
	if err != nil {
		t.Error(err)
		return
	}
	var c v1.ChannelData
	err = json.Unmarshal([]byte(channelstr), &c)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = api.CloseChannel(c.ChannelAddress, true)
	if err != nil {
		t.Error(err)
		return
	}
	//c := channels[0]
	//api.Transfers(c.TokenAddress.String(), c.PartnerAddress.String(), "1", "0", 0)
}

func TestFormat(t *testing.T) {
	a := utils.NewRandomAddress()
	t.Logf("a=%q,a=%v,a=%s", a, a, a)
}
