package encoding

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/SmartMeshFoundation/Photon/params"

	"bytes"

	"errors"

	"reflect"

	"math/big"

	"encoding/json"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

const TestPrivkey = "4359f525e2b373089be5fe8f9a4e8ffb6d30e2960918be426217921e1b2547f7"

func init() {
	params.InitForUnitTest()
}
func GetTestPrivKey() *ecdsa.PrivateKey {
	key, _ := hex.DecodeString(TestPrivkey)
	privkey, _ := crypto.ToECDSA(key)
	return privkey
}

func GetTestPubKey() ecdsa.PublicKey {
	priv := GetTestPrivKey()
	return priv.PublicKey
}
func GetTestAddress() common.Address {
	return crypto.PubkeyToAddress(GetTestPubKey())
}

/*

def test_signature():
    ping = Ping(nonce=0)
    ping.Sign(PRIVKEY, ADDRESS)
    print binascii.b2a_hex(ping.encode())
    assert ping.sender == ADDRESS
*/

func TestSignature(t *testing.T) {
	ping := NewPing(0x33)
	ping.Signature = SignMessage(GetTestPrivKey(), ping)
	data := ping.Pack()
	ping2 := new(Ping)
	ping2.UnPack(data)
	if ping.Nonce != ping2.Nonce {
		t.Errorf("expect equal nonce but ping=%d,ping2=%d\n", ping.Nonce, ping2.Nonce)
	}
	if !bytes.Equal(ping.Signature, ping2.Signature) {
		t.Errorf("unequal signature for Ping")
	}
	sender, err := VerifyMessage(ping2.Pack())
	if err != nil {
		t.Errorf("sigature verify error")
	}
	testAddr := GetTestAddress()
	if !bytes.Equal(sender[:], testAddr[:]) {
		t.Errorf("sender's address is  error when extract from signature")
	}
	ping = new(Ping)
	if len(ping.Pack()) > 65 {
		t.Errorf("length error before signature")
	}
	err = ping.Sign(GetTestPrivKey(), ping)
	if err != nil {
		t.Error(err)
	}
	if len(ping.Pack()) < 65 {
		t.Errorf("length error after signature")
	}
	t.Log(hex.Dump(ping.Pack()))
}

func TestType(t *testing.T) {
	var p Messager = new(Ping)
	var pi interface{}
	pi = p
	if _, ok := pi.(*CmdStruct); ok {
		t.Log("is type  cmd struct")
	}
	//if _, ok := p.(*CmdStruct); ok {
	//	T.Log("struct is type cmd struct")
	//}
}
func TestEnvelopeMessage(t *testing.T) {
	bp := &BalanceProof{
		Nonce:             11,
		ChannelIdentifier: utils.Sha3([]byte("123")),
		TransferAmount:    big.NewInt(12),
		OpenBlockNumber:   3,
		Locksroot:         utils.EmptyHash,
	}
	p := NewDirectTransfer(bp)
	var sm SignedMessager = p
	err := p.Sign(GetTestPrivKey(), p)
	if err != nil {
		t.Error(err)
	}
	data := p.Pack()
	err = sm.verifySignature(data)
	if err != nil {
		t.Error(err)
	}
	p2 := new(DirectTransfer)
	err = p2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	if p2.ChannelIdentifier != p.ChannelIdentifier ||
		p.Locksroot != p2.Locksroot || p.Nonce != p.Nonce ||
		!bytes.Equal(p2.Signature, p.Signature) {
		t.Error(errors.New("data pack unpack error"))
	}
}

func TestHash(t *testing.T) {
	ping := NewPing(32)
	ping.Sign(GetTestPrivKey(), ping)
	data := ping.Pack()
	msgHash := utils.Sha3(data)
	ping2 := NewPing(0)
	ping2.UnPack(data)
	if utils.Sha3(ping2.Pack()) != msgHash {
		t.Error("hash error")
	}
}

func TestDirectTransfer(t *testing.T) {
	bp := &BalanceProof{
		Nonce:             11,
		ChannelIdentifier: utils.Sha3([]byte("123")),
		TransferAmount:    big.NewInt(12),
		OpenBlockNumber:   3,
		Locksroot:         utils.EmptyHash,
	}
	d1 := NewDirectTransfer(bp)
	d1.Data = []byte("123")
	d1.Sign(GetTestPrivKey(), d1)
	d2 := new(DirectTransfer)
	err := d2.UnPack(d1.Pack())
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, d1, d2)
	t.Logf("d1=%s data=%s\n", utils.StringInterface(d1, 3), string(d1.Data))
	t.Logf("d2=%s data=%s\n", utils.StringInterface(d2, 3), string(d2.Data))
}

func TestMediatedTransfer(t *testing.T) {
	bp := &BalanceProof{
		Nonce:             11,
		ChannelIdentifier: utils.Sha3([]byte("123")),
		TransferAmount:    big.NewInt(12),
		OpenBlockNumber:   3,
		Locksroot:         utils.EmptyHash,
	}
	lock := &mtree.Lock{
		Amount:         big.NewInt(34),
		Expiration:     4589895, //expiration block number
		LockSecretHash: utils.ShaSecret([]byte("hashlock")),
	}
	m1 := NewMediatedTransfer(bp, lock, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(33), []common.Address{utils.NewRandomAddress()})
	m1.Sign(GetTestPrivKey(), m1)
	data := m1.Pack()
	m2 := new(MediatedTransfer)
	m2.UnPack(data)
	spew.Dump("m1", m1)
	spew.Dump("m2", m2)
	if !reflect.DeepEqual(m1, m2) {
		t.Error("not equal")
	}
}

func TestNewAnnounceDisposedTransfer(t *testing.T) {
	bp := &AnnounceDisposedProof{
		ChannelIDInMessage: ChannelIDInMessage{
			ChannelIdentifier: utils.Sha3([]byte("123")),
			OpenBlockNumber:   3,
		},
		Lock: &mtree.Lock{
			Amount:         big.NewInt(34),
			Expiration:     4589895, //expiration block number
			LockSecretHash: utils.ShaSecret([]byte("hashlock")),
		},
	}
	m1 := NewAnnounceDisposed(bp, 1, "success")
	err := m1.Sign(GetTestPrivKey(), m1)
	if err != nil {
		t.Error(err)
		return
	}
	data := m1.Pack()
	m2 := new(AnnounceDisposed)
	err = m2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump("m1", m1)
	spew.Dump("m2", m2)
	if !reflect.DeepEqual(m1, m2) {
		t.Error("not equal")
	}
}

func TestNewSecret(t *testing.T) {
	bp := &BalanceProof{
		Nonce:             11,
		ChannelIdentifier: utils.Sha3([]byte("123")),
		TransferAmount:    big.NewInt(12),
		OpenBlockNumber:   3,
		Locksroot:         utils.EmptyHash,
	}
	s1 := NewUnlock(bp, utils.ShaSecret([]byte("xxx")))
	s1.Sign(GetTestPrivKey(), s1)
	data := s1.Pack()
	s2 := new(UnLock)
	err := s2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Error("not equal")
	}
}

func TestNewRevealSecret(t *testing.T) {
	s1 := NewRevealSecret(utils.ShaSecret([]byte("xxx")))
	s1.Data = []byte("123")
	s1.Sign(GetTestPrivKey(), s1)
	data := s1.Pack()
	s2 := new(RevealSecret)
	err := s2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Error("not equal")
	}
}
func TestErrorNotify(t *testing.T) {
	p1key, _ := utils.MakePrivateKeyAddress()

	m := NewErrorNotify(InvalidNonceErrorNotify, []byte{1, 2, 3})
	err := m.Sign(p1key, m)
	if err != nil {
		t.Error(err)
		return
	}
	data := m.Pack()
	//t.Logf("data=\n%s", hex.Dump(data))
	m2 := new(ErrorNotify)
	err = m2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, m, m2)
}
func TestNewSecretRequest(t *testing.T) {
	s1 := NewSecretRequest(utils.ShaSecret([]byte("xxx")), big.NewInt(506))
	s1.Sign(GetTestPrivKey(), s1)
	data := s1.Pack()
	s2 := new(SecretRequest)
	err := s2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Error("not equal")
	}
}
func TestNewRemoveExpiredHashlockTransfer(t *testing.T) {
	bp := &BalanceProof{
		Nonce:             11,
		ChannelIdentifier: utils.Sha3([]byte("123")),
		TransferAmount:    big.NewInt(12),
		OpenBlockNumber:   3,
		Locksroot:         utils.EmptyHash,
	}
	s1 := NewRemoveExpiredHashlockTransfer(bp, utils.ShaSecret([]byte("xxx")))
	s1.Sign(GetTestPrivKey(), s1)
	data := s1.Pack()
	s2 := new(RemoveExpiredHashlockTransfer)
	err := s2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Error("not equal")
	}
}
func TestLock_AsBytes(t *testing.T) {
	lock := &mtree.Lock{
		Amount:         big.NewInt(34),
		Expiration:     4589895, //expiration block number
		LockSecretHash: utils.ShaSecret([]byte("hashlock")),
	}
	t.Log("\n", hex.Dump(lock.AsBytes()))
	lock2 := new(mtree.Lock)
	lock2.FromBytes(lock.AsBytes())
	if !reflect.DeepEqual(lock, lock2) {
		t.Error("not equal")
	}
	//T.Log(lock.AsBytes())
}

func TestNewAnnounceDisposedTransferResponse(t *testing.T) {
	bp := &BalanceProof{
		Nonce:             11,
		ChannelIdentifier: utils.Sha3([]byte("123")),
		TransferAmount:    big.NewInt(12),
		OpenBlockNumber:   3,
		Locksroot:         utils.NewRandomHash(),
	}
	m := NewAnnounceDisposedResponse(bp, utils.NewRandomHash())
	err := m.Sign(GetTestPrivKey(), m)
	if err != nil {
		t.Error(err)
		return
	}
	data := m.Pack()
	m2 := new(AnnounceDisposedResponse)
	err = m2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(m, m2) {
		t.Error("not equal")
	}
}

func TestWithdrawRequest(t *testing.T) {
	p1key, p1addr := utils.MakePrivateKeyAddress()
	_, p2addr := utils.MakePrivateKeyAddress()
	bp := new(WithdrawRequestData)
	bp.ChannelIdentifier = utils.NewRandomHash()
	bp.OpenBlockNumber = 3
	bp.Participant1 = p1addr
	bp.Participant1Balance = big.NewInt(10)
	bp.Participant1Withdraw = big.NewInt(3)
	bp.Participant2 = p2addr
	m := NewWithdrawRequest(bp)
	err := m.Sign(p1key, m)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("m=%s\n", utils.StringInterface(m, 3))
	data := m.Pack()
	m2 := new(WithdrawRequest)
	err = m2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, m, m2)
}

func TestWithdrawResponse(t *testing.T) {
	_, p1addr := utils.MakePrivateKeyAddress()
	p2key, p2addr := utils.MakePrivateKeyAddress()
	bp := new(WithdrawReponseData)
	bp.ChannelIdentifier = utils.NewRandomHash()
	bp.ChannelIdentifier = utils.NewRandomHash()
	bp.OpenBlockNumber = 3
	bp.Participant1 = p1addr
	bp.Participant1Balance = big.NewInt(10)
	bp.Participant1Withdraw = big.NewInt(3)
	bp.Participant2 = p2addr

	fmt.Printf("addr1=%s,addr2=%s\n", utils.APex2(p1addr), utils.APex2(p2addr))
	m := NewWithdrawResponse(bp, 1, "testxxxxx")
	err := m.Sign(p2key, m)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("m=%s\n", utils.StringInterface(m, 3))
	data := m.Pack()
	m2 := new(WithdrawResponse)
	err = m2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, m, m2)
}
func TestSettleRequest(t *testing.T) {
	p1key, p1addr := utils.MakePrivateKeyAddress()
	_, p2addr := utils.MakePrivateKeyAddress()
	bp := new(SettleRequestData)
	bp.ChannelIdentifier = utils.NewRandomHash()
	bp.OpenBlockNumber = 3
	bp.Participant1 = p1addr
	bp.Participant1Balance = big.NewInt(10)
	bp.Participant2 = p2addr
	bp.Participant2Balance = big.NewInt(30)
	fmt.Printf("addr1=%s,addr2=%s\n", utils.APex2(p1addr), utils.APex2(p2addr))
	m := NewSettleRequest(bp)
	err := m.Sign(p1key, m)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("m=%s\n", utils.StringInterface(m, 3))
	data := m.Pack()
	m2 := new(SettleRequest)
	err = m2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, m, m2)
}
func TestSettleResponse(t *testing.T) {
	_, p1addr := utils.MakePrivateKeyAddress()
	p2key, p2addr := utils.MakePrivateKeyAddress()
	bp := new(SettleResponseData)
	bp.ChannelIdentifier = utils.NewRandomHash()
	bp.OpenBlockNumber = 3
	bp.Participant1 = p1addr
	bp.Participant1Balance = big.NewInt(10)
	bp.Participant2 = p2addr
	bp.Participant2Balance = big.NewInt(30)
	fmt.Printf("addr1=%s,addr2=%s\n", utils.APex2(p1addr), utils.APex2(p2addr))
	m := NewSettleResponse(bp, 1, "test1111111111111")
	err := m.Sign(p2key, m)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("m=%s\n", utils.StringInterface(m, 3))
	data := m.Pack()
	m2 := new(SettleResponse)
	err = m2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, m, m2)
}

type testStruct struct {
	T  int
	Bt *big.Int
}

func TestBigInt(t *testing.T) {
	tt := &testStruct{
		T: 3,
	}
	tt.Bt = new(big.Int)
	tt.Bt.SetString("0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa", 0)
	t.Log(tt.Bt.String())
	data, _ := json.Marshal(tt)
	t.Log(string(data))
	var tt2 testStruct
	json.Unmarshal(data, &tt2)
	if !reflect.DeepEqual(tt.Bt.Bytes(), tt2.Bt.Bytes()) {
		t.Error("not equal")
	}
}
