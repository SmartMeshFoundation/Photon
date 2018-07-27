package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func LeavingTokenNetwork(url string, Token string, OnlyReceivingChannels bool) (Status string, err error) {
	var resp *http.Response
	var count int
	var payload LeavingTokenNetworkPayload
	payload.OnlyReceivingChannels = OnlyReceivingChannels
	p, _ := json.Marshal(payload)
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/connections/" + Token
		req, _ := http.NewRequest(http.MethodDelete, fullurl, bytes.NewReader(p))
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

func LeavingTokenNetworkTest(url string) {
	var err error
	var Status string
	start := time.Now()
	ShowTime()
	log.Println("Start Leaving TokenNetworkAddres Network")
	Tokens, _, _ := QueryingRegisteredTokens(url)

	//test the token which doesn't exist.
	Status, err = LeavingTokenNetwork(url, "0x00000", true)
	ShowError(err)
	//display the details of the error
	ShowLeavingTokenNetworkMsgDetail(Status)
	if Status == "500 Internal Server Error" {
		log.Println("Test pass:Leaving a not exist TokenNetwork")
	} else {
		log.Println("Test failed:Leaving a not exist TokenNetwork")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	//test the token which has registered.
	for i := 0; i < len(Tokens); i++ {
		Status, err = LeavingTokenNetwork(url, Tokens[i], false)
		ShowError(err)
		//display the details of the error
		ShowLeavingTokenNetworkMsgDetail(Status)
		if Status == "200 OK" {
			fmt.Printf("Test pass:Leaving TokenNetwork [%s]\n", Tokens[i])
		} else {
			fmt.Printf("Test failed:Leaving TokenNetwork [%s]\n", Tokens[i])
			if HalfLife {
				log.Fatal("HalfLife,exit")
			}
		}
	}
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowLeavingTokenNetworkMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Successfully leaving a token network")

	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,leaving a TokenNetwork Failure:%s\n", Status)
	}
}
