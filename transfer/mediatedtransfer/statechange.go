package mediatedtransfer

import (
	"encoding/gob"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
ActionInitInitiatorStateChange start a mediated transfer
 Note: The init states must contain all the required data for trying doing
 useful work, ie. there must /not/ be an event for requesting new data.
*/
type ActionInitInitiatorStateChange struct {
	OurAddress      common.Address        //This node address.
	Tranfer         *LockedTransferState  //A state object containing the transfer details.
	Routes          *transfer.RoutesState //The current available routes.
	RandomGenerator utils.SecretGenerator //A generator for secrets.
	BlockNumber     int64                 //The current block number.
	Db              channeltype.Db        //get the latest channel state
}

//ActionInitMediatorStateChange  Initial state for a new mediator.
type ActionInitMediatorStateChange struct {
	OurAddress  common.Address             //This node address.
	FromTranfer *LockedTransferState       //The received MediatedTransfer.
	Routes      *transfer.RoutesState      //The current available routes.
	FromRoute   *transfer.RouteState       //The route from which the MediatedTransfer was received.
	BlockNumber int64                      //The current block number.
	Message     *encoding.MediatedTransfer //the message trigger this statechange
	Db          channeltype.Db             //get the latest channel state
}

//ActionInitTargetStateChange Initial state for a new target.
type ActionInitTargetStateChange struct {
	OurAddress  common.Address       //This node address.
	FromTranfer *LockedTransferState //The received MediatedTransfer.
	FromRoute   *transfer.RouteState //The route from which the MediatedTransfer was received.
	BlockNumber int64
	Message     *encoding.MediatedTransfer //the message trigger this statechange
	Db          channeltype.Db             //get the latest channel state
}

/*
ActionCancelRouteStateChange Cancel the current route.
 Notes:
        Used to cancel a specific route but not the transfer, may be used for
        timeouts.
*/
type ActionCancelRouteStateChange struct {
	Identifier uint64
}

//ReceiveSecretRequestStateChange A SecretRequest message received.
type ReceiveSecretRequestStateChange struct {
	Identifier uint64
	Amount     *big.Int
	Hashlock   common.Hash
	Sender     common.Address
	Message    *encoding.SecretRequest //the message trigger this statechange
}

//ReceiveSecretRevealStateChange A SecretReveal message received
type ReceiveSecretRevealStateChange struct {
	Secret  common.Hash
	Sender  common.Address
	Message *encoding.RevealSecret //the message trigger this statechange
}

//ReceiveTransferRefundStateChange A AnnounceDisposed message received.
type ReceiveTransferRefundStateChange struct {
	Sender   common.Address
	Transfer *LockedTransferState
	Message  *encoding.AnnounceDisposed //the message trigger this statechange
}

//ReceiveBalanceProofStateChange A balance proof `identifier` was received.
type ReceiveBalanceProofStateChange struct {
	Identifier   uint64
	NodeAddress  common.Address
	BalanceProof *transfer.BalanceProofState
	Message      encoding.EnvelopMessager //the message trigger this statechange
}

/*
ContractSecretRevealStateChange A lock was withdrawn via the blockchain.

    Used when a hash time lock was withdrawn and a log ChannelSecretRevealed is
    emited by the netting channel.

    Note:
        For this state change the contract caller is not important but only the
        receiving address. `receiver` is the address to which the lock's token
        was transferred, this may be either of the channel participants.

        If the channel was used for a mediated transfer that was refunded, this
        event must be used twice, once for each receiver.
*/
type ContractSecretRevealStateChange struct {
	Secret common.Hash
}

type ContractReceiveChannelWithdrawStateChange struct {
	ChannelAddress *contracts.ChannelUniqueID
	//剩余的 balance 有意义?目前提供的 Event 并不知道 Participant1是谁,所以没啥用.
	Participant1        common.Address
	Participant1Balance *big.Int
	Participant2        common.Address
	Participant2Balance *big.Int
}

//ContractReceiveClosedStateChange a channel was closed
type ContractReceiveClosedStateChange struct {
	ChannelAddress    *contracts.ChannelUniqueID
	ClosingAddress    common.Address
	ClosedBlock       int64 //block number when close
	LocksRoot         common.Hash
	TransferredAmount *big.Int
}

//ContractReceiveSettledStateChange a channel was settled
type ContractReceiveSettledStateChange struct {
	ChannelAddress *contracts.ChannelUniqueID
	SettledBlock   int64
}

//ContractReceiveCooperativeSettledStateChange a channel was cooperatively settled
type ContractReceiveCooperativeSettledStateChange struct {
	ChannelAddress *contracts.ChannelUniqueID
	SettledBlock   int64
}

//ContractReceiveBalanceStateChange new deposit on channel
type ContractReceiveBalanceStateChange struct {
	ChannelAddress     *contracts.ChannelUniqueID
	ParticipantAddress common.Address
	Balance            *big.Int
	BlockNumber        int64
}

//ContractReceiveNewChannelStateChange new channel created on block chain
type ContractReceiveNewChannelStateChange struct {
	ChannelAddress *contracts.ChannelUniqueID
	Participant1   common.Address
	Participant2   common.Address
	SettleTimeout  int
}

//ContractReceiveTokenAddedStateChange a new token registered
type ContractReceiveTokenAddedStateChange struct {
	RegistryAddress     common.Address
	TokenAddress        common.Address
	TokenNetworkAddress common.Address
}

//ContractBalanceProofUpdatedStateChange contrct TransferUpdated event
type ContractBalanceProofUpdatedStateChange struct {
	ChannelAddress    *contracts.ChannelUniqueID
	Participant       common.Address
	LocksRoot         common.Hash
	TransferredAmount *big.Int
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
	gob.Register(&ContractSecretRevealStateChange{})
	gob.Register(&ContractReceiveClosedStateChange{})
	gob.Register(&ContractReceiveSettledStateChange{})
	gob.Register(&ContractReceiveBalanceStateChange{})
	gob.Register(&ContractReceiveNewChannelStateChange{})
	gob.Register(&ContractReceiveTokenAddedStateChange{})
	gob.Register(&ContractBalanceProofUpdatedStateChange{})
}
