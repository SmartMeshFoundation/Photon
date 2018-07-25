package contracttest

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestCooperativeSettleRight : 正确调用测试
func TestCooperativeSettleRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2, a3 := env.getThreeRandomAccount(t)
	depositA1 := big.NewInt(20)
	depositA2 := big.NewInt(10)
	balanceA1 := big.NewInt(4)
	balanceA2 := big.NewInt(26)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin + 10)
	cs := getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	// cases
	tx, err := env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, balanceA1, a2.Address, balanceA2, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxSuccess(t, &count, tx, err)
	t.Logf("CooperativeSettle 正确调用测试完成,case数量 : %d", count)
}

// TestCooperativeSettleException : 异常调用测试
func TestCooperativeSettleException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Logf("CooperativeSettle 异常调用测试完成,case数量 : %d", count)

}

// TestCooperativeSettleEdge : 边界测试
func TestCooperativeSettleEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2, a3 := env.getThreeRandomAccount(t)
	depositA1 := big.NewInt(20)
	depositA2 := big.NewInt(10)
	balanceA1 := big.NewInt(4)
	balanceA2 := big.NewInt(26)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin + 10)
	cs := getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	signA1 := cs.sign(a1.Key)
	signA2 := cs.sign(a2.Key)
	// cases
	// edge balance
	// partner1_balance = -1
	tx, err := env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, big.NewInt(-1), a2.Address, balanceA2, signA1, signA2)
	assertTxFail(t, &count, tx, err)
	// partner2_balance = -1
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, a2.Address, big.NewInt(-1), signA1, signA2)
	assertTxFail(t, &count, tx, err)

	// edge address
	// partner1_address = "0x0"
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, common.HexToAddress("0x0"), balanceA1, a2.Address, balanceA2, signA1, signA2)
	assertTxFail(t, &count, tx, err)
	// partner2_address = "0x0"
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, common.HexToAddress("0x0"), balanceA2, signA1, signA2)
	assertTxFail(t, &count, tx, err)
	// partner1_address = "0x0000000000000000000000000000000000000000"
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, EmptyAccountAddress, balanceA1, a2.Address, balanceA2, signA1, signA2)
	assertTxFail(t, &count, tx, err)
	// partner2_address = "0x0000000000000000000000000000000000000000"
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, EmptyAccountAddress, balanceA2, signA1, signA2)
	assertTxFail(t, &count, tx, err)

	// edge sign
	// partner1_sign = "0x0"
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, a2.Address, balanceA2, common.Hex2Bytes("0x0"), signA2)
	assertTxFail(t, &count, tx, err)
	// partner2_sign = "0x0"
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, a2.Address, balanceA2, signA1, common.Hex2Bytes("0x0"))
	assertTxFail(t, &count, tx, err)
	t.Logf("CooperativeSettle 边界测试完成,case数量 : %d", count)
}

// TestCooperativeSettleAttack : 恶意调用测试
func TestCooperativeSettleAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Logf("CooperativeSettle 恶意调用测试完成,case数量 : %d", count)
}