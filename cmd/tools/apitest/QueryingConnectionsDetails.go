package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//query the details of the TokenNetworkAddres network connection
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

//test for querying the details of the TokenNetworkAddres network connection
func QueryingConnectionsDetailsTest(url string) {
	var err error
	var Status string
	Infos := make(map[string]*ConnectionsDetails)
	start := time.Now()
	ShowTime()
	log.Println("Start Querying Connecions Details")
	Infos, Status, err = QueryingConnectionsDetails(url)
	ShowError(err)
	//display the details of the error
	ShowQueryingConnectionsDetailsMsgDetail(Status)
	if Infos != nil {
		//for k, v := range Infos {
		//	log.Println("TokenNetworkAddres:", k, " Funds:", v.Funds, " SumDeposits:", v.SumDeposits, " Channels:", v.Channels)
		//}
	}
	switch Status {
	case "200 OK":
		log.Println("Test pass:QueryingConnectionsDetails")

	default:
		fmt.Printf("Test failed:QueryingConnectionsDetails:%s\n", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowQueryingConnectionsDetailsMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Successful query")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingConnectionsDetails Failure:%s\n", Status)
	}
}
