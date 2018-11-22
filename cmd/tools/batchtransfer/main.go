package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SmartMeshFoundation/Photon/internal/debug"

	"github.com/SmartMeshFoundation/Photon/log"
	"gopkg.in/urfave/cli.v1"
)

/*
batch transfer is a tool for test transfer
*/

//test for closing the specified channel for the node
func transfer(urlstr, tokenAddr, target string, amount, identifier int, isdirect bool) (result string, err error) {
	var payload string
	payload = "{\"amount\":%d,\"is_direct\":%v,\"sync\":true}"
	payload = fmt.Sprintf(payload, amount, isdirect)
	//fmt.Printf("payload=%s\n", payload)
	// MatrixHTTPClient is a custom http client
	var client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err2 := net.DialTimeout(netw, addr, time.Second*3)
				if err2 != nil {
					//fmt.Println("dail timeout", err2)
					return nil, err2
				}
				return c, nil
			},
			Proxy: func(_ *http.Request) (*url.URL, error) {
				proxyurl := os.Getenv("http_proxy")
				if len(proxyurl) > 0 {
					return url.Parse(proxyurl)
				}
				return nil, nil
			},
			MaxIdleConnsPerHost:   100,
			ResponseHeaderTimeout: time.Second * 300000,
		},
	}
	fullurl := fmt.Sprintf("%s/api/1/transfers/%s/%s", urlstr, tokenAddr, target)
	log.Trace(fmt.Sprintf("transfer %d start @%s", identifier, time.Now().Format("15:04:05.000")))
	req, err := http.NewRequest("POST", fullurl, strings.NewReader(payload))
	if err != nil {
		log.Error(fmt.Sprintf("reqest err %s\n", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "name=anny")
	_, body, err := doRequest(client, req)
	if err != nil {
		log.Error(fmt.Sprintf("%d %d err %s\n", amount, identifier, err))
		return
	}
	result = string(body)
	//log.Printf("%d %d,response %s", amount, identifier, string(body))
	return
}

func doRequest(c *http.Client, req *http.Request) (Status string, body []byte, err error) {
	var buf [4096]byte
	var n = 0
	if req.Body != nil {
		n, err = req.Body.Read(buf[:])
		if err != nil {
			log.Error(fmt.Sprintf("req %s body err %s", req.URL.String(), err))
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(buf[:n]))
	}
	if n > 0 {
		log.Trace(fmt.Sprintf("send -> %s %s:\n%s\n", req.Method, req.URL.String(), string(buf[:n])))
	} else {
		log.Trace(fmt.Sprintf("send -> %s %s\n", req.Method, req.URL.String()))
	}
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	Status = resp.Status
	n, err = resp.Body.Read(buf[:])
	if err != nil && err != io.EOF {
		log.Error(fmt.Sprintf("req %s read err %s", req.URL.String(), err))
		return
	}
	body = buf[:n]
	if len(body) > 0 {
		//log.Printf("receive <- :\n%s\n", string(body))
	}
	err = resp.Body.Close()
	return
}
func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token",
			Usage: "transfer from tokens, it's a comma sperated token list. for example 0x57e682b80257aad73c4f3ad98d20435b4e1644d8762ef1ea1ff2806c27a5fa3d,0x56d05806656D3Ea8F2142a2FF2eA6dFE5D625254 ",
		},
		cli.StringFlag{
			Name:  "target",
			Usage: "transfer target",
		},
		cli.IntFlag{
			Name:  "number",
			Usage: `transfer number`,
		},
		cli.BoolFlag{
			Name:  "direct",
			Usage: "use direct transfer",
		},
		cli.StringFlag{
			Name:  "photon",
			Usage: "photon host port",
			Value: "http://127.0.0.1:5001",
		},
	}
	app.Flags = append(app.Flags, debug.Flags...)
	app.Action = mainctx
	app.Name = "batchtransfer"
	app.Version = "0.1"
	app.Before = func(ctx *cli.Context) error {
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Trace(fmt.Sprintf("run err %s\n", err))
	}
}

func mainctx(ctx *cli.Context) {
	err := debug.Setup(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("setup err=%s", err))
		return
	}
	tokens := strings.Split(ctx.String("token"), ",")
	if len(tokens) == 0 {
		log.Error(fmt.Sprintf("token arg error"))
		return
	}
	log.Info(fmt.Sprintf("tokens=%s", tokens))
	tokensNumber := len(tokens)
	number := ctx.Int("number")
	wg := sync.WaitGroup{}
	wg.Add(number)
	wg2 := sync.WaitGroup{}
	wg2.Add(number)
	start := time.Now()
	for i := 1; i <= number; i++ {
		go func(index int) {
			wg.Done()
			wg.Wait()
			start2 := time.Now()
			result, err := transfer(ctx.String("photon"), tokens[index%tokensNumber], ctx.String("target"), index, index, ctx.Bool("direct"))
			end := time.Now()
			if err != nil {
				log.Error(fmt.Sprintf("transfer:%d finished err=%s, result=%s, take time=%s", index, err, result, end.Sub(start2)))
			} else {
				log.Trace(fmt.Sprintf("transfer:%d finished err=%s, result=%s, take time=%s", index, err, result, end.Sub(start2)))
			}
			wg2.Done()
		}(i)
	}
	wg2.Wait()
	log.Info("all finished\n")
	duration := time.Now().Sub(start)

	log.Info(fmt.Sprintf("tps=%f", float64(number)/float64(float64(duration)/float64(time.Second))))
}
