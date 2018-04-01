package main

import (
	"fmt"
	"github.com/larspensjo/config"
	"log"
	"net/http"
	"time"
)

//本地注释：注册新Token
func RegisteringOneToken(url string, Token string) (Status string, err error) {
	var resp *http.Response
	var count int

	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, url+"/api/1/tokens/"+Token, nil)
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//req.Header.Set("Cookie", "name=anny")
		resp, err = client.Do(req)
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

//本地注释：测试注册新Token
func RegisteringOneTokenTest(url string) {
	start := time.Now()
	ShowTime()
	fmt.Println("Start Registering Token Test")
	fmt.Println("Start Registering a already registered token")
	Tokens, _, _ := QueryingRegisteredTokens(url)
	Status, err := RegisteringOneToken(url, Tokens[0])
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowRegisteringOneTokenMsgDetail(Status)
	switch Status {
	case "409 Conflict":
		fmt.Println("Test pass:Registering a already registered token")

	default:
		fmt.Println("Test failed:Registering  a already registered Token:", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	//registering a nonexistent token
	fmt.Println("Start Registering a nonexistent token")
	Status, err = RegisteringOneToken(url, "0xffffffffffffffffffffffffffffffffffffffff")
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowRegisteringOneTokenMsgDetail(Status)
	switch Status {
	case "201 Created":
		fmt.Println("Test pass:Success Registering a new token")

	default:
		fmt.Println("Test failed:Registering new Token:", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	fmt.Println("Start Registering a new token")
	c, err := config.ReadDefault("./ApiTest.INI")

	if err != nil {
		fmt.Println("config.ReadDefault error:", err)
		return
	}

	EthRpcEndpoint, err := c.String("common", "eth_rpc_endpoint")
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	KeyStorePath, err := c.String("common", "keystore_path")
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	NewTokenName, _, _ := CreateNewToken(EthRpcEndpoint, KeyStorePath)

	Status, err = RegisteringOneToken(url, NewTokenName)
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowRegisteringOneTokenMsgDetail(Status)
	switch Status {
	case "201 Created":
		fmt.Println("Test pass:Success Registering a new token")

	default:
		fmt.Println("Test failed:Registering new Token:", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}

	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowRegisteringOneTokenMsgDetail(Status string) {
	switch Status {
	case "201 Created":
		fmt.Println("Success registering token")
	case "409 Conflict":
		fmt.Println("Token already registered")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Println("Unknown error,RegisteringOneToken!:", Status)
	}
}
