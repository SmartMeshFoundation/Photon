package main

import (
	"log"

	"github.com/huamou/config"
)

//detail for testing transfer
func TransferTest(NewTokenName string) (code int) {

	log.Println("==============================================================================================")
	log.Println("Start TransferTest")

	c, err := config.ReadDefault("./apitest.ini")

	if err != nil {
		log.Println("config.ReadDefault error:", err)
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

	//deposit to new token account

	//get the node address
	Node1Address, Status, err := QueryingNodeAddress(Node1Url)
	Node2Address, Status, err := QueryingNodeAddress(Node2Url)
	Node3Address, Status, err := QueryingNodeAddress(Node3Url)
	Node4Address, Status, err := QueryingNodeAddress(Node4Url)
	Node5Address, Status, err := QueryingNodeAddress(Node5Url)
	Node6Address, Status, err := QueryingNodeAddress(Node6Url)

	log.Println("Create Channels:A(100)-B(50) A(100)-C(50) B(100)-D(50) C(100)-D(50) D(100)-E(50) E(100)-F(50)")
	//A-B establish channel
	Channel, Status, err := OpenChannel(Node1Url, Node2Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node2Url, Channel.ChannelAddress, 50)
	//A-C establish channel
	Channel, Status, err = OpenChannel(Node1Url, Node3Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node3Url, Channel.ChannelAddress, 50)
	//B-D establish channel
	Channel, Status, err = OpenChannel(Node2Url, Node4Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node4Url, Channel.ChannelAddress, 50)
	//C-D establish channel
	Channel, Status, err = OpenChannel(Node3Url, Node4Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node4Url, Channel.ChannelAddress, 50)
	//D-E establish channel
	Channel, Status, err = OpenChannel(Node4Url, Node5Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node5Url, Channel.ChannelAddress, 50)
	//E-F establish channel
	Channel, Status, err = OpenChannel(Node5Url, Node6Address.OurAddress, NewTokenName, 100, 1000)
	Deposit2Channel(Node6Url, Channel.ChannelAddress, 50)
	//D-A establish channel
	//Channel, Status, err = OpenChannel(Node4Url, Node1Address.OurAddress, NewTokenNames[0], 100, 1000)
	//Deposit2Channel(Node1Url, Channel.ChannelAddress, 50)

	log.Println("Create Channels Complete")
	var Amount int32

	Amount = 5
	log.Println("Transfer ", Amount, " tokens from A to B")
	//A->F 5Token
	TransferResult, Status, err := InitiatingTransfer(Node1Url, NewTokenName, Node2Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node2Address.OurAddress, NewTokenName, Amount)

	Amount = 6
	log.Println("Transfer ", Amount, " tokens from A to C")
	//A->F 5Token
	TransferResult, Status, err = InitiatingTransfer(Node1Url, NewTokenName, Node3Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node3Address.OurAddress, NewTokenName, Amount)

	Amount = 7
	log.Println("Transfer ", Amount, " tokens from A to D")
	//A->F 5Token
	TransferResult, Status, err = InitiatingTransfer(Node1Url, NewTokenName, Node4Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node4Address.OurAddress, NewTokenName, Amount)

	Amount = 8
	log.Println("Transfer ", Amount, " tokens from A to F")
	//A->F 5Token
	TransferResult, Status, err = InitiatingTransfer(Node1Url, NewTokenName, Node6Address.OurAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	ResultJudge(TransferResult, Status, err, Node1Address.OurAddress, Node6Address.OurAddress, NewTokenName, Amount)

	log.Println("Finish TransferTest")

	//log.Println("==============================================================================================")

	return 0
}
