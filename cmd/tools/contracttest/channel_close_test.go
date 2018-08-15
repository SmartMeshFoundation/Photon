package contracttest

import (
	"math/big"
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// TestChannelCloseRight : 正确调用测试
func TestChannelCloseRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	a1, a2 := env.getTwoRandomAccount(t)
	depositA1 := big.NewInt(10)
	depositA2 := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 10
	cooperativeSettleChannelIfExists(a1, a2)
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// cases
	// close right
	bp := createPartnerBalanceProof(a1, a2, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)
	tx, err := env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxSuccess(t, &count, tx, err)
	// close twice
	tx, err = env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)
	t.Log(endMsg("ChannelClose 正确调用测试", count, a1, a2))
}

// TestChannelCloseException : 异常调用测试
func TestChannelCloseException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2, a3 := env.getThreeRandomAccount(t)
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(10)
	depositA2 := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 10
	// cases
	// close nonexistent channel
	tx, err := env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(1), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)

	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	bp := createPartnerBalanceProof(a1, a2, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)
	// close with wrong sender
	tx, err = env.TokenNetwork.CloseChannel(a3.Auth, a2.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	assertTxFail(t, &count, tx, err)
	// close with wrong signature
	bp.sign(a3.Key)
	tx, err = env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	_, _, _, state, _, _ := getChannelInfo(a1, a2)
	fmt.Printf("state=========%d\n", state)
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
	t.Log(endMsg("ChannelClose 恶意调用测试", count))
}
