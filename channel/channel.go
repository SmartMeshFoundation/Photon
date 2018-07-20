package channel

import (
	"errors"

	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/fee"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
ch is the living representation of  channel on blockchain.
it contains all the transfers between two participants.
*/
type Channel struct {
	OurState             *EndState
	PartnerState         *EndState
	ExternState          *ExternalState
	ChannelIdentifier    contracts.ChannelUniqueID //this channel
	TokenAddress         common.Address
	RevealTimeout        int
	SettleTimeout        int
	ReceivedTransfers    []encoding.SignedMessager
	SentTransfers        []encoding.SignedMessager
	IsCloseEventComplete bool        //channel close event has been processed  completely  ,  crash when processing close event
	feeCharger           fee.Charger //calc fee for each transfer?
	State                channeltype.State
}

/*
NewChannel returns the living channel.
channelAddress must be a valid contract adress
settleTimeout must be valid, it cannot too small.
*/
func NewChannel(ourState, partenerState *EndState, externState *ExternalState, tokenAddr common.Address, channelIdentifier *contracts.ChannelUniqueID,
	revealTimeout, settleTimeout int) (c *Channel, err error) {
	if settleTimeout <= revealTimeout {
		err = errors.New("reveal_timeout can not be larger-or-equal to settle_timeout")
		return
	}
	if revealTimeout < 3 {
		/*
						 To guarantee that tokens won't be lost the expiration needs to
			             decrease at each hop, this is what forces the next hop to reveal
			             the secret with enough time for this node to unlock the lock with
			             the previous.

			             This /should be/ at least:

			               reveal_timeout = blocks_to_learn + blocks_to_mine * 2

			             Where:

			             - `blocks_to_learn` is the estimated worst case for a given block
			             to propagate to the full network. This is the time to learn a
			             secret revealed throught the blockchain.
			             - `blocks_to_mine * 2` is the estimated worst case for a given
			             transfer to be included in a block. This is the time to close a
			             channel and then to unlock a lock on chain.

		*/
		err = errors.New("reveal_timeout must be at least 3")
		return
	}
	c = &Channel{
		OurState:          ourState,
		PartnerState:      partenerState,
		ExternState:       externState,
		ChannelIdentifier: *channelIdentifier,
		TokenAddress:      tokenAddr,
		RevealTimeout:     revealTimeout,
		SettleTimeout:     settleTimeout,
		State:             channeltype.StateOpened,
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
	return channeltype.CanTransferMap[c.State] && c.Distributable().Cmp(utils.BigInt0) > 0
}

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
HandleClosed handles this channel was closed on blockchain
*/
func (c *Channel) HandleClosed(blockNumber int64, closingAddress common.Address) {
	balanceProof := c.PartnerState.BalanceProofState
	//the channel was closed, update our half of the state if we need to
	if closingAddress != c.OurState.Address {
		c.ExternState.UpdateTransfer(balanceProof)
	}
	unlockProofs := c.PartnerState.GetKnownUnlocks()
	if len(unlockProofs) > 0 {
		result := c.ExternState.Unlock(unlockProofs, c.PartnerState.TransferAmount())
		go func() {
			err := <-result.Result
			if err != nil {
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
	hashlock := utils.Sha3(secret[:])
	ourKnown := c.OurState.IsKnown(hashlock)
	partenerKnown := c.PartnerState.IsKnown(hashlock)
	if !ourKnown && !partenerKnown {
		return fmt.Errorf("secret doesn't correspond to a registered hashlock. hashlock %s token %s",
			utils.Pex(hashlock[:]), utils.HPex(c.ChannelIdentifier.ChannelIdentifier))
	}
	if ourKnown {
		lock := c.OurState.getLockByHashlock(hashlock)
		log.Debug(fmt.Sprintf("secret registered node=%s,from=%s,to=%s,token=%s,hashlock=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.OurState.Address[:]),
			utils.Pex(c.PartnerState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(hashlock[:]), lock.Amount))
		err := c.OurState.RegisterSecret(secret)
		return err
	}
	if partenerKnown {
		lock := c.PartnerState.getLockByHashlock(hashlock)
		log.Debug(fmt.Sprintf("secret registered node=%s,from=%s,to=%s,token=%s,hashlock=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.PartnerState.Address[:]),
			utils.Pex(c.OurState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(hashlock[:]), lock.Amount))
		return c.PartnerState.RegisterSecret(secret)
	}
	return nil
}

//链上对应的密码注册了
func (c *Channel) RegisterRevealedSecretHash(lockSecretHash common.Hash, blockNumber int64) error {
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
		err := c.OurState.RegisterRevealedSecretHash(lockSecretHash, blockNumber)
		if err == nil {
			//todo 需要发送给对方 unlock 消息,在哪里发比较合适呢? stateManager 还是这里?
		}
		return err
	}
	if partenerKnown {
		lock := c.PartnerState.getLockByHashlock(lockSecretHash)
		log.Debug(fmt.Sprintf("lockSecretHash registered node=%s,from=%s,to=%s,token=%s,lockSecretHash=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.PartnerState.Address[:]),
			utils.Pex(c.OurState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(lockSecretHash[:]), lock.Amount))
		return c.PartnerState.RegisterRevealedSecretHash(lockSecretHash, blockNumber)
	}
	return nil
}

//RegisterTransfer register a signed transfer, updating the channel's state accordingly.
//这些消息会改变 channel 的balance Proof
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
		err = fmt.Errorf("ch address mismatch,expect=%s,got=%s", c.ChannelIdentifier, evMsg)
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
	*/
	isInvalidNonce := (evMsg.Nonce < 1 || (fromState.nonce() != 0 && evMsg.Nonce != fromState.nonce()+1))
	//If a node data is damaged, then the channel will not work, so the data must not be damaged.
	if isInvalidNonce {
		//c may occur on normal operation
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
		err = fmt.Errorf("Negative transfer")
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
func (c *Channel) registerDirectTransfer(tr *encoding.DirectTransfer, blockNumber int64) (err error) {
	fromState, toState, err := c.PreCheckRecievedTransfer(tr)
	if err != nil {
		return
	}
	/*
		这次转账金额是多少
	*/
	amount := new(big.Int).Set(tr.TransferAmount)
	amount = amount.Sub(amount, fromState.TransferAmount())
	/*
		转账金额是负数或者超过了可以给的金额,都是错的
	*/
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
func (c *Channel) registerMediatedTranser(tr *encoding.MediatedTransfer, blockNumber int64) (err error) {
	fromState, toState, err := c.PreCheckRecievedTransfer(tr)
	if err != nil {
		return
	}
	/*
		这次转账金额是多少
	*/
	amount := tr.PaymentAmount
	/*
		转账金额是负数或者超过了可以给的金额,都是错的
	*/
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
	isSender := tr.Sender == c.OurState.Address
	if isSender && expiresAfterSettle { //After receiving this lock, the party can close or updatetransfer on the chain, so that if the party does not have a password, he still can't get the money.
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
func (c *Channel) RegisterAnnounceDisposedResponse(response *encoding.AnnounceDisposedResponse, blockNumber int64) (err error) {
	return c.registerRemoveLock(response, blockNumber, response.LockSecretHash, false)
}
func (c *Channel) registerRemoveLock(messager encoding.EnvelopMessager, blockNumber int64, locksSecretHash common.Hash, mustExpired bool) (err error) {
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
	_, newtree, newlocksroot, err := fromState.TryRemoveHashLock(locksSecretHash, blockNumber, mustExpired)
	if err != nil {
		return err
	}
	/*
		locksroot必须一致.
	*/
	if newlocksroot != msg.Locksroot {
		return &InvalidLocksRootError{ExpectedLocksroot: newlocksroot, GotLocksroot: msg.Locksroot}
	}
	fromState.tree = newtree
	err = fromState.registerRemoveLock(messager, locksSecretHash)
	//if err == nil {
	//	c.ExternState.db.RemoveLock(c.ChannelIdentifier, fromState.Address, tr.LockSecretHash)
	//}
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
func (c *Channel) GetNextNonce() int64 {
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
	currentLocksroot := to.tree.MerkleRoot()
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
		return nil, fmt.Errorf("Insufficient funds")
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
必须先收到对方的AnnouceDisposedTransfer, 然后才能移除.
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
var ErrWithdrawButHasLocks = errors.New("cannot withdraw when has lock")

//ErrSettleButHasLocks 不能在有锁的情况下发起 settle 请求
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
	d.Participant1Balance = c.OurState.Balance(c.PartnerState)
	d.Participant2 = c.PartnerState.Address
	d.Participant2Balance = c.PartnerState.Balance(c.OurState)
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

/*
RegisterWithdrawRequest
1. 验证信息准确
2. 通道状态要切换到StateWithdraw

*/
func (c *Channel) RegisterWithdrawRequest(tr *encoding.WithdrawRequest) (err error) {
	err = c.preCheckSettleDataInMessage(tr, &tr.SettleDataInMessage)
	if err != nil {
		return
	}
	c.State = channeltype.StateWithdraw
	return nil
}

/*
我已经验证过了,对方的 withdrawRequest 是合理,可以接受的,
这里只是构建数据就可以了.
有可能在我收到对方 withdrawRequest 过程中,我在发起一笔交易,
当然这笔交易会失败,因为对方肯定不会接受.,就算对方接受了,也没有任何意义.不可能拿到此笔钱
所以 withdraw 和 cooperative settle都会影响到现在正在进行的交易,这些 statemanager 也需要处理.
*/
func (c *Channel) CreateWithdrawResponse(req *encoding.WithdrawRequest, withdrawAmount *big.Int) (w *encoding.WithdrawResponse, err error) {
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
	wd.Participant2 = c.OurState.Address
	wd.Participant1 = c.PartnerState.Address
	wd.Participant2Balance = c.OurState.Balance(c.PartnerState)
	wd.Participant1Balance = c.PartnerState.Balance(c.OurState)
	wd.Participant1Withdraw = req.Participant1Withdraw
	wd.Participant2Withdraw = withdrawAmount
	if withdrawAmount.Cmp(wd.Participant2Balance) > 0 {
		err = fmt.Errorf("withdraw amount too large,current=%s,withdraw=%s", w.Participant2Balance, withdrawAmount)
		return
	}
	w = encoding.NewWithdrawResponse(wd)
	/*
		再次验证信息正确性,
	*/
	if req.Participant1Balance.Cmp(w.Participant1Balance) != 0 ||
		req.Participant2Balance.Cmp(w.Participant2Balance) != 0 {
		panic(fmt.Sprintf("withdrawequest=%s,\nwithdrawresponse=%s", req, w))
	}
	return
}
func (c *Channel) RegisterWithdrawResponse(tr *encoding.WithdrawResponse) error {
	err := c.preCheckSettleDataInMessage(tr, &tr.SettleDataInMessage)
	if err != nil {
		return err
	}
	c.State = channeltype.StateWithdraw
	return nil
}

/*
CreateCooperativeSettleRequest 一定要不持有任何锁,否则双方可能对金额分配有争议.
*/
func (c *Channel) CreateCooperativeSettleRequest() (s *encoding.SettleRequest, err error) {
	/*
		SettleRequest 一旦发出去就只能关闭通道
		无论是通过 cooperative settle 成功,造成通道关闭重开
		还是自己主动发起 close/settle.
		所以只要有一方持有锁,对于通道金额有争议,都不能发起 cooperative settle
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
func (c *Channel) RegisterCooperativeSettleRequest(msg *encoding.SettleRequest) error {
	err := c.preCheckSettleDataInMessage(msg, &msg.SettleDataInMessage)
	if err != nil {
		return err
	}
	c.State = channeltype.StateCooprativeSettle
	return nil
}

/*
我已经验证过了,对方的 settleRequest 是合理,可以接受的,
这里只是构建数据就可以了.
有可能在我收到对方 settleRequest 过程中,我在发起一笔交易,
当然这笔交易会失败,因为对方肯定不会接受.,就算对方接受了,也没有任何意义.不可能拿到此笔钱
所以 withdraw 和 cooperative settle都会影响到现在正在进行的交易,这些 statemanager 也需要处理.
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
	if req.Participant1Balance.Cmp(d.Participant1Balance) != 0 ||
		req.Participant2Balance.Cmp(d.Participant2Balance) != 0 {
		panic(fmt.Sprintf("settle request=%s,\n settle re=%s", req, res))
	}
	return
}
func (c *Channel) RegisterCooperativeSettleResponse(msg *encoding.SettleResponse) error {
	err := c.preCheckSettleDataInMessage(msg, &msg.SettleDataInMessage)
	if err != nil {
		return err
	}
	c.State = channeltype.StateCooprativeSettle
	return nil
}

/*
由于 withdraw 和 合作settle 需要事先没有任何锁,因此必须先标记不进行任何交易
等现有交易完成以后再
*/
func (c *Channel) PrepareForWithdraw() error {
	if c.State != channeltype.StateOpened {
		return fmt.Errorf("state must be opened when withdraw, but state is %s", c.State)
	}
	c.State = channeltype.StatePrepareForWithdraw
	return nil
}

/*
由于 withdraw 和 合作settle 需要事先没有任何锁,因此必须先标记不进行任何交易
等现有交易完成以后再
*/
func (c *Channel) PrepareForCooperativeSettle() error {
	if c.State != channeltype.StateOpened {
		return fmt.Errorf("state must be opened when cooperative settle, but state is %s", c.State)
	}
	c.State = channeltype.StatePrepareForSettle
	return nil
}

/*
等待一段时间以后发现不能合作关闭通道,可以撤销
也可以直接选择调用 close
*/
func (c *Channel) CancelWithdrawOrCooperativeSettle() error {
	if c.ExternState.ClosedBlock != 0 {
		return fmt.Errorf("no need cancel because of channel is closed")
	}
	if c.State != channeltype.StatePrepareForSettle && c.State != channeltype.StatePrepareForWithdraw {
		return fmt.Errorf("state is %s,cannot cancel withdraw or cooperative", c.State)
	}
	c.State = channeltype.StateOpened
	return nil
}

/*
只有在任何锁的情况下才能进行 withdraw 和cooperative settle
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
	c.State = channeltype.StateClosing
	bp := c.PartnerState.BalanceProofState
	result = c.ExternState.Close(bp)
	return
}
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

// String fmt.Stringer
func (c *Channel) String() string {
	return fmt.Sprintf("{ContractBalance=%s,Balance=%s,Distributable=%s,locked=%s,transferAmount=%s}",
		c.ContractBalance(), c.Balance(), c.Distributable(), c.Locked(), c.TransferAmount())
}

/*
todo 优化保存在数据库中的 channel 信息,不必要保存这么多
*/
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
		OurLeaves:              c.OurState.tree.Leaves,
		PartnerLeaves:          c.PartnerState.tree.Leaves,
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
