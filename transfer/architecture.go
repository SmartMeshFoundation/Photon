package transfer

/*
# Quick overview
# --------------
#
# Goals:
# - Reliable failure recovery.
#
# Approach:
# - Use a write-ahead-log for state changes. Under a node restart the
# latest state snapshot can be recovered and the pending state changes
# reaplied.
#
# Requirements:
# - The function call `state_transition(curr_state, state_change)` must be
# deterministic, the recovery depends on the re-execution of the state changes
# from the WAL and must produce the same result.
# - StateChange must be idenpotent because the partner node might be recovering
# from a failure and a Event might be produced more than once.
#
# Requirements that are enforced:
# - A state_transition function must not produce a result that must be further
# processed, i.e. the state change must be self contained and the result state
# tree must be serializable to produce a snapshot. To enforce this inputs and
# outputs are separated under different class hierarquies (StateChange and Event).
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
	FuncStateTransition FuncStateTransition
	CurrentState        State
}

//todo StateManager 看起来完全没用
func NewStateManager(stateTransition FuncStateTransition, currentState State) *StateManager {
	return &StateManager{
		FuncStateTransition: stateTransition, CurrentState: currentState,
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
			    # the state objects must be treated as immutable, so make a copy of the
		        # current state and pass the copy to the state machine to be modified.
		        next_state = deepcopy(self.current_state)
			todo why clone?
	*/
	transitionResult := this.FuncStateTransition(this.CurrentState, stateChange)
	this.CurrentState, events = transitionResult.NewState, transitionResult.Events
	return
}

/*

   def dispatch(self, state_change):
       """ Apply the `state_change` in the current machine and return the
       resulting events.

       Args:
           state_change (StateChange): An object representation of a state
           change.

       Return:
           Event: A list of events produced by the state transition, it's
           the upper layer's responsibility to decided how to handle these
           events.
       """
       assert isinstance(state_change, StateChange)

       # the state objects must be treated as immutable, so make a copy of the
       # current state and pass the copy to the state machine to be modified.
       next_state = deepcopy(self.current_state)

       # update the current state by applying the change
       iteration = self.state_transition(
           next_state,
           state_change,
       )

       assert isinstance(iteration, TransitionResult)

       self.current_state, events = iteration

       assert isinstance(self.current_state, (State, type(None)))
       assert all(isinstance(e, Event) for e in events)

       return events

   def __eq__(self, other):
       if isinstance(other, StateManager):
           return (
               self.state_transition == other.state_transition and
               self.current_state == other.current_state
           )
       return False

   def __ne__(self, other):
       return not self.__eq__(other)

*/
