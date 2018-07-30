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
	t.Log(endMsg("ChannelClose 正确调用测试", count))
}

// TestChannelCloseException : 异常调用测试
func TestChannelCloseException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoRandomAccount(t)
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(10)
	depositA2 := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 10
	// cases
	// close nonexistent channel
	tx, err := env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(1), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)
	// close settled channel
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	cooperativeSettleChannelIfExists(a1, a2)
	tx, err = env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(1), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelClose 异常调用测试", count))

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
	t.Log(endMsg("ChannelClose 恶意调用测试", count))
}
