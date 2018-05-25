package channel

import (
	"fmt"
	"math/big"

	"errors"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
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
	NettingChannel                 *rpc.NettingChannelContractProxy
	bcs                            *rpc.BlockChainService
	OpenedBlock                    int64
	ClosedBlock                    int64
	SettledBlock                   int64
	ChanClosed                     chan struct{}
	ChanSettled                    chan struct{}
	ChannelAddress                 common.Address
	lock                           sync.Mutex
	db                             Db
}

//NewChannelExternalState create a new channel external state
func NewChannelExternalState(fun FuncRegisterChannelForHashlock,
	nettingChannel *rpc.NettingChannelContractProxy, channelAddress common.Address, bcs *rpc.BlockChainService, db Db) *ExternalState {
	var err error
	cs := &ExternalState{
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
		panic(fmt.Sprintf("call contract error: %s", err))
	}
	cs.ClosedBlock, _ = nettingChannel.Closed()
	cs.SettledBlock = 0
	return cs
}

//SetClosed set the closed blocknubmer of this channel
func (e *ExternalState) SetClosed(blocknumber int64) bool {
	if e.ClosedBlock != 0 {
		return false
	}
	e.ClosedBlock = blocknumber
	e.ChanClosed <- struct{}{}
	return true
}

//SetSettled set the settled number of this channel
func (e *ExternalState) SetSettled(blocknumber int64) bool {
	if e.SettledBlock != 0 && e.SettledBlock != blocknumber {
		return false
	}
	e.SettledBlock = blocknumber
	//bai:write many times to channel ,error todo ?
	e.ChanSettled <- struct{}{}
	return true
}

//Close call close function of smart contract
//todo fix somany duplicate codes
func (e *ExternalState) Close(balanceProof *transfer.BalanceProofState) error {
	e.lock.Lock()
	defer e.lock.Unlock()
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
	tx, err := e.NettingChannel.GetContract().Close(e.bcs.Auth, uint64(Nonce),
		TransferAmount, LocksRoot, MessageHash, Signature)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(rpc.GetCallContext(), e.bcs.Client, tx)
	if err != nil {
		return err
	}
	//log.Trace(fmt.Sprintf("receipt=%s", receipt))
	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("tx execution failed")
	}
	return nil
}

//UpdateTransfer call updateTransfer of contract
func (e *ExternalState) UpdateTransfer(bp *transfer.BalanceProofState) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if bp != nil {
		log.Info(fmt.Sprintf("UpdateTransfer %s called ,BalanceProofState=%s",
			utils.APex(e.ChannelAddress), utils.StringInterface(bp, 3)))
		tx, err := e.NettingChannel.GetContract().UpdateTransfer(e.bcs.Auth, uint64(bp.Nonce), bp.TransferAmount, bp.LocksRoot,
			bp.MessageHash, bp.Signature)
		if err != nil {
			return err
		}
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.bcs.Client, tx)
		if err != nil {
			return err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Info(fmt.Sprintf("updatetransfer failed %s,receipt=%s", utils.APex(e.ChannelAddress), receipt))
			return errors.New("tx execution failed")
		}
		log.Info(fmt.Sprintf("updatetransfer success %s,balanceproof=%s", utils.APex(e.ChannelAddress), utils.StringInterface1(bp)))
	}
	return nil
}

//WithDraw call withdraw function of contract
func (e *ExternalState) WithDraw(unlockproofs []*UnlockProof) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	log.Info(fmt.Sprintf("withdraw called %s", utils.APex(e.ChannelAddress)))
	failed := false
	for _, proof := range unlockproofs {
		if e.db.IsThisLockHasWithdraw(e.ChannelAddress, proof.Secret) {
			log.Info(fmt.Sprintf("withdraw secret has been used %s-%s", utils.APex(e.ChannelAddress), utils.HPex(proof.Secret)))
			continue
		}
		tx, err := e.NettingChannel.GetContract().Withdraw(e.bcs.Auth, e.bcs.NodeAddress, proof.LockEncoded, transfer.Proof2Bytes(proof.MerkleProof), proof.Secret)
		lock := new(encoding.Lock)
		lock.FromBytes(proof.LockEncoded)
		if err != nil {
			failed = true
			log.Info(fmt.Sprintf("withdraw failed %s on channel %s,lock=%s", err, utils.APex2(e.ChannelAddress), utils.StringInterface(lock, 7)))
			continue
			//return err
		}
		receipt, err := bind.WaitMined(rpc.GetCallContext(), e.bcs.Client, tx)
		if err != nil {
			log.Info(fmt.Sprintf("WithDraw failed with error:%s", err))
			failed = true
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Info(fmt.Sprintf("withdraw failed %s,receipt=%s", utils.APex2(e.ChannelAddress), receipt))
			failed = true
			//return errors.New("withdraw execution failed ,maybe reverted?")
		} else {
			/*
				allow try withdraw next time if not success?
			*/
			e.db.WithdrawThisLock(e.ChannelAddress, proof.Secret)
			log.Info(fmt.Sprintf("withdraw success %s,proof=%s", utils.APex2(e.ChannelAddress), utils.StringInterface1(proof)))
		}
	}
	if failed {
		return fmt.Errorf("there are errors when withdraw on channel %s  for %s", utils.APex2(e.ChannelAddress), utils.APex2(e.bcs.NodeAddress))
	}
	return nil
}

//Settle call settle function of contract
func (e *ExternalState) Settle() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	log.Info(fmt.Sprintf("settle called %s", utils.APex(e.ChannelAddress)))
	tx, err := e.NettingChannel.GetContract().Settle(e.bcs.Auth)
	if err != nil {
		log.Info(fmt.Sprintf("settle failed %s", utils.APex(e.ChannelAddress)))
		return err
		//return err
	}
	receipt, err := bind.WaitMined(rpc.GetCallContext(), e.bcs.Client, tx)
	if err != nil {
		log.Info(fmt.Sprintf("settle WaitMined failed with error:%s", err))
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("settle failed %s,receipt=%s", utils.APex(e.ChannelAddress), receipt))
		return errors.New("settle execution failed ,maybe reverted?")
	}
	log.Info(fmt.Sprintf("settle success %s", utils.APex(e.ChannelAddress)))
	return nil
}

//Deposit call deposit of contract
func (e *ExternalState) Deposit(amount *big.Int) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	log.Info(fmt.Sprintf("Deposit called %s", utils.APex(e.ChannelAddress)))
	tx, err := e.NettingChannel.GetContract().Deposit(e.bcs.Auth, amount)
	if err != nil {
		log.Info(fmt.Sprintf("Deposit failed %s", utils.APex(e.ChannelAddress)))
		return err
		//return err
	}
	receipt, err := bind.WaitMined(rpc.GetCallContext(), e.bcs.Client, tx)
	if err != nil {
		log.Info(fmt.Sprintf("Deposit WaitMined failed with error:%s", err))
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("Deposit failed %s,receipt=%s", utils.APex(e.ChannelAddress), receipt))
		return errors.New("Deposit execution failed ,maybe reverted?")
	}
	log.Info(fmt.Sprintf("Deposit success %s", utils.APex(e.ChannelAddress)))
	return nil
}
