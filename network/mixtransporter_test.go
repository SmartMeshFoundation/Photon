package network

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewMixDiscovery(t *testing.T) {
	key, _ := crypto.GenerateKey()

	ts, d, _ := NewMixTranspoter(key, "test", "127.0.0.0.1", 5001, nil, nil, &dummyPolicy{})
	b := ts.switchToIce()
	if b {
		t.Error("should fail because default is ice")
	}
	b = ts.switchToUDP()
	if !b {
		t.Error("should success")
	}
	b = ts.switchToUDP()
	if b {
		t.Error("should fail,because it's udp already")
	}
	b = d.switchToUDP()
	if !b {
		t.Error("should success")
	}
	b = d.switchToIce()
	if !b {
		t.Error("should success")
	}
	b = d.switchToIce()
	if b {
		t.Error("should fail")
	}
}
