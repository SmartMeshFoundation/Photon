package mobile

import (
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/Photon/log"

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
	if testing.Short() {
		return
	}
	mainimpl.GoVersion = "test"
	mainimpl.GitCommit = utils.NewRandomAddress().String()[2:]
	mainimpl.BuildDate = "test"
	nodeAddr := common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa")
	api, err := StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK"), nil)
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
	callID, err := api.Deposit(partnerAddr.String(), tokens[0].String(), 300, "3", true)
	if err != nil {
		t.Error(err)
		return
	}
	var channelstr string
	var crs string
	var cr callResult
	now := time.Now()
	for {
		if time.Since(now) > 60*time.Second {
			break
		}
		time.Sleep(time.Second)
		crs = api.GetCallResult(callID)
		err = json.Unmarshal([]byte(crs), &cr)
		if err != nil {
			t.Error(err)
			return
		}
		if cr.Status == statusDealing {
			continue
		}
		if cr.Status == statusError {
			t.Error(cr.Message)
			return
		}
		channelstr = cr.Message
		break
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
		time.Sleep(time.Second)
		crs = api.GetCallResult(callID)
		err = json.Unmarshal([]byte(crs), &cr)
		if err != nil {
			t.Error(err)
			return
		}
		if cr.Status == statusDealing {
			continue
		}
		if cr.Status == statusError {
			t.Error(cr.Message)
			return
		}
		channelstr = cr.Message
		break
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
	crs = api.GetCallResult(callID)
	err = json.Unmarshal([]byte(crs), &cr)
	if err != nil {
		t.Error(err)
		return
	}
	if cr.Message != "not found" {
		t.Error("should be not found")
	}
	api.Stop()
}

func TestFormat(t *testing.T) {
	a := utils.NewRandomAddress()
	t.Logf("a=%q,a=%v,a=%s", a, a, a)
}

type testHandler struct {
}

func (t *testHandler) OnError(errCode int, failure string) {
	log.Info(fmt.Sprintf("onError errcode=%d,failure=%s", errCode, failure))
}
func (t *testHandler) OnStatusChange(s string) {
	log.Info(fmt.Sprintf("OnStatusChange %s", s))
}

//OnReceivedTransfer  receive a transfer
func (t *testHandler) OnReceivedTransfer(tr string) {
	log.Info(fmt.Sprintf("OnReceivedTransfer %s", tr))
}

//OnSentTransfer a transfer sent success
func (t *testHandler) OnSentTransfer(tr string) {
	log.Info(fmt.Sprintf("OnSentTransfer %s", tr))
}

/* OnNotify get some important message Photon want to notify upper application
level: 0:info,1:warn,2:error
info: type InfoStruct struct {
	Type    int
	Message interface{}
	}
当info.Type=0 表示Message是一个string,1表示Message是TransferStatus
*/
func (t *testHandler) OnNotify(level int, info string) {

}

func TestMobileNotify(t *testing.T) {
	if testing.Short() {
		return
	}
	mainimpl.GoVersion = "test"
	mainimpl.GitCommit = utils.NewRandomAddress().String()[2:]
	mainimpl.BuildDate = "test"
	nodeAddr := common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa")
	api, err := StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK"), nil)
	if err != nil {
		t.Error(err)
		return
	}
	if api.Address() != common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa").String() {
		t.Error("address error")
	}
	sub, err := api.Subscribe(&testHandler{})
	if err != nil {
		t.Error(err)
		return
	}
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second * 3)
			api.NotifyNetworkDown()
		}
	}()
	time.Sleep(time.Minute * 1)
	sub.Unsubscribe()
	api.Stop()
}
