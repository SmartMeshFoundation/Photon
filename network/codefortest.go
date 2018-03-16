package network

import (
	"time"

	"fmt"

	"math/rand"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

var testdiscovery DiscoveryInterface

//var testtransport Transporter

func init() {
	testdiscovery = NewDiscovery()
	//testtransport = MakeTestUDPTransport(rand.Intn(50000))
}

type DummyProtocol struct {
}

func (this *DummyProtocol) Receive(data []byte, host string, port int) {
	log.Debug(fmt.Sprintf("receive from %s:%d data len=%d", host, port, len(data)))
}
func MakeTestUDPTransport(port int) *UDPTransport {
	//todo how to get my ip
	return NewUDPTransportWithHostPort("127.0.0.1", port, nil, NewTokenBucket(10, 2, time.Now))
}

func GetTestDiscovery() DiscoveryInterface {
	return testdiscovery
}
func MakeTestRaidenProtocol() *RaidenProtocol {
	port := rand.New(utils.RandSrc).Intn(50000)
	privkey, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(MakeTestUDPTransport(port), testdiscovery, privkey)
	testdiscovery.Register(rp.nodeAddr, "127.0.0.1", port)
	return rp
}

func NewTestIceTransport(key *ecdsa.PrivateKey, name string) *IceTransport {
	InitIceTransporter("182.254.155.208:3478", "bai", "bai", "119.28.43.121:5222")
	it := NewIceTransporter(key, name)
	return it
}
func MakeTestIceRaidenProtocol(name string) *RaidenProtocol {
	key, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(NewTestIceTransport(key, name), NewIceHelperDiscovery(), key)
	return rp
}
