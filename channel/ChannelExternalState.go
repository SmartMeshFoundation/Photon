package channel

import (
	"fmt"
	"math/big"

	"errors"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//FuncRegisterChannelForHashlock is the callback for notify a new hashlock comes
type FuncRegisterChannelForHashlock func(channel *Channel, hashlock common.Hash)

/*
ExternalState operation on nettingchannelcontract
*/
type ExternalState struct {
	funcRegisterChannelForHashlock FuncRegisterChannelForHashlock
	TokenNetwork                   *rpc.TokenNetworkProxy
	auth                           *bind.TransactOpts
	privKey                        *ecdsa.PrivateKey
	Client                         *helper.SafeEthClient
	ClosedBlock                    int64
	SettledBlock                   int64
	ChannelIdentifier              contracts.ChannelUniqueID
	MyAddress                      common.Address
	PartnerAddress                 common.Address
	db                             channeltype.Db
}

//NewChannelExternalState create a new channel external state
func NewChannelExternalState(fun FuncRegisterChannelForHashlock,
	tokenNetwork *rpc.TokenNetworkProxy, channelAddress *contracts.ChannelUniqueID, privkey *ecdsa.PrivateKey, client *helper.SafeEthClient, db channeltype.Db, closedBlock int64, MyAddress, PartnerAddress common.Address) *ExternalState {
	cs := &ExternalState{
		funcRegisterChannelForHashlock: fun,
		TokenNetwork:                   tokenNetwork,
		auth:                           bind.NewKeyedTransactor(privkey),
		privKey:                        privkey,
		Client:                         client,
		ChannelIdentifier:              *channelAddress,
		db:                             db,
		ClosedBlock:                    closedBlock,
		SettledBlock:                   0,
		MyAddress:                      MyAddress,
		PartnerAddress:                 PartnerAddress,
	}
	return cs
}

//SetClosed set the closed blocknubmer of this channel
func (e *ExternalState) SetClosed(blocknumber int64) bool {
	if e.ClosedBlock != 0 {
		return false
	}
	e.ClosedBlock = blocknumber
	return true
}

//SetSettled set the settled number of this channel
func (e *ExternalState) SetSettled(blocknumber int64) bool {
	if e.SettledBlock != 0 && e.SettledBlock != blocknumber {
		return false
	}
	e.SettledBlock = blocknumber
	return true
}

//Close call close function of smart contract
//todo fix somany duplicate codes
func (e *ExternalState) Close(balanceProof *transfer.BalanceProofState) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	if e.ClosedBlock != 0 {
		result.Result <- fmt.Errorf("%s already closed,closeBlock=%d", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), e.ClosedBlock)
		return
	}
	//start tx close and wait.
	go func() {
		var Nonce int64
		TransferAmount := utils.BigInt0
		var LocksRoot = utils.EmptyHash
		//var ChannelAddress common.Address = utils.EmptyAddress
		var MessageHash = utils.EmptyHash
		var Signature []byte
		if balanceProof != nil {
			Nonce = balanceProof.Nonce
			TransferAmount = balanceProof.TransferAmount
			LocksRoot = balanceProof.LocksRoot
			//ChannelAddress = balanceProof.ChannelAddress
			MessageHash = balanceProof.MessageHash
			Signature = balanceProof.Signature
		}
		tx, err := e.TokenNetwork.GetContract().CloseChannel(e.auth,
			e.PartnerAddress,
			TransferAmount, LocksRoot, uint64(Nonce),
			MessageHash, Signature)
		if err != nil {
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("Close channel %s, txhash=%s", e.ChannelIdentifier.String(), tx.Hash().String()))
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
		if err != nil {
			result.Result <- err
			return
		}
		//log.Trace(fmt.Sprintf("receipt=%s", receipt))
		if receipt.Status != types.ReceiptStatusSuccessful {
			result.Result <- errors.New("tx execution failed")
			return
		}
		result.Result <- nil
		return
	}()
	return
}

//UpdateTransfer call updateTransfer of contract
func (e *ExternalState) UpdateTransfer(bp *transfer.BalanceProofState) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	if bp == nil {
		result.Result <- errors.New("bp is nil")
		return
	}
	go func() {
		log.Info(fmt.Sprintf("UpdateTransfer %s called ,BalanceProofState=%s",
			utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.StringInterface(bp, 3)))
		tx, err := e.TokenNetwork.GetContract().UpdateBalanceProof(e.auth, e.PartnerAddress, bp.TransferAmount, bp.LocksRoot, uint64(bp.Nonce),
			bp.MessageHash, bp.Signature)
		if err != nil {
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("UpdateTransfer %s, txhash=%s", e.ChannelIdentifier.String(), tx.Hash().String()))
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
		if err != nil {
			result.Result <- err
			return
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Info(fmt.Sprintf("updatetransfer failed %s,receipt=%s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), receipt))
			result.Result <- errors.New("tx execution failed")
			return
		}
		log.Info(fmt.Sprintf("updatetransfer success %s,balanceproof=%s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.StringInterface1(bp)))
		result.Result <- nil
	}()

	return
}

/*
Unlock call withdraw function of contract
调用者要确保不包含自己声明放弃过的锁
*/
func (e *ExternalState) Unlock(unlockproofs []*channeltype.UnlockProof, argTransferdAmount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	transferAmount := new(big.Int).Set(argTransferdAmount)
	go func() {
		log.Info(fmt.Sprintf("withdraw called %s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier)))
		failed := false
		for _, proof := range unlockproofs {
			if e.db.IsThisLockHasWithdraw(e.ChannelIdentifier.ChannelIdentifier, proof.Lock.LockSecretHash) {
				log.Info(fmt.Sprintf("withdraw secret has been used %s  %s", e.ChannelIdentifier, utils.HPex(proof.Lock.LockSecretHash)))
				continue
			}
			tx, err := e.TokenNetwork.GetContract().Unlock(
				e.auth,
				e.PartnerAddress,
				transferAmount,
				big.NewInt(proof.Lock.Expiration),
				proof.Lock.Amount,
				proof.Lock.LockSecretHash,
				mtree.Proof2Bytes(proof.MerkleProof))
			lock := proof.Lock
			if err != nil {
				failed = true
				log.Info(fmt.Sprintf("withdraw failed %s on channel %s,lock=%s", err, utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.StringInterface(lock, 7)))
				continue
				//return err
			}
			log.Info(fmt.Sprintf("withdraw on %s ,txhash=%s", e.ChannelIdentifier.String(), tx.Hash().String()))
			receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
			if err != nil {
				log.Info(fmt.Sprintf("Unlock failed with error:%s", err))
				failed = true
			}
			if receipt.Status != types.ReceiptStatusSuccessful {
				log.Info(fmt.Sprintf("withdraw failed %s,receipt=%s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), receipt))
				failed = true
			} else {
				/*
					allow try withdraw next time if not success?
				*/
				e.db.WithdrawThisLock(e.ChannelIdentifier.ChannelIdentifier, proof.Lock.LockSecretHash)
				log.Info(fmt.Sprintf("withdraw success %s,proof=%s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.StringInterface1(proof)))
				/*
					一旦 unlock 成功,那么 transferAmount 就会发生变化,下次必须用新的 transferAmount
				*/
				transferAmount = transferAmount.Add(transferAmount, proof.Lock.Amount)
			}
		}
		if failed {
			result.Result <- fmt.Errorf("there are errors when withdraw on channel %s  for %s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.APex2(e.MyAddress))
		} else {
			result.Result <- nil
		}
	}()
	return
}

//Settle call settle function of contract
func (e *ExternalState) Settle(MyTransferAmount, PartnerTransferAmount *big.Int, MyLocksroot, PartnerLocksroot common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	if e.SettledBlock != 0 {
		result.Result <- fmt.Errorf("channel %s already settled", e.ChannelIdentifier)
		return
	}
	go func() {
		log.Info(fmt.Sprintf("settle called %s", e.ChannelIdentifier))
		tx, err := e.TokenNetwork.GetContract().SettleChannel(e.auth, e.MyAddress, MyTransferAmount, MyLocksroot, e.PartnerAddress, PartnerTransferAmount, PartnerLocksroot)
		if err != nil {
			err = fmt.Errorf("settle failed %s,err=%s", e.ChannelIdentifier, err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("Settle Channel %s, err %s", e.ChannelIdentifier.String(), tx.Hash().String()))
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
		if err != nil {
			err = fmt.Errorf("%s settle WaitMined failed with error:%s", e.ChannelIdentifier, err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			err = fmt.Errorf("settle failed %s,receipt=%s", e.ChannelIdentifier, receipt)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("settle success %s", e.ChannelIdentifier))
		result.Result <- nil
	}()
	return
}

//Deposit call deposit of contract
func (e *ExternalState) Deposit(amount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		log.Info(fmt.Sprintf("Deposit called %s", e.ChannelIdentifier))
		tx, err := e.TokenNetwork.GetContract().Deposit(e.auth, e.MyAddress, e.PartnerAddress, amount)
		if err != nil {
			err = fmt.Errorf("%s Deposit failed %s", e.ChannelIdentifier, err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("Deposit to %s, txhash=%s", e.ChannelIdentifier.String(), tx.Hash().String()))
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
		if err != nil {
			err = fmt.Errorf("Deposit WaitMined failed with error:%s", err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			err = fmt.Errorf("Deposit failed %s,receipt=%s", e.ChannelIdentifier, receipt)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("Deposit success %s", e.ChannelIdentifier))
		result.Result <- nil
	}()
	return
}

/*
WithDraw on contract
应该是在我收到 withdrawresponse消息 以后调用
*/
func (e *ExternalState) WithDraw(myBalance, partnerBalance *big.Int, myWithdraw, partnerWithDraw *big.Int, mySignature, PartnerSignature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		log.Info(fmt.Sprintf("WithDraw called %s", e.ChannelIdentifier))
		tx, err := e.TokenNetwork.GetContract().WithDraw(
			e.auth,
			e.MyAddress, myBalance, myWithdraw,
			e.PartnerAddress, partnerBalance, partnerWithDraw,
			mySignature, PartnerSignature,
		)
		if err != nil {
			err = fmt.Errorf("%s WithDraw failed %s", e.ChannelIdentifier, err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("WithDraw to %s, txhash=%s", e.ChannelIdentifier.String(), tx.Hash().String()))
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
		if err != nil {
			err = fmt.Errorf("WithDraw WaitMined failed with error:%s", err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			err = fmt.Errorf("WithDraw failed %s,receipt=%s", e.ChannelIdentifier, receipt)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("WithDraw success %s", e.ChannelIdentifier))
		result.Result <- nil
	}()
	return
}

/*
PunishObsoleteUnlock 惩罚对手 unlock 一个声明放弃了的锁.
*/
func (e *ExternalState) PunishObsoleteUnlock(lockhash, additionalHash common.Hash, cheaterSignature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		log.Info(fmt.Sprintf("PunishObsoleteUnlock called %s", e.ChannelIdentifier))
		tx, err := e.TokenNetwork.GetContract().PunishObsoleteUnlock(e.auth, e.MyAddress, e.PartnerAddress, lockhash, additionalHash, cheaterSignature)
		if err != nil {
			err = fmt.Errorf("%s PunishObsoleteUnlock failed %s", e.ChannelIdentifier, err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("PunishObsoleteUnlock to %s, txhash=%s", e.ChannelIdentifier.String(), tx.Hash().String()))
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
		if err != nil {
			err = fmt.Errorf("PunishObsoleteUnlock WaitMined failed with error:%s", err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			err = fmt.Errorf("PunishObsoleteUnlock failed %s,receipt=%s", e.ChannelIdentifier, receipt)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("PunishObsoleteUnlock success %s", e.ChannelIdentifier))
		result.Result <- nil
	}()
	return
}

/*
CooperativeSettle 收到对方 cooperativeSettleReponse 消息以后调用
*/
func (e *ExternalState) CooperativeSettle(myBalance, partnerBalance *big.Int, mySignature, PartnerSignature []byte) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		log.Info(fmt.Sprintf("CooperativeSettle called %s", e.ChannelIdentifier))
		tx, err := e.TokenNetwork.GetContract().CooperativeSettle(e.auth, e.MyAddress, myBalance, e.PartnerAddress, partnerBalance, mySignature, PartnerSignature)
		if err != nil {
			err = fmt.Errorf("%s CooperativeSettle failed %s", e.ChannelIdentifier, err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("CooperativeSettle to %s, txhash=%s", e.ChannelIdentifier.String(), tx.Hash().String()))
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.Client, tx)
		if err != nil {
			err = fmt.Errorf("CooperativeSettle WaitMined failed with error:%s", err)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			err = fmt.Errorf("CooperativeSettle failed %s,receipt=%s", e.ChannelIdentifier, receipt)
			log.Info(err.Error())
			result.Result <- err
			return
		}
		log.Info(fmt.Sprintf("CooperativeSettle success %s", e.ChannelIdentifier))
		result.Result <- nil
	}()
	return
}
