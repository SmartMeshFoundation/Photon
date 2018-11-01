package channel

import (
	"errors"

	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/fee"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
Channel is the living representation of  channel on blockchain.
it contains all the transfers between two participants.
*/
type Channel struct {
	OurState          *EndState
	PartnerState      *EndState
	ExternState       *ExternalState
	ChannelIdentifier contracts.ChannelUniqueID //this channel
	TokenAddress      common.Address
	RevealTimeout     int
	SettleTimeout     int
	feeCharger        fee.Charger //calc fee for each transfer?
	State             channeltype.State
}

/*
NewChannel returns the living channel.
channelIdentifier must be a valid contract adress
settleTimeout must be valid, it cannot too small.
*/
func NewChannel(ourState, partnerState *EndState, externState *ExternalState, tokenAddr common.Address, channelIdentifier *contracts.ChannelUniqueID,
	revealTimeout, settleTimeout int) (c *Channel, err error) {
	if settleTimeout <= revealTimeout {
		err = fmt.Errorf("reveal_timeout can not be larger-or-equal to settle_timeout, reveal_timeout=%d,settle_timeout=%d", revealTimeout, settleTimeout)
		return
	}
	if revealTimeout < 3 {
		err = errors.New("reveal_timeout must be at least 3")
		return
	}
	/*
		考虑到避免分叉问题而延迟投递new channel事件的情况,构造channel时实时查询状态
	*/
	var stateOnChain = uint8(channeltype.StateOpened)
	_, _, _, stateOnChain, _, err = externState.TokenNetwork.GetChannelInfo(ourState.Address, partnerState.Address)
	if err != nil {
		log.Error(fmt.Sprintf("receive new channel,but can not get channel info from chain, err = %s", err.Error()))
	}
	c = &Channel{
		OurState:          ourState,
		PartnerState:      partnerState,
		ExternState:       externState,
		ChannelIdentifier: *channelIdentifier,
		TokenAddress:      tokenAddr,
		RevealTimeout:     revealTimeout,
		SettleTimeout:     settleTimeout,
		State:             channeltype.State(stateOnChain),
	}
	if externState.ClosedBlock != 0 {
		c.State = channeltype.StateClosed
	}
	if externState.SettledBlock != 0 {
		c.State = channeltype.StateSettled
	}
	return
}

/*
Distributable return the available amount of the token that our end of the channel can transfer to the partner.
*/
func (c *Channel) Distributable() *big.Int {
	return c.OurState.Distributable(c.PartnerState)
}

/*
CanTransfer  a closed channel and has no Balance channel cannot
transfer tokens to partner.
*/
func (c *Channel) CanTransfer() bool {
	return channeltype.CanTransferMap[c.State]
}

//CanContinueTransfer unfinished transfer can continue?
func (c *Channel) CanContinueTransfer() bool {
	return !channeltype.TransferCannotBeContinuedMap[c.State]
}

/*
ContractBalance Return the total amount of token we deposited in the channel
*/
func (c *Channel) ContractBalance() *big.Int {
	return c.OurState.ContractBalance
}

/*
TransferAmount Return how much we transferred to partner.
*/
func (c *Channel) TransferAmount() *big.Int {
	return c.OurState.TransferAmount()
}

/*
Balance Return our current Balance.

OurBalance is equal to `initial_deposit + received_amount - sent_amount`,
were both `receive_amount` and `sent_amount` are unlocked.
*/
func (c *Channel) Balance() *big.Int {
	x := new(big.Int)
	x.Sub(c.OurState.ContractBalance, c.OurState.TransferAmount())
	x.Add(x, c.PartnerState.TransferAmount())
	return x
}

/*
PartnerBalance return partner current Balance.
OurBalance is equal to `initial_deposit + received_amount - sent_amount`,
were both `receive_amount` and `sent_amount` are unlocked.
*/
func (c *Channel) PartnerBalance() *big.Int {
	x := new(big.Int)
	x.Sub(c.PartnerState.ContractBalance, c.PartnerState.TransferAmount())
	x.Add(x, c.OurState.TransferAmount())
	return x
}

/*
Locked return the current amount of our token that is locked waiting for a
        secret.

The locked value is equal to locked transfers that have been
initialized but their secret has not being revealed.
*/
func (c *Channel) Locked() *big.Int {
	return c.OurState.amountLocked()
}

/*
Outstanding is the tokens on road...
*/
func (c *Channel) Outstanding() *big.Int {
	return c.PartnerState.amountLocked()
}

/*
GetSettleExpiration returns how many blocks have to wait before settle.
*/
func (c *Channel) GetSettleExpiration(blocknumer int64) int64 {
	ClosedBlock := c.ExternState.ClosedBlock
	if ClosedBlock != 0 {
		return ClosedBlock + int64(c.SettleTimeout)
	}
	return blocknumer + int64(c.SettleTimeout)
}

/*
HandleBalanceProofUpdated 有可能对方使用了旧的信息,这样的话将会导致我无法 settle 通道
*/
/*
 *	HandleBalanceProofUpdated : It handles events that channel partners submitting BalanceProof that is not the most recent,
 * 		which leads to inability to settle channel.
 */
func (c *Channel) HandleBalanceProofUpdated(updatedParticipant common.Address, transferAmount *big.Int, locksRoot common.Hash) {
	endStateContractUpdated := c.OurState
	if updatedParticipant == c.PartnerState.Address {
		endStateContractUpdated = c.PartnerState
	}
	endStateContractUpdated.SetContractTransferAmount(transferAmount)
	endStateContractUpdated.SetContractLocksroot(locksRoot)
}

/*
HandleChannelPunished 发生了 Punish 事件,意味着受益方合约上的信息发生了变化.
*/
/*
 *	HandleChannelPunished : Punish event occurs,
 * 		which means that information on contract of beneficiary has been changed.
 */
func (c *Channel) HandleChannelPunished(beneficiaries common.Address) {
	var beneficiaryState, cheaterState *EndState
	if beneficiaries == c.OurState.Address {
		beneficiaryState = c.OurState
		cheaterState = c.PartnerState
	} else if beneficiaries == c.PartnerState.Address {
		beneficiaryState = c.PartnerState
		cheaterState = c.OurState
	} else {
		panic(fmt.Sprintf("channel=%s,but participant =%s",
			c.ChannelIdentifier.String(),
			beneficiaries.String(),
		))
	}
	beneficiaryState.SetContractTransferAmount(utils.BigInt0)
	beneficiaryState.SetContractLocksroot(utils.EmptyHash)
	beneficiaryState.SetContractNonce(0xfffffff)
	beneficiaryState.ContractBalance = beneficiaryState.ContractBalance.Add(
		beneficiaryState.ContractBalance, cheaterState.ContractBalance,
	)
	cheaterState.ContractBalance = new(big.Int).Set(utils.BigInt0)
}

/*
HandleClosed handles this channel was closed on blockchain
1. 更新NonClosing 一方的 ContractTransferAmount 和 LocksRoot,
2. 对方可能用旧的BalanceProof, 所以未必与我保存的 TransferAmount 和 LocksRoot一致
3. 如果我不是关闭方,那么需要更新对方的 BalanceProof
4. 我持有的知道密码的锁,需要解锁.
*/
/*
 *	HandleClosed : It handles events of closing channel.
 *
 *		1. Update ContractTransferAmount & LocksRoot of the non-closing participant.
 *		2. That participant may submit used BalanceProof, in which TransferAmount & LocksRoot are not consistent with mine.
 *		3. If I am not the non-closing participant, then update the BalanceProof of my channel partner.
 *		4. All locks I am holding that have known secrets must be unlocked.
 */
func (c *Channel) HandleClosed(closingAddress common.Address, transferredAmount *big.Int, locksRoot common.Hash) {
	endStateUpdatedOnContract := c.PartnerState
	balanceProof := c.PartnerState.BalanceProofState
	//依据合约上保存的 ContractTransferAmount 以及 LocksRoot 来更新我本地的
	//the channel was closed, update our half of the state if we need to
	if closingAddress != c.OurState.Address {
		c.ExternState.UpdateTransfer(balanceProof)
		endStateUpdatedOnContract = c.OurState
	}
	endStateUpdatedOnContract.SetContractTransferAmount(transferredAmount)
	endStateUpdatedOnContract.SetContractLocksroot(locksRoot)
	/*
		校验数据,如果没有用最新的数据来更新链上信息,有可能是一种攻击,也有可能是我本地的数据是错误的.
	*/
	// Verify data, if no more update message, which might be attack, or which might be local storage error.
	if endStateUpdatedOnContract.TransferAmount().Cmp(endStateUpdatedOnContract.contractTransferAmount()) != 0 {
		log.Error(fmt.Sprintf("Channel %s closed,but contract transfer amount is %s, and local stored %s's transfer amount is %s",
			utils.HPex(c.ChannelIdentifier.ChannelIdentifier), endStateUpdatedOnContract.contractTransferAmount(),
			utils.APex2(endStateUpdatedOnContract.Address), endStateUpdatedOnContract.TransferAmount(),
		))
		//todo 报告错误给最上层,可能是一个 bug? 一种攻击?,还是我自己存储数据有问题
		// todo throw error to the uppermost layer, maybe a bug? an attack? or just local storage error
	}
	if endStateUpdatedOnContract.locksRoot() != endStateUpdatedOnContract.contractLocksRoot() {
		log.Error(fmt.Sprintf("channel %s closed,but contract locksroot is %s, and local stored %s's locksroot is %s",
			utils.HPex(c.ChannelIdentifier.ChannelIdentifier), utils.HPex(endStateUpdatedOnContract.contractLocksRoot()),
			utils.APex2(endStateUpdatedOnContract.Address), utils.HPex(endStateUpdatedOnContract.locksRoot()),
		))
		//todo 报告错误给最上层,可能是一个 bug? 一种攻击?,还是我自己存储数据有问题
		// todo throw error to the uppermost layer, maybe a bug? an attack? or just local storage error.
	}
	unlockProofs := c.PartnerState.GetKnownUnlocks()
	if len(unlockProofs) > 0 {
		result := c.ExternState.Unlock(unlockProofs, c.PartnerState.contractTransferAmount())
		go func() {
			err := <-result.Result
			if err != nil {
				//todo 需要回报错误给Photon 调用者
				// todo need to report error to Photon
				log.Info(fmt.Sprintf("Unlock failed because of %s", err))
			}
		}()
	}

	c.State = channeltype.StateClosed
}

/*
HandleSettled handles this channel was settled on blockchain
there is nothing tod rightnow
*/
func (c *Channel) HandleSettled(blockNumber int64) {
	c.State = channeltype.StateSettled
}

//HandleWithdrawed 需要重新分配初始化整个通道的信息
/*
 *	HandleWithdrawed : function to handle withdraw message.
 *		This function will re-allocate the messages that initialize the whole payment channel.
 */
func (c *Channel) HandleWithdrawed(newOpenBlockNumber int64, participant1, participant2 common.Address, participant1Balance, participant2Balance *big.Int) {
	var p1, p2 *EndState
	if c.OurState.Address == participant1 && c.PartnerState.Address == participant2 {
		p1 = c.OurState
		p2 = c.PartnerState
	} else if c.OurState.Address == participant2 && c.PartnerState.Address == participant1 {
		p1 = c.PartnerState
		p2 = c.OurState
	} else {
		panic(fmt.Sprintf("channel event error, ourAddress=%s,partnerAddress=%s,p1=%s,p2=%s",
			c.OurState.Address.String(), c.PartnerState.Address.String(),
			participant1.String(), participant2.String(),
		))
	}
	if len(p1.Lock2UnclaimedLocks) > 0 || len(p2.Lock2UnclaimedLocks) > 0 {
		log.Warn(fmt.Sprintf("channel %s receive contract withdraw event, but has unclaimed locks."+
			"p1lock=%s,p2lock=%s", c.ChannelIdentifier.String(), utils.StringInterface(p1.Lock2UnclaimedLocks, 3),
			utils.StringInterface(p2.Lock2UnclaimedLocks, 3)))
	}
	/*
		通道所有的历史交易直接抛弃,并且不会在 settle 历史中保存,
	*/
	// all history record in channel should be abandoned, and do not store them in channel settle history.
	c.ChannelIdentifier.OpenBlockNumber = newOpenBlockNumber
	c.State = channeltype.StateOpened
	c.ExternState.ChannelIdentifier.OpenBlockNumber = newOpenBlockNumber
	c.ExternState.ClosedBlock = 0
	c.ExternState.SettledBlock = 0
	p1.ContractBalance = participant1Balance
	p1.BalanceProofState = transfer.NewEmptyBalanceProofState()
	p1.Lock2PendingLocks = make(map[common.Hash]channeltype.PendingLock)
	p1.Lock2UnclaimedLocks = make(map[common.Hash]channeltype.UnlockPartialProof)
	p1.Tree = mtree.NewMerkleTree(nil)
	p2.ContractBalance = participant2Balance
	p2.BalanceProofState = transfer.NewEmptyBalanceProofState()
	p2.Lock2PendingLocks = make(map[common.Hash]channeltype.PendingLock)
	p2.Lock2UnclaimedLocks = make(map[common.Hash]channeltype.UnlockPartialProof)
	p2.Tree = mtree.NewMerkleTree(nil)

}

/*
GetStateFor returns the latest status of one participant
*/
func (c *Channel) GetStateFor(nodeAddress common.Address) (*EndState, error) {
	if c.OurState.Address == nodeAddress {
		return c.OurState, nil
	}
	if c.PartnerState.Address == nodeAddress {
		return c.PartnerState, nil
	}
	return nil, fmt.Errorf("GetStateFor Unknown address %s", nodeAddress)
}

/*
RegisterSecret Register a secret to this channel

        This wont claim the lock (update the transferred_amount), it will only
        save the secret in case that a proof needs to be created. This method
        can be used for any of the ends of the channel.

        Note:
            When a secret is revealed a message could be in-transit containing
            the older lockroot, for this reason the recipient cannot update
            their locksroot at the moment a secret was revealed.

            The protocol is to register the secret so that it can compute a
            proof of Balance, if necessary, forward the secret to the sender
            and wait for the update from it. It's the sender's duty to order the
            current in-transit (and possible the transfers in queue) transfers
            and the secret/locksroot update.

            The channel and its queue must be changed in sync, a transfer must
            not be created while we update the balance_proof.

        Args:
            secret: The secret that releases a locked transfer.
*/
func (c *Channel) RegisterSecret(secret common.Hash) error {
	hashlock := utils.ShaSecret(secret[:])
	ourKnown := c.OurState.IsKnown(hashlock)
	partenerKnown := c.PartnerState.IsKnown(hashlock)
	if !ourKnown && !partenerKnown {
		return fmt.Errorf("secret doesn't correspond to a registered hashlock. hashlock %s token %s",
			utils.Pex(hashlock[:]), utils.HPex(c.ChannelIdentifier.ChannelIdentifier))
	}
	if ourKnown {
		lock := c.OurState.getLockByHashlock(hashlock)
		log.Debug(fmt.Sprintf("secret registered node=%s,from=%s,to=%s,token=%s,hashlock=%s, secret=%s, amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.OurState.Address[:]),
			utils.Pex(c.PartnerState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(hashlock[:]), utils.Pex(secret[:]), lock.Amount))
		err := c.OurState.RegisterSecret(secret)
		return err
	}
	if partenerKnown {
		lock := c.PartnerState.getLockByHashlock(hashlock)
		log.Debug(fmt.Sprintf("secret registered node=%s,from=%s,to=%s,token=%s,hashlock=%s, secret=%s, amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.PartnerState.Address[:]),
			utils.Pex(c.OurState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(hashlock[:]), utils.Pex(secret[:]), lock.Amount))
		err := c.PartnerState.RegisterSecret(secret)
		if err != nil {
			return err
		}
	}
	return nil
}

//RegisterRevealedSecretHash 链上对应的密码注册了
// RegisterRevealedSecretHash : secret has been registered on chain.
func (c *Channel) RegisterRevealedSecretHash(lockSecretHash, secret common.Hash, blockNumber int64) error {
	ourKnown := c.OurState.IsKnown(lockSecretHash)
	partenerKnown := c.PartnerState.IsKnown(lockSecretHash)
	if !ourKnown && !partenerKnown {
		return fmt.Errorf("LockSecretHash doesn't correspond to a registered lockSecretHash. lockSecretHash %s token %s",
			utils.Pex(lockSecretHash[:]), utils.HPex(c.ChannelIdentifier.ChannelIdentifier))
	}
	if ourKnown {
		lock := c.OurState.getLockByHashlock(lockSecretHash)
		log.Debug(fmt.Sprintf("lockSecretHash registered node=%s,from=%s,to=%s,token=%s,lockSecretHash=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.OurState.Address[:]),
			utils.Pex(c.PartnerState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(lockSecretHash[:]), lock.Amount))
		err := c.OurState.RegisterRevealedSecretHash(lockSecretHash, secret, blockNumber)
		if err == nil {
			//todo 需要发送给对方 unlock 消息,在哪里发比较合适呢? stateManager 还是这里?
			// todo need to send his partner unlock message. Where to send? In stateManager or right in this function ?
		}
		return err
	}
	if partenerKnown {
		lock := c.PartnerState.getLockByHashlock(lockSecretHash)
		log.Debug(fmt.Sprintf("lockSecretHash registered node=%s,from=%s,to=%s,token=%s,lockSecretHash=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.PartnerState.Address[:]),
			utils.Pex(c.OurState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(lockSecretHash[:]), lock.Amount))
		return c.PartnerState.RegisterRevealedSecretHash(lockSecretHash, secret, blockNumber)
	}
	return nil
}

//RegisterTransfer register a signed transfer, updating the channel's state accordingly.
//这些消息会改变 channel 的balance Proof
/*
 *	RegisterTransfer : register a signed transfer, updating the channel's state accordingly.
 *		This transfer will change BalanceProof of this channel.
 */
func (c *Channel) RegisterTransfer(blocknumber int64, tr encoding.EnvelopMessager) error {
	var err error
	switch msg := tr.(type) {
	case *encoding.MediatedTransfer:
		err = c.registerMediatedTranser(msg, blocknumber)
	case *encoding.DirectTransfer:
		err = c.registerDirectTransfer(msg, blocknumber)
	case *encoding.UnLock:
		err = c.registerUnlock(msg, blocknumber)
	case *encoding.AnnounceDisposedResponse:
		err = c.RegisterAnnounceDisposedResponse(msg, blocknumber)
	case *encoding.RemoveExpiredHashlockTransfer:
		err = c.RegisterRemoveExpiredHashlockTransfer(msg, blocknumber)
	default:
		return fmt.Errorf("receive unkonw transfer %s", tr)
	}
	return err
}

/*
PreCheckRecievedTransfer pre check received message(directtransfer,mediatedtransfer,refundtransfer) is valid or not
*/
func (c *Channel) PreCheckRecievedTransfer(tr encoding.EnvelopMessager) (fromState *EndState, toState *EndState, err error) {
	evMsg := tr.GetEnvelopMessage()
	if !c.isValidEnvelopMessage(evMsg) {
		err = fmt.Errorf("ch address mismatch,expect=%s,got=%s", c.ChannelIdentifier.String(), evMsg)
		return
	}
	if tr.GetSender() == c.OurState.Address {
		fromState = c.OurState
		toState = c.PartnerState
	} else if tr.GetSender() == c.PartnerState.Address {
		fromState = c.PartnerState
		toState = c.OurState
	} else {
		err = fmt.Errorf("received transfer from unknown address =%s", utils.APex(tr.GetSender()))
		return
	}
	/*
			  nonce is changed only when a transfer is un/registered, if the test
		     fails either we are out of sync, a message out of order, or it's a
		     forged transfer
			Strictly monotonic value used to order transfers. The nonce starts at 1
	*/
	isInvalidNonce := evMsg.Nonce < 1 || evMsg.Nonce != fromState.nonce()+1
	//If a node data is damaged, then the channel will not work, so the data must not be damaged.
	if isInvalidNonce {
		/*
			may occur on normal operation
			todo: give a example
		*/
		log.Info(fmt.Sprintf("invalid nonce node=%s,from=%s,to=%s,expected nonce=%d,nonce=%d",
			utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]),
			utils.Pex(toState.Address[:]), fromState.nonce()+1, evMsg.Nonce))
		err = rerr.InvalidNonce(utils.StringInterface(tr, 3))
		return
	}
	//  transfer amount should never decrese.
	if evMsg.TransferAmount.Cmp(fromState.TransferAmount()) < 0 {
		log.Error(fmt.Sprintf("NEGATIVE TRANSFER node=%s,from=%s,to=%s,transfer=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
			utils.StringInterface(tr, 3))) //for nest struct
		err = fmt.Errorf("negative transfer")
		return
	}
	return
}

/*
收到 unlock 消息:
1. nonce ,channel 要对
2. 验证密码有对应的锁
3. transferAmount 要想等
4. locksroot 要对,只是去掉了一个锁
*/
/*
 *	registerUnlock : function to receive unlock message.
 *
 *		1. value of nonce and channel should be correct.
 *		2. verify that the secret actually unlock a related hashlock in Unlock message.
 *		3. transferAmount should be equal to the one in BalanceProof.
 *		4. locksroot should be correct, but the hashlock verified in step 2 has been removed.
 */
func (c *Channel) registerUnlock(tr *encoding.UnLock, blockNumber int64) (err error) {
	fromState, _, err := c.PreCheckRecievedTransfer(tr)
	if err != nil {
		return
	}
	err = fromState.registerSecretMessage(tr)
	return err
}

/*
收到 DirectTransfer 消息:
1. nonce ,channel 要对
2. locksroot 要不变
3. 金额要增长,相等都是错的.
4. 账户要有这么多钱转
*/
/*
 *	registerDirectTransfer : function to register direct transfer.
 *
 *		1. nonce and channel should be correct.
 *		2. locksroot should not have any change.
 *		3. transferAmount should increase, if no change, then throw error.
 *		4. sufficient tokens should remain in accounts in order to process transfer.
 */
func (c *Channel) registerDirectTransfer(tr *encoding.DirectTransfer, blockNumber int64) (err error) {
	fromState, toState, err := c.PreCheckRecievedTransfer(tr)
	if err != nil {
		return
	}
	/*
		这次转账金额是多少
	*/
	// the amount of tokens this transfer takes.
	amount := new(big.Int).Set(tr.TransferAmount)
	amount = amount.Sub(amount, fromState.TransferAmount())
	/*
		转账金额是负数或者超过了可以给的金额,都是错的
	*/
	// It is error that token amount is negative or above available balance.
	if amount.Cmp(utils.BigInt0) <= 0 {
		return fmt.Errorf("direct transfer amount <0,amount=%s,message=%s", amount, tr)
	}
	if amount.Cmp(fromState.Distributable(toState)) > 0 {
		return fmt.Errorf("direct transfer amount too large,amount=%s,availabe=%s", amount, fromState.Distributable(toState))
	}
	err = fromState.registerDirectTransfer(tr)
	return err
}

/*
收到 MediatedTransfer 消息:
1. nonce,channel 要对
2. locksroot 要对,只是新增加了一个锁
3. transferAmount 要相等
4. 金额要够,
*/
/*
 *	registerMediatedTransfer : function to register MediatedTransfer.
 *
 *		1. nonce and channel should be correct.
 *		2. locksroot should be correct but with one more lock.
 *		3. transferAmount should be equal.
 *		4. there should be sufficient fund deposited in
 */
func (c *Channel) registerMediatedTranser(tr *encoding.MediatedTransfer, blockNumber int64) (err error) {
	fromState, toState, err := c.PreCheckRecievedTransfer(tr)
	if err != nil {
		return
	}
	/*
		这次转账金额是多少
	*/
	// the amount of tokens
	amount := tr.PaymentAmount
	/*
		转账金额是负数或者超过了可以给的金额,都是错的
	*/
	// fault occurs that token amount is negative or above available amount.
	if amount.Cmp(utils.BigInt0) <= 0 {
		return fmt.Errorf("mediated transfer amount <0,amount=%s,message=%s", amount, tr)
	}
	if amount.Cmp(fromState.Distributable(toState)) > 0 {
		return rerr.ErrInsufficientBalance
	}
	/*
				  For mediators: This is registering the *mediator* paying
		            transfer. The expiration of the lock must be `reveal_timeout`
		            blocks smaller than the *received* paying transfer. This cannot
		            be checked by the paying channel alone.

		            For the initiators: As there is no backing transfer, the
		            expiration is arbitrary, using the channel settle_timeout as an
		            upper limit because the node receiving the transfer will use it
		            as an upper bound while mediating.

		            For the receiver: A lock that expires after the settle period
		            just means there is more time to withdraw it.
	*/
	endSettlePeriod := c.GetSettleExpiration(blockNumber)
	expiresAfterSettle := tr.Expiration > endSettlePeriod
	/*
		我不能接收超过 settle timeout 的交易,这样对我不安全
		我也不能发出超过 settle timeout 的交易,这样不符合规则
		为什么我接收超过 settle timeout 的交易不安全?
		交易: A-B-C-D
		AB: settle timeout 1000
		BC settle timeout 10
		CD settle timeout 1000
		假设 当前块为20000,B 收到了来自 A 超时区块为超时时间为21000
		B给 C 超时时间21000,C给 D 超时时间21000
		那么 BD 可以合谋,D 告诉 B 密码, B close/settle 通道,然后 D 可以链上注册密码,取走相应 token
	*/
	/*
	 *	I can not receive transfers after settle_timeout, not secure.
	 *	I can not send transfers after settle_timeout, not abide to contract rules.
	 *	Why my receiving transfers after settle_timeout is not secure?
	 *
	 *	Transfer : A-B-C-D
	 *	AB : settle_timeout 1000
	 *	BC : settle_timeout 10
	 *	CD : settle_timeout	1000
	 *
	 *	Assume that current block height is 20000, and transfer expiration that B received A is 21000.
	 *	transfer expiration in BC 21000, transfer expiration in CD 21000
	 * 	then BC can collude and D reveal secret to B, after B close/settle channel, D can register the secret on-chain
	 *	and steal tokens in BC.
	 */
	if expiresAfterSettle { //After receiving this lock, the party can close or updatetransfer on the chain, so that if the party does not have a password, he still can't get the money.
		log.Error(fmt.Sprintf("Lock expires after the settlement period. node=%s,from=%s,to=%s,lockexpiration=%d,currentblock=%d,end_settle_period=%d",
			utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
			tr.Expiration, blockNumber, endSettlePeriod))
		return fmt.Errorf("lock expires after the settlement period")
	}
	err = fromState.registerMediatedMessage(tr)
	if err == nil {
		c.ExternState.funcRegisterChannelForHashlock(c, tr.LockSecretHash)
	}
	return err
}

/*
RegisterRemoveExpiredHashlockTransfer register a request to remove a expired hashlock and this hashlock must be sent out from the sender.
*/
func (c *Channel) RegisterRemoveExpiredHashlockTransfer(tr *encoding.RemoveExpiredHashlockTransfer, blockNumber int64) (err error) {
	return c.registerRemoveLock(tr, blockNumber, tr.LockSecretHash, true)
}

/*
RegisterAnnounceDisposedResponse 从我这里发出或者收到来自对方的announceDisposedTransferResponse,
注意收到对方消息的话,一定要验证事先发出去过AnnounceDisposedTransfer.
*/
/*
 *	RegisterAnnounceDisposedResponse : function to register AnnounceDisposedRespnse, and send out or receive announceDisposedTransferResponse from channel partner.
 *
 *		Note that everytime a participant receives message from his partner, he must verify the AnnounceDisposedTransfer he sent out beforehand.
 */
func (c *Channel) RegisterAnnounceDisposedResponse(response *encoding.AnnounceDisposedResponse, blockNumber int64) (err error) {
	return c.registerRemoveLock(response, blockNumber, response.LockSecretHash, false)
}
func (c *Channel) registerRemoveLock(messager encoding.EnvelopMessager, blockNumber int64, lockSecretHash common.Hash, mustExpired bool) (err error) {
	msg := messager.GetEnvelopMessage()
	fromState, _, err := c.PreCheckRecievedTransfer(messager)
	if err != nil {
		return
	}
	/*
		transfer amount should not change.
	*/
	if msg.TransferAmount.Cmp(fromState.TransferAmount()) != 0 {
		err = errTransferAmountMismatch
		return
	}
	_, newtree, newlocksroot, err := fromState.TryRemoveHashLock(lockSecretHash, blockNumber, mustExpired)
	if err != nil {
		return err
	}
	/*
		locksroot必须一致.
	*/
	if newlocksroot != msg.Locksroot {
		return &InvalidLocksRootError{ExpectedLocksroot: newlocksroot, GotLocksroot: msg.Locksroot}
	}
	fromState.Tree = newtree
	err = fromState.registerRemoveLock(messager, lockSecretHash)
	if err == nil {
		c.ExternState.db.RemoveLock(c.ChannelIdentifier.ChannelIdentifier, fromState.Address, lockSecretHash)
	}
	return err
}

func (c *Channel) isValidEnvelopMessage(evMsg *encoding.EnvelopMessage) bool {
	return evMsg.ChannelIdentifier == c.ChannelIdentifier.ChannelIdentifier &&
		evMsg.OpenBlockNumber == c.ChannelIdentifier.OpenBlockNumber
}

func (c *Channel) isChannelIdentifierValid(id *contracts.ChannelUniqueID) bool {
	return c.ChannelIdentifier.ChannelIdentifier == id.ChannelIdentifier &&
		c.ChannelIdentifier.OpenBlockNumber == c.ChannelIdentifier.OpenBlockNumber
}

//GetNextNonce change nonce  means banlance proof state changed
func (c *Channel) GetNextNonce() uint64 {
	if c.OurState.nonce() != 0 {
		return c.OurState.nonce() + 1
	}
	// 0 must not be used since in the netting contract it represents null.
	return 1
}

/*
CreateDirectTransfer return a DirectTransfer message.

This message needs to be signed and registered with the channel before
sent.
*/
func (c *Channel) CreateDirectTransfer(amount *big.Int) (tr *encoding.DirectTransfer, err error) {
	if !c.CanTransfer() {
		return nil, fmt.Errorf("transfer not possible, no funding or channel closed")
	}
	from := c.OurState
	to := c.PartnerState
	distributable := from.Distributable(to)
	if amount.Cmp(utils.BigInt0) <= 0 || amount.Cmp(distributable) > 0 {
		log.Debug(fmt.Sprintf("Insufficient funds : amount=%s, Distributable=%s", amount, distributable))
		return nil, rerr.ErrInsufficientFunds
	}
	transferAmount := new(big.Int).Add(from.TransferAmount(), amount)
	currentLocksroot := to.Tree.MerkleRoot()
	nonce := c.GetNextNonce()
	bp := encoding.NewBalanceProof(nonce, transferAmount, currentLocksroot, &c.ChannelIdentifier)
	tr = encoding.NewDirectTransfer(bp)
	return
}

/*
CreateMediatedTransfer return a MediatedTransfer message.

This message needs to be signed and registered with the channel before
sent.

Args:
    initiator : The node that requested the transfer.
    target : The final destination node of the transfer
    amount : How much of a token is being transferred.
    expiration : The maximum block number until the transfer
        message can be received.
	fee: 手续费
*/
func (c *Channel) CreateMediatedTransfer(initiator, target common.Address, fee *big.Int, amount *big.Int, expiration int64, lockSecretHash common.Hash) (tr *encoding.MediatedTransfer, err error) {
	if !c.CanTransfer() {
		return nil, fmt.Errorf("transfer not possible, no funding or channel closed")
	}
	if amount.Cmp(utils.BigInt0) <= 0 || amount.Cmp(c.Distributable()) > 0 {
		log.Info(fmt.Sprintf("Insufficient funds  amount=%s,Distributable=%s", amount, c.Distributable()))
		return nil, fmt.Errorf("insufficient funds")
	}
	from := c.OurState
	lock := &mtree.Lock{
		Amount:         amount,
		Expiration:     expiration,
		LockSecretHash: lockSecretHash,
	}
	_, updatedLocksroot := from.computeMerkleRootWith(lock)
	transferAmount := from.TransferAmount()
	nonce := c.GetNextNonce()
	bp := encoding.NewBalanceProof(nonce, transferAmount, updatedLocksroot, &c.ChannelIdentifier)
	tr = encoding.NewMediatedTransfer(bp, lock, target, initiator, fee)
	return
}

//CreateUnlock creates  a unlock message
func (c *Channel) CreateUnlock(lockSecretHash common.Hash) (tr *encoding.UnLock, err error) {
	from := c.OurState
	lock, secret, err := from.getSecretByLockSecretHash(lockSecretHash)
	if err != nil {
		return nil, fmt.Errorf("no such lock for lockSecretHash:%s", utils.HPex(lockSecretHash))
	}
	_, locksrootWithPendingLockRemoved, err := from.computeMerkleRootWithout(lock)
	if err != nil {
		return
	}
	transferAmount := new(big.Int).Add(from.TransferAmount(), lock.Amount)
	nonce := c.GetNextNonce()
	bp := encoding.NewBalanceProof(nonce, transferAmount, locksrootWithPendingLockRemoved, &c.ChannelIdentifier)
	tr = encoding.NewUnlock(bp, secret)
	return
}

/*
CreateRemoveExpiredHashLockTransfer create this transfer to notify my patner that this hashlock is expired and i want to remove it .
*/
func (c *Channel) CreateRemoveExpiredHashLockTransfer(lockSecretHash common.Hash, blockNumber int64) (tr *encoding.RemoveExpiredHashlockTransfer, err error) {
	_, _, newlocksroot, err := c.OurState.TryRemoveHashLock(lockSecretHash, blockNumber, true)
	if err != nil {
		return
	}
	nonce := c.GetNextNonce()
	transferAmount := c.OurState.TransferAmount()
	bp := encoding.NewBalanceProof(nonce, transferAmount, newlocksroot, &c.ChannelIdentifier)
	tr = encoding.NewRemoveExpiredHashlockTransfer(bp, lockSecretHash)
	return
}

/*
CreateAnnounceDisposedResponse 必须先收到对方的AnnouceDisposedTransfer, 然后才能移除.
*/
/*
 *	CreateAnnounceDisposedResponse : function to create message of AnnounceDisposedResponse.
 *	Note that a channel participant must first receive AnnounceDisposedTransfer, then he can
 */
func (c *Channel) CreateAnnounceDisposedResponse(lockSecretHash common.Hash, blockNumber int64) (tr *encoding.AnnounceDisposedResponse, err error) {
	_, _, newlocksroot, err := c.OurState.TryRemoveHashLock(lockSecretHash, blockNumber, false)
	if err != nil {
		return
	}
	nonce := c.GetNextNonce()
	transferAmount := c.OurState.TransferAmount()
	bp := encoding.NewBalanceProof(nonce, transferAmount, newlocksroot, &c.ChannelIdentifier)
	tr = encoding.NewAnnounceDisposedResponse(bp, lockSecretHash)
	return
}

/*
CreateAnnouceDisposed  声明我放弃收到的某个锁
*/
/*
 *	CreateAnnouceDisposed : function to create message of AnnounceDisposed
 *	Note that it claims that I have abandoned a lock.
 */
func (c *Channel) CreateAnnouceDisposed(lockSecretHash common.Hash, blockNumber int64) (tr *encoding.AnnounceDisposed, err error) {
	lock, _, _, err := c.PartnerState.TryRemoveHashLock(lockSecretHash, blockNumber, false)
	if err != nil {
		return
	}
	rp := &encoding.AnnounceDisposedProof{
		Lock: lock,
	}
	rp.ChannelIdentifier = c.ChannelIdentifier.ChannelIdentifier
	rp.OpenBlockNumber = c.ChannelIdentifier.OpenBlockNumber
	tr = encoding.NewAnnounceDisposed(rp)
	return
}

//ErrWithdrawButHasLocks 不能在有锁的情况下发起 withdraw 请求
/*
 *	ErrWithdrawButHasLocks : we can't send a request for withdraw when there are locks.
 */
var ErrWithdrawButHasLocks = errors.New("cannot withdraw when has lock")

//ErrSettleButHasLocks 不能在有锁的情况下发起 settle 请求
/*
 *	ErrSettleButHasLocks : we can't send a request for settle when there are locks.
 */
var ErrSettleButHasLocks = errors.New("cannot cooperative settle when has lock")

var errInvalidChannelIdentifier = errors.New("channel identifier is invalid")
var errInvalidSender = errors.New("messager's sender is not a participant of channel")
var errParticipant = errors.New("participant error")
var errBalance = errors.New("balance not match")

func (c *Channel) preCheckChannelID(tr encoding.SignedMessager, id *encoding.ChannelIDInMessage) error {
	if c.ChannelIdentifier.ChannelIdentifier != id.ChannelIdentifier ||
		c.ChannelIdentifier.OpenBlockNumber != id.OpenBlockNumber {
		return errInvalidChannelIdentifier
	}
	if tr.GetSender() != c.OurState.Address && tr.GetSender() != c.PartnerState.Address {
		return errInvalidSender
	}
	return nil
}

/*
RegisterAnnouceDisposed 收到对方的 AnnouceDisposed 消息
签名验证已经进行过了.
*/
/*
 *	RegisterAnnouceDisposed : function to register message of AnnounceDisposed.
 *  Note that signature verification has been undergone.
 */
func (c *Channel) RegisterAnnouceDisposed(tr *encoding.AnnounceDisposed) (err error) {
	err = c.preCheckChannelID(tr, &tr.ChannelIDInMessage)
	if err != nil {
		return
	}
	var state = c.PartnerState
	if tr.GetSender() == c.PartnerState.Address {
		state = c.OurState
	}
	mlock := tr.Lock
	lock := state.getLockByHashlock(mlock.LockSecretHash)
	if lock == nil || mlock.LockSecretHash != lock.LockSecretHash ||
		mlock.Expiration != lock.Expiration ||
		mlock.Amount.Cmp(lock.Amount) != 0 {
		return fmt.Errorf("RegisterAnnouceDisposed lock not match,receive=%s, mine=%s", mlock, lock)
	}
	return nil
}

/*
CreateWithdrawRequest 一定要不持有任何锁,否则双方可能对金额分配有争议.
*/
/*
 *	CreateWithdrawRequest : function to create message of request withdraw.
 *	Note that there must not be any lock, or conflict will reside in token allocation.
 */
func (c *Channel) CreateWithdrawRequest(withdrawAmount *big.Int) (w *encoding.WithdrawRequest, err error) {
	/*
		withdraw 一旦发出去就只能关闭通道
		无论是通过 withdraw 成功,造成通道关闭重开
		还是自己主动发起 close/settle.
		所以只要有一方持有锁,对于通道金额有争议,都不能发起 withdraw
	*/
	if len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2UnclaimedLocks) > 0 {
		err = ErrWithdrawButHasLocks
	}
	d := new(encoding.WithdrawRequestData)
	d.ChannelIdentifier = c.ChannelIdentifier.ChannelIdentifier
	d.OpenBlockNumber = c.ChannelIdentifier.OpenBlockNumber
	d.Participant1 = c.OurState.Address
	d.Participant2 = c.PartnerState.Address
	d.Participant1Balance = c.OurState.Balance(c.PartnerState)
	d.Participant1Withdraw = withdrawAmount
	if withdrawAmount.Cmp(d.Participant1Balance) > 0 {
		err = fmt.Errorf("withdraw amount too large,current=%s,withdraw=%s", w.Participant1Balance, withdrawAmount)
		return
	}
	w = encoding.NewWithdrawRequest(d)
	return
}

func (c *Channel) preCheckSettleDataInMessage(tr encoding.SignedMessager, sd *encoding.SettleDataInMessage) (err error) {
	if c.ChannelIdentifier.ChannelIdentifier != sd.ChannelIdentifier ||
		c.ChannelIdentifier.OpenBlockNumber != sd.OpenBlockNumber {
		return errInvalidChannelIdentifier
	}
	var state1, state2 *EndState
	if tr.GetSender() == c.OurState.Address {
		state1 = c.OurState
		state2 = c.PartnerState
	} else if tr.GetSender() == c.PartnerState.Address {
		state1 = c.PartnerState
		state2 = c.OurState
	} else {
		return errInvalidSender
	}
	/*
		state1 ,state2和 participant1,participant2没有对应关系,需要自己找出来.
	*/
	if (state1.Address != sd.Participant1 && state1.Address != sd.Participant2) ||
		(state2.Address != sd.Participant1 && state2.Address != sd.Participant2) ||
		sd.Participant1 == sd.Participant2 {
		return errParticipant
	}
	if state1.Address == sd.Participant1 {
		if state1.Balance(state2).Cmp(sd.Participant1Balance) != 0 ||
			state2.Balance(state1).Cmp(sd.Participant2Balance) != 0 {
			return errBalance
		}
	} else {
		if state2.Balance(state1).Cmp(sd.Participant1Balance) != 0 ||
			state1.Balance(state2).Cmp(sd.Participant2Balance) != 0 {
			return errBalance
		}
	}

	return nil
}

func (c *Channel) hasAnyLock() bool {
	if len(c.PartnerState.Lock2UnclaimedLocks) > 0 ||
		len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2UnclaimedLocks) > 0 ||
		len(c.OurState.Lock2PendingLocks) > 0 {
		return true
	}
	return false
}

/*RegisterWithdrawRequest :
1. 验证信息准确
2. 通道状态要切换到StateWithdraw
*/
/*
 *	RegisterWithdrawRequest : function to register WithdrawRequest.
 *
 *		1. verify the information is correct.
 *		2. channel state must switch to StateWithdraw.
 */
func (c *Channel) RegisterWithdrawRequest(tr *encoding.WithdrawRequest) (err error) {
	if c.ChannelIdentifier.ChannelIdentifier != tr.ChannelIdentifier ||
		c.ChannelIdentifier.OpenBlockNumber != tr.OpenBlockNumber {
		return errInvalidChannelIdentifier
	}
	if tr.GetSender() != c.PartnerState.Address {
		return errInvalidSender
	}
	if c.PartnerState.Balance(c.OurState).Cmp(tr.Participant1Balance) != 0 {
		return errBalance
	}
	/*
		有可能在我收到 request 的前一刻,我正在发出一笔交易,
		如果我是中间节点,相当于我收到了 announce disposed 一样处理
		如果我是发起方,认为此交易立即失败.
	*/
	if len(c.PartnerState.Lock2UnclaimedLocks) > 0 ||
		len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2UnclaimedLocks) > 0 {
		return errors.New("cannot withdraw when has unlock")
	}
	c.State = channeltype.StateWithdraw
	return nil
}

//HasAnyUnkonwnSecretTransferOnRoad 是否还有任何我发出的交易,并且对方不知道密码的
/*
 *	HasAnyUnknownSecretTransferOnRoad : function to check whether there is any transfer sent out from me
 * 		that my partner has no idea about the secret.
 */
func (c *Channel) HasAnyUnkonwnSecretTransferOnRoad() bool {
	return len(c.OurState.Lock2PendingLocks) > 0
}

/*
CreateWithdrawResponse :
我已经验证过了,对方的 withdrawRequest 是合理,可以接受的,
这里只是构建数据就可以了.
有可能在我收到对方 withdrawRequest 过程中,我在发起一笔交易,
当然这笔交易会失败,因为对方肯定不会接受.,就算对方接受了,也没有任何意义.不可能拿到此笔钱
所以 withdraw 和 cooperative settle都会影响到现在正在进行的交易,这些 statemanager 也需要处理.
*/
/*
 *	CreateWithdrawResponse : function to create message of WithdrawResponse.
 *
 *	Note that there is possibilities that I send out another transfer when receiving `withdrawRequest` from my partner.
 * 	With no doubt that this transfer will fail because my partner has no chance to accept it. Even he accepts it, he still
 * 	can not get the token.
 * 	So withdraw and cooperative settle may both impact ongoing transfers which statemanager should deal with.
 */
func (c *Channel) CreateWithdrawResponse(req *encoding.WithdrawRequest) (w *encoding.WithdrawResponse, err error) {
	if len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2PendingLocks) > 0 {
		log.Warn(fmt.Sprintf("CreateWithdrawResponse ,but i'm sending transfer on road,these transfer should canceled immediately"))
	}
	if len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2UnclaimedLocks) > 0 {
		panic("should no locks for partner state when  CreateWithdrawResponse")
	}
	wd := new(encoding.WithdrawReponseData)
	wd.ChannelIdentifier = c.ChannelIdentifier.ChannelIdentifier
	wd.OpenBlockNumber = c.ChannelIdentifier.OpenBlockNumber
	wd.Participant1 = c.PartnerState.Address
	wd.Participant2 = c.OurState.Address
	wd.Participant1Balance = c.PartnerState.Balance(c.OurState)
	wd.Participant1Withdraw = req.Participant1Withdraw
	w = encoding.NewWithdrawResponse(wd)
	/*
		再次验证信息正确性,
	*/
	// re-verify message to ensure correctness.
	if req.Participant1Balance.Cmp(w.Participant1Balance) != 0 {
		panic(fmt.Sprintf("withdrawequest=%s,\nwithdrawresponse=%s", req, w))
	}
	return
}

//RegisterWithdrawResponse check withdraw response
//外部应该验证响应与请求是一致的
/*
 *	RegisterWithdrawResponse : function to check withdraw response.
 *
 *	Explicit verify that withdraw response should be consistent with withdraw request.
 */
func (c *Channel) RegisterWithdrawResponse(tr *encoding.WithdrawResponse) error {
	if c.ChannelIdentifier.ChannelIdentifier != tr.ChannelIdentifier ||
		c.ChannelIdentifier.OpenBlockNumber != tr.OpenBlockNumber {
		return errInvalidChannelIdentifier
	}
	if tr.GetSender() != c.PartnerState.Address {
		return errInvalidSender
	}
	if c.OurState.Balance(c.PartnerState).Cmp(tr.Participant1Balance) != 0 {
		return errBalance
	}
	if len(c.PartnerState.Lock2UnclaimedLocks) > 0 ||
		len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2UnclaimedLocks) > 0 ||
		len(c.OurState.Lock2PendingLocks) > 0 {
		return errors.New("cannot withdraw when has unlock")
	}
	c.State = channeltype.StateWithdraw
	return nil
}

/*
CreateCooperativeSettleRequest 一定要不持有任何锁,否则双方可能对金额分配有争议.
*/
/*
 *	CreateCooperativeSettleRequest : function to create message of CooperativeSettleRequest.
 *	Note that there should be no lock, or both participants may have conflict with token allocation.
 */
func (c *Channel) CreateCooperativeSettleRequest() (s *encoding.SettleRequest, err error) {
	/*
		SettleRequest 一旦发出去就只能关闭通道
		无论是通过 cooperative settle 成功,造成通道关闭重开
		还是自己主动发起 close/settle.
		所以只要有一方持有锁,对于通道金额有争议,都不能发起 cooperative settle
	*/
	/*
	 *	Once SettleRequest sent out, channel has to be closed.
	 *	Channel reopens after being closed, via cooperative settle,
	 *	or participant send close/settle.
	 *	No matter which is the case, if one participant holds locks and has dispute about token amount,
	 *	they can not do cooperativesettle.
	 */
	if len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2UnclaimedLocks) > 0 {
		err = ErrWithdrawButHasLocks
	}
	wd := new(encoding.SettleRequestData)
	wd.ChannelIdentifier = c.ChannelIdentifier.ChannelIdentifier
	wd.OpenBlockNumber = c.ChannelIdentifier.OpenBlockNumber
	wd.Participant1 = c.OurState.Address
	wd.Participant2 = c.PartnerState.Address
	wd.Participant1Balance = c.OurState.Balance(c.PartnerState)
	wd.Participant2Balance = c.PartnerState.Balance(c.OurState)
	s = encoding.NewSettleRequest(wd)
	return
}

//RegisterCooperativeSettleRequest check settle request and update state
func (c *Channel) RegisterCooperativeSettleRequest(msg *encoding.SettleRequest) error {
	err := c.preCheckSettleDataInMessage(msg, &msg.SettleDataInMessage)
	if err != nil {
		return err
	}
	/*
		不能持有任何锁,除了在收到 settle request 前一刻,我正在发出交易
		如果我是交易发起方,认为交易理解失败
		如果我是交易的中间节点,就相当于收到了对方的 annouce disposed 一样处理.
		这需要我保存 settle request,如果 cooperative settle 失败怎么处理呢?!!
	*/
	/*
	 *	Can't hold any lock, except that I am sending out transfers right before settlerequest.
	 *	If I am the transfer initiator, assume sending transfer fails.
	 *	If I am the mediator, then handle transfer like announce disposed event.
	 *	which needs settle request, if cooperative settle failed, then how to deal with that?
	 */
	if len(c.PartnerState.Lock2UnclaimedLocks) > 0 ||
		len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2UnclaimedLocks) > 0 {
		return errors.New("cannot cooperative settle when has unlock")
	}
	c.State = channeltype.StateCooprativeSettle
	return nil
}

/*
CreateCooperativeSettleResponse :
我已经验证过了,对方的 settleRequest 是合理,可以接受的,
这里只是构建数据就可以了.
有可能在我收到对方 settleRequest 过程中,我在发起一笔交易,
当然这笔交易会失败,因为对方肯定不会接受.,就算对方接受了,也没有任何意义.不可能拿到此笔钱
所以 withdraw 和 cooperative settle都会影响到现在正在进行的交易,这些 statemanager 也需要处理.
*/
/*
 *	CreateCooperativeSettleResponse : function to create message of CooperativeSettleResponse.
 *	Note that a channel participant may send out another transfer to his partner, while receiving partner's settleRequest.
 *	With no doubt that this new transfer will fail because his channel partner has no chance to accept it.
 *	Even he accepts it, he cannot get that token.
 * 	So withdraw and cooperative settle may both impact ongoing transfers, which statemanager should handle.
 */
func (c *Channel) CreateCooperativeSettleResponse(req *encoding.SettleRequest) (res *encoding.SettleResponse, err error) {
	if len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2PendingLocks) > 0 {
		log.Warn(fmt.Sprintf("CreateCooperativeSettleResponse ,but i'm sending transfer on road,these transfer should canceled immediately"))
	}
	if len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2UnclaimedLocks) > 0 {
		panic("should no locks for partner state when  CreateWithdrawResponse")
	}
	d := new(encoding.SettleResponseData)
	d.ChannelIdentifier = c.ChannelIdentifier.ChannelIdentifier
	d.OpenBlockNumber = c.ChannelIdentifier.OpenBlockNumber
	d.Participant2 = c.OurState.Address
	d.Participant2Balance = c.OurState.Balance(c.PartnerState)
	d.Participant1 = c.PartnerState.Address
	d.Participant1Balance = c.PartnerState.Balance(c.OurState)

	res = encoding.NewSettleResponse(d)
	/*
		再次验证信息正确性,
	*/
	// Re-verify message correctness.
	if req.Participant1Balance.Cmp(d.Participant1Balance) != 0 ||
		req.Participant2Balance.Cmp(d.Participant2Balance) != 0 {
		panic(fmt.Sprintf("settle request=%s,\n settle re=%s", req, res))
	}
	return
}

//RegisterCooperativeSettleResponse check settle response and update state
func (c *Channel) RegisterCooperativeSettleResponse(msg *encoding.SettleResponse) error {
	err := c.preCheckSettleDataInMessage(msg, &msg.SettleDataInMessage)
	if err != nil {
		return err
	}
	c.State = channeltype.StateCooprativeSettle
	return nil
}

/*
PrepareForWithdraw :
由于 withdraw 和 合作settle 需要事先没有任何锁,因此必须先标记不进行任何交易
等现有交易完成以后再
*/
/*
 *	PrepareForWithdraw : function to change channel state to StatePrepareForWithdraw
 *	Note that because withdraw and cooperative settle require no lock,
 *	hence we should tag that any new transfer is forbidden, and after ongoing transfers finish,
 * 	we can do channel withdraw.
 */
func (c *Channel) PrepareForWithdraw() error {
	if c.State != channeltype.StateOpened {
		return fmt.Errorf("state must be opened when withdraw, but state is %s", c.State)
	}
	c.State = channeltype.StatePrepareForWithdraw
	return nil
}

/*
PrepareForCooperativeSettle :
由于 withdraw 和 合作settle 需要事先没有任何锁,因此必须先标记不进行任何交易
等现有交易完成以后再
*/
/*
 *	PrepareForCooperativeSettle : function to switch channel state to StatePrepareForCooperativeSettle.
 *
 *	Note that because withdraw and cooperative settle require no lock,
 *	hence we should tag that any new transfer is forbidden, and after ongoing transfers finish,
 * 	we can do channel withdraw.
 */
func (c *Channel) PrepareForCooperativeSettle() error {
	if c.State != channeltype.StateOpened {
		return fmt.Errorf("state must be opened when cooperative settle, but state is %s", c.State)
	}
	c.State = channeltype.StatePrepareForCooperativeSettle
	return nil
}

/*
CancelWithdrawOrCooperativeSettle 等待一段时间以后发现不能合作关闭通道,可以撤销
也可以直接选择调用 close
*/
/*
 *	CancelWithdrawCooperativeSettle : function to switch channel state to StateOpened.
 *
 *	Note that if we wait for some amount of time and found that we cannot cooperative settle, then we can cancel that
 *	Or directly invoke close.
 */
func (c *Channel) CancelWithdrawOrCooperativeSettle() error {
	if c.ExternState.ClosedBlock != 0 {
		return fmt.Errorf("no need cancel because of channel is closed")
	}
	if c.State != channeltype.StatePrepareForCooperativeSettle && c.State != channeltype.StatePrepareForWithdraw {
		return fmt.Errorf("state is %s,cannot cancel withdraw or cooperative", c.State)
	}
	c.State = channeltype.StateOpened
	return nil
}

/*
CanWithdrawOrCooperativeSettle 只有在任何锁的情况下才能进行 withdraw 和cooperative settle
*/
/*
 *	CanWithdrawOrCooperativeSettle : function to check whether we can process Withdraw / CooperativeSettle.
 *
 *	Note that we can do withdraw / CooperativeSettle without lock.
 */
func (c *Channel) CanWithdrawOrCooperativeSettle() bool {
	if len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.OurState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2PendingLocks) > 0 ||
		len(c.PartnerState.Lock2UnclaimedLocks) > 0 {
		return false
	}
	return true
}

//Close async close this channel
func (c *Channel) Close() (result *utils.AsyncResult) {
	if c.State != channeltype.StateOpened {
		log.Warn(fmt.Sprintf("try to close channel %s,but it's state is %s", utils.HPex(c.ChannelIdentifier.ChannelIdentifier), c.State))
	}
	if c.State == channeltype.StateClosed ||
		c.State == channeltype.StateSettled {
		result = utils.NewAsyncResult()
		result.Result <- fmt.Errorf("channel %s already closed or settled", utils.HPex(c.ChannelIdentifier.ChannelIdentifier))
		return
	}
	/*
		在关闭的过程中崩溃了,或者关闭 tx 失败了,这些都可能发生.所以不能因为 state 不对,就不允许 close
		标记的目的是为了阻止继续接受或者发起交易.
	*/
	/*
	 *	Things happen like crash while channel is closing, or failure when closing transaction.
	 *	We cannot forbid close just because channel state is abnormal.
	 *	State tag is used to prevent further receiving or sending transfers.
	 */
	c.State = channeltype.StateClosing
	bp := c.PartnerState.BalanceProofState
	result = c.ExternState.Close(bp)
	return
}

//Settle async settle this channel
func (c *Channel) Settle() (result *utils.AsyncResult) {
	if c.State != channeltype.StateClosed {
		return utils.NewAsyncResultWithError(fmt.Errorf("settle only valid when a channel is closed,now is %s", c.State))
	}
	//不需要修改状态, settle 失败以后还可以继续调用 settle.
	//c.State = channeltype.StateSettling
	var MyTransferAmount, PartnerTransferAmount *big.Int
	var MyLocksroot, PartnerLocksroot common.Hash
	if c.OurState.BalanceProofState != nil {
		MyTransferAmount = c.OurState.BalanceProofState.ContractTransferAmount
		MyLocksroot = c.OurState.BalanceProofState.ContractLocksRoot
	} else {
		MyTransferAmount = utils.BigInt0
	}
	if c.PartnerState.BalanceProofState != nil {
		PartnerTransferAmount = c.PartnerState.BalanceProofState.ContractTransferAmount
		PartnerLocksroot = c.PartnerState.BalanceProofState.ContractLocksRoot
	} else {
		PartnerTransferAmount = utils.BigInt0
	}
	return c.ExternState.Settle(MyTransferAmount, PartnerTransferAmount, MyLocksroot, PartnerLocksroot)
}

//GetNeedRegisterSecrets find all secres need to reveal on secret
func (c *Channel) GetNeedRegisterSecrets(blockNumber int64) (secrets []common.Hash) {
	for _, l := range c.PartnerState.Lock2UnclaimedLocks {
		if l.Lock.Expiration > blockNumber-int64(c.RevealTimeout) && l.Lock.Expiration < blockNumber {
			//底层负责处理重复的问题
			// lower layer takes charge of handling issues that repeatitively happen.
			secrets = append(secrets, l.Secret)
		}
	}
	return
}

/*
CooperativeSettleChannel 收到对方的 settle response, 关闭通道即可.
*/
/*
 *	CooperativeSettleChannel : function to undergo CooperativeSettle
 *
 *	Note that once a channel participant receives his partner's settle response, just close this channel.
 */
func (c *Channel) CooperativeSettleChannel(res *encoding.SettleResponse) (result *utils.AsyncResult) {
	w, err := c.CreateCooperativeSettleRequest()
	if err != nil {
		panic(err)
	}
	err = w.Sign(c.ExternState.privKey, w)
	if err != nil {
		panic(err)
	}
	return c.ExternState.TokenNetwork.CooperativeSettleAsync(
		res.Participant1, res.Participant2,
		res.Participant1Balance, res.Participant2Balance,
		w.Participant1Signature, res.Participant2Signature)
}

//CooperativeSettleChannelOnRequest 收到对方的 settle requet, 但是由于某些原因,需要我自己立即关闭通道
/*
 *	CooperativeSettleChannelOnRequest : function to handle channel cooperative channel request.
 *
 *	Note that this is case that a channel participant receives a cooperative settle request,
 *	but for some reasons that he has to close the channel immediately.
 */
func (c *Channel) CooperativeSettleChannelOnRequest(partnerSignature []byte, res *encoding.SettleResponse) (result *utils.AsyncResult) {
	return c.ExternState.TokenNetwork.CooperativeSettleAsync(
		res.Participant1, res.Participant2,
		res.Participant1Balance, res.Participant2Balance,
		partnerSignature, res.Participant2Signature,
	)
}

/*
Withdraw 收到对方的 withdraw response,
需要先验证参数有效
*/
/*
 *	Withdraw : function to undergo channel withdraw.
 *
 *	Note that this function has to work after verify parameter is valid.
 */
func (c *Channel) Withdraw(res *encoding.WithdrawResponse) (result *utils.AsyncResult) {
	//没有保存,需要重新签名.
	// No record, need to re-write signature.
	w, err := c.CreateWithdrawRequest(res.Participant1Withdraw)
	if err != nil {
		panic(err)
	}
	err = w.Sign(c.ExternState.privKey, w)
	if err != nil {
		panic(err)
	}
	return c.ExternState.TokenNetwork.WithdrawAsync(
		res.Participant1, res.Participant2,
		res.Participant1Balance, res.Participant1Withdraw,
		w.Participant1Signature, res.Participant2Signature,
	)
}

// String fmt.Stringer
func (c *Channel) String() string {
	return fmt.Sprintf("{ContractBalance=%s,Balance=%s,Distributable=%s,locked=%s,transferAmount=%s}",
		c.ContractBalance(), c.Balance(), c.Distributable(), c.Locked(), c.TransferAmount())
}

// NewChannelSerialization serialize the channel to save to database
func NewChannelSerialization(c *Channel) *channeltype.Serialization {
	var ourSecrets, partnerSecrets []common.Hash
	for _, s := range c.OurState.Lock2UnclaimedLocks {
		ourSecrets = append(ourSecrets, s.Secret)
	}
	for _, s := range c.PartnerState.Lock2UnclaimedLocks {
		partnerSecrets = append(partnerSecrets, s.Secret)
	}
	s := &channeltype.Serialization{
		Key:                    c.ChannelIdentifier.ChannelIdentifier[:],
		ChannelIdentifier:      &c.ChannelIdentifier,
		TokenAddressBytes:      c.TokenAddress[:],
		PartnerAddressBytes:    c.PartnerState.Address[:],
		OurAddress:             c.OurState.Address,
		RevealTimeout:          c.RevealTimeout,
		OurBalanceProof:        c.OurState.BalanceProofState,
		PartnerBalanceProof:    c.PartnerState.BalanceProofState,
		OurLeaves:              c.OurState.Tree.Leaves,
		PartnerLeaves:          c.PartnerState.Tree.Leaves,
		OurKnownSecrets:        ourSecrets,
		PartnerKnownSecrets:    partnerSecrets,
		State:                  c.State,
		SettleTimeout:          c.SettleTimeout,
		OurContractBalance:     c.OurState.ContractBalance,
		PartnerContractBalance: c.PartnerState.ContractBalance,
		ClosedBlock:            c.ExternState.ClosedBlock,
		SettledBlock:           c.ExternState.SettledBlock,
	}
	return s
}
