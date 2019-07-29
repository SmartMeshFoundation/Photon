package mdns

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	sa, err := NewMdnsService(ctx, 3000, "ima", 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	sb, err := NewMdnsService(ctx, 3001, "imb", 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	_ = sb

	n := &DiscoveryNotifee{
		m: make(map[string]*net.UDPAddr),
	}

	sa.RegisterNotifee(n)

	time.Sleep(params.Cfg.MDNSQueryInterval * 2 * 5)
	req := require.New(t)
	//req.Len(n.m, 1, "found b ")
	req.NotNil(n.m["imb"], "found b ")
}

func TestGetBestIP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sa, err := NewMdnsService(ctx, 3000, "ima", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	m := sa.(*mdnsService)
	m.service.IPs = []net.IP{
		net.ParseIP("10.0.0.17").To4(),
		net.ParseIP("192.168.124.13").To4(),
	}
	remotes := []net.IP{
		net.ParseIP("192.168.124.11").To4(),
		net.ParseIP("192.168.122.1").To4(),
	}
	ip := m.getBestMatchIP(remotes)
	assert.EqualValues(t, ip.String(), net.ParseIP("192.168.124.11").String())
}
