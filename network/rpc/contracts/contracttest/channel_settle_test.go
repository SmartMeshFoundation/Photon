package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// TestChannelSettleRight : 正确调用测试
func TestChannelSettleRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	// cases
	runNoBpSettleTest(a1, a2, t, &count)                        // 无交易通道直接settle
	runReceiverSettleAfterSingleTransferTest(a1, a2, t, &count) // 一次交易后,收款方close后settle
	runPayerSettleAfterSingleTransferTest(a1, a2, t, &count)    // 一次交易后,付款方close,收款方先updateBalanceProof后settle
	runSettleWithUnregisteredLock(a1, a2, t, &count)            // 带未解的锁settle
	t.Log(endMsg("ChannelSettle 正确调用测试", count))
}

// TestChannelSettleException : 异常调用测试
func TestChannelSettleException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelSettle 异常调用测试", count))

}

// TestChannelSettleEdge : 边界测试
func TestChannelSettleEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelSettle 边界测试", count))
}

// TestChannelSettleAttack : 恶意调用测试
func TestChannelSettleAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelSettle 恶意调用测试", count))
}

// cases
// 无交易的channel直接settle
func runNoBpSettleTest(a1 *Account, a2 *Account, t *testing.T, count *int) {
	// get pre token balance
	preTokenBalanceA1, preTokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(10)
	depositA2 := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 1
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// close
	tx, err := env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// wait for settle
	waitToSettle(a1, a2)
	// settle
	tx, err = env.TokenNetwork.SettleChannel(a1.Auth, a1.Address, big.NewInt(0), utils.EmptyHash, a2.Address, big.NewInt(0), utils.EmptyHash)
	assertTxSuccess(t, count, tx, err)
	// get token balance after settle
	tokenBalanceA1, tokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, preTokenBalanceA1, tokenBalanceA1)
	assertEqual(t, count, preTokenBalanceA2, tokenBalanceA2)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}

// 一次交易后,收款方close后settle
func runReceiverSettleAfterSingleTransferTest(a1 *Account, a2 *Account, t *testing.T, count *int) {
	// get pre token balance
	preTokenBalanceA1, preTokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(1)
	depositA2 := big.NewInt(10)
	transferAmountA2 := big.NewInt(5)
	testSettleTimeout := TestSettleTimeoutMin + 1
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// a1 close
	bf2 := createPartnerBalanceProof(a1, a2, transferAmountA2, utils.EmptyHash, utils.EmptyHash, 1)
	tx, err := env.TokenNetwork.CloseChannel(
		a1.Auth, a2.Address, bf2.TransferAmount, bf2.LocksRoot, bf2.Nonce, bf2.AdditionalHash, bf2.Signature)
	assertTxSuccess(t, nil, tx, err)
	// wait for settle
	waitToSettle(a1, a2)
	// a1 settle
	tx, err = env.TokenNetwork.SettleChannel(a1.Auth,
		a1.Address, big.NewInt(0), utils.EmptyHash,
		a2.Address, transferAmountA2, utils.EmptyHash)
	assertTxSuccess(t, count, tx, err)
	// get token balance after settle
	tokenBalanceA1, tokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, preTokenBalanceA1.Add(preTokenBalanceA1, transferAmountA2), tokenBalanceA1)
	assertEqual(t, count, preTokenBalanceA2.Sub(preTokenBalanceA2, transferAmountA2), tokenBalanceA2)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}

// 一次交易后, 付款方close,收款方先updateBalanceProof后settle
func runPayerSettleAfterSingleTransferTest(a1 *Account, a2 *Account, t *testing.T, count *int) {
	// get pre token balance
	preTokenBalanceA1, preTokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(1)
	depositA2 := big.NewInt(10)
	transferAmountA2 := big.NewInt(5)
	testSettleTimeout := TestSettleTimeoutMin + 1
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// a2 close
	tx, err := env.TokenNetwork.CloseChannel(
		a2.Auth, a1.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// a1 updateProof
	bf := createPartnerBalanceProof(a1, a2, transferAmountA2, utils.EmptyHash, utils.EmptyHash, 1)
	tx, err = env.TokenNetwork.UpdateBalanceProof(a1.Auth, a2.Address, bf.TransferAmount, bf.LocksRoot, bf.Nonce, bf.AdditionalHash, bf.Signature)
	assertTxSuccess(t, count, tx, err)
	// wait for settle
	waitToSettle(a1, a2)
	// a1 settle
	tx, err = env.TokenNetwork.SettleChannel(a1.Auth,
		a1.Address, big.NewInt(0), utils.EmptyHash,
		a2.Address, transferAmountA2, utils.EmptyHash)
	assertTxSuccess(t, count, tx, err)
	// get token balance after settle
	tokenBalanceA1, tokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, preTokenBalanceA1.Add(preTokenBalanceA1, transferAmountA2), tokenBalanceA1)
	assertEqual(t, count, preTokenBalanceA2.Sub(preTokenBalanceA2, transferAmountA2), tokenBalanceA2)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}

// 带未解的锁settle
func runSettleWithUnregisteredLock(a1 *Account, a2 *Account, t *testing.T, count *int) {
	// get pre token balance
	preTokenBalanceA1, preTokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(a1, a2)
	depositA1 := big.NewInt(25)
	lockAmountA1 := big.NewInt(10)
	depositA2 := big.NewInt(30)
	transferAmountA2 := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 1
	openChannelAndDeposit(a1, a2, depositA1, depositA2, testSettleTimeout)
	// a1提交proof并close, 收到a2共20个token
	bp1 := createPartnerBalanceProof(a1, a2, transferAmountA2, utils.EmptyHash, utils.EmptyHash, 1)
	tx, err := env.TokenNetwork.CloseChannel(a1.Auth, a2.Address, bp1.TransferAmount, bp1.LocksRoot, bp1.Nonce, bp1.AdditionalHash, bp1.Signature)
	assertTxSuccess(t, nil, tx, err)
	// a2提交proof,锁定a1共10个token
	locks, _ := createLock(getLatestBlockNumber().Number.Int64()+10, lockAmountA1)
	bp2 := createPartnerBalanceProof(a2, a1, big.NewInt(0), mtree.NewMerkleTree(locks).MerkleRoot(), utils.EmptyHash, 1)
	tx, err = env.TokenNetwork.UpdateBalanceProof(a2.Auth, a1.Address, bp2.TransferAmount, bp2.LocksRoot, bp2.Nonce, bp2.AdditionalHash, bp2.Signature)
	assertTxSuccess(t, nil, tx, err)
	// wait to settle
	waitToSettle(a1, a2)
	// a1 settle
	tx, err = env.TokenNetwork.SettleChannel(a1.Auth,
		a1.Address, big.NewInt(0), mtree.NewMerkleTree(locks).MerkleRoot(),
		a2.Address, transferAmountA2, utils.EmptyHash)
	assertTxSuccess(t, count, tx, err)
	// check balance
	tokenBalanceA1, tokenBalanceA2 := getTokenBalance(a1), getTokenBalance(a2)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	assertEqual(t, count, preTokenBalanceA1.Add(preTokenBalanceA1, transferAmountA2), tokenBalanceA1)
	assertEqual(t, count, preTokenBalanceA2.Sub(preTokenBalanceA2, transferAmountA2), tokenBalanceA2)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}
