package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

// TestUpdateBalanceProofRight : 正确调用测试
func TestUpdateBalanceProofRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, big.NewInt(15), big.NewInt(10), testSettleTimeout)
	// partner close channel
	tx, err := env.TokenNetwork.CloseChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// create balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(10), utils.EmptyHash, utils.EmptyHash, 5)

	// update right, MUST SUCCESS
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, &count, tx, err)

	// continue update right, MUST SUCCESS
	bpPartner.Nonce = bpPartner.Nonce + 1
	bpPartner.TransferAmount = bpPartner.TransferAmount.Add(bpPartner.TransferAmount, big.NewInt(10))
	bpPartner.sign(partner.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, &count, tx, err)

	// check self state after update
	depositSelf, balanceHashSelf, nonceSelf, err := env.TokenNetwork.GetChannelParticipantInfo(nil, self.Address, partner.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, &count, big.NewInt(15), depositSelf)
	assertEqual(t, &count, utils.EmptyHash[:24], balanceHashSelf[:])
	assertEqual(t, &count, uint64(0), nonceSelf)

	// check partner state after update
	depositPartner, balanceHashPartner, noncePartner, err := env.TokenNetwork.GetChannelParticipantInfo(nil, partner.Address, self.Address)
	localBpPartnerBalanceHash := bpPartner.BalanceData.Hash()
	assertSuccess(t, nil, err)
	assertEqual(t, &count, big.NewInt(10), depositPartner)
	assertEqual(t, &count, localBpPartnerBalanceHash[:24], balanceHashPartner[:])
	assertEqual(t, &count, bpPartner.Nonce, noncePartner)

	// settle
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(self.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	t.Log(endMsg("UpdateBalanceProof 正确调用测试", count))
}

// TestUpdateBalanceProofException : 异常调用测试
func TestUpdateBalanceProofException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	// create balance proof with nonexistent channel
	bpPartnerFake := createPartnerBalanceProof(self, partner, big.NewInt(10), utils.EmptyHash, utils.EmptyHash, 5)
	// update balance proof on nonexistent channel, MUST FAIL
	tx, err := env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartnerFake.TransferAmount, bpPartnerFake.LocksRoot, bpPartnerFake.Nonce, bpPartnerFake.AdditionalHash, bpPartnerFake.Signature)
	assertTxFail(t, &count, tx, err)
	// create channel
	openChannelAndDeposit(self, partner, big.NewInt(25), big.NewInt(0), testSettleTimeout)
	// create balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(10), utils.EmptyHash, utils.EmptyHash, 5)

	// update balance proof on open channel, MUST FAIL
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)

	// partner close channel
	tx, err = env.TokenNetwork.CloseChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)

	// update with no off-chain transfer, MUST FAIL
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxFail(t, &count, tx, err)
	bpPartnerEmpty := createPartnerBalanceProof(self, partner, big.NewInt(0), utils.EmptyHash, utils.EmptyHash, 0)
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartnerEmpty.TransferAmount, bpPartnerEmpty.LocksRoot, bpPartnerEmpty.Nonce, bpPartnerEmpty.AdditionalHash, bpPartnerEmpty.Signature)
	assertTxFail(t, &count, tx, err)

	// right update, MUST SUCCESS
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, &count, tx, err)

	// update with same nonce, MUST FAIL
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)

	// update with lower nonce, MUST FAIL
	bpPartner.Nonce = bpPartner.Nonce - 1
	bpPartner.sign(partner.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)

	// after settle timeout, MUST FAIL
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)

	// settle
	tx, err = env.TokenNetwork.SettleChannel(self.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	t.Log(endMsg("UpdateBalanceProof 异常调用测试", count))
}

// TestUpdateBalanceProofEdge : 边界测试
func TestUpdateBalanceProofEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, big.NewInt(15), big.NewInt(10), testSettleTimeout)
	// partner close channel
	tx, err := env.TokenNetwork.CloseChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// create balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(10), utils.EmptyHash, utils.EmptyHash, 5)
	// self update proof with wrong parm
	// EmptyAddress
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, EmptyAccountAddress, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)
	// wrong transfer amount
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, big.NewInt(1), bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)
	// wrong locksroot
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, common.HexToHash("0x123"), bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)
	// wrong nonce
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, 3, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)
	// wrong additional hash
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, common.HexToHash("0x123"), nil)
	assertTxFail(t, &count, tx, err)
	// wrong signature
	bpPartner.sign(self.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)
	// wait for settle
	waitToSettle(self, partner)
	// settle
	tx, err = env.TokenNetwork.SettleChannel(self.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, big.NewInt(0), utils.EmptyHash)
	assertTxSuccess(t, nil, tx, err)
	t.Log(endMsg("UpdateBalanceProof 边界测试", count))
}

// TestUpdateBalanceProofAttack : 恶意调用测试
func TestUpdateBalanceProofAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, big.NewInt(15), big.NewInt(10), testSettleTimeout)
	// partner close channel
	tx, err := env.TokenNetwork.CloseChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// create balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(10), utils.EmptyHash, utils.EmptyHash, 5)
	// update right, MUST SUCCESS
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, &count, tx, err)
	// settle
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(self.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)

	// reopen and update with old balance proof, MUST FAIL
	openChannelAndDeposit(self, partner, big.NewInt(15), big.NewInt(10), testSettleTimeout)
	tx, err = env.TokenNetwork.CloseChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// update with old balance proof, MUST FAIL
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxFail(t, &count, tx, err)
	// update with new balance proof, MUST SUCCESS
	bpPartner = createPartnerBalanceProof(self, partner, big.NewInt(10), utils.EmptyHash, utils.EmptyHash, 5)
	tx, err = env.TokenNetwork.UpdateBalanceProof(self.Auth, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot, bpPartner.Nonce, bpPartner.AdditionalHash, bpPartner.Signature)
	assertTxSuccess(t, &count, tx, err)

	// settle
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(self.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartner.TransferAmount, bpPartner.LocksRoot)
	assertTxSuccess(t, nil, tx, err)
	t.Log(endMsg("UpdateBalanceProof 恶意调用测试", count))
}

// TestChannelUnlockDelegateAttack : 授权调用测试
func TestUpdateBalanceProofDelegate(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 30
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	third := env.getRandomAccountExcept(t, self, partner)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, big.NewInt(15), big.NewInt(11), testSettleTimeout)
	// partner close channel
	tx, err := env.TokenNetwork.CloseChannel(partner.Auth, self.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	// create balance proof
	bpPartner := createPartnerBalanceProof(self, partner, big.NewInt(10), utils.EmptyHash, utils.Sha3([]byte("123")), 3)
	bpPartnerSelf := &BalanceProofUpdateForContracts{
		BalanceProofForContract: *bpPartner,
	}
	// update at the first half of settle window
	bpPartnerSelf.sign(self.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProofDelegate(third.Auth, partner.Address, self.Address, bpPartnerSelf.TransferAmount, bpPartnerSelf.LocksRoot, bpPartnerSelf.Nonce, bpPartnerSelf.AdditionalHash, bpPartnerSelf.Signature, bpPartnerSelf.NonClosingSignature)
	assertTxFail(t, &count, tx, err)

	// wait to the last half of settle window
	waitToUpdateBalanceProofDelegate(self, partner)

	// update with wrong signature
	bpPartnerSelf.sign(third.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProofDelegate(third.Auth, partner.Address, self.Address, bpPartnerSelf.TransferAmount, bpPartnerSelf.LocksRoot, bpPartnerSelf.Nonce, bpPartnerSelf.AdditionalHash, bpPartnerSelf.Signature, bpPartnerSelf.NonClosingSignature)
	assertTxFail(t, &count, tx, err)
	bpPartnerSelf.sign(partner.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProofDelegate(third.Auth, partner.Address, self.Address, bpPartnerSelf.TransferAmount, bpPartnerSelf.LocksRoot, bpPartnerSelf.Nonce, bpPartnerSelf.AdditionalHash, bpPartnerSelf.Signature, bpPartnerSelf.NonClosingSignature)
	assertTxFail(t, &count, tx, err)

	// update right
	bpPartnerSelf.sign(self.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProofDelegate(third.Auth, partner.Address, self.Address, bpPartnerSelf.TransferAmount, bpPartnerSelf.LocksRoot, bpPartnerSelf.Nonce, bpPartnerSelf.AdditionalHash, bpPartnerSelf.Signature, bpPartnerSelf.NonClosingSignature)
	assertTxSuccess(t, &count, tx, err)

	// settle
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(self.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, bpPartnerSelf.TransferAmount, bpPartnerSelf.LocksRoot)
	assertTxSuccess(t, nil, tx, err)

	// update after settle
	bpPartner = createPartnerBalanceProof(self, partner, big.NewInt(11), utils.EmptyHash, utils.Sha3([]byte("123")), 4)
	bpPartnerSelf = &BalanceProofUpdateForContracts{
		BalanceProofForContract: *bpPartner,
	}
	bpPartnerSelf.sign(self.Key)
	tx, err = env.TokenNetwork.UpdateBalanceProofDelegate(third.Auth, partner.Address, self.Address, bpPartnerSelf.TransferAmount, bpPartnerSelf.LocksRoot, bpPartnerSelf.Nonce, bpPartnerSelf.AdditionalHash, bpPartnerSelf.Signature, bpPartnerSelf.NonClosingSignature)
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("UpdateBalanceProof 授权调用测试", count, self, partner, third))
}
