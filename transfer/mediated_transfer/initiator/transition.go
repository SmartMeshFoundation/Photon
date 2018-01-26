package initiator

import (
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/transfer"
	mt "github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/log"
)

const NameInitiatorTransition = "InitiatorTransition"

/*
""" Clear current state and try a new route.

    - Discards the current secret
    - Add the current route to the canceled list
    - Add the current message to the canceled transfers
    """
*/
func CancelCurrentRoute(state *mt.InitiatorState) *transfer.TransitionResult {
	if state.RevealSecret != nil {
		panic("cannot cancel a transfer with a RevealSecret in flight")
	}
	state.Routes.CanceledRoutes = append(state.Routes.CanceledRoutes, state.Route)
	state.CanceledTransfers = append(state.CanceledTransfers, state.Message)
	state.Transfer.Secret = utils.EmptyHash
	state.Message = nil
	state.Route = nil
	state.SecretRequest = nil

	return TryNewRoute(state)
}

//Cancel the current in-transit message
func UserCancelTransfer(state *mt.InitiatorState) *transfer.TransitionResult {
	if state.RevealSecret != nil {
		panic("cannot cancel a transfer with a RevealSecret in flight")
	}
	state.Transfer.Secret = utils.EmptyHash
	state.Transfer.Hashlock = utils.EmptyHash
	state.Message = nil
	state.Route = nil
	state.SecretRequest = nil
	state.RevealSecret = nil
	cancel := &transfer.EventTransferSentFailed{
		Identifier: state.Transfer.Identifier,
		Reason:     "user canceled transfer",
	}
	return &transfer.TransitionResult{
		NewState: nil,
		Events:   []transfer.Event{cancel},
	}
}

func TryNewRoute(state *mt.InitiatorState) *transfer.TransitionResult {
	if state.Route != nil {
		panic("cannot try a new route while one is being used")
	}
	/*
			  TODO:
		     - Route ranking. An upper layer should rate each route to optimize
		       the fee price/quality of each route and add a rate from in the range
		       [0.0,1.0].
		     - Add in a policy per route:
		       - filtering, e.g. so the user may have a per route maximum transfer
		         value based on fixed value or reputation.
		       - reveal time computation
		       - These policy details are better hidden from this implementation and
		         changes should be applied through the use of Route state changes.

		     Find a single route that may fulfill the request, this uses a single
		     route intentionally
	*/
	var tryRoute *transfer.RouteState = nil
	for len(state.Routes.AvailableRoutes) > 0 {
		route := state.Routes.AvailableRoutes[0]
		state.Routes.AvailableRoutes = state.Routes.AvailableRoutes[1:]
		if route.AvaibleBalance.Cmp(state.Transfer.Amount) < 0 {
			state.Routes.IgnoredRoutes = append(state.Routes.IgnoredRoutes, route)
		} else {
			tryRoute = route
			break
		}
	}
	var unlockFailed *mt.EventUnlockFailed = nil
	if state.Message != nil {
		unlockFailed = &mt.EventUnlockFailed{
			Identifier: state.Transfer.Identifier,
			Hashlock:   state.Transfer.Hashlock,
			Reason:     "route was canceled",
		}
	}
	if tryRoute == nil {
		/*
					 No available route has sufficient balance for the current transfer,
			         cancel it.

			         At this point we can just discard all the state data, this is only
			         valid because we are the initiator and we know that the secret was
			         not released.
		*/
		transferFailed := &transfer.EventTransferSentFailed{
			Identifier: state.Transfer.Identifier,
			Reason:     "no route available",
		}
		events := []transfer.Event{transferFailed}
		if unlockFailed != nil {
			events = append(events, unlockFailed)
		}
		return &transfer.TransitionResult{
			NewState: nil,
			Events:   events,
		}
	} else {
		state.Route = tryRoute
		secret := state.RandomGenerator()
		hashlock := utils.Sha3(secret[:])
		/*
					  The initiator doesn't need to learn the secret, so there is no need
			         to decrement reveal_timeout from the lock timeout.

			         The lock_expiration could be set to a value larger than
			         settle_timeout, this is not useful since the next hop will take this
			         channel settle_timeout as an upper limit for expiration.

			         The two nodes will most likely disagree on latest block, as far as
			         the expiration goes this is no problem.
		*/
		lockExpiration := state.BlockNumber + int64(tryRoute.SettleTimeout)
		tr := &mt.LockedTransferState{
			Identifier: state.Transfer.Identifier,
			Amount:     new(big.Int).Set(state.Transfer.Amount),
			Token:      state.Transfer.Token,
			Initiator:  state.Transfer.Initiator,
			Target:     state.Transfer.Target,
			Expiration: lockExpiration,
			Hashlock:   hashlock,
			Secret:     secret,
		}
		msg := mt.NewEventSendMediatedTransfer(tr, tryRoute.HopNode)
		//& mt.EventSendMediatedTransfer{
		//	Identifier: tr.Identifier,
		//	Token:      tr.Token,
		//	Amount:     tr.Amount,
		//	HashLock:   tr.Hashlock,
		//	Initiator:  state.OurAddress,
		//	Target:     tr.Target,
		//	Expiration: lockExpiration,
		//	Receiver:   tryRoute.HopNode,
		//}
		state.Transfer = tr
		state.Message = msg
		events := []transfer.Event{msg}
		if unlockFailed != nil {
			events = append(events, unlockFailed)
		}
		return &transfer.TransitionResult{
			NewState: state,
			Events:   events,
		}
	}
}

func HandleBlock(state *mt.InitiatorState, stateChange *transfer.BlockStateChange) *transfer.TransitionResult {
	if state.BlockNumber < stateChange.BlockNumber {
		state.BlockNumber = stateChange.BlockNumber
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

func HandleRouteChange(state *mt.InitiatorState, stateChange *transfer.ActionRouteChangeStateChange) *transfer.TransitionResult {
	mt.UpdateRoute(state.Routes, stateChange)
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

func HandleTransferRefund(state *mt.InitiatorState, stateChange *mt.ReceiveTransferRefundStateChange) *transfer.TransitionResult {
	if stateChange.Sender == state.Route.HopNode {
		return CancelCurrentRoute(state)
	} else {
		return &transfer.TransitionResult{state, nil}
	}
}

func HandleCancelRoute(state *mt.InitiatorState, stateChange *mt.ActionCancelRouteStateChange) *transfer.TransitionResult {
	if stateChange.Identifier == state.Transfer.Identifier {
		return CancelCurrentRoute(state)
	} else {
		return &transfer.TransitionResult{state, nil}
	}
}

func HandleCancelTransfer(state *mt.InitiatorState) *transfer.TransitionResult {
	return UserCancelTransfer(state)
}

func HandleSecretRequest(state *mt.InitiatorState, stateChange *mt.ReceiveSecretRequestStateChange) *transfer.TransitionResult {
	isValid := stateChange.Sender == state.Transfer.Target &&
		stateChange.Hashlock == state.Transfer.Hashlock &&
		stateChange.Identifier == state.Transfer.Identifier &&
		stateChange.Amount.Cmp(state.Transfer.Amount) == 0
	isInvalid := stateChange.Sender == state.Transfer.Target &&
		stateChange.Hashlock == state.Transfer.Hashlock && !isValid
	if isValid {
		/*
		   Reveal the secret to the target node and wait for its confirmation,
		   at this point the transfer is not cancellable anymore either the lock
		   timeouts or a secret reveal is received.

		   Note: The target might be the first hop

		*/
		tr := state.Transfer
		revealSecret := &mt.EventSendRevealSecret{
			Identifier: tr.Identifier,
			Secret:     tr.Secret,
			Token:      tr.Token,
			Receiver:   tr.Target,
			Sender:     state.OurAddress,
		}
		state.RevealSecret = revealSecret
		return &transfer.TransitionResult{state, []transfer.Event{revealSecret}}
	} else if isInvalid {
		return CancelCurrentRoute(state)
	} else {
		return &transfer.TransitionResult{state, nil}
	}
	return nil
}

/*
Send a balance proof to the next hop with the current mediated transfer
    lock removed and the balance updated.
*/
func HandleSecretReveal(state *mt.InitiatorState, st *mt.ReceiveSecretRevealStateChange) *transfer.TransitionResult {
	if st.Sender == state.Route.HopNode {
		/*
					   next hop learned the secret, unlock the token locally and send the
			         withdraw message to next hop
		*/
		tr := state.Transfer
		unlockLock := &mt.EventSendBalanceProof{
			Identifier:     tr.Identifier,
			ChannelAddress: state.Route.ChannelAddress,
			Token:          tr.Token,
			Receiver:       state.Route.HopNode,
			Secret:         tr.Secret,
		}
		transferSuccess := &transfer.EventTransferSentSuccess{
			Identifier: tr.Identifier,
			Amount:     tr.Amount,
			Target:     tr.Target,
		}
		unlockSuccess := &mt.EventUnlockSuccess{
			Identifier: tr.Identifier,
			Hashlock:   tr.Hashlock,
		}
		return &transfer.TransitionResult{
			NewState: nil,
			Events:   []transfer.Event{unlockLock, transferSuccess, unlockSuccess},
		}
	} else {
		return &transfer.TransitionResult{state, nil}
	}
}

/*
	State machine for a node starting a mediated transfer.

    Args:
        state: The current State that is transitioned from.
        state_change: The state_change that will be applied.
*/
func StateTransition(originalState transfer.State, st transfer.StateChange) *transfer.TransitionResult {
	/*
	   	 TODO: Add synchronization for expired locks.
	        Transfers added to the canceled list by an ActionCancelRoute are stale in
	        the channels merkle tree, while this doesn't increase the messages sizes
	        nor does it interfere with the guarantees of finality it increases memory
	        usage for each end, since the full merkle tree must be saved to compute
	        it's root.
	*/
	it := &transfer.TransitionResult{originalState, nil}
	state, ok := originalState.(*mt.InitiatorState)
	if !ok {
		if originalState != nil {
			panic("InitiatorState StateTransition get type error")
		}
		state = nil //originalState is nil
	}
	if state == nil { //这里可能是一个坑,关于interface和nil比较的问题
		staii, ok := st.(*mt.ActionInitInitiatorStateChange)
		if ok {
			var routes transfer.RoutesState
			err := utils.DeepCopy(&routes, staii.Routes)
			if err != nil {
				panic(fmt.Sprintf("deepcopy error:#v", err))
			}
			state = &mt.InitiatorState{
				OurAddress:      staii.OurAddress,
				Transfer:        staii.Tranfer,
				Routes:          &routes,
				BlockNumber:     staii.BlockNumber,
				RandomGenerator: staii.RandomGenerator,
			}
			return TryNewRoute(state)
		} else {
			//todo fix, find a way to remove this identifier from raiden.Identifier2StateManagers
			log.Warn(fmt.Sprintf("originalState,statechange should not be here originalState=\n%s\n,statechange=\n%s",
				utils.StringInterface1(originalState), utils.StringInterface1(st)))
		}
	} else if state.RevealSecret == nil {
		switch st2 := st.(type) {
		case *transfer.BlockStateChange:
			it = HandleBlock(state, st2)
			//目前没用,
		case *transfer.ActionRouteChangeStateChange:
			it = HandleRouteChange(state, st2)
		case *mt.ReceiveSecretRequestStateChange:
			it = HandleSecretRequest(state, st2)
		case *mt.ReceiveTransferRefundStateChange:
			it = HandleTransferRefund(state, st2)
			//目前没用
		case *mt.ActionCancelRouteStateChange:
			it = HandleCancelRoute(state, st2)
			//目前也没用
		case *transfer.ActionCancelTransferStateChange:
			it = HandleCancelTransfer(state)
		default:
			log.Warn(fmt.Sprintf("RevealSecret is nil,cannot handle %s", utils.StringInterface(st, 3)))
		}
	} else {
		switch st2 := st.(type) {
		case *transfer.BlockStateChange:
			it = HandleBlock(state, st2)
		case *mt.ReceiveSecretRevealStateChange:
			it = HandleSecretReveal(state, st2)
		default:
			log.Warn(fmt.Sprintf("RevealSecret is not nil,cannot handle %s", utils.StringInterface(st, 3)))
		}
	}
	return it
}
