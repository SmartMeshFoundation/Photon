package initiator

import (
	"testing"

	"math/big"

	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/utils/utest"
	"github.com/ethereum/go-ethereum/common"
	assert2 "github.com/stretchr/testify/assert"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

var x = big.NewInt(0)

func assert(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert2.EqualValues(t, expected, actual, msgAndArgs...)
}
func makeInitStateChange(routes []*transfer.RouteState, target common.Address, amount *big.Int, blocknumber int64, ourAddress common.Address, identifier uint64, token common.Address) *mediated_transfer.ActionInitInitiatorStateChange {
	tr := &mediated_transfer.LockedTransferState{
		Identifier: identifier,
		Amount:     amount,
		Initiator:  ourAddress,
		Target:     target,
		Token:      token,
	}
	initStateChange := &mediated_transfer.ActionInitInitiatorStateChange{
		OurAddress:      ourAddress,
		Tranfer:         tr,
		Routes:          transfer.NewRoutesState(routes),
		RandomGenerator: utils.RandomSecretGenerator,
		BlockNumber:     blocknumber,
	}
	return initStateChange
}
func makeInitiatorState(routes []*transfer.RouteState, target common.Address, amount *big.Int, blocknumber int64, ourAddress common.Address, identifier uint64, token common.Address) (initState *mediated_transfer.InitiatorState) {
	initStateChange := makeInitStateChange(routes, target, amount, blocknumber, ourAddress, identifier, token)
	it := StateTransition(nil, initStateChange)
	initState = it.NewState.(*mediated_transfer.InitiatorState)
	return initState
}
func TestNextRoute(t *testing.T) {
	target := utest.HOP1
	routes := []*transfer.RouteState{
		utest.MakeRoute(utest.HOP2, utest.UnitTransferAmount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP3, x.Sub(utest.UnitTransferAmount, big.NewInt(1)), utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP4, utest.UnitTransferAmount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	state := makeInitiatorState(routes, target, utest.UnitTransferAmount, 0, utest.ADDR, 0, utest.UnitTokenAddress)
	assert(t, state.Route, routes[0])
	assert(t, state.Routes.AvailableRoutes, routes[1:])
	assert(t, state.Routes.IgnoredRoutes == nil, true)
	assert(t, state.Routes.RefundedRoutes == nil, true)
	assert(t, state.Routes.CanceledRoutes == nil, true)

	//open this will panic,how to test panic?
	//err := TryNewRoute(state)
	//assert.Equal(t, err != nil, true)
	state.Routes.CanceledRoutes = append(state.Routes.CanceledRoutes, state.Route)
	state.Route = nil
	TryNewRoute(state)
	/*
	   HOP3 should be ignored because it doesnt have enough balance
	*/
	assert(t, len(state.Routes.IgnoredRoutes), 1)
	assert(t, len(state.Routes.AvailableRoutes), 0)
	assert(t, state.Routes.RefundedRoutes == nil, true)
	assert(t, state.Routes.CanceledRoutes != nil, true)
}

func TestInitWithUsableRoutes(t *testing.T) {
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddrses := utest.ADDR
	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	initStateChange := makeInitStateChange(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddrses, 0, utest.UnitTokenAddress)
	expiration := blockNumber + int64(utest.Hop1Timeout)
	initiatorStateMachine := transfer.NewStateManager(StateTransition, nil, NameInitiatorTransition, 3, utils.NewRandomAddress())
	assert(t, initiatorStateMachine.CurrentState, nil)
	events := initiatorStateMachine.Dispatch(initStateChange)
	initiatorState := initiatorStateMachine.CurrentState.(*mediated_transfer.InitiatorState)
	assert(t, initiatorState.OurAddress, ourAddrses)
	tr := initiatorState.Transfer
	assert(t, tr.Amount, amount)
	assert(t, tr.Target, targetAddress)
	assert(t, tr.Hashlock, utils.Sha3(tr.Secret[:]))
	assert(t, len(events) > 0, true)

	var mtrs []*mediated_transfer.EventSendMediatedTransfer
	for _, e := range events {
		if e2, ok := e.(*mediated_transfer.EventSendMediatedTransfer); ok {
			mtrs = append(mtrs, e2)
		}
	}
	assert(t, len(mtrs), 1)

	mtr := mtrs[0]
	assert(t, mtr.Token, utest.UnitTokenAddress)
	assert(t, mtr.Amount, amount, "transfer amount mismatch")
	assert(t, mtr.Expiration, expiration, "transfer expiration mismatch")
	assert(t, mtr.HashLock != utils.EmptyHash, true)
	assert(t, mtr.Receiver, mediatorAddress)
	assert(t, initiatorState.Route, routes[0])
	assert(t, len(initiatorState.Routes.AvailableRoutes), 0)
	assert(t, len(initiatorState.Routes.IgnoredRoutes), 0)
	assert(t, len(initiatorState.Routes.RefundedRoutes), 0)
	assert(t, len(initiatorState.Routes.CanceledRoutes), 0)
}

func TestInitWithoutRoutes(t *testing.T) {
	blockNumber := utest.UnitBlockNumber
	targetAddress := utest.HOP2
	ourAddrses := utest.ADDR
	routes := []*transfer.RouteState{}
	initStateChange := makeInitStateChange(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddrses, 0, utest.UnitTokenAddress)
	initiatorStateMachine := transfer.NewStateManager(StateTransition, nil, NameInitiatorTransition, 3, utils.NewRandomAddress())
	assert(t, initiatorStateMachine.CurrentState, nil)
	events := initiatorStateMachine.Dispatch(initStateChange)

	assert(t, len(events), 1)
	assert(t, initiatorStateMachine.CurrentState, nil)
	_, ok := events[0].(*transfer.EventTransferSentFailed)
	assert(t, ok, true)
}

func TestStateWaitSecretRequestValid(t *testing.T) {
	identifier := utest.UnitIdentifier
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR

	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, identifier, utest.UnitTokenAddress)
	hashlock := currentState.Transfer.Hashlock
	stateChange := &mediated_transfer.ReceiveSecretRequestStateChange{
		Identifier: identifier,
		Amount:     amount,
		Hashlock:   hashlock,
		Sender:     targetAddress,
	}
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, 3, utils.NewRandomAddress())
	events := sm.Dispatch(stateChange)
	assert(t, len(events), 1)
	_, ok := events[0].(*mediated_transfer.EventSendRevealSecret)
	assert(t, ok, true)
}
func TestStateWaitUnlockValid(t *testing.T) {
	identifier := utest.UnitIdentifier
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, identifier, token)
	secret := currentState.Transfer.Secret
	assert(t, secret != utils.EmptyHash, true)

	//setup the state for the wait unlock
	currentState.RevealSecret = &mediated_transfer.EventSendRevealSecret{
		Identifier: identifier,
		Secret:     secret,
		Token:      token,
		Receiver:   targetAddress,
		Sender:     ourAddress,
	}
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, 3, utils.NewRandomAddress())
	stateChange := &mediated_transfer.ReceiveSecretRevealStateChange{
		Secret: secret,
		Sender: mediatorAddress,
	}
	events := sm.Dispatch(stateChange)
	assert(t, len(events), 3)
	var EventSendBalanceProof *mediated_transfer.EventSendBalanceProof
	var EventTransferSentSuccess *transfer.EventTransferSentSuccess
	var EventUnlockSuccess *mediated_transfer.EventUnlockSuccess
	for _, e := range events {
		switch e2 := e.(type) {
		case *mediated_transfer.EventSendBalanceProof:
			EventSendBalanceProof = e2
		case *transfer.EventTransferSentSuccess:
			EventTransferSentSuccess = e2
		case *mediated_transfer.EventUnlockSuccess:
			EventUnlockSuccess = e2
		}
	}
	assert(t, EventSendBalanceProof != nil, true)
	assert(t, EventTransferSentSuccess != nil, true)
	assert(t, EventUnlockSuccess != nil, true)

	assert(t, EventSendBalanceProof.Receiver, mediatorAddress)
	assert(t, EventTransferSentSuccess.Identifier, identifier)
	assert(t, sm.CurrentState, nil, "state must be cleaned")
}

func TestStateWaitUnlockInvalid(t *testing.T) {
	identifier := utest.UnitIdentifier
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, identifier, token)
	secret := currentState.Transfer.Secret
	assert(t, secret != utils.EmptyHash, true)

	//setup the state for the wait unlock
	currentState.RevealSecret = &mediated_transfer.EventSendRevealSecret{
		Identifier: identifier,
		Secret:     secret,
		Token:      token,
		Receiver:   targetAddress,
		Sender:     ourAddress,
	}
	var beforeState mediated_transfer.InitiatorState
	utils.DeepCopy(&beforeState, currentState)

	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, 3, utils.NewRandomAddress())
	stateChange := &mediated_transfer.ReceiveSecretRevealStateChange{
		Secret: secret,
		Sender: utest.ADDR, //wrong sender
	}
	events := sm.Dispatch(stateChange)
	assert(t, len(events), 0)
	assert(t, currentState.RevealSecret != nil, true)
	assert(t, sm.CurrentState, currentState)
	assertStateEqual(t, currentState, &beforeState)
}
func TestRefundTransferNextRoute(t *testing.T) {
	identifier := utest.UnitIdentifier
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP2, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, identifier, token)
	tr := utest.MakeTransfer(amount, ourAddress, targetAddress, blockNumber+int64(utest.UnitSettleTimeout), utils.EmptyHash, utils.EmptyHash, identifier, utest.UnitTokenAddress)
	stateChange := &mediated_transfer.ReceiveTransferRefundStateChange{
		Sender:   mediatorAddress,
		Transfer: tr,
	}
	var priorState mediated_transfer.InitiatorState
	utils.DeepCopy(&priorState, currentState)
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, 3, utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 1)
	_, ok := events[0].(*mediated_transfer.EventSendMediatedTransfer)
	assert(t, ok, true, "No mediated transfer event emitted, should have tried a new route")
	assert(t, sm.CurrentState != nil, true)
	assert(t, currentState.Routes.CanceledRoutes[0], priorState.Route)
}
func TestRefundTransferNoMoreRoutes(t *testing.T) {
	identifier := utest.UnitIdentifier
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, identifier, token)
	tr := utest.MakeTransfer(amount, ourAddress, targetAddress, blockNumber+int64(utest.UnitSettleTimeout), utils.EmptyHash, utils.EmptyHash, identifier, utest.UnitTokenAddress)
	stateChange := &mediated_transfer.ReceiveTransferRefundStateChange{
		Sender:   mediatorAddress,
		Transfer: tr,
	}
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, 3, utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 1)
	_, ok := events[0].(*transfer.EventTransferSentFailed)
	assert(t, ok, true)
	assert(t, sm.CurrentState == nil, true)
}
func TestRefundTransferInvalidSender(t *testing.T) {
	identifier := utest.UnitIdentifier
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP2, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, identifier, token)
	tr := utest.MakeTransfer(amount, ourAddress, targetAddress, blockNumber+int64(utest.UnitSettleTimeout), utils.EmptyHash, utils.EmptyHash, identifier, utest.UnitTokenAddress)
	stateChange := &mediated_transfer.ReceiveTransferRefundStateChange{
		Sender:   ourAddress,
		Transfer: tr,
	}
	var priorState mediated_transfer.InitiatorState
	utils.DeepCopy(&priorState, currentState)
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, 3, utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 0)
	assert(t, sm.CurrentState != nil, true)
	assertStateEqual(t, currentState, &priorState)
}

func TestCancelTransfer(t *testing.T) {
	identifier := utest.UnitIdentifier
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*transfer.RouteState{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP2, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomAddress()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, identifier, token)

	stateChange := &transfer.ActionCancelTransferStateChange{
		Identifier: identifier,
	}

	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, 3, utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 1)
	_, ok := events[0].(*transfer.EventTransferSentFailed)
	assert(t, ok, true)
	assert(t, sm.CurrentState == nil, true)
}

func assertStateEqual(t *testing.T, currentState, beforeState *mediated_transfer.InitiatorState) {
	//assert(t, reflect.DeepEqual(currentState, beforeState), true)
	assert(t, currentState.Transfer, beforeState.Transfer)
	assert(t, currentState.RevealSecret, beforeState.RevealSecret)
	assert(t, currentState.Message, beforeState.Message)
	//不比较这个是因为gob在处理空数组和nil的时候不一致
	//assert(t, currentState.Routes, beforeState.Routes)
	assert(t, currentState.Route, beforeState.Route)
	assert(t, currentState.CanceledTransfers, beforeState.CanceledTransfers)
	assert(t, currentState.SecretRequest, beforeState.SecretRequest)
	assert(t, currentState.OurAddress, beforeState.OurAddress)
	assert(t, currentState.BlockNumber, beforeState.BlockNumber)
	//assert(t, currentState, beforeState)
}

/*
def assert_state_equal(state1, state2):
    """ Weak equality check between two InitiatorState instances """
    assert state1.__class__ == state2.__class__
    assert state1.our_address == state2.our_address
    assert state1.block_number == state2.block_number
    assert state1.routes == state2.routes
    assert state1.route == state2.route
    assert state1.transfer == state2.transfer
    assert state1.random_generator.secrets == state2.random_generator.secrets
    assert state1.canceled_transfers == state2.canceled_transfers

*/
