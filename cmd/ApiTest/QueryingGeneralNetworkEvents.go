package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/*
The network registry used is the default registry. The default registry is preconfigured and can be edited from the raiden configuration file.
You can query for registry network events by making a GET request to the following endpoint.
*/
func QueryingGeneralNetworkEvents(url string) (Events []GeneralNetworkEvents, Status string, err error) {
	var resp *http.Response
	var count int
	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1/events/network")
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

func QueryingGeneralNetworkEventsTest(url string) {
	start := time.Now()
	ShowTime()
	fmt.Println("Start Querying General Network Events")
	Events, Status, err := QueryingGeneralNetworkEvents(url)
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowQueryingGeneralNetworkEventsMsgDetail(Status)
	if Events != nil {
		//for i := 0; i < len(Events); i++ {
		//	fmt.Println("EventType:", Events[i].EventType, " TokenAddress:", Events[i].TokenAddress, " ChannelManagerAddress:", Events[i].ChannelManagerAddress)
		//}
	}

	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：显示错误详细信息
func ShowQueryingGeneralNetworkEventsMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successful query")
	case "400 Bad Request":
		fmt.Println("The provided query string is malformed")
	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingGeneralNetworkEvents Failure:%s\n", Status)
	}
}
