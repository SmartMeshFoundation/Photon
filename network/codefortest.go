package network

import (
	"crypto/ecdsa"
	"time"

	"fmt"

	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/ethereum/go-ethereum/crypto"
)

//dummyProtocol only print received message
type dummyProtocol struct {
	name string
	data chan []byte
}

func newDummyProtocol(name string) *dummyProtocol {
	return &dummyProtocol{
		name: name,
		data: make(chan []byte, 20),
	}
}
func (p *dummyProtocol) receive(data []byte) {
	log.Debug(fmt.Sprintf("%s receive  data len=%d,data=\n%s", p.name, len(data), hex.Dump(data)))
	p.data <- data
}

//MakeTestUDPTransport test only
func MakeTestUDPTransport(name string, port int) *UDPTransport {
	t, err := NewUDPTransport(name, "127.0.0.1", port, nil, NewTokenBucket(10, 2, time.Now))
	if err != nil {
		panic(err)
	}
	return t
}

//MakeTestXMPPTransport create a test xmpp transport
func MakeTestXMPPTransport(name string, key *ecdsa.PrivateKey) *XMPPTransport {
	return NewXMPPTransport(name, params.DefaultTestXMPPServer, key, DeviceTypeOther)
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
func MakeTestRaidenProtocol(name string) *RaidenProtocol {
	privkey, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(MakeTestXMPPTransport(name, privkey), privkey, &testBlockNumberGetter{})
	return rp
}

//MakeTestDiscardExpiredTransferRaidenProtocol test only
func MakeTestDiscardExpiredTransferRaidenProtocol(name string) *RaidenProtocol {
	privkey, _ := crypto.GenerateKey()
	rp := NewRaidenProtocol(MakeTestXMPPTransport(name, privkey), privkey, newTimeBlockNumberGetter(time.Now()))
	return rp
}
