package cases

import (
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
)

// Deposit2ChannelTest : test case for deposit to channel
func Deposit2ChannelTest(env *models.RaidenEnvReader, allowFail bool) {

	testDepositToNotExistChannel(env, allowFail)
	testDepositToChannelByState(env, allowFail, "opened", 200)
	testDepositToChannelByState(env, allowFail, "closed", 408)
	testDepositToChannelByState(env, allowFail, "settled", 408)

}

func testDepositToNotExistChannel(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "Deposit to not-exist channel",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: " Deposit2Channel",
			FullURL: env.RandomNode().Host + "/api/1/channels/0xffffffffffffffffffffffffffffffffffffffff",
			Method:  http.MethodPatch,
			Payload: "{\"balance\":5}",
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 409,
	}
	case1.Run()
}

func testDepositToChannelByState(env *models.RaidenEnvReader, allowFail bool, channelState string, targetStatusCode int) {
	// prepare data
	caseName := fmt.Sprintf("Deposit to %s channel", channelState)
	var node *models.RaidenNode
	var channels []models.Channel
	for _, n := range env.RaidenNodes {
		channels = env.GetChannelsOfNodeByState(n.AccountAddress, channelState)
		if len(channels) > 0 {
			node = n
			break
		}
	}
	if channels == nil || len(channels) == 0 {
		log.Printf("Case [%-40s] FAILED because no suitable env !!!", caseName)
		Logger.Printf("Case [%-40s] FAILED because no suitable env !!!", caseName)
		if !allowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		return
	}
	// run case
	case1 := &APITestCase{
		CaseName:  caseName,
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "Deposit2Channel",
			FullURL: node.Host + "/api/1/channels/" + channels[0].ChannelAddress,
			Method:  http.MethodPatch,
			Payload: "{\"balance\":5}",
			Timeout: time.Second * 180,
		},
		TargetStatusCode: targetStatusCode,
	}
	case1.Run()
}
