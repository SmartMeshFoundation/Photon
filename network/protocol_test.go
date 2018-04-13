package network

import (
	"testing"

	"math/rand"

	"bytes"

	"time"

	"os"

	"errors"

	"math/big"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

//need a valid account on blockchain and it needs gas
func TestDiscovery(t *testing.T) {
	bcs := rpc.MakeTestBlockChainService()
	discover := NewContractDiscovery(bcs.NodeAddress, params.ROPSTEN_DISCOVERY_ADDRESS, bcs.Client, bcs.Auth)
	host, port, err := discover.Get(bcs.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	host = "0.0.0.0"
	port = rand.New(utils.RandSrc).Intn(50000)
	if err := discover.Register(bcs.NodeAddress, host, port); err != nil {
		t.Error(err)
		return
	}
	newhost, newport, err := discover.Get(bcs.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	if host != newhost || newport != newport {
		t.Error("register Host Port failer")
		return
	}
}
func TestNewHttpDiscovery(t *testing.T) {
	return //http discovery has been obsolete
	dis := NewHttpDiscovery()
	host := "127.0.0.1"
	port := rand.New(utils.RandSrc).Intn(50000)
	addr := utils.NewRandomAddress()
	err := dis.Register(addr, host, port)
	if err != nil {
		t.Error(err)
	}
	host2, port2, err := dis.Get(addr)
	if err != nil || host2 != host || port2 != port {
		t.Error(err)
	}
	address, err := dis.NodeIdByHostPort(host, port)
	if err != nil || address != addr {
		t.Error(err)
	}

}

var lastreceive [][]byte
var lastsend [][]byte

func registercallback() {
	RegisterReceiveCallback(func(sender common.Address, hostport string, msg []byte) {
		lastreceive = append(lastreceive, msg)
	})
	RegisterSendCallback(func(sender common.Address, hostport string, msg []byte) {
		lastsend = append(lastsend, msg)
	})
}

func TestRaidenProtocolSendReceive(t *testing.T) {
	log.Trace("log...")
	p1 := MakeTestRaidenProtocol()
	p2 := MakeTestRaidenProtocol()
	registercallback()
	p1.Start()
	p2.Start()
	ping := encoding.NewPing(32)
	ping.Sign(p1.privKey, ping)
	err := p1.SendAndWait(p2.nodeAddr, ping, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	if len(lastsend) != 2 || len(lastreceive) != 2 {
		t.Error("send receive packet numer error")
		return
	}
	spew.Dump("send:", lastsend)
	spew.Dump("receive", lastreceive)
	if !bytes.Equal(lastsend[0], lastreceive[0]) {
		t.Error("first packet not match")
	}
	if !bytes.Equal(lastsend[1], lastreceive[1]) {
		t.Error("second packet not match")
	}
}
func TestRaidenProtocolSendReceiveTimeout(t *testing.T) {
	log.Trace("log...")
	p1 := MakeTestRaidenProtocol()
	p2 := MakeTestRaidenProtocol()
	registercallback()
	p1.Start()
	ping := encoding.NewPing(32)
	ping.Sign(p1.privKey, ping)
	err := p1.SendAndWait(p2.nodeAddr, ping, time.Second*2)
	if err == nil {
		t.Error(errors.New("should timeout"))
		return
	}
	if len(lastsend) != int(time.Second*2/p1.retryInterval)+1 || len(lastreceive) != 0 {
		t.Error("send receive packet numer error")
		//return
	}
	spew.Dump("send:", lastsend)
	spew.Dump("receive", lastreceive)
}
func TestRaidenProtocolSendReceiveNormalMessage(t *testing.T) {
	var msg encoding.SignedMessager
	p1 := MakeTestRaidenProtocol()
	p2 := MakeTestRaidenProtocol()
	p1.Start()
	p2.Start()
	revealSecretMsg := encoding.NewRevealSecret(utils.Sha3([]byte{12}))
	revealSecretMsg.Sign(p1.privKey, revealSecretMsg)
	go func() {
		m := <-p2.ReceivedMessageChannel
		t.Logf("received msg :%#v", m)
		msg = m.Msg
		p2.ReceivedMessageResultChannel <- nil
	}()
	err := p1.SendAndWait(p2.nodeAddr, revealSecretMsg, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	revealSecretMsg2, ok := msg.(*encoding.RevealSecret)
	if !ok {
		t.Errorf("recevied message type error")
		return
	}
	if revealSecretMsg.Secret != revealSecretMsg2.Secret {
		t.Errorf("secret not match")
	}
}

func TestNew(t *testing.T) {
	msger := encoding.MessageMap[encoding.SECRET_CMDID]
	msg := New(msger)
	spew.Dump(msg)
	switch m2 := msg.(type) {
	case *encoding.Secret:
		t.Log("m2 type correct ", m2.CmdId)
	default:
		t.Error("type convert error")
	}
}

func TestIceRaidenProtocolSendReceiveNormalMessage(t *testing.T) {
	var msg encoding.SignedMessager
	var wg = sync.WaitGroup{}
	p1 := MakeTestIceRaidenProtocol("client1")
	p2 := MakeTestIceRaidenProtocol("client2")
	p1.Start()
	p2.Start()
	revealSecretMsg := encoding.NewRevealSecret(utils.Sha3([]byte{12}))
	revealSecretMsg.Sign(p1.privKey, revealSecretMsg)
	go func() {
		m := <-p2.ReceivedMessageChannel
		t.Logf("client2 received msg :%#v", m)
		msg = m.Msg
		p2.ReceivedMessageResultChannel <- nil
		secretRequest := encoding.NewSecretRequest(33, utils.EmptyHash, big.NewInt(12))
		secretRequest.Sign(p2.privKey, secretRequest)
		err := p2.SendAndWait(p1.nodeAddr, secretRequest, time.Minute)
		if err != nil {
			t.Error(err)
		}
	}()
	go func() {
		m := <-p1.ReceivedMessageChannel
		t.Logf("client1 received msg:%#v", m)
		p1.ReceivedMessageResultChannel <- nil
		time.Sleep(time.Millisecond * 10)
		wg.Done()
	}()
	wg.Add(1)
	err := p1.SendAndWait(p2.nodeAddr, revealSecretMsg, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	revealSecretMsg2, ok := msg.(*encoding.RevealSecret)
	if !ok {
		t.Errorf("recevied message type error")
		return
	}
	if revealSecretMsg.Secret != revealSecretMsg2.Secret {
		t.Errorf("secret not match")
	}
	wg.Wait()
}

func TestRaidenProtocolSendMediatedTransferExpired(t *testing.T) {
	log.Trace("log...")
	p1 := MakeTestDiscardExpiredTransferRaidenProtocol()
	registercallback()
	p1.Start()
	expiration := 7 //7 second
	lock := encoding.Lock{
		Expiration: int64(expiration),
		Amount:     big.NewInt(10),
		HashLock:   utils.Sha3([]byte("test")),
	}
	reciever := utils.NewRandomAddress()
	mtr := encoding.NewMediatedTransfer(1, 1, utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0, reciever, utils.EmptyHash, &lock,
		utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0)
	mtr.Sign(p1.privKey, mtr)
	err := p1.SendAndWait(reciever, mtr, time.Second*5)
	if err != errTimeout {
		t.Errorf("should time out but get %s", err)
		return
	}
	lock.Expiration = 3
	mtr2 := encoding.NewMediatedTransfer(1, 1, utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0, reciever, utils.EmptyHash, &lock,
		utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0)
	mtr2.Sign(p1.privKey, mtr2)
	err = p1.SendAndWait(reciever, mtr2, time.Second*5)
	if err != errExpired {
		t.Error(errors.New("should expired before timeout"))
		return
	}

}
