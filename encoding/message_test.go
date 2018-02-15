package encoding

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"bytes"

	"errors"

	"reflect"

	"math/big"

	"encoding/json"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const TestPrivkey = "4359f525e2b373089be5fe8f9a4e8ffb6d30e2960918be426217921e1b2547f7"

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
    ping.sign(PRIVKEY, ADDRESS)
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
	tokenaddress := utils.NewRandomAddress()
	channel := utils.NewRandomAddress()
	p := NewDirectTransfer(32, 11, tokenaddress, channel,
		big.NewInt(11), utils.EmptyAddress,
		utils.EmptyHash)
	var sm SignedMessager = p
	err := p.Sign(GetTestPrivKey(), p)
	if err != nil {
		t.Error(err)
	}
	data := p.Pack()
	err = sm.VerifySignature(data)
	if err != nil {
		t.Error(err)
	}
	p2 := new(DirectTransfer)
	err = p2.UnPack(data)
	if err != nil {
		t.Error(err)
		return
	}
	if p2.Channel != p.Channel || p.Token != p2.Token ||
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
	d1 := NewDirectTransfer(22, 32, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(12), utils.NewRandomAddress(), utils.Sha3([]byte("abd")))
	d1.Sign(GetTestPrivKey(), d1)
	d2 := new(DirectTransfer)
	err := d2.UnPack(d1.Pack())
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(d1, d2) {
		t.Error("not equal")
	}
	//T.Log(utils.StringInterface(d1, 3))
	//if utils.StringInterface(d1, 3) != utils.StringInterface(d2, 3) {
	//
	//}
}

func TestMediatedTransfer(t *testing.T) {
	lock := &Lock{
		Amount:     big.NewInt(34),
		Expiration: 4589895, //expiration block number
		HashLock:   utils.Sha3([]byte("hashlock")),
	}
	m1 := NewMediatedTransfer(11, 32, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(33), utils.NewRandomAddress(),
		utils.Sha3([]byte("ddd")), lock, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(33))
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

func TestNewRefundTransfer(t *testing.T) {
	lock := &Lock{
		Amount:     big.NewInt(34),
		Expiration: 4589895, //expiration block number
		HashLock:   utils.Sha3([]byte("hashlock")),
	}
	m1 := NewRefundTransfer(11, 32, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(33), utils.NewRandomAddress(),
		utils.Sha3([]byte("ddd")), lock, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(33))
	m1.Sign(GetTestPrivKey(), m1)
	data := m1.Pack()
	m2 := new(RefundTransfer)
	m2.UnPack(data)
	spew.Dump("m1", m1)
	spew.Dump("m2", m2)
	if !reflect.DeepEqual(m1, m2) {
		t.Error("not equal")
	}
}

func TestNewSecret(t *testing.T) {
	s1 := NewSecret(30, 40, utils.NewRandomAddress(), big.NewInt(50), utils.Sha3([]byte("oo")), utils.Sha3([]byte("xxx")))
	s1.Sign(GetTestPrivKey(), s1)
	data := s1.Pack()
	s2 := new(Secret)
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
	s1 := NewRevealSecret(utils.Sha3([]byte("xxx")))
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

func TestNewSecretRequest(t *testing.T) {
	s1 := NewSecretRequest(606, utils.Sha3([]byte("xxx")), big.NewInt(506))
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

func TestLock_AsBytes(t *testing.T) {
	lock := &Lock{
		Amount:     big.NewInt(34),
		Expiration: 4589895, //expiration block number
		HashLock:   utils.Sha3([]byte("hashlock")),
	}
	t.Log("\n", hex.Dump(lock.AsBytes()))
	lock2 := new(Lock)
	lock2.FromBytes(lock.AsBytes())
	if !reflect.DeepEqual(lock, lock2) {
		t.Error("not equal")
	}
	//T.Log(lock.AsBytes())
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
