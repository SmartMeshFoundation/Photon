package main

import (
	"log"
	"time"

	"github.com/larspensjo/config"
)

//本地注释：API测试和使用范例
func ApiTest() {
	c, err := config.ReadDefault("./ApiTest.INI")
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
	log.Println("Start Test goRaiden Api")
	start := time.Now()
	//本地注释：测试查询某节点地址
	QueryingNodeAddressTest(Node1Url)
	//本地注释：测试查询某节点所有通道
	QueryingNodeAllChannelsTest(Node1Url)
	//本地注释：测试查询某节点指定通道
	QueryingNodeSpecificChannelTest(Node1Url)
	//本地注释：测试查询系统注册的Token
	QueryingRegisteredTokensTest(Node1Url)
	//本地注释：测试查询节点指定Token有通道的伙伴地址
	QueryingAllPartnersForOneTokensTest(Node1Url)
	//本地注释：测试注册新Token到雷电网
	RegisteringOneTokenTest(Node1Url)
	//本地注释：测试交换Token测试 节点1和节点2,Token 不定，数量2:1
	TokenSwapsTest(Node1Url, Node2Url)
	//本地注释：测试在节点1和节点2 建立Channel，Token为查询到的第一个注册Token
	OpenChannelTest(Node1Url, Node2Url)
	//本地注释：测试关闭节点指定通道
	CloseChannelTest(Node1Url)
	//本地注释：测试Settle节点指定通道
	SettleChannelTest(Node1Url)
	//本地注释：测试向指定通道充值
	Deposit2ChannelTest(Node1Url)
	//本地注释：测试指定Token在雷电网最大限额？
	Connecting2TokenNetworkTest(Node1Url, 2000)
	//本地注释：测试离开指定Token,非常耗时
	LeavingTokenNetworkTest(Node1Url)
	//本地注释：测试查询Token网络连接详情
	QueryingConnectionsDetailsTest(Node1Url)
	//本地注释：测试在节点1和节点2尝试每个Token交易
	InitiatingTransferTest(Node1Url, Node2Url)
	InitiatingTransferTest(Node1Url, Node3Url)
	InitiatingTransferTest(Node1Url, Node4Url)
	InitiatingTransferTest(Node1Url, Node5Url)
	InitiatingTransferTest(Node1Url, Node6Url)
	//本地注释：查询网络事件
	QueryingGeneralNetworkEventsTest(Node1Url)
	//本地注释：查询Token网络事件
	QueryingTokenNetworkEventsTest(Node1Url)
	//本地注释：查询通道事件
	QueryingChannelEventsTest(Node1Url)
	duration := time.Since(start)
	log.Println("Total time used:", duration.Seconds(), " seconds")
}
