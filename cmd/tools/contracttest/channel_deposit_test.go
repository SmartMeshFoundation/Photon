package contracttest

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestChannelDepositRight : 正确调用测试
func TestChannelDepositRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	a1, a2, a3 := env.getThreeRandomAccount(t)
	cooperativeSettleChannelIfExists(a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 10
	depositA1 := big.NewInt(200)
	depositA2 := big.NewInt(300)
	// 创建
	tx, err := env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assertTxSuccess(t, nil, tx, err)
	balanceA1, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, 0, big.NewInt(0).Cmp(balanceA1))
	// cases
	// deposit a1
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, a2.Address, depositA1)
	assertTxSuccess(t, &count, tx, err)
	// check a1 balance
	balanceA1, _, _, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, depositA1, balanceA1)
	// third deposit a1
	tx, err = env.TokenNetwork.Deposit(a3.Auth, a1.Address, a2.Address, depositA1)
	assertTxSuccess(t, &count, tx, err)
	//  check a1 balance
	balanceA1, _, _, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, depositA1.Add(depositA1, depositA1), balanceA1)

	// check a2 balance
	balanceA2, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, 0, big.NewInt(0).Cmp(balanceA2))
	// deposit a2
	tx, err = env.TokenNetwork.Deposit(a2.Auth, a2.Address, a1.Address, depositA2)
	assertTxSuccess(t, &count, tx, err)
	// check a2 balance
	balanceA2, _, _, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, nil, depositA2, balanceA2)
	t.Log(endMsg("ChannelDeposit 正确调用测试", count, a1, a2, a3))
}

// TestChannelDepositException : 异常调用测试
func TestChannelDepositException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelDeposit 异常调用测试", count))

}

// TestChannelDepositEdge : 边界测试
func TestChannelDepositEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	a1, a2 := env.getTwoRandomAccount(t)
	cooperativeSettleChannelIfExists(a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 10
	depositA1 := big.NewInt(200)
	// 创建
	tx, err := env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assertTxSuccess(t, nil, tx, err)
	// cases
	// self地址错误
	tx, err = env.TokenNetwork.Deposit(a1.Auth, common.HexToAddress("-1"), a2.Address, depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, common.HexToAddress(""), a2.Address, depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, common.HexToAddress("0x0"), a2.Address, depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, EmptyAccountAddress, a2.Address, depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, FakeAccountAddress, a2.Address, depositA1)
	assertTxFail(t, &count, tx, err)
	// partner地址错误
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, common.HexToAddress("-1"), depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, common.HexToAddress(""), depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, common.HexToAddress("0x0"), depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, EmptyAccountAddress, depositA1)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, FakeAccountAddress, depositA1)
	assertTxFail(t, &count, tx, err)
	// deposit 0
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, a2.Address, big.NewInt(0))
	assertTxFail(t, &count, tx, err)
	// deposit > approve
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, a2.Address, big.NewInt(50000001))
	assertTxFail(t, &count, tx, err)
	// deposit > balance
	approve(a1, big.NewInt(500000000002))
	tx, err = env.TokenNetwork.Deposit(a1.Auth, a1.Address, a2.Address, big.NewInt(500000000001))
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelDeposit 边界测试", count, a1, a2))
}

// TestChannelDepositAttack : 恶意调用测试
func TestChannelDepositAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelDeposit 恶意调用测试", count))
}
