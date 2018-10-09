package contracttest

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

// TestTokenFallbackRight : 正确调用测试
func TestTokenFallbackRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	depositAmountA1 := big.NewInt(10)
	depositAmountA2 := big.NewInt(20)
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	cooperativeSettleChannelIfExists(a1, a2)
	tokenBalanceA1, tokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)

	// a1 openAndDeposit by fall back
	buf := new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(big.NewInt(1)))   // func code of openAndDeposit
	buf.Write(to32bytes(a1.Address[:]))               // participant
	buf.Write(to32bytes(a2.Address[:]))               // partner
	buf.Write(utils.BigIntTo32Bytes(big.NewInt(100))) //settle_timeout
	tx, err := env.Token.Transfer(a1.Auth, env.TokenNetworkAddress, depositAmountA1, buf.Bytes())
	assertTxSuccess(t, &count, tx, err)

	// a2 deposit by fall back
	buf = new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(big.NewInt(2))) // func code of deposit
	buf.Write(to32bytes(a2.Address[:]))             // participant
	buf.Write(to32bytes(a1.Address[:]))             // partner
	tx, err = env.Token.Transfer(a2.Auth, env.TokenNetworkAddress, depositAmountA2, buf.Bytes())
	assertTxSuccess(t, &count, tx, err)

	// a1 deposit 10 by approve and call
	depositAmount := big.NewInt(10)
	buf = new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(big.NewInt(2))) // func code of deposit
	buf.Write(to32bytes(a1.Address[:]))
	buf.Write(to32bytes(a2.Address[:]))
	tx, err = env.Token.ApproveAndCall(a1.Auth, env.TokenNetworkAddress, depositAmount, buf.Bytes())
	assertTxSuccess(t, &count, tx, err)
	depositAmountA1.Add(depositAmountA1, depositAmount)

	// a2 deposit 10 by approve and call
	buf = new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(big.NewInt(2))) // func code of deposit
	buf.Write(to32bytes(a2.Address[:]))
	buf.Write(to32bytes(a1.Address[:]))
	tx, err = env.Token.ApproveAndCall(a2.Auth, env.TokenNetworkAddress, depositAmount, buf.Bytes())
	assertTxSuccess(t, &count, tx, err)
	depositAmountA2.Add(depositAmountA2, depositAmount)

	// check state
	// 查询通道
	_, settleBlockNumber, _, state, settleTimeout, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, &count, ChannelStateOpened, state)
	assertEqual(t, nil, 0, settleBlockNumber)
	assertEqual(t, nil, settleTimeout, settleTimeout)
	// 查询通道双方信息
	deposit, balanceHash, nonce, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, nil, depositAmountA1, deposit)
	assertEqual(t, nil, 0, nonce)
	assertEqual(t, nil, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))

	deposit, balanceHash, nonce, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assertSuccess(t, &count, err)
	assertEqual(t, nil, depositAmountA2, deposit)
	assertEqual(t, nil, 0, nonce)
	assertEqual(t, nil, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))

	// check a1,a2's token balance
	tokenBalanceA1New, tokenBalanceA2New := getTokenBalance(a1), getTokenBalance(a2)
	assertEqual(t, &count, tokenBalanceA1.Sub(tokenBalanceA1, depositAmountA1), tokenBalanceA1New)
	assertEqual(t, &count, tokenBalanceA2.Sub(tokenBalanceA2, depositAmountA2), tokenBalanceA2New)
	t.Log(endMsg("TokenFallback 正确调用测试", count))
}

// TestTokenFallbackException : 异常调用测试
func TestTokenFallbackException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("TokenFallback 异常调用测试", count))

}

// TestTokenFallbackEdge : 边界测试
func TestTokenFallbackEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("TokenFallback 边界测试", count))
}

// TestTokenFallbackAttack : 恶意调用测试
func TestTokenFallbackAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("TokenFallback 恶意调用测试", count))
}

func to32bytes(src []byte) []byte {
	dst := common.BytesToHash(src)
	return dst[:]
}
