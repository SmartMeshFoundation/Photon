package cases

import (
	"log"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon/utils"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
)

// Deposit2ChannelTest : test case for deposit to channel
func Deposit2ChannelTest(env *models.PhotonEnvReader, allowFail bool) {

	testDepositToNotExistChannel(env, allowFail)
	testDepositToChannelByState(env, allowFail, contracts.ChannelStateOpened, 200)
	testDepositToChannelByState(env, true, contracts.ChannelStateClosed, 408)
	testDepositToChannelByState(env, true, contracts.ChannelStateSettledOrNotExist, 408)

}

func testDepositToNotExistChannel(env *models.PhotonEnvReader, allowFail bool) {
	payload := newOpenChannelPayload(utils.NewRandomAddress().String(), utils.NewRandomAddress().String(), 30, 0, false)
	case1 := &APITestCase{
		CaseName:  "Deposit to not-exist channel",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: " Deposit2Channel",
			FullURL: env.RandomNode().Host + "/api/1/deposit",
			Method:  http.MethodPut,
			Payload: string(payload),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 409,
	}
	case1.Run()
}

func testDepositToChannelByState(env *models.PhotonEnvReader, allowFail bool, channelState int, targetStatusCode int) {
	// prepare data
	caseName := fmt.Sprintf("Deposit to %d channel", channelState)
	var node *models.PhotonNode
	var channels []models.Channel
	for _, n := range env.PhotonNodes {
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
	c := channels[0]
	payload := newOpenChannelPayload(c.PartnerAddress, c.TokenAddress, 5, 0, false)
	// run case
	case1 := &APITestCase{
		CaseName:  caseName,
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "Deposit2Channel",
			FullURL: node.Host + "/api/1/deposit",
			Method:  http.MethodPut,
			Payload: string(payload),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: targetStatusCode,
	}
	case1.Run()
}
