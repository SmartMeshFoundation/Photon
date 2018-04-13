package transfer

import (
	"testing"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestNewBalanceProofStateFromEnvelopMessage(t *testing.T) {
	var id uint64 = 32
	var nonce int64 = 78
	var transferAmount = big.NewInt(9999)
	secret := encoding.NewSecret(id, nonce, utils.EmptyAddress, transferAmount, utils.EmptyHash, utils.EmptyHash)
	state := NewBalanceProofStateFromEnvelopMessage(secret)
	if state.Nonce != nonce || state.TransferAmount.Cmp(transferAmount) != 0 || state.LocksRoot != utils.EmptyHash {
		t.Error("NewBalanceProofStateFromEnvelopMessage error")
	}
}
