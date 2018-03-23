package mediated_transfer

import (
	"encoding/gob"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
State of a transfer that is time hash locked.

    Args:
        identifier (int): A unique identifer for the transfer.
        amount (int): Amount of `token` being transferred.
        token (address): Token being transferred.
        target (address): Transfer target address.
        expiration (int): The absolute block number that the lock expires.
        hashlock (bin): The hashlock.
        secret (bin): The secret that unlocks the lock, may be None.
*/
type LockedTransferState struct {
	Identifier uint64         //A unique identifer for the transfer.
	Amount     *big.Int       // Amount of `token` being transferred.
	Token      common.Address //Token being transferred.
	Initiator  common.Address //Transfer initiator
	Target     common.Address //Transfer target address.
	Expiration int64          //The absolute block number that the lock expires.
	Hashlock   common.Hash    // The hashlock.
	Secret     common.Hash    //The secret that unlocks the lock, may be None.
}

func (self *LockedTransferState) AlmostEqual(other *LockedTransferState) bool {
	//expiration maybe different
	return self.Identifier == other.Identifier &&
		self.Amount.Cmp(other.Amount) == 0 &&
		self.Token == other.Token &&
		self.Target == other.Target &&
		self.Hashlock == other.Hashlock &&
		self.Secret == other.Secret
}

//Create LockedTransferState from a MediatedTransfer message.
func LockedTransferFromMessage(msg *encoding.MediatedTransfer) *LockedTransferState {
	return &LockedTransferState{
		Identifier: msg.Identifier,
		Amount:     new(big.Int).Set(msg.Amount),
		Token:      msg.Token,
		Initiator:  msg.Initiator,
		Target:     msg.Target,
		Expiration: msg.Expiration,
		Hashlock:   msg.HashLock,
	}
}

/*
State of a node initiating a mediated transfer.

    Args:
        our_address (address): This node address.
        transfer (LockedTransferState): The description of the mediated transfer.
        routes (RoutesState): Routes available for this transfer.
        block_number (int): Latest known block number.
        random_generator (generator): A generator that yields valid secrets.
*/

type InitiatorState struct {
	OurAddress        common.Address        //This node address.
	Transfer          *LockedTransferState  // The description of the mediated transfer.
	Routes            *transfer.RoutesState //Routes available for this transfer.
	BlockNumber       int64                 //Latest known block number.
	RandomGenerator   utils.SecretGenerator //A generator that yields valid secrets.
	Message           transfer.Event        // current message in-transit todo this type?
	Route             *transfer.RouteState  //current route being used
	SecretRequest     *encoding.SecretRequest
	RevealSecret      *EventSendRevealSecret
	CanceledTransfers []transfer.Event
	Db                channel.ChannelDb
}

/*
State of a node mediating a transfer.

    Args:
        our_address (address): This node address.
        routes (RoutesState): Routes available for this transfer.
        block_number (int): Latest known block number.
        hashlock (bin): The hashlock used for this transfer.
*/
type MediatorState struct {
	OurAddress  common.Address        //This node address.
	Routes      *transfer.RoutesState //Routes available for this transfer.
	BlockNumber int64                 //Latest known block number.
	Hashlock    common.Hash           //  The hashlock used for this transfer.
	Secret      common.Hash
	/*
			keeping all transfers in a single list byzantine behavior for secret
		        reveal and simplifies secret setting
	*/
	TransfersPair []*MediationPairState
	HasRefunded   bool //此节点已经发生了refund，肯定不能再用了。
	Db            channel.ChannelDb
}

/*
 Set the secret to all mediated transfers.

    It doesn't matter if the secret was learned through the blockchain or a
    secret reveal message.
*/
func (this *MediatorState) SetSecret(secret common.Hash) {
	this.Secret = secret
	for _, p := range this.TransfersPair {
		p.PayerTransfer.Secret = secret
		p.PayeeTransfer.Secret = secret
	}
}

const STATE_SECRET_REQUEST = "secret_request"
const STATE_REVEAL_SECRET = "reveal_secret"
const STATE_BALANCE_PROOF = "balance_proof"
const STATE_WAITING_CLOSE = "waiting_close"

//State of mediated transfer target.
type TargetState struct {
	OurAddress   common.Address
	FromRoute    *transfer.RouteState
	FromTransfer *LockedTransferState
	BlockNumber  int64
	Secret       common.Hash
	State        string // default secret_request
	Db           channel.ChannelDb
}

/*
State for a mediated transfer.

    A mediator will pay payee node knowing that there is a payer node to cover
    the token expenses. This state keeps track of the routes and transfer for
    the payer and payee, and the current state of the payment.
*/
type MediationPairState struct {
	PayeeRoute    *transfer.RouteState
	PayeeTransfer *LockedTransferState
	PayeeState    string
	PayerRoute    *transfer.RouteState
	PayerTransfer *LockedTransferState
	PayerState    string
}

//all payee's valid state
// Initial state.
const STATE_PAYEE_PENDING = "payee_pending"

//The payee is following the raiden protocol and has sent a SecretReveal.
const STATE_PAYEE_SECRET_REVEALED = "payee_secret_revealed"

/*
   The corresponding refund transfer was withdrawn on-chain, the payee has
       /not/ withdrawn the lock yet, it only learned the secret through the
       blockchain.
       Note: This state is reachable only if there is a refund transfer, that
       is represented by a different MediationPairState, and the refund
       transfer is at 'payer_contract_withdraw'.
*/
const STATE_PAYEE_REFUND_WITHDRAW = "payee_refund_withdraw"

/*
 The payee received the token on-chain. A transition to this state is
valid from all but the `payee_expired` state.
*/
const STATE_PAYEE_CONTRACT_WITHDRAW = "payee_contract_withdraw"

/*
   This node has sent a SendBalanceProof to the payee with the balance
    updated.
*/
const STATE_PAYEE_BALANCE_PROOF = "payee_balance_proof"

//The lock has expired.
const STATE_PAYEE_EXPIRED = "payee_expired"

var ValidPayeeStateMap = map[string]bool{
	STATE_PAYEE_PENDING:           true,
	STATE_PAYEE_SECRET_REVEALED:   true,
	STATE_PAYEE_REFUND_WITHDRAW:   true,
	STATE_PAYEE_CONTRACT_WITHDRAW: true,
	STATE_PAYEE_BALANCE_PROOF:     true,
	STATE_PAYEE_EXPIRED:           true,
}

var ValidPayerStateMap = map[string]bool{
	STATE_PAYER_PENDING:           true,
	STATE_PAYER_SECRET_REVEALED:   true,
	STATE_PAYER_WAITING_CLOSE:     true,
	STATE_PAYER_WAITING_WITHDRAW:  true,
	STATE_PAYER_CONTRACT_WITHDRAW: true,
	STATE_PAYER_BALANCE_PROOF:     true,
	STATE_PAYER_EXPIRED:           true,
}

//payer's state
const STATE_PAYER_PENDING = "payer_pending"

//SendRevealSecret was sent
const STATE_PAYER_SECRET_REVEALED = "payer_secret_revealed"

//ContractSendChannelClose was sent
const STATE_PAYER_WAITING_CLOSE = "payer_waiting_close"

//ContractSendWithdraw was sent
const STATE_PAYER_WAITING_WITHDRAW = "payer_waiting_withdraw"

// ContractReceiveWithdraw for the above send received
const STATE_PAYER_CONTRACT_WITHDRAW = "payer_contract_withdraw"

// ReceiveBalanceProof was received
const STATE_PAYER_BALANCE_PROOF = "payer_balance_proof"

//None of the above happened and the lock expired
const STATE_PAYER_EXPIRED = "payer_expired"

/*
  Args:
           payer_route (RouteState): The details of the route with the payer.
           payer_transfer (LockedTransferState): The transfer this node
               *received* that will cover the expenses.

           payee_route (RouteState): The details of the route with the payee.
           payee_transfer (LockedTransferState): The transfer this node *sent*
               that will be withdrawn by the payee.
*/
func NewMediationPairState(payerRoute, payeeRoute *transfer.RouteState, payerTransfer, payeeTransfer *LockedTransferState) *MediationPairState {
	return &MediationPairState{
		PayerRoute:    payerRoute,
		PayerTransfer: payerTransfer,
		PayerState:    STATE_PAYER_PENDING,
		PayeeRoute:    payeeRoute,
		PayeeTransfer: payeeTransfer,
		PayeeState:    STATE_PAYEE_PENDING,
	}
}

func init() {
	gob.Register(&LockedTransferState{})
	gob.Register(&InitiatorState{})
	gob.Register(&MediatorState{})
	gob.Register(&TargetState{})
	gob.Register(&MediationPairState{})
}
