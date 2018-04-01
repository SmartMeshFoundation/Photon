package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//本地注释： 连接TokenNetwork
func Connecting2TokenNetwork(url string, Token string, Funds int32) (Status string, err error) {
	var resp *http.Response
	var count int
	var payload Connecting2TokenNetworkPayload
	payload.Funds = Funds
	p, _ := json.Marshal(payload)
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/connections/" + Token
		req, _ := http.NewRequest(http.MethodPut, fullurl, bytes.NewReader(p))
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

//本地注释： 测试连接TokenNetwork
func Connecting2TokenNetworkTest(url string, Funds int32) {
	var err error
	var Status string
	start := time.Now()
	ShowTime()
	log.Println("Start Connecting2TokenNetwork")
	Tokens, _, _ := QueryingRegisteredTokens(url)

	//本地注释：测试不存在的Token
	log.Println("Start Connecting to a not exist TokenNetwork")
	Status, err = Connecting2TokenNetwork(url, "0xffffffffffffffffffffffffffffffffffffffff", Funds)
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowConnecting2TokenNetworkMsgDetail(Status)
	if Status == "500 Internal Server Error" {
		log.Println("Test pass:Connecting to a not exist TokenNetwork")
	} else {
		log.Println("Test failed:Connecting to a not exist TokenNetwork")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	log.Println("Start Connecting to a registered TokenNetwork")
	//本地注释：测试已经注册的Token
	for i := 0; i < len(Tokens); i++ {
		Status, err = Connecting2TokenNetwork(url, Tokens[i], Funds)
		ShowError(err)
		//本地注释：显示错误详细信息
		ShowConnecting2TokenNetworkMsgDetail(Status)
		if Status == "204 No Content" {
			fmt.Printf("Test pass:Connecting2TokenNetwork [%s]\n", Tokens[i])
		} else {
			fmt.Printf("Test failed:Connecting2TokenNetwork [%s]\n", Tokens[i])
			if HalfLife {
				log.Fatal("HalfLife,exit")
			}
		}
	}

	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowConnecting2TokenNetworkMsgDetail(Status string) {
	switch Status {
	case "204 No Content":
		log.Println("Successful connection creation")
	case "402 Payment required":
		log.Println("Any of the channel deposits fail due to insufficient ETH balance")
	case "408 Request Timeout":
		log.Println("A timeout happened during any of the transactions")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,Connecting2TokenNetwork Failure:%s\n", Status)
	}
}
