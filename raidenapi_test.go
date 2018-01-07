package raiden_network

import (
	"testing"

	"github.com/SmartMeshFoundation/raiden-network/utils"
)

func TestSwapKeyAsMapKey(t *testing.T) {
	key1 := SwapKey{
		Identifier: 32,
		FromToken:  utils.NewRandomAddress(),
		FromAmount: 300,
	}
	key2 := key1
	m := make(map[SwapKey]bool)
	m[key1] = true
	if m[key2] != true {
		t.Error("expect equal")
	}
	key2.Identifier = 3
	if m[key2] == true {
		t.Error("should not equal")
	}
}
