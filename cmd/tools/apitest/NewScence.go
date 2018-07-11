package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/huamou/config"
)

//deploy the new scence
func NewScene() (NewTokenName string) {

	log.Println("==============================================================================================")
	log.Println("Start NewScene")
	c, err := config.ReadDefault("./apitest.ini")

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

	//create a new token and deposit to the account

	EthRpcEndpoint := c.RdString("common", "eth_rpc_endpoint", "ws://127.0.0.1:8546")

	KeyStorePath := c.RdString("common", "keystore_path", "/smtwork/privnet3/data/keystore")
	conn, err := ethclient.Dial(EthRpcEndpoint)
	if err != nil {
		log.Fatal(err)
		return
	}
	registryAddress := c.RdString("common", "registry_contract_address", "")
	NewTokenName = CreateTokenAndChannels(KeyStorePath, conn, common.HexToAddress(registryAddress), true)
	log.Println("New TokenNetworkAddres=", NewTokenName)

	//start the raiden client
	datadir := c.RdString("common", "datadir", "/smtwork/share/.smartraiden")
	os.RemoveAll(datadir)

	//time.Sleep(10 * time.Second)

	//test for registering new token to raiden network
	RegisteringOneToken(Node1Url, NewTokenName)
	//Status, err := RegisteringOneToken(Node1Url, NewTokenName)
	//ShowError(err)
	////display the details of the error
	//ShowRegisteringOneTokenMsgDetail(Status)
	//switch Status {
	//case "201 Created":
	//	log.Println("Success Registering a new token:", NewTokenName)
	//default:
	//	log.Println("Failed  Registering new TokenNetworkAddres:", Status)
	//	os.Exit(-1)
	//}

	log.Println("Finish NewScene")

	return
}
