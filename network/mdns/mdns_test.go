package mdns

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/Photon/params"

	"github.com/stretchr/testify/require"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

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

	time.Sleep(params.DefaultMDNSQueryInterval * 2)
	req := require.New(t)
	//req.Len(n.m, 1, "found b ")
	req.NotNil(n.m["imb"], "found b ")
}
