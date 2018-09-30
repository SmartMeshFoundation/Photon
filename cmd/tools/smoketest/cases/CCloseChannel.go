package cases

import (
	"log"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
)

// CloseChannelTest : test case for close a channel
func CloseChannelTest(env *models.RaidenEnvReader, allowFail bool) {
	caseName := "CloseChannel"
	// prepare data
	var node *models.RaidenNode
	var channels []models.Channel
	for _, n := range env.RaidenNodes {
		channels = env.GetChannelsOfNodeByState(n.AccountAddress, contracts.ChannelStateOpened)
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
	// find the min settle_timeout one
	var channel *models.Channel
	for _, c := range channels {
		if c.SettleTimeout == 35 {
			channel = &c
		}
	}
	if channel == nil {
		channel = &(channels[0])
	}
	// run case
	case1 := &APITestCase{
		CaseName:  caseName,
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "CloseChannel",
			FullURL: node.Host + "/api/1/channels/" + channel.ChannelIdentifier,
			Method:  http.MethodPatch,
			Payload: "{\"state\":\"closed\"}",
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}
