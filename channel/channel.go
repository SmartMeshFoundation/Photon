package channel

import (
	"errors"

	"fmt"

	"math/big"

	"encoding/gob"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/fee"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
Channel is the living representation of  channel on blockchain.
it contains all the transfers between two participants.
*/
type Channel struct {
	OurState             *EndState
	PartnerState         *EndState
	ExternState          *ExternalState
	TokenAddress         common.Address
	MyAddress            common.Address //this channel
	RevealTimeout        int
	SettleTimeout        int
	ReceivedTransfers    []encoding.SignedMessager
	SentTransfers        []encoding.SignedMessager
	IsCloseEventComplete bool        //channel close event has been processed  completely  ,  crash when processing close event
	feeCharger           fee.Charger //calc fee for each transfer?
}

/*
NewChannel returns the living channel.
channelAddress must be a valid contract adress
settleTimeout must be valid, it cannot too small.
*/
func NewChannel(ourState, partenerState *EndState, externState *ExternalState,
	tokenAddress, channelAddress common.Address, bcs *rpc.BlockChainService,
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
		OurState:      ourState,
		PartnerState:  partenerState,
		ExternState:   externState,
		TokenAddress:  tokenAddress,
		MyAddress:     channelAddress,
		RevealTimeout: revealTimeout,
		SettleTimeout: settleTimeout,
	}
	return
}

//State returns the state of this channel
func (c *Channel) State() string {
	if c.ExternState.SettledBlock != 0 {
		return transfer.ChannelStateSettled
	}
	if c.ExternState.ClosedBlock != 0 {
		return transfer.ChannelStateClosed
	}
	return transfer.ChannelStateOpened
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
	return c.State() == transfer.ChannelStateOpened && c.Distributable().Cmp(utils.BigInt0) > 0
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
	err := c.ExternState.WithDraw(unlockProofs)
	if err != nil {
		log.Error(fmt.Sprintf("withdraw on %s failed, channel is gone, error:%s", utils.APex(c.MyAddress), err))
	}
	c.IsCloseEventComplete = true
}

/*
HandleSettled handles this channel was settled on blockchain
there is nothing tod rightnow
*/
func (c *Channel) HandleSettled(blockNumber int64) {

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
		return fmt.Errorf("Secret doesn't correspond to a registered hashlock. hashlock %s token %s",
			utils.Pex(hashlock[:]), utils.Pex(c.TokenAddress[:]))
	}
	if ourKnown {
		lock := c.OurState.getLockByHashlock(hashlock)
		log.Debug(fmt.Sprintf("secret registered node=%s,from=%s,to=%s,token=%s,hashlock=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.OurState.Address[:]),
			utils.Pex(c.PartnerState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(hashlock[:]), lock.Amount))
		return c.OurState.RegisterSecret(secret)
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

/*
PreCheckRecievedTransfer pre check received message(directtransfer,mediatedtransfer,refundtransfer) is valid or not
*/
func (c *Channel) PreCheckRecievedTransfer(blockNumber int64, tr encoding.EnvelopMessager) (fromState *EndState, toState *EndState, err error) {
	evMsg := tr.GetEnvelopMessage()
	if evMsg.Channel != c.MyAddress {
		err = fmt.Errorf("Channel address mismatch")
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
RegisterRemoveExpiredHashlockTransfer register a request to remove a expired hashlock and this hashlock must be sent out from the sender.
*/
func (c *Channel) RegisterRemoveExpiredHashlockTransfer(tr *encoding.RemoveExpiredHashlockTransfer, blockNumber int64) (err error) {
	fromState, _, err := c.PreCheckRecievedTransfer(blockNumber, tr)
	if err != nil {
		return
	}
	/*
		transfer amount should not change.
	*/
	if tr.TransferAmount.Cmp(fromState.TransferAmount()) != 0 {
		err = errTransferAmountMismatch
		return
	}
	_, newtree, newlocksroot, err := fromState.TryRemoveExpiredHashLock(tr.HashLock, blockNumber)
	if err != nil {
		return err
	}
	/*
		only remove a expired hashlock
	*/
	if newlocksroot != tr.Locksroot {
		return &InvalidLocksRootError{ExpectedLocksroot: newlocksroot, GotLocksroot: tr.Locksroot}
	}
	fromState.TreeState = transfer.NewMerkleTreeState(newtree)
	err = fromState.registerRemoveExpiredHashlockTransfer(tr)
	if err == nil {
		c.ExternState.db.RemoveLock(c.MyAddress, fromState.Address, tr.HashLock)
	}
	return err
}

/*
CreateRemoveExpiredHashLockTransfer create this transfer to notify my patner that this hashlock is expired and i want to remove it .
*/
func (c *Channel) CreateRemoveExpiredHashLockTransfer(hashlock common.Hash, blockNumber int64) (tr *encoding.RemoveExpiredHashlockTransfer, err error) {
	_, _, newlocksroot, err := c.OurState.TryRemoveExpiredHashLock(hashlock, blockNumber)
	if err != nil {
		return
	}
	nonce := c.GetNextNonce()
	transferAmount := c.OurState.TransferAmount()
	tr = encoding.NewRemoveExpiredHashlockTransfer(0, nonce, c.MyAddress, transferAmount, newlocksroot, hashlock)
	return
}

//RegisterTransfer register a signed transfer, updating the channel's state accordingly.
func (c *Channel) RegisterTransfer(blocknumber int64, tr encoding.EnvelopMessager) error {
	var err error
	if tr.GetSender() == c.OurState.Address {
		err = c.RegisterTransferFromTo(blocknumber, tr, c.OurState, c.PartnerState)
		if err != nil {
			return err
		}
		c.SentTransfers = append(c.SentTransfers, tr)
		return nil
	} else if tr.GetSender() == c.PartnerState.Address {
		err = c.RegisterTransferFromTo(blocknumber, tr, c.PartnerState, c.OurState)
		if err != nil {
			return err
		}
		c.ReceivedTransfers = append(c.ReceivedTransfers, tr)
		return nil
	} else {
		log.Warn(fmt.Sprintf("Received a transfer from party that is not a part of the channel node=%s,from=%s, channel=%s",
			utils.Pex(c.OurState.Address[:]), utils.APex(tr.GetSender()), utils.APex(tr.GetEnvelopMessage().Channel)))
		return rerr.UnknownAddress(utils.StringInterface(tr, 3))
	}
}

/*
RegisterTransferFromTo Validates and register a signed transfer, updating the channel's state accordingly.
Note:
            The transfer must be registered before it is sent, not on
            acknowledgement. That is necessary for two reasons:

            - Guarantee that the transfer is valid.
            - Avoid sending a new transaction without funds.

        Raises:
            ErrInsufficientBalance: If the transfer is negative or above the Distributable amount.
            InvalidLocksRoot: If locksroot check fails.
            InvalidNonce: If the expected nonce does not match.
            ValueError: If there is an address mismatch (token or node address).
*/
func (c *Channel) RegisterTransferFromTo(blockNumber int64, tr encoding.EnvelopMessager, fromState *EndState, toState *EndState) error {
	var err error
	evMsg := tr.GetEnvelopMessage()
	if evMsg.Channel != c.MyAddress {
		return fmt.Errorf("Channel address mismatch")
	}
	if tr.GetSender() != fromState.Address {
		return fmt.Errorf("Unsigned transfer")
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
		return rerr.InvalidNonce(utils.StringInterface(tr, 3))
	}
	/*
					 if the locksroot is out-of-sync (because a transfer was created while
				    a Secret was in traffic) the Balance _will_ be wrong, so first check
				    the locksroot and then the Balance
		During building this transfer and registering transfer, we receive a secret.
	*/
	if encoding.IsLockedTransfer(tr) {
		mtr := encoding.GetMtrFromLockedTransfer(tr)
		lock := mtr.GetLock()
		if fromState.IsKnown(lock.HashLock) {
			//c may occur on normal operation
			log.Info(fmt.Sprintf("duplicated lock node=%s,from=%s,to=%s,hashlock=%s,received_locksroot=%s",
				utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]),
				utils.Pex(toState.Address[:]), utils.Pex(lock.HashLock[:]),
				utils.Pex(mtr.Locksroot[:])))
			return fmt.Errorf("hashlock is already registered")
		}
		/*
					  As a receiver: Check that all locked transfers are registered in
			            the locksroot, if any hashlock is missing there is no way to
			            claim it while the channel is closing
		*/
		_, expectedLocksroot := fromState.computeMerkleRootWith(mtr.GetLock())
		if expectedLocksroot != mtr.Locksroot {
			//c should not happen
			log.Warn(fmt.Sprintf("locksroot mismatch node=%s,from=%s,to=%s,hashlock=%s,expectedlocksroot=%s,receivedlocksroot=%s",
				utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
				utils.Pex(lock.HashLock[:]), utils.Pex(expectedLocksroot[:]), utils.Pex(mtr.Locksroot[:])))
			return &InvalidLocksRootError{expectedLocksroot, mtr.Locksroot}
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
		expiresAfterSettle := mtr.Expiration > endSettlePeriod
		isSender := mtr.Sender == c.OurState.Address
		if isSender && expiresAfterSettle { //After receiving this lock, the party can close or updatetransfer on the chain, so that if the party does not have a password, he still can't get the money.
			log.Error(fmt.Sprintf("Lock expires after the settlement period. node=%s,from=%s,to=%s,lockexpiration=%d,currentblock=%d,end_settle_period=%d",
				utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
				mtr.Expiration, blockNumber, endSettlePeriod))
			return fmt.Errorf("lock expires after the settlement period")
		}
	}
	// only check the Balance if the locksroot matched
	if evMsg.TransferAmount.Cmp(fromState.TransferAmount()) < 0 {
		log.Error(fmt.Sprintf("NEGATIVE TRANSFER node=%s,from=%s,to=%s,transfer=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
			utils.StringInterface(tr, 3))) //for nest struct
		return fmt.Errorf("Negative transfer")
	}
	amount := new(big.Int).Sub(evMsg.TransferAmount, fromState.TransferAmount())
	distributable := fromState.Distributable(toState)
	if tr.Cmd() == encoding.DirectTransferCmdID {
		if amount.Cmp(distributable) > 0 {
			return rerr.ErrInsufficientBalance
		}
	} else if encoding.IsLockedTransfer(tr) {
		mtr := encoding.GetMtrFromLockedTransfer(tr)
		if new(big.Int).Add(amount, mtr.Amount).Cmp(distributable) > 0 {
			return rerr.ErrInsufficientBalance
		}
	} else if tr.Cmd() == encoding.SecretCmdID {
		sec := tr.(*encoding.Secret)
		hashlock := utils.Sha3(sec.Secret[:])
		lock := fromState.getLockByHashlock(hashlock)
		if lock == nil {
			err = fmt.Errorf("channel %s receive secret message,but has no related hashlock,msg=%s", utils.APex(c.MyAddress), utils.StringInterface(sec, 3))
			log.Error(fmt.Sprintf("getLockByHashlock err %s", err))
			return err
		}
		transferAmount := new(big.Int).Add(fromState.TransferAmount(), lock.Amount)
		/*
			 tr.transferred_amount could be larger than the previous
				             transferred_amount + lock.amount, that scenario is a bug of the
				             payer
		*/
		if sec.TransferAmount.Cmp(transferAmount) != 0 {
			return fmt.Errorf("invalid transferred_amount, expected: %s got: %s",
				transferAmount, sec.TransferAmount)
		}
	}
	/*
			   all checks need to be done before the internal state of the channel
		         is changed, otherwise if a check fails and the state was changed the
		         channel will be left trashed
	*/
	if encoding.IsLockedTransfer(tr) {
		mtr := encoding.GetMtrFromLockedTransfer(tr)
		mroot := fromState.TreeState.Tree.MerkleRoot()
		log.Debug(fmt.Sprintf("REGISTERED LOCK node=%s,from=%s,to=%s,currentlocksroot=%s,lockamouont=%s,lock_expiration=%d,lock_hashlock=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
			utils.Pex(mroot[:]), mtr.Amount, mtr.Expiration, mtr.HashLock.String()))
		err = fromState.registerLockedTransfer(tr)
		if err != nil {
			return err
		}
		/*
			 register c channel as waiting for the secret (the secret can
			be revealed through a message or a blockchain log)
		*/
		c.ExternState.funcRegisterChannelForHashlock(c, mtr.HashLock)
	}
	if tr.Cmd() == encoding.DirectTransferCmdID {
		err = fromState.registerDirectTransfer(tr.(*encoding.DirectTransfer))
		if err != nil {
			return err
		}
	}
	if tr.Cmd() == encoding.SecretCmdID {
		err = fromState.registerSecretMessage(tr.(*encoding.Secret))
		if err != nil {
			return err
		}
	}
	mroot := fromState.TreeState.Tree.MerkleRoot()
	log.Debug(fmt.Sprintf("'REGISTERED TRANSFER node=%s,from=%s,to=%s,transfer_amount=%s,nonce=%d,current_locksroot=%s,\ntransfer=%s",
		utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
		fromState.TransferAmount(), fromState.nonce(), utils.Pex(mroot[:]), utils.StringInterface(tr, 3)))
	return nil
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
func (c *Channel) CreateDirectTransfer(amount *big.Int, identifier uint64) (tr *encoding.DirectTransfer, err error) {
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
	tranferAmount := new(big.Int).Add(from.TransferAmount(), amount)
	currentLocksroot := to.TreeState.Tree.MerkleRoot()
	nonce := c.GetNextNonce()
	tr = encoding.NewDirectTransfer(identifier, nonce, c.TokenAddress, c.MyAddress, tranferAmount, to.Address, currentLocksroot)
	return
}

/*
CreateMediatedTransfer return a MediatedTransfer message.

This message needs to be signed and registered with the channel before
sent.

Args:
    transfer_initiator (address): The node that requested the transfer.
    transfer_target (address): The final destination node of the transfer
    amount (float): How much of a token is being transferred.
    expiration (int): The maximum block number until the transfer
        message can be received.
*/
func (c *Channel) CreateMediatedTransfer(initiator, target common.Address, fee *big.Int, amount *big.Int, identifier uint64, expiration int64, hashlock common.Hash) (tr *encoding.MediatedTransfer, err error) {
	if !c.CanTransfer() {
		return nil, fmt.Errorf("transfer not possible, no funding or channel closed")
	}
	if amount.Cmp(utils.BigInt0) <= 0 || amount.Cmp(c.Distributable()) > 0 {
		log.Info(fmt.Sprintf("Insufficient funds  amount=%s,Distributable=%s", amount, c.Distributable()))
		return nil, fmt.Errorf("Insufficient funds")
	}
	from := c.OurState
	to := c.PartnerState
	lock := &encoding.Lock{
		Amount:     amount,
		Expiration: expiration,
		HashLock:   hashlock,
	}
	_, updatedLocksroot := from.computeMerkleRootWith(lock)
	transferAmount := from.TransferAmount()
	nonce := c.GetNextNonce()
	tr = encoding.NewMediatedTransfer(identifier, nonce, c.TokenAddress, c.MyAddress,
		transferAmount, to.Address, updatedLocksroot, lock, target, initiator, fee)
	return
}

/*
CreateRefundTransfer is similar as CreateMediatedTransfer
*/
func (c *Channel) CreateRefundTransfer(initiator, target common.Address, fee *big.Int, amount *big.Int, identifier uint64, expiration int64, hashlock common.Hash) (tr *encoding.RefundTransfer, err error) {
	mtr, err := c.CreateMediatedTransfer(initiator, target, fee, amount, identifier, expiration, hashlock)
	if err != nil {
		return
	}
	tr = encoding.NewRefundTransferFromMediatedTransfer(mtr)
	return
}

//CreateSecret creates  a secret message
func (c *Channel) CreateSecret(identifer uint64, secret common.Hash) (tr *encoding.Secret, err error) {
	hashlock := utils.Sha3(secret[:])
	from := c.OurState
	lock := from.getLockByHashlock(hashlock)
	if lock == nil {
		return nil, fmt.Errorf("no such lock for secret:%s", utils.HPex(secret))
	}
	_, locksrootWithPendingLockRemoved, err := from.computeMerkleRootWithout(lock)
	if err != nil {
		return
	}
	transferAmount := new(big.Int).Add(from.TransferAmount(), lock.Amount)
	nonce := c.GetNextNonce()
	tr = encoding.NewSecret(identifer, nonce, c.MyAddress, transferAmount, locksrootWithPendingLockRemoved, secret)
	return
}

// String fmt.Stringer
func (c *Channel) String() string {
	return fmt.Sprintf("{ContractBalance=%s,Balance=%s,Distributable=%s,locked=%s,transferAmount=%s}",
		c.ContractBalance(), c.Balance(), c.Distributable(), c.Locked(), c.TransferAmount())
}

// Serialization is the living channel in the database
type Serialization struct {
	ChannelAddress             common.Address
	ChannelAddressString       string `storm:"id"` //only for storm, because of save bug
	TokenAddress               common.Address
	PartnerAddress             common.Address
	TokenAddressString         string `storm:"index"`
	PartnerAddressString       string `storm:"index"`
	OurAddress                 common.Address
	RevealTimeout              int
	OurBalanceProof            *transfer.BalanceProofState
	PartnerBalanceProof        *transfer.BalanceProofState
	OurLeaves                  []common.Hash
	PartnerLeaves              []common.Hash
	OurLock2PendingLocks       map[common.Hash]PendingLock
	OurLock2UnclaimedLocks     map[common.Hash]UnlockPartialProof
	PartnerLock2PendingLocks   map[common.Hash]PendingLock
	PartnerLock2UnclaimedLocks map[common.Hash]UnlockPartialProof
	State                      string
	OurBalance                 *big.Int
	PartnerBalance             *big.Int
	OurContractBalance         *big.Int
	PartnerContractBalance     *big.Int
	OurAmountLocked            *big.Int
	PartnerAmountLocked        *big.Int
	ClosedBlock                int64
	SettledBlock               int64
	SettleTimeout              int
}

// NewChannelSerialization serialize the channel to save to database
func NewChannelSerialization(c *Channel) *Serialization {
	s := &Serialization{
		ChannelAddress:             c.MyAddress,
		ChannelAddressString:       c.MyAddress.String(),
		TokenAddress:               c.TokenAddress,
		TokenAddressString:         c.TokenAddress.String(),
		PartnerAddress:             c.PartnerState.Address,
		PartnerAddressString:       c.PartnerState.Address.String(),
		OurAddress:                 c.OurState.Address,
		RevealTimeout:              c.RevealTimeout,
		OurBalanceProof:            c.OurState.BalanceProofState,
		PartnerBalanceProof:        c.PartnerState.BalanceProofState,
		OurLeaves:                  c.OurState.TreeState.Tree.Layers[transfer.LayerLeaves],
		PartnerLeaves:              c.PartnerState.TreeState.Tree.Layers[transfer.LayerLeaves],
		OurLock2PendingLocks:       c.OurState.Lock2PendingLocks,
		OurLock2UnclaimedLocks:     c.OurState.Lock2UnclaimedLocks,
		PartnerLock2PendingLocks:   c.PartnerState.Lock2PendingLocks,
		PartnerLock2UnclaimedLocks: c.PartnerState.Lock2UnclaimedLocks,
		State:                  c.State(),
		OurBalance:             c.Balance(),
		PartnerBalance:         c.PartnerBalance(),
		SettleTimeout:          c.SettleTimeout,
		OurContractBalance:     c.ContractBalance(),
		PartnerContractBalance: c.PartnerState.ContractBalance,
		OurAmountLocked:        c.OurState.amountLocked(),
		PartnerAmountLocked:    c.PartnerState.amountLocked(),
		ClosedBlock:            c.ExternState.ClosedBlock,
		SettledBlock:           c.ExternState.SettledBlock,
	}
	return s
}
func init() {
	gob.Register(&Serialization{})
}
