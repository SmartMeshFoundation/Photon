package rpc

import (
	"errors"
	"fmt"
	"math/big"

	"bytes"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//TokenNetworkProxy proxy of TokenNetwork Contract
type TokenNetworkProxy struct {
	Address common.Address //this contract address
	bcs     *BlockChainService
	ch      *contracts.TokenNetwork
}

//NewChannel create new channel ,block until a new channel create
func (t *TokenNetworkProxy) NewChannel(participantAddress, partnerAddress common.Address, settleTimeout int) (err error) {
	tx, err := t.ch.OpenChannel(t.bcs.Auth, participantAddress, partnerAddress, uint64(settleTimeout))
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("NewChannel txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("NewChannel failed %s,receipt=%s", utils.APex(t.Address), receipt))
		err = errors.New("NewChannel tx execution failed")
		return
	}
	log.Info(fmt.Sprintf("NewChannel success %s, partnerAddress=%s", utils.APex(t.Address), utils.APex(partnerAddress)))
	return
}

//NewChannelAsync create channel async
func (t *TokenNetworkProxy) NewChannelAsync(participantAddress, partnerAddress common.Address, settleTimeout int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.NewChannel(participantAddress, partnerAddress, settleTimeout)
		result.Result <- err
	}()
	return result
}
func to32bytes(src []byte) []byte {
	dst := common.BytesToHash(src)
	return dst[:]
}
func makeNewChannelAndDepositData(participantAddress, partnerAddress common.Address, settleTimeout int) []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(1))) //open and deposit
	_, err = buf.Write(to32bytes(participantAddress[:]))
	_, err = buf.Write(to32bytes(partnerAddress[:]))
	_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(int64(settleTimeout)))) //settle_timeout
	if err != nil {
		log.Error(fmt.Sprintf("buf write err %s", err))
	}
	return buf.Bytes()
}
func (t *TokenNetworkProxy) newChannelAndDepositByApproveAndCall(token *TokenProxy, participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	data := makeNewChannelAndDepositData(participantAddress, partnerAddress, settleTimeout)
	return token.ApproveAndCall(t.Address, amount, data)
}
func (t *TokenNetworkProxy) newChannelAndDepositByFallback(token *TokenProxy, participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	data := makeNewChannelAndDepositData(participantAddress, partnerAddress, settleTimeout)
	return token.TransferWithFallback(t.Address, amount, data)
}
func (t *TokenNetworkProxy) newChannelAndDepositByApprove(token *TokenProxy, participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	err = token.Approve(t.Address, amount)
	if err != nil {
		return err
	}
	tx, err := t.GetContract().OpenChannelWithDeposit(t.bcs.Auth, participantAddress, partnerAddress, uint64(settleTimeout), amount)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("OpenChannelWithDeposit  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Warn(fmt.Sprintf("OpenChannelWithDeposit failed %s", receipt))
		return errors.New("OpenChannelWithDeposit tx execution failed")
	}
	log.Info(fmt.Sprintf("OpenChannelWithDeposit success %s ", utils.APex(t.Address)))
	return nil
}

//NewChannelAndDeposit create new channel ,block until a new channel create
func (t *TokenNetworkProxy) NewChannelAndDeposit(participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	tokenAddr, err := t.ch.Token(nil)
	if err != nil {
		return
	}
	token, err := t.bcs.Token(tokenAddr)
	if err != nil {
		return
	}
	err = t.newChannelAndDepositByFallback(token, participantAddress, partnerAddress, settleTimeout, amount)
	if err == nil {
		log.Trace(fmt.Sprintf("%s-%s newChannelAndDepositByFallback success", utils.APex(tokenAddr), utils.APex(participantAddress)))
		return
	}
	err = t.newChannelAndDepositByApproveAndCall(token, participantAddress, partnerAddress, settleTimeout, amount)
	if err == nil {
		log.Trace(fmt.Sprintf("%s-%s newChannelAndDepositByApproveAndCall success", utils.APex(tokenAddr), utils.APex(participantAddress)))
		return
	}
	return t.newChannelAndDepositByApprove(token, participantAddress, partnerAddress, settleTimeout, amount)
}

//NewChannelAndDepositAsync create channel async
func (t *TokenNetworkProxy) NewChannelAndDepositAsync(participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.NewChannelAndDeposit(participantAddress, partnerAddress, settleTimeout, amount)
		result.Result <- err
	}()
	return result
}

/*GetChannelInfo Returns the channel specific data.
@param participant1 Address of one of the channel participants.
@param participant2 Address of the other channel participant.
@return ch state and settle_block_number.
if state is 1, settleBlockNumber is settle timeout, if state is 2,settleBlockNumber is the min block number ,settle can be called.
*/
func (t *TokenNetworkProxy) GetChannelInfo(participant1, participant2 common.Address) (channelID common.Hash, settleBlockNumber, openBlockNumber uint64, state uint8, settleTimeout uint64, err error) {
	return t.ch.GetChannelInfo(t.bcs.getQueryOpts(), participant1, participant2)
}

//GetChannelParticipantInfo Returns Info of this channel.
//@return The address of the token.
func (t *TokenNetworkProxy) GetChannelParticipantInfo(participant, partner common.Address) (deposit *big.Int, balanceHash common.Hash, nonce uint64, err error) {
	deposit, h, nonce, err := t.ch.GetChannelParticipantInfo(t.bcs.getQueryOpts(), participant, partner)
	balanceHash = common.BytesToHash(h[:])
	return
}

//GetContract return contract
func (t *TokenNetworkProxy) GetContract() *contracts.TokenNetwork {
	return t.ch
}

//CloseChannel close channel
func (t *TokenNetworkProxy) CloseChannel(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (err error) {
	tx, err := t.GetContract().CloseChannel(t.bcs.Auth, partnerAddr, transferAmount, locksRoot, uint64(nonce), extraHash, signature)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("CloseChannel  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("CloseChannel failed %s", receipt))
		return errors.New("CloseChannel tx execution failed")
	}
	log.Info(fmt.Sprintf("CloseChannel success %s ,partner=%s", utils.APex(t.Address), utils.APex(partnerAddr)))
	return nil
}

//CloseChannelAsync close channel async
func (t *TokenNetworkProxy) CloseChannelAsync(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.CloseChannel(partnerAddr, transferAmount, locksRoot, nonce, extraHash, signature)
		result.Result <- err
	}()
	return
}

//UpdateBalanceProof update balance proof of partner
func (t *TokenNetworkProxy) UpdateBalanceProof(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (err error) {
	tx, err := t.GetContract().UpdateBalanceProof(t.bcs.Auth, partnerAddr, transferAmount, locksRoot, nonce, extraHash, signature)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("UpdateBalanceProof  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("UpdateBalanceProof failed %s", receipt))
		return errors.New("UpdateBalanceProof tx execution failed")
	}
	log.Info(fmt.Sprintf("UpdateBalanceProof success %s ,partner=%s", utils.APex(t.Address), utils.APex(partnerAddr)))
	return nil
}

//UpdateBalanceProofAsync update balance proof async
func (t *TokenNetworkProxy) UpdateBalanceProofAsync(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.UpdateBalanceProof(partnerAddr, transferAmount, locksRoot, nonce, extraHash, signature)
		result.Result <- err
	}()

	return
}

//Unlock a partner's lock
func (t *TokenNetworkProxy) Unlock(partnerAddr common.Address, transferAmount *big.Int, lock *mtree.Lock, proof []byte) (err error) {
	tx, err := t.GetContract().Unlock(t.bcs.Auth, partnerAddr, transferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, proof)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("Unlock  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("Unlock failed %s", receipt))
		return errors.New("Unlock tx execution failed")
	}
	log.Info(fmt.Sprintf("Unlock success %s ,partner=%s", utils.APex(t.Address), utils.APex(partnerAddr)))
	return nil
}

//UnlockAsync a partner's lock async
func (t *TokenNetworkProxy) UnlockAsync(partnerAddr common.Address, transferAmount *big.Int, lock *mtree.Lock, proof []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.Unlock(partnerAddr, transferAmount, lock, proof)
		result.Result <- err
	}()
	return
}

//SettleChannel settle a channel
func (t *TokenNetworkProxy) SettleChannel(p1Addr, p2Addr common.Address, p1Amount, p2Amount *big.Int, p1Locksroot, p2Locksroot common.Hash) (err error) {
	tx, err := t.GetContract().SettleChannel(t.bcs.Auth, p1Addr, p1Amount, p1Locksroot, p2Addr, p2Amount, p2Locksroot)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("SettleChannel  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Warn(fmt.Sprintf("SettleChannel failed %s", receipt))
		return errors.New("SettleChannel tx execution failed")
	}
	log.Info(fmt.Sprintf("SettleChannel success %s ", utils.APex(t.Address)))
	return nil
}

//SettleChannelAsync settle a channel async
func (t *TokenNetworkProxy) SettleChannelAsync(p1Addr, p2Addr common.Address, p1Amount, p2Amount *big.Int, p1Locksroot, p2Locksroot common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.SettleChannel(p1Addr, p2Addr, p1Amount, p2Amount, p1Locksroot, p2Locksroot)
		result.Result <- err
	}()
	return
}
func makeDepositData(participantAddress, partnerAddress common.Address) []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(2))) //open and deposit
	_, err = buf.Write(to32bytes(participantAddress[:]))
	_, err = buf.Write(to32bytes(partnerAddress[:]))
	if err != nil {
		log.Error(fmt.Sprintf("buf write err %s", err))
	}
	return buf.Bytes()
}
func (t *TokenNetworkProxy) depositByFallback(token *TokenProxy, participant, partner common.Address, amount *big.Int) (err error) {
	data := makeDepositData(participant, partner)
	return token.TransferWithFallback(t.Address, amount, data)
}
func (t *TokenNetworkProxy) depositByApproveAndCall(token *TokenProxy, participant, partner common.Address, amount *big.Int) (err error) {
	data := makeDepositData(participant, partner)
	return token.ApproveAndCall(t.Address, amount, data)
}
func (t *TokenNetworkProxy) depositByApprove(token *TokenProxy, participant, partner common.Address, amount *big.Int) (err error) {
	err = token.Approve(t.Address, amount)
	if err != nil {
		return
	}
	tx, err := t.GetContract().Deposit(t.bcs.Auth, participant, partner, amount)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("Deposit  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Warn(fmt.Sprintf("Deposit failed %s", receipt))
		return errors.New("Deposit tx execution failed")
	}
	log.Info(fmt.Sprintf("Deposit success %s ", utils.APex(t.Address)))
	return nil
}

//Deposit  to  a channel
func (t *TokenNetworkProxy) Deposit(participant, partner common.Address, amount *big.Int) (err error) {
	tokenAddr, err := t.ch.Token(nil)
	if err != nil {
		return
	}
	token, err := t.bcs.Token(tokenAddr)
	if err != nil {
		return
	}
	err = t.depositByFallback(token, participant, partner, amount)
	if err == nil {
		log.Trace(fmt.Sprintf("%s-%s depositByFallback success", utils.APex(tokenAddr), utils.APex(partner)))
		return
	}
	err = t.depositByApproveAndCall(token, participant, partner, amount)
	if err == nil {
		log.Trace(fmt.Sprintf("%s-%s depositByApproveAndCall success", utils.APex(tokenAddr), utils.APex(partner)))
		return
	}
	return t.depositByApprove(token, participant, partner, amount)
}

//DepositAsync to  a channel async
func (t *TokenNetworkProxy) DepositAsync(participant, partner common.Address, amount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.Deposit(participant, partner, amount)
		result.Result <- err
	}()
	return
}

//Withdraw  to  a channel
func (t *TokenNetworkProxy) Withdraw(p1Addr, p2Addr common.Address, p1Balance,
	p1Withdraw *big.Int, p1Signature, p2Signature []byte) (err error) {
	tx, err := t.GetContract().WithDraw(t.bcs.Auth, p1Addr, p2Addr, p1Balance, p1Withdraw,
		p1Signature, p2Signature,
	)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("Withdraw  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Warn(fmt.Sprintf("Withdraw failed %s", receipt))
		return errors.New("Withdraw tx execution failed")
	}
	log.Info(fmt.Sprintf("Withdraw success %s ", utils.APex(t.Address)))
	return nil
}

//WithdrawAsync   a channel async
func (t *TokenNetworkProxy) WithdrawAsync(p1Addr, p2Addr common.Address, p1Balance,
	p1Withdraw *big.Int, p1Signature, p2Signature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.Withdraw(p1Addr, p2Addr, p1Balance, p1Withdraw, p1Signature, p2Signature)
		result.Result <- err
	}()
	return
}

//PunishObsoleteUnlock  to  a channel
func (t *TokenNetworkProxy) PunishObsoleteUnlock(beneficiary, cheater common.Address, lockhash, extraHash common.Hash, cheaterSignature []byte) (err error) {
	tx, err := t.GetContract().PunishObsoleteUnlock(t.bcs.Auth, beneficiary, cheater, lockhash, extraHash, cheaterSignature)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("PunishObsoleteUnlock  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Warn(fmt.Sprintf("PunishObsoleteUnlock failed %s", receipt))
		return errors.New("PunishObsoleteUnlock tx execution failed")
	}
	log.Info(fmt.Sprintf("PunishObsoleteUnlock success %s ", utils.APex(t.Address)))
	return nil
}

//PunishObsoleteUnlockAsync   a channel async
func (t *TokenNetworkProxy) PunishObsoleteUnlockAsync(beneficiary, cheater common.Address, lockhash, extraHash common.Hash, cheaterSignature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.PunishObsoleteUnlock(beneficiary, cheater, lockhash, extraHash, cheaterSignature)
		result.Result <- err
	}()
	return
}

//CooperativeSettle  settle  a channel
func (t *TokenNetworkProxy) CooperativeSettle(p1Addr, p2Addr common.Address, p1Balance, p2Balance *big.Int, p1Signature, p2Signatue []byte) (err error) {
	tx, err := t.GetContract().CooperativeSettle(t.bcs.Auth, p1Addr, p1Balance, p2Addr, p2Balance, p1Signature, p2Signatue)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("CooperativeSettle  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Warn(fmt.Sprintf("CooperativeSettle failed %s", receipt))
		return errors.New("CooperativeSettle tx execution failed")
	}
	log.Info(fmt.Sprintf("CooperativeSettle success %s ", utils.APex(t.Address)))
	return nil
}

//CooperativeSettleAsync  settle  a channel async
func (t *TokenNetworkProxy) CooperativeSettleAsync(p1Addr, p2Addr common.Address, p1Balance, p2Balance *big.Int, p1Signature, p2Signatue []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		err := t.CooperativeSettle(p1Addr, p2Addr, p1Balance, p2Balance, p1Signature, p2Signatue)
		result.Result <- err
	}()
	return
}
