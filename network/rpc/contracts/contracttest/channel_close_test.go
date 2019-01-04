package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/Photon/utils"
)

// TestChannelCloseRight : 正确调用测试
// TestChannelCloseRight : Test for correct function call.
func TestChannelCloseRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(10)
	depositA2 := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 10
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// cases
	// close right
	bpA2 := createPartnerBalanceProof(a1, a2, big.NewInt(1), utils.EmptyHash, utils.EmptyHash, 1)
	tx, err := env.TokenNetwork.PrepareSettle(a1.Auth, env.TokenAddress, a2.Address, bpA2.TransferAmount, bpA2.LocksRoot, bpA2.Nonce, bpA2.AdditionalHash, bpA2.Signature)
	assertTxSuccess(t, &count, tx, err)
	// check the state
	_, _, _, state, _, _ := getChannelInfo(a1, a2)
	assertEqual(t, &count, ChannelStateClosed, state)
	balanceA1, _, nonceA1, err := env.TokenNetwork.GetChannelParticipantInfo(nil, env.TokenAddress, a1.Address, a2.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, &count, depositA1, balanceA1)
	assertEqual(t, &count, 0, nonceA1)
	// close twice
	tx, err = env.TokenNetwork.PrepareSettle(a1.Auth, env.TokenAddress, a2.Address, bpA2.TransferAmount, bpA2.LocksRoot, bpA2.Nonce, bpA2.AdditionalHash, bpA2.Signature)
	assertTxFail(t, &count, tx, err)

	// settle for other cases
	waitToSettle(a1, a2)
	tx, err = env.TokenNetwork.Settle(a1.Auth, env.TokenAddress, a1.Address, big.NewInt(0), utils.EmptyHash, a2.Address, bpA2.TransferAmount, bpA2.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	t.Log(endMsg("ChannelClose 正确调用测试", count, a1, a2))
}

// TestChannelCloseException : 异常调用测试
// TestChannelCloseException : Test for abnormal function call.
func TestChannelCloseException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	a3 := env.getRandomAccountExcept(t, a1, a2)
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(10)
	depositA2 := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 10
	// cases
	// close nonexistent channel
	tx, err := env.TokenNetwork.PrepareSettle(a1.Auth, env.TokenAddress, a2.Address, big.NewInt(1), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)

	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	bp := createPartnerBalanceProof(a1, a2, big.NewInt(1), utils.EmptyHash, utils.EmptyHash, 1)
	// close with wrong sender
	tx, err = env.TokenNetwork.PrepareSettle(a3.Auth, env.TokenAddress, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)
	// close with wrong signature
	bp.sign(a3.Key)
	tx, err = env.TokenNetwork.PrepareSettle(a2.Auth, env.TokenAddress, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)
	// close settled channel
	cooperativeSettleChannelIfExists(a1, a2)
	tx, err = env.TokenNetwork.PrepareSettle(a1.Auth, env.TokenAddress, a2.Address, big.NewInt(1), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelClose 异常调用测试", count, a1, a2, a3))

}

// TestChannelCloseEdge : 边界测试
// TestChannelCloseEdge : Edge Test.
func TestChannelCloseEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelClose 边界测试", count))
}

// TestChannelCloseAttack : 恶意调用测试
// TestChannelCloseAttack : Test for Potential Attack.
func TestChannelCloseAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(10)
	trasAmtA1 := big.NewInt(0)
	locksrootA1 := utils.EmptyHash

	depositA2 := big.NewInt(20)
	trasAmtA2 := big.NewInt(1)
	nonceA2 := uint64(1)
	locksrootA2 := utils.EmptyHash

	testSettleTimeout := TestSettleTimeoutMin + 10
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// close with self balance proof
	bpA2 := createPartnerBalanceProof(a1, a2, trasAmtA2, locksrootA2, locksrootA2, nonceA2)
	tx, err := env.TokenNetwork.PrepareSettle(a2.Auth, env.TokenAddress, a1.Address, bpA2.TransferAmount, bpA2.LocksRoot, bpA2.Nonce, bpA2.AdditionalHash, bpA2.Signature)
	assertTxFail(t, &count, tx, err)

	// close new channel with old balance proof
	tx, err = env.TokenNetwork.PrepareSettle(a1.Auth, env.TokenAddress, a2.Address, bpA2.TransferAmount, bpA2.LocksRoot, bpA2.Nonce, bpA2.AdditionalHash, bpA2.Signature)
	assertTxSuccess(t, nil, tx, err) // close
	waitToSettle(a1, a2)
	tx, err = env.TokenNetwork.Settle(a1.Auth, env.TokenAddress, a1.Address, trasAmtA1, locksrootA1, a2.Address, trasAmtA2, locksrootA2)
	assertTxSuccess(t, nil, tx, err)                                       // settle
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout) // reopen
	tx, err = env.TokenNetwork.PrepareSettle(a1.Auth, env.TokenAddress, a2.Address, bpA2.TransferAmount, bpA2.LocksRoot, bpA2.Nonce, bpA2.AdditionalHash, bpA2.Signature)
	assertTxFail(t, &count, tx, err) // close with old balance proof

	t.Log(endMsg("ChannelClose 恶意调用测试", count))
}
