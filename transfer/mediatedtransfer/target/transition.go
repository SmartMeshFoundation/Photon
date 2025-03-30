package target

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/Photon/utils"
)

//NameTargetTransition name for state manager
const NameTargetTransition = "TargetTransition"

func init() {
}

/*
Emits the event for closing the netting channel if from_transfer needs
    to be settled on-chain.
*/
func eventsForRegisterSecret(state *mediatedtransfer.TargetState) (events []transfer.Event) {
	fromTransfer := state.FromTransfer
	fromRoute := state.FromRoute
	safeToWait := mediator.IsSafeToWait(fromTransfer, fromRoute.RevealTimeout(), state.BlockNumber)
	secretKnown := fromTransfer.Secret != utils.EmptyHash
	payerClosed := state.FromRoute.State() == channeltype.StateClosed
	if safeToWait {
		safeToWait = !payerClosed //只要通道关闭,就应该立即注册密码,不要等过期了.
	}
	if !safeToWait && secretKnown {
		state.State = mediatedtransfer.StateWaitingRegisterSecret
		channelClose := &mediatedtransfer.EventContractSendRegisterSecret{
			Secret: fromTransfer.Secret,
		}
		events = append(events, channelClose)
	}
	return
}

//handleInitTarget Handle an ActionInitTarget state change.
func handleInitTarget(st *mediatedtransfer.ActionInitTargetStateChange) *transfer.TransitionResult {
	tr := st.FromTranfer
	route := st.FromRoute
	blockNumber := st.BlockNumber
	state := &mediatedtransfer.TargetState{
		OurAddress:               st.OurAddress,
		FromRoute:                route,
		FromTransfer:             tr,
		BlockNumber:              blockNumber,
		Db:                       st.Db,
		IsEffectiveChain:         st.IsEffectiveChain,
		EffectiveChangeTimestamp: st.EffectiveChangeTimestamp,
	}
	safeToWait := mediator.IsSafeToWait(tr, route.RevealTimeout(), blockNumber)
	/*
			  if there is not enough time to safely withdraw the token on-chain
		     silently let the transfer expire.
	*/
	if safeToWait {
		secretRequest := &mediatedtransfer.EventSendSecretRequest{
			ChannelIdentifier: route.ChannelIdentifier,
			LockSecretHash:    tr.LockSecretHash,
			Amount:            tr.Amount,
			Receiver:          tr.Initiator,
		}
		return &transfer.TransitionResult{
			NewState: state,
			Events:   []transfer.Event{secretRequest},
		}
	}
	//如果超时了,那就什么都不做,等待相关各方自己取消?
	// If timeout, then do nothing and wait to cancel this lock via participants themselves?
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

//handleSecretRegisteredOnChain this state manager has finished
func handleSecretRegisteredOnChain(state *mediatedtransfer.TargetState, st *mediatedtransfer.ContractSecretRevealOnChainStateChange) (it *transfer.TransitionResult) {
	var events []transfer.Event
	validSecret := st.LockSecretHash == state.FromTransfer.LockSecretHash
	if validSecret {
		/*
			无论是否超时,都应该结束了,
			没有超时,交易成功结束
			超时,交易失败结束
		*/
		/*
		 *	Not timeout, transfer finishes successfully.
		 *	timeout, transfer failed.
		 */
		state.State = mediatedtransfer.StateSecretRegistered
		ev := &mediatedtransfer.EventRemoveStateManager{
			Key: utils.Sha3(st.LockSecretHash[:], state.FromTransfer.Token[:]),
		}
		events = append(events, ev)
		state.Secret = st.Secret
		state.FromTransfer.Secret = st.Secret
		//链上注册没有过期,并且通道已经关闭,说明我还需要再次unlock
		if st.BlockNumber < state.FromTransfer.Expiration && state.FromRoute.State() == channeltype.StateClosed {
			events = append(events, &mediatedtransfer.EventContractSendUnlock{
				LockSecretHash:    st.LockSecretHash,
				ChannelIdentifier: state.FromRoute.ChannelIdentifier,
			})
		}
	} else {
		panic("should not here")
	}
	it = &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
	return
}

// Validate and handle a ReceiveSecretReveal state change.
func handleSecretReveal(state *mediatedtransfer.TargetState, st *mediatedtransfer.ReceiveSecretRevealStateChange) (it *transfer.TransitionResult) {
	validSecret := utils.ShaSecret(st.Secret[:]) == state.FromTransfer.LockSecretHash
	// 判断是否超时,如果已经该锁已经超时,不发送secret给上家
	isExpired := state.BlockNumber > state.FromTransfer.Expiration
	var events []transfer.Event
	if validSecret && !isExpired {
		tr := state.FromTransfer
		route := state.FromRoute
		state.State = mediatedtransfer.StateRevealSecret
		// 仅在第一次收到reveal secret消息的时候,保留data字段
		if tr.Secret == utils.EmptyHash {
			tr.Data = string(st.Message.Data)
		}
		tr.Secret = st.Secret
		reveal := &mediatedtransfer.EventSendRevealSecret{
			LockSecretHash: tr.LockSecretHash,
			Secret:         tr.Secret,
			Token:          tr.Token,
			Receiver:       route.HopNode(),
			Sender:         state.OurAddress,
		}
		events = append(events, reveal)
	} else {
		// TODO: event for byzantine behavior
	}
	it = &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
	return
}

/*
我收到了对方的 unlock 消息以后,就算是彻底结束了.
*/
/*
 *	handleBalanceProof : function to handle event of BalanceProof.
 *
 *	Note that once this participant receives unlock message from his channel partner, the function ends.
 */
func handleBalanceProof(state *mediatedtransfer.TargetState, st *mediatedtransfer.ReceiveUnlockStateChange) (it *transfer.TransitionResult) {
	var events []transfer.Event
	//TODO: byzantine behavior event when the sender doesn't match
	if st.NodeAddress == state.FromRoute.HopNode() && state.FromTransfer.LockSecretHash == st.LockSecretHash {
		state.State = mediatedtransfer.StateBalanceProof
		ev := &mediatedtransfer.EventRemoveStateManager{
			Key: utils.Sha3(state.FromTransfer.LockSecretHash[:], state.FromTransfer.Token[:]),
		}
		events = append(events, ev)
	}
	it = &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
	return
}

/*
After Photon learns about a new block this function must be called to
    handle expiration of the hash time lock.
*/
func handleBlock(state *mediatedtransfer.TargetState, st *transfer.BlockStateChange) (it *transfer.TransitionResult) {
	if state.BlockNumber < st.BlockNumber {
		state.BlockNumber = st.BlockNumber
	}
	/*
	   only emit the close event once

	*/
	var events []transfer.Event
	if state.State != mediatedtransfer.StateWaitingRegisterSecret && state.State != mediatedtransfer.StateSecretRegistered {
		events = eventsForRegisterSecret(state)
	}
	it = &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
	return
}

//Clear the state if the transfer was either completed or failed
func clearIfFinalized(previt *transfer.TransitionResult) (it *transfer.TransitionResult) {
	if previt.NewState == nil {
		return previt
	}
	state, ok := previt.NewState.(*mediatedtransfer.TargetState)
	if !ok {
		panic(fmt.Sprintf("clearIfFinalized for targetstate type error:%s", utils.StringInterface1(previt)))
	}
	it = previt
	if state.FromTransfer.Secret == utils.EmptyHash && state.BlockNumber > state.FromTransfer.Expiration {
		failed := &mediatedtransfer.EventWithdrawFailed{
			LockSecretHash:    state.FromTransfer.LockSecretHash,
			ChannelIdentifier: state.FromRoute.ChannelIdentifier,
			Reason:            "lock expired",
		}
		it = &transfer.TransitionResult{
			NewState: nil,
			Events:   append(it.Events, failed),
		}
	} else if state.State == mediatedtransfer.StateBalanceProof {
		//这些事件对应的处理都没有
		// these events have no related handle solution
		transferSuccess := &transfer.EventTransferReceivedSuccess{
			LockSecretHash:    state.FromTransfer.LockSecretHash,
			Amount:            state.FromTransfer.Amount,
			Initiator:         state.FromTransfer.Initiator,
			ChannelIdentifier: state.FromRoute.ChannelIdentifier,
			Data:              state.FromTransfer.Data,
		}
		unlockSuccess := &mediatedtransfer.EventWithdrawSuccess{
			LockSecretHash: state.FromTransfer.LockSecretHash,
		}
		it = &transfer.TransitionResult{
			NewState: nil,
			Events:   append(it.Events, transferSuccess, unlockSuccess),
		}
	}
	// 一旦锁过期,就结束了,注销StateManager
	// Once locks expired, remove StateManager.
	if state.BlockNumber > state.FromTransfer.Expiration {
		it.Events = append(it.Events, &mediatedtransfer.EventRemoveStateManager{
			Key: utils.Sha3(state.FromTransfer.LockSecretHash[:], state.FromTransfer.Token[:]),
		})
	}
	return it
}

// StateTransiton is State machine for the target node of a target transfer.
func StateTransiton(originalState transfer.State, stateChange transfer.StateChange) (it *transfer.TransitionResult) {
	it = &transfer.TransitionResult{
		NewState: originalState,
		Events:   nil,
	}
	if originalState == nil {
		ait, ok := stateChange.(*mediatedtransfer.ActionInitTargetStateChange)
		if ok {
			it = handleInitTarget(ait)
		}
	} else {
		state, ok := originalState.(*mediatedtransfer.TargetState)
		if !ok {
			panic(fmt.Sprintf("targetstate StateTransiton type error:%s", utils.StringInterface1(originalState)))
		}
		switch st2 := stateChange.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
		case *mediatedtransfer.ContractSecretRevealOnChainStateChange:
			it = handleSecretRegisteredOnChain(state, st2)
		case *mediatedtransfer.ReceiveSecretRevealStateChange:
			if state.FromTransfer.Secret == utils.EmptyHash {
				//可能会反复收到 reveal secret, 比如 token swap的时候,再比如存在环路的时候
				// Maybe we can receive reveal secret over and over again,
				// such as when using token swap, or circuit exist.
				it = handleSecretReveal(state, st2)
			}
		case *mediatedtransfer.ReceiveUnlockStateChange:
			//有可能在不知道密码的情况下直接收到 unlock 消息,比如
			// Maybe we can receive unlock message without receiving secret.
			it = handleBalanceProof(state, st2)
		case *transfer.EffectiveChainStateChange:
			state.IsEffectiveChain = st2.IsEffective
			state.EffectiveChangeTimestamp = st2.LastBlockNumberTimestamp
			log.Info(fmt.Sprintf("TargetStateManager with lockSecretHash=%s EffctiveChainState change to %v", state.FromTransfer.LockSecretHash.String(), state.IsEffectiveChain))
		default:
			log.Error(fmt.Sprintf("target state manager receive unkown state change,if this transfer is a token swap ,it's ok.  %s", utils.StringInterface(stateChange, 3)))
		}
	}
	return clearIfFinalized(it)
}
