package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

//settle the specified channel for the node
func SettleChannel(url string, Channel string) (Status string, err error) {
	var resp *http.Response
	var payload string
	var count int
	payload = "{\"state\":\"settled\"}"
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/channels/" + Channel
		req, _ := http.NewRequest("PATCH", fullurl, strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "name=anny")
		resp, err = client.Do(req)
		//body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if resp != nil {
		Status = resp.Status
	} else {
		Status = "null"
	}
	if count >= MaxTry {
		Status = "504 TimeOut"
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//test for settling the specified channel for the node
//there are four kinds of situations. first, the channel  does not exist; second, the channel is opened;third, the channel has closed;fourth, the channel is settled.
func SettleChannelTest(url string) {
	var err error
	var ChannelAddress string
	var i int
	var Status string
	start := time.Now()
	ShowTime()
	log.Println("Start Settle Channel")

	//Settle the channel that does not exist
	ChannelAddress = "0x00000"
	Status, err = CloseChannel(url, ChannelAddress)
	ShowError(err)
	//display the details of the error
	ShowSettleChannelMsgDetail(Status)
	if Status == "409 Conflict" {
		log.Println("Test pass:Settle a not exist Channel")
	} else {
		log.Println("Test failed:Settle a not exist Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}

	//settle the channel which is opened
	Channels, _, _ := QueryingNodeAllChannels(url)
	//query the first opened channel
	for i = 0; i < len(Channels); i++ {
		if Channels[i].State == "opened" {
			ChannelAddress = Channels[i].ChannelAddress
			break
		}
	}
	if i > len(Channels) {
		goto Testclosed
	}
	Status, err = CloseChannel(url, ChannelAddress)
	ShowError(err)
	ShowSettleChannelMsgDetail(Status)
	if Status == "200 OK" {
		log.Println("Test pass:Settle a opened Channel")
	} else {
		log.Println("Test failed:Settle a opened Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
Testclosed:
	//settle a channel which has closed
	//query the first closed channel
	for i = 0; i < len(Channels); i++ {
		if Channels[i].State == "closed" {
			ChannelAddress = Channels[i].ChannelAddress
			break
		}
	}
	if i > len(Channels) {
		goto Testsettled
	}
	Status, err = CloseChannel(url, ChannelAddress)
	ShowError(err)
	//display the details of the error
	ShowCloseChannelMsgDetail(Status)
	if Status == "200 OK" {
		log.Println("Test pass:Settle a closed Channel")
	} else {
		log.Println("Test failed:Settle a closed Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
Testsettled:
	//settle a settled channel
	//query the first settled channel
	for i = 0; i < len(Channels); i++ {
		if Channels[i].State == "settled" {
			ChannelAddress = Channels[i].ChannelAddress
			break
		}
	}
	if i > len(Channels) {
		goto EndTest
	}
	Status, err = CloseChannel(url, ChannelAddress)
	ShowError(err)
	ShowSettleChannelMsgDetail(Status)
	if Status == "200 OK" {
		log.Println("Test pass:Settle a settled Channel")
	} else {
		log.Println("Test failed:Settle a settled Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
EndTest:
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowSettleChannelMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Settle Channel Success!")
	case "400 Bad Request":
		log.Println("The provided json is in some way malformed!")
	case "409 Conflict":
		log.Println("Provided channel does not exist，or is inside settlement period")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,Settle Channel Failure! %s\n", Status)
	}
}
