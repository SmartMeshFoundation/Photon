package network

import (
	"testing"

	"time"

	"os"

	"errors"

	"math/big"

	"sync"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
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
	var err error
	log.Trace("log...")
	p2 := MakeTestRaidenProtocol("p2")
	p1 := MakeTestRaidenProtocol("p1")

	//err := SetMatrixDB(p1, p2.nodeAddr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//err = SetMatrixDB(p2, p1.nodeAddr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	p1.Start()
	p2.StopAndWait()
	ping := encoding.NewPing(32)
	ping.Sign(p1.privKey, ping)
	err = p1.SendAndWait(p2.nodeAddr, ping, time.Second*2)
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
	revealSecretMsg := encoding.NewRevealSecret(utils.ShaSecret([]byte{12}))
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
	if revealSecretMsg.LockSecret != revealSecretMsg2.LockSecret {
		t.Errorf("secret not match")
	}
}

func TestNew(t *testing.T) {
	msger := encoding.MessageMap[encoding.UnlockCmdID]
	msg := New(msger)
	spew.Dump(msg)
	switch m2 := msg.(type) {
	case *encoding.UnLock:
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
	revealSecretMsg := encoding.NewRevealSecret(utils.ShaSecret([]byte{12}))
	revealSecretMsg.Sign(p1.privKey, revealSecretMsg)
	go func() {
		m := <-p2.ReceivedMessageChan
		t.Logf("client2 received msg :%#v", m)
		msg = m.Msg
		p2.ReceivedMessageResultChan <- nil
		secretRequest := encoding.NewSecretRequest(utils.EmptyHash, big.NewInt(12))
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
	if revealSecretMsg.LockSecret != revealSecretMsg2.LockSecret {
		t.Errorf("secret not match")
	}
	wg.Wait()
}

func TestRaidenProtocolSendMediatedTransferExpired(t *testing.T) {
	log.Trace("log...")
	p1 := MakeTestDiscardExpiredTransferRaidenProtocol("p1")
	p1.Start()
	expiration := 7 //7 second
	lock := mtree.Lock{
		Expiration:     int64(expiration),
		Amount:         big.NewInt(10),
		LockSecretHash: utils.ShaSecret([]byte("test")),
	}
	reciever := utils.NewRandomAddress()
	bp := encoding.NewBalanceProof(1, utils.BigInt0, utils.EmptyHash, &contracts.ChannelUniqueID{
		ChannelIdentifier: utils.NewRandomHash(),
		OpenBlockNumber:   3,
	})
	mtr := encoding.NewMediatedTransfer(bp, &lock,
		utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0)
	mtr.Sign(p1.privKey, mtr)
	err := p1.SendAndWait(reciever, mtr, time.Second*5)
	fmt.Println(err)
	if err != errTimeout {
		t.Errorf("should time out but get %s", err)
		return
	}
	lock.Expiration = 3
	p1.ChannelStatusGetter = &testChannelStatusGetterInvalid{}
	mtr2 := encoding.NewMediatedTransfer(bp, &lock,
		utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0)
	mtr2.Sign(p1.privKey, mtr2)
	err = p1.SendAndWait(reciever, mtr2, time.Second*5)
	fmt.Println(err)
	if err != errExpired {
		t.Error(errors.New("should expired before timeout"))
		return
	}

}
