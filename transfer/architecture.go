package transfer

import (
	"encoding/gob"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/ethereum/go-ethereum/common"
)

/*
 Quick overview
 --------------

 Goals:
 - Reliable failure recovery.

 Approach:
 - Use a write-ahead-log for state changes. Under a node restart the
 latest state snapshot can be recovered and the pending state changes
 reaplied.

 Requirements:
 - The function call `state_transition(curr_state, state_change)` must be
 deterministic, the recovery depends on the re-execution of the state changes
 from the WAL and must produce the same result.
 - StateChange must be idenpotent because the partner node might be recovering
 from a failure and a Event might be produced more than once.

 Requirements that are enforced:
 - A state_transition function must not produce a result that must be further
 processed, i.e. the state change must be self contained and the result state
 tree must be serializable to produce a snapshot. To enforce this inputs and
 outputs are separated under different class hierarquies (StateChange and Event).
*/
/*
 """ An isolated state, modified by StateChange messages.

    Notes:
    - Don't duplicate the same state data in two different States, instead use
    identifiers.
    - State objects may be nested.
    - State classes don't have logic by design.
    - Each iteration must operate on fresh copy of the state, treating the old
          objects as immutable.
    - This class is used as a marker for states.
    """
*/
type State interface{}

/*
Events produced by the execution of a state change.

    Nomenclature convention:
    - 'Send' prefix for protocol messages.
    - 'ContractSend' prefix for smart contract function calls.
    - 'Event' for node events.

    Notes:
    - This class is used as a marker for events.
    - These objects don't have logic by design.
    - Separate events are preferred because there is a decoupling of what the
      upper layer will use the events for.
*/
type Event interface{}

/*
Declare the transition to be applied in a state object.

    StateChanges are incoming events that change this node state (eg. a
    blockchain event, a new packet, an error). It is not used for the node to
    communicate with the outer world.

    Nomenclature convention:
    - 'Receive' prefix for protocol messages.
    - 'ContractReceive' prefix for smart contract logs.
    - 'Action' prefix for other interactions.

    Notes:
    - These objects don't have logic by design.
    - This class is used as a marker for state changes.
*/
type StateChange interface{}
type TransitionResult struct {
	NewState State
	Events   []Event
}

//def state_transition(state, state_change):
/*
 The mutable storage for the application state, this storage can do
    state transitions by applying the StateChanges to the current State.
*/
type FuncStateTransition func(state State, stateChange StateChange) *TransitionResult
type StateManager struct {
	ID                  int64 `storm:"id,increment"`
	FuncStateTransition FuncStateTransition
	CurrentState        State
	Name                string
	LastActive          time.Time //the latest message sent time
	ManagerState        string    `storm:"index"` //state for initiator and target ,distingush operation from crash
	Identifier          uint64    //transfer identifier
	TokenAddress        common.Address
	ChannelAddress      common.Address //channel address from initiator A-B-C channel A-B
	ChannelAddressTo    common.Address //mediated transfer will send to. A-B-C channel B-C
	ChannelAddresRefund common.Address //node received a refund transfer, should save and forget.

	LastReceivedMessage    interface{}       //message received status, except reveal secret,may be init statechange
	LastSendMessage        encoding.Messager //sending message.
	IsBalanceProofSent     bool              //mediatedtransfer must both true for finish
	IsBalanceProofReceived bool              //mediatedtransfer must both true for finish
}

const StateManager_State_Init = "ManagerInit"
const StateManager_ReceivedMessage = "ManagerReceivedOneMessage"
const StateManager_SendMessage = "ManagerSendMessage" //may sending several message, for example reveal secret
const StateManager_ReceivedMessageProcessComplete = "ManagerReceivedMessageComplete"
const StateManager_SendMessageSuccesss = "ManagerSendMessageSuccess"
const StateManager_TransferComplete = "ManagerTransferComplete"

type MessageTag struct {
	stateManager           *StateManager //message related statemanager, this field should not save to database because of cycle reference
	ReceiveProcessComplete bool          //Whether the receipt of the message has been processed,
	SendingMessageComplete bool          //Whether the message sent has received ACK
	IsASendingMessage      bool          //this message is on sending or receiveing?
	MessageId              string        //messageId for ping message
	EchoHash               common.Hash
	Receiver               common.Address
}

func (this *MessageTag) GetStateManager() *StateManager {
	return this.stateManager
}
func (this *MessageTag) SetStateManager(stateManager *StateManager) {
	this.stateManager = stateManager
}
func NewStateManager(stateTransition FuncStateTransition, currentState State, name string, identifier uint64, tokenAddress common.Address) *StateManager {
	return &StateManager{
		FuncStateTransition: stateTransition,
		CurrentState:        currentState,
		Name:                name,
		ManagerState:        StateManager_State_Init,
		LastActive:          time.Now(),
		Identifier:          identifier,
		TokenAddress:        tokenAddress,
	}
}

/*
Apply the `state_change` in the current machine and return the
        resulting events.

        Args:
            state_change (StateChange): An object representation of a state
            change.

        Return:
            Event: A list of events produced by the state transition, it's
            the upper layer's responsibility to decided how to handle these
            events.
*/
func (this *StateManager) Dispatch(stateChange StateChange) (events []Event) {

	/*
			     the state objects must be treated as immutable, so make a copy of the
		         current state and pass the copy to the state machine to be modified.
		        next_state = deepcopy(self.current_state)
			todo why clone?
	*/
	transitionResult := this.FuncStateTransition(this.CurrentState, stateChange)
	this.CurrentState, events = transitionResult.NewState, transitionResult.Events
	return
}

func init() {
	gob.Register(&StateManager{})
	gob.Register(&TransitionResult{})
	gob.Register(&MessageTag{})
}
