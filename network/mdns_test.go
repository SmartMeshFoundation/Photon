package network

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/hashicorp/mdns"
)

func TestMDns(t *testing.T) {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"port is 5001"}
	service, err := mdns.NewMDNSService(host, "photon_node", "", "", 8000, nil, info)
	if err != nil {
		t.Error(err)
		return
	}

	// Create the mDNS server, defer shutdown
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %s\n", utils.StringInterface(entry, 3))
		}
	}()

	// Start the lookup
	mdns.Lookup("photon_node", entriesCh)
	close(entriesCh)

	time.Sleep(time.Second * 20)
}
