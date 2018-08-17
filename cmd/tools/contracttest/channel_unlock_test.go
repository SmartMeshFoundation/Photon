package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// TestChannelUnlockRight : 正确调用测试
func TestChannelUnlockRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	// cases
	runRightUnlockTest(a1, a2, t, &count)              // 正确unlock
	runUnlockWithWrongLocksrootTest(a1, a2, t, &count) // 错误的locksroot来unlock
	t.Log(endMsg("ChannelUnlock 正确调用测试", count))
}

// TestChannelUnlockException : 异常调用测试
func TestChannelUnlockException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelUnlock 异常调用测试", count))

}

// TestChannelUnlockEdge : 边界测试
func TestChannelUnlockEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelUnlock 边界测试", count))
}

// TestChannelUnlockAttack : 恶意调用测试
func TestChannelUnlockAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelUnlock 恶意调用测试", count))
}

// 正确unlock
func runRightUnlockTest(self, partner *Account, t *testing.T, count *int) {
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	lockAmounts := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(5)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	// create new channel
	cooperativeSettleChannelIfExists(self, partner)
	testSettleTimeout := TestSettleTimeoutMin + 1
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	// build locks
	locks, secrets := createLockByArray(expireBlockNumber, lockAmounts)
	mp := mtree.NewMerkleTree(locks)
	// register secrets
	registrySecrets(self, secrets)
	// self close channel with right locks
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), mp.MerkleRoot(), utils.EmptyHash, 3)
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)
	// unlock with right locks
	partnerTransferAmount := bpPartner.TransferAmount
	for _, lock := range locks {
		proof := mp.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(self.Auth, partner.Address, partnerTransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxSuccess(t, count, tx, err)
		partnerTransferAmount = partnerTransferAmount.Add(partnerTransferAmount, lock.Amount)
	}
	// settled for cases after this
	waitForSettle(testSettleTimeout)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, partnerTransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
}

// 错误的locksroot来unlock
func runUnlockWithWrongLocksrootTest(self, partner *Account, t *testing.T, count *int) {
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	lockAmounts := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(5)}
	fakeLockAmount := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(6)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	// create new channel
	cooperativeSettleChannelIfExists(self, partner)
	testSettleTimeout := TestSettleTimeoutMin + 1
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	// build locks
	locks, secrets := createLockByArray(expireBlockNumber, lockAmounts)
	mp := mtree.NewMerkleTree(locks)
	// register secrets
	registrySecrets(self, secrets)
	// self close channel with right locks
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), mp.MerkleRoot(), utils.EmptyHash, 3)
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)
	// unlock with wrong locks
	fakeLocks, fakeSecrets := createLockByArray(expireBlockNumber, fakeLockAmount)
	registrySecrets(self, fakeSecrets)
	for _, lock := range fakeLocks {
		proof := mp.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(self.Auth, partner.Address, bpPartner.TransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxFail(t, count, tx, err)
	}
	// settled for cases after this
	waitForSettle(testSettleTimeout)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
}
