package network

import (
	"testing"

	"math/rand"

	"bytes"

	"time"

	"os"

	"errors"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}

//need a valid account on blockchain and it needs gas
func TestDiscovery(t *testing.T) {
	bcs := rpc.MakeTestBlockChainService()
	discover := NewContractDiscovery(bcs.NodeAddress, bcs.Client, bcs.Auth)
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
