package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//本地注释：查询某节点所有通道
func QueryingNodeAllChannels(url string) (Channels []NodeChannel, Status string, err error) {
	var resp *http.Response
	var count int
	var result bool
	result = true
	//Channels = nil
	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1/channels")
		if err == nil {
			Status = resp.Status
			//io.Copy(os.Stdout, resp.Body)
			if resp.Status != "200 OK" {
				result = false
			}
			break
		}
		time.Sleep(time.Second)
	}
	if count >= MaxTry {
		Status = "504 TimeOut"
		result = false
	}

	if result {
		p, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(p, &Channels)
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//本地注释：测试查询某节点所有通道
func QueryingNodeAllChannelsTest(url string) (Channels []NodeChannel) {
	start := time.Now()
	ShowTime()
	log.Println("Start Querying Node All Channels")
	Channels, Status, err := QueryingNodeAllChannels(url)
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowQueryingNodeAllChannelsMsgDetail(Status)
	switch Status {
	case "200 OK":
		log.Println("Test pass:querying node1 all channels Success!")
	default:
		log.Println("Test failed:querying node1 all channels Success!%s", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
	return
}

//本地注释：显示错误详细信息
func ShowQueryingNodeAllChannelsMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Successful Query")
	case "404 Not Found":
		log.Println("The channel does not exist")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingNodeAllChannels Failure:%s\n", Status)
	}
}
