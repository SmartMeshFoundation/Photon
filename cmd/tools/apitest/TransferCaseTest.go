package main

import (
	"log"
	"strconv"
	"time"

	"github.com/huamou/config"
)

func Transfer(NewTokenName string, IniFile string) {
	log.Println("==============================================================================================")
	log.Println("Start TransferTest")
	SChannels, TransCase, DChannels, Result := TransferParmReader(IniFile)
	TransferCase(NewTokenName, SChannels, TransCase, DChannels, Result)
	log.Println("Finish TransferTest")

}

//test transfer
func TransferCase(NewTokenName string, SChannels []TransferCaseChannel, TransCase TransferCaseTransfer, DChannels []TransferCaseChannel, Result bool) (code int) {
	c, err := config.ReadDefault("./../../testdata/TransCase/apitest.ini")
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

	//establish the channel
	for i := 0; i < len(SChannels); i++ {

		log.Println("Create Channels:", SChannels[i].Node1Url, "(", SChannels[i].Balance1, ")-", SChannels[i].Node2Url, "(", SChannels[i].Balance2, ")")
		switch SChannels[i].Node1Url {
		case "A":
			SChannels[i].Node1Url = Node1Url
		case "B":
			SChannels[i].Node1Url = Node2Url
		case "C":
			SChannels[i].Node1Url = Node3Url
		case "D":
			SChannels[i].Node1Url = Node4Url
		case "E":
			SChannels[i].Node1Url = Node5Url
		case "F":
			SChannels[i].Node1Url = Node6Url
		default:
			log.Fatal("Node number error", SChannels[i].Node1Url)
		}
		switch SChannels[i].Node2Url {
		case "A":
			SChannels[i].Node2Url = Node1Url
		case "B":
			SChannels[i].Node2Url = Node2Url
		case "C":
			SChannels[i].Node2Url = Node3Url
		case "D":
			SChannels[i].Node2Url = Node4Url
		case "E":
			SChannels[i].Node2Url = Node5Url
		case "F":
			SChannels[i].Node2Url = Node6Url
		default:
			log.Fatal("Node number error", SChannels[i].Node2Url)
		}

		NodeBAddress, _, _ := QueryingNodeAddress(SChannels[i].Node2Url)
		Channel, _, _ := OpenChannel(SChannels[i].Node1Url, NodeBAddress.OurAddress, NewTokenName, SChannels[i].Balance1, 200)
		Deposit2Channel(SChannels[i].Node2Url, Channel.ChannelAddress, SChannels[i].Balance2)

	}

	log.Println("Create Channels Complete")

	var Amount int32

	Amount = TransCase.Balance
	log.Println("Transfer ", Amount, " tokens from ", TransCase.Node1Url, " to ", TransCase.Node2Url)

	switch TransCase.Node1Url {
	case "A":
		TransCase.Node1Url = Node1Url
	case "B":
		TransCase.Node1Url = Node2Url
	case "C":
		TransCase.Node1Url = Node3Url
	case "D":
		TransCase.Node1Url = Node4Url
	case "E":
		TransCase.Node1Url = Node5Url
	case "F":
		TransCase.Node1Url = Node6Url
	default:
		log.Fatal("Node number error", TransCase.Node1Url)
	}

	switch TransCase.Node2Url {
	case "A":
		TransCase.Node2Url = Node1Url
	case "B":
		TransCase.Node2Url = Node2Url
	case "C":
		TransCase.Node2Url = Node3Url
	case "D":
		TransCase.Node2Url = Node4Url
	case "E":
		TransCase.Node2Url = Node5Url
	case "F":
		TransCase.Node2Url = Node6Url
	default:
		log.Fatal("Node number error", TransCase.Node2Url)
	}

	NodeAAddress, _, _ := QueryingNodeAddress(TransCase.Node1Url)
	NodeBAddress, _, _ := QueryingNodeAddress(TransCase.Node2Url)
	//log.Println("TransCase.Node1Url=", TransCase.Node1Url, " NewTokenName=", NewTokenName, " NodeBAddress.OurAddress=", NodeBAddress.OurAddress, " Amount=", Amount)
	TransferResult, Status, err := InitiatingTransfer(TransCase.Node1Url, NewTokenName, NodeBAddress.OurAddress, Amount)
	//ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
	//
	ResultJudge(TransferResult, Status, err, NodeAAddress.OurAddress, NodeBAddress.OurAddress, NewTokenName, Amount)
	//
	time.Sleep(6 * time.Second)
	if CheckChannel(DChannels) {
		log.Println("Transfer case test success!")
	} else {
		log.Println("Transfer case test failure!")
	}
	//time.Sleep(2 * time.Hour)
	return 0
}

//read transfer parameters
func TransferParmReader(IniFile string) (SChannels []TransferCaseChannel, TransCase TransferCaseTransfer, DChannels []TransferCaseChannel, Result bool) {
	c, err := config.ReadDefault(IniFile)
	if err != nil {
		log.Println(IniFile, " Read error:", err)
		return
	}
	ChannelCount, _ := c.Int("MAIN", "ChannelCount")
	TransCase.Balance, _ = c.Int32("MAIN", "Amount")
	TransCase.Node1Url, _ = c.String("MAIN", "NOTE1")
	TransCase.Node2Url, _ = c.String("MAIN", "NOTE2")
	Result, _ = c.Bool("MAIN", "Result")

	var seg string
	var aTrans TransferCaseChannel
	for i := 0; i < ChannelCount; i++ {
		seg = "SChannle" + strconv.Itoa(i+1)
		SChannels = append(SChannels, aTrans)
		SChannels[i].Node1Url = c.RdString(seg, "NOTE1", "A")
		SChannels[i].Node2Url = c.RdString(seg, "NOTE2", "B")
		SChannels[i].Balance1 = c.RdInt32(seg, "Balance1", 100)
		SChannels[i].Balance2 = c.RdInt32(seg, "Balance2", 100)
		seg = "DChannle" + strconv.Itoa(i+1)
		DChannels = append(DChannels, aTrans)
		DChannels[i].Node1Url = c.RdString(seg, "NOTE1", "A")
		DChannels[i].Node2Url = c.RdString(seg, "NOTE2", "B")
		DChannels[i].Balance1 = c.RdInt32(seg, "Balance1", 100)
		DChannels[i].Balance2 = c.RdInt32(seg, "Balance2", 100)
		DChannels[i].LockedBalance1 = c.RdInt32(seg, "LockedBalance1", 0)
		DChannels[i].LockedBalance2 = c.RdInt32(seg, "LockedBalance2", 0)
	}

	return
}

func CheckChannel(DChannels []TransferCaseChannel) (checked bool) {

	checked = false
	c, err := config.ReadDefault("./../../testdata/TransCase/apitest.ini")
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

	//query the channel
	for i := 0; i < len(DChannels); i++ {

		switch DChannels[i].Node1Url {
		case "A":
			DChannels[i].Node1Url = Node1Url
		case "B":
			DChannels[i].Node1Url = Node2Url
		case "C":
			DChannels[i].Node1Url = Node3Url
		case "D":
			DChannels[i].Node1Url = Node4Url
		case "E":
			DChannels[i].Node1Url = Node5Url
		case "F":
			DChannels[i].Node1Url = Node6Url
		default:
			log.Fatal("Node number error", DChannels[i].Node1Url)
		}
		switch DChannels[i].Node2Url {
		case "A":
			DChannels[i].Node2Url = Node1Url
		case "B":
			DChannels[i].Node2Url = Node2Url
		case "C":
			DChannels[i].Node2Url = Node3Url
		case "D":
			DChannels[i].Node2Url = Node4Url
		case "E":
			DChannels[i].Node2Url = Node5Url
		case "F":
			DChannels[i].Node2Url = Node6Url
		default:
			log.Fatal("Node number error", DChannels[i].Node2Url)
		}

		//NodeAAddress, _, _ := QueryingNodeAddress(DChannels[i].Node1Url)
		NodeBAddress, _, _ := QueryingNodeAddress(DChannels[i].Node2Url)
		Channels, _, _ := QueryingNodeAllChannels(DChannels[i].Node1Url)
		var l int
		for l = 0; l < len(Channels); l++ {
			if Channels[l].PartnerAddress == NodeBAddress.OurAddress {
				if Channels[l].Balance != DChannels[i].Balance1 {
					return false
				}
				if Channels[l].LockedAmount != DChannels[i].LockedBalance1 {
					return false
				}
				if Channels[l].PatnerBalance != DChannels[i].Balance2 {
					return false
				}
				if Channels[l].PartnerLockedAmount != DChannels[i].LockedBalance2 {
					return false
				}
				break
			}
		}
		if l >= len(Channels) {
			return false
		}
	}

	return true
}
