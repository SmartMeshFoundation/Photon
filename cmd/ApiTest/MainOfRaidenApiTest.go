package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//deploy the new scence,register new token,register new raiden network,initiate raiden testing node
	NewTokenName := NewScene()
	//test transfer
	Transfer(NewTokenName, "./../../testdata/TransCase/case1.ini")
	//test API item by item
	//ApiTest()
}
