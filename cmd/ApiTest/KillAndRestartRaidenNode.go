package main

import (
	"math/rand"
	"time"
)

func KillAndRestartStartraiden() {
	var pstr []string
	//公共参数
	//public parameter
	param := new(RaidenParam)
	param.datadir = "d:\\share\\goraiden"
	param.keystore_path = "d:\\privnet3\\data\\keystore"
	param.discovery_contract_address = "0x5f014DA6ea514405f641341e42aC0e61B8190653"
	param.registry_contract_address = "0xAEABE46207c1f31f44C3F5876383B808d4280456"
	param.password_file = "d:\\share\\pass"
	param.nat = "none"
	param.eth_rpc_endpoint = "ws://127.0.0.1:8546"
	param.conditionquit = "{\"QuitEvent\":\"RefundTransferRecevieAckxx}"
	param.debug = true

	r := rand.Intn(3)
	//节点1
	//NODE 1
	if r == 1 {
		pstr = append(pstr, "/im")
		pstr = append(pstr, "goraiden1.exe")
		pstr = append(pstr, "/f")
		Exec_shell("taskkill.exe", pstr)
		time.Sleep(8 * time.Second)
		param.api_address = "127.0.0.1:5001"
		param.listen_address = "127.0.0.1:40001"
		param.address = "0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa"
		pstr = param.getParam()
		//log.Println(pstr)
		Exec_shell("goraiden1.exe", pstr)
	}
	//节点2
	//NODE 2
	if r == 2 {
		pstr = append(pstr, "/im")
		pstr = append(pstr, "goraiden2.exe")
		pstr = append(pstr, "/f")
		Exec_shell("taskkill.exe", pstr)
		time.Sleep(8 * time.Second)
		param.api_address = "127.0.0.1:5002"
		param.listen_address = "127.0.0.1:40002"
		param.address = "0x33df901abc22dcb7f33c2a77ad43cc98fbfa0790"
		pstr = param.getParam()
		Exec_shell("goraiden2.exe", pstr)
	}
	//节点3
	//NODE 3
	if r == 3 {
		pstr = append(pstr, "/im")
		pstr = append(pstr, "goraiden3.exe")
		pstr = append(pstr, "/f")
		Exec_shell("taskkill.exe", pstr)
		time.Sleep(8 * time.Second)
		param.api_address = "127.0.0.1:5003"
		param.listen_address = "127.0.0.1:40003"
		param.address = "0x8c1b2e9e838e2bf510ec7ff49cc607b718ce8401"
		pstr = param.getParam()
		Exec_shell("goraiden3.exe", pstr)
	}
	time.Sleep(30 * time.Second)
}
