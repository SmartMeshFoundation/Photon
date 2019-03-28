package network

import (
	"testing"
	"time"

	"github.com/SmartMeshFoundation/Photon/params"

	"bytes"

	"github.com/SmartMeshFoundation/Photon/utils"
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
	addr1 := utils.NewRandomAddress()
	addr2 := utils.NewRandomAddress()
	udp1 := MakeTestUDPTransport(addr1.String(), 40000)
	udp2 := MakeTestUDPTransport(addr2.String(), 40001)

	d1 := newDummyProtocol("u1")
	d2 := newDummyProtocol("u2")
	udp1.RegisterProtocol(d1)
	udp2.RegisterProtocol(d2)
	udp1.Start()
	udp2.Start()
	defer udp1.Stop()
	defer udp2.Stop()
	time.Sleep(params.DefaultMDNSQueryInterval * 2)
	deviceType, isOnline := udp1.NodeStatus(addr2)
	if deviceType != DeviceTypeOther || !isOnline {
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
	case <-time.After(params.DefaultMDNSQueryInterval):
		t.Error("timeout")
	case data2 := <-d2.data:
		if !bytes.Equal(data2, data) {
			t.Error("not equal")
		}
	}
}

func TestUDPTransportWithMDNS(t *testing.T) {
	addr1 := utils.NewRandomAddress()
	addr2 := utils.NewRandomAddress()
	udp1 := MakeTestUDPTransport(addr1.String(), 40000)
	udp2 := MakeTestUDPTransport(addr2.String(), 40001)

	d1 := newDummyProtocol("u1")
	d2 := newDummyProtocol("u2")
	udp1.RegisterProtocol(d1)
	udp2.RegisterProtocol(d2)
	udp1.Start()
	udp2.Start()
	defer udp1.Stop()
	defer udp2.Stop()
	time.Sleep(params.DefaultMDNSQueryInterval * 2) //休息一段时间,等待对方上线通知
	deviceType, isOnline := udp1.NodeStatus(addr2)
	if deviceType != DeviceTypeOther || !isOnline {
		t.Error("node status error")
		return
	}
	deviceType, isOnline = udp2.NodeStatus(addr1)
	if deviceType != DeviceTypeOther || !isOnline {
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
