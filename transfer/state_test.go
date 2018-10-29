package transfer

import (
	"testing"

	"math/big"

	"bytes"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/utils"
)

func TestNewBalanceProofStateFromEnvelopMessage(t *testing.T) {
	bp := &encoding.BalanceProof{
		ChannelIdentifier: utils.NewRandomHash(),
		OpenBlockNumber:   3,
		Nonce:             30,
		TransferAmount:    big.NewInt(10),
		Locksroot:         utils.NewRandomHash(),
	}
	secret := encoding.NewUnlock(bp, utils.NewRandomHash())
	state := NewBalanceProofStateFromEnvelopMessage(secret)
	if state.Nonce != bp.Nonce || state.TransferAmount.Cmp(bp.TransferAmount) != 0 || bytes.Compare(state.LocksRoot[:], bp.Locksroot[:]) != 0 {
		t.Errorf("NewBalanceProofStateFromEnvelopMessage error,state=%s,bp=%s", utils.StringInterface(state, 3), utils.StringInterface(bp, 3))
	}
}
