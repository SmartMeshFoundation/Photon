package crashnode

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	mt "github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

/*
只处理两件事情
1. 我发出的锁,需要过期以后发送 RemoveExpiredLock
2. 收到的对方的锁,如果我知道密码,那么应该在接近到期的时候去注册密码
3. 收到 RevealSecretOnChain 要判断是否需要发送 Unlock 消息
4. 收到 unlock 消息,移除响应的锁.
*/
// This piece of code only process two things
// 1. all locks sent by me should send `RemoveExpiredLock` after expired.
// 2. Every transfer I receive from my partner, if secret is known, I should register this secret before expiration.
// 3. Once received RevealSecretOnChain, we need to check whether we should send Unlock.
// 4. Once received unlock message, we should remove relevant lock.

//NameCrashNodeTransition name for state manager
const NameCrashNodeTransition = "CrashTransition"

/*
make sure not call this when transfer already finished , state is nil means finished.
请确保重启以后,先处理完链上待处理的事件,比如 SettleChannel,SecretRegisterOnChain 等,然后再处理新块事件,
否则可能会产生冲突,导致错误.
如何处理我发送的锁:
1. 如果过期,发送 RemovedExpiredLock
 我们不能先发送 RemovedExpiredLock, 然后再发送 Unlock 消息. 重新启动以后,密码在链上注册了,但是我没处理,就发送了 RemovedExpiredLock,这是不合理的.
如何处理我收到的锁:
1. 如果过期,直接忽略,标记为已经处理
2. 如果我知道密码,并且快要过期,那么应该去链上注册密码
*/
/*
 *	handleBlock : function to handle crash case.
 *
 *	Note that we should make sure not call this when transfer already finished, state is nil means finished.
 *	And we should make sure we need to first handle pending events, such as SettleChannel, SecretRegisterOnChain, etc, then handle events of new block events.
 *	Or there may be conflicts and error occurs.
 *	How to handle locks that I have sent out :
 *		1. If locks expired, send RemoveDExpiredLock
 * 			- It is unreasonable case that we first send RemovedExpiredLock, then send Unlock, once node restarts, and secret has been registered on chain,
 *			but I send RemovedExpiredLock without handling those messages.
 *	How to handle locks that I have received :
 *		1. If locks expired, just miss them, tag them with message of handled.
 *		2. If I know the secret and it is close to expiration, then I should register the secret on chain.
 */
func handleBlock(state *mt.CrashState, stateChange *transfer.BlockStateChange) *transfer.TransitionResult {
	var events []transfer.Event
	var removedSentIndex []int
	for i, l := range state.SentLocks {
		if stateChange.BlockNumber > l.Lock.Expiration {
			removedSentIndex = append(removedSentIndex, i)
			events = append(events, &mt.EventUnlockFailed{
				LockSecretHash:    l.Lock.LockSecretHash,
				ChannelIdentifier: l.Channel.ChannelIdentifier.ChannelIdentifier,
				Reason:            "lock expired",
			})
		}
	}
	var removedReceviedIndex []int
	for i, l := range state.ReceivedLocks {
		if !l.Channel.PartnerState.IsKnown(l.Lock.LockSecretHash) {
			//有可能有些锁因为收到对方的 AnnounceDisposedResponse, 对方的 unlock 消息直接解锁了,移除即可.
			removedReceviedIndex = append(removedReceviedIndex, i)
			continue
		}
		if stateChange.BlockNumber > l.Lock.Expiration {
			removedReceviedIndex = append(removedReceviedIndex, i)
		} else {
			/*
				有可能在我启动以后才收到相应的密码
			*/
			secret, found := l.Channel.PartnerState.GetSecret(state.LockSecretHash)
			if found && l.Lock.Expiration > stateChange.BlockNumber-int64(l.Channel.RevealTimeout) && l.Lock.Expiration < stateChange.BlockNumber {
				//临近过期了,需要通知链上注册
				events = append(events, &mt.EventContractSendRegisterSecret{
					Secret: secret,
				})
				removedReceviedIndex = append(removedReceviedIndex, i)
			}
		}
	}
	if len(removedReceviedIndex) > 0 || len(removedSentIndex) > 0 {
		for _, i := range removedSentIndex {
			state.ProcessedSentLocks = append(state.ProcessedSentLocks, state.SentLocks[i])
		}
		for _, i := range removedReceviedIndex {
			state.ProcessedReceivedLocks = append(state.ProcessedReceivedLocks, state.ReceivedLocks[i])
		}
		state.SentLocks = removeSliceFromSlice(state.SentLocks, removedSentIndex)
		state.ReceivedLocks = removeSliceFromSlice(state.ReceivedLocks, removedReceviedIndex)
		events = append(events, checkFinish(state)...)
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

/*
从一个 slice 中移除一组 slice, 带移除的元素通过下标指定.
*/
/*
 *	removeSliceFromSlice : remove one slice denoted by among a set of slices,
 */
func removeSliceFromSlice(src []*mt.LockAndChannel, removedIndex []int) (dest []*mt.LockAndChannel) {
	for i, l := range src {
		found := false
		for _, j := range removedIndex {
			if i == j {
				found = true
				break
			}
		}
		if !found {
			dest = append(dest, l)
		}
	}
	return
}

/*
密码在链上注册了,只要在有效期范围内,就相当于收到了对方的 reveal secret, 主动给对方发送 unlock 消息.
*/
/*
 *	handleSecretRevealOnChain : function to handle events of reveal secret on chain.
 *
 *	Note that once secret has been registered on chain, and before expiration, we assume that all other nodes in this payment channel
 *	have received this secret and they will send unlock to their partner.
 */
func handleSecretRevealOnChain(state *mt.CrashState, st *mt.ContractSecretRevealOnChainStateChange) (it *transfer.TransitionResult) {
	it = &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
	if st.LockSecretHash != state.LockSecretHash {
		//无论是不是 token swap, 都应该知道 locksecrethash,否则肯定是实现有问题
		panic(fmt.Sprintf("my locksecrethash=%s,received=%s", state.LockSecretHash.String(), st.LockSecretHash.String()))
	}
	var removedIndex []int
	for i, l := range state.SentLocks {
		if l.Lock.Expiration >= st.BlockNumber {
			if l.Lock.LockSecretHash != st.LockSecretHash {
				panic(fmt.Sprintf("receive SecretRevealOnChain for different lockseret mine=%s,statechange=%s", l.Lock.LockSecretHash.String(), st.LockSecretHash.String()))
			}
			//send unlock event
			it.Events = append(it.Events, transferSuccessEvents(l)...)
			removedIndex = append(removedIndex, i)
		}
	}
	if len(removedIndex) > 0 {
		for _, i := range removedIndex {
			state.ProcessedSentLocks = append(state.ProcessedSentLocks, state.SentLocks[i])
		}
		state.SentLocks = removeSliceFromSlice(state.SentLocks, removedIndex)
		it.Events = append(it.Events, checkFinish(state)...)
	}

	/*
			至于我收到的锁,密码注册了,我也不用管,
		1. 如果对方发送给我 unlock, 那么我移除就可以了,
		2. 如果不给我发 unlock, 我到时候主动去链上注册密码(等于什么都不做)
		2. 如果过期了,我移除就可以了
			todo 如果做到这点,我们就必须保证,如果重启以后,所有的事件处理完了,然后再去处理新块事件,如何做到呢?
	*/
	return
}

/*
 todo 是否会发生 RevealSecretOnChain 针对的是旧的 channel 的 Lock?
比如我关掉了很长一段时间,然后重新启动,这时候收到了 RevealSecretOnChain, 但是我持有的锁完全是过去的某个锁呢?
*/
func transferSuccessEvents(l *mt.LockAndChannel) (events []transfer.Event) {
	unlockLock := &mt.EventSendBalanceProof{
		LockSecretHash:    l.Lock.LockSecretHash,
		ChannelIdentifier: l.Channel.ChannelIdentifier.ChannelIdentifier,
		Token:             l.Channel.TokenAddress,
		Receiver:          l.Channel.PartnerState.Address,
	}

	unlockSuccess := &mt.EventUnlockSuccess{
		LockSecretHash: l.Lock.LockSecretHash,
	}
	events = []transfer.Event{unlockLock, unlockSuccess}
	return events
}

/*
这个 stateManager 是否彻底完成了
如果发送和收到的锁都处理完毕了,
发送的锁都移动到 ProcessedSentLock,因为过期或者发送了 Unlock消息
收到的锁都移动到了 ProcessedReceviedLock, 因为链上注册了密码.
*/
/*
 *	checkFinish : function to check whether statemanager has been finished.
 *
 *	Note that if locks that are sent or received have been handled,
 * 	then locks that this participant sent moves to ProcessedSentLock, because they are expired or unlocked.
 *	And locks that this participant received moves to ProcessedReceivedLock, because secret has been registered on chain.
 */
func checkFinish(state *mt.CrashState) (events []transfer.Event) {
	if len(state.SentLocks) == 0 && len(state.ReceivedLocks) == 0 {
		events = append(events, &mt.EventRemoveStateManager{
			Key: utils.Sha3(state.LockSecretHash[:], state.Token[:]),
		})
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
func handleBalanceProof(state *mt.CrashState, st *mt.ReceiveUnlockStateChange) (it *transfer.TransitionResult) {
	it = &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
	for i, l := range state.ReceivedLocks {
		if st.NodeAddress == l.Channel.PartnerState.Address && l.Lock.LockSecretHash == st.LockSecretHash {
			state.ReceivedLocks = append(state.ReceivedLocks[:i], state.ReceivedLocks[i+1:]...)
			state.ProcessedReceivedLocks = append(state.ProcessedReceivedLocks, l)
		}
	}
	return
}

/*
我可能在重启以后收到来自对方的 AnnounceDisposed 消息,如果一切正确,我应该发送AnnounceDisposedResponse
*/
/*
 *	handleAnnounceDisposed : function to handle event of AnnounceDisposed
 *
 *	Note that this participant is probable to receive AnnounceDisposed from his partner after he reconnects,
 *	if everything is correct, then he should send AnnounceDisposedResponse
 */
func handleAnnounceDisposed(state *mt.CrashState, st *mt.ReceiveAnnounceDisposedStateChange) (it *transfer.TransitionResult) {
	for i, l := range state.SentLocks {
		if l.Lock.Equal(st.Lock) && st.Token == l.Channel.TokenAddress && st.Sender == l.Channel.PartnerState.Address {
			state.SentLocks = append(state.SentLocks[:i], state.SentLocks[i+1:]...)
			state.ProcessedSentLocks = append(state.ProcessedSentLocks, l)
			return &transfer.TransitionResult{
				NewState: state,
				Events: []transfer.Event{&mt.EventSendAnnounceDisposedResponse{
					LockSecretHash: l.Lock.LockSecretHash,
					Token:          state.Token,
					Receiver:       st.Sender,
				}},
			}
		}
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

/*
StateTransition is State machine for a node who restart after a crash
    originalState: The current State that is transitioned from.
    st: The state_change that will be applied.
*/
func StateTransition(originalState transfer.State, st transfer.StateChange) *transfer.TransitionResult {
	it := &transfer.TransitionResult{
		NewState: originalState,
		Events:   nil,
	}
	state, ok := originalState.(*mt.CrashState)
	if !ok {
		if originalState != nil {
			panic("crash StateTransition get type error")
		}
		state = nil //originalState is nil
	}
	if state == nil {
		staii, ok := st.(*mt.ActionInitCrashRestartStateChange)
		if ok {
			state = &mt.CrashState{
				OurAddress:     staii.OurAddress,
				LockSecretHash: staii.LockSecretHash,
				SentLocks:      staii.SentLocks,
				ReceivedLocks:  staii.ReceivedLocks,
				Token:          staii.Token,
			}
			it.NewState = state
			return it
		}
		/*
			作为交易发起方,发送完 Unlock 消息,对方确认收到,就应该认为这次交易彻底完成了
		*/
		//todo fix, find a way to remove this identifier from raiden.Transfer2StateManager
		log.Error(fmt.Sprintf("originalState,statechange should not be here originalState=\n%s\n,statechange=\n%s",
			utils.StringInterface1(originalState), utils.StringInterface1(st)))
	} else {
		switch st2 := st.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
		case *mt.ContractSecretRevealOnChainStateChange:
			it = handleSecretRevealOnChain(state, st2)
		case *mt.ReceiveUnlockStateChange:
			it = handleBalanceProof(state, st2)
		case *mt.ReceiveAnnounceDisposedStateChange:
			//有可能重启以后收到对方的 AnnounceDisposed 消息,我也需要正常处理,移除锁.
			it = handleAnnounceDisposed(state, st2)
		default:
			//我重启了,应该会收到来自对方的消息,但是我都不会处理.
			log.Warn(fmt.Sprintf("crash state manager received unkown state change %s", utils.StringInterface(st, 3)))
		}
	}
	return it
}
