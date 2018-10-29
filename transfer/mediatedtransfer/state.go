package mediatedtransfer

import (
	"encoding/gob"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/transfer/route"
	"github.com/ethereum/go-ethereum/common"
)

/*
LockedTransferState is State of a transfer that is time hash locked.
*/
type LockedTransferState struct {
	TargetAmount   *big.Int       //amount target should recevied
	Amount         *big.Int       // Amount of `token` being transferred.
	Token          common.Address //Token being transferred.
	Initiator      common.Address //Transfer initiator
	Target         common.Address //Transfer target address.
	Expiration     int64          //The absolute block number that the lock expires.
	LockSecretHash common.Hash    // The hashlock.
	Secret         common.Hash    //The secret that unlocks the lock, may be None.
	Fee            *big.Int       // how much fee left for other hop node.
}

//AlmostEqual if two state equals?
func (l *LockedTransferState) AlmostEqual(other *LockedTransferState) bool {
	//expiration maybe different
	return l.TargetAmount.Cmp(other.TargetAmount) == 0 &&
		l.Token == other.Token &&
		l.Target == other.Target &&
		l.LockSecretHash == other.LockSecretHash &&
		l.Secret == other.Secret
}

//LockedTransferFromMessage Create LockedTransferState from a MediatedTransfer message.
func LockedTransferFromMessage(msg *encoding.MediatedTransfer, tokenAddress common.Address) *LockedTransferState {
	return &LockedTransferState{
		TargetAmount:   new(big.Int).Sub(msg.PaymentAmount, msg.Fee),
		Amount:         new(big.Int).Set(msg.PaymentAmount),
		Initiator:      msg.Initiator,
		Target:         msg.Target,
		Expiration:     msg.Expiration,
		LockSecretHash: msg.LockSecretHash,
		Fee:            msg.Fee,
		Token:          tokenAddress,
	}
}

/*
LockAndChannel the lock and associated channel
*/
type LockAndChannel struct {
	Lock    *mtree.Lock
	Channel *channel.Channel
}

//CrashState is state of a node after restarting, but found some transfer is not finished.
type CrashState struct {
	OurAddress             common.Address
	LockSecretHash         common.Hash
	Token                  common.Address //which token
	SentLocks              []*LockAndChannel
	ReceivedLocks          []*LockAndChannel
	ProcessedSentLocks     []*LockAndChannel
	ProcessedReceivedLocks []*LockAndChannel
}

/*
InitiatorState is State of a node initiating a mediated transfer.
*/
type InitiatorState struct {
	OurAddress        common.Address       //This node address.
	Transfer          *LockedTransferState // The description of the mediated transfer.
	Routes            *route.RoutesState   //Routes available for this transfer.
	BlockNumber       int64                //Latest known block number.
	LockSecretHash    common.Hash
	Secret            common.Hash
	Message           *EventSendMediatedTransfer // current message in-transit
	Route             *route.State               //current route being used
	SecretRequest     *encoding.SecretRequest
	RevealSecret      *EventSendRevealSecret
	CanceledTransfers []*EventSendMediatedTransfer
	Db                channeltype.Db
}

/*
MediatorState is State of a node mediating a transfer.
*/
type MediatorState struct {
	OurAddress  common.Address     //This node address.
	Routes      *route.RoutesState //Routes available for this transfer.
	BlockNumber int64              //Latest known block number.
	Hashlock    common.Hash        //  The hashlock used for this transfer.
	Secret      common.Hash
	/*
			keeping all transfers in a single list byzantine behavior for secret
		        reveal and simplifies secret setting
	*/
	TransfersPair  []*MediationPairState
	LockSecretHash common.Hash
	Token          common.Address
	Db             channeltype.Db
}

/*
SetSecret Set the secret to all mediated transfers.

    It doesn't matter if the secret was learned through the blockchain or a
    secret reveal message.
*/
func (m *MediatorState) SetSecret(secret common.Hash) {
	m.Secret = secret
	for _, p := range m.TransfersPair {
		p.PayerTransfer.Secret = secret
		p.PayeeTransfer.Secret = secret
	}
}

//StateSecretRequest receive  secret request
const StateSecretRequest = "secret_request"

//StateRevealSecret receive reveal secret
const StateRevealSecret = "reveal_secret"

//StateBalanceProof receive balance proof
const StateBalanceProof = "balance_proof"

//StateWaitingRegisterSecret wait register secret on chain
const StateWaitingRegisterSecret = "waiting_register_secret"

/*
StateSecretRegistered 密码已经在链上披露了
整个交易的所有参与方都可以认为这笔交易从彻底完成了,
无论是发起方还是中间节点以及接收方
发起方收到SecretRegistered, 应该立即发送 unlock 消息
中间节点也是一样,收到 secretRegistered 以后,必须立即给下家发送 unlock 消息,无论有没有收到上家的 unlock 消息
*/
// StateSecretRegistered : channel state meaning that secret has been registered on-chain.
// Then all channel participant assume that this transfer completes, no matter transfer initiator, mediator, or recipient.
// And they should immediately send unlock to their partner once they have received SecretRegistered.
const StateSecretRegistered = "secret_registered"

//TargetState State of mediated transfer target.
type TargetState struct {
	OurAddress   common.Address
	FromRoute    *route.State
	FromTransfer *LockedTransferState
	BlockNumber  int64
	Secret       common.Hash
	State        string // default secret_request
	Db           channeltype.Db
}

/*
MediationPairState State for a mediated transfer.

    A mediator will pay payee node knowing that there is a payer node to cover
    the token expenses. This state keeps track of the routes and transfer for
    the payer and payee, and the current state of the payment.
*/
type MediationPairState struct {
	PayeeRoute    *route.State
	PayeeTransfer *LockedTransferState
	PayeeState    string
	PayerRoute    *route.State
	PayerTransfer *LockedTransferState
	PayerState    string
}

//all payee's valid state

//StatePayeePending Initial state.
const StatePayeePending = "payee_pending"

//StatePayeeSecretRevealed The payee is following the photon protocol and has sent a SecretReveal.
//#nosec
const StatePayeeSecretRevealed = "payee_secret_revealed"

/*
StatePayeeBalanceProof   This node has sent a SendBalanceProof to the payee with the balance
    updated.
*/
const StatePayeeBalanceProof = "payee_balance_proof"

//StatePayeeExpired The lock has expired.
const StatePayeeExpired = "payee_expired"

//ValidPayeeStateMap payee's valid state
var ValidPayeeStateMap = map[string]bool{
	StatePayeePending:        true,
	StatePayeeSecretRevealed: true,
	StatePayeeBalanceProof:   true,
	StatePayeeExpired:        true,
}

//ValidPayerStateMap payer's valid state
var ValidPayerStateMap = map[string]bool{
	StatePayerPending:               true,
	StatePayerSecretRevealed:        true,
	StatePayerWaitingRegisterSecret: true,
	StatePayerBalanceProof:          true,
	StatePayerExpired:               true,
}

//payer's state

//StatePayerPending payer pending
const StatePayerPending = "payer_pending"

//StatePayerSecretRevealed RevealSecret was sent
const StatePayerSecretRevealed = "payer_secret_revealed"

//StatePayerWaitingRegisterSecret register secret on chain
const StatePayerWaitingRegisterSecret = "payer_waiting_register_secret"

//StatePayerBalanceProof ReceiveBalanceProof was received
const StatePayerBalanceProof = "payer_balance_proof"

//StatePayerExpired None of the above happened and the lock expired
const StatePayerExpired = "payer_expired"

/*
NewMediationPairState create mediated state
           payerRoute : The details of the route with the payer.
           payerTransfer  : The transfer this node
               *received* that will cover the expenses.
           payeeRoute  : The details of the route with the payee.
           payeeTransfer  : The transfer this node *sent*
               that will be withdrawn by the payee.
*/
func NewMediationPairState(payerRoute, payeeRoute *route.State, payerTransfer, payeeTransfer *LockedTransferState) *MediationPairState {
	return &MediationPairState{
		PayerRoute:    payerRoute,
		PayerTransfer: payerTransfer,
		PayerState:    StatePayerPending,
		PayeeRoute:    payeeRoute,
		PayeeTransfer: payeeTransfer,
		PayeeState:    StatePayeePending,
	}
}

func init() {
	gob.Register(&LockedTransferState{})
	gob.Register(&InitiatorState{})
	gob.Register(&MediatorState{})
	gob.Register(&TargetState{})
	gob.Register(&MediationPairState{})
}
