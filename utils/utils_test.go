package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

func TestNewRandomAddress(t *testing.T) {
	addr := NewRandomAddress()
	fmt.Println(addr)
	fmt.Printf("addrs=%s\n", addr)
	fmt.Printf("addrs=%s", addr.String())
	//spew.Dump(addr)
	//t.Logf("addrs=%s\n", addr)
	//t.Logf("addrv=%v\n", addr)
}
