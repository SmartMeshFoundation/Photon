package mobile

import (
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	"github.com/stretchr/testify/assert"

	"github.com/SmartMeshFoundation/Photon/notify"
	"github.com/SmartMeshFoundation/Photon/restful/v1"

	"github.com/SmartMeshFoundation/Photon/log"

	"os"
	"path"

	"encoding/json"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestMobile(t *testing.T) {
	ast := assert.New(t)
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
	//t.Log(channelsstr)
	err = json.Unmarshal([]byte(channelsstr), &channels)
	if err != nil {
		t.Error(err)
		return
	}
	sub, err := api.Subscribe(newTestHandler())
	if err != nil {
		t.Error(err)
		return
	}
	defer sub.Unsubscribe()
	partnerAddr := utils.NewRandomAddress()
	channelstr, err := api.Deposit(partnerAddr.String(), tokens[0].String(), 300, "3", true)
	if err != nil {
		t.Error(err)
		return
	}
	ast.EqualValues(channelstr, "")
	channelIdentifier := utils.CalcChannelID(tokens[0], api.api.Photon.Chain.GetRegistryAddress(), nodeAddr, partnerAddr)
	//等待交易被打包
	for i := 0; i < 60; i++ {
		channelstr, err = api.GetOneChannel(channelIdentifier.String())
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	ast.NotEqual(channelstr, "")
	var c channeltype.ChannelDataDetail
	err = json.Unmarshal([]byte(channelstr), &c)
	if err != nil {
		t.Error(err)
		return
	}

	channelstr, err = api.CloseChannel(channelIdentifier.String(), true)
	if err != nil {
		t.Error(err)
		return
	}
	err = json.Unmarshal([]byte(channelstr), &c)
	if err != nil {
		t.Error(err)
		return
	}
	ast.EqualValues(c.State, channeltype.StateClosing)
	api.Stop()
}

func TestFormat(t *testing.T) {
	a := utils.NewRandomAddress()
	t.Logf("a=%q,a=%v,a=%s", a, a, a)
}

type testHandler struct {
	lastNotify map[int]interface{}
}

func newTestHandler() *testHandler {
	return &testHandler{
		lastNotify: make(map[int]interface{}),
	}
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
当info.Type=0 表示Message是一个string,1表示Message是TransferStatus,2 表示channel callid,3 表示channel status
*/
func (t *testHandler) OnNotify(level int, info string) {
	log.Info(fmt.Sprintf("notify leve=%d,info=%s", level, info))
	if level == 0 {
		var is notify.InfoStruct
		err := json.Unmarshal([]byte(info), &is)
		if err != nil {
			log.Error("unmarshal error %s", err)
		} else {
			t.lastNotify[is.Type] = is.Message
		}
	}
}

func (t *testHandler) clearHistory() {
	for k := range t.lastNotify {
		delete(t.lastNotify, k)
	}
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
	sub, err := api.Subscribe(newTestHandler())
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
