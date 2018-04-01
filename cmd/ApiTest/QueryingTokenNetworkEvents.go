package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func QueryingTokenNetworkEvents(url string, Token string) (Events []TokenNetworkEvents, Status string, err error) {
	var resp *http.Response
	var count int

	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1//events/tokens/" + Token)
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
	if resp == nil {
		Events = nil
	} else {
		if resp.Status == "200 OK" {
			p, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(p, &Events)
		} else {
			Events = nil
		}
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

func QueryingTokenNetworkEventsTest(url string) {
	//var Events []TokenNetworkEvents
	var Status string
	var err error
	start := time.Now()
	ShowTime()
	fmt.Println("Start Querying Token Network Events")
	Tokens, _, _ := QueryingRegisteredTokens(url)
	//本地注释：测试不存在的Token
	_, Status, err = QueryingTokenNetworkEvents(url, "0xffffffffffffffffffffffffffffffffffffffff")
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowQueryingTokenNetworkEventsMsgDetail(Status)
	switch Status {
	case "404 Not Found":
		fmt.Println("Test pass:Querying  nonexistent Tokens")
	default:
		fmt.Println("Test failed:Querying  nonexistent Tokens:", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}

	for i := 0; i < len(Tokens); i++ {
		_, Status, err = QueryingTokenNetworkEvents(url, Tokens[i])
		ShowError(err)
		//本地注释：显示错误详细信息
		ShowQueryingTokenNetworkEventsMsgDetail(Status)
		switch Status {
		case "200 OK":
			fmt.Println("Test pass:QueryingRegisteredTokens:", Tokens[i])
		default:
			fmt.Printf("Test failed:QueryingRegisteredTokens:", Status)
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
func ShowQueryingTokenNetworkEventsMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successful query")
	case "400 Bad Request":
		fmt.Println("The provided query string is malformed")
	case "404 Not Found":
		fmt.Println("The token does not exist")
	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingGeneralNetworkEvents Failure:%s\n", Status)
	}
}
