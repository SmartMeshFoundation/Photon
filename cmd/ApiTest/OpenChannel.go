package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//establish the Channel
func OpenChannel(url string, PartnerAddress string, TokenAddress string, Balance int32, SettleTimeout int32) (Channel NodeChannel, Status string, err error) {
	var count int
	var resp *http.Response
	var newchannel OpenChannelPayload
	newchannel.PartnerAddress = PartnerAddress
	newchannel.TokenAddress = TokenAddress
	newchannel.Balance = Balance
	newchannel.SettleTimeout = SettleTimeout
	p, _ := json.Marshal(newchannel)
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, url+"/api/1/channels", bytes.NewReader(p))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", "name=anny")
		resp, err = client.Do(req)
		//body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			//io.Copy(os.Stdout, resp.Body)
			if resp != nil {
				p, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(p, &Channel)
			}
			break
		}
		time.Sleep(time.Second)
	}

	if count >= MaxTry {
		Status = "504 TimeOut"
	} else {
		Status = resp.Status
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//establish the channel between the node1 and node2
func OpenChannelTest(url string, url2 string) {
	start := time.Now()
	ShowTime()
	log.Println("Start Open Channel")
	Address, _, _ := QueryingNodeAddress(url2)
	PartnerAddress := Address.OurAddress
	Tokens, _, _ := QueryingRegisteredTokens(url)
	TokenAddress := Tokens[0]
	Balance := int32(100)
	SettleTimeout := int32(800)
	_, Status, err := OpenChannel(url, PartnerAddress, TokenAddress, Balance, SettleTimeout)
	ShowError(err)
	ShowOpenChannelMsgDetail(Status)
	switch Status {
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
