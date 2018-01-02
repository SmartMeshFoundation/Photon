package transfer

import (
	"testing"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/utils"
)

func TestNewBalanceProofStateFromEnvelopMessage(t *testing.T) {
	var id uint64 = 32
	var nonce int64 = 78
	var transferAmount int64 = 9999
	secret := encoding.NewSecret(id, nonce, utils.EmptyAddress, transferAmount, utils.EmptyHash, utils.EmptyHash)
	state := NewBalanceProofStateFromEnvelopMessage(secret)
	if state.Nonce != nonce || state.TransferAmount != transferAmount || state.LocksRoot != utils.EmptyHash {
		t.Error("NewBalanceProofStateFromEnvelopMessage error")
	}
}
