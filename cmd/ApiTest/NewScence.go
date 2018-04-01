package main

import (
	"fmt"
	"github.com/larspensjo/config"
	"time"
)

//本地注释：布置新场景
func NewScene() (NewTokenName string) {

	fmt.Println("==============================================================================================")
	fmt.Println("Start NewScene")
	c, err := config.ReadDefault("./ApiTest.INI")

	if err != nil {
		fmt.Println("config.ReadDefault error:", err)
		return
	}

	Node1Url, err := c.String("NOTE1", "api_address")
	Node2Url, err := c.String("NOTE2", "api_address")
	Node3Url, err := c.String("NOTE3", "api_address")
	Node4Url, err := c.String("NOTE4", "api_address")
	Node5Url, err := c.String("NOTE5", "api_address")
	Node6Url, err := c.String("NOTE6", "api_address")

	Node1Url = "http://" + Node1Url
	Node2Url = "http://" + Node2Url
	Node3Url = "http://" + Node3Url
	Node4Url = "http://" + Node4Url
	Node5Url = "http://" + Node5Url
	Node6Url = "http://" + Node6Url

	//本地注释：创建新Token并向账户充值

	EthRpcEndpoint, err := c.String("common", "eth_rpc_endpoint")
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	KeyStorePath, err := c.String("common", "keystore_path")
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	NewTokenName, RegistryAddress, _ := CreateNewToken(EthRpcEndpoint, KeyStorePath)
	//fmt.Println("New Token1=", NewTokenNames[0])
	//fmt.Println("New Token2=", NewTokenNames[1])
	//fmt.Println("registryAddress=", RegistryAddress.String())

	//本地注释：启动雷电客户端
	Startraiden(RegistryAddress.String())

	time.Sleep(10 * time.Second)

	//本地注释：测试注册新Token到雷电网
	Status, err := RegisteringOneToken(Node1Url, NewTokenName)
	ShowError(err)
	//本地注释：显示错误详细信息
	ShowRegisteringOneTokenMsgDetail(Status)
	switch Status {
	case "201 Created":
		fmt.Println("Success Registering a new token:", NewTokenName)
	default:
		fmt.Println("Failed  Registering new Token:", Status)
		//return -1
	}

	fmt.Println("Finish NewScene")

	return
}
