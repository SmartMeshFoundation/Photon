package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

// TODO start from withdraw
// TestCooperativeSettleRight : 正确调用测试
func TestCooperativeSettleRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2, a3 := env.getThreeRandomAccount(t)
	// cases
	// 正常调用
	depositA1 := big.NewInt(20)
	depositA2 := big.NewInt(10)
	balanceA1 := big.NewInt(4)
	balanceA2 := big.NewInt(26)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	cs := getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	tx, err := env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxSuccess(t, &count, tx, err)

	// 一方金额为0
	depositA1 = big.NewInt(20)
	depositA2 = big.NewInt(10)
	balanceA1 = big.NewInt(0)
	balanceA2 = big.NewInt(30)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	cs = getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxSuccess(t, &count, tx, err)

	// 所有金额为0
	depositA1 = big.NewInt(0)
	depositA2 = big.NewInt(0)
	balanceA1 = big.NewInt(0)
	balanceA2 = big.NewInt(0)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	cs = getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxSuccess(t, &count, tx, err)

	// 双方先withdraw然后CooperativeSettle
	depositA1 = big.NewInt(20)
	depositA2 = big.NewInt(10)
	withdrawA1 := big.NewInt(10)
	withdrawA2 := big.NewInt(5)
	balanceA1 = big.NewInt(2)
	balanceA2 = big.NewInt(13)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	withdraw(a1, withdrawA1, depositA1, a2, withdrawA2, depositA2)
	cs = getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxSuccess(t, &count, tx, err)
	t.Log(endMsg("CooperativeSettle 正确调用测试", count, a1, a2, a3))
}

// TestCooperativeSettleException : 异常调用测试
func TestCooperativeSettleException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2, a3 := env.getThreeRandomAccount(t)
	depositA1 := big.NewInt(20)
	depositA2 := big.NewInt(10)
	balanceA1 := big.NewInt(4)
	balanceA2 := big.NewInt(26)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	cs := getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	// cases
	// 重复CooperativeSettle
	tx, err := env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxSuccess(t, nil, tx, err)
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)

	// close状态下CooperativeSettle
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	// get new cs
	cs = getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	// a1 close channel
	tx, err = env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// close状态下CooperativeSettle
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("CooperativeSettle 异常调用测试", count, a1, a2, a3))
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
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
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
	t.Log(endMsg("CooperativeSettle 边界测试", count, a1, a2, a3))
}

// TestCooperativeSettleAttack : 恶意调用测试
func TestCooperativeSettleAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2, a3 := env.getThreeRandomAccount(t)
	depositA1 := big.NewInt(20)
	depositA2 := big.NewInt(10)
	balanceA1 := big.NewInt(4)
	balanceA2 := big.NewInt(26)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	cs := getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	// cases
	// 第三方签名
	tx, err := env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a3.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a3.Key))
	assertTxFail(t, &count, tx, err)
	// 错误余额
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant2Balance, a2.Address, cs.Participant1Balance, cs.sign(a3.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)
	wrongBalance := big.NewInt(cs.Participant1Balance.Int64() + 1)
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, wrongBalance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, wrongBalance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)
	wrongBalance = big.NewInt(cs.Participant1Balance.Int64() - 1)
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, wrongBalance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, wrongBalance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)
	// CooperativeSettle后重新创建channel,使用老的余额证明来CooperativeSettle新的channel
	depositA1 = big.NewInt(20)
	depositA2 = big.NewInt(10)
	balanceA1 = big.NewInt(4)
	balanceA2 = big.NewInt(26)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	cs = getCooperativeSettleParams(a1, a2, balanceA1, balanceA2)
	cs.Participant1Balance = balanceA1
	cs.Participant2Balance = balanceA2
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxSuccess(t, nil, tx, err)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, TestSettleTimeoutMin+10)
	tx, err = env.TokenNetwork.CooperativeSettle(
		a3.Auth, a1.Address, cs.Participant1Balance, a2.Address, cs.Participant2Balance, cs.sign(a1.Key), cs.sign(a2.Key))
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("CooperativeSettle 恶意调用测试", count, a1, a2, a3))
}
