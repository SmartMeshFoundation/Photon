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
	selfLockAmounts := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	// get pre token balance
	//preTokenBalanceSelf, preTokenBalancePartner := getTokenBalance(self), getTokenBalance(partner)
	//preTokenBalanceContract := getTokenBalanceByAddess(env.TokenNetworkAddress)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)

	// self close channel
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)
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
		ChannelIdentifier:   bpSelf.ChannelIdentifier,
		OpenBlockNumber:     bpSelf.OpenBlockNumber,
		TokenNetworkAddress: env.TokenNetworkAddress,
		ChainID:             bpSelf.ChainID,
		BeneficiaryAddress:  self.Address,
		LockHash:            lock.Hash(),
		AdditionalHash:      utils.EmptyHash,
		MerkleProof:         mtree.Proof2Bytes(proof),
	}
	tx, err = env.TokenNetwork.PunishObsoleteUnlock(self.Auth, self.Address, partner.Address, ou.LockHash, ou.AdditionalHash, ou.sign(partner.Key))
	assertTxSuccess(t, &count, tx, err)
	t.Log(endMsg("ChannelPunish 正确调用测试", count, self, partner))
}

// TestChannelPunishException : 异常调用测试
func TestChannelPunishException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelPunish 异常调用测试", count))

}

// TestChannelPunishEdge : 边界测试
func TestChannelPunishEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelPunish 边界测试", count))
}

// TestChannelPunishAttack : 恶意调用测试
func TestChannelPunishAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
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
