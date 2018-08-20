package contracttest

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

// TestChannelPunishRight : 正确调用测试
func TestChannelPunishRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 30
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	selfLockAmounts := []*big.Int{big.NewInt(1)}
	// get pre token balance
	preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)

	// self close channel
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(1), utils.EmptyHash, utils.EmptyHash, 1)
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)

	// partner update proof with locks
	locksSelf, secretsSelf := createLockByArray(expireBlockNumber, selfLockAmounts)
	registrySecrets(self, secretsSelf)
	mpSelf := mtree.NewMerkleTree(locksSelf)
	bpSelf := createPartnerBalanceProof(partner, self, big.NewInt(3), mpSelf.MerkleRoot(), utils.EmptyHash, 2)
	tx, err = env.TokenNetwork.UpdateBalanceProof(partner.Auth, self.Address, bpSelf.TransferAmount, bpSelf.LocksRoot, bpSelf.Nonce, bpSelf.AdditionalHash, bpSelf.Signature)
	assertTxSuccess(t, nil, tx, err)

	// partner unlock
	lock := locksSelf[0]
	proof := mpSelf.MakeProof(lock.Hash())
	tx, err = env.TokenNetwork.Unlock(partner.Auth, self.Address, bpSelf.TransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
	assertTxSuccess(t, nil, tx, err)

	// self punish partner
	ou := &ObseleteUnlockForContract{
		ChannelIdentifier:  bpSelf.ChannelIdentifier,
		OpenBlockNumber:    bpSelf.OpenBlockNumber,
		ChainID:            bpSelf.ChainID,
		BeneficiaryAddress: self.Address,
		LockHash:           lock.Hash(),
		AdditionalHash:     utils.EmptyHash,
		MerkleProof:        mtree.Proof2Bytes(proof),
	}
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxSuccess(t, &count, tx, err)

	// settled for cases after this
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)

	// get token balance after settle
	tokenBalanceSelf, tokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	tokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// check balance, self gets all token and partner gets 0
	assertEqual(t, &count, preTokenBalanceSelf.Add(preTokenBalanceSelf, depositPartner), tokenBalanceSelf)
	assertEqual(t, &count, preTokenBalancePartner.Sub(preTokenBalancePartner, depositPartner), tokenBalancePartner)
	assertEqual(t, &count, preTokenBalanceContract, tokenBalanceContract)

	t.Log(endMsg("ChannelPunish 正确调用测试", count, self, partner))
}

// TestChannelPunishException : 异常调用测试
func TestChannelPunishException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	// open channel
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 30
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)

	// create self locks
	selfLockAmounts := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	locksSelf, secretsSelf := createLockByArray(expireBlockNumber, selfLockAmounts)
	registrySecrets(self, secretsSelf)
	mpSelf := mtree.NewMerkleTree(locksSelf)
	lock := locksSelf[0]
	proof := mpSelf.MakeProof(lock.Hash())
	// create self balance proof with locks
	bpSelf := createPartnerBalanceProof(partner, self, big.NewInt(3), mpSelf.MerkleRoot(), utils.EmptyHash, 2)
	// create partner balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)
	// create punish param
	ou := &ObseleteUnlockForContract{
		ChannelIdentifier:  bpSelf.ChannelIdentifier,
		OpenBlockNumber:    bpSelf.OpenBlockNumber,
		ChainID:            bpSelf.ChainID,
		BeneficiaryAddress: self.Address,
		LockHash:           lock.Hash(),
		AdditionalHash:     utils.EmptyHash,
		MerkleProof:        mtree.Proof2Bytes(proof),
	}

	// 1. self punish partner on open channel, MUST FAIL
	tx, err := env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// self close channel
	tx, err = env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)

	// 2. self punish partner without partner update balance proof,MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// partner update proof with locks
	tx, err = env.TokenNetwork.UpdateBalanceProof(partner.Auth, self.Address, bpSelf.TransferAmount, bpSelf.LocksRoot, bpSelf.Nonce, bpSelf.AdditionalHash, bpSelf.Signature)
	assertTxSuccess(t, nil, tx, err)

	// 3. self punish partner without partner unlock,MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// settled for cases after this
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, bpSelf.TransferAmount, bpSelf.LocksRoot, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)

	// 4. self punish partner after settled,MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("ChannelPunish 异常调用测试", count))

}

// TestChannelPunishEdge : 边界测试
func TestChannelPunishEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	// open channel
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 30
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)

	// create self locks, the first amount = 0
	selfLockAmounts := []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	locksSelf, secretsSelf := createLockByArray(expireBlockNumber, selfLockAmounts)
	registrySecrets(self, secretsSelf)
	mpSelf := mtree.NewMerkleTree(locksSelf)
	lock := locksSelf[0]
	proof := mpSelf.MakeProof(lock.Hash())
	// create self balance proof with locks
	bpSelf := createPartnerBalanceProof(partner, self, big.NewInt(3), mpSelf.MerkleRoot(), utils.EmptyHash, 2)
	// create partner balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)

	// self close channel
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)

	// partner update proof with locks
	tx, err = env.TokenNetwork.UpdateBalanceProof(partner.Auth, self.Address, bpSelf.TransferAmount, bpSelf.LocksRoot, bpSelf.Nonce, bpSelf.AdditionalHash, bpSelf.Signature)
	assertTxSuccess(t, nil, tx, err)

	// partner unlock
	tx, err = env.TokenNetwork.Unlock(partner.Auth, self.Address, bpSelf.TransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
	assertTxSuccess(t, nil, tx, err)

	// create right param
	ou := &ObseleteUnlockForContract{
		ChannelIdentifier:  bpSelf.ChannelIdentifier,
		OpenBlockNumber:    bpSelf.OpenBlockNumber,
		ChainID:            bpSelf.ChainID,
		BeneficiaryAddress: self.Address,
		LockHash:           lock.Hash(),
		AdditionalHash:     utils.EmptyHash,
	}
	// self punish partner with EmptyAddress, MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, EmptyAccountAddress, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, EmptyAccountAddress, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// self punish partner with wrong AdditionalHash, MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, common.HexToHash("0x123"), ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// self punish partner with another BeneficiaryAddress, MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, partner.Address, self.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// wrong signature
	// self punish partner with wrong ChannelIdentifier, MUST FAIL
	ou.ChannelIdentifier = contracts.ChannelIdentifier(utils.EmptyHash)
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)
	ou.ChannelIdentifier = bpSelf.ChannelIdentifier

	// self punish partner with wrong OpenBlockNumber, MUST FAIL
	ou.OpenBlockNumber = 123
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)
	ou.OpenBlockNumber = bpSelf.OpenBlockNumber

	// self punish partner with another ChainId, MUST FAIL
	ou.ChainID = big.NewInt(99999)
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)
	ou.ChainID = bpSelf.ChainID

	// self punish partner with wrong LockHash, MUST FAIL
	ou.LockHash = locksSelf[1].Hash()
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)
	ou.LockHash = locksSelf[0].Hash()

	// self punish partner with wrong signer, MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(self.Key))
	assertTxFail(t, &count, tx, err)

	// settled for cases after this
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, bpSelf.TransferAmount.Add(bpSelf.TransferAmount, lock.Amount), bpSelf.LocksRoot, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)

	t.Log(endMsg("ChannelPunish 边界测试", count))
}

// TestChannelPunishAttack : 恶意调用测试
func TestChannelPunishAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	// open channel
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	testSettleTimeout := TestSettleTimeoutMin + 30
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)

	// create self locks, the first amount = 0
	selfLockAmounts := []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1)}
	expireBlockNumber := getLatestBlockNumber().Number.Int64() + 100
	locksSelf, secretsSelf := createLockByArray(expireBlockNumber, selfLockAmounts)
	registrySecrets(self, secretsSelf)
	mpSelf := mtree.NewMerkleTree(locksSelf)
	lock := locksSelf[0]
	proof := mpSelf.MakeProof(lock.Hash())
	// create self balance proof with locks
	bpSelf := createPartnerBalanceProof(partner, self, big.NewInt(3), mpSelf.MerkleRoot(), utils.EmptyHash, 2)
	// create partner balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)

	// self close channel
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, nil, tx, err)

	// partner update proof with locks
	tx, err = env.TokenNetwork.UpdateBalanceProof(partner.Auth, self.Address, bpSelf.TransferAmount, bpSelf.LocksRoot, bpSelf.Nonce, bpSelf.AdditionalHash, bpSelf.Signature)
	assertTxSuccess(t, nil, tx, err)

	// partner unlock
	tx, err = env.TokenNetwork.Unlock(partner.Auth, self.Address, bpSelf.TransferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, mtree.Proof2Bytes(proof))
	assertTxSuccess(t, nil, tx, err)

	// self punish partner
	ou := &ObseleteUnlockForContract{
		ChannelIdentifier: bpSelf.ChannelIdentifier,
		OpenBlockNumber:   bpSelf.OpenBlockNumber,
		ChainID:           bpSelf.ChainID,
		LockHash:          lock.Hash(),
		AdditionalHash:    utils.EmptyHash,
	}
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxSuccess(t, nil, tx, err)

	// settled for cases after this
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)

	// reopen
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	// create self locks, the first amount = 0
	locksSelfNew, secretsSelfNew := createLockByArray(expireBlockNumber, selfLockAmounts)
	registrySecrets(self, secretsSelfNew)
	mpSelfNew := mtree.NewMerkleTree(locksSelfNew)
	lockNew := locksSelfNew[0]
	proofNew := mpSelfNew.MakeProof(lockNew.Hash())
	// create self balance proof with locks
	bpSelfNew := createPartnerBalanceProof(partner, self, big.NewInt(3), mpSelfNew.MerkleRoot(), utils.EmptyHash, 2)
	// self close channel
	bpPartnerNew := createPartnerBalanceProof(self, partner, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)
	tx, err = env.TokenNetwork.CloseChannel(self.Auth, partner.Address, bpPartnerNew.TransferAmount, bpPartnerNew.LocksRoot, bpPartnerNew.Nonce, bpPartnerNew.AdditionalHash, bpPartnerNew.Signature)
	assertTxSuccess(t, nil, tx, err)
	// partner update proof with locks
	tx, err = env.TokenNetwork.UpdateBalanceProof(partner.Auth, self.Address, bpSelfNew.TransferAmount, bpSelfNew.LocksRoot, bpSelfNew.Nonce, bpSelfNew.AdditionalHash, bpSelfNew.Signature)
	assertTxSuccess(t, nil, tx, err)
	// partner unlock
	tx, err = env.TokenNetwork.Unlock(partner.Auth, self.Address, bpSelfNew.TransferAmount, big.NewInt(lockNew.Expiration), lockNew.Amount, lockNew.LockSecretHash, mtree.Proof2Bytes(proofNew))
	assertTxSuccess(t, nil, tx, err)

	// 1. self punish partner with old signature on old lock, MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// 2. self punish partner with old signature on new lock, MUST FAIL
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, lockNew.Hash(), ou.AdditionalHash, ou.sign(partner.Key))
	assertTxFail(t, &count, tx, err)

	// settled for cases after this
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(partner.Auth, self.Address, bpSelfNew.TransferAmount.Add(bpSelfNew.TransferAmount, lockNew.Amount), bpSelfNew.LocksRoot, partner.Address, bpPartnerNew.TransferAmount, bpPartnerNew.LocksRoot)
	assertTxSuccess(t, nil, tx, err)

	t.Log(endMsg("ChannelPunish 恶意调用测试", count))
}

// ObseleteUnlockForContract :
type ObseleteUnlockForContract struct {
	ChannelIdentifier            contracts.ChannelIdentifier
	BeneficiaryAddress           common.Address
	LockHash                     common.Hash
	BeneficiaryTransferredAmount *big.Int
	BeneficiaryNonce             *big.Int
	AdditionalHash               common.Hash
	TokenNetworkAddress          common.Address
	ChainID                      *big.Int
	OpenBlockNumber              uint64
	MerkleProof                  []byte
}

func (w *ObseleteUnlockForContract) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(w.LockHash[:])
	buf.Write(w.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, w.OpenBlockNumber)
	//buf.Write(w.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(w.ChainID))
	buf.Write(w.AdditionalHash[:])
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}
