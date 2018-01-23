package network

import (
	"time"

	"fmt"

	"math/rand"

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
	return NewUDPTransportWithHostPort("192.168.0.102", port, nil, NewTokenBucket(10, 2, time.Now))
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
