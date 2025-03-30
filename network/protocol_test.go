package network

import (
	"testing"

	"github.com/SmartMeshFoundation/Photon/params"

	"time"

	"os"

	"errors"

	"math/big"

	"sync"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	params.InitForUnitTest()
}

func TestPhotonProtocolSendReceive(t *testing.T) {
	if testing.Short() {
		return
	}
	log.Trace("log...")
	p1, p1Key := MakeTestPhotonProtocol("p1")
	p2, p2Key := MakeTestPhotonProtocol("p2")
	p2Addr := crypto.PubkeyToAddress(p2Key.PublicKey)
	p1.Start(true)
	p2.Start(true)
	ping := encoding.NewPing(32)
	ping.Sign(p1Key, ping)
	err := p1.SendAndWait(p2Addr, ping, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestPhotonProtocolSendReceiveTimeout(t *testing.T) {
	if testing.Short() {
		return
	}
	var err error
	log.Trace("log...")
	p2, p2Key := MakeTestPhotonProtocol("p2")
	p1, p1Key := MakeTestPhotonProtocol("p1")
	p2Addr := crypto.PubkeyToAddress(p2Key.PublicKey)

	//err := SetMatrixDB(p1, p2Addr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//err = SetMatrixDB(p2, p1Addr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	p1.Start(true)
	p2.StopAndWait()
	ping := encoding.NewPing(32)
	ping.Sign(p1Key, ping)
	err = p1.SendAndWait(p2Addr, ping, time.Minute)
	if err == nil {
		t.Error(errors.New("should timeout"))
		return
	}
}
func TestPhotonProtocolSendReceiveNormalMessage(t *testing.T) {
	if testing.Short() {
		return
	}
	var msg encoding.SignedMessager
	p1, p1Key := MakeTestPhotonProtocol("p1")
	p2, p2Key := MakeTestPhotonProtocol("p2")
	p2Addr := crypto.PubkeyToAddress(p2Key.PublicKey)
	p1.Start(true)
	p2.Start(true)
	revealSecretMsg := encoding.NewRevealSecret(utils.ShaSecret([]byte{12}))
	revealSecretMsg.Sign(p1Key, revealSecretMsg)
	go func() {
		m := <-p2.ReceivedMessageChan
		t.Logf("received msg :%#v", m)
		msg = m.Msg
		p2.ReceivedMessageResultChan <- nil
	}()
	err := p1.SendAndWait(p2Addr, revealSecretMsg, time.Minute)
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

func TestPhotonProtocolSendReceiveNormalMessage2(t *testing.T) {
	if testing.Short() {
		return
	}
	var msg encoding.SignedMessager
	var wg = sync.WaitGroup{}
	p1, p1Key := MakeTestPhotonProtocol("p1")
	p2, p2Key := MakeTestPhotonProtocol("p2")
	p1Addr := crypto.PubkeyToAddress(p1Key.PublicKey)
	p2Addr := crypto.PubkeyToAddress(p2Key.PublicKey)
	p1.Start(true)
	p2.Start(true)
	revealSecretMsg := encoding.NewRevealSecret(utils.ShaSecret([]byte{12}))
	revealSecretMsg.Sign(p1Key, revealSecretMsg)
	go func() {
		m := <-p2.ReceivedMessageChan
		t.Logf("client2 received msg :%#v", m)
		msg = m.Msg
		p2.ReceivedMessageResultChan <- nil
		secretRequest := encoding.NewSecretRequest(utils.EmptyHash, big.NewInt(12))
		secretRequest.Sign(p2Key, secretRequest)
		err := p2.SendAndWait(p1Addr, secretRequest, time.Minute)
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
	err := p1.SendAndWait(p2Addr, revealSecretMsg, time.Minute)
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

func TestPhotonProtocolSendMediatedTransferExpired(t *testing.T) {
	if testing.Short() {
		return
	}
	log.Trace("log...")
	_, testOpenBlockNumber := (&testChannelStatusGetter{}).GetChannelStatus(utils.EmptyHash)
	p1, p1Key := MakeTestPhotonProtocol("p1")
	p1.Start(true)
	expiration := 7 //7 second
	lock := mtree.Lock{
		Expiration:     int64(expiration),
		Amount:         big.NewInt(10),
		LockSecretHash: utils.ShaSecret([]byte("test")),
	}
	reciever := utils.NewRandomAddress()
	bp := encoding.NewBalanceProof(1, utils.BigInt0, utils.EmptyHash, &contracts.ChannelUniqueID{
		ChannelIdentifier: utils.NewRandomHash(),
		OpenBlockNumber:   testOpenBlockNumber,
	})
	mtr := encoding.NewMediatedTransfer(bp, &lock,
		utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0, []common.Address{utils.NewRandomAddress()})
	mtr.Sign(p1Key, mtr)
	err := p1.SendAndWait(reciever, mtr, time.Minute)
	fmt.Println(err)
	if err != errTimeout {
		t.Errorf("should time out but get %s", err)
		return
	}
	lock.Expiration = 3
	p1.ChannelStatusGetter = &testChannelStatusGetterInvalid{}
	mtr2 := encoding.NewMediatedTransfer(bp, &lock,
		utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0, []common.Address{utils.NewRandomAddress()})
	mtr2.Sign(p1Key, mtr2)
	err = p1.SendAndWait(reciever, mtr2, time.Minute)
	fmt.Println(err)
	if err != errExpired {
		t.Error(errors.New("should expired before timeout"))
		return
	}

}
