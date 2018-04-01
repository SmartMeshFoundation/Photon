package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

//本地注释：有返回的交易，输入本节点url,Token地址，目标地址，交易数量，返回状态信息和error
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
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

//本地注释：有信息输出的交易，输入本节点url,Token地址，目标地址，交易数量，无返回值
func InitiatingTransfer2Msg(url string, Token string, TargetAddress string, Amount int32) {
	_, Status, err := InitiatingTransfer(url, Token, TargetAddress, Amount)
	ShowError(err)
	ShowInitiatingTransferMsgDetail(Status)
}

//本地注释：测试两个节点交易测试，输入两个交易节点的url地址
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
		fmt.Printf("Token:[%s]\n", Token)
		_, Status, err := InitiatingTransfer(url, Token, TargetAddress, Amount)
		ShowError(err)
		ShowInitiatingTransferMsgDetail(Status)
	}
	duration := time.Since(start)
	ShowTime()
	fmt.Println("time used:", duration.Nanoseconds()/1000000, " ms")
}

//本地注释：验证交易结果
func ResultJudge(TransferResult TransferResponse, Status string, err error, InitiatorAddress string, TargetAddress string, TokenAddress string, Amount int32) {
	if Status != "200 OK" {
		fmt.Println("Transfer failed:", Status)
	} else {
		var e bool
		e = false
		if TransferResult.InitiatorAddress != InitiatorAddress {
			fmt.Println("Transfer failed:InitiatorAddress mismatching")
			e = true
		}
		if TransferResult.TargetAddress != TargetAddress {
			fmt.Println("Transfer failed:TargetAddress mismatching")
			e = true
		}
		if TransferResult.TokenAddress != TokenAddress {
			fmt.Println("Transfer failed:TokenAddress mismatching")
			e = true
		}
		if TransferResult.Amount != Amount {
			fmt.Println("Transfer failed:Amount mismatching")
			e = true
		}
		if e == false {
			fmt.Println("Transfer Success")
		}
	}
}

//本地注释：显示错误详细信息
func ShowInitiatingTransferMsgDetail(Status string) {
	switch Status {
	case "200 OK":
		fmt.Println("Successful Transfer creation.")
	case "400 Bad Request":
		fmt.Println("The provided json is in some way malformed")
	case "402 Payment required":
		fmt.Println("The transfer can’t start due to insufficient balance")
	case "408 Timeout":
		fmt.Println("A timeout happened during the transfer")
	case "409 Conflict":
		fmt.Println("The address or the amount is invalid or if there is no path to the target")
	case "500 Server Error":
		fmt.Println("Internal Raiden node error")
	case "504 TimeOut":
		fmt.Println("No response,timeout")
	default:
		fmt.Printf("Unknown error,InitiatingTransfer Failure:%s\n", Status)
	}
}
