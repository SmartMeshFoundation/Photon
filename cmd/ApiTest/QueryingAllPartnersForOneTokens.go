package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//本地注释：查询节点指定Token有通道的伙伴地址
func QueryingAllPartnersForOneToken(url string, token string) (Partners []TokenPartnerPayload, Status string, err error) {
	var resp *http.Response
	var count int
	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1/tokens/" + token + "/partners")
		if err == nil {
			Status = resp.Status
			break
		}
		time.Sleep(time.Second)
	}
	if count >= MaxTry {
		Status = "504 Timeout"
	}
	if resp == nil {
		Partners = nil
	} else {
		if resp.Status == "200 OK" {
			p, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(p, &Partners)
		} else {
			Partners = nil
		}
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//本地注释：测试查询节点指定Token有通道的伙伴地址
func QueryingAllPartnersForOneTokensTest(url string) {
	start := time.Now()
	ShowTime()
	fmt.Println("Start Querying All Partners For a Existent Token")
	Tokens, _, _ := QueryingRegisteredTokens(url)
	//fmt.Printf("!!!!!!!!!!!!Token:%s\n", Tokens[0])
	if Tokens == nil {
		fmt.Println("Warning:No registered Token!")
	} else {
		_, Status, err := QueryingAllPartnersForOneToken(url, Tokens[0])
		ShowError(err)
		//本地注释：显示错误详细信息
		ShowQueryingAllPartnersForOneTokenMsgDetail(Status)
		switch Status {
		case "200 OK":
			fmt.Println("Test pass:Querying All Partners For a Tokens Success!")
		default:
			fmt.Printf("Test failed:Querying All Partners For a Tokens Failure! %s\n", Status)
			if HalfLife {
				log.Fatal("HalfLife,exit")
			}
		}
	}

	fmt.Println("Start Querying All Partners For a Nonexistent Token")
	_, Status, err := QueryingAllPartnersForOneToken(url, "0x0000")
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowQueryingAllPartnersForOneTokenMsgDetail(Status)
	switch Status {
	case "404 Not Found":
		fmt.Println("Test pass:Querying All Partners For a Nonexistent Tokens Success!")
	default:
		fmt.Printf("Test failed:Querying All Partners For a Nonexistent Tokens Failure! %s\n", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowQueryingAllPartnersForOneTokenMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successful query")
	case "302 Redirect":
		fmt.Println("The user accesses the channel link endpoint")
	case "404 Not Found":
		fmt.Println("The token does not exist")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingAllPartnersForOneToken Failure:%s\n", Status)
	}
}
