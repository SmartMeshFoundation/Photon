package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	log.Println("Start Querying General Network Events")
	Events, Status, err := QueryingGeneralNetworkEvents(url)
	ShowError(err)
	//display the details of the error
	ShowQueryingGeneralNetworkEventsMsgDetail(Status)
	if Events != nil {
		//for i := 0; i < len(Events); i++ {
		//	log.Println("EventType:", Events[i].EventType, " TokenAddress:", Events[i].TokenAddress, " ChannelManagerAddress:", Events[i].ChannelManagerAddress)
		//}
	}

	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowQueryingGeneralNetworkEventsMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Successful query")
	case "400 Bad Request":
		log.Println("The provided query string is malformed")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingGeneralNetworkEvents Failure:%s\n", Status)
	}
}
