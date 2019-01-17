package mediator

import (
	"fmt"

	"math/big"

	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/route"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//NameMediatorTransition name for state manager
const NameMediatorTransition = "MediatorTransition"

/*
 Reduce the lock expiration by some additional blocks to prevent this exploit:
 The payee could reveal the secret on it's lock expiration block, the lock
 would be valid and the previous lock can be safely unlocked so the mediator
 would follow the secret reveal with a balance-proof, at this point the secret
 is known, the payee transfer is payed, and if the payer expiration is exactly
 reveal_timeout blocks away the mediator will be forced to close the channel
 to be safe.
*/
var stateSecretKnownMaps = map[string]bool{
	mediatedtransfer.StatePayeeSecretRevealed: true,
	mediatedtransfer.StatePayeeBalanceProof:   true,

	mediatedtransfer.StatePayerSecretRevealed:        true,
	mediatedtransfer.StatePayerWaitingRegisterSecret: true,
	mediatedtransfer.StatePayerBalanceProof:          true,
}
var stateTransferPaidMaps = map[string]bool{
	mediatedtransfer.StatePayeeBalanceProof: true,
	mediatedtransfer.StatePayerBalanceProof: true,
}

var stateTransferFinalMaps = map[string]bool{
	mediatedtransfer.StatePayeeExpired:      true,
	mediatedtransfer.StatePayeeBalanceProof: true,

	mediatedtransfer.StatePayerExpired:      true,
	mediatedtransfer.StatePayerBalanceProof: true,
}

//True if the lock has not expired.
func isLockValid(tr *mediatedtransfer.LockedTransferState, blockNumber int64) bool {
	return blockNumber <= tr.Expiration
}

/*
IsSafeToWait returns True if there are more than enough blocks to safely settle on chain and
    waiting is safe.
*/
func IsSafeToWait(tr *mediatedtransfer.LockedTransferState, revealTimeout int, blockNumber int64) bool {
	// A node may wait for a new balance proof while there are reveal_timeout
	// left, at that block and onwards it is not safe to wait.
	return blockNumber < tr.Expiration-int64(revealTimeout)
}

//IsValidRefund returns True if the refund transfer matches the original transfer.
func IsValidRefund(originTr *mediatedtransfer.LockedTransferState, originRoute *route.State, st *mediatedtransfer.ReceiveAnnounceDisposedStateChange) bool {
	//Ignore a refund from the target
	if st.Sender == originTr.Target {
		return false
	}
	if st.Sender != originRoute.HopNode() {
		return false
	}
	return originTr.Amount.Cmp(st.Lock.Amount) == 0 &&
		originTr.LockSecretHash == st.Lock.LockSecretHash &&
		originTr.Token == st.Token &&
		/*
			必须严格相等,我们要的是锁的 hash
		*/
		originTr.Expiration == st.Lock.Expiration
}

/*
True if this node needs to register secret on chain

    Only close the channel to withdraw on chain if the corresponding payee node
    has received, this prevents attacks were the payee node burns it's payment
    to force a close with the payer channel.
*/
func isSecretRegisterNeeded(tr *mediatedtransfer.MediationPairState, blockNumber int64) bool {
	payeeReceived := stateTransferPaidMaps[tr.PayeeState]
	payerPayed := stateTransferPaidMaps[tr.PayerState]
	channelClosed := tr.PayerRoute.State() == channeltype.StateClosed
	AlreadyRegisterring := tr.PayerState == mediatedtransfer.StatePayerWaitingRegisterSecret
	safeToWait := IsSafeToWait(tr.PayerTransfer, tr.PayerRoute.RevealTimeout(), blockNumber)
	//如果payer已经支付了不用注册密码
	//如果已经在链上注册密码了,不用注册
	//如果快要过期了,链上注册
	//如果通道已经关闭了,尽快注册.
	return ((payeeReceived && !safeToWait) || channelClosed) && !AlreadyRegisterring && !payerPayed
}

//Return the transfer pairs that are not at a final state.
func getPendingTransferPairs(pairs []*mediatedtransfer.MediationPairState) (pendingPairs []*mediatedtransfer.MediationPairState) {
	for _, pair := range pairs {
		if !stateTransferFinalMaps[pair.PayeeState] || !stateTransferFinalMaps[pair.PayerState] {
			pendingPairs = append(pendingPairs, pair)
		}
	}
	return pendingPairs
}

func getExpiredTransferPairs(pairs []*mediatedtransfer.MediationPairState) (pendingPairs []*mediatedtransfer.MediationPairState) {
	for _, pair := range pairs {
		if pair.PayeeState == mediatedtransfer.StatePayeeExpired {
			pendingPairs = append(pendingPairs, pair)
		}
	}
	return pendingPairs
}

/*
Return the timeout blocks, it's the base value from which the payee's
    lock timeout must be computed.

    The payee lock timeout is crucial for safety of the mediated transfer, the
    value must be chosen so that the payee hop is forced to reveal the secret
    with sufficient time for this node to claim the received lock from the
    payer hop.

    The timeout blocks must be the smallest of:

    - payerTransfer.expiration: The payer lock expiration, to force the payee
      to reveal the secret before the lock expires.
    - payerRoute.settleTimeout: Lock expiration must be lower than
      the settlement period since the lock cannot be claimed after the channel is
      settled.
    - payerRoute.ClosedBlock: If the channel is closed then the settlement
      period is running and the lock expiration must be lower than number of
      blocks left.
*/
func getTimeoutBlocks(payerRoute *route.State, payerTransfer *mediatedtransfer.LockedTransferState, blockNumber int64) int64 {
	blocksUntilSettlement := int64(payerRoute.SettleTimeout())
	if payerRoute.ClosedBlock() != 0 {
		if blockNumber < payerRoute.ClosedBlock() {
			panic("ClosedBlock bigger than the lastest blocknumber")
		}
		blocksUntilSettlement -= blockNumber - payerRoute.ClosedBlock()
	}
	if blocksUntilSettlement > payerTransfer.Expiration-blockNumber {
		blocksUntilSettlement = payerTransfer.Expiration - blockNumber
	}
	log.Debug(fmt.Sprintf("get transfer lockSecretHash=%s, expiration=%d, now=%d, blocksUntilSettlement=%d",
		utils.HPex(payerTransfer.LockSecretHash), payerTransfer.Expiration, blockNumber, blocksUntilSettlement))
	return blocksUntilSettlement
}

//Check invariants that must hold.
//return error is better for production environment
func sanityCheck(state *mediatedtransfer.MediatorState) {
	if len(state.TransfersPair) == 0 {
		return
	}
	//if a transfer is paid we must know the secret
	for _, pair := range state.TransfersPair {
		if stateTransferPaidMaps[pair.PayerState] && state.Secret == utils.EmptyHash {
			panic(fmt.Sprintf("payer:a transfer is paid but we don't know the secret, payerstate=%s,payeestate=%s", pair.PayerState, pair.PayeeState))
		}
		if stateTransferPaidMaps[pair.PayeeState] && state.Secret == utils.EmptyHash {
			panic(fmt.Sprintf("payee:a transfer is paid but we don't know the secret,payerstate=%s,payeestate=%s", pair.PayerState, pair.PayeeState))
		}
	}
	//the "transitivity" for these values is checked below as part of
	//almost_equal check
	if len(state.TransfersPair) > 0 {
		firstPair := state.TransfersPair[0]
		if state.Hashlock != firstPair.PayerTransfer.LockSecretHash {
			panic("sanity check failed:state.LockSecretHash!=firstPair.PayerTransfer.LockSecretHash")
		}
		if state.Secret != utils.EmptyHash {
			if firstPair.PayerTransfer.Secret != state.Secret {
				panic("sanity check failed:firstPair.PayerTransfer.Secret!=state.Secret")
			}
		}
	}
	for _, p := range state.TransfersPair {
		if !p.PayerTransfer.AlmostEqual(p.PayeeTransfer) {
			panic("sanity check failed:PayerTransfer.AlmostEqual(p.PayeeTransfer)")
		}
		if p.PayerTransfer.Expiration < p.PayeeTransfer.Expiration {
			panic("sanity check failed:PayerTransfer.Expiration<=p.PayeeTransfer.Expiration")
		}
		if !mediatedtransfer.ValidPayerStateMap[p.PayerState] {
			panic(fmt.Sprint("sanity check failed: payerstate invalid :", p.PayerState))
		}
		if !mediatedtransfer.ValidPayeeStateMap[p.PayeeState] {
			panic(fmt.Sprint("sanity check failed: payee invalid :", p.PayeeState))
		}
	}
	pairs2 := state.TransfersPair[0 : len(state.TransfersPair)-1]
	for i := range pairs2 {
		original := state.TransfersPair[i]
		refund := state.TransfersPair[i+1]
		if !original.PayeeTransfer.AlmostEqual(refund.PayerTransfer) {
			panic("sanity check failed:original.PayeeTransfer.AlmostEqual(refund.PayerTransfer)")
		}
		if original.PayeeRoute.HopNode() == refund.PayerRoute.HopNode() {
			panic("sanity check failed:original.PayeeRoute.HopNode==refund.PayerRoute.HopNode")
		}
		if original.PayeeTransfer.Expiration < refund.PayerTransfer.Expiration {
			panic("sanity check failed:original.PayeeTransfer.Expiration>refund.PayerTransfer.Expiration")
		}
	}
}

//Clear the state if all transfer pairs have finalized
func clearIfFinalized(result *transfer.TransitionResult) *transfer.TransitionResult {
	if result.NewState == nil {
		return result
	}
	state := result.NewState.(*mediatedtransfer.MediatorState)
	isAllFinalized := true
	for _, p := range state.TransfersPair {
		if !stateTransferPaidMaps[p.PayeeState] || !stateTransferPaidMaps[p.PayerState] {
			isAllFinalized = false
			break
		}
	}
	if isAllFinalized {
		// 此处可以移除state manager了
		return &transfer.TransitionResult{
			NewState: nil,
			Events: append(result.Events, &mediatedtransfer.EventRemoveStateManager{
				Key: utils.Sha3(state.LockSecretHash[:], state.Token[:]),
			}),
		}
	}
	return result
}

/*
Finds the first route available that may be used.
        rss  : Current available routes that may be used,
            it's assumed that the available_routes list is ordered from best to
            worst.
        timeoutBlocks : Base number of available blocks used to compute
            the lock timeout.
        transferAmount : The amount of tokens that will be transferred
            through the given route.
    Returns:
         The next route.
*/
/*
找到一个满足下面提交的 route
1.通道金额足够
2.上家给出的 费用足够
3.时间还足够安全
*/

func nextRoute(fromRoute *route.State, rss *route.RoutesState, timeoutBlocks int, transferAmount, fee *big.Int) *route.State {
	for len(rss.AvailableRoutes) > 0 {
		route := rss.AvailableRoutes[0]
		rss.AvailableRoutes = rss.AvailableRoutes[1:]
		lockTimeout := timeoutBlocks - route.RevealTimeout()
		/*
				1.通道金额足够
				2. 给出的收费也够
				3. 时间还安全
				4. 通道可以发起交易
				5. 不能使用再次使用上家做下一跳.
			 有可能形成环路的时候,上家已经在我认为可用的路由节点中,但是实际上就是从他发过来的 lockedTransfer
		*/
		if route.CanTransfer() && route.AvailableBalance().Cmp(transferAmount) >= 0 && lockTimeout > 0 && fee.Cmp(route.Fee) >= 0 && route.HopNode() != fromRoute.HopNode() {
			return route
		}
		rss.IgnoredRoutes = append(rss.IgnoredRoutes, route)
	}
	return nil
}

/*
Given a payer transfer tries a new route to proceed with the mediation.

        payerRoute  : The previous route in the path that provides
            the token for the mediation.
        payerTransfer : The transfer received from the
            payerRoute.
        routesState  : Current available routes that may be used,
            it's assumed that the available_routes list is ordered from best to
            worst.
        timeoutBlocks : Base number of available blocks used to compute
            the lock timeout.
        blockNumber  : The current block number.
*/
func nextTransferPair(payerRoute *route.State, payerTransfer *mediatedtransfer.LockedTransferState,
	routesState *route.RoutesState, timeoutBlocks int, blockNumber int64) (
	transferPair *mediatedtransfer.MediationPairState, events []transfer.Event) {
	if timeoutBlocks <= 0 {
		panic("timeoutBlocks<=0")
	}
	if int64(timeoutBlocks) > payerTransfer.Expiration-blockNumber {
		panic("timeoutBlocks >payerTransfer.Expiration-blockNumber")
	}
	payeeRoute := nextRoute(payerRoute, routesState, timeoutBlocks, payerTransfer.Amount, payerTransfer.Fee)
	if payeeRoute != nil {
		/*
					有可能 payeeroute 的 settle timeout 比较小,从而导致我指定的lockexpiration 特别大,从而对我不利.
				这个地方需要一个例子,例子是什么呢?
			如果 timeoutBlocks 超过 settle timout 会有什么问题呢

		*/
		if timeoutBlocks >= payeeRoute.SettleTimeout() {
			timeoutBlocks = payeeRoute.SettleTimeout()
		}
		//不再减少时间,没有必要了,只要这个时间不超过 payee 的 settle timeout 即可
		lockTimeout := timeoutBlocks //- payeeRoute.RevealTimeout()
		lockExpiration := int64(lockTimeout) + blockNumber
		payeeTransfer := &mediatedtransfer.LockedTransferState{
			TargetAmount:   payerTransfer.TargetAmount,
			Amount:         big.NewInt(0).Sub(payerTransfer.Amount, payeeRoute.Fee),
			Token:          payerTransfer.Token,
			Initiator:      payerTransfer.Initiator,
			Target:         payerTransfer.Target,
			Expiration:     lockExpiration,
			LockSecretHash: payerTransfer.LockSecretHash,
			Secret:         payerTransfer.Secret,
			Fee:            big.NewInt(0).Sub(payerTransfer.Fee, payeeRoute.Fee),
		}
		if payeeRoute.HopNode() == payeeTransfer.Target {
			//i'm the last hop,so take the rest of the fee
			payeeTransfer.Fee = utils.BigInt0
			payeeTransfer.Amount = payerTransfer.TargetAmount
		}
		//todo log how many tokens fee for this transfer .
		transferPair = mediatedtransfer.NewMediationPairState(payerRoute, payeeRoute, payerTransfer, payeeTransfer)
		eventSendMediatedTransfer := mediatedtransfer.NewEventSendMediatedTransfer(payeeTransfer, payeeRoute.HopNode())
		eventSendMediatedTransfer.FromChannel = payerRoute.ChannelIdentifier
		events = []transfer.Event{eventSendMediatedTransfer}
	}
	return
}

/*
Set the state of a transfer *sent* to a payee and check the secret is
    being revealed backwards.

        The elements from transfers_pair are changed in place, the list must
        contain all the known transfers to properly check reveal order.
*/
func setPayeeStateAndCheckRevealOrder(transferPair []*mediatedtransfer.MediationPairState, payeeAddress common.Address,
	newPayeeState string) []transfer.Event {
	if !mediatedtransfer.ValidPayeeStateMap[newPayeeState] {
		panic(fmt.Sprintf("invalid payeestate:%s", newPayeeState))
	}
	WrongRevealOrder := false
	for j := len(transferPair) - 1; j >= 0; j-- {
		back := transferPair[j]
		if back.PayeeRoute.HopNode() == payeeAddress {
			back.PayeeState = newPayeeState
			break
		} else if !stateSecretKnownMaps[back.PayeeState] {
			WrongRevealOrder = true
		}
	}
	if WrongRevealOrder {
		/*
					   TODO: Append an event for byzantine behavior.
			         XXX: With the current events_for_withdraw implementation this may
			         happen, should the notification about byzantine behavior removed or
			         fix the events_for_withdraw function fixed?
		*/
		return nil
	}
	return nil
}

/*
Set the state of expired transfers, and return the failed events
按照现在的规则
payer的 expiration>=payee 的 expiration 可以相等
*/
/*
 *	set the state of expired transfers, and return the failed events.
 *	According to current rule, that expiration between payer and payee can be equal.
 */
func setExpiredPairs(transfersPairs []*mediatedtransfer.MediationPairState, blockNumber int64) (events []transfer.Event, allExpired bool) {
	pendingTransfersPairs := getPendingTransferPairs(transfersPairs)
	allExpired = len(pendingTransfersPairs) == 0
	for _, pair := range pendingTransfersPairs {
		if blockNumber > pair.PayerTransfer.Expiration {
			if pair.PayeeState != mediatedtransfer.StatePayeeExpired {
				//未必一定, 两者的 expiration 很可能是相同的
				// we cannot be certain, both expiration can be identical.
			}
			if pair.PayeeTransfer.Expiration > pair.PayerTransfer.Expiration {
				panic("PayeeTransfer.Expiration>=pair.PayerTransfer.Expiration")
			}
			if pair.PayerState != mediatedtransfer.StatePayerExpired {
				pair.PayerState = mediatedtransfer.StatePayerExpired
				withdrawFailed := &mediatedtransfer.EventWithdrawFailed{
					LockSecretHash:    pair.PayerTransfer.LockSecretHash,
					ChannelIdentifier: pair.PayerRoute.ChannelIdentifier,
					Reason:            "lock expired",
				}
				events = append(events, withdrawFailed)
			}
		}
		/*
			考虑到分叉攻击,延迟一定块数之后才发送remove
			这里需要考虑如下问题:
			1. 本该过期之后,收到payee的reveal secret,这时候肯定不能给对方发送BalanceProof
			2. 本该过期之后,收到payee的AnnouceDisposed, 这时候谨慎起见,忽略这种消息
			3. 本该过期之后,收到链上的密码注册事件,因为链上事件给出的块数已经过期了,所以会被忽略
		*/
		if blockNumber-params.ForkConfirmNumber > pair.PayeeTransfer.Expiration {
			/*
			   For safety, the correct behavior is:

			   - If the payee has been paid, then the payer must pay too.

			     And the corollary:

			   - If the payer transfer has expired, then the payee transfer must
			     have expired too.

			   The problem is that this corollary cannot be asserted. If a user
			   is running Photon without a monitoring service, then it may go
			   offline after having paid a transfer to a payee, but without
			   getting a balance proof of the payer, and once it comes back
			   online the transfer may have expired.
			*/
			if pair.PayeeTransfer.Expiration > pair.PayerTransfer.Expiration {
				panic("PayeeTransfer.Expiration>=pair.PayerTransfer.Expiration")
			}
			if pair.PayeeState != mediatedtransfer.StatePayeeExpired {
				pair.PayeeState = mediatedtransfer.StatePayeeExpired
				unlockFailed := &mediatedtransfer.EventUnlockFailed{
					LockSecretHash:    pair.PayeeTransfer.LockSecretHash,
					ChannelIdentifier: pair.PayeeRoute.ChannelIdentifier,
					Reason:            "lock expired",
				}
				events = append(events, unlockFailed)
			}
		}
	}
	//// 考虑到分叉攻击,延迟一定块数之后才发送remove
	//expiredPairs := getExpiredTransferPairs(transfersPairs)
	//for _, pair := range expiredPairs {
	//	if blockNumber-params.ForkConfirmNumber > pair.PayeeTransfer.Expiration {
	//		unlockFailed := &mediatedtransfer.EventUnlockFailed{
	//			LockSecretHash:    pair.PayeeTransfer.LockSecretHash,
	//			ChannelIdentifier: pair.PayeeRoute.ChannelIdentifier,
	//			Reason:            "lock expired",
	//		}
	//		events = append(events, unlockFailed)
	//	}
	//}
	return
}

/*
Refund the transfer.

        refundRoute   The original route that sent the mediated
            transfer to this node.
        refundTransfer (LockedTransferState): The original mediated transfer
            from the refundRoute.
    Returns:
        create a annouceDisposed event
*/
func eventsForRefund(refundRoute *route.State, refundTransfer *mediatedtransfer.LockedTransferState) (events []transfer.Event) {
	/*
		原封不动声明放弃此锁即可
	*/
	// abandon this lock just fine.
	rtr2 := &mediatedtransfer.EventSendAnnounceDisposed{
		Token:          refundTransfer.Token,
		Amount:         new(big.Int).Set(refundTransfer.Amount),
		LockSecretHash: refundTransfer.LockSecretHash,
		Expiration:     refundTransfer.Expiration,
		Receiver:       refundRoute.HopNode(),
	}
	events = append(events, rtr2)
	return
}

/*
Reveal the secret backwards.

    This node is named N, suppose there is a mediated transfer with two refund
    transfers, one from B and one from C:

        A-N-B...B-N-C..C-N-D

    Under normal operation N will first learn the secret from D, then reveal to
    C, wait for C to inform the secret is known before revealing it to B, and
    again wait for B before revealing the secret to A.

    If B somehow sent a reveal secret before C and D, then the secret will be
    revealed to A, but not C and D, meaning the secret won't be propagated
    forward. Even if D sent a reveal secret at about the same time, the secret
    will only be revealed to B upon confirmation from C.

    Even though B somehow learnt the secret out-of-order N is safe to proceed
    with the protocol, the transitBlocks configuration adds enough time for
    the reveal secrets to propagate backwards and for B to send the balance
    proof. If the proof doesn't arrive in time and the lock's expiration is at
    risk, N won't lose tokens since it knows the secret can go on-chain at any
    time.
这些transfersPair,都是我介入的中介传输
如果我收到密码的时候已经临近上家密码的 reveal timeout,这样的做法会不会造成连锁反应,造成所有的人都主动去链上注册密码呢?
我要不要发送 reveal secret 给上家呢?
*/
/*
 *	Reveal the secret backwards.
 *
 *   This node is named N, suppose there is a mediated transfer with two refund
 *   transfers, one from B and one from C:
 *
 *       A-N-B...B-N-C..C-N-D
 *
 *   Under normal operation N will first learn the secret from D, then reveal to
 *  C, wait for C to inform the secret is known before revealing it to B, and
 *   again wait for B before revealing the secret to A.
 *
 *   If B somehow sent a reveal secret before C and D, then the secret will be
 *   revealed to A, but not C and D, meaning the secret won't be propagated
 *   forward. Even if D sent a reveal secret at about the same time, the secret
 *   will only be revealed to B upon confirmation from C.
 *
 *   Even though B somehow learnt the secret out-of-order N is safe to proceed
 *   with the protocol, the transitBlocks configuration adds enough time for
 *   the reveal secrets to propagate backwards and for B to send the balance
 *   proof. If the proof doesn't arrive in time and the lock's expiration is at
 *   risk, N won't lose tokens since it knows the secret can go on-chain at any
 *   time.
 *
 *   All these transfersPair are mediated transfers I involve in.
 *	 If the time is near reveal timeout set by previous node, does this leads to interlock effect that all participants register secret on-chain?
 */
func eventsForRevealSecret(transfersPair []*mediatedtransfer.MediationPairState, ourAddress common.Address, blockNumber int64) (events []transfer.Event) {
	for j := len(transfersPair) - 1; j >= 0; j-- {
		pair := transfersPair[j]
		isPayeeSecretKnown := stateSecretKnownMaps[pair.PayeeState]
		isPayerSecretKnown := stateSecretKnownMaps[pair.PayerState]
		// 判断是否超时,如果已经该锁已经超时,不发送secret给上家
		isExpired := blockNumber > pair.PayerTransfer.Expiration
		tr := pair.PayerTransfer
		if isPayeeSecretKnown && !isPayerSecretKnown && !isExpired {
			pair.PayerState = mediatedtransfer.StatePayerSecretRevealed
			revealSecret := &mediatedtransfer.EventSendRevealSecret{
				LockSecretHash: tr.LockSecretHash,
				Secret:         tr.Secret,
				Token:          tr.Token,
				Receiver:       pair.PayerRoute.HopNode(),
				Sender:         ourAddress,
			}
			events = append(events, revealSecret)
		}
		if tr.Fee.Cmp(big.NewInt(0)) > 0 {
			events = append(events, &mediatedtransfer.EventSaveFeeChargeRecord{
				LockSecretHash: tr.LockSecretHash,
				TokenAddress:   tr.Token,
				TransferFrom:   tr.Initiator,
				TransferTo:     tr.Target,
				TransferAmount: tr.TargetAmount,
				InChannel:      pair.PayerRoute.ChannelIdentifier,
				OutChannel:     pair.PayeeRoute.ChannelIdentifier,
				Fee:            new(big.Int).Sub(pair.PayerTransfer.Fee, pair.PayeeTransfer.Fee),
				Timestamp:      time.Now().Unix(),
			})
		}
	}
	return events
}

//Send the balance proof to nodes that know the secret.
func eventsForBalanceProof(transfersPair []*mediatedtransfer.MediationPairState, blockNumber int64) (events []transfer.Event) {
	for j := len(transfersPair) - 1; j >= 0; j-- {
		pair := transfersPair[j]
		payeeKnowsSecret := stateSecretKnownMaps[pair.PayeeState]
		payeePayed := stateTransferPaidMaps[pair.PayeeState]
		payeeChannelOpen := pair.PayeeRoute.State() == channeltype.StateOpened
		/*
			如果上家通道已经关闭,不主动发送unlock,尽量让下家去注册密码,避免自己上链注册密码
		*/
		// If previous channel closed, don't send unlock and let next node to register secret
		payerChannelOpen := pair.PayerRoute.State() == channeltype.StateOpened

		/*
			如果我收到密码的时候已经临近上家密码的 reveal timeout,那么安全的做法就是什么都不做.
				强迫这个交易失败,或者强迫下家去注册密码,对方就不应该在临近过期的时候才告诉我密码,应该提早告诉.
		*/
		// If the time I receive secret from my previous node is near reveal_timeout, then we better do nothing.
		// to force failure of this transfer or next node to register secret.
		// My partner should not reveal this secret to me near reveal_timeout instead that he should reveal it ahead.
		payerTransferInDanger := blockNumber > pair.PayerTransfer.Expiration-int64(pair.PayerRoute.RevealTimeout())

		lockValid := isLockValid(pair.PayeeTransfer, blockNumber)
		if payerChannelOpen && payeeChannelOpen && payeeKnowsSecret && !payeePayed && lockValid && !payerTransferInDanger {
			pair.PayeeState = mediatedtransfer.StatePayeeBalanceProof
			tr := pair.PayeeTransfer
			balanceProof := &mediatedtransfer.EventSendBalanceProof{
				LockSecretHash:    tr.LockSecretHash,
				ChannelIdentifier: pair.PayeeRoute.ChannelIdentifier,
				Token:             tr.Token,
				Receiver:          pair.PayeeRoute.HopNode(),
			}
			unlockSuccess := &mediatedtransfer.EventUnlockSuccess{
				LockSecretHash: pair.PayerTransfer.LockSecretHash,
			}
			events = append(events, balanceProof, unlockSuccess)
		}
	}
	return
}

/*
Close the channels that are in the unsafe region prior to an on-chain
    withdraw
发出注册密码事件应该是整个交易的事情,而不是某个 pair 的事
如果发生了崩溃恢复怎么处理呢?
*/
/*
 *	Close the channels that are in the unsafe region prior to an on-chain withdraw
 *	All channel participants should be responsbile to send reveal secret when necessary.
 */
func eventsForRegisterSecret(transfersPair []*mediatedtransfer.MediationPairState, blockNumber int64) (events []transfer.Event) {
	pendings := getPendingTransferPairs(transfersPair)
	needRegisterSecret := false
	for j := len(pendings) - 1; j >= 0; j-- {
		pair := pendings[j]
		if isSecretRegisterNeeded(pair, blockNumber) {
			//只需发出一次注册请求,所有的 pair 状态都应该修改为StatePayerWaitingRegisterSecret
			// we only need to send reveal secret once, all pairs state should switch to StatePayerWaitingRegisterSecret.
			if needRegisterSecret {
				pair.PayerState = mediatedtransfer.StatePayerWaitingRegisterSecret
			} else {
				needRegisterSecret = true
				pair.PayerState = mediatedtransfer.StatePayerWaitingRegisterSecret
				registerSecretEvent := &mediatedtransfer.EventContractSendRegisterSecret{
					Secret: pair.PayeeTransfer.Secret,
				}
				events = append(events, registerSecretEvent)
			}
		}
	}
	return
}

/*
Set the state of the `payeeAddress` transfer, check the secret is
    being revealed backwards, and if necessary send out RevealSecret,
    SendBalanceProof, and Withdraws.
*/
func secretLearned(state *mediatedtransfer.MediatorState, secret common.Hash, payeeAddress common.Address, newPayeeState string) *transfer.TransitionResult {
	if !stateSecretKnownMaps[newPayeeState] {
		panic(fmt.Sprintf("%s not in STATE_SECRET_KNOWN", newPayeeState))
	}
	if state.Secret == utils.EmptyHash {
		state.SetSecret(secret)
	}
	var events []transfer.Event
	eventsWrongOrder := setPayeeStateAndCheckRevealOrder(state.TransfersPair, payeeAddress, newPayeeState)
	eventsSecretReveal := eventsForRevealSecret(state.TransfersPair, state.OurAddress, state.BlockNumber)
	eventBalanceProof := eventsForBalanceProof(state.TransfersPair, state.BlockNumber)
	eventsRegisterSecretEvent := eventsForRegisterSecret(state.TransfersPair, state.BlockNumber)
	events = append(events, eventsWrongOrder...)
	events = append(events, eventsSecretReveal...)
	events = append(events, eventBalanceProof...)
	events = append(events, eventsRegisterSecretEvent...)
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

/*
Try a new route or fail back to a refund.

    The mediator can safely try a new route knowing that the tokens from
    payer_transfer will cover the expenses of the mediation. If there is no
    route available that may be used at the moment of the call the mediator may
    send a refund back to the payer, allowing the payer to try a different
    route.
*/
func mediateTransfer(state *mediatedtransfer.MediatorState, payerRoute *route.State, payerTransfer *mediatedtransfer.LockedTransferState) *transfer.TransitionResult {
	var transferPair *mediatedtransfer.MediationPairState
	var events []transfer.Event

	timeoutBlocks := int(getTimeoutBlocks(payerRoute, payerTransfer, state.BlockNumber))
	//log.Trace(fmt.Sprintf("timeoutBlocks=%d,payerroute=%s,payertransfer=%s,blocknumber=%d",
	//	timeoutBlocks, utils.StringInterface(payerRoute, 3), utils.StringInterface(payerTransfer, 3),
	//	state.BlockNumber,
	//))
	if timeoutBlocks > 0 {
		transferPair, events = nextTransferPair(payerRoute, payerTransfer, state.Routes, timeoutBlocks, state.BlockNumber)
	}
	if transferPair == nil {
		log.Warn("no usable route, reject")
		/*
			回退此交易,相当于没收到过一样处理

		*/
		/*
		 *	Exit this transfer, like never received it.
		 */
		originalTransfer := payerTransfer
		originalRoute := payerRoute
		refundEvents := eventsForRefund(originalRoute, originalTransfer)
		return &transfer.TransitionResult{
			NewState: state,
			Events:   refundEvents,
		}
	}
	/*
		这里根据reveal_timeout限制持有锁的数量,如果同时持有对方的锁大于某个值,说明对方很有可能不是诚信节点,此时不再接收对方发来的交易,
		这是为了避免上下家合作利用时间差来攻击我,造成我损失钱的情况
		在这里做可以直接通过announce disposed拒绝该笔交易,既达到拒绝的目的,也不至于让交易卡住
		暂时直接使用通道的reveal_timeout作为阈值
	*/
	/*
		Here, the number of locks held is limited according to reveal_timeout. If the locks held at the same time are more than a certain value,
		the other party may not be a trustworthy node and no longer receive transactions from the other party.
		This is to avoid the cooperation between the payer and payee in using time difference to attack me, causing me to lose money.

		It is possible to reject the transaction directly through announcement disposed, so that the purpose of rejection is achieved and the transaction is not stuck.
		Reveal_timeout is used directly as a threshold value temporarily.
	*/
	payerChannel := transferPair.PayerRoute.Channel()
	if len(payerChannel.PartnerState.Lock2PendingLocks)+len(payerChannel.PartnerState.Lock2UnclaimedLocks) > payerChannel.RevealTimeout {
		log.Warn(fmt.Sprintf("holding too much lock of %s, reject new mediated transfer from him", utils.APex2(payerChannel.PartnerState.Address)))
		return &transfer.TransitionResult{
			NewState: state,
			Events:   eventsForRefund(payerRoute, payerTransfer),
		}
	}
	/*
				   the list must be ordered from high to low expiration, expiration
		         handling depends on it
	*/
	state.TransfersPair = append(state.TransfersPair, transferPair)
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

/*

 */
func cancelCurrentRoute(state *mediatedtransfer.MediatorState, refundChannelIdentify common.Hash) *transfer.TransitionResult {
	var it = &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
	l := len(state.TransfersPair)
	if l <= 0 {
		log.Error(fmt.Sprintf("recevie refund ,but has no transfer pair ,must be a attack!!"))
		return it
	}
	transferPair := state.TransfersPair[l-1]
	state.TransfersPair = state.TransfersPair[:l-1] //移除最后一个
	/*
		if refund msg came from payer, panic, something must wrong!
	*/
	if refundChannelIdentify == transferPair.PayerRoute.ChannelIdentifier {
		panic("receive refund/withdraw/cooperateSettle from payer,that should happen")
	}
	/*
		这里需要判断下payer通道的状态,如果该通道的状态已经不为open了,就不应该继续尝试新路由.
		因为我已经在通道close时提交过balance proof,如果继续这笔交易,并且最终交易成功的话,我需要自己注册密码并且上链unlock
		这不仅需要花费gas,而且存在没来得及unlock以至于丢钱的风险,所以没有必要继续,直接拒绝
	*/
	/*
		Here I need to determine the state of the payer channel, and if the state of the channel is no longer open, you should not continue to try new routing.
		Because I've already submitted balance proof at channel close, and if this transaction continues and eventually succeeds,
		I need to register secret and unlock it.
		Not only does it cost gas, but there's a risk of losing money by not unlocking it in time, so there's no need to go on.
		Refuse by announce disposed.
	*/
	if transferPair.PayerRoute.ClosedBlock() != 0 {
		log.Warn("channel already closed, stop trying new route")
		it.Events = eventsForRefund(transferPair.PayerRoute, transferPair.PayerTransfer)
		return it
	}
	it = mediateTransfer(state, transferPair.PayerRoute, transferPair.PayerTransfer)
	return it
}

/*
又收到了一个 mediatedtransfer
*/
// receive another mediatedTransfer
func handleMediatedTransferAgain(state *mediatedtransfer.MediatorState, st *mediatedtransfer.MediatorReReceiveStateChange) *transfer.TransitionResult {
	return mediateTransfer(state, st.FromRoute, st.FromTransfer)
}

/*
After Photon learns about a new block this function must be called to
    handle expiration of the hash time locks.
        state : The current state.

    Return:
        TransitionResult: The resulting iteration
*/
func handleBlock(state *mediatedtransfer.MediatorState, st *transfer.BlockStateChange) *transfer.TransitionResult {
	blockNumber := state.BlockNumber
	if blockNumber < st.BlockNumber {
		blockNumber = st.BlockNumber
	}
	state.BlockNumber = blockNumber
	closeEvents := eventsForRegisterSecret(state.TransfersPair, blockNumber)
	unlockfailEvents, allExpired := setExpiredPairs(state.TransfersPair, blockNumber)
	var events []transfer.Event
	events = append(events, closeEvents...)
	events = append(events, unlockfailEvents...)
	if allExpired {
		//所有的mediatedtransfer 都已经过期了,放心移除这个 stateManager 吧
		// All mediatedTransfers are expired, feel safe to remove this StateManager.
		events = append(events, &mediatedtransfer.EventRemoveStateManager{
			Key: utils.Sha3(state.LockSecretHash[:], state.Token[:]),
		})
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

/*
Validate and handle a ReceiveTransferRefund state change.

    A node might participate in mediated transfer more than once because of
    refund transfers, eg. A-B-C-F-B-D-T, B tried to mediate the transfer through
    C, which didn't have an available route to proceed and refunds B, at this
    point B is part of the path again and will try a new partner to proceed
    with the mediation through D, D finally reaches the target T.

    In the above scenario B has two pairs of payer and payee transfers:

        payer:A payee:C from the first SendMediatedTransfer
        payer:F payee:D from the following SendRefundTransfer

        state : Current state.
        st : The state change.

    Returns:
        TransitionResult: The resulting iteration.
*/
func handleAnnouceDisposed(state *mediatedtransfer.MediatorState, st *mediatedtransfer.ReceiveAnnounceDisposedStateChange) *transfer.TransitionResult {
	it := &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
	if state.Secret != utils.EmptyHash {
		panic("refunds are not allowed if the secret is revealed")
	}
	/*
			  The last sent transfer is the only one thay may be refunded, all the
		     previous ones are refunded already.

	*/
	l := len(state.TransfersPair)
	if l <= 0 {
		log.Error(fmt.Sprintf("recevie refund ,but has no transfer pair ,must be a attack!!"))
		return it
	}
	transferPair := state.TransfersPair[l-1]
	payeeTransfer := transferPair.PayeeTransfer
	payeeRoute := transferPair.PayeeRoute

	/*
			A-B-C-F-B-G-D
			B首先收到了来自 C 的refund 怎么处理?网络错误?网络攻击?

		考虑支持分叉,有可能在锁过期之后收到来自payee的AnnounceDisposed,
		这时候虽然可以尝试找其他路径,但是更安全的做法是忽略,当然即使不忽略也不应该有安全问题.
		这里采取更安全的做法: 如果payee的锁已经过期,直接忽略.

	*/
	// A-B-C-F-B-G-D
	// If B first receives refund of C, how to deal with that?
	if IsValidRefund(payeeTransfer, payeeRoute, st) {
		if payeeTransfer.Expiration > state.BlockNumber {
			/*
					假定队列中的是
				AB BC
				EB BF
				这时候收到了来自 F 的 refund, 那么应该是认为 payeeTransfer 无效了,
				相当于刚刚收到了来自 E 的 transfer, 然后去重新找路径

			*/
			/*
			 *	Assume that order in this queue is
			 *	AB BC EB BF
			 *	which means we receive refund of F, then we should assume that payeeTransfer invalid,
			 *  which acts like receiving transfer of E, then begin to find a route again.
			 */
			it = cancelCurrentRoute(state, st.Message.ChannelIdentifier)
			ev := &mediatedtransfer.EventSendAnnounceDisposedResponse{
				Token:          state.Token,
				LockSecretHash: st.Lock.LockSecretHash,
				Receiver:       st.Sender,
			}
			it.Events = append(it.Events, ev)
		} else {
			log.Warn(fmt.Sprintf("receive expired EventSendAnnounceDisposedResponse,expiration=%d,currentblock=%d,response=%s",
				payeeTransfer.Expiration, state.BlockNumber, utils.StringInterface(st, 3),
			))
		}

	}
	return it
}

/*
Validate and handle a ReceiveSecretReveal state change.

    The Secret must propagate backwards through the chain of mediators, this
    function will record the learned secret, check if the secret is propagating
    backwards (for the known paths), and send the SendBalanceProof/RevealSecret if
    necessary.
*/
func handleSecretReveal(state *mediatedtransfer.MediatorState, st *mediatedtransfer.ReceiveSecretRevealStateChange) *transfer.TransitionResult {
	secret := st.Secret
	if utils.ShaSecret(secret[:]) != state.Hashlock {
		panic("must a implementation error")
	}
	return secretLearned(state, secret, st.Sender, mediatedtransfer.StatePayeeSecretRevealed)
}

/*
收到密码在链上注册,
此次交易就彻底结束了,
已经过期的交易,对方应该已经收到了 remove expired,
没有过期的交易都还在这里
*/
/*
 *	Received secret has been registered on-chain,
 *	this transfer completes.
 *	For those expired transfers, partner should receive remove expired Hashlock
 *	For those unexpired transfers, they are here.
 */
func handleSecretRevealOnChain(state *mediatedtransfer.MediatorState, st *mediatedtransfer.ContractSecretRevealOnChainStateChange) *transfer.TransitionResult {
	var it = &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
	var events []transfer.Event
	if state.LockSecretHash != st.LockSecretHash {
		panic("impementation error")
	}
	state.SetSecret(st.Secret)
	for _, pair := range state.TransfersPair {
		if true {
			tr := pair.PayeeTransfer
			route := pair.PayeeRoute
			//针对下家,有效的链上注册,应该发送unlock
			// 如果通道已经关闭,可以不发,但是发了也没什么坏处,只不过是消息多发几遍,然后丢弃而已.
			//因为如果通道已经关闭,说明下家已经提交了BalanceProof,他也不会接受新的balance proof了.
			if tr.Expiration >= st.BlockNumber {
				//没有超时,就应该发送 unlock 消息,不用关心现在通道是什么状态,就是 settle 了问题也不大.
				// if not be reveal_timeout, then we should send unlock, do not care current channel state.
				ev := &mediatedtransfer.EventSendBalanceProof{
					LockSecretHash:    tr.LockSecretHash,
					ChannelIdentifier: route.ChannelIdentifier,
					Token:             tr.Token,
					Receiver:          route.HopNode(),
				}
				events = append(events, ev)
				pair.PayeeState = mediatedtransfer.StatePayeeBalanceProof
				//至于 payer 一方,不发送也不影响我所得,需要浪费 gas 进行链上兑现.
				// As for payer, he will not be impacted even he does not send BalanceProof, but cost gas to on-chain secret register.

				// 没有超时即确认收到了手续费,记录流水
				if tr.Fee.Cmp(big.NewInt(0)) > 0 {
					events = append(events, &mediatedtransfer.EventSaveFeeChargeRecord{
						LockSecretHash: tr.LockSecretHash,
						TokenAddress:   tr.Token,
						TransferFrom:   tr.Initiator,
						TransferTo:     tr.Target,
						TransferAmount: tr.TargetAmount,
						InChannel:      pair.PayerRoute.ChannelIdentifier,
						OutChannel:     pair.PayeeRoute.ChannelIdentifier,
						Fee:            pair.PayerRoute.Fee,
						Timestamp:      time.Now().Unix(),
					})
				}
			}
		}
		if true { //针对上家通道
			route := pair.PayerRoute
			//针对上家,如果通道已经关闭了,说明自己已经提交过balance proof,并且尝试unlock过,当然unlock是失败的
			if route.State() == channeltype.StateClosed {
				events = append(events, &mediatedtransfer.EventContractSendUnlock{
					LockSecretHash:    st.LockSecretHash,
					ChannelIdentifier: route.ChannelIdentifier,
				})
				pair.PayerState = mediatedtransfer.StatePayerBalanceProof
			}
		}
	}
	it.Events = events
	return it
}

// Handle a ReceiveBalanceProof state change.
func handleBalanceProof(state *mediatedtransfer.MediatorState, st *mediatedtransfer.ReceiveUnlockStateChange) *transfer.TransitionResult {
	var events []transfer.Event
	//走到这里,说明密码对了,金额也是对的,balanceProof也是对的,但是我有多个pair,需要找到对应的那个pair
	for _, pair := range state.TransfersPair {
		//如果不做检查,假设我是B,比如A-B-C-D-B-E ,那么B收到来自A的unlock消息以后,就会把unlock消息同时发送给C,E,这是错误的.
		if pair.PayerRoute.HopNode() != st.NodeAddress {
			continue
		}
		if pair.PayeeState != mediatedtransfer.StatePayeeBalanceProof {
			/*
				如果收到unlock的时候,还没有给下家发送unlock,补发
			*/
			// If not receiving unlock from next node when receiving unlock, re-send it.
			pair.PayeeState = mediatedtransfer.StatePayeeBalanceProof
			tr := pair.PayeeTransfer
			balanceProof := &mediatedtransfer.EventSendBalanceProof{
				LockSecretHash:    tr.LockSecretHash,
				ChannelIdentifier: pair.PayeeRoute.ChannelIdentifier,
				Token:             tr.Token,
				Receiver:          pair.PayeeRoute.HopNode(),
			}
			unlockSuccess := &mediatedtransfer.EventUnlockSuccess{
				LockSecretHash: pair.PayerTransfer.LockSecretHash,
			}
			events = append(events, balanceProof, unlockSuccess)
		}

		withdraw := &mediatedtransfer.EventWithdrawSuccess{
			LockSecretHash: pair.PayeeTransfer.LockSecretHash,
		}
		events = append(events, withdraw)
		pair.PayerState = mediatedtransfer.StatePayerBalanceProof
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

//StateTransition is State machine for a node mediating a transfer.
func StateTransition(originalState transfer.State, stateChange transfer.StateChange) (it *transfer.TransitionResult) {
	/*
			  Notes:
		     - A user cannot cancel a mediated transfer after it was initiated, she
		       may only reject to mediate before hand. This is because the mediator
		       doesn't control the secret reveal and needs to wait for the lock
		       expiration before safely discarding the transfer.
	*/
	it = &transfer.TransitionResult{
		NewState: originalState,
		Events:   nil,
	}
	state, ok := originalState.(*mediatedtransfer.MediatorState)
	if !ok {
		if originalState != nil {
			panic("MediatorState StateTransition get type error ")
		}
		state = nil
	}
	if state == nil {
		if aim, ok := stateChange.(*mediatedtransfer.ActionInitMediatorStateChange); ok {
			state = &mediatedtransfer.MediatorState{
				OurAddress:     aim.OurAddress,
				Routes:         aim.Routes,
				BlockNumber:    aim.BlockNumber,
				Hashlock:       aim.FromTranfer.LockSecretHash,
				Db:             aim.Db,
				Token:          aim.FromTranfer.Token,
				LockSecretHash: aim.FromTranfer.LockSecretHash,
			}
			it = mediateTransfer(state, aim.FromRoute, aim.FromTranfer)
		}
	} else {
		switch st2 := stateChange.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
		case *mediatedtransfer.ReceiveAnnounceDisposedStateChange:
			if state.Secret == utils.EmptyHash {
				it = handleAnnouceDisposed(state, st2)
			} else {
				log.Error(fmt.Sprintf("mediator state manager ,already knows secret,but recevied announce disposed, must be a error"))
			}

		case *mediatedtransfer.ReceiveSecretRevealStateChange:
			it = handleSecretReveal(state, st2)
		case *mediatedtransfer.ContractSecretRevealOnChainStateChange:
			it = handleSecretRevealOnChain(state, st2)
		case *mediatedtransfer.ReceiveUnlockStateChange:
			it = handleBalanceProof(state, st2)
			if state.Secret == utils.EmptyHash {
				log.Warn(fmt.Sprintf("mediated state manager recevie unlock,but i don't know secret,this maybe a error "))
			}
		case *mediatedtransfer.MediatorReReceiveStateChange:
			if state.Secret == utils.EmptyHash {
				it = handleMediatedTransferAgain(state, st2)
			} else {
				log.Error(fmt.Sprintf("already known secret,but recevie medaited tranfer again:%s", st2.Message))
			}
		/*
			only receive from channel with payee,
			never receive from channel with payer
		*/
		case *mediatedtransfer.ContractCooperativeSettledStateChange:
			it = cancelCurrentRoute(state, st2.ChannelIdentifier)
		case *mediatedtransfer.ContractChannelWithdrawStateChange:
			it = cancelCurrentRoute(state, st2.ChannelIdentifier.ChannelIdentifier)
		default:
			log.Info(fmt.Sprintf("unknown statechange :%s", utils.StringInterface(st2, 3)))
		}
	}
	// this is the place for paranoia
	if it.NewState != nil {
		sanityCheck(it.NewState.(*mediatedtransfer.MediatorState))
	}
	return clearIfFinalized(it)
}
