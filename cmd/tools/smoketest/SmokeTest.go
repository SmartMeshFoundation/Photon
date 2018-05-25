package main

import (
	"math/rand"
	"time"

	"log"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/cases"
)

// SmokeTest of smartraiden api
func SmokeTest() {
	log.Println("SmokeTest start ...")
	caseLogger := cases.Logger
	caseLogger.Println("==============================================================================================")
	caseLogger.Println("Start Test goRaiden Api")
	caseLogger.Println("==================================================")
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	runSmokeCases()

	duration := time.Since(start)
	caseLogger.Println("Total time used:", duration.Seconds(), " seconds")
	if len(cases.FailCases) > 0 {
		caseLogger.Printf(" %d Fail Cases :", len(cases.FailCases))
		for _, c := range cases.FailCases {
			caseLogger.Println(c)
		}
	}
	log.Println("SmokeTest done. Check logs at " + cases.LogPath)
}

func runSmokeCases() {
	// cases about query api
	cases.QueryNodeAddressTest(env, allowFail)
	cases.QueryRegisteredTokenTest(env, allowFail)
	cases.QueryAllPartnersForOneTokenTest(env, allowFail)
	cases.QueryNodeAllChannelsTest(env, allowFail)
	cases.QueryNodeSpecificChannelTest(env, allowFail)
	cases.QueryGeneralNetworkEventsTest(env, allowFail)
	cases.QueryTokenNetworkEventsTest(env, allowFail)
	cases.QueryChannelEventsTest(env, allowFail)

	// cases about transfer
	cases.InitiatingTransferTest(env, allowFail)
	cases.TokenSwapsTest(env, allowFail)

	// cases about token
	cases.RegisteringTokenTest(env, allowFail)
	cases.Connecting2TokenNetworkTest(env, allowFail)
	cases.LeavingTokenNetworkTest(env, allowFail)

	// cases about channel
	cases.OpenChannelTest(env, allowFail)
	env.RefreshChannels()
	cases.CloseChannelTest(env, allowFail)
	env.RefreshChannels()
	cases.SettleChannelTest(env, allowFail)
	env.RefreshChannels()
	cases.Deposit2ChannelTest(env, allowFail)
}
