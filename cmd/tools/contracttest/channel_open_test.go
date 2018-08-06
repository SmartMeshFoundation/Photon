package contracttest

import (
	"testing"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
)

// TestOpenChannelRight : 正确调用测试
func TestOpenChannelRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	a1, a2 := env.getTwoRandomAccount(t)
	cooperativeSettleChannelIfExists(a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 10
	count := 0
	// cases
	// 正确创建
	tx, err := env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assertTxSuccess(t, &count, tx, err)
	// 查询通道
	_, _, _, state, _, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, &count, ChannelStateOpened, state)
	// 查询通道双方信息
	deposit, balanceHash, nonce, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, nil, int64(0), deposit.Int64())
	assertEqual(t, nil, uint64(0), nonce)
	assertEqual(t, nil, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))

	deposit, balanceHash, nonce, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, nil, int64(0), deposit.Int64())
	assertEqual(t, nil, uint64(0), nonce)
	assertEqual(t, nil, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))
	t.Log(endMsg("OpenChannel 正确调用测试", count, a1, a2))
}

// TestOpenChannelRight : 异常调用测试
func TestOpenChannelException(t *testing.T) {
	InitEnv(t, "./env.INI")
	a1, a2 := env.getTwoRandomAccount(t)
	cooperativeSettleChannelIfExists(a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 10
	count := 0
	tx, err := env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assertTxSuccess(t, nil, tx, err)
	// cases
	// 重复创建
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a2.Address, a1.Address, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("OpenChannel 异常调用测试", count, a1, a2))
}

// TestOpenChannelEdge : 边界测试
func TestOpenChannelEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoRandomAccount(t)
	cooperativeSettleChannelIfExists(a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 5
	// cases
	// settles_timeout = 0
	tx, err := env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, 0)
	assertTxFail(t, &count, tx, err)
	// self地址为0x0
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, common.StringToAddress("0x0"), a2.Address, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// self地址为""
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, common.StringToAddress(""), a2.Address, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// self地址为0x03432
	tx, err = env.TokenNetwork.OpenChannel(a2.Auth, FakeAccountAddress, a2.Address, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// self地址为0x0000000000000000000000000000000000000000
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, EmptyAccountAddress, a2.Address, testSettleTimeout)
	assertTxFail(t, &count, tx, err)

	// partner地址为0x0
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, common.StringToAddress("0x0"), testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// partner地址为""
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, common.StringToAddress(""), testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// partner地址为0x03432
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, FakeAccountAddress, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// partner地址为0x0000000000000000000000000000000000000000
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, EmptyAccountAddress, testSettleTimeout)
	assertTxFail(t, &count, tx, err)

	// 通道双方地址相同
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a1.Address, testSettleTimeout)
	assertTxFail(t, &count, tx, err)
	// settle_timeout = 5
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMin-1)
	assertTxFail(t, &count, tx, err)
	// settle_timeout = 2700001
	tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMax+1)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("OpenChannel 边界测试", count, a1, a2))
}

// TestOpenChannelAttack : 恶意调用测试
func TestOpenChannelAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("OpenChannel 恶意调用测试", count))

}
