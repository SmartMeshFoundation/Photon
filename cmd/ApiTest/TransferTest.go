package main

import (
	"fmt"
	"github.com/larspensjo/config"
)

//本地注释：详细测试交易
func TransferTest(NewTokenName string) (code int) {

	fmt.Println("==============================================================================================")
	fmt.Println("Start TransferTest")

	c, err := config.ReadDefault("./ApiTest.INI")

	if err != nil {
		fmt.Println("config.ReadDefault error:", err)
		return
	}

	Node1Url, err := c.String("NOTE1", "api_address")
	Node2Url, err := c.String("NOTE2", "api_address")
	Node3Url, err := c.String("NOTE3", "api_address")
	Node4Url, err := c.String("NOTE4", "api_address")
	Node5Url, err := c.String("NOTE5", "api_address")
	Node6Url, err := c.String("NOTE6", "api_address")

	Node1Url = "http://" + Node1Url
	Node2Url = "http://" + Node2Url
	Node3Url = "http://" + Node3Url
	Node4Url = "http://" + Node4Url
	Node5Url = "http://" + Node5Url
	Node6Url = "http://" + Node6Url

	//本地注释：新Token账户充值

	//本地注释：获取节点地址
	Node1Address, Status, err := QueryingNodeAddress(Node1Url)
	Node2Address, Status, err := QueryingNodeAddress(Node2Url)
	Node3Address, Status, err := QueryingNodeAddress(Node3Url)
	Node4Address, Status, err := QueryingNodeAddress(Node4Url)
	Node5Address, Status, err := QueryingNodeAddress(Node5Url)
	Node6Address, Status, err := QueryingNodeAddress(Node6Url)

	fmt.Println("Create Channels:A(100)-B(50) A(100)-C(50) B(100)-D(50) C(100)-D(50) D(100)-E(50) E(100)-F(50)")
	//本地注释：A-B建立通道
	Channel, Status, err := OpenChannel(Node1Url, Node2Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node2Url, Channel.ChannelAddress, 50)
	//本地注释：A-C建立通道
	Channel, Status, err = OpenChannel(Node1Url, Node3Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node3Url, Channel.ChannelAddress, 50)
	//本地注释：B-D建立通道
	Channel, Status, err = OpenChannel(Node2Url, Node4Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node4Url, Channel.ChannelAddress, 50)
	//本地注释：C-D建立通道
	Channel, Status, err = OpenChannel(Node3Url, Node4Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node4Url, Channel.ChannelAddress, 50)
	//本地注释：D-E建立通道
	Channel, Status, err = OpenChannel(Node4Url, Node5Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node5Url, Channel.ChannelAddress, 50)
	//本地注释：E-F建立通道
	Channel, Status, err = OpenChannel(Node5Url, Node6Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node6Url, Channel.ChannelAddress, 50)
	////本地注释：D-A建立通道
	//Channel, Status, err = OpenChannel(Node4Url, Node1Address.OurAddress, NewTokenNames[0], 100, 1000)
	//Deposit2Channel(Node1Url, Channel.ChannelAddress, 50)

	fmt.Println("Create Channels Complete")
	var Amount int32

	Amount = 5
	fmt.Println("Transfer ", Amount, " tokens from A to B")
	//本地注释：A->F 5Token
	TransferResult, Status, err := InitiatingTransfer(Node1Url, NewTokenName, Node2Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node2Address.OurAddress, NewTokenName, Amount)

	Amount = 6
	fmt.Println("Transfer ", Amount, " tokens from A to C")
	//本地注释：A->F 5Token
	TransferResult, Status, err = InitiatingTransfer(Node1Url, NewTokenName, Node3Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node3Address.OurAddress, NewTokenName, Amount)

	Amount = 7
	fmt.Println("Transfer ", Amount, " tokens from A to D")
	//本地注释：A->F 5Token
	TransferResult, Status, err = InitiatingTransfer(Node1Url, NewTokenName, Node4Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node4Address.OurAddress, NewTokenName, Amount)

	Amount = 8
	fmt.Println("Transfer ", Amount, " tokens from A to F")
	//本地注释：A->F 5Token
	TransferResult, Status, err = InitiatingTransfer(Node1Url, NewTokenName, Node6Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node6Address.OurAddress, NewTokenName, Amount)

	fmt.Println("Finish TransferTest")

	//fmt.Println("==============================================================================================")

	return 0
}
