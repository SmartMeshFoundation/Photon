package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//exchange the token
func TokenSwaps(url string, TargetAddress string, Role string, SendingAmount int32, SendingToken string, ReceivingToken string, ReceivingAmount int32, sn int64) (Status string, err error, id int64) {

	var count int
	var resp *http.Response
	var payload TokenSwapsPayload

	if sn == 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sn = r.Int63n(9223372036854775807)
	}
	id = sn

	payload.Role = Role
	payload.SendingAmount = SendingAmount
	payload.SendingToken = SendingToken
	payload.ReceivingAmount = ReceivingAmount
	payload.ReceivingToken = ReceivingToken

	p, _ := json.Marshal(payload)

	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/token_swaps/" + TargetAddress + "/" + strconv.FormatInt(sn, 10)

		req, _ := http.NewRequest("PUT", fullurl, bytes.NewReader(p))
		//req.Header.Set("Content-Type", "application/json")
		//req.Header.Set("Cookie", "name=anny")
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

//test for exchanging the token,token number in each node is indefinite,but the ratio 2:1.
func TokenSwapsTest(url1 string, url2 string) {

	var err error
	var Status string
	var id int64
	start := time.Now()
	ShowTime()
	log.Println("Start Token Swaps")
	Address, _, _ := QueryingNodeAddress(url2)
	Tokens, _, _ := QueryingRegisteredTokens(url1)
	if len(Tokens) < 2 {
		log.Println("Registered Tokens <2")
		return
	}

	Status, err, id = TokenSwaps(url1, Address.OurAddress, "maker", 2, Tokens[0], Tokens[1], 1, 0)

	ShowError(err)
	ShowTokenSwapsMsgDetail(Status)
	if Status != "201 Created" && Status != "200 OK" {
		return
	}

	Address, _, _ = QueryingNodeAddress(url1)
	Status, err, id = TokenSwaps(url2, Address.OurAddress, "taker", 1, Tokens[1], Tokens[0], 2, id)

	ShowError(err)
	ShowTokenSwapsMsgDetail(Status)

	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//display the details of the error
func ShowTokenSwapsMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fallthrough
	case "201 Created":
		log.Println("Successful Creation!")
	case "400 Bad Request":
		log.Println("The provided json is in some way malformed!")
	case "408 Request Timeout":
		log.Println("The token swap operation times out!")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	default:
		log.Println("Unknown error:", Status)
	}
}
