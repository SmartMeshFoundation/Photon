package mobile

import (
	"testing"

	"os"
	"path"

	"encoding/json"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/restful/v1"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestMobile(t *testing.T) {
	mainimpl.GoVersion = "test"
	mainimpl.GitCommit = utils.NewRandomAddress().String()[2:]
	mainimpl.BuildDate = "test"
	nodeAddr := common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa")
	api, err := StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK_REGISTRY"), nil)
	if err != nil {
		t.Error(err)
		return
	}
	if api.Address() != common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa").String() {
		t.Error("address error")
	}
	if testing.Short() {
		return
	}
	var tokens []common.Address
	tokensstr := api.Tokens()
	err = json.Unmarshal([]byte(tokensstr), &tokens)
	if err != nil {
		t.Error(err)
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

	partnerAddr := utils.NewRandomAddress()
	callID, err := api.OpenChannel(partnerAddr.String(), tokens[0].String(), 30, "3")
	if err != nil {
		t.Error(err)
		return
	}
	var channelstr string
	now := time.Now()
	for {
		if time.Since(now) > 40*time.Second {
			break
		}
		channelstr, err = api.GetCallResult(callID)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
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
	callID, err = api.CloseChannel(c.ChannelIdentifier, true)
	if err != nil {
		t.Error(err)
		return
	}
	now = time.Now()
	for {
		if time.Since(now) > 20*time.Second {
			break
		}
		channelstr, err = api.GetCallResult(callID)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		t.Error(err)
		return
	}
	err = json.Unmarshal([]byte(channelstr), &c)
	if err != nil {
		t.Error(err)
		return
	}
	if c.State != contracts.ChannelStateClosed {
		t.Error(err)
		return
	}
	_, err = api.GetCallResult(callID)
	if err.Error() != "not found" {
		t.Error(err)
	}
	api.Stop()
}

func TestFormat(t *testing.T) {
	a := utils.NewRandomAddress()
	t.Logf("a=%q,a=%v,a=%s", a, a, a)
}
