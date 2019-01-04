package contracttest

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestChannelDepositRight : 正确调用测试
// TestChannelDepositRight : normal function call
func TestChannelDepositRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	a3 := env.getRandomAccountExcept(t, a1, a2)
	cooperativeSettleChannelIfExists(a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 10
	depositA1 := big.NewInt(200)
	depositA2 := big.NewInt(300)
	// 创建
	// cases
	// deposit a1
	tx, err := env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, a2.Address, depositA1, testSettleTimeout)
	assertTxSuccess(t, &count, tx, err)
	// check a1 balance
	balanceA1, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, env.TokenAddress, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, depositA1, balanceA1)
	// third deposit a1
	tx, err = env.TokenNetwork.Deposit(a3.Auth, env.TokenAddress, a1.Address, a2.Address, depositA1, testSettleTimeout)
	assertTxSuccess(t, &count, tx, err)
	//  check a1 balance
	balanceA1, _, _, err = env.TokenNetwork.GetChannelParticipantInfo(nil, env.TokenAddress, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, depositA1.Add(depositA1, depositA1), balanceA1)

	// check a2 balance
	balanceA2, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, env.TokenAddress, a2.Address, a1.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, 0, big.NewInt(0).Cmp(balanceA2))
	// deposit a2
	tx, err = env.TokenNetwork.Deposit(a2.Auth, env.TokenAddress, a2.Address, a1.Address, depositA2, testSettleTimeout)
	assertTxSuccess(t, &count, tx, err)
	// check a2 balance
	balanceA2, _, _, err = env.TokenNetwork.GetChannelParticipantInfo(nil, env.TokenAddress, a2.Address, a1.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, depositA2, balanceA2)
	t.Log(endMsg("ChannelDeposit 正确调用测试", count, a1, a2, a3))
}

// TestChannelDepositException : 异常调用测试
// TestChannelDepositException : abnormal function call
func TestChannelDepositException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelDeposit 异常调用测试", count))

}

// TestChannelDepositEdge : 边界测试
// TestChannelDepositEdge : edge test
func TestChannelDepositEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 10
	depositA1 := big.NewInt(200)
	// 创建
	// cases
	// self地址错误
	// self address fault.
	tx, err := env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, common.HexToAddress("-1"), a2.Address, depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, common.HexToAddress(""), a2.Address, depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, common.HexToAddress("0x0"), a2.Address, depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, EmptyAccountAddress, a2.Address, depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// partner地址错误
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, common.HexToAddress("-1"), depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, common.HexToAddress(""), depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, common.HexToAddress("0x0"), depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, EmptyAccountAddress, depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// deposit 0
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, a2.Address, big.NewInt(0), testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// deposit > approve
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, a2.Address, big.NewInt(50000001), testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// deposit > balance
	balance, err := env.Token.BalanceOf(nil, a1.Address)
	assertSuccess(t, nil, err)
	depositA1 = balance.Add(balance, big.NewInt(10000))
	tx, err = env.TokenNetwork.Deposit(a1.Auth, env.TokenAddress, a1.Address, a2.Address, depositA1, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelDeposit 边界测试", count, a1, a2))
}

// TestChannelDepositAttack : 恶意调用测试
// TestChannelDepositAttack : test for potential attack.
func TestChannelDepositAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelDeposit 恶意调用测试", count))
}
