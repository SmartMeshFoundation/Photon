package main

import (
	"log"
	"time"

	"github.com/huamou/config"
)

//API tests and use examples
func APITest() {
	c, err := config.ReadDefault("./apitest.ini")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	Node1Url := c.RdString("NODE1", "api_address", "127.0.0.1:5001")
	Node2Url := c.RdString("NODE2", "api_address", "127.0.0.1:5002")
	Node3Url := c.RdString("NODE3", "api_address", "127.0.0.1:5003")
	Node4Url := c.RdString("NODE4", "api_address", "127.0.0.1:5004")
	Node5Url := c.RdString("NODE5", "api_address", "127.0.0.1:5005")
	Node6Url := c.RdString("NODE6", "api_address", "127.0.0.1:5006")

	Node1Url = "http://" + Node1Url
	Node2Url = "http://" + Node2Url
	Node3Url = "http://" + Node3Url
	Node4Url = "http://" + Node4Url
	Node5Url = "http://" + Node5Url
	Node6Url = "http://" + Node6Url

	log.Println("==============================================================================================")
	log.Println("Start Test goRaiden API")
	start := time.Now()
	//test for querying  a node address
	QueryingNodeAddressTest(Node1Url)
	//test for querying all channels for a node
	QueryingNodeAllChannelsTest(Node1Url)
	//test for querying a node specified channel
	QueryingNodeSpecificChannelTest(Node1Url)
	//test for querying registered Token
	QueryingRegisteredTokensTest(Node1Url)
	//test for querying the Partner address in the channel of special Token
	QueryingAllPartnersForOneTokensTest(Node1Url)
	//test for registering  new Token to Raiden Network
	RegisteringOneTokenTest(Node1Url)
	//test for exchanging the token,token number in each node is indefinite,but the ratio 2:1.
	TokenSwapsTest(Node1Url, Node2Url)
	//test for establishing ch between node 1 and node 2
	OpenChannelTest(Node1Url)
	//test for closing the specified channel for the node
	CloseChannelTest(Node1Url)
	//test for settling the specified channel for the node
	SettleChannelTest(Node1Url)
	//test for depositing  to the specified channel
	Deposit2ChannelTest(Node1Url)
	//test for connecting to a TokenNetwork
	Connecting2TokenNetworkTest(Node1Url, 2000)
	//test for leaving the TokenNetwork
	LeavingTokenNetworkTest(Node1Url)
	//test for querying the details of the Token network connection
	QueryingConnectionsDetailsTest(Node1Url)
	//test for Token transaction between node 1 and other node
	InitiatingTransferTest(Node1Url, Node2Url)
	InitiatingTransferTest(Node1Url, Node3Url)
	InitiatingTransferTest(Node1Url, Node4Url)
	InitiatingTransferTest(Node1Url, Node5Url)
	InitiatingTransferTest(Node1Url, Node6Url)
	//test for querying network events
	QueryingGeneralNetworkEventsTest(Node1Url)
	//test for querying Token network events
	QueryingTokenNetworkEventsTest(Node1Url)
	//test for querying channel event
	QueryingChannelEventsTest(Node1Url)
	duration := time.Since(start)
	log.Println("Total time used:", duration.Seconds(), " seconds")
}
