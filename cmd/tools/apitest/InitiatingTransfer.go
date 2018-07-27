package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//a transaction which return values, input the url of the node,token address,target address,amount,return the status information and error
func InitiatingTransfer(url string, Token string, TargetAddress string, Amount int32) (TransferResult TransferResponse, Status string, err error) {
	var resp *http.Response
	var count int
	var payload TransferRequest
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	Intsn := r.Int63n(9223372036854775807)
	payload.Amount = Amount
	payload.Identifier = Intsn
	p, _ := json.Marshal(payload)
	for count = 0; count < MaxTry; count = count + 1 {
		client := &http.Client{}
		fullurl := url + "/api/1/transfers/" + Token + "/" + TargetAddress
		req, _ := http.NewRequest(http.MethodPost, fullurl, bytes.NewReader(p))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "name=anny")
		resp, err = client.Do(req)
		if err == nil {
			if resp != nil {
				p, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(p, &TransferResult)
			}
			break
		}
		time.Sleep(time.Second)
	}

	if count >= MaxTry {
		Status = "504 TimeOut"
	}
	if resp != nil {
		Status = resp.Status
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return
}

//a transaction which output information, input the url of the node,token address,target address,amount,no return value
func InitiatingTransfer2Msg(url string, Token string, TargetAddress string, Amount int32) {
	_, Status, err := InitiatingTransfer(url, Token, TargetAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
}

//test the transfer between two nodes,input the urls of the two nodes
func InitiatingTransferTest(url string, TargetUrl string) {
	var Token string
	var TargetAddress string
	var Amount int32
	start := time.Now()
	ShowTime()
	fmt.Printf("Start Initiating Transfer from  [%s] to  [%s]\n", url, TargetUrl)
	TargetNodeAddress, _, _ := QueryingNodeAddress(TargetUrl)
	TargetAddress = TargetNodeAddress.OurAddress
	fmt.Printf("TargetAddress:[%s]\n", TargetAddress)
	Tokens, _, _ := QueryingRegisteredTokens(url)
	Amount = 1
	for i := 0; i < len(Tokens); i++ {
		Token = Tokens[i]
		fmt.Printf("TokenNetworkAddres:[%s]\n", Token)
		_, Status, err := InitiatingTransfer(url, Token, TargetAddress, Amount)
		ShowError(err)
		ShowInitiatingTransferMsgDetail(Status)
	}
	duration := time.Since(start)
	ShowTime()
	log.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//verify the result of the transfer
func ResultJudge(TransferResult TransferResponse, Status string, err error, InitiatorAddress string, TargetAddress string, TokenAddress string, Amount int32) {
	if Status != "200 OK" {
		log.Println("Transfer failed:", Status)
	} else {
		var e bool
		e = false
		if TransferResult.InitiatorAddress != InitiatorAddress {
			log.Println("Transfer failed:InitiatorAddress mismatching")
			e = true
		}
		if TransferResult.TargetAddress != TargetAddress {
			log.Println("Transfer failed:TargetAddress mismatching")
			e = true
		}
		if TransferResult.TokenAddress != TokenAddress {
			log.Println("Transfer failed:TokenAddress mismatching")
			e = true
		}
		if TransferResult.Amount != Amount {
			log.Println("Transfer failed:PaymentAmount mismatching")
			e = true
		}
		if e == false {
			log.Println("Transfer Success")
		}
	}
}

//display the details of the error
func ShowInitiatingTransferMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		log.Println("Successful Transfer creation.")
	case "400 Bad Request":
		log.Println("The provided json is in some way malformed")
	case "402 Payment required":
		log.Println("The transfer can’t start due to insufficient balance")
	case "408 Timeout":
		log.Println("A timeout happened during the transfer")
	case "409 Conflict":
		log.Println("The address or the amount is invalid or if there is no path to the target")
	case "500 Server Error":
		log.Println("Internal Raiden node error")
	case "504 TimeOut":
		log.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,InitiatingTransfer Failure:%s\n", Status)
	}
}
