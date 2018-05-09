package network

import (
	"testing"
	"time"

	"bytes"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
	udp1 := MakeTestUDPTransport(40000)
	udp2 := MakeTestUDPTransport(40001)
	registercallback()
	udp1.RegisterProtocol(new(DummyProtocol))
	udp2.RegisterProtocol(new(DummyProtocol))
	err := udp1.Send(utils.EmptyAddress, udp2.Host, udp2.Port, []byte("abcdefg"))
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Millisecond * 10)
	if (len(lastsend) != len(lastreceive)) || (len(lastsend) != 1) {
		t.Errorf("send data error")
	}
	if !bytes.Equal(lastsend[0], lastreceive[0]) {
		t.Errorf("send receive error")
	}
}
