package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//deposit to the specified channel
func Deposit2Channel(url string, Channel string, Balance int32) (Status string, err error) {
	var resp *http.Response
	var count int
	var payload Desposit2ChannelPayload
	payload.Balance = Balance
	p, _ := json.Marshal(payload)
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/channels/" + Channel
		req, _ := http.NewRequest(http.MethodPatch, fullurl, bytes.NewReader(p))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "name=anny")
		resp, err = client.Do(req)
		//body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
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

//test for depositing  to the specified channel
func Deposit2ChannelTest(url string) {
	var err error
	var ChannelAddress string
	var Status string
	var Balance int32
	var i int
	start := time.Now()
	ShowTime()
	log.Println("Start Deposit2Channel")
	Balance = 100
	log.Println("Deposit to a not exist Channel")
	//deposit to the channel that doesn't exist
	ChannelAddress = "0xffffffffffffffffffffffffffffffffffffffff"
	Status, err = Deposit2Channel(url, ChannelAddress, Balance)
	ShowError(err)
	//display the details of the error
	ShowDeposit2ChannelMsgDetail(Status)
	if Status == "409 Conflict" {
		log.Println("Test pass:Deposit to a not exist Channel")
	} else {
		log.Println("Test failed:Deposit to a not exist Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}

	log.Println("Deposit to a opened Channel")
	//query all channels of the node
	Channels, _, _ := QueryingNodeAllChannels(url)
	//deposit to the opened channel
	for i = 0; i < len(Channels); i++ {
		if Channels[i].State == "opened" {
			ChannelAddress = Channels[i].ChannelAddress
			Status, err = Deposit2Channel(url, ChannelAddress, Balance)
			ShowError(err)
			//display the details of the error
			ShowDeposit2ChannelMsgDetail(Status)
			if Status == "200 OK" {
				fmt.Printf("Test pass:Deposit to a open Channel:%s\n", ChannelAddress)
			} else {
				fmt.Printf("Test failed:Deposit to a open Channel:%s\n", ChannelAddress)
				if HalfLife {
					log.Fatal("HalfLife,exit")
				}
			}
		}
	}
	log.Println("Deposit to a closed Channel")
	// deposit the closed channel
	for i = 0; i < len(Channels); i++ {
		if Channels[i].State == "closed" {
			ChannelAddress = Channels[i].ChannelAddress
			Status, err = Deposit2Channel(url, ChannelAddress, Balance)
			ShowError(err)
			//display the details of the error
			ShowDeposit2ChannelMsgDetail(Status)
			if Status == "408 Request Timeout" {
				fmt.Printf("Test pass:Deposit to a closed Channel:%s\n", ChannelAddress)
			} else {
				fmt.Printf("Test failed:Deposit to a closed Channel:%s\n", ChannelAddress)
				if HalfLife {
					log.Fatal("HalfLife,exit")
				}
			}
		}
	}
	log.Println("Deposit to a settled Channel")
	// deposit the settled channel
	for i = 0; i < len(Channels); i++ {
		if Channels[i].State == "settled" {
			ChannelAddress = Channels[i].ChannelAddress
			Status, err = Deposit2Channel(url, ChannelAddress, Balance)
			ShowError(err)
			//display the details of the error
			ShowDeposit2ChannelMsgDetail(Status)
			if Status == "408 Request Timeout" {
				fmt.Printf("Test pass:Deposit to a settled Channel:%s\n", ChannelAddress)
			} else {
				fmt.Printf("Test failed:Deposit to a settled Channel:%s\n", ChannelAddress)
				if HalfLife {
					log.Fatal("HalfLife,exit")
				}
			}
		}
	}

	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowDeposit2ChannelMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Success Deposit")
	case "400 Bad Request":
		log.Println("The provided json is in some way malformed!")
	case "402 Payment required":
		log.Println("Insufficient balance to do a deposit")
	case "408 Request Timeout":
		log.Println("The deposit event was not read in time by the ethereum node")
	case "409 Conflict":
		log.Println("Provided channel does not exist")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,Deposit Failure! %s\n", Status)
	}
}
