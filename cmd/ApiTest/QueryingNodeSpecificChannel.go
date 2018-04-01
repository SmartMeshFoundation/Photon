package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//本地注释：查询某节点指定通道
func QueryingNodeSpecificChannel(url string, channel string) (Channel NodeChannel, Status string, err error) {
	var resp *http.Response
	var count int
	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1/channels/" + channel)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if count >= MaxTry {
		Status = "504 TimeOut"
	}
	if resp != nil {
		p, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(p, &Channel)
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//本地注释：测试查询某节点指定通道
func QueryingNodeSpecificChannelTest(url string) {

	var existedchannel string
	existedchannel = ""
	start := time.Now()
	ShowTime()
	fmt.Println("Start Querying Node Existed Specific Channel")
	Channels, _, _ := QueryingNodeAllChannels(url)
	if Channels != nil {
		if len(Channels) >= 1 {
			existedchannel = Channels[0].ChannelAddress
		}
	}
	//fmt.Printf("Existed Specific Channel:%s\n", existedchannel)
	_, Status, err := QueryingNodeSpecificChannel(url, existedchannel)
	ShowError(err)
	ShowQueryingNodeSpecificChannelMsgDetail(Status)
	switch Status {
	case "200 OK":
		fmt.Println("Test pass:querying node existed Specific channel Success!")
	default:
		fmt.Println("Test failed:querying node1 existed channels Failure:", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	fmt.Println("Start Querying Node  not existed Specific Channel")
	_, Status, err = QueryingNodeSpecificChannel(url, "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	ShowError(err)
	ShowQueryingNodeSpecificChannelMsgDetail(Status)
	switch Status {
	case "404 Not Found":
		fmt.Println("Test pass:querying node not existed Specific channel Success!")
	default:
		fmt.Println("Test failed:querying node1 not existed channels Failure:", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowQueryingNodeSpecificChannelMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successful Query")
	case "404 Not Found":
		fmt.Println("The channel does not exist")
	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Println("Unknown error,QueryingNodeSpecificChannel:", Status)
	}
}
