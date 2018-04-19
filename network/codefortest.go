package network

import (
	"time"

	"fmt"

	"math/rand"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
	return NewUDPTransportWithHostPort("127.0.0.1", port, nil, NewTokenBucket(10, 2, time.Now))
}

func GetTestDiscovery() DiscoveryInterface {
	return testdiscovery
}

type testBlockNumberGetter struct{}

func (t *testBlockNumberGetter) GetBlockNumber() int64 {
	return 0
}

type TimeBlockNumberGetter struct {
	t time.Time
}

func NewTimeBlockNumberGetter(t time.Time) *TimeBlockNumberGetter {
	return &TimeBlockNumberGetter{t}
}
func (t *TimeBlockNumberGetter) GetBlockNumber() int64 {
	/*
		assume 1s a block
	*/
	return int64(time.Now().Sub(t.t) / time.Second)
}
func MakeTestRaidenProtocol() *RaidenProtocol {
	port := rand.New(utils.RandSrc).Intn(50000)
	privkey, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(MakeTestUDPTransport(port), testdiscovery, privkey, &testBlockNumberGetter{})
	testdiscovery.Register(rp.nodeAddr, "127.0.0.1", port)
	return rp
}
func MakeTestDiscardExpiredTransferRaidenProtocol() *RaidenProtocol {
	port := rand.New(utils.RandSrc).Intn(50000)
	privkey, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(MakeTestUDPTransport(port), testdiscovery, privkey, NewTimeBlockNumberGetter(time.Now()))
	testdiscovery.Register(rp.nodeAddr, "127.0.0.1", port)
	return rp
}
func NewTestIceTransport(key *ecdsa.PrivateKey, name string) *IceTransport {
	InitIceTransporter("182.254.155.208:3478", "bai", "bai", "139.199.6.114:5222")
	it := NewIceTransporter(key, name)
	return it
}
func MakeTestIceRaidenProtocol(name string) *RaidenProtocol {
	key, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(NewTestIceTransport(key, name), NewIceHelperDiscovery(), key, &testBlockNumberGetter{})
	return rp
}
