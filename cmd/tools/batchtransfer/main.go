package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"fmt"
	"time"
)

/*
batch transfer is a tool for test transfer
*/

//test for closing the specified channel for the node
func transfer(url, tokenAddr, target string, amount, identifier int) (err error) {
	var payload string
	payload = "{\"amount\":%d,\"identifier\":%d}"
	payload = fmt.Sprintf(payload, amount, identifier)
	client := &http.Client{}
	fullurl := fmt.Sprintf("%s/api/1/transfers/%s/%s", url, tokenAddr, target)
	req, err := http.NewRequest("POST", fullurl, strings.NewReader(payload))
	if err != nil {
		log.Printf("reqest err %s\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "name=anny")
	_, body, err := doRequest(client, req)
	if err != nil {
		log.Printf("%d %d err %s\n", amount, identifier, err)
		return
	}
	log.Printf("%d %d,response %s", amount, identifier, string(body))
	return
}

func doRequest(c *http.Client, req *http.Request) (Status string, body []byte, err error) {
	var buf [4096]byte
	var n = 0
	if req.Body != nil {
		n, err = req.Body.Read(buf[:])
		if err != nil {
			log.Printf("req %s body err %s", req.URL.String(), err)
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(buf[:n]))
	}
	if n > 0 {
		log.Printf("send -> %s %s:\n%s\n", req.Method, req.URL.String(), string(buf[:n]))
	} else {
		log.Printf("send -> %s %s\n", req.Method, req.URL.String())
	}
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	Status = resp.Status
	n, err = resp.Body.Read(buf[:])
	if err != nil {
		log.Printf("req %s read err %s", req.URL.String(), err)
	}
	body = buf[:n]
	if len(body) > 0 {
		log.Printf("receive <- :\n%s\n", string(body))
	}
	err = req.Body.Close()
	err = resp.Body.Close()
	return
}
func main() {
	for i := 1; i <= 10; i++ {
		log.Printf("start %d\n", i)
		go transfer("http://127.0.0.1:5001", "0x2c6978089905bbE7437e0071294A54852E5666D6", "0x33Df901ABc22DcB7F33c2a77aD43CC98FbFa0790", i, i)
	}
	time.Sleep(time.Second * 5)
	log.Printf("finished\n")
}
