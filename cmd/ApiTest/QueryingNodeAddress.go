package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"time"
)

func QueryingNodeAddress(url string) (Address NodeAddressPayload, Status string, err error) {
	var resp *http.Response
	var count int
	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1/address")
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if count >= MaxTry {
		Status = "504 TimeOut"
		Address.OurAddress = ""
	} else {
		p, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(p))
		json.Unmarshal(p, &Address)
		Status = resp.Status
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//本地注释：测试查询某节点地址
func QueryingNodeAddressTest(url string) {
	start := time.Now()
	ShowTime()
	fmt.Println("Start Querying Node Address")
	Address, Status, err := QueryingNodeAddress(url)
	ShowError(err)
	ShowQueryingNodeAddressMsgDetail(Status)
	switch Status {
	case "200 OK":
		fmt.Printf("Test pass:Querying Node Address Success!Node Address=%s\n", Address.OurAddress)
	default:
		fmt.Printf("Test failed:Querying Node Address Failure! %s\n", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowQueryingNodeAddressMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successful query")
	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Println("Unknown error,QueryingNodeAddress Failure:", Status)
	}
}
