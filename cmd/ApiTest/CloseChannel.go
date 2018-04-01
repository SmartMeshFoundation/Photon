package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

//本地注释：测试关闭节点指定通道
func CloseChannel(url string, Channel string) (Status string, err error) {
	var resp *http.Response
	var payload string
	var count int
	payload = "{\"state\":\"closed\"}"
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/channels/" + Channel
		req, _ := http.NewRequest("PATCH", fullurl, strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

//本地注释：测试关闭节点指定通道
//本地注释：分四种情况 一是不存在的通道 二是opened的通道 三是已经closed的通道，四是settled的的通道
func CloseChannelTest(url string) {
	var err error
	var ChannelAddress string
	var i int
	var Status string
	start := time.Now()
	ShowTime()
	fmt.Println("Start Close Channel")

	//本地注释：关闭一个不存在的通道
	ChannelAddress = "0x00000"
	Status, err = CloseChannel(url, ChannelAddress)
	ShowError(err)
	ShowCloseChannelMsgDetail(Status)
	if Status == "409 Conflict" {
		fmt.Println("Test pass:Close a not exist Channel")
	} else {
		fmt.Println("Test failed:Close a not exist Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}

	//本地注释：关闭一个open通道
	Channels, _, _ := QueryingNodeAllChannels(url)
	//本地注释：查询第一个open通道
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
		fmt.Println("Test pass:Close a opened Channel")
	} else {
		fmt.Println("Test failed:Close a opened Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	//本地注释：关闭一个closed通道 关闭一个settled通道 长时间没有响应，跳过
	goto EndTest
Testclosed:
	//本地注释：关闭一个closed通道
	//本地注释：查询第一个closed通道
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
		fmt.Println("Test pass:Close a closed Channel")
	} else {
		fmt.Println("Test failed:Close a closed Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
Testsettled:
	//本地注释：关闭一个settled通道
	//本地注释：查询第一个settled通道
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
		fmt.Println("Test pass:Close a settled Channel")
	} else {
		fmt.Println("Test failed:Close a settled Channel")
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
EndTest:
	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowCloseChannelMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Close Channel Success!")
	case "400 Bad Request":
		fmt.Println("The provided json is in some way malformed!")
	case "409 Conflict":
		fmt.Println("Provided channel does not exist")
	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,Close Channel Failure! %s\n", Status)
	}
}
