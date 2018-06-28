package initiator

import (
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	mt "github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

//NameInitiatorTransition name for state manager
const NameInitiatorTransition = "InitiatorTransition"

/*
""" Clear current state and try a new route.

    - Discards the current secret
    - Add the current route to the canceled list
    - Add the current message to the canceled transfers
    """
*/
func cancelCurrentRoute(state *mt.InitiatorState) *transfer.TransitionResult {
	if state.RevealSecret != nil {
		panic("cannot cancel a transfer with a RevealSecret in flight")
	}
	if state.Route != nil {
		state.LastRefundChannelAddress = state.Route.ChannelAddress
	}
	state.Routes.CanceledRoutes = append(state.Routes.CanceledRoutes, state.Route)
	state.CanceledTransfers = append(state.CanceledTransfers, state.Message)
	state.Transfer.Secret = utils.EmptyHash
	state.Message = nil
	state.Route = nil
	state.SecretRequest = nil

	return tryNewRoute(state)
}

//Cancel the current in-transit message
func userCancelTransfer(state *mt.InitiatorState) *transfer.TransitionResult {
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
		Identifier:     state.Transfer.Identifier,
		Reason:         "user canceled transfer",
		Target:         state.Transfer.Target,
		ChannelAddress: state.LastRefundChannelAddress,
	}
	return &transfer.TransitionResult{
		NewState: nil,
		Events:   []transfer.Event{cancel},
	}
}

func tryNewRoute(state *mt.InitiatorState) *transfer.TransitionResult {
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
	var tryRoute *transfer.RouteState
	for len(state.Routes.AvailableRoutes) > 0 {
		route := state.Routes.AvailableRoutes[0]
		if state.Db != nil {
			ch, err := state.Db.GetChannelByAddress(route.ChannelAddress)
			if err != nil {
				log.Error(fmt.Sprintf("get channle %s status from db err %s", utils.APex(route.ChannelAddress), err))
			} else {
				route.AvaibleBalance = ch.OurBalance.Sub(ch.OurBalance, ch.OurAmountLocked)
				route.State = ch.State
				route.ClosedBlock = ch.ClosedBlock
			}
		} else {
			log.Error(" db is nil can only be ignored when you are run testing...")
		}

		state.Routes.AvailableRoutes = state.Routes.AvailableRoutes[1:]
		if route.AvaibleBalance.Cmp(new(big.Int).Add(state.Transfer.TargetAmount, route.Fee)) < 0 {
			state.Routes.IgnoredRoutes = append(state.Routes.IgnoredRoutes, route)
		} else {
			tryRoute = route
			break
		}
	}
	var unlockFailed *mt.EventUnlockFailed
	if state.Message != nil { //目前无论是发起时还是取消路由，Message都会被设置为nil，所以这个事件永远也不会发生。
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
			Identifier:     state.Transfer.Identifier,
			Reason:         "no route available",
			Target:         state.Transfer.Target,
			ChannelAddress: state.LastRefundChannelAddress,
		}
		events := []transfer.Event{transferFailed}
		if unlockFailed != nil {
			events = append(events, unlockFailed)
		}
		var newState *mt.InitiatorState
		if len(state.CanceledTransfers) != 0 {
			newState = state //make sure not finish
		}
		return &transfer.TransitionResult{
			NewState: newState,
			Events:   events,
		}
	}
	state.Route = tryRoute
	secret, hashlock := state.RandomGenerator()
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
	if lockExpiration > state.Transfer.Expiration && state.Transfer.Expiration != 0 {
		lockExpiration = state.Transfer.Expiration
	}
	tr := &mt.LockedTransferState{
		Identifier:   state.Transfer.Identifier,
		TargetAmount: state.Transfer.TargetAmount,
		Amount:       new(big.Int).Add(state.Transfer.TargetAmount, tryRoute.TotalFee),
		Token:        state.Transfer.Token,
		Initiator:    state.Transfer.Initiator,
		Target:       state.Transfer.Target,
		Expiration:   lockExpiration,
		Hashlock:     hashlock,
		Secret:       secret,
		Fee:          tryRoute.TotalFee,
	}
	msg := mt.NewEventSendMediatedTransfer(tr, tryRoute.HopNode)
	state.Transfer = tr
	state.Message = msg
	log.Trace(fmt.Sprintf("send mediated transfer id=%d,amount=%s,token=%s,target=%s,secret=%s", tr.Identifier, tr.Amount, utils.APex(tr.Token), utils.APex(tr.Target), tr.Secret.String()))
	events := []transfer.Event{msg}
	if unlockFailed != nil {
		events = append(events, unlockFailed)
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}
func expiredHashLockEvents(state *mt.InitiatorState) (events []transfer.Event) {
	if state.BlockNumber > state.Transfer.Expiration {
		if state.Route != nil && !state.Db.IsThisLockRemoved(state.Route.ChannelAddress, state.OurAddress, state.Transfer.Hashlock) {
			unlockFailed := &mt.EventUnlockFailed{
				Identifier:     state.Transfer.Identifier,
				Hashlock:       state.Transfer.Hashlock,
				ChannelAddress: state.Route.ChannelAddress,
				Reason:         "lock expired",
			}
			events = append(events, unlockFailed)
		}
	}
	for i, tr := range state.CanceledTransfers {
		route := state.Routes.CanceledRoutes[i]
		if state.BlockNumber > tr.Expiration && !state.Db.IsThisLockRemoved(route.ChannelAddress, state.OurAddress, tr.HashLock) {
			unlockFailed := &mt.EventUnlockFailed{
				Identifier:     tr.Identifier,
				Hashlock:       tr.HashLock,
				ChannelAddress: route.ChannelAddress,
				Reason:         "lock expired",
			}
			events = append(events, unlockFailed)
		}
	}
	return
}

/*
make sure not call this when transfer already finished , state is nil means finished.
*/
func handleBlock(state *mt.InitiatorState, stateChange *transfer.BlockStateChange) *transfer.TransitionResult {
	if state.BlockNumber < stateChange.BlockNumber {
		state.BlockNumber = stateChange.BlockNumber
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   expiredHashLockEvents(state),
	}
}

func handleRouteChange(state *mt.InitiatorState, stateChange *transfer.ActionRouteChangeStateChange) *transfer.TransitionResult {
	mt.UpdateRoute(state.Routes, stateChange)
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}
func handleTransferRefund(state *mt.InitiatorState, stateChange *mt.ReceiveTransferRefundStateChange) *transfer.TransitionResult {
	if stateChange.Sender == state.Route.HopNode && mediator.IsValidRefund(state.Transfer, stateChange.Transfer, state.Route.HopNode, stateChange.Sender) {
		return cancelCurrentRoute(state)
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

func handleCancelRoute(state *mt.InitiatorState, stateChange *mt.ActionCancelRouteStateChange) *transfer.TransitionResult {
	if stateChange.Identifier == state.Transfer.Identifier {
		return cancelCurrentRoute(state)
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

func handleCancelTransfer(state *mt.InitiatorState) *transfer.TransitionResult {
	return userCancelTransfer(state)
}

func handleSecretRequest(state *mt.InitiatorState, stateChange *mt.ReceiveSecretRequestStateChange) *transfer.TransitionResult {
	isValid := stateChange.Sender == state.Transfer.Target &&
		stateChange.Hashlock == state.Transfer.Hashlock &&
		stateChange.Identifier == state.Transfer.Identifier &&
		stateChange.Amount.Cmp(state.Transfer.TargetAmount) == 0
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
		return &transfer.TransitionResult{
			NewState: state,
			Events:   []transfer.Event{revealSecret},
		}
	} else if isInvalid {
		return cancelCurrentRoute(state)
	} else {
		return &transfer.TransitionResult{
			NewState: state,
			Events:   nil,
		}
	}
}

/*
Send a balance proof to the next hop with the current mediated transfer
    lock removed and the balance updated.
*/
func handleSecretReveal(state *mt.InitiatorState, st *mt.ReceiveSecretRevealStateChange) *transfer.TransitionResult {
	/*
		考虑到崩溃恢复情形,可能崩溃了很久. 如果这时候交易还继续进行,显然不合理.
	*/
	if state.BlockNumber >= state.Transfer.Expiration {
		return &transfer.TransitionResult{
			NewState: state,
			Events:   nil,
		}
	}
	if state.Route != nil && state.Transfer != nil && st.Sender == state.Route.HopNode && st.Secret == state.Transfer.Secret {
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
			Identifier:     tr.Identifier,
			Amount:         tr.Amount,
			Target:         tr.Target,
			ChannelAddress: state.Route.ChannelAddress,
		}
		unlockSuccess := &mt.EventUnlockSuccess{
			Identifier: tr.Identifier,
			Hashlock:   tr.Hashlock,
		}
		return &transfer.TransitionResult{
			NewState: nil,
			Events:   []transfer.Event{unlockLock, transferSuccess, unlockSuccess},
		}
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

/*
StateTransition is State machine for a node starting a mediated transfer.
    originalState: The current State that is transitioned from.
    st: The state_change that will be applied.
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
	it := &transfer.TransitionResult{
		NewState: originalState,
		Events:   nil,
	}
	state, ok := originalState.(*mt.InitiatorState)
	if !ok {
		if originalState != nil {
			panic("InitiatorState StateTransition get type error")
		}
		state = nil //originalState is nil
	}
	if state == nil {
		staii, ok := st.(*mt.ActionInitInitiatorStateChange)
		if ok {
			var routes transfer.RoutesState
			err := utils.DeepCopy(&routes, staii.Routes)
			if err != nil {
				panic(fmt.Sprintf("deepcopy error:%#v", err))
			}
			state = &mt.InitiatorState{
				OurAddress:      staii.OurAddress,
				Transfer:        staii.Tranfer,
				Routes:          &routes,
				BlockNumber:     staii.BlockNumber,
				RandomGenerator: staii.RandomGenerator,
				Db:              staii.Db,
			}
			return tryNewRoute(state)
		}
		//todo fix, find a way to remove this identifier from raiden.Identifier2StateManagers
		//log.Warn(fmt.Sprintf("originalState,statechange should not be here originalState=\n%s\n,statechange=\n%s",
		//	utils.StringInterface1(originalState), utils.StringInterface1(st)))
	} else if state.RevealSecret == nil {
		switch st2 := st.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
			//目前没用,
		case *transfer.ActionRouteChangeStateChange:
			it = handleRouteChange(state, st2)
		case *mt.ReceiveSecretRequestStateChange:
			it = handleSecretRequest(state, st2)
		case *mt.ReceiveTransferRefundStateChange:
			it = handleTransferRefund(state, st2)
			//目前没用
		case *mt.ActionCancelRouteStateChange:
			it = handleCancelRoute(state, st2)
			//目前也没用
		case *transfer.ActionCancelTransferStateChange:
			it = handleCancelTransfer(state)
		case *mt.ReceiveSecretRevealStateChange:
			//只要密码正确,就应该发送secret ,流程上可能有问题,但是结果是没错的(只有在token swap的时候才会走到这一步) . 因为按照协议层要求,同一个消息不会重复发送, 导致在tokenswap的时候maker不可能重复发送reveal secret
			if st2.Secret == state.Transfer.Secret {
				log.Warn(fmt.Sprintf("send balance proof before send a reveal secret message, this is only for token swap taker,state=%s", utils.StringInterface(state, 3)))
			}
			it = handleSecretReveal(state, st2)
		default:
			log.Warn(fmt.Sprintf("RevealSecret is nil,cannot handle %s", utils.StringInterface(st, 3)))
		}
	} else {
		switch st2 := st.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
		case *mt.ReceiveSecretRevealStateChange:
			it = handleSecretReveal(state, st2)
		default:
			log.Warn(fmt.Sprintf("RevealSecret is not nil,cannot handle %s", utils.StringInterface(st, 3)))
		}
	}
	return it
}
