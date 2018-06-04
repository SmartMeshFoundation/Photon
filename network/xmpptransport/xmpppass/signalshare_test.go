package xmpppass

import (
	"testing"

	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestCreatePasswordAndVerify(t *testing.T) {
	key, _ := crypto.GenerateKey()
	sig, err := CreatePassword(key)
	if err != nil {
		t.Error(err)
		return
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	fmt.Printf("addr=%s,sig=%s\n", addr.String(), sig)
	err = VerifySignature(addr.String(), sig)
	if err != nil {
		t.Error(err)
	}
}
