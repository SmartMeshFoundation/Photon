package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// TestChannelCloseRight : 正确调用测试
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
	bp := createPartnerBalanceProof(a1, a2, big.NewInt(1), utils.EmptyHash, utils.EmptyHash, 1)
	tx, err := env.TokenNetwork.CloseChannel(a2.Auth, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxSuccess(t, &count, tx, err)
	// check the state
	_, _, _, state, _, _ := getChannelInfo(a1, a2)
	assertEqual(t, &count, ChannelStateClosed, state)
	balanceA1, _, nonceA1, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, &count, depositA1, balanceA1)
	assertEqual(t, &count, bp.Nonce, nonceA1)
	// close twice
	tx, err = env.TokenNetwork.CloseChannel(a2.Auth, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelClose 正确调用测试", count, a1, a2))
}

// TestChannelCloseException : 异常调用测试
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
	tx, err := env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(1), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)

	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	bp := createPartnerBalanceProof(a1, a2, big.NewInt(1), utils.EmptyHash, utils.EmptyHash, 1)
	// close with wrong sender
	tx, err = env.TokenNetwork.CloseChannel(a3.Auth, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)
	// close with wrong signature
	bp.sign(a3.Key)
	tx, err = env.TokenNetwork.CloseChannel(a2.Auth, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)
	// close settled channel
	cooperativeSettleChannelIfExists(a1, a2)
	tx, err = env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(1), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelClose 异常调用测试", count, a1, a2, a3))

}

// TestChannelCloseEdge : 边界测试
func TestChannelCloseEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelClose 边界测试", count))
}

// TestChannelCloseAttack : 恶意调用测试
func TestChannelCloseAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(10)
	trasAmtA1 := big.NewInt(1)
	nonceA1 := uint64(1)
	locksrootA1 := utils.EmptyHash

	depositA2 := big.NewInt(20)
	trasAmtA2 := big.NewInt(0)
	locksrootA2 := utils.EmptyHash

	testSettleTimeout := TestSettleTimeoutMin + 10
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// close with self balance proof
	bp := createPartnerBalanceProof(a1, a2, trasAmtA1, locksrootA1, locksrootA1, nonceA1)
	tx, err := env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)

	// close new channel with old balance proof
	tx, err = env.TokenNetwork.CloseChannel(a2.Auth, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxSuccess(t, nil, tx, err) // close
	waitForSettle(testSettleTimeout)
	tx, err = env.TokenNetwork.SettleChannel(a1.Auth, a1.Address, trasAmtA1, locksrootA1, a2.Address, trasAmtA2, locksrootA2)
	assertTxSuccess(t, nil, tx, err)                                       // settle
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout) // reopen
	tx, err = env.TokenNetwork.CloseChannel(a2.Auth, a1.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err) // close with old balance proof

	t.Log(endMsg("ChannelClose 恶意调用测试", count))
}
