package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//query the Partner address in the channel of special Token
func QueryingAllPartnersForOneToken(url string, token string) (Partners []TokenPartnerPayload, Status string, err error) {
	var resp *http.Response
	var count int
	for count = 0; count < MaxTry; count = count + 1 {
		resp, err = http.Get(url + "/api/1/tokens/" + token + "/partners")
		if err == nil {
			Status = resp.Status
			break
		}
		time.Sleep(time.Second)
	}
	if count >= MaxTry {
		Status = "504 Timeout"
	}
	if resp == nil {
		Partners = nil
	} else {
		if resp.Status == "200 OK" {
			p, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(p, &Partners)
		} else {
			Partners = nil
		}
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//test for querying the Partner address in the channel of special Token
func QueryingAllPartnersForOneTokensTest(url string) {
	start := time.Now()
	ShowTime()
	log.Println("Start Querying All Partners For a Existent TokenNetworkAddres")
	Tokens, _, _ := QueryingRegisteredTokens(url)
	//fmt.Printf("!!!!!!!!!!!!Token:%s\n", Tokens[0])
	if Tokens == nil {
		log.Println("Warning:No registered TokenNetworkAddres!")
	} else {
		_, Status, err := QueryingAllPartnersForOneToken(url, Tokens[0])
		ShowError(err)
		//display the details of the error
		ShowQueryingAllPartnersForOneTokenMsgDetail(Status)
		switch Status {
		case "200 OK":
			log.Println("Test pass:Querying All Partners For a Tokens Success!")
		default:
			fmt.Printf("Test failed:Querying All Partners For a Tokens Failure! %s\n", Status)
			if HalfLife {
				log.Fatal("HalfLife,exit")
			}
		}
	}

	log.Println("Start Querying All Partners For a Nonexistent TokenNetworkAddres")
	_, Status, err := QueryingAllPartnersForOneToken(url, "0x0000")
	ShowError(err)
	//display the details of the error
	ShowQueryingAllPartnersForOneTokenMsgDetail(Status)
	switch Status {
	case "404 Not Found":
		log.Println("Test pass:Querying All Partners For a Nonexistent Tokens Success!")
	default:
		fmt.Printf("Test failed:Querying All Partners For a Nonexistent Tokens Failure! %s\n", Status)
		if HalfLife {
			log.Fatal("HalfLife,exit")
		}
	}
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowQueryingAllPartnersForOneTokenMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Successful query")
	case "302 Redirect":
		log.Println("The user accesses the channel link endpoint")
	case "404 Not Found":
		log.Println("The token does not exist")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,QueryingAllPartnersForOneToken Failure:%s\n", Status)
	}
}
