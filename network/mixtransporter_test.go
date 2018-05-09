package network

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewMixDiscovery(t *testing.T) {
	key, _ := crypto.GenerateKey()

	ts, d := NewMixTranspoter(key, "test", "127.0.0.0.1", 5001, nil, nil, &DummyPolicy{})
	b := ts.switchToIce()
	if b {
		t.Error("should fail because default is ice")
	}
	b = ts.switchToUdp()
	if !b {
		t.Error("should success")
	}
	b = ts.switchToUdp()
	if b {
		t.Error("should fail,because it's udp already")
	}
	b = d.switchToUdp()
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
