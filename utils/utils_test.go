package utils

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestBigIntTo32Bytes(t *testing.T) {
	s := "0000000000000000000000000000000000000000000000000000000000003322"
	expect, _ := hex.DecodeString(s)
	if !bytes.Equal(expect, BigIntTo32Bytes(big.NewInt(0x3322))) {
		t.Errorf("0x3322 to 32s should be  %s", s)
	}
}
