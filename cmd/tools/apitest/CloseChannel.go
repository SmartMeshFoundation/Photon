package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

//test for closing the specified channel for the node
func CloseChannel(url string, Channel string) (Status string, err error) {
	var payload string
	var count int
	payload = "{\"state\":\"closed\"}"
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/channels/" + Channel
		req, _ := http.NewRequest("PATCH", fullurl, strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "name=anny")
		Status, _, err = DoRequest(client, req)
		//body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if count >= MaxTry {
		Status = "504 TimeOut"
	}
	return
}

//test for closing the specified channel for the node
//there are four kinds of situations. first, the channel  does not exist; second, the channel is opened;third, the channel has closed;fourth, the channel is settled.
func CloseChannelTest(url string) {
	var err error
	var ChannelAddress string
	var i int
	var Status string
	start := time.Now()
	ShowTime()
	log.Println("Start Close Channel")

	//close the channel that does not exist
	ChannelAddress = "0x00000"
	Status, err = CloseChannel(url, ChannelAddress)
	ShowError(err)
	ShowCloseChannelMsgDetail(Status)
	if Status == "409 Conflict" {
		log.Println("Test pass:Close a not exist Channel")
	} else {
		log.Println("Test failed:Close a not exist Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}

	//close the channel which is opened
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
	ShowCloseChannelMsgDetail(Status)
	if Status == "200 OK" {
		log.Println("Test pass:Close a opened Channel")
	} else {
		log.Println("Test failed:Close a opened Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	//close a channel which has closed,close a settled channel,if no reaction ,ignore it!
	goto EndTest
Testclosed:
	//close a channel which has closed
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
	ShowCloseChannelMsgDetail(Status)
	if Status == "200 OK" {
		log.Println("Test pass:Close a closed Channel")
	} else {
		log.Println("Test failed:Close a closed Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
Testsettled:
	// close a settled channel
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
	ShowCloseChannelMsgDetail(Status)
	if Status == "200 OK" {
		log.Println("Test pass:Close a settled Channel")
	} else {
		log.Println("Test failed:Close a settled Channel")
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
func ShowCloseChannelMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Close Channel Success!")
	case "400 Bad Request":
		log.Println("The provided json is in some way malformed!")
	case "409 Conflict":
		log.Println("Provided channel does not exist")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,Close Channel Failure! %s\n", Status)
	}
}
