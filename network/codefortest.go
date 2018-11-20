package network

import (
	"crypto/ecdsa"
	"math/rand"
	"time"

	"fmt"

	"encoding/hex"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/ethereum/go-ethereum/common"
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
func randomPort() int {
	/* #nosec */
	return rand.Int()%1000 + 40000
}

//MakeTestXMPPTransport create a test xmpp transport
func MakeTestXMPPTransport(name string, key *ecdsa.PrivateKey) *XMPPTransport {
	return NewXMPPTransport(name, params.DefaultTestXMPPServer, key, DeviceTypeOther)
}

//MakeTestMixTransport creat a test mix transport
func MakeTestMixTransport(name string, key *ecdsa.PrivateKey) *MixTransport {
	port := randomPort()
	t, err := NewMixTranspoter(name, params.DefaultTestXMPPServer, "127.0.0.1", port, key, nil, NewTokenBucket(10, 2, time.Now), DeviceTypeOther)
	if err != nil {
		panic(err)
	}
	log.Debug(fmt.Sprintf("udp listen in port %d\n", port))
	return t
}

type testChannelStatusGetter struct{}

func (t *testChannelStatusGetter) GetChannelStatus(channelIdentifier common.Hash) (int, int64) {
	return channeltype.StateOpened, 0
}

type testChannelStatusGetterInvalid struct{}

func (t *testChannelStatusGetterInvalid) GetChannelStatus(channelIdentifier common.Hash) (int, int64) {
	return channeltype.StateInValid, 0
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

//MakeTestPhotonProtocol test only
func MakeTestPhotonProtocol(name string) *PhotonProtocol {
	////#nosec
	privkey, _ := crypto.GenerateKey()
	rp := NewPhotonProtocol(MakeTestXMPPTransport(name, privkey), privkey, &testChannelStatusGetter{})
	return rp
}

//MakeTestDiscardExpiredTransferPhotonProtocol test only
func MakeTestDiscardExpiredTransferPhotonProtocol(name string) *PhotonProtocol {
	//#nosec
	privkey, _ := crypto.GenerateKey()
	rp := NewPhotonProtocol(MakeTestXMPPTransport(name, privkey), privkey, &testChannelStatusGetter{})
	return rp
}

//SubscribeNeighbor subscribe neighbor's online and offline status
func SubscribeNeighbor(p *PhotonProtocol, addr common.Address) error {
	xt := p.Transport.(*XMPPTransport)
	return xt.conn.SubscribeNeighbour(addr)
}
