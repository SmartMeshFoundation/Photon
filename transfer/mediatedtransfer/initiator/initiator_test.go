package initiator

import (
	"testing"

	"github.com/SmartMeshFoundation/Photon/params"

	"math/big"

	"os"

	"encoding/json"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/transfer/route"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/Photon/utils/utest"
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
func makeInitStateChange(routes []*route.State, target common.Address, amount *big.Int, blocknumber int64, ourAddress common.Address, token common.Address) *mediatedtransfer.ActionInitInitiatorStateChange {
	tr := &mediatedtransfer.LockedTransferState{
		Amount:       amount,
		Initiator:    ourAddress,
		Target:       target,
		Token:        token,
		TargetAmount: amount,
		Fee:          utils.BigInt0,
	}
	initStateChange := &mediatedtransfer.ActionInitInitiatorStateChange{
		OurAddress:  ourAddress,
		Tranfer:     tr,
		Routes:      route.NewRoutesState(routes),
		BlockNumber: blocknumber,
	}
	initStateChange.Secret = utils.NewRandomHash()
	initStateChange.LockSecretHash = utils.ShaSecret(initStateChange.Secret[:])
	return initStateChange
}
func makeInitiatorState(routes []*route.State, target common.Address, amount *big.Int, blocknumber int64, ourAddress common.Address, token common.Address) (initState *mediatedtransfer.InitiatorState) {
	initStateChange := makeInitStateChange(routes, target, amount, blocknumber, ourAddress, token)
	it := StateTransition(nil, initStateChange)
	initState = it.NewState.(*mediatedtransfer.InitiatorState)
	return initState
}
func TestNextRoute(t *testing.T) {
	target := utest.HOP1
	routes := []*route.State{
		utest.MakeRoute(utest.HOP2, utest.UnitTransferAmount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP3, x.Sub(utest.UnitTransferAmount, big.NewInt(1)), utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP4, utest.UnitTransferAmount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	state := makeInitiatorState(routes, target, utest.UnitTransferAmount, 0, utest.ADDR, utest.UnitTokenAddress)
	assert(t, state.Route, routes[0])
	assert(t, state.Routes.AvailableRoutes, routes[1:])
	assert(t, state.Routes.IgnoredRoutes == nil, true)
	assert(t, state.Routes.RefundedRoutes == nil, true)
	assert(t, state.Routes.CanceledRoutes == nil, true)

	//open this will panic,how to test panic?
	//err := tryNewRoute(state)
	//assert.Equal(t, err != nil, true)
	state.Routes.CanceledRoutes = append(state.Routes.CanceledRoutes, state.Route)
	state.Route = nil
	tryNewRoute(state)
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
	identifier := utils.Sha3([]byte("3"))
	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	initStateChange := makeInitStateChange(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddrses, utest.UnitTokenAddress)
	expiration := blockNumber + int64(utest.Hop1Timeout)
	initiatorStateMachine := transfer.NewStateManager(StateTransition, nil, NameInitiatorTransition, identifier, utils.NewRandomAddress())
	assert(t, initiatorStateMachine.CurrentState, nil)
	events := initiatorStateMachine.Dispatch(initStateChange)
	initiatorState := initiatorStateMachine.CurrentState.(*mediatedtransfer.InitiatorState)
	assert(t, initiatorState.OurAddress, ourAddrses)
	tr := initiatorState.Transfer
	assert(t, tr.Amount, amount)
	assert(t, tr.Target, targetAddress)
	assert(t, tr.LockSecretHash, utils.ShaSecret(tr.Secret[:]))
	assert(t, len(events) > 0, true)

	var mtrs []*mediatedtransfer.EventSendMediatedTransfer
	for _, e := range events {
		if e2, ok := e.(*mediatedtransfer.EventSendMediatedTransfer); ok {
			mtrs = append(mtrs, e2)
		}
	}
	assert(t, len(mtrs), 1)

	mtr := mtrs[0]
	assert(t, mtr.Token, utest.UnitTokenAddress)
	assert(t, mtr.Amount, amount, "transfer amount mismatch")
	assert(t, mtr.Expiration, expiration-int64(params.DefaultRevealTimeout), "transfer expiration mismatch")
	assert(t, mtr.LockSecretHash != utils.EmptyHash, true)
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
	routes := []*route.State{}
	initStateChange := makeInitStateChange(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddrses, utest.UnitTokenAddress)
	initiatorStateMachine := transfer.NewStateManager(StateTransition, nil, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())
	assert(t, initiatorStateMachine.CurrentState, nil)
	events := initiatorStateMachine.Dispatch(initStateChange)

	assert(t, len(events), 2)
	assert(t, initiatorStateMachine.CurrentState, nil)
	_, ok := events[0].(*transfer.EventTransferSentFailed)
	assert(t, ok, true)
}

func TestStateWaitSecretRequestValid(t *testing.T) {
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR

	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, utest.UnitTokenAddress)
	hashlock := currentState.Transfer.LockSecretHash
	stateChange := &mediatedtransfer.ReceiveSecretRequestStateChange{
		Amount:         amount,
		LockSecretHash: hashlock,
		Sender:         targetAddress,
	}
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())
	events := sm.Dispatch(stateChange)
	assert(t, len(events), 1)
	_, ok := events[0].(*mediatedtransfer.EventSendRevealSecret)
	assert(t, ok, true)
}
func TestStateWaitUnlockValid(t *testing.T) {
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, token)
	secret := currentState.Transfer.Secret
	assert(t, secret != utils.EmptyHash, true)

	//setup the state for the wait unlock
	currentState.RevealSecret = &mediatedtransfer.EventSendRevealSecret{
		Secret:   secret,
		Token:    token,
		Receiver: targetAddress,
		Sender:   ourAddress,
	}
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())
	stateChange := &mediatedtransfer.ReceiveSecretRevealStateChange{
		Secret: secret,
		Sender: mediatorAddress,
	}
	events := sm.Dispatch(stateChange)
	assert(t, len(events), 4)
	var EventSendBalanceProof *mediatedtransfer.EventSendBalanceProof
	var EventTransferSentSuccess *transfer.EventTransferSentSuccess
	var EventUnlockSuccess *mediatedtransfer.EventUnlockSuccess
	for _, e := range events {
		switch e2 := e.(type) {
		case *mediatedtransfer.EventSendBalanceProof:
			EventSendBalanceProof = e2
		case *transfer.EventTransferSentSuccess:
			EventTransferSentSuccess = e2
		case *mediatedtransfer.EventUnlockSuccess:
			EventUnlockSuccess = e2
		}
	}
	assert(t, EventSendBalanceProof != nil, true)
	assert(t, EventTransferSentSuccess != nil, true)
	assert(t, EventUnlockSuccess != nil, true)

	assert(t, EventSendBalanceProof.Receiver, mediatorAddress)
	assert(t, sm.CurrentState, nil, "state must be cleaned")
}

func TestStateWaitUnlockInvalid(t *testing.T) {
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, token)
	secret := currentState.Transfer.Secret
	assert(t, secret != utils.EmptyHash, true)

	//setup the state for the wait unlock
	currentState.RevealSecret = &mediatedtransfer.EventSendRevealSecret{
		LockSecretHash: currentState.LockSecretHash,
		Secret:         secret,
		Token:          token,
		Receiver:       targetAddress,
		Sender:         ourAddress,
	}
	var beforeState mediatedtransfer.InitiatorState
	utils.DeepCopy(&beforeState, currentState)
	beforeState.Route.Path = []common.Address{}
	beforeState.Message.Path = []common.Address{}

	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())
	stateChange := &mediatedtransfer.ReceiveSecretRevealStateChange{
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
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP2, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, token)
	stateChange := &mediatedtransfer.ReceiveAnnounceDisposedStateChange{
		Sender:  mediatorAddress,
		Token:   token,
		Message: nil,
		Lock: &mtree.Lock{
			Expiration:     currentState.Transfer.Expiration,
			LockSecretHash: currentState.LockSecretHash,
			Amount:         amount,
		},
	}
	var priorState mediatedtransfer.InitiatorState
	utils.DeepCopy(&priorState, currentState)
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 2)
	_, ok := events[0].(*mediatedtransfer.EventSendMediatedTransfer)
	assert(t, ok, true, "No mediated transfer event emitted, should have tried a new route")
	assert(t, sm.CurrentState != nil, true)
	//assert(t, currentState.Routes.CanceledRoutes[0], priorState.Route)
}
func TestRefundTransferNoMoreRoutes(t *testing.T) {
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, token)
	stateChange := &mediatedtransfer.ReceiveAnnounceDisposedStateChange{
		Sender:  mediatorAddress,
		Token:   token,
		Message: nil,
		Lock: &mtree.Lock{
			Expiration:     currentState.Transfer.Expiration,
			LockSecretHash: currentState.LockSecretHash,
			Amount:         amount,
		},
	}
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 3)
	_, ok := events[0].(*transfer.EventTransferSentFailed)
	assert(t, ok, true)
	assert(t, sm.CurrentState == nil, true)
}
func TestRefundTransferInvalidSender(t *testing.T) {
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP2, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, token)
	stateChange := &mediatedtransfer.ReceiveAnnounceDisposedStateChange{
		Sender:  ourAddress,
		Token:   token,
		Message: nil,
		Lock: &mtree.Lock{
			Expiration:     currentState.Transfer.Expiration,
			LockSecretHash: currentState.LockSecretHash,
			Amount:         amount,
		},
	}
	var priorState mediatedtransfer.InitiatorState
	utils.DeepCopy(&priorState, currentState)
	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 0)
	assert(t, sm.CurrentState != nil, true)
	assertStateEqual(t, currentState, &priorState)
}

func TestCancelTransfer(t *testing.T) {
	amount := utest.UnitTransferAmount
	blockNumber := utest.UnitBlockNumber
	mediatorAddress := utest.HOP1
	targetAddress := utest.HOP2
	ourAddress := utest.ADDR
	token := utest.UnitTokenAddress

	routes := []*route.State{
		utest.MakeRoute(mediatorAddress, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP2, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	currentState := makeInitiatorState(routes, targetAddress, utest.UnitTransferAmount, blockNumber, ourAddress, token)

	stateChange := &transfer.ActionCancelTransferStateChange{
		LockSecretHash: currentState.LockSecretHash,
	}

	sm := transfer.NewStateManager(StateTransition, currentState, NameInitiatorTransition, utils.ShaSecret([]byte("3")), utils.NewRandomAddress())

	events := sm.Dispatch(stateChange)
	assert(t, len(events), 1)
	_, ok := events[0].(*transfer.EventTransferSentFailed)
	assert(t, true, ok)
}

func assertStateEqual(t *testing.T, currentState, beforeState *mediatedtransfer.InitiatorState) {
	//assert(t, reflect.DeepEqual(currentState, beforeState), true)
	assert(t, currentState.Transfer, beforeState.Transfer)
	assert(t, currentState.RevealSecret, beforeState.RevealSecret)
	assert(t, currentState.Message, beforeState.Message)
	//不比较这个是因为gob在处理空数组和nil的时候不一致
	//assert(t, currentState.Routes, beforeState.Routes)
	//私有成员变量 gob 无法编码
	//assert(t, currentState.Route, beforeState.Route)
	assert(t, currentState.CanceledTransfers, beforeState.CanceledTransfers)
	assert(t, currentState.SecretRequest, beforeState.SecretRequest)
	assert(t, currentState.OurAddress, beforeState.OurAddress)
	assert(t, currentState.BlockNumber, beforeState.BlockNumber)
	//assert(t, currentState, beforeState)
}

// marshalIndent :
func marshalIndent(v interface{}) string {
	buf, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
