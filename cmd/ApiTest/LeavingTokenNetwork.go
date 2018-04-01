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
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	fmt.Println("Start Leaving Token Network")
	Tokens, _, _ := QueryingRegisteredTokens(url)

	//本地注释：测试不存在的Token
	Status, err = LeavingTokenNetwork(url, "0x00000", true)
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowLeavingTokenNetworkMsgDetail(Status)
	if Status == "500 Internal Server Error" {
		fmt.Println("Test pass:Leaving a not exist TokenNetwork")
	} else {
		fmt.Println("Test failed:Leaving a not exist TokenNetwork")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	//本地注释：测试已经注册的Token
	for i := 0; i < len(Tokens); i++ {
		Status, err = LeavingTokenNetwork(url, Tokens[i], false)
		ShowError(err)
		//本地注释：显示错误详细信息
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
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowLeavingTokenNetworkMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successfully leaving a token network")

	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,leaving a TokenNetwork Failure:%s\n", Status)
	}
}
