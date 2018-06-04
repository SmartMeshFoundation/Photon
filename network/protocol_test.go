package network

import (
	"testing"

	"time"

	"os"

	"errors"

	"math/big"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/davecgh/go-spew/spew"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

func TestRaidenProtocolSendReceive(t *testing.T) {
	log.Trace("log...")
	p1 := MakeTestRaidenProtocol("p1")
	p2 := MakeTestRaidenProtocol("p2")
	p1.Start()
	p2.Start()
	ping := encoding.NewPing(32)
	ping.Sign(p1.privKey, ping)
	err := p1.SendAndWait(p2.nodeAddr, ping, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestRaidenProtocolSendReceiveTimeout(t *testing.T) {
	log.Trace("log...")
	p1 := MakeTestRaidenProtocol("p1")
	p2 := MakeTestRaidenProtocol("p2")
	p1.Start()
	ping := encoding.NewPing(32)
	ping.Sign(p1.privKey, ping)
	err := p1.SendAndWait(p2.nodeAddr, ping, time.Second*2)
	if err == nil {
		t.Error(errors.New("should timeout"))
		return
	}
}
func TestRaidenProtocolSendReceiveNormalMessage(t *testing.T) {
	var msg encoding.SignedMessager
	p1 := MakeTestRaidenProtocol("p1")
	p2 := MakeTestRaidenProtocol("p2")
	p1.Start()
	p2.Start()
	revealSecretMsg := encoding.NewRevealSecret(utils.Sha3([]byte{12}))
	revealSecretMsg.Sign(p1.privKey, revealSecretMsg)
	go func() {
		m := <-p2.ReceivedMessageChan
		t.Logf("received msg :%#v", m)
		msg = m.Msg
		p2.ReceivedMessageResultChan <- nil
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
	msger := encoding.MessageMap[encoding.SecretCmdID]
	msg := New(msger)
	spew.Dump(msg)
	switch m2 := msg.(type) {
	case *encoding.Secret:
		t.Log("m2 type correct ", m2.CmdID)
	default:
		t.Error("type convert error")
	}
}

func TestRaidenProtocolSendReceiveNormalMessage2(t *testing.T) {
	var msg encoding.SignedMessager
	var wg = sync.WaitGroup{}
	p1 := MakeTestRaidenProtocol("p1")
	p2 := MakeTestRaidenProtocol("p2")
	p1.Start()
	p2.Start()
	revealSecretMsg := encoding.NewRevealSecret(utils.Sha3([]byte{12}))
	revealSecretMsg.Sign(p1.privKey, revealSecretMsg)
	go func() {
		m := <-p2.ReceivedMessageChan
		t.Logf("client2 received msg :%#v", m)
		msg = m.Msg
		p2.ReceivedMessageResultChan <- nil
		secretRequest := encoding.NewSecretRequest(33, utils.EmptyHash, big.NewInt(12))
		secretRequest.Sign(p2.privKey, secretRequest)
		err := p2.SendAndWait(p1.nodeAddr, secretRequest, time.Minute)
		if err != nil {
			t.Error(err)
		}
	}()
	go func() {
		m := <-p1.ReceivedMessageChan
		t.Logf("client1 received msg:%#v", m)
		p1.ReceivedMessageResultChan <- nil
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
	p1 := MakeTestDiscardExpiredTransferRaidenProtocol("p1")
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
