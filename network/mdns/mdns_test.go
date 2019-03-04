package mdns

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SmartMeshFoundation/Photon/log"
)

type DiscoveryNotifee struct {
	m map[string]*net.UDPAddr
}

func (n *DiscoveryNotifee) HandlePeerFound(id string, addr *net.UDPAddr) {
	log.Info(fmt.Sprintf("peer found id=%s,addr=%s", id, addr))
	n.m[id] = addr
}

func TestMdnsDiscovery(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sa, err := NewMdnsService(ctx, 3000, "ima", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	sb, err := NewMdnsService(ctx, 3001, "imb", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	_ = sb

	n := &DiscoveryNotifee{
		m: make(map[string]*net.UDPAddr),
	}

	sa.RegisterNotifee(n)

	time.Sleep(time.Second * 15)
	req := require.New(t)
	//req.Len(n.m, 1, "found b ")
	req.NotNil(n.m["imb"], "found b ")
}
