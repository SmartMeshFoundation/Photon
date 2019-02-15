package initiator

import (
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/transfer"
	mt "github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/Photon/transfer/route"
	"github.com/SmartMeshFoundation/Photon/utils"
)

//NameInitiatorTransition name for state manager
const NameInitiatorTransition = "InitiatorTransition"

/*
Clear current state and try a new route.

- Discards the current secret
- Add the current route to the canceled list
- Add the current message to the canceled transfers
*/
func cancelCurrentRoute(state *mt.InitiatorState) *transfer.TransitionResult {
	if state.RevealSecret != nil {
		panic("cannot cancel a transfer with a RevealSecret in flight")
	}
	state.Routes.CanceledRoutes = append(state.Routes.CanceledRoutes, state.Route)
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
	//state.Transfer.LockSecretHash = utils.EmptyHash // need by remove
	state.Message = nil
	//state.Route = nil // need by remove
	state.SecretRequest = nil
	state.RevealSecret = nil
	cancel := &transfer.EventTransferSentFailed{
		LockSecretHash: state.Transfer.LockSecretHash,
		Reason:         "user canceled transfer",
		Target:         state.Transfer.Target,
		Token:          state.Transfer.Token,
	}
	/*
		need state exist to send remove msg after expired
	*/
	return &transfer.TransitionResult{
		NewState: state,
		Events:   []transfer.Event{cancel},
	}
}

func tryNewRoute(state *mt.InitiatorState) *transfer.TransitionResult {
	if state.Route != nil {
		panic("cannot try a new route while one is being used")
	}
	var tryRoute *route.State
	for len(state.Routes.AvailableRoutes) > 0 {
		r := state.Routes.AvailableRoutes[0]
		state.Routes.AvailableRoutes = state.Routes.AvailableRoutes[1:]
		//if !r.CanTransfer() /*交易发起方不应该考虑收费*/ || r.AvailableBalance().Cmp(new(big.Int).Add(state.Transfer.TargetAmount, r.Fee)) < 0 {
		if !r.CanTransfer() || r.AvailableBalance().Cmp(state.Transfer.TargetAmount) < 0 {
			state.Routes.IgnoredRoutes = append(state.Routes.IgnoredRoutes, r)
		} else {
			tryRoute = r
			break
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
			LockSecretHash: state.Transfer.LockSecretHash,
			Reason:         "no route available",
			Target:         state.Transfer.Target,
			Token:          state.Transfer.Token,
		}
		events := []transfer.Event{transferFailed}
		removeManager := &mt.EventRemoveStateManager{
			Key: utils.Sha3(state.LockSecretHash[:], state.Transfer.Token[:]),
		}
		events = append(events, removeManager)
		return &transfer.TransitionResult{
			NewState: nil,
			Events:   events,
		}
	}
	/*
				  The initiator doesn't need to learn the secret, so there is no need
		         to decrement reveal_timeout from the lock timeout.

		         The lock_expiration could be set to a value larger than
		         settle_timeout, this is not useful since the next hop will take this
		         channel settle_timeout as an upper limit for expiration.

		         The two nodes will most likely disagree on latest block, as far as
		         the expiration goes this is no problem.
	*/
	lockExpiration := state.BlockNumber + int64(tryRoute.SettleTimeout()) - int64(params.DefaultRevealTimeout) // - revealTimeout for test
	if lockExpiration > state.Transfer.Expiration && state.Transfer.Expiration != 0 {
		lockExpiration = state.Transfer.Expiration
	}
	tr := &mt.LockedTransferState{
		TargetAmount:   state.Transfer.TargetAmount,
		Amount:         new(big.Int).Add(state.Transfer.TargetAmount, tryRoute.TotalFee),
		Token:          state.Transfer.Token,
		Initiator:      state.Transfer.Initiator,
		Target:         state.Transfer.Target,
		Expiration:     lockExpiration,
		LockSecretHash: state.LockSecretHash,
		Secret:         state.Secret,
		Fee:            tryRoute.TotalFee,
		Data:           state.Transfer.Data,
	}
	msg := mt.NewEventSendMediatedTransfer(tr, tryRoute.HopNode())
	if len(state.Routes.CanceledRoutes) > 0 {
		/*
			保存上次尝试的路由信息,否则当发起方收到AnnounceDisposed的时候,尝试新路由时,会出现异常
		*/
		// Store route info of previous one, or when receiving AnnounceDisposed message and try a new route, error occurs.
		msg.FromChannel = state.Routes.CanceledRoutes[len(state.Routes.CanceledRoutes)-1].ChannelIdentifier
	}
	state.Route = tryRoute
	state.Transfer = tr
	state.Message = msg
	log.Trace(fmt.Sprintf("send mediated transfer id=%s,amount=%s,token=%s,target=%s,secret=%s,data=%s", utils.HPex(tr.LockSecretHash), tr.Amount, utils.APex(tr.Token), utils.APex(tr.Target), tr.Secret.String(), tr.Data))
	events := []transfer.Event{msg}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}
func expiredHashLockEvents(state *mt.InitiatorState) (events []transfer.Event) {
	if state.BlockNumber-params.ForkConfirmNumber > state.Transfer.Expiration {
		if state.Route != nil && !state.Db.IsThisLockRemoved(state.Route.ChannelIdentifier, state.OurAddress, state.Transfer.LockSecretHash) {
			unlockFailed := &mt.EventUnlockFailed{
				LockSecretHash:    state.Transfer.LockSecretHash,
				ChannelIdentifier: state.Route.ChannelIdentifier,
				Reason:            "lock expired",
			}
			transferFailed := &transfer.EventTransferSentFailed{
				LockSecretHash: state.Transfer.LockSecretHash,
				Reason:         "no route available",
				Target:         state.Transfer.Target,
				Token:          state.Transfer.Token,
			}
			events = append(events, unlockFailed, transferFailed)
		}
	}
	return
}

/*
make sure not call this when transfer already finished , state is nil means finished.
*/
func handleBlock(state *mt.InitiatorState, stateChange *transfer.BlockStateChange) *transfer.TransitionResult {
	var events []transfer.Event
	if state.BlockNumber < stateChange.BlockNumber {
		state.BlockNumber = stateChange.BlockNumber
	}
	// 考虑到分叉攻击,延迟一定块数之后才发送remove
	if state.BlockNumber-params.ForkConfirmNumber > state.Transfer.Expiration {
		// 超时
		// 如果我没有发送过密码,直接发送remove expired lock,然后移除state manager
		// 如果我已经发送过密码,那么超时说明我没有收到reveal secret 或 链上密码注册事件,此时我认为交易超时失败,可不可以发送RemoveExpiredHashlock,由通道自己决定.然后移除state manager
		// timeout
		// If I have not sent secret, then just send removeExpiredLock, and remove stateManager.
		// If I have already sent secret, then assume transfer timeout failure, send remove expired, and remove state manager.
		events = expiredHashLockEvents(state)
		events = append(events, &mt.EventRemoveStateManager{
			Key: utils.Sha3(state.LockSecretHash[:], state.Transfer.Token[:]),
		})
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

func handleRefund(state *mt.InitiatorState, stateChange *mt.ReceiveAnnounceDisposedStateChange) *transfer.TransitionResult {
	if mediator.IsValidRefund(state.Transfer, state.Route, stateChange) {
		it := cancelCurrentRoute(state)
		ev := &mt.EventSendAnnounceDisposedResponse{
			LockSecretHash: stateChange.Lock.LockSecretHash,
			Token:          state.Transfer.Token,
			Receiver:       stateChange.Sender,
		}
		it.Events = append(it.Events, ev)
		return it
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

func handleCancelRoute(state *mt.InitiatorState, stateChange *mt.ActionCancelRouteStateChange) *transfer.TransitionResult {
	if stateChange.LockSecretHash == state.Transfer.LockSecretHash {
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
		stateChange.LockSecretHash == state.Transfer.LockSecretHash &&
		stateChange.Amount.Cmp(state.Transfer.TargetAmount) == 0
	//如果收到secret request时候已经过期了,应该让这个交易失败,而不是告诉对方密码
	if isValid && !state.CancelByExceptionSecretRequest && state.BlockNumber < state.Transfer.Expiration {
		/*
		   Reveal the secret to the target node and wait for its confirmation,
		   at this point the transfer is not cancellable anymore either the lock
		   timeouts or a secret reveal is received.

		   Note: The target might be the first hop

		*/
		tr := state.Transfer
		revealSecret := &mt.EventSendRevealSecret{
			LockSecretHash: tr.LockSecretHash,
			Secret:         tr.Secret,
			Token:          tr.Token,
			Receiver:       tr.Target,
			Sender:         state.OurAddress,
			Data:           tr.Data,
		}
		state.RevealSecret = revealSecret
		return &transfer.TransitionResult{
			NewState: state,
			Events:   []transfer.Event{revealSecret},
		}
	}
	/*
		如果收到无效的secret request,当成有人做出了攻击行为来处理,设置标志位拒绝后续所有secret request,
		然后什么都不做,等待remove
	*/
	state.CancelByExceptionSecretRequest = true
	/*
		BUG : 每次交易密码不会发生变化,如果尝试其他路径,可能会被恶意利用
	*/
	//if isInvalid {
	//	return cancelCurrentRoute(state)
	//}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

/*
密码在链上注册了,只要在有效期范围内,就相当于收到了对方的 reveal secret, 主动给对方发送 unlock 消息.
*/
/*
 *	handleSecretRevealOnChain : function to handle event of RevealSecretOnChain.
 *
 *	Note : Once the secret has been registered on chain, all nodes act like they receives reveal secret from their partner,
 *	then send unlock to their partner.
 */
func handleSecretRevealOnChain(state *mt.InitiatorState, st *mt.ContractSecretRevealOnChainStateChange) *transfer.TransitionResult {
	if st.LockSecretHash != state.LockSecretHash {
		//无论是不是 token swap, 都应该知道 locksecrethash,否则肯定是实现有问题
		// we should know locksecrethash no matter whether it is token swap, otherwise implementation has problem.
		panic(fmt.Sprintf("my locksecrethash=%s,received=%s", state.LockSecretHash.String(), st.LockSecretHash.String()))
	}
	log.Trace(fmt.Sprintf("Check lock's expiration, state.Transfer.Expiration=%d, st.BlockNumber=%d\n", state.Transfer.Expiration, st.BlockNumber))
	if state.Transfer.Expiration < st.BlockNumber {
		//对于我来说这笔交易已经超期了. 应该发出 移除此锁消息.
		// As to me this transfer expired, should send RemoveExpiredLock message.
		events := expiredHashLockEvents(state)
		events = append(events, &mt.EventRemoveStateManager{
			Key: utils.Sha3(state.LockSecretHash[:], state.Transfer.Token[:]),
		})
		return &transfer.TransitionResult{
			NewState: state,
			Events:   events,
		}
	}
	//认为交易成功了
	// assume transfer succeed.
	return &transfer.TransitionResult{
		NewState: state,
		Events:   transferSuccessEvents(state),
	}
}

func transferSuccessEvents(state *mt.InitiatorState) (events []transfer.Event) {
	tr := state.Transfer
	unlockLock := &mt.EventSendBalanceProof{
		LockSecretHash:    tr.LockSecretHash,
		ChannelIdentifier: state.Route.ChannelIdentifier,
		Token:             tr.Token,
		Receiver:          state.Route.HopNode(),
	}
	transferSuccess := &transfer.EventTransferSentSuccess{
		LockSecretHash:    tr.LockSecretHash,
		Amount:            tr.Amount,
		Target:            tr.Target,
		ChannelIdentifier: state.Route.ChannelIdentifier,
		Token:             tr.Token,
		Data:              tr.Data,
	}
	unlockSuccess := &mt.EventUnlockSuccess{
		LockSecretHash: tr.LockSecretHash,
	}
	removeManager := &mt.EventRemoveStateManager{
		Key: utils.Sha3(tr.LockSecretHash[:], tr.Token[:]),
	}
	events = []transfer.Event{unlockLock, transferSuccess, unlockSuccess, removeManager}
	return events
}

/*
Send a balance proof to the next hop with the current mediated transfer
    lock removed and the balance updated.
*/
func handleSecretReveal(state *mt.InitiatorState, st *mt.ReceiveSecretRevealStateChange) *transfer.TransitionResult {
	/*
		考虑到崩溃恢复情形,可能崩溃了很久. 如果这时候交易还继续进行,显然不合理.
	*/
	// Consider that crash happened for a long time, if transfer still goes on, that's not reasonable.
	if state.BlockNumber >= state.Transfer.Expiration {
		return &transfer.TransitionResult{
			NewState: state,
			Events:   nil,
		}
	}
	if st.Sender == state.Route.HopNode() && st.Secret == state.Transfer.Secret {
		/*
					   next hop learned the secret, unlock the token locally and send the
			         unlock message to next hop
		*/
		return &transfer.TransitionResult{
			NewState: nil,
			Events:   transferSuccessEvents(state),
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
			state = &mt.InitiatorState{
				OurAddress:                     staii.OurAddress,
				Transfer:                       staii.Tranfer,
				Routes:                         staii.Routes,
				BlockNumber:                    staii.BlockNumber,
				LockSecretHash:                 staii.LockSecretHash,
				Secret:                         staii.Secret,
				Db:                             staii.Db,
				CancelByExceptionSecretRequest: false,
			}
			return tryNewRoute(state)
		}
		/*
			作为交易发起方,发送完 Unlock 消息,对方确认收到,就应该认为这次交易彻底完成了
		*/
		// As transfer initiator, we assume that this transfer completes once we send unlock and my partner receive it.
		log.Warn(fmt.Sprintf("originalState,statechange should not be here originalState=\n%s\n,statechange=\n%s",
			utils.StringInterface1(originalState), utils.StringInterface1(st)))
	} else {
		switch st2 := st.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
			//只要密码正确,就应该发送secret ,流程上可能有问题,但是结果是没错的(只有在token swap的时候才会走到这一步) . 因为按照协议层要求,同一个消息不会重复发送, 导致在tokenswap的时候maker不可能重复发送reveal secret
			/*
					关于 token swap
					由于 同样的 reveal secret 双方都发送和接收了两遍,会有冗余的情况发生.
				 maker:
				1. maker发送给对方 reveal secret 的时候,同一个 lock 对应的两个 statemanager 都要知道密码,
				因为有可能对方是恶意的,一个恶意的实现就是, maker 发出的secret request 对方根本不响应,造成自己有一个 state manager 不知道密码,
				从而造成损失.
			*/
			/*
			 *	As long as secret correct, then we should send secret. There might be problematic about this procedure but result is correct.
			 *	Because according to protocol layer, same message won't send repeatedly, which leads to maker can't send reveal secret in tokenswap.
			 *
			 *	As to token swap, maybe redundency occurs because both participants send / receive revealsecret twice.
			 *
			 *	maker :
			 *		1. when maker sends reveal secret to his partner, two statemanager of a lock should know the secret.
			 *			Because maybe partner is fraudulent node, and he never responds to secret request, which leads to one stateManager without secret.
			 *
			 */
		case *mt.ReceiveSecretRevealStateChange:
			it = handleSecretReveal(state, st2)
		case *mt.ContractSecretRevealOnChainStateChange:
			it = handleSecretRevealOnChain(state, st2)
		case *mt.ReceiveSecretRequestStateChange:
			if state.RevealSecret == nil {
				it = handleSecretRequest(state, st2)
			} else {
				log.Warn(fmt.Sprintf("recevie secret request but initiator have already sent reveal secret"))
			}
		case *mt.ReceiveAnnounceDisposedStateChange:
			if state.RevealSecret == nil {
				it = handleRefund(state, st2)
			} else {
				log.Warn(fmt.Sprintf("secret already revealed ,but initiator recevied announce disposed %s", utils.StringInterface(st, 3)))
			}
		case *mt.ActionCancelRouteStateChange:
			if state.RevealSecret == nil {
				it = handleCancelRoute(state, st2)
			} else {
				panic(fmt.Sprintf("secret already revealed,route cannot canceled"))
			}
		case *transfer.ActionCancelTransferStateChange:
			if state.RevealSecret == nil {
				it = handleCancelTransfer(state)
			} else {
				panic(fmt.Sprintf("secret already revealed,transfer cannot canceled"))
			}
		case *mt.ContractCooperativeSettledStateChange:
			it = cancelCurrentRoute(state)
		case *mt.ContractChannelWithdrawStateChange:
			it = cancelCurrentRoute(state)
		default:
			log.Error(fmt.Sprintf("initiator received unkown state change %s", utils.StringInterface(st, 3)))
		}
	}
	return it
}
