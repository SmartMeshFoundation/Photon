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
	tokenNetwork *rpc.TokenNetworkProxy, channelIdentifier *contracts.ChannelUniqueID, privkey *ecdsa.PrivateKey, client *helper.SafeEthClient, db channeltype.Db, closedBlock int64, MyAddress, PartnerAddress common.Address) *ExternalState {
	cs := &ExternalState{
		funcRegisterChannelForHashlock: fun,
		TokenNetwork:                   tokenNetwork,
		auth:                           bind.NewKeyedTransactor(privkey),
		privKey:                        privkey,
		Client:                         client,
		ChannelIdentifier:              *channelIdentifier,
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
	if e.ClosedBlock != 0 {
		result = utils.NewAsyncResult()
		result.Result <- fmt.Errorf("%s already closed,closeBlock=%d", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), e.ClosedBlock)
		return
	}
	//start tx close and wait.
	var Nonce uint64
	TransferAmount := utils.BigInt0
	var LocksRoot = utils.EmptyHash
	//var ChannelIdentifier common.Address = utils.EmptyAddress
	var MessageHash = utils.EmptyHash
	var Signature []byte
	if balanceProof != nil {
		Nonce = balanceProof.Nonce
		TransferAmount = balanceProof.TransferAmount
		LocksRoot = balanceProof.LocksRoot
		//ChannelIdentifier = balanceProof.ChannelIdentifier
		MessageHash = balanceProof.MessageHash
		Signature = balanceProof.Signature
	}
	result = e.TokenNetwork.CloseChannelAsync(e.PartnerAddress, TransferAmount, LocksRoot, Nonce, MessageHash, Signature)
	return
}

//UpdateTransfer call updateTransfer of contract
func (e *ExternalState) UpdateTransfer(bp *transfer.BalanceProofState) (result *utils.AsyncResult) {
	if bp == nil {
		result = utils.NewAsyncResult()
		result.Result <- errors.New("bp is nil")
		return
	}
	log.Info(fmt.Sprintf("UpdateTransfer %s called ,BalanceProofState=%s",
		utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.StringInterface(bp, 3)))
	result = e.TokenNetwork.UpdateBalanceProofAsync(e.PartnerAddress, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.MessageHash, bp.Signature)
	return
}

/*
Unlock call withdraw function of contract
调用者要确保不包含自己声明放弃过的锁
*/
/*
 *	Unlock : function to unlock.
 *
 *	Note that caller has to ensure that there aren't locks that claimed abandoned by him contained.
 */
func (e *ExternalState) Unlock(unlockproofs []*channeltype.UnlockProof, argTransferdAmount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	transferAmount := new(big.Int).Set(argTransferdAmount)
	go func() {
		log.Info(fmt.Sprintf("Unlock called %s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier)))
		failed := false
		for _, proof := range unlockproofs {
			if e.db.IsThisLockHasUnlocked(e.ChannelIdentifier.ChannelIdentifier, proof.Lock.LockSecretHash) {
				log.Info(fmt.Sprintf("Unlock secret has been used %s  %s", e.ChannelIdentifier.String(), utils.HPex(proof.Lock.LockSecretHash)))
				continue
			}
			err := e.TokenNetwork.Unlock(e.PartnerAddress, transferAmount, proof.Lock, mtree.Proof2Bytes(proof.MerkleProof))
			if err != nil {
				failed = true
			} else {
				/*
					allow try withdraw next time if not success?
				*/
				e.db.UnlockThisLock(e.ChannelIdentifier.ChannelIdentifier, proof.Lock.LockSecretHash)
				log.Info(fmt.Sprintf("Unlock success %s,proof=%s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.StringInterface1(proof)))
				/*
					一旦 unlock 成功,那么 transferAmount 就会发生变化,下次必须用新的 transferAmount
				*/
				// Once unlock succeed, then transferAmount is going to change
				// next time we must use a new transferAmount.
				transferAmount = transferAmount.Add(transferAmount, proof.Lock.Amount)
			}
		}
		if failed {
			result.Result <- fmt.Errorf("there are errors when Unlock on channel %s  for %s", utils.HPex(e.ChannelIdentifier.ChannelIdentifier), utils.APex2(e.MyAddress))
		} else {
			result.Result <- nil
		}
	}()
	return
}

//Settle call settle function of contract
func (e *ExternalState) Settle(MyTransferAmount, PartnerTransferAmount *big.Int, MyLocksroot, PartnerLocksroot common.Hash) (result *utils.AsyncResult) {
	if e.SettledBlock != 0 {
		result = utils.NewAsyncResult()
		result.Result <- fmt.Errorf("channel %s already settled", e.ChannelIdentifier.String())
		return
	}
	log.Info(fmt.Sprintf("settle called %s,myTransferAmount=%s,partnerTransferAmount=%s,mylocksRoot=%s,partnerLocksroot=%s",
		e.ChannelIdentifier.String(), MyTransferAmount, PartnerTransferAmount,
		utils.HPex(MyLocksroot), utils.HPex(PartnerLocksroot),
	))
	result = e.TokenNetwork.SettleChannelAsync(e.MyAddress, e.PartnerAddress,
		MyTransferAmount, PartnerTransferAmount,
		MyLocksroot, PartnerLocksroot,
	)
	return
}

//Deposit call deposit of contract
func (e *ExternalState) Deposit(tokenAddress common.Address, amount *big.Int) (result *utils.AsyncResult) {
	result = e.TokenNetwork.DepositAsync(e.MyAddress, e.PartnerAddress, amount)
	return
}

/*
PunishObsoleteUnlock 惩罚对手 unlock 一个声明放弃了的锁.
*/
/*
 *	PunishObsoleteUnlock : function to punishment channel participant who unlocks a transfer lock that has been claimed abandoned.
 */
func (e *ExternalState) PunishObsoleteUnlock(lockhash, additionalHash common.Hash, cheaterSignature []byte) (result *utils.AsyncResult) {
	log.Info(fmt.Sprintf("PunishObsoleteUnlock called %s", e.ChannelIdentifier.String()))
	result = e.TokenNetwork.PunishObsoleteUnlockAsync(e.MyAddress, e.PartnerAddress, lockhash, additionalHash, cheaterSignature)
	return
}
