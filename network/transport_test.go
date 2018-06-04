package network

import (
	"testing"
	"time"

	"net"

	"bytes"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket(t *testing.T) {
	capacity := 2.0
	fillRate := 2.0
	tokenRefill := 1.0 / fillRate
	timeFunc := func() time.Time {
		return time.Unix(1, 0)
	}
	bucket := NewTokenBucket(capacity, fillRate, timeFunc)
	assert.Equal(t, bucket.Consume(1), time.Duration(0))
	assert.Equal(t, bucket.Consume(1), time.Duration(0))
	for i := 1; i < 9; i++ {
		assert.Equal(t, time.Duration(float64(i)*tokenRefill*float64(time.Second)), bucket.Consume(1))
	}
}

func TestUDPTransport(t *testing.T) {
	udp1 := MakeTestUDPTransport("u1", 40000)
	udp2 := MakeTestUDPTransport("u2", 40001)
	addr1 := utils.NewRandomAddress()
	addr2 := utils.NewRandomAddress()
	nodes := map[common.Address]*net.UDPAddr{
		addr1: udp1.UAddr,
		addr2: udp2.UAddr,
	}
	udp1.setHostPort(nodes)
	udp2.setHostPort(nodes)
	d1 := newDummyProtocol("u1")
	d2 := newDummyProtocol("u2")
	udp1.RegisterProtocol(d1)
	udp2.RegisterProtocol(d2)
	udp1.Start()
	udp2.Start()
	defer udp1.Stop()
	defer udp2.Stop()
	deviceType, isOnline := udp1.NodeStatus(addr1)
	if deviceType != DeviceTypeMobile || !isOnline {
		t.Error("node status error")
		return
	}
	deviceType, isOnline = udp1.NodeStatus(utils.NewRandomAddress())
	if isOnline {
		t.Error("should unkown")
		return
	}
	data := []byte("abc")
	err := udp1.Send(addr2, data)
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Millisecond * 100):
		t.Error("timeout")
	case data2 := <-d2.data:
		if !bytes.Equal(data2, data) {
			t.Error("not equal")
		}
	}
}
