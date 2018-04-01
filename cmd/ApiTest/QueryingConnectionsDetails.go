package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//本地注释：查询Token网络连接详情
func QueryingConnectionsDetails(url string) (Infos map[string]*ConnectionsDetails, Status string, err error) {
	var resp *http.Response
	var count int

	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1/connections")
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
		Infos = nil
	} else {
		if resp.Status == "200 OK" {
			p, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(p, &Infos)
		} else {
			Infos = nil
		}
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//本地注释：测试查询Token网络连接详情
func QueryingConnectionsDetailsTest(url string) {
	var err error
	var Status string
	Infos := make(map[string]*ConnectionsDetails)
	start := time.Now()
	ShowTime()
	fmt.Println("Start Querying Connecions Details")
	Infos, Status, err = QueryingConnectionsDetails(url)
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowQueryingConnectionsDetailsMsgDetail(Status)
	if Infos != nil {
		//for k, v := range Infos {
		//	fmt.Println("Token:", k, " Funds:", v.Funds, " SumDeposits:", v.SumDeposits, " Channels:", v.Channels)
		//}
	}
	switch Status {
	case "200 OK":
		fmt.Println("Test pass:QueryingConnectionsDetails")

	default:
		fmt.Printf("Test failed:QueryingConnectionsDetails:%s\n", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowQueryingConnectionsDetailsMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successful query")
	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingConnectionsDetails Failure:%s\n", Status)
	}
}
