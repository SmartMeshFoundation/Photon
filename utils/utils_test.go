package utils

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"crypto/sha256"

	"github.com/ethereum/go-ethereum/common"
)

func TestBigIntTo32Bytes(t *testing.T) {
	s := "0000000000000000000000000000000000000000000000000000000000003322"
	expect, _ := hex.DecodeString(s)
	if !bytes.Equal(expect, BigIntTo32Bytes(big.NewInt(0x3322))) {
		t.Errorf("0x3322 to 32s should be  %s", s)
	}
}

func TestSha3(t *testing.T) {
	data := []byte("abc")
	hashsha3 := Sha3(data)
	hash256 := sha256.Sum256(data)
	hashsecret := ShaSecret(data)
	t.Logf("sha3=%s", hashsha3.String())
	t.Logf("sha256=%s", common.Bytes2Hex(hash256[:]))
	t.Logf("hashsecret=%s", common.Bytes2Hex(hashsecret[:]))
}

func TestShaSecret(t *testing.T) {
	cases := map[string]string{
		"b4e12862b4433c3eab351e07ccedd64cdd16273c057ca81ab2c4c7b93cdba2ce": "9aa5bea8995976cb7ee6adeeafa51445b23ef199a22785c12ca968cd13774d7f",
	}
	for k, v := range cases {
		secret, _ := hex.DecodeString(k)
		secrethash, _ := hex.DecodeString(v)
		calchash := ShaSecret(secret[:])
		if !bytes.Equal(secrethash[:], calchash[:]) {
			t.Error("not equal")
		}
	}
}
