package channel

import (
	"fmt"
	"math/big"

	"errors"

	"sync"

	"github.com/SmartMeshFoundation/raiden-network/abi/bind"
	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type FuncRegisterChannelForHashlock func(channel *Channel, hashlock common.Hash)

/*
operation on nettingchannelcontract
*/
type ChannelExternalState struct {
	funcRegisterChannelForHashlock FuncRegisterChannelForHashlock
	NettingChannel                 *rpc.NettingChannelContractProxy
	bcs                            *rpc.BlockChainService
	OpenedBlock                    int64
	ClosedBlock                    int64
	SettledBlock                   int64
	ChanClosed                     chan struct{}
	ChanSettled                    chan struct{}
	IsCallClose                    bool
	IsCallSettle                   bool
	ChannelAddress                 common.Address
	lock                           sync.Mutex
	db                             ChannelDb
}

func NewChannelExternalState(fun FuncRegisterChannelForHashlock,
	nettingChannel *rpc.NettingChannelContractProxy, channelAddress common.Address, bcs *rpc.BlockChainService, db ChannelDb) *ChannelExternalState {
	var err error
	cs := &ChannelExternalState{
		funcRegisterChannelForHashlock: fun,
		NettingChannel:                 nettingChannel,
		bcs:                            bcs,
		ChanClosed:                     make(chan struct{}, 1),
		ChanSettled:                    make(chan struct{}, 1),
		ChannelAddress:                 channelAddress,
		db:                             db,
	}
	cs.OpenedBlock, err = nettingChannel.Opened()
	if err != nil {
		//todo don't panic for network error
		panic(fmt.Sprintf("call contract error:", err))
	}
	cs.ClosedBlock, _ = nettingChannel.Closed()
	cs.SettledBlock = 0
	return cs
}

func (this *ChannelExternalState) SetClosed(blocknumber int64) bool {
	if this.ClosedBlock != 0 {
		return false
	}
	this.ClosedBlock = blocknumber
	this.ChanClosed <- struct{}{}
	return true
}
func (this *ChannelExternalState) SetSettled(blocknumber int64) bool {
	if this.SettledBlock != 0 && this.SettledBlock != blocknumber {
		return false
	}
	this.SettledBlock = blocknumber
	//bai:write many times to channel ,error todo ?
	this.ChanSettled <- struct{}{}
	return true
}

//todo fix somany duplicate codes
//call close function of smart contract
func (this *ChannelExternalState) Close(balanceProof *transfer.BalanceProofState) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if !this.IsCallClose {
		this.IsCallClose = true
		var Nonce int64 = 0
		TransferAmount := utils.BigInt0
		var LocksRoot common.Hash = utils.EmptyHash
		//var ChannelAddress common.Address = utils.EmptyAddress
		var MessageHash common.Hash = utils.EmptyHash
		var Signature []byte = nil
		if balanceProof != nil {
			Nonce = balanceProof.Nonce
			TransferAmount = balanceProof.TransferAmount
			LocksRoot = balanceProof.LocksRoot
			//ChannelAddress = balanceProof.ChannelAddress
			MessageHash = balanceProof.MessageHash
			Signature = balanceProof.Signature
		}
		tx, err := this.NettingChannel.GetContract().Close(this.bcs.Auth, uint64(Nonce),
			TransferAmount, LocksRoot, MessageHash, Signature)
		if err != nil {
			return err
		}
		receipt, err := bind.WaitMined(rpc.GetCallContext(), this.bcs.Client, tx)
		if err != nil {
			return err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return errors.New("tx execution failed")
		}
	}
	return nil
}

func (this *ChannelExternalState) UpdateTransfer(bp *transfer.BalanceProofState) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if bp != nil {
		log.Info(fmt.Sprintf("UpdateTransfer %s called ,BalanceProofState=%s",
			this.ChannelAddress.String(), bp))
		tx, err := this.NettingChannel.GetContract().UpdateTransfer(this.bcs.Auth, uint64(bp.Nonce), bp.TransferAmount, bp.LocksRoot,
			bp.MessageHash, bp.Signature)
		if err != nil {
			return err
		}
		receipt, err := bind.WaitMined(rpc.GetCallContext(), this.bcs.Client, tx)
		if err != nil {
			return err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Info("updatetransfer failed %s,receipt=%s", this.ChannelAddress.String(), receipt.String())
			return errors.New("tx execution failed")
		} else {
			log.Info("updatetransfer success %s,balanceproof=%s", this.ChannelAddress.String(), utils.StringInterface1(bp))
		}
	}
	return nil
}

func (this *ChannelExternalState) WithDraw(unlockproofs []*UnlockProof) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	log.Info(fmt.Sprintf("withdraw called %s", this.ChannelAddress.String()))
	failed := false
	for _, proof := range unlockproofs {
		if this.db.IsThisLockHasWithdraw(this.ChannelAddress, proof.Secret) {
			log.Info(fmt.Sprintf("withdraw secret has been used %s-%s", utils.APex(this.ChannelAddress), utils.HPex(proof.Secret)))
			continue
		}
		tx, err := this.NettingChannel.GetContract().Withdraw(this.bcs.Auth, proof.LockEncoded, transfer.Proof2Bytes(proof.MerkleProof), proof.Secret)
		lock := new(encoding.Lock)
		lock.FromBytes(proof.LockEncoded)
		if err != nil {
			failed = true
			log.Info(fmt.Sprintf("withdraw failed %s,lock=%s", this.ChannelAddress.String(), utils.StringInterface1(lock)))
			continue
			//return err
		}
		receipt, err := bind.WaitMined(rpc.GetCallContext(), this.bcs.Client, tx)
		if err != nil {
			log.Info(fmt.Sprintf("WithDraw failed with error:%s", err))
			failed = true
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Info("withdraw failed %s,receipt=%s", this.ChannelAddress.String(), receipt.String())
			failed = true
			//return errors.New("withdraw execution failed ,maybe reverted?")
		} else {
			/*
				allow try withdraw next time if not success?
			*/
			this.db.WithdrawThisLock(this.ChannelAddress, proof.Secret)
			log.Info("withdraw success %s,proof=%s", this.ChannelAddress.String(), utils.StringInterface1(proof))
		}
	}
	if failed {
		return fmt.Errorf("there are errors when withdraw on channel %s for %s", this.ChannelAddress, this.bcs.NodeAddress)
	}
	return nil
}

func (this *ChannelExternalState) Settle() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.IsCallSettle {
		return nil
	}
	this.IsCallSettle = true
	log.Info(fmt.Sprintf("settle called %s", this.ChannelAddress.String()))
	tx, err := this.NettingChannel.GetContract().Settle(this.bcs.Auth)
	if err != nil {
		log.Info(fmt.Sprintf("settle failed %s", this.ChannelAddress.String()))
		return err
		//return err
	}
	receipt, err := bind.WaitMined(rpc.GetCallContext(), this.bcs.Client, tx)
	if err != nil {
		log.Info(fmt.Sprintf("settle WaitMined failed with error:%s", err))
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info("settle failed %s,receipt=%s", this.ChannelAddress.String(), receipt.String())
		return errors.New("settle execution failed ,maybe reverted?")
	} else {
		log.Info("settle success %s", this.ChannelAddress.String())
	}
	return nil
}

func (this *ChannelExternalState) Deposit(amount *big.Int) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	log.Info(fmt.Sprintf("Deposit called %s", this.ChannelAddress.String()))
	tx, err := this.NettingChannel.GetContract().Deposit(this.bcs.Auth, amount)
	if err != nil {
		log.Info(fmt.Sprintf("Deposit failed %s", this.ChannelAddress.String()))
		return err
		//return err
	}
	receipt, err := bind.WaitMined(rpc.GetCallContext(), this.bcs.Client, tx)
	if err != nil {
		log.Info(fmt.Sprintf("Deposit WaitMined failed with error:%s", err))
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info("Deposit failed %s,receipt=%s", this.ChannelAddress.String(), receipt.String())
		return errors.New("Deposit execution failed ,maybe reverted?")
	} else {
		log.Info(fmt.Sprintf("Deposit success %s", this.ChannelAddress.String()))
	}
	return nil
}
