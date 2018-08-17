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
	runRightUnlockTest(a1, a2, t, &count) // 1. 双方正确unlock 2.重复unlock
	t.Log(endMsg("ChannelUnlock 正确调用测试", count))
}

// TestChannelUnlockException : 异常调用测试
func TestChannelUnlockException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	//cases
	runUnlockWithWrongLocksrootTest(a1, a2, t, &count) // 1. 错误的locksroot来unlock 2. 在settle之后unlock
	runUnlockAfterExpirationTest(a1, a2, t, &count)    // 在锁过期后unlock
	runUnlockAfterSettleTimeoutTest(a1, a2, t, &count) // 在锁过期之前,settleTimeout之后解锁
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
	// prepare
	a1, a2 := env.getTwoAccountWithoutChannelClose(t)
	//cases
	runUnlockWithTamperedProofTest(a1, a2, t, &count) // 1. 自己unlock自己的锁 2. 篡改密码后unlock
	t.Log(endMsg("ChannelUnlock 恶意调用测试", count))
}

// TestChannelUnlockDelegateAttack : 授权调用测试
func TestChannelUnlockDelegate(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	third := env.getRandomAccountExcept(t, self, partner)
	//cases
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	partnerLockAmounts := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	selfLockAmounts := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	testSettleTimeout := TestSettleTimeoutMin + 30
	// get pre token balance
	preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	// get channel info
	channelID, _, openBlockNumber, _, _, chainID := getChannelInfo(self, partner)
	// build locks
	locksSelf, secretsSelf := createLockByArray(expireBlockNumber, selfLockAmounts)
	mpSelf := mtree.NewMerkleTree(locksSelf)
	locksPartner, secretsPartner := createLockByArray(expireBlockNumber, partnerLockAmounts)
	mpPartner := mtree.NewMerkleTree(locksPartner)
	// register secrets
	registrySecrets(self, secretsSelf)
	registrySecrets(self, secretsPartner)
	// self close channel with partner's lock
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), mpPartner.MerkleRoot(), utils.EmptyHash, 3)
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)
	// partner update proof with self's lock
	bpSelf := createPartnerBalanceProof(partner, self, big.NewInt(0), mpSelf.MerkleRoot(), utils.EmptyHash, 4)
	tx, err = env.TokenNetwork.UpdateBalanceProof(partner.Auth, self.Address, bpSelf.TransferAmount, bpSelf.LocksRoot, bpSelf.Nonce, bpSelf.AdditionalHash, bpSelf.Signature)
	assertTxSuccess(t, nil, tx, err)
	// partner' lock
	partnerTransferAmount := bpPartner.TransferAmount
	for _, lock := range locksPartner {
		uf := &UnlockDelegateForContract{
			Agent:             third.Address,
			Expiraition:       lock.Expiration,
			Amount:            lock.Amount,
			SecretHash:        lock.LockSecretHash,
			ChannelIdentifier: channelID,
			ChainID:           chainID,
			OpenBlockNumber:   openBlockNumber,
			MerkleProof:       mtree.Proof2Bytes(mpPartner.MakeProof(lock.Hash())),
		}
		// third unlock partner's lock on behalf of partner with partner's sig, MUST FAIL
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			self.Address,
			partner.Address,
			partnerTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(partner.Key),
		)
		assertTxFail(t, &count, tx, err)
		// third unlock partner's lock on behalf of partner with self's sig, MUST FAIL
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			self.Address,
			partner.Address,
			partnerTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(self.Key),
		)
		assertTxFail(t, &count, tx, err)
		// third unlock partner's lock on behalf of self with partner's sig, MUST FAIL
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			partner.Address,
			self.Address,
			partnerTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(partner.Key),
		)
		assertTxFail(t, &count, tx, err)
		// third unlock partner's lock on behalf of self with self's sig, MUST SUCCESS
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			partner.Address,
			self.Address,
			partnerTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(self.Key),
		)
		assertTxSuccess(t, &count, tx, err)
		partnerTransferAmount = partnerTransferAmount.Add(partnerTransferAmount, lock.Amount)
	}
	// self' lock
	selfTransferAmount := bpSelf.TransferAmount
	for _, lock := range locksSelf {
		uf := &UnlockDelegateForContract{
			Agent:             third.Address,
			Expiraition:       lock.Expiration,
			Amount:            lock.Amount,
			SecretHash:        lock.LockSecretHash,
			ChannelIdentifier: channelID,
			ChainID:           chainID,
			OpenBlockNumber:   openBlockNumber,
			MerkleProof:       mtree.Proof2Bytes(mpSelf.MakeProof(lock.Hash())),
		}
		// third unlock self's lock on behalf of self with self's sig, MUST FAIL
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			partner.Address,
			self.Address,
			selfTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(self.Key),
		)
		assertTxFail(t, &count, tx, err)
		// third unlock self's lock on behalf of self with partner's sig, MUST FAIL
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			partner.Address,
			self.Address,
			selfTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(partner.Key),
		)
		assertTxFail(t, &count, tx, err)
		// third unlock self's lock on behalf of partner with self's sig, MUST FAIL
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			self.Address,
			partner.Address,
			selfTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(self.Key),
		)
		assertTxFail(t, &count, tx, err)
		// third unlock self's lock on behalf of partner with partner's sig, MUST SUCCESS
		tx, err = env.TokenNetwork.UnlockDelegate(third.Auth,
			self.Address,
			partner.Address,
			selfTransferAmount,
			big.NewInt(lock.Expiration),
			lock.Amount,
			lock.LockSecretHash,
			uf.MerkleProof,
			uf.sign(partner.Key),
		)
		assertTxSuccess(t, &count, tx, err)
		selfTransferAmount = selfTransferAmount.Add(selfTransferAmount, lock.Amount)
	}

	// settled for cases after this
	waitForSettle(testSettleTimeout)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, selfTransferAmount, bpSelf.LocksRoot, partner.Address, partnerTransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	// get token balance
	tokenBalanceSelf, tokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, &count, big.NewInt(4), selfTransferAmount)
	assertEqual(t, &count, big.NewInt(3), partnerTransferAmount)
	assertEqual(t, &count, preTokenBalanceSelf.Add(preTokenBalanceSelf, partnerTransferAmount).Sub(preTokenBalanceSelf, selfTransferAmount), tokenBalanceSelf)
	assertEqual(t, &count, preTokenBalancePartner.Add(preTokenBalancePartner, selfTransferAmount).Sub(preTokenBalancePartner, partnerTransferAmount), tokenBalancePartner)
	assertEqual(t, &count, preTokenBalanceContract, tokenBalanceContract)
	t.Log(endMsg("ChannelUnlock 授权调用测试", count))
}

// 1. 双方正确unlock
// 2. 重复unlock
func runRightUnlockTest(self, partner *Account, t *testing.T, count *int) {
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	partnerLockAmounts := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	selfLockAmounts := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	// get pre token balance
	preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(self, partner)
	testSettleTimeout := TestSettleTimeoutMin + 30
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	// build locks
	locksSelf, secretsSelf := createLockByArray(expireBlockNumber, selfLockAmounts)
	mpSelf := mtree.NewMerkleTree(locksSelf)
	locksPartner, secretsPartner := createLockByArray(expireBlockNumber, partnerLockAmounts)
	mpPartner := mtree.NewMerkleTree(locksPartner)
	// register secrets
	registrySecrets(self, secretsSelf)
	registrySecrets(self, secretsPartner)
	// self close channel with partner's lock
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), mpPartner.MerkleRoot(), utils.EmptyHash, 3)
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)
	// partner update proof with self's lock
	bpSelf := createPartnerBalanceProof(partner, self, big.NewInt(0), mpSelf.MerkleRoot(), utils.EmptyHash, 4)
	tx, err = env.TokenNetwork.UpdateBalanceProof(partner.Auth, self.Address, bpSelf.TransferAmount, bpSelf.LocksRoot, bpSelf.Nonce, bpSelf.AdditionalHash, bpSelf.Signature)
	assertTxSuccess(t, nil, tx, err)
	// self unlock with partner's lock -------Case1
	partnerTransferAmount := bpPartner.TransferAmount
	for _, lock := range locksPartner {
		proof := mpPartner.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(self.Auth, partner.Address, partnerTransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxSuccess(t, count, tx, err)
		partnerTransferAmount = partnerTransferAmount.Add(partnerTransferAmount, lock.Amount)
	}
	// partner unlock with self's lock -------Case1
	selfTransferAmount := bpSelf.TransferAmount
	for _, lock := range locksSelf {
		proof := mpSelf.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(partner.Auth, self.Address, selfTransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxSuccess(t, count, tx, err)
		selfTransferAmount = selfTransferAmount.Add(selfTransferAmount, lock.Amount)
	}
	// partner unlock with self's lock repeat -------Case2
	for _, lock := range locksSelf {
		proof := mpSelf.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(partner.Auth, self.Address, selfTransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxFail(t, count, tx, err)
	}
	// settled for cases after this
	waitForSettle(testSettleTimeout)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, selfTransferAmount, bpSelf.LocksRoot, partner.Address, partnerTransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	// get token balance
	tokenBalanceSelf, tokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, big.NewInt(4), selfTransferAmount)
	assertEqual(t, count, big.NewInt(3), partnerTransferAmount)
	assertEqual(t, count, preTokenBalanceSelf.Add(preTokenBalanceSelf, partnerTransferAmount).Sub(preTokenBalanceSelf, selfTransferAmount), tokenBalanceSelf)
	assertEqual(t, count, preTokenBalancePartner.Add(preTokenBalancePartner, selfTransferAmount).Sub(preTokenBalancePartner, partnerTransferAmount), tokenBalancePartner)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}

// 1. 错误的locksroot来unlock
// 2. 在settle之后unlock
func runUnlockWithWrongLocksrootTest(self, partner *Account, t *testing.T, count *int) {
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	lockAmounts := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(5)}
	fakeLockAmount := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(6)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	// get pre token balance
	preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
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
	// unlock with wrong locks -------Case1
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
	// get token balance
	tokenBalanceSelf, tokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, preTokenBalanceSelf, tokenBalanceSelf)
	assertEqual(t, count, preTokenBalancePartner, tokenBalancePartner)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
	// unlock after settle -------Case2
	mp = mtree.NewMerkleTree(locks)
	proof := mp.MakeProof(locks[0].Hash())
	tx, err = env.TokenNetwork.Unlock(self.Auth, partner.Address, bpPartner.TransferAmount, big.NewInt(locks[0].Expiration), locks[0].Amount, locks[0].LockSecretHash, mtree.Proof2Bytes(proof))
	assertTxFail(t, count, tx, err)
}

// 在锁过期后unlock
func runUnlockAfterExpirationTest(self, partner *Account, t *testing.T, count *int) {
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	lockAmounts := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(5)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 3
	testSettleTimeout := TestSettleTimeoutMin + 50
	// get pre token balance
	preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(self, partner)
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
	// wait for lock to expiration
	waitForSettle(uint64(3))
	// unlock with right locks
	for _, lock := range locks {
		proof := mp.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(self.Auth, partner.Address, bpPartner.TransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxFail(t, count, tx, err)
	}
	// settled for cases after this
	waitForSettle(testSettleTimeout)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	// get token balance
	tokenBalanceSelf, tokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, preTokenBalanceSelf, tokenBalanceSelf)
	assertEqual(t, count, preTokenBalancePartner, tokenBalancePartner)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}

// 1. 自己unlock自己的锁
// 2. 篡改密码后unlock
func runUnlockWithTamperedProofTest(self, partner *Account, t *testing.T, count *int) {
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	lockAmounts := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(5)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	// get pre token balance
	preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
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
	// partner unlock on behalf of partner -------Case1
	proof := mp.MakeProof(locks[0].Hash())
	tx, err = env.TokenNetwork.Unlock(partner.Auth, self.Address, bpPartner.TransferAmount, big.NewInt(locks[0].Expiration), locks[0].Amount, locks[0].LockSecretHash, mtree.Proof2Bytes(proof))
	assertTxFail(t, count, tx, err)
	// self unlock after change secret -------Case2
	locks, secrets = createLockByArray(expireBlockNumber, lockAmounts)
	mp = mtree.NewMerkleTree(locks)
	for _, lock := range locks {
		proof := mp.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(self.Auth, partner.Address, bpPartner.TransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxFail(t, count, tx, err)
	}
	// settled for cases after this
	waitForSettle(testSettleTimeout)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	// get token balance
	tokenBalanceSelf, tokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, preTokenBalanceSelf, tokenBalanceSelf)
	assertEqual(t, count, preTokenBalancePartner, tokenBalancePartner)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}

// 在锁过期之前,settleTimeout之后解锁
func runUnlockAfterSettleTimeoutTest(self, partner *Account, t *testing.T, count *int) {
	// transaction data
	depositSelf := big.NewInt(60)
	depositPartner := big.NewInt(60)
	lockAmounts := []*big.Int{big.NewInt(1), big.NewInt(3), big.NewInt(5)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	testSettleTimeout := TestSettleTimeoutMin + 5
	// get pre token balance
	preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// create new channel
	cooperativeSettleChannelIfExists(self, partner)
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
	// wait for settleTimeout
	waitForSettle(testSettleTimeout)
	// unlock with right locks
	for _, lock := range locks {
		proof := mp.MakeProof(lock.Hash())
		tx, err = env.TokenNetwork.Unlock(self.Auth, partner.Address, bpPartner.TransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
		assertTxFail(t, count, tx, err)
	}
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	// get token balance
	tokenBalanceSelf, tokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance
	assertEqual(t, count, preTokenBalanceSelf, tokenBalanceSelf)
	assertEqual(t, count, preTokenBalancePartner, tokenBalancePartner)
	assertEqual(t, count, preTokenBalanceContract, tokenBalanceContract)
}
