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
	// 1 cases about query api
	cases.QueryNodeAddressTest(env, allowFail)
	cases.QueryRegisteredTokenTest(env, allowFail)
	cases.QueryAllPartnersForOneTokenTest(env, allowFail)
	cases.QueryNodeAllChannelsTest(env, allowFail)
	cases.QueryNodeSpecificChannelTest(env, allowFail)
	cases.QueryGeneralNetworkEventsTest(env, allowFail)
	cases.QueryTokenNetworkEventsTest(env, allowFail)
	cases.QueryChannelEventsTest(env, allowFail)

	// cases about token
	cases.RegisteringTokenTest(env, allowFail)
	cases.Connecting2TokenNetworkTest(env, allowFail)
	cases.LeavingTokenNetworkTest(env, allowFail)

	// cases about channel
	cases.OpenChannelTest(env, allowFail)
	cases.Deposit2ChannelTest(env, allowFail)
	//CCloseChannel(&node, "Case12")
	//CSettleChannel(&node, "Case13")

	// cases about transfer
	//CInitiatingTransfer(&nodes[1], &nodes[3], "Case14", 1)
	//CInitiatingTransfer(&nodes[1], &nodes[5], "Case15", 2)
	//
	//CTokenSwaps(&nodes[2], &nodes[1], "Case16", 2, 1, "taker")
	//CTokenSwaps(&nodes[1], &nodes[2], "Case16", 1, 2, "maker")
	//
	//CTokenSwaps(&nodes[5], &nodes[1], "Case17", 2, 1, "taker")
	//CTokenSwaps(&nodes[1], &nodes[5], "Case17", 1, 2, "maker")

}
