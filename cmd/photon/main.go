package main

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl"
)

func main() {
	if _, err := mainimpl.StartMain(); err != nil {
		fmt.Printf("quit with err %s\n", err)
	}
}
