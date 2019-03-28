package cases

import (
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
)

// query api default timeout
var queryTimeOut = time.Second * 30

// QueryNodeAddressTest :
func QueryNodeAddressTest(env *models.PhotonEnvReader, allowFail bool) {
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
func QueryRegisteredTokenTest(env *models.PhotonEnvReader, allowFail bool) {
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
func QueryAllPartnersForOneTokenTest(env *models.PhotonEnvReader, allowFail bool) {
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
func QueryNodeAllChannelsTest(env *models.PhotonEnvReader, allowFail bool) {
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
func QueryNodeSpecificChannelTest(env *models.PhotonEnvReader, allowFail bool) {
	// prepare data for this case
	var node *models.PhotonNode
	var channels []models.Channel
	for _, n := range env.PhotonNodes {
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
			FullURL: node.Host + "/api/1/channels/" + channels[0].ChannelIdentifier,
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// QueryGeneralNetworkEventsTest :
func QueryGeneralNetworkEventsTest(env *models.PhotonEnvReader, allowFail bool) {
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
func QueryTokenNetworkEventsTest(env *models.PhotonEnvReader, allowFail bool) {
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
func QueryChannelEventsTest(env *models.PhotonEnvReader, allowFail bool) {
	// prepare data for this case
	var node *models.PhotonNode
	var channels []models.Channel
	for _, n := range env.PhotonNodes {
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
		log.Println("Case [QueryChannelEventsTest] FAILED because no suitable env !!!")
		Logger.Println("Case [QueryChannelEventsTest] FAILED because no suitable env !!!")
		return
	}

	// run case
	case1 := &APITestCase{
		CaseName:  "QueryChannelEvents",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "QueryChannelEvents",
			FullURL: node.Host + "/api/1/events/channels/" + channels[0].ChannelIdentifier + "?from_block=1",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// GetSentTransfersTest :
func GetSentTransfersTest(env *models.PhotonEnvReader, allowFail bool) {
	// prepare data for this case
	var node *models.PhotonNode
	var channels []models.Channel
	for _, n := range env.PhotonNodes {
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
		log.Println("Case [GetSentTransfersTest] FAILED because no suitable env !!!")
		Logger.Println("Case [GetSentTransfersTest] FAILED because no suitable env !!!")
		return
	}
	// run case
	case1 := &APITestCase{
		CaseName:  "GetSentTransferDetails",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "GetSentTransferDetails",
			FullURL: node.Host + "/api/1/querysenttransfer",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// GetReceivedTransfersTest :
func GetReceivedTransfersTest(env *models.PhotonEnvReader, allowFail bool) {
	// prepare data for this case
	var node *models.PhotonNode
	var channels []models.Channel
	for _, n := range env.PhotonNodes {
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
		log.Println("Case [GetReceivedTransfersTest] FAILED because no suitable env !!!")
		Logger.Println("Case [GetReceivedTransfersTest] FAILED because no suitable env !!!")
		return
	}
	// run case
	case1 := &APITestCase{
		CaseName:  "GetReceivedTransfers",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "GetReceivedTransfers",
			FullURL: node.Host + "/api/1/queryreceivedtransfer",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// GetRandomSecretTest :
func GetRandomSecretTest(env *models.PhotonEnvReader, allowFail bool) {
	// run case
	case1 := &APITestCase{
		CaseName:  "GetRandomSecre",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "GetRandomSecre",
			FullURL: env.PhotonNodes[0].Host + "/api/1/secret",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// GetBalanceByTokenAddressTest :
func GetBalanceByTokenAddressTest(env *models.PhotonEnvReader, allowFail bool) {
	// run case
	case1 := &APITestCase{
		CaseName:  "GetBalanceByTokenAddress",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "GetBalanceByTokenAddress",
			FullURL: env.PhotonNodes[0].Host + "/api/1/balance",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}

// ChannelFor3rdPartyTest :
func ChannelFor3rdPartyTest(env *models.PhotonEnvReader, allowFail bool) {
	// run case
	node := env.PhotonNodes[0]
	channels := env.GetChannelsOfNodeByState(node.AccountAddress, contracts.ChannelStateOpened)
	if channels == nil || len(channels) == 0 {
		Logger.Println("Current env can not afford this case !!!")
		if !allowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		log.Println("Case [ChannelFor3rdPartyTest] FAILED because no suitable env !!!")
		Logger.Println("Case [ChannelFor3rdPartyTest] FAILED because no suitable env !!!")
		return
	}
	case1 := &APITestCase{
		CaseName:  "ChannelFor3rdParty",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "ChannelFor3rdParty",
			FullURL: fmt.Sprintf("%s/api/1/thirdparty/%s/%s", node.Host, channels[0].ChannelIdentifier, env.PhotonNodes[1].AccountAddress),
			Method:  http.MethodGet,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}
