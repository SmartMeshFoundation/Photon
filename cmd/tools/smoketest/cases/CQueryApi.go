package cases

import (
	"log"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
)

// query api default timeout
var queryTimeOut = time.Second * 30

// QueryNodeAddressTest :
func QueryNodeAddressTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "QueryNodeAddress",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryNodeAddress",
			FullURL: env.RandomNode().Host + "/api/1/address",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryRegisteredTokenTest :
func QueryRegisteredTokenTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "QueryRegisteredToken",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryRegisteredToken",
			FullURL: env.RandomNode().Host + "/api/1/tokens",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryAllPartnersForOneTokenTest :
func QueryAllPartnersForOneTokenTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "QueryAllPartnersForOneToken",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryAllPartnersForOneToken",
			FullURL: env.RandomNode().Host + "/api/1/tokens/" + env.RandomToken().Address + "/partners",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryNodeAllChannelsTest :
func QueryNodeAllChannelsTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "QueryNodeAllChannels",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryNodeAllChannels",
			FullURL: env.RandomNode().Host + "/api/1/channels",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryNodeSpecificChannelTest :
func QueryNodeSpecificChannelTest(env *models.RaidenEnvReader, allowFail bool) {
	// prepare data for this case
	var node *models.RaidenNode
	var channels []models.Channel
	for _, n := range env.RaidenNodes {
		channels = env.GetChannelsOfNode(n.AccountAddress)
		if len(channels) > 0 {
			node = n
			break
		}
	}
	if channels == nil || len(channels) == 0 {
		log.Println("Case [QueryNodeSpecificChannel] FAILED because no suitable env !!!")
		Logger.Println("Case [QueryNodeSpecificChannel] FAILED because no suitable env !!!")
		if !allowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		return
	}
	// run case
	case1 := &APITestCase{
		CaseName:  "QueryNodeSpecificChannel",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryNodeSpecificChannel",
			FullURL: node.Host + "/api/1/channels/" + channels[0].ChannelAddress,
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryGeneralNetworkEventsTest :
func QueryGeneralNetworkEventsTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "QueryGeneralNetworkEvents",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryGeneralNetworkEvents",
			FullURL: env.RandomNode().Host + "/api/1/events/network",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryTokenNetworkEventsTest :
func QueryTokenNetworkEventsTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "QueryTokenNetworkEvents",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryTokenNetworkEvents",
			FullURL: env.RandomNode().Host + "/api/1/events/tokens/" + env.RandomToken().Address,
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryChannelEventsTest :
func QueryChannelEventsTest(env *models.RaidenEnvReader, allowFail bool) {
	// prepare data for this case
	var node *models.RaidenNode
	var channels []models.Channel
	for _, n := range env.RaidenNodes {
		channels = env.GetChannelsOfNode(n.AccountAddress)
		if len(channels) > 0 {
			node = n
			break
		}
	}
	if channels == nil || len(channels) == 0 {
		Logger.Println("Current env can not afford this case !!!")
		if !allowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
	}
	if channels == nil || len(channels) == 0 {
		log.Println("Case [QueryNodeSpecificChannel] FAILED because no suitable env !!!")
		Logger.Println("Case [QueryNodeSpecificChannel] FAILED because no suitable env !!!")
		if !allowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		return
	}
	// run case
	case1 := &APITestCase{
		CaseName:  "QueryChannelEvents",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryChannelEvents",
			FullURL: node.Host + "/api/1/events/channels/" + channels[0].ChannelAddress + "?from_block=1",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}
