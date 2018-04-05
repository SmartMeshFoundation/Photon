package main

import (
	"encoding/json"
	//"github.com/labstack/gommon/log"
	"io/ioutil"
	"log"
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
		//log.Println(string(p))
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
	log.Println("Start Querying Node Address")
	Address, Status, err := QueryingNodeAddress(url)
	ShowError(err)
	ShowQueryingNodeAddressMsgDetail(Status)
	switch Status {
	case "200 OK":
		log.Println("Test pass:Querying Node Address Success!Node Address=", Address.OurAddress)
	default:
		log.Println("Test failed:Querying Node Address Failure!", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowQueryingNodeAddressMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Successful query")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		log.Println("Unknown error,QueryingNodeAddress Failure:", Status)
	}
}
