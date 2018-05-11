package mediated_transfer

import (
	"encoding/gob"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
 Note: The init states must contain all the required data for trying doing
 useful work, ie. there must /not/ be an event for requesting new data.
*/
type ActionInitInitiatorStateChange struct {
	OurAddress      common.Address        //This node address.
	Tranfer         *LockedTransferState  //A state object containing the transfer details.
	Routes          *transfer.RoutesState //The current available routes.
	RandomGenerator utils.SecretGenerator //A generator for secrets.
	BlockNumber     int64                 //The current block number.
	Db              channel.Db            //get the latest channel state
}

// Initial state for a new mediator.
type ActionInitMediatorStateChange struct {
	OurAddress  common.Address             //This node address.
	FromTranfer *LockedTransferState       //The received MediatedTransfer.
	Routes      *transfer.RoutesState      //The current available routes.
	FromRoute   *transfer.RouteState       //The route from which the MediatedTransfer was received.
	BlockNumber int64                      //The current block number.
	Message     *encoding.MediatedTransfer //the message trigger this statechange
	Db          channel.Db                 //get the latest channel state
}

//Initial state for a new target.
type ActionInitTargetStateChange struct {
	OurAddress  common.Address       //This node address.
	FromTranfer *LockedTransferState //The received MediatedTransfer.
	FromRoute   *transfer.RouteState //The route from which the MediatedTransfer was received.
	BlockNumber int64
	Message     *encoding.MediatedTransfer //the message trigger this statechange
	Db          channel.Db                 //get the latest channel state
}

/*
Cancel the current route.
 Notes:
        Used to cancel a specific route but not the transfer, may be used for
        timeouts.
*/
type ActionCancelRouteStateChange struct {
	Identifier uint64
}

// A SecretRequest message received.
type ReceiveSecretRequestStateChange struct {
	Identifier uint64
	Amount     *big.Int
	Hashlock   common.Hash
	Sender     common.Address
	Message    *encoding.SecretRequest //the message trigger this statechange
}

//A SecretReveal message received
type ReceiveSecretRevealStateChange struct {
	Secret  common.Hash
	Sender  common.Address
	Message *encoding.RevealSecret //the message trigger this statechange
}

// A RefundTransfer message received.
type ReceiveTransferRefundStateChange struct {
	Sender   common.Address
	Transfer *LockedTransferState
	Message  *encoding.RefundTransfer //the message trigger this statechange
}

//A balance proof `identifier` was received.
type ReceiveBalanceProofStateChange struct {
	Identifier   uint64
	NodeAddress  common.Address
	BalanceProof *transfer.BalanceProofState
	Message      encoding.EnvelopMessager //the message trigger this statechange
}

/*
A lock was withdrawn via the blockchain.

    Used when a hash time lock was withdrawn and a log ChannelSecretRevealed is
    emited by the netting channel.

    Note:
        For this state change the contract caller is not important but only the
        receiving address. `receiver` is the address to which the lock's token
        was transferred, this may be either of the channel participants.

        If the channel was used for a mediated transfer that was refunded, this
        event must be used twice, once for each receiver.
*/
type ContractReceiveWithdrawStateChange struct {
	ChannelAddress common.Address
	Secret         common.Hash
	Receiver       common.Address //this address use secret to withdraw onchain
}

type ContractReceiveClosedStateChange struct {
	ChannelAddress common.Address
	ClosingAddress common.Address
	ClosedBlock    int64 //block number when close
}

type ContractReceiveSettledStateChange struct {
	ChannelAddress common.Address
	SettledBlock   int64
}

type ContractReceiveBalanceStateChange struct {
	ChannelAddress     common.Address
	TokenAddress       common.Address
	ParticipantAddress common.Address
	Balance            *big.Int //todo type error?
	BlockNumber        int64
}

type ContractReceiveNewChannelStateChange struct {
	ManagerAddress common.Address
	ChannelAddress common.Address
	Participant1   common.Address
	Participant2   common.Address
	SettleTimeout  int
}

type ContractReceiveTokenAddedStateChange struct {
	RegistryAddress common.Address
	TokenAddress    common.Address
	ManagerAddress  common.Address
}

func init() {
	gob.Register(&ActionInitInitiatorStateChange{})
	gob.Register(&ActionInitMediatorStateChange{})
	gob.Register(&ActionInitTargetStateChange{})
	gob.Register(&ActionCancelRouteStateChange{})
	gob.Register(&ReceiveSecretRequestStateChange{})
	gob.Register(&ReceiveSecretRevealStateChange{})
	gob.Register(&ReceiveTransferRefundStateChange{})
	gob.Register(&ReceiveBalanceProofStateChange{})
	gob.Register(&ContractReceiveWithdrawStateChange{})
	gob.Register(&ContractReceiveClosedStateChange{})
	gob.Register(&ContractReceiveSettledStateChange{})
	gob.Register(&ContractReceiveBalanceStateChange{})
	gob.Register(&ContractReceiveNewChannelStateChange{})
	gob.Register(&ContractReceiveTokenAddedStateChange{})
}
