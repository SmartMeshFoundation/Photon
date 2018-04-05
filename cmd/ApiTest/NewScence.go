package main

import (
	"log"
	"time"

	"github.com/larspensjo/config"
)

//本地注释：布置新场景
func NewScene() (NewTokenName string) {

	log.Println("==============================================================================================")
	log.Println("Start NewScene")
	c, err := config.ReadDefault("./ApiTest.INI")

	if err != nil {
		log.Println("config.ReadDefault error:", err)
		return
	}

	Node1Url := c.RdString("NODE1", "api_address", "127.0.0.1:5001")
	Node2Url := c.RdString("NODE2", "api_address", "127.0.0.1:5002")
	Node3Url := c.RdString("NODE3", "api_address", "127.0.0.1:5003")
	Node4Url := c.RdString("NODE4", "api_address", "127.0.0.1:5004")
	Node5Url := c.RdString("NODE5", "api_address", "127.0.0.1:5005")
	Node6Url := c.RdString("NODE6", "api_address", "127.0.0.1:5006")

	Node1Url = "http://" + Node1Url
	Node2Url = "http://" + Node2Url
	Node3Url = "http://" + Node3Url
	Node4Url = "http://" + Node4Url
	Node5Url = "http://" + Node5Url
	Node6Url = "http://" + Node6Url

	//本地注释：创建新Token并向账户充值

	EthRpcEndpoint := c.RdString("common", "eth_rpc_endpoint", "ws://127.0.0.1:8546")

	KeyStorePath := c.RdString("common", "keystore_path", "/smtwork/privnet3/data/keystore")

	NewTokenName, RegistryAddress, _ := CreateNewToken(EthRpcEndpoint, KeyStorePath)
	log.Println("New Token=", NewTokenName)
	log.Println("registryAddress=", RegistryAddress.String())

	//本地注释：启动雷电客户端
	Startraiden(RegistryAddress.String())

	time.Sleep(10 * time.Second)

	//本地注释：测试注册新Token到雷电网
	RegisteringOneToken(Node1Url, NewTokenName)
	//Status, err := RegisteringOneToken(Node1Url, NewTokenName)
	//ShowError(err)
	////本地注释：显示错误详细信息
	//ShowRegisteringOneTokenMsgDetail(Status)
	//switch Status {
	//case "201 Created":
	//	log.Println("Success Registering a new token:", NewTokenName)
	//default:
	//	log.Println("Failed  Registering new Token:", Status)
	//	os.Exit(-1)
	//}

	log.Println("Finish NewScene")

	return
}
