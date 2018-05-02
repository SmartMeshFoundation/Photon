package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

//establish the Channel
func OpenChannel(url string, PartnerAddress string, TokenAddress string, Balance int32, SettleTimeout int32) (Channel NodeChannel, Status string, err error) {
	var count int
	var newchannel OpenChannelPayload
	newchannel.PartnerAddress = PartnerAddress
	newchannel.TokenAddress = TokenAddress
	newchannel.Balance = Balance
	newchannel.SettleTimeout = SettleTimeout
	p, _ := json.Marshal(newchannel)
	for count = 0; count < MaxTry; count = count + 1 {
		var body []byte
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, url+"/api/1/channels", bytes.NewReader(p))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "name=anny")
		Status, body, err = DoRequest(client, req)
		if err == nil {
			err = json.Unmarshal(body, &Channel)
			break
		}
		time.Sleep(time.Second)
	}

	if count >= MaxTry {
		Status = "504 TimeOut"
	}
	return
}

//establish the channel between the node1 and node2
func OpenChannelTest(url string) {
	start := time.Now()
	ShowTime()
	log.Println("Start Open Channel")
	PartnerAddress := utils.NewRandomAddress().String()
	Tokens, _, _ := QueryingRegisteredTokens(url)
	log.Printf("all tokens =%#v\n", Tokens)
	TokenAddress := Tokens[0]
	Balance := int32(100)
	SettleTimeout := int32(800)
	_, Status, err := OpenChannel(url, PartnerAddress, TokenAddress, Balance, SettleTimeout)
	ShowError(err)
	ShowOpenChannelMsgDetail(Status)
	switch Status {
	case "200 OK":
		fallthrough
	case "201 Created":
		log.Println("Test pass: successful Creation channels")
	default:
		fmt.Printf("Test failed: %s\n", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowOpenChannelMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fallthrough
	case "201 Created":
		log.Println("Successful Creation channels")
	case "400 Bad Request":
		log.Println("The provided json is in some way malformed")
	case "408 Request Timeout":
		log.Println("The deposit event was not read in time by the ethereum node")
	case "409 Conflict":
		log.Println("The input is invalid, such as too low a settle timeout")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		log.Println("Unknown error,OpenChannel:", Status)
	}
}
