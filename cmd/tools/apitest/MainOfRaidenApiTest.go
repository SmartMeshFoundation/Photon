package main

import (
	"flag"
	"log"
)

var dofast = flag.Bool("fast", true, "skip create token and channels")

func init() {
	flag.Parse()
}
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	Startraiden("")
	if !*dofast {
		//deploy the new scence,register new token,register new raiden network,initiate raiden testing node
		NewTokenName := NewScene()
		//test transfer
		Transfer(NewTokenName, "./../../testdata/TransCase/case1.ini")
	}
	//test API item by item
	ApiTest()
}
