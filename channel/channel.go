package channel

import (
	"errors"

	"fmt"

	"math/big"

	"encoding/gob"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/rerr"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type Channel struct {
	OurState             *ChannelEndState
	PartnerState         *ChannelEndState
	ExternState          *ChannelExternalState
	TokenAddress         common.Address
	MyAddress            common.Address //this channel
	RevealTimeout        int
	SettleTimeout        int
	ReceivedTransfers    []encoding.SignedMessager
	SentTransfers        []encoding.SignedMessager
	IsCloseEventComplete bool //channel close event has been processed  completely  ,  crash when processing close event
}

func NewChannel(ourState, partenerState *ChannelEndState, externState *ChannelExternalState,
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

func (c *Channel) State() string {
	if c.ExternState.SettledBlock != 0 {
		return transfer.CHANNEL_STATE_SETTLED
	}
	if c.ExternState.ClosedBlock != 0 {
		return transfer.CHANNEL_STATE_CLOSED
	}
	return transfer.CHANNEL_STATE_OPENED
}

/*
 Return the available amount of the token that our end of the
        channel can transfer to the partner.
*/
func (c *Channel) Distributable() *big.Int {
	return c.OurState.Distributable(c.PartnerState)
}
func (c *Channel) CanTransfer() bool {
	return c.State() == transfer.CHANNEL_STATE_OPENED && c.Distributable().Cmp(utils.BigInt0) > 0
}

//Return the total amount of token we deposited in the channel
func (c *Channel) ContractBalance() *big.Int {
	return c.OurState.ContractBalance
}

//Return how much we transferred to partner.
func (c *Channel) TransferAmount() *big.Int {
	return c.OurState.TransferAmount()
}

/*
		Return our current balance.

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
		Return partner current balance.

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
Return the current amount of our token that is locked waiting for a
        secret.

        The locked value is equal to locked transfers that have been
        initialized but their secret has not being revealed.
*/
func (c *Channel) Locked() *big.Int {
	return c.OurState.AmountLocked()
}

/*
token on road...
*/
func (c *Channel) Outstanding() *big.Int {
	return c.PartnerState.AmountLocked()
}

func (c *Channel) GetSettleExpiration(blocknumer int64) int64 {
	ClosedBlock := c.ExternState.ClosedBlock
	if ClosedBlock != 0 {
		return ClosedBlock + int64(c.SettleTimeout)
	} else {
		return blocknumer + int64(c.SettleTimeout)
	}
}

func (c *Channel) HandleClosed(blockNumber int64, closingAddress common.Address) {
	balanceProof := c.PartnerState.BalanceProofState
	//the channel was closed, update our half of the state if we need to
	if closingAddress != c.OurState.Address {
		c.ExternState.UpdateTransfer(balanceProof)
	}
	unlockProofs := c.PartnerState.GetKnownUnlocks()
	err := c.ExternState.WithDraw(unlockProofs)
	if err != nil {
		log.Error(fmt.Sprintf("withdraw on % failed, channel is gone, error:%s", c.MyAddress.String(), err))
	}
	c.IsCloseEventComplete = true
}

//there is nothing tod rightnow
func (c *Channel) HandleSettled(blockNumber int64) {

}

func (c *Channel) GetStateFor(nodeAddress common.Address) (*ChannelEndState, error) {
	if c.OurState.Address == nodeAddress {
		return c.OurState, nil
	}
	if c.PartnerState.Address == nodeAddress {
		return c.PartnerState, nil
	}
	return nil, fmt.Errorf("GetStateFor Unknown address %s", nodeAddress)
}

/*
Register a secret.

        This wont claim the lock (update the transferred_amount), it will only
        save the secret in case that a proof needs to be created. This method
        can be used for any of the ends of the channel.

        Note:
            When a secret is revealed a message could be in-transit containing
            the older lockroot, for this reason the recipient cannot update
            their locksroot at the moment a secret was revealed.

            The protocol is to register the secret so that it can compute a
            proof of balance, if necessary, forward the secret to the sender
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
		lock := c.OurState.GetLockByHashlock(hashlock)
		log.Debug(fmt.Sprintf("secret registered node=%s,from=%s,to=%s,token=%s,hashlock=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.OurState.Address[:]),
			utils.Pex(c.PartnerState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(hashlock[:]), lock.Amount))
		c.OurState.RegisterSecret(secret)
	}
	if partenerKnown {
		lock := c.PartnerState.GetLockByHashlock(hashlock)
		log.Debug(fmt.Sprintf("secret registered node=%s,from=%s,to=%s,token=%s,hashlock=%s,amount=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(c.PartnerState.Address[:]),
			utils.Pex(c.OurState.Address[:]), utils.APex(c.TokenAddress),
			utils.Pex(hashlock[:]), lock.Amount))
		c.PartnerState.RegisterSecret(secret)
	}
	return nil
}

//Register a signed transfer, updating the channel's state accordingly.
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
Validates and register a signed transfer, updating the channel's state accordingly.
Note:
            The transfer must be registered before it is sent, not on
            acknowledgement. That is necessary for two reasons:

            - Guarantee that the transfer is valid.
            - Avoid sending a new transaction without funds.

        Raises:
            InsufficientBalance: If the transfer is negative or above the distributable amount.
            InvalidLocksRoot: If locksroot check fails.
            InvalidNonce: If the expected nonce does not match.
            ValueError: If there is an address mismatch (token or node address).
*/
func (c *Channel) RegisterTransferFromTo(blockNumber int64, tr encoding.EnvelopMessager, fromState *ChannelEndState, toState *ChannelEndState) error {
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
	isInvalidNonce := (evMsg.Nonce < 1 || (fromState.Nonce() != 0 && evMsg.Nonce != fromState.Nonce()+1))
	//如果一个node数据损坏了,那么这个channel将不能工作了? 所以数据一定不能损坏
	if isInvalidNonce {
		//c may occur on normal operation
		log.Info(fmt.Sprintf("invalid nonce node=%s,from=%s,to=%s,expected nonce=%d,nonce=%d",
			utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]),
			utils.Pex(toState.Address[:]), fromState.Nonce(), evMsg.Nonce))
		return rerr.InvalidNonce(utils.StringInterface(tr, 3))
	}
	/*
				 if the locksroot is out-of-sync (because a transfer was created while
			    a Secret was in traffic) the balance _will_ be wrong, so first check
			    the locksroot and then the balance
		在创建完这个transfer和register这个transfer之间,收到了一个secret.
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
		_, expectedLocksroot := fromState.ComputeMerkleRootWith(mtr.GetLock())
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
		if isSender && expiresAfterSettle { //对方收到这个lock以后,可以到链上close或者updatetransfer,这样如果对方没有密码,自己仍然不能取到钱.
			log.Error(fmt.Sprintf("Lock expires after the settlement period. node=%s,from=%s,to=%s,lockexpiration=%d,currentblock=%d,end_settle_period=%d",
				utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
				mtr.Expiration, blockNumber, endSettlePeriod))
			return fmt.Errorf("Lock expires after the settlement period.")
		}
	}
	// only check the balance if the locksroot matched
	if evMsg.TransferAmount.Cmp(fromState.TransferAmount()) < 0 {
		log.Error(fmt.Sprintf("NEGATIVE TRANSFER node=%s,from=%s,to=%s,transfer=%s",
			utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
			utils.StringInterface(tr, 3))) //for nest struct
		return fmt.Errorf("Negative transfer")
	}
	amount := new(big.Int).Sub(evMsg.TransferAmount, fromState.TransferAmount())
	distributable := fromState.Distributable(toState)
	if tr.Cmd() == encoding.DIRECTTRANSFER_CMDID {
		if amount.Cmp(distributable) > 0 {
			return rerr.InsufficientBalance
		}
	} else if encoding.IsLockedTransfer(tr) {
		mtr := encoding.GetMtrFromLockedTransfer(tr)
		if new(big.Int).Add(amount, mtr.Amount).Cmp(distributable) > 0 {
			return rerr.InsufficientBalance
		}
	} else if tr.Cmd() == encoding.SECRET_CMDID {
		sec := tr.(*encoding.Secret)
		hashlock := utils.Sha3(sec.Secret[:])
		lock := fromState.GetLockByHashlock(hashlock)
		transferAmount := new(big.Int).Add(fromState.TransferAmount(), lock.Amount)
		/*
			 tr.transferred_amount could be larger than the previous
				             transferred_amount + lock.amount, that scenario is a bug of the
				             payer
		*/
		if sec.TransferAmount.Cmp(transferAmount) != 0 {
			return fmt.Errorf("invalid transferred_amount, expected: %s got: %s",
				transferAmount, sec.TransferAmount.Int64())
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
		err = fromState.RegisterLockedTransfer(tr)
		if err != nil {
			return err
		}
		/*
			 register c channel as waiting for the secret (the secret can
			be revealed through a message or a blockchain log)
		*/
		c.ExternState.funcRegisterChannelForHashlock(c, mtr.HashLock)
	}
	if tr.Cmd() == encoding.DIRECTTRANSFER_CMDID {
		err = fromState.RegisterDirectTransfer(tr.(*encoding.DirectTransfer))
		if err != nil {
			return err
		}
	}
	if tr.Cmd() == encoding.SECRET_CMDID {
		err = fromState.RegisterSecretMessage(tr.(*encoding.Secret))
		if err != nil {
			return err
		}
	}
	mroot := fromState.TreeState.Tree.MerkleRoot()
	log.Debug(fmt.Sprintf("'REGISTERED TRANSFER node=%s,from=%s,to=%s,transfer_amount=%s,nonce=%d,current_locksroot=%s,\ntransfer=%s",
		utils.Pex(c.OurState.Address[:]), utils.Pex(fromState.Address[:]), utils.Pex(toState.Address[:]),
		fromState.TransferAmount(), fromState.Nonce(), utils.Pex(mroot[:]), utils.StringInterface(tr, 3)))
	return nil
}

//change nonce  means banlance proof state changed
func (c *Channel) GetNextNonce() int64 {
	if c.OurState.Nonce() != 0 {
		return c.OurState.Nonce() + 1
	}
	// 0 must not be used since in the netting contract it represents null.
	return 1
}

/*
 Return a DirectTransfer message.

        This message needs to be signed and registered with the channel before
        sent.
*/
func (c *Channel) CreateDirectTransfer(amount *big.Int, identifier uint64) (tr *encoding.DirectTransfer, err error) {
	if !c.CanTransfer() {
		return nil, fmt.Errorf("Transfer not possible, no funding or channel closed.")
	}
	from := c.OurState
	to := c.PartnerState
	distributable := from.Distributable(to)
	if amount.Cmp(utils.BigInt0) <= 0 || amount.Cmp(distributable) > 0 {
		log.Debug(fmt.Sprintf("Insufficient funds : amount=%s, distributable=%s", amount, distributable))
		return nil, rerr.InsufficientFunds
	}
	tranferAmount := new(big.Int).Add(from.TransferAmount(), amount)
	currentLocksroot := to.TreeState.Tree.MerkleRoot()
	nonce := c.GetNextNonce()
	tr = encoding.NewDirectTransfer(identifier, nonce, c.TokenAddress, c.MyAddress, tranferAmount, to.Address, currentLocksroot)
	return
}

/*
Return a MediatedTransfer message.

        This message needs to be signed and registered with the channel before
        sent.

        Args:
            transfer_initiator (address): The node that requested the transfer.
            transfer_target (address): The final destination node of the transfer
            amount (float): How much of a token is being transferred.
            expiration (int): The maximum block number until the transfer
                message can be received.
*/
func (c *Channel) CreateMediatedTransfer(transfer_initiator, transfer_target common.Address, fee *big.Int, amount *big.Int, identifier uint64, expiration int64, hashlock common.Hash) (tr *encoding.MediatedTransfer, err error) {
	if !c.CanTransfer() {
		return nil, fmt.Errorf("Transfer not possible, no funding or channel closed.")
	}
	if amount.Cmp(utils.BigInt0) <= 0 || amount.Cmp(c.Distributable()) > 0 {
		log.Debug("Insufficient funds  amount=%s,distributable=%s", amount, c.Distributable())
		return nil, fmt.Errorf("Insufficient funds")
	}
	from := c.OurState
	to := c.PartnerState
	lock := &encoding.Lock{
		Amount:     amount,
		Expiration: expiration,
		HashLock:   hashlock,
	}
	_, updatedLocksroot := from.ComputeMerkleRootWith(lock)
	transferAmount := from.TransferAmount()
	nonce := c.GetNextNonce()
	tr = encoding.NewMediatedTransfer(identifier, nonce, c.TokenAddress, c.MyAddress,
		transferAmount, to.Address, updatedLocksroot, lock, transfer_target, transfer_initiator, fee)
	return
}

/*
similar as CreateMediatedTransfer
*/
func (c *Channel) CreateRefundTransfer(transfer_initiator, transfer_target common.Address, fee *big.Int, amount *big.Int, identifier uint64, expiration int64, hashlock common.Hash) (tr *encoding.RefundTransfer, err error) {
	mtr, err := c.CreateMediatedTransfer(transfer_initiator, transfer_target, fee, amount, identifier, expiration, hashlock)
	if err != nil {
		return
	}
	tr = encoding.NewRefundTransferFromMediatedTransfer(mtr)
	return
}

func (c *Channel) CreateSecret(identifer uint64, secret common.Hash) (tr *encoding.Secret, err error) {
	hashlock := utils.Sha3(secret[:])
	from := c.OurState
	lock := from.GetLockByHashlock(hashlock)
	if lock == nil {
		return nil, fmt.Errorf("no such lock for secret:%s", utils.HPex(secret))
	}
	_, locksrootWithPendingLockRemoved, err := from.ComputeMerkleRootWithout(lock)
	if err != nil {
		return
	}
	transferAmount := new(big.Int).Add(from.TransferAmount(), lock.Amount)
	nonce := c.GetNextNonce()
	tr = encoding.NewSecret(identifer, nonce, c.MyAddress, transferAmount, locksrootWithPendingLockRemoved, secret)
	return
}

func (c *Channel) StateTransition(st transfer.StateChange) (err error) {
	switch st2 := st.(type) {
	case *transfer.BlockStateChange:
		if c.State() == transfer.CHANNEL_STATE_CLOSED {
			settlementEnd := c.ExternState.ClosedBlock + int64(c.SettleTimeout)
			if st2.BlockNumber > settlementEnd {
				err = c.ExternState.Settle()
			}
		}
	case *mediated_transfer.ContractReceiveClosedStateChange:
		if st2.ChannelAddress == c.MyAddress {
			if !c.IsCloseEventComplete {
				c.ExternState.SetClosed(st2.ClosedBlock)
				c.HandleClosed(st2.ClosedBlock, st2.ClosingAddress)
			} else {
				log.Warn(fmt.Sprint("channel closed on a different block or close event happened twice channel=%s,closedblock=%s,thisblock=%sn",
					c.MyAddress, c.ExternState.ClosedBlock, st2.ClosedBlock))
			}
		}
	case *mediated_transfer.ContractReceiveSettledStateChange:
		//settled channel should be removed. todo bai fix it
		if st2.ChannelAddress == c.MyAddress {
			if c.ExternState.SetSettled(st2.SettledBlock) {
				c.HandleSettled(st2.SettledBlock)
			} else {
				log.Warn(fmt.Sprintf("channel is already settled on a different block channeladdress=%s,settleblock=%d,thisblock=%d",
					c.MyAddress.String(), c.ExternState.SettledBlock, st2.SettledBlock))
			}
		}
	case *mediated_transfer.ContractReceiveBalanceStateChange:
		participant := st2.ParticipantAddress
		balance := st2.Balance
		var channelState *ChannelEndState
		channelState, err = c.GetStateFor(participant)
		if err != nil {
			return
		}
		if channelState.ContractBalance.Cmp(balance) != 0 {
			err = channelState.UpdateContractBalance(balance)
		}
	}
	return
}
func (c *Channel) String() string {
	return fmt.Sprintf("{ContractBalance=%s,Balance=%s,Distributable=%s,locked=%s,transferAmount=%s}",
		c.ContractBalance(), c.Balance(), c.Distributable(), c.Locked(), c.TransferAmount())
}

type ChannelSerialization struct {
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
	SettleTimeout              int
}

func NewChannelSerialization(c *Channel) *ChannelSerialization {
	s := &ChannelSerialization{
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
	}
	return s
}
func init() {
	gob.Register(&ChannelSerialization{})
	//gob.Register(&Channel{})
	//gob.Register(&ChannelEndState{})
	//gob.Register(&ChannelExternalState{})
}
