package rpc

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"bytes"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts/test/tokens/smttoken"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

//RegistryProxy 只是为了表达方便,兼容以前代码,todo 完全去掉registry信息
type RegistryProxy struct {
	Address common.Address
	ch      *contracts.TokensNetwork
}

//TokenNetworkByToken get token
func (r *RegistryProxy) TokenNetworkByToken(token common.Address) (bool, error) {
	return r.ch.RegisteredToken(nil, token)
}

//GetContractVersion query contract version
func (r *RegistryProxy) GetContractVersion() (string, error) {
	return r.ch.ContractVersion(nil)
}

//GetContract return Contract interface
func (r *RegistryProxy) GetContract() *contracts.TokensNetwork {
	return r.ch
}

//TokenNetworkProxy proxy of TokenNetwork Contract
type TokenNetworkProxy struct {
	*RegistryProxy
	bcs   *BlockChainService
	token common.Address
}

func to32bytes(src []byte) []byte {
	dst := common.BytesToHash(src)
	return dst[:]
}
func makeNewChannelAndDepositData(participantAddress, partnerAddress common.Address, settleTimeout int) []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(to32bytes(participantAddress[:]))
	_, err = buf.Write(to32bytes(partnerAddress[:]))
	_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(int64(settleTimeout)))) //settle_timeout
	if err != nil {
		log.Error(fmt.Sprintf("buf write err %s", err))
	}
	return buf.Bytes()
}

//注意此函数并不会等待交易打包,只要交易进入了缓冲池就返回
func (t *TokenNetworkProxy) newChannelAndDepositByApproveAndCall(token *TokenProxy, participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	data := makeNewChannelAndDepositData(participantAddress, partnerAddress, settleTimeout)
	depositTXParams := &models.DepositTXParams{
		TokenAddress:       t.token,
		ParticipantAddress: participantAddress,
		PartnerAddress:     partnerAddress,
		Amount:             amount,
		SettleTimeout:      uint64(settleTimeout)}
	return token.ApproveAndCall(t.Address, amount, data, depositTXParams)
}

//注意这个函数并不会等待交易打包完成才返回,只要确定交易进入了缓冲池就返回
func (t *TokenNetworkProxy) newChannelAndDepositByFallback(token *TokenProxy, participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	data := makeNewChannelAndDepositData(participantAddress, partnerAddress, settleTimeout)
	depositTXParams := &models.DepositTXParams{
		TokenAddress:       t.token,
		ParticipantAddress: participantAddress,
		PartnerAddress:     partnerAddress,
		Amount:             amount,
		SettleTimeout:      uint64(settleTimeout)}
	return token.TransferWithFallback(t.Address, amount, data, depositTXParams)
}

/*
todo 目前这个处理流程有问题,必须要将相应的信息存入数据库中
*/
func (t *TokenNetworkProxy) newChannelAndDepositByApprove(token *TokenProxy, participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	log.Info(fmt.Sprintf("newChannelAndDepositByApprove participant=%s,partner=%s,settletimeout=%d,amount=%s,token=%s",
		utils.APex2(participantAddress), utils.APex2(partnerAddress), settleTimeout, amount, utils.APex2(t.token),
	))
	tx, err := token.Token.Approve(t.bcs.Auth, t.Address, amount)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(token.Address, t.Address, participantAddress, partnerAddress)
	txParams := &models.DepositTXParams{
		TokenAddress:       t.token,
		ParticipantAddress: participantAddress,
		PartnerAddress:     partnerAddress,
		Amount:             amount,
		SettleTimeout:      uint64(settleTimeout),
	}
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeApproveDeposit, channelID, 0, txParams)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("Approve %s, txhash=%s", utils.APex(t.Address), tx.Hash().String()))
	//go func() {
	//	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//	if err != nil {
	//		log.Error(fmt.Sprintf("Approve waitmined err,txhash=%s,err=%s", tx.Hash().String(), err))
	//		return
	//	}
	//	if receipt.Status != types.ReceiptStatusSuccessful {
	//		log.Error(fmt.Sprintf("Approve failed %s,receipt=%s", utils.APex(t.Address), receipt))
	//		return
	//	}
	//	log.Info(fmt.Sprintf("Approve success %s,spender=%s,value=%d", utils.APex(t.Address), utils.APex(t.Address), amount))
	//
	//	tx, err = t.GetContract().Deposit(t.bcs.Auth, t.token, participantAddress, partnerAddress, amount, uint64(settleTimeout))
	//	if err != nil {
	//		return
	//	}
	//	log.Info(fmt.Sprintf("OpenChannelWithDeposit  txhash=%s", tx.Hash().String()))
	//	receipt, err = bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//	if err != nil {
	//		log.Error(fmt.Sprintf("OpenChannelWithDeposit waitmined err, txhash=%s,err=%s", tx.Hash().String(), err))
	//		return
	//	}
	//	if receipt.Status != types.ReceiptStatusSuccessful {
	//		log.Error(fmt.Sprintf("OpenChannelWithDeposit failed %s", receipt))
	//		return
	//	}
	//	log.Info(fmt.Sprintf("OpenChannelWithDeposit success %s txhash=%s", utils.APex(t.Address), tx.Hash().String()))
	//	return
	//}()
	return nil

}

//NewChannelAndDeposit create new channel ,block until a new channel create
//func (t *TokenNetworkProxy) NewChannelAndDeposit(participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
//	log.Trace(fmt.Sprintf("NewChannelAndDeposit participant=%s,partner=%s,settletimeout=%d,amount=%s",
//		utils.APex2(participantAddress), utils.APex2(partnerAddress), settleTimeout, amount,
//	))
//	tokenAddr := t.token
//	if err != nil {
//		return
//	}
//	token, err := t.bcs.Token(tokenAddr)
//	if err != nil {
//		return
//	}
//	err = t.newChannelAndDepositByFallback(token, participantAddress, partnerAddress, settleTimeout, amount)
//	if err == nil {
//		log.Trace(fmt.Sprintf("%s-%s newChannelAndDepositByFallback success", utils.APex(tokenAddr), utils.APex(participantAddress)))
//		return
//	}
//	err = t.newChannelAndDepositByApproveAndCall(token, participantAddress, partnerAddress, settleTimeout, amount)
//	if err == nil {
//		log.Trace(fmt.Sprintf("%s-%s newChannelAndDepositByApproveAndCall success", utils.APex(tokenAddr), utils.APex(participantAddress)))
//		return
//	}
//	return t.newChannelAndDepositByApprove(token, participantAddress, partnerAddress, settleTimeout, amount)
//}

/*

 */
func (t *TokenNetworkProxy) newChannelAndDepositOnSMTToken(tokenAddress common.Address, participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	log.Info(fmt.Sprintf("deposit on SMTToken address=%s", tokenAddress.String()))
	smtTokenProxy, err := smttoken.NewSMTToken(tokenAddress, t.bcs.Client)
	if err != nil {
		log.Error(fmt.Sprintf("smttoken.NewSMTToken err = %s", err))
		return rerr.ContractCallError(err)
	}
	data := makeNewChannelAndDepositData(participantAddress, partnerAddress, settleTimeout)
	// 在Auth中设置金额,不用t.bcs.Auth,避免影响其他交易
	auth := bind.NewKeyedTransactor(t.bcs.PrivKey)
	auth.Value = amount
	tx, err := smtTokenProxy.BuyAndTransfer(auth, data)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	txParams := &models.DepositTXParams{
		TokenAddress:       tokenAddress,
		ParticipantAddress: participantAddress,
		PartnerAddress:     partnerAddress,
		Amount:             amount,
		SettleTimeout:      uint64(settleTimeout),
	}
	channelID := utils.CalcChannelID(txParams.TokenAddress, t.bcs.RegistryProxy.Address, txParams.ParticipantAddress, txParams.PartnerAddress)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeDeposit, channelID, 0, txParams)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	return
}

/*NewChannelAndDepositAsync create channel async
创建通道并存款和存款分两种情况,
一,只有一个Tx就能完成的情况,那么和关闭通道,settle通道处理流程是一样的
二,需要两个Tx,先Approve然后调用deposit,那么就需要详细规划
1. 首先approve初步验证没问题,就把相应的deposit信息存入数据库中
2. 收到approve事件以后,调取数据库中的信息
3. 继续deposit,并从数据库中删除记录, 如果失败,则需要专门通知用户失败了,
还要考虑重复的Deposit,如果数据库中有相应记录,不允许继续创建通道也不允许继续存款,这会覆盖上一次的操作.
*/
func (t *TokenNetworkProxy) NewChannelAndDepositAsync(participantAddress, partnerAddress common.Address, settleTimeout int, amount *big.Int) (err error) {
	log.Trace(fmt.Sprintf("NewChannelAndDeposit participant=%s,partner=%s,settletimeout=%d,amount=%s",
		utils.APex2(participantAddress), utils.APex2(partnerAddress), settleTimeout, amount,
	))
	tokenAddr := t.token
	if err != nil {
		return
	}
	token, err := t.bcs.Token(tokenAddr)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 获取tokenName,如果是SMTToken,即主链币代理合约, todo 移除这个查询,否则会造成不必要的网络访问,并且造成api阻塞
	name, err := token.Token.Name(nil)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	if name == params.SMTTokenName {
		return t.newChannelAndDepositOnSMTToken(tokenAddr, participantAddress, partnerAddress, settleTimeout, amount)
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

/*GetChannelInfo Returns the channel specific data.
@param participant1 Address of one of the channel participants.
@param participant2 Address of the other channel participant.
@return ch state and settle_block_number.
if state is 1, settleBlockNumber is settle timeout, if state is 2,settleBlockNumber is the min block number ,settle can be called.
*/
func (t *TokenNetworkProxy) GetChannelInfo(participant1, participant2 common.Address) (channelID common.Hash, settleBlockNumber, openBlockNumber uint64, state uint8, settleTimeout uint64, err error) {
	return t.ch.GetChannelInfo(t.bcs.getQueryOpts(), t.token, participant1, participant2)
}

//GetChannelParticipantInfo Returns Info of this channel.
//@return The address of the token.
func (t *TokenNetworkProxy) GetChannelParticipantInfo(participant, partner common.Address) (deposit *big.Int, balanceHash common.Hash, nonce uint64, err error) {
	deposit, h, nonce, err := t.ch.GetChannelParticipantInfo(t.bcs.getQueryOpts(), t.token, participant, partner)
	balanceHash = common.BytesToHash(h[:])
	return
}

//GetContract return contract
func (t *TokenNetworkProxy) GetContract() *contracts.TokensNetwork {
	return t.ch
}

//CloseChannel close channel
func (t *TokenNetworkProxy) CloseChannel(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (err error) {
	tx, err := t.GetContract().PrepareSettle(t.bcs.Auth, t.token, partnerAddr, transferAmount, locksRoot, uint64(nonce), extraHash, signature)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, t.bcs.Auth.From, partnerAddr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeClose, channelID, 0, &models.ChannelCloseOrChannelUpdateBalanceProofTXParams{
		TokenAddress:       t.token,
		ParticipantAddress: t.bcs.Auth.From,
		PartnerAddress:     partnerAddr,
		TransferAmount:     transferAmount,
		LocksRoot:          locksRoot,
		Nonce:              nonce,
		ExtraHash:          extraHash,
		Signature:          signature,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("CloseChannel  txhash=%s", tx.Hash().String()))
	//receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//if err != nil {
	//	return rerr.ErrTxWaitMined.AppendError(err)
	//}
	//if receipt.Status != types.ReceiptStatusSuccessful {
	//	log.Info(fmt.Sprintf("CloseChannel failed %s", receipt))
	//	return rerr.ErrTxReceiptStatus.Append("CloseChannel tx execution failed")
	//}
	//log.Info(fmt.Sprintf("CloseChannel success %s ,partner=%s", utils.APex(t.Address), utils.APex(partnerAddr)))
	return nil
}

//CloseChannelAsync close channel async 认为只要交易进入了缓冲池中,肯定会成功.
func (t *TokenNetworkProxy) CloseChannelAsync(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (err error) {
	tx, err := t.GetContract().PrepareSettle(t.bcs.Auth, t.token, partnerAddr, transferAmount, locksRoot, uint64(nonce), extraHash, signature)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, t.bcs.Auth.From, partnerAddr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeClose, channelID, 0, &models.ChannelCloseOrChannelUpdateBalanceProofTXParams{
		TokenAddress:       t.token,
		ParticipantAddress: t.bcs.Auth.From,
		PartnerAddress:     partnerAddr,
		TransferAmount:     transferAmount,
		LocksRoot:          locksRoot,
		Nonce:              nonce,
		ExtraHash:          extraHash,
		Signature:          signature,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("CloseChannel  txhash=%s", tx.Hash().String()))
	//go func() {
	//	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//	if err != nil {
	//		log.Error(fmt.Sprintf("CloseChannel error ,partner=%s,error=%s", err, utils.APex2(partnerAddr)))
	//		return
	//	}
	//	if receipt.Status != types.ReceiptStatusSuccessful {
	//		log.Error(fmt.Sprintf("CloseChannel failed %s", receipt))
	//		return
	//	}
	//	log.Info(fmt.Sprintf("CloseChannel success %s ,partner=%s", utils.APex(t.Address), utils.APex(partnerAddr)))
	//}()

	return nil
}

//UpdateBalanceProof update balance proof of partner
func (t *TokenNetworkProxy) UpdateBalanceProof(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (err error) {
	tx, err := t.GetContract().UpdateBalanceProof(t.bcs.Auth, t.token, partnerAddr, transferAmount, locksRoot, nonce, extraHash, signature)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, t.bcs.Auth.From, partnerAddr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeUpdateBalanceProof, channelID, 0, &models.ChannelCloseOrChannelUpdateBalanceProofTXParams{
		TokenAddress:       t.token,
		ParticipantAddress: t.bcs.Auth.From,
		PartnerAddress:     partnerAddr,
		TransferAmount:     transferAmount,
		LocksRoot:          locksRoot,
		Nonce:              nonce,
		ExtraHash:          extraHash,
		Signature:          signature,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("UpdateBalanceProof  txhash=%s", tx.Hash().String()))
	//receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//if err != nil {
	//	return rerr.ErrTxWaitMined.AppendError(err)
	//}
	//if receipt.Status != types.ReceiptStatusSuccessful {
	//	log.Info(fmt.Sprintf("UpdateBalanceProof failed %s", receipt))
	//	return rerr.ErrTxReceiptStatus.Append("UpdateBalanceProof tx execution failed")
	//}
	//log.Info(fmt.Sprintf("UpdateBalanceProof success %s ,partner=%s", utils.APex(t.Address), utils.APex(partnerAddr)))
	return nil
}

//UpdateBalanceProofAsync update balance proof async
func (t *TokenNetworkProxy) UpdateBalanceProofAsync(partnerAddr common.Address, transferAmount *big.Int, locksRoot common.Hash, nonce uint64, extraHash common.Hash, signature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		/*
			异步的链上操作应该是分成两步的,第一步是合约直接调用获取TxHash,如果没问题,那么我们认为第二步的WaitMined一定不会出问题
		*/
		err := t.UpdateBalanceProof(partnerAddr, transferAmount, locksRoot, nonce, extraHash, signature)
		result.Result <- err
	}()

	return
}

//Unlock a partner's lock
func (t *TokenNetworkProxy) Unlock(partnerAddr common.Address, transferAmount *big.Int, lock *mtree.Lock, proof []byte) (err error) {
	tx, err := t.GetContract().Unlock(t.bcs.Auth, t.token, partnerAddr, transferAmount, big.NewInt(lock.Expiration), lock.Amount, lock.LockSecretHash, proof)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, t.bcs.Auth.From, partnerAddr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeUnlock, channelID, 0, &models.UnlockTXParams{
		TokenAddress:       t.token,
		ParticipantAddress: t.bcs.Auth.From,
		PartnerAddress:     partnerAddr,
		TransferAmount:     transferAmount,
		Expiration:         big.NewInt(lock.Expiration),
		Amount:             lock.Amount,
		LockSecretHash:     lock.LockSecretHash,
		Proof:              proof,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("Unlock  txhash=%s", tx.Hash().String()))
	//receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//if err != nil {
	//	return rerr.ErrTxWaitMined.AppendError(err)
	//}
	//if receipt.Status != types.ReceiptStatusSuccessful {
	//	log.Info(fmt.Sprintf("Unlock failed %s", receipt))
	//	return rerr.ErrTxReceiptStatus.Append("Unlock tx execution failed")
	//}
	//log.Info(fmt.Sprintf("Unlock success %s ,partner=%s", utils.APex(t.Address), utils.APex(partnerAddr)))
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
	tx, err := t.GetContract().Settle(t.bcs.Auth, t.token, p1Addr, p1Amount, p1Locksroot, p2Addr, p2Amount, p2Locksroot)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, p1Addr, p2Addr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeSettle, channelID, 0, &models.ChannelSettleTXParams{
		TokenAddress:     t.token,
		P1Address:        p1Addr,
		P1TransferAmount: p1Amount,
		P1LocksRoot:      p1Locksroot,
		P2Address:        p2Addr,
		P2TransferAmount: p2Amount,
		P2LocksRoot:      p2Locksroot,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("SettleChannel  txhash=%s", tx.Hash().String()))
	//receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//if err != nil {
	//	return rerr.ErrTxWaitMined.AppendError(err)
	//}
	//if receipt.Status != types.ReceiptStatusSuccessful {
	//	log.Warn(fmt.Sprintf("SettleChannel failed %s", receipt))
	//	return rerr.ErrTxReceiptStatus.Append("SettleChannel tx execution failed")
	//}
	//log.Info(fmt.Sprintf("SettleChannel success %s ", utils.APex(t.Address)))
	return nil
}

//SettleChannelAsync settle a channel async 进入缓冲池就认为成功了
func (t *TokenNetworkProxy) SettleChannelAsync(p1Addr, p2Addr common.Address, p1Amount, p2Amount, p1Balance, p2Balance *big.Int, p1Locksroot, p2Locksroot common.Hash) (err error) {
	tx, err := t.GetContract().Settle(t.bcs.Auth, t.token, p1Addr, p1Amount, p1Locksroot, p2Addr, p2Amount, p2Locksroot)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, p1Addr, p2Addr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeSettle, channelID, 0, &models.ChannelSettleTXParams{
		TokenAddress:     t.token,
		P1Address:        p1Addr,
		P1TransferAmount: p1Amount,
		P1LocksRoot:      p1Locksroot,
		P2Address:        p2Addr,
		P2TransferAmount: p2Amount,
		P2LocksRoot:      p2Locksroot,
		P1Balance:        p1Balance,
		P2Balance:        p2Balance,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("SettleChannel  txhash=%s", tx.Hash().String()))
	//go func() {
	//	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//	if err != nil {
	//		log.Error(fmt.Sprintf("SettleChannel waitmined err %s", err))
	//		return
	//	}
	//	if receipt.Status != types.ReceiptStatusSuccessful {
	//		log.Error(fmt.Sprintf("SettleChannel failed %s", receipt))
	//		return
	//	}
	//	log.Info(fmt.Sprintf("SettleChannel success %s ", utils.APex(t.Address)))
	//}()
	return nil
}

//Withdraw  to  a channel
func (t *TokenNetworkProxy) Withdraw(p1Addr, p2Addr common.Address, p1Balance,
	p1Withdraw *big.Int, p1Signature, p2Signature []byte) (err error) {
	tx, err := t.GetContract().WithDraw(t.bcs.Auth, t.token, p1Addr, p2Addr, p1Balance, p1Withdraw,
		p1Signature, p2Signature,
	)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, p1Addr, p2Addr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeWithdraw, channelID, 0, &models.ChannelWithDrawTXParams{
		TokenAddress: t.token,
		P1Address:    p1Addr,
		P2Address:    p2Addr,
		P1Balance:    p1Balance,
		P1Withdraw:   p1Withdraw,
		P1Signature:  p1Signature,
		P2Signature:  p2Signature,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("Withdraw  txhash=%s", tx.Hash().String()))
	//receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//if err != nil {
	//	return rerr.ErrTxWaitMined.AppendError(err)
	//}
	//if receipt.Status != types.ReceiptStatusSuccessful {
	//	log.Warn(fmt.Sprintf("Withdraw failed %s", receipt))
	//	return rerr.ErrTxReceiptStatus.Append("Withdraw tx execution failed")
	//}
	//log.Info(fmt.Sprintf("Withdraw success %s ", utils.APex(t.Address)))
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
	tx, err := t.GetContract().PunishObsoleteUnlock(t.bcs.Auth, t.token, beneficiary, cheater, lockhash, extraHash, cheaterSignature)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, beneficiary, cheater)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypePunish, channelID, 0, &models.PunishObsoleteUnlockTXParams{
		TokenAddress:     t.token,
		Beneficiary:      beneficiary,
		Cheater:          cheater,
		LockHash:         lockhash,
		ExtraHash:        extraHash,
		CheaterSignature: cheaterSignature,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("PunishObsoleteUnlock  txhash=%s", tx.Hash().String()))
	//receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//if err != nil {
	//	return rerr.ErrTxWaitMined.AppendError(err)
	//}
	//if receipt.Status != types.ReceiptStatusSuccessful {
	//	log.Warn(fmt.Sprintf("PunishObsoleteUnlock failed %s", receipt))
	//	return rerr.ErrTxReceiptStatus.Append("PunishObsoleteUnlock tx execution failed")
	//}
	//log.Info(fmt.Sprintf("PunishObsoleteUnlock success %s ", utils.APex(t.Address)))
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
	tx, err := t.GetContract().CooperativeSettle(t.bcs.Auth, t.token, p1Addr, p1Balance, p2Addr, p2Balance, p1Signature, p2Signatue)
	if err != nil {
		return rerr.ContractCallError(err)
	}
	// 保存TXInfo并注册到bcs中监控其执行结果
	channelID := utils.CalcChannelID(t.token, t.Address, p1Addr, p2Addr)
	txInfo, err := t.bcs.TXInfoDao.NewPendingTXInfo(tx, models.TXInfoTypeCooperateSettle, channelID, 0, &models.ChannelCooperativeSettleTXParams{
		TokenAddress: t.token,
		P1Address:    p1Addr,
		P1Balance:    p1Balance,
		P2Address:    p2Addr,
		P2Balance:    p2Balance,
		P1Signature:  p1Signature,
		P2Signature:  p2Signatue,
	})
	if err != nil {
		return rerr.ContractCallError(err)
	}
	t.bcs.RegisterPendingTXInfo(txInfo)
	//log.Info(fmt.Sprintf("CooperativeSettle  txhash=%s", tx.Hash().String()))
	//receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	//if err != nil {
	//	return rerr.ErrTxWaitMined.AppendError(err)
	//}
	//if receipt.Status != types.ReceiptStatusSuccessful {
	//	log.Warn(fmt.Sprintf("CooperativeSettle failed %s", receipt))
	//	return rerr.ErrTxReceiptStatus.Append("CooperativeSettle tx execution failed")
	//}
	//log.Info(fmt.Sprintf("CooperativeSettle success %s ", utils.APex(t.Address)))
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
