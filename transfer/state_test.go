package transfer

import (
	"testing"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestNewBalanceProofStateFromEnvelopMessage(t *testing.T) {
	var nonce int64 = 78
	var transferAmount = big.NewInt(9999)
	bp := &encoding.BalanceProof{
		TokenAddress:      utils.NewRandomAddress(),
		ChannelIdentifier: utils.NewRandomHash(),
		OpenBlockNumber:   3,
		Nonce:             30,
		TransferAmount:    big.NewInt(10),
		Locksroot:         utils.NewRandomHash(),
	}
	secret := encoding.NewUnlock(bp, utils.NewRandomHash())
	state := NewBalanceProofStateFromEnvelopMessage(secret)
	if state.Nonce != nonce || state.TransferAmount.Cmp(transferAmount) != 0 || state.LocksRoot != utils.EmptyHash {
		t.Error("NewBalanceProofStateFromEnvelopMessage error")
	}
}
