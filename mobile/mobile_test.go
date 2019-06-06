package mobile

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/Photon/dto"
	v1 "github.com/SmartMeshFoundation/Photon/restful/v1"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	"github.com/stretchr/testify/assert"

	"github.com/SmartMeshFoundation/Photon/notify"

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
	other := NewStrings(1)
	other.Set(0, "--xmpp")
	//other = nil
	api, resultstr := StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK"), other)
	err := dto.ParseResult(resultstr, nil)
	if err != nil {
		t.Error(err)
		return
	}
	defer api.Stop()

	var s string
	err = dto.ParseResult(api.Address(), &s)
	if err != nil {
		t.Error(err)
		return
	}
	ast.EqualValues(s, common.HexToAddress("0x1a9eC3b0b807464e6D3398a59d6b0a369Bf422fA").String())

	var tokens []common.Address
	err = dto.ParseResult(api.Tokens(), &tokens)
	if err != nil {
		t.Error(err)
		return
	}
	var channels []*v1.ChannelData
	err = dto.ParseResult(api.GetChannelList(), &channels)
	//t.Log(resultstr)
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
	var c channeltype.ChannelDataDetail
	resultstr = api.Deposit(nodeAddr.String(), tokens[0].String(), 300, "3", true)
	err = dto.ParseResult(resultstr, &c)
	ast.NotNil(err)
	err = dto.ParseResult(resultstr, &c)
	resultstr = api.Deposit(partnerAddr.String(), tokens[0].String(), 300, "3", true)
	err = dto.ParseResult(resultstr, &c)
	ast.Nil(err)

	channelIdentifier := utils.CalcChannelID(tokens[0], api.api.Photon.Chain.GetRegistryAddress(), nodeAddr, partnerAddr)
	//等待交易被打包
	for i := 0; i < 60; i++ {
		resultstr = api.GetOneChannel(channelIdentifier.String())
		err = dto.ParseResult(resultstr, &c)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
		continue
	}
	if err != nil {
		t.Error(err)
		return
	}

	resultstr = api.CloseChannel(channelIdentifier.String(), true)
	dto.ParseResult(resultstr, &c)
	ast.EqualValues(c.State, channeltype.StateClosing)

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
	api, resultstr := StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK"), nil)
	err := dto.ParseResult(resultstr, nil)
	if err != nil {
		t.Error(err)
		return
	}
	defer api.Stop()
	res := api.Address()
	var resp dto.APIResponse
	json.Unmarshal([]byte(res), &resp)
	var addr string
	json.Unmarshal([]byte(resp.Data), &addr)
	if addr != common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa").String() {
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
}

func testresult() (result string) {
	defer func() {
		log.Info(fmt.Sprintf("result=%s", result))
	}()
	return "ok"
}

func TestResult(t *testing.T) {
	t.Logf("result=%s", testresult())
}

func TestSimpleApi(t *testing.T) {
	ast := assert.New(t)
	if testing.Short() {
		return
	}
	a, err := NewSimpleAPI("/Users/bai/sm/Photon/cmd/photon/.photon", "0x292650fee408320D888e06ed89D938294Ea42f99")
	ast.Nil(err)

	r := a.BalanceAvailabelOnPhoton("0x333")
	v := big.NewInt(0)
	err = dto.ParseResult(r, v)
	ast.Nil(err)
	ast.EqualValues(v, big.NewInt(0))
	r = a.BalanceAvailabelOnPhoton("0x6601F810eaF2fa749EEa10533Fd4CC23B8C791dc")
	fmt.Printf("r=%s", r)
	err = dto.ParseResult(r, v)
	ast.Nil(err)
	ast.True(v.Cmp(big.NewInt(0)) > 0)
	a.Stop()
}

func TestDuplicateStartup(t *testing.T) {
	if testing.Short() {
		return
	}
	mainimpl.GoVersion = "test"
	mainimpl.GitCommit = utils.NewRandomAddress().String()[2:]
	mainimpl.BuildDate = "test"
	nodeAddr := common.HexToAddress("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa")
	other := NewStrings(1)
	other.Set(0, "--xmpp")
	//other = nil
	api, resultstr := StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK"), other)
	err := dto.ParseResult(resultstr, nil)
	if err != nil {
		t.Error(err)
		return
	}
	_, resultstr = StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK"), other)
	err = dto.ParseResult(resultstr, nil)
	if err == nil {
		t.Error("can not startup twice")
		return
	}
	api.Stop()
	api, resultstr = StartUp(nodeAddr.String(), "../testdata/keystore", rpc.TestRPCEndpoint, path.Join(os.TempDir(), utils.RandomString(10)), "../testdata/keystore/pass", "0.0.0.0:5001", "127.0.0.1:40001", "", os.Getenv("TOKEN_NETWORK"), other)
	err = dto.ParseResult(resultstr, nil)
	if err != nil {
		t.Error(err)
		return
	}
}
