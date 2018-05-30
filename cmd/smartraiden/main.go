package main

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl"
)

func main() {
	if err := mainimpl.StartMain(); err != nil {
		fmt.Printf("quit with err %s", err)
	}
}
