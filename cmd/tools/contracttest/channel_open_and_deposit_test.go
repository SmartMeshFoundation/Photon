package contracttest

import (
	"math/big"
	"testing"

	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// TestChannelOpenAndDepositRight : 正确调用测试
// TestChannelOpenAndDepositRight : normal function call
func TestChannelOpenAndDepositRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	settleTimeout := TestSettleTimeoutMin + 10
	depositAmountA1 := big.NewInt(20)
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)
	// get token balance
	balanceA1 := getTokenBalance(a1)
	// run right
	tx, err := env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, depositAmountA1)
	assertTxSuccess(t, &count, tx, err)

	// get token balance new
	balanceA1New := getTokenBalance(a1)
	// check token balance
	assertEqual(t, &count, balanceA1.Sub(balanceA1, depositAmountA1), balanceA1New)
	// 查询通道
	// check channel
	_, settleBlockNumber, _, state, settleTimeout, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, &count, ChannelStateOpened, state)
	assertEqual(t, nil, uint64(0), settleBlockNumber)
	assertEqual(t, nil, settleTimeout, settleTimeout)
	// 查询通道双方信息
	// check channel participants info
	deposit, balanceHash, nonce, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, nil, depositAmountA1, deposit)
	assertEqual(t, nil, uint64(0), nonce)
	assertEqual(t, nil, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))

	deposit, balanceHash, nonce, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, nil, int64(0), deposit.Int64())
	assertEqual(t, nil, uint64(0), nonce)
	assertEqual(t, nil, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))
	t.Log(endMsg("ChannelOpenAndDeposit 正确调用测试", count))
}

// TestChannelOpenAndDepositException : 异常调用测试
// TestChannelOpenAndDepositException : abnormal function call
func TestChannelOpenAndDepositException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	settleTimeout := TestSettleTimeoutMin + 10
	depositAmountA1 := big.NewInt(20)
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)
	openChannelAndDeposit(a1, a2, depositAmountA1, big.NewInt(0), settleTimeout)
	// run when channel open, MUST FAIL
	tx, err := env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, depositAmountA1)
	assertTxFail(t, &count, tx, err)

	// run when channel close, MUST FAIL
	tx, err = env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, depositAmountA1)
	assertTxFail(t, &count, tx, err)

	// run when channel settled, MUST SUCCESS
	waitToSettle(a1, a2)
	tx, err = env.TokenNetwork.SettleChannel(a1.Auth, a1.Address, big.NewInt(0), utils.EmptyHash, a2.Address, big.NewInt(0), utils.EmptyHash)
	assertTxSuccess(t, nil, tx, err)
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, depositAmountA1)
	assertTxSuccess(t, &count, tx, err)
	t.Log(endMsg("ChannelOpenAndDeposit 异常调用测试", count))

}

// TestChannelOpenAndDepositEdge : 边界测试
// TestChannelOpenAndDepositEdge : edge test
func TestChannelOpenAndDepositEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	settleTimeout := TestSettleTimeoutMin + 10
	depositAmountA1 := big.NewInt(20)
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)

	//run with EmptyAddress, MUST FAIL
	tx, err := env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, EmptyAccountAddress, a2.Address, settleTimeout, depositAmountA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, EmptyAccountAddress, settleTimeout, depositAmountA1)
	assertTxFail(t, &count, tx, err)

	// run with wrong settleTimeout, MUST FAIL
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, 0, depositAmountA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMin-1, depositAmountA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMax+1, depositAmountA1)
	assertTxFail(t, &count, tx, err)

	// run with wrong deposit amount, MUST FAIL
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, big.NewInt(-1))
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, utils.MaxBigUInt256)
	assertTxFail(t, &count, tx, err)

	// run with deposit amount 0, MUST FAIL
	tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, big.NewInt(0))
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelOpenAndDeposit 边界测试", count))
}

// TestChannelOpenAndDepositAttack : 恶意调用测试
// TestChannelOpenAndDepositAttack : test for potential attack
func TestChannelOpenAndDepositAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelOpenAndDeposit 恶意调用测试", count))
}
