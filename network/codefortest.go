package network

import (
	"time"

	"fmt"

	"math/rand"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

var testdiscovery DiscoveryInterface

//var testtransport Transporter

func init() {
	testdiscovery = NewDiscovery()
	//testtransport = MakeTestUDPTransport(rand.Intn(50000))
}

//DummyProtocol only print received message
type DummyProtocol struct {
}

func (p *DummyProtocol) receive(data []byte, host string, port int) {
	log.Debug(fmt.Sprintf("receive from %s:%d data len=%d", host, port, len(data)))
}

//MakeTestUDPTransport test only
func MakeTestUDPTransport(port int) *UDPTransport {
	return newUDPTransportWithHostPort("127.0.0.1", port, nil, NewTokenBucket(10, 2, time.Now))
}

//GetTestDiscovery test only
func GetTestDiscovery() DiscoveryInterface {
	return testdiscovery
}

type testBlockNumberGetter struct{}

func (t *testBlockNumberGetter) GetBlockNumber() int64 {
	return 0
}

type timeBlockNumberGetter struct {
	t time.Time
}

//newTimeBlockNumberGetter test only
func newTimeBlockNumberGetter(t time.Time) *timeBlockNumberGetter {
	return &timeBlockNumberGetter{t}
}

//GetBlockNumber pseudo blockNumber by time
func (t *timeBlockNumberGetter) GetBlockNumber() int64 {
	/*
		assume 1s a block
	*/
	return int64(time.Now().Sub(t.t) / time.Second)
}

//MakeTestRaidenProtocol test only
func MakeTestRaidenProtocol() *RaidenProtocol {
	port := rand.New(utils.RandSrc).Intn(50000)
	privkey, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(MakeTestUDPTransport(port), testdiscovery, privkey, &testBlockNumberGetter{})
	testdiscovery.Register(rp.nodeAddr, "127.0.0.1", port)
	return rp
}

//MakeTestDiscardExpiredTransferRaidenProtocol test only
func MakeTestDiscardExpiredTransferRaidenProtocol() *RaidenProtocol {
	port := rand.New(utils.RandSrc).Intn(50000)
	privkey, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(MakeTestUDPTransport(port), testdiscovery, privkey, newTimeBlockNumberGetter(time.Now()))
	testdiscovery.Register(rp.nodeAddr, "127.0.0.1", port)
	return rp
}

//NewTestIceTransport test only
func NewTestIceTransport(key *ecdsa.PrivateKey, name string) *IceTransport {
	InitIceTransporter("182.254.155.208:3478", "bai", "bai", "139.199.6.114:5222")
	it, _ := NewIceTransporter(key, name)
	return it
}

//MakeTestIceRaidenProtocol test only
func MakeTestIceRaidenProtocol(name string) *RaidenProtocol {
	key, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(NewTestIceTransport(key, name), NewIceHelperDiscovery(), key, &testBlockNumberGetter{})
	return rp
}
