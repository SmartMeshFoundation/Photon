package mediated_transfer

import (
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
)

// A mediated transfer that must be sent to `node_address`.
type EventSendMediatedTransfer struct {
	Identifier uint64
	Token      common.Address
	Amount     int64
	HashLock   common.Hash
	Initiator  common.Address
	Target     common.Address
	Expiration int64
	Receiver   common.Address
}

func NewEventSendMediatedTransfer(transfer *LockedTransferState, receiver common.Address) *EventSendMediatedTransfer {
	return &EventSendMediatedTransfer{
		Identifier: transfer.Identifier,
		Token:      transfer.Token,
		Amount:     transfer.Amount,
		HashLock:   transfer.Hashlock,
		Initiator:  transfer.Initiator,
		Target:     transfer.Target,
		Expiration: transfer.Expiration,
		Receiver:   receiver,
	}
}

/*
Sends a RevealSecret to another node.

    This event is used once the secret is known locally and an action must be
    performed on the receiver:

        - For receivers in the payee role, it informs the node that the lock has
            been released and the token can be withdrawn, either on-chain or
            off-chain.
        - For receivers in the payer role, it tells the payer that the payee
            knows the secret and wants to withdraw the lock off-chain, so the payer
            may unlock the lock and send an up-to-date balance proof to the payee,
            avoiding on-chain payments which would require the channel to be
            closed.

    For any mediated transfer:

        - The initiator will only perform the payer role.
        - The target will only perform the payee role.
        - The mediators will have `n` channels at the payee role and `n` at the
          payer role, where `n` is equal to `1 + number_of_refunds`.

    Note:
        The payee must only update its local balance once the payer sends an
        up-to-date balance-proof message. This is a requirement for keeping the
        nodes synchronized. The reveal secret message flows from the receiver
        to the sender, so when the secret is learned it is not yet time to
        update the balance.
*/
type EventSendRevealSecret struct {
	Identifier uint64
	Secret     common.Hash
	Token      common.Address
	Receiver   common.Address
	Sender     common.Address
}

/*
 Event to send a balance-proof to the counter-party, used after a lock
    is unlocked locally allowing the counter-party to withdraw.

    Used by payers: The initiator and mediator nodes.

    Note:
        This event has a dual role, it serves as a synchronization and as
        balance-proof for the netting channel smart contract.

        Nodes need to keep the last known merkle root synchronized. This is
        required by the receiving end of a transfer in order to properly
        validate. The rule is "only the party that owns the current payment
        channel may change it" (remember that a netting channel is composed of
        two uni-directional channels), as a consequence the merkle root is only
        updated by the receiver once a balance proof message is received.
*/
type EventSendBalanceProof struct {
	Identifier     uint64
	ChannelAddress common.Address
	Token          common.Address
	Receiver       common.Address
	Secret         common.Hash //Secret is not required for the balance proof to dispatch the message
}

/*
Event used by a target node to request the secret from the initiator
    (`receiver`).
*/
type EventSendSecretRequest struct {
	Identifer uint64
	Amount    int64
	Hashlock  common.Hash
	Receiver  common.Address
}

/*
Event used to cleanly backtrack the current node in the route.

    This message will pay back the same amount of token from the receiver to
    the sender, allowing the sender to try a different route without the risk
    of losing token.
*/
type EventSendRefundTransfer struct {
	Identifier uint64
	Token      common.Address
	Amount     int64
	HashLock   common.Hash
	Initiator  common.Address
	Target     common.Address
	Expiration int64
	Receiver   common.Address
}

/*
Event emitted to close the netting channel.

    This event is used when a node needs to prepare the channel to withdraw
    on-chain.
*/
type EventContractSendChannelClose struct {
	ChannelAddress common.Address
	Token          common.Address
}

//Event emitted when the lock must be withdrawn on-chain.
type EventContractSendWithdraw struct {
	Transfer       *LockedTransferState
	ChannelAddress common.Address
}

//Event emitted when a lock unlock succeded ,emit this event after receive a revealsecret message
type EventUnlockSuccess struct {
	Identifier uint64
	Hashlock   common.Hash
}

//Event emitted when a lock unlock failed.
type EventUnlockFailed struct {
	Identifier uint64
	Hashlock   common.Hash
	Reason     string
}

//Event emitted when a lock withdraw succeded.
type EventWithdrawSuccess struct {
	Identifier uint64
	Hashlock   common.Hash
}

//Event emitted when a lock withdraw failed.
type EventWithdrawFailed struct {
	Identifier uint64
	Hashlock   common.Hash
	Reason     string
}

func init() {
	gob.Register(&EventSendMediatedTransfer{})
	gob.Register(&EventSendRevealSecret{})
	gob.Register(&EventSendBalanceProof{})
	gob.Register(&EventSendSecretRequest{})
	gob.Register(&EventSendRefundTransfer{})
	gob.Register(&EventContractSendChannelClose{})
	gob.Register(&EventContractSendWithdraw{})
	gob.Register(&EventUnlockSuccess{})
	gob.Register(&EventUnlockFailed{})
	gob.Register(&EventWithdrawSuccess{})
	gob.Register(&EventWithdrawFailed{})
}
