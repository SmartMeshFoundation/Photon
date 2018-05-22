package mediator

import (
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
const transitBlocks = 2 // TODO: make this a configuration variable

var stateSecretKnownMaps = map[string]bool{
	mediatedtransfer.StatePayeeSecretRevealed:   true,
	mediatedtransfer.StatePayeeRefundWithdraw:   true,
	mediatedtransfer.StatePayeeContractWithdraw: true,
	mediatedtransfer.StatePayeeBalanceProof:     true,

	mediatedtransfer.StatePayerSecretRevealed:   true,
	mediatedtransfer.StatePayerWaitingClose:     true,
	mediatedtransfer.StatePayerWaitingWithdraw:  true,
	mediatedtransfer.StatePayerContractWithdraw: true,
	mediatedtransfer.StatePayerBalanceProof:     true,
}
var stateTransferPaidMaps = map[string]bool{
	mediatedtransfer.StatePayeeContractWithdraw: true,
	mediatedtransfer.StatePayeeBalanceProof:     true,
	mediatedtransfer.StatePayerContractWithdraw: true,
	mediatedtransfer.StatePayerBalanceProof:     true,
}

//TODO: fix expired state, it is not final
var stateTransferFinalMaps = map[string]bool{
	mediatedtransfer.StatePayeeExpired:          true,
	mediatedtransfer.StatePayeeContractWithdraw: true,
	mediatedtransfer.StatePayeeBalanceProof:     true,

	mediatedtransfer.StatePayerExpired:          true,
	mediatedtransfer.StatePayerContractWithdraw: true,
	mediatedtransfer.StatePayerBalanceProof:     true,
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
func IsValidRefund(originTr, refundTr *mediatedtransfer.LockedTransferState, refundSender common.Address) bool {
	//Ignore a refund from the target
	if refundSender == originTr.Target {
		return false
	}
	return originTr.Identifier == refundTr.Identifier &&
		originTr.Amount.Cmp(refundTr.Amount) == 0 &&
		originTr.Hashlock == refundTr.Hashlock &&
		originTr.Target == refundTr.Target &&
		originTr.Token == refundTr.Token &&
		originTr.Initiator == refundTr.Initiator &&
		/*
					 A larger-or-equal expiration is byzantine behavior that favors this
			         node, neverthless it's being ignored since the only reason for the
			         other node to use an invalid expiration is to play the protocol.
		*/
		originTr.Expiration > refundTr.Expiration
}

/*
True if this node needs to close the channel to withdraw on-chain.

    Only close the channel to withdraw on chain if the corresponding payee node
    has received, this prevents attacks were the payee node burns it's payment
    to force a close with the payer channel.
*/
func isChannelCloseNeeded(tr *mediatedtransfer.MediationPairState, blockNumber int64) bool {
	payeeReceived := stateTransferPaidMaps[tr.PayeeState]
	payerPayed := stateTransferPaidMaps[tr.PayerState]
	payerChannelOpen := tr.PayerRoute.State == transfer.ChannelStateOpened
	AlreadyClosing := tr.PayerState == mediatedtransfer.StatePayerWaitingClose
	safeToWait := IsSafeToWait(tr.PayerTransfer, tr.PayerRoute.RevealTimeout, blockNumber)

	return payeeReceived && !payerPayed && payerChannelOpen && !AlreadyClosing && !safeToWait
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

/*
Return the timeout blocks, it's the base value from which the payee's
    lock timeout must be computed.

    The payee lock timeout is crucial for safety of the mediated transfer, the
    value must be chosen so that the payee hop is forced to reveal the secret
    with sufficient time for this node to claim the received lock from the
    payer hop.

    The timeout blocks must be the smallest of:

    - payer_transfer.expiration: The payer lock expiration, to force the payee
      to reveal the secret before the lock expires.
    - payer_route.settle_timeout: Lock expiration must be lower than
      the settlement period since the lock cannot be claimed after the channel is
      settled.
    - payer_route.closed_block: If the channel is closed then the settlement
      period is running and the lock expiration must be lower than number of
      blocks left.
这个函数实际上是说下一跳最多还有多少块时间可利用。
*/
func getTimeoutBlocks(payerRoute *transfer.RouteState, payerTransfer *mediatedtransfer.LockedTransferState, blockNumber int64) int64 {
	blocksUntilSettlement := int64(payerRoute.SettleTimeout)
	if payerRoute.ClosedBlock != 0 {
		if blockNumber < payerRoute.ClosedBlock {
			panic("ClosedBlock bigger than the lastest blocknumber")
		}
		blocksUntilSettlement -= blockNumber - payerRoute.ClosedBlock
	}
	if blocksUntilSettlement > payerTransfer.Expiration-blockNumber {
		blocksUntilSettlement = payerTransfer.Expiration - blockNumber
	}
	return blocksUntilSettlement - transitBlocks
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
			panic("payer:a transfer is paid but we don't know the secret")
		}
		if stateTransferPaidMaps[pair.PayeeState] && state.Secret == utils.EmptyHash {
			panic("payee:a transfer is paid but we don't know the secret")
		}
	}
	//the "transitivity" for these values is checked below as part of
	//almost_equal check
	if len(state.TransfersPair) > 0 {
		firstPair := state.TransfersPair[0]
		if state.Hashlock != firstPair.PayerTransfer.Hashlock {
			panic("sanity check failed:state.Hashlock!=firstPair.PayerTransfer.Hashlock")
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
		if p.PayerTransfer.Expiration <= p.PayeeTransfer.Expiration {
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
		if original.PayeeRoute.HopNode != refund.PayerRoute.HopNode {
			panic("sanity check failed:original.PayeeRoute.HopNode!=refund.PayerRoute.HopNode")
		}
		if original.PayeeTransfer.Expiration <= refund.PayerTransfer.Expiration {
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
	/*
			  TODO: clear the expired transfer, this will need some sort of
		     synchronization among the nodes
	*/
	isAllFinalized := true
	for _, p := range state.TransfersPair {
		if !stateTransferPaidMaps[p.PayeeState] || !stateTransferPaidMaps[p.PayerState] {
			isAllFinalized = false
			break
		}
	}
	if isAllFinalized {
		return &transfer.TransitionResult{
			NewState: nil,
			Events:   result.Events,
		}
	}
	return result
}

/*
Finds the first route available that may be used.

    Args:
        rss (RoutesState): Current available routes that may be used,
            it's assumed that the available_routes list is ordered from best to
            worst.
        timeout_blocks (int): Base number of available blocks used to compute
            the lock timeout.
        transfer_amount (int): The amount of tokens that will be transferred
            through the given route.

    Returns:
        (RouteState): The next route.
*/
func nextRoute(rss *transfer.RoutesState, timeoutBlocks int, transferAmount *big.Int) *transfer.RouteState {
	for len(rss.AvailableRoutes) > 0 {
		route := rss.AvailableRoutes[0]
		rss.AvailableRoutes = rss.AvailableRoutes[1:]
		lockTimeout := timeoutBlocks - route.RevealTimeout
		if route.AvaibleBalance.Cmp(transferAmount) >= 0 && lockTimeout > 0 {
			return route
		}
		rss.IgnoredRoutes = append(rss.IgnoredRoutes, route)
	}
	return nil
}

/*
Given a payer transfer tries a new route to proceed with the mediation.

    Args:
        payer_route (RouteState): The previous route in the path that provides
            the token for the mediation.
        payer_transfer (LockedTransferState): The transfer received from the
            payer_route.
        routes_state (RoutesState): Current available routes that may be used,
            it's assumed that the available_routes list is ordered from best to
            worst.
        timeout_blocks (int): Base number of available blocks used to compute
            the lock timeout.
        block_number (int): The current block number.
*/
func nextTransferPair(payerRoute *transfer.RouteState, payerTransfer *mediatedtransfer.LockedTransferState,
	routesState *transfer.RoutesState, timeoutBlocks int, blockNumber int64) (
	transferPair *mediatedtransfer.MediationPairState, events []transfer.Event) {
	if timeoutBlocks <= 0 {
		panic("timeoutBlocks<=0")
	}
	if int64(timeoutBlocks) > payerTransfer.Expiration-blockNumber {
		panic("timeoutBlocks >payerTransfer.Expiration-blockNumber")
	}
	payeeRoute := nextRoute(routesState, timeoutBlocks, payerTransfer.Amount)
	if payeeRoute != nil {
		if payeeRoute.RevealTimeout >= timeoutBlocks {
			panic("payeeRoute.RevealTimeout>=timeoutBlocks")
		}
		/*
			有可能 payeeroute 的 settle timeout 比较小,从而导致我指定的lockexpiration 特别大,从而对我不利.
		*/
		if timeoutBlocks >= payeeRoute.SettleTimeout {
			timeoutBlocks = payeeRoute.SettleTimeout
		}
		lockTimeout := timeoutBlocks - payeeRoute.RevealTimeout
		lockExpiration := int64(lockTimeout) + blockNumber
		payeeTransfer := &mediatedtransfer.LockedTransferState{
			Identifier:   payerTransfer.Identifier,
			TargetAmount: payerTransfer.TargetAmount,
			Amount:       big.NewInt(0).Sub(payerTransfer.Amount, payeeRoute.Fee),
			Token:        payerTransfer.Token,
			Initiator:    payerTransfer.Initiator,
			Target:       payerTransfer.Target,
			Expiration:   lockExpiration,
			Hashlock:     payerTransfer.Hashlock,
			Secret:       payerTransfer.Secret,
			Fee:          big.NewInt(0).Sub(payerTransfer.Fee, payeeRoute.Fee),
		}
		if payeeTransfer.Fee.Cmp(utils.BigInt0) < 0 || payeeTransfer.Amount.Cmp(utils.BigInt0) < 0 {
			//no enough fee to route.
			return
		}
		if payeeRoute.HopNode == payeeTransfer.Target {
			//i'm the last hop,so take the rest of the fee
			payeeTransfer.Fee = utils.BigInt0
			payeeTransfer.Amount = payerTransfer.TargetAmount
		}
		//log how many tokens fee for this transfer . todo
		transferPair = mediatedtransfer.NewMediationPairState(payerRoute, payeeRoute, payerTransfer, payeeTransfer)
		events = []transfer.Event{mediatedtransfer.NewEventSendMediatedTransfer(payeeTransfer, payeeRoute.HopNode)}
	}
	return
}

/*
Set the state of a transfer *sent* to a payee and check the secret is
    being revealed backwards.

    Note:
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
		if back.PayeeRoute.HopNode == payeeAddress {
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

// Set the state of expired transfers, and return the failed events
func setExpiredPairs(transfersPairs []*mediatedtransfer.MediationPairState, blockNumber int64) (events []transfer.Event) {
	pendingTransfersPairs := getPendingTransferPairs(transfersPairs)
	for _, pair := range pendingTransfersPairs {
		if blockNumber > pair.PayerTransfer.Expiration {
			if pair.PayeeState != mediatedtransfer.StatePayeeExpired {
				log.Error("PayeeState!=mediatedtransfer.StatePayeeExpired")
				return
			}
			if pair.PayeeTransfer.Expiration >= pair.PayerTransfer.Expiration {
				log.Error("PayeeTransfer.Expiration>=pair.PayerTransfer.Expiration")
				return
			}
			if pair.PayerState != mediatedtransfer.StatePayerExpired {
				pair.PayerState = mediatedtransfer.StatePayerExpired
				withdrawFailed := &mediatedtransfer.EventWithdrawFailed{
					Identifier:     pair.PayerTransfer.Identifier,
					Hashlock:       pair.PayerTransfer.Hashlock,
					ChannelAddress: pair.PayerRoute.ChannelAddress,
					Reason:         "lock expired",
				}
				events = append(events, withdrawFailed)
			}
		} else if blockNumber > pair.PayeeTransfer.Expiration {
			/*
			   For safety, the correct behavior is:

			   - If the payee has been paid, then the payer must pay too.

			     And the corollary:

			   - If the payer transfer has expired, then the payee transfer must
			     have expired too.

			   The problem is that this corollary cannot be asserted. If a user
			   is running Raiden without a monitoring service, then it may go
			   offline after having paid a transfer to a payee, but without
			   getting a balance proof of the payer, and once it comes back
			   online the transfer may have expired.

			   assert pair.payee_state == 'payee_expired'
			*/
			//if stateTransferPaidMaps[pair.PayeeState] {
			//	panic("pair.payee_state should not in STATE_TRANSFER_PAID")
			//}
			if pair.PayeeTransfer.Expiration >= pair.PayerTransfer.Expiration {
				panic("PayeeTransfer.Expiration>=pair.PayerTransfer.Expiration")
			}
			if pair.PayeeState != mediatedtransfer.StatePayeeExpired {
				pair.PayeeState = mediatedtransfer.StatePayeeExpired
				unlockFailed := &mediatedtransfer.EventUnlockFailed{
					Identifier:     pair.PayeeTransfer.Identifier,
					Hashlock:       pair.PayeeTransfer.Hashlock,
					ChannelAddress: pair.PayeeRoute.ChannelAddress,
					Reason:         "lock expired",
				}
				events = append(events, unlockFailed)
			}
		}
	}
	return
}

/*
Refund the transfer.

    Args:
        refund_route (RouteState): The original route that sent the mediated
            transfer to this node.
        refund_transfer (LockedTransferState): The original mediated transfer
            from the refund_route.
        timeout_blocks (int): The number of blocks available from the /latest
            transfer/ received by this node, this transfer might be the
            original mediated transfer (if no route was available) or a refund
            transfer from a down stream node.
        block_number (int): The current block number.

    Returns:
        An empty list if there are not enough blocks to safely create a refund,
        or a list with a refund event.
*/
func eventsForRefundTransfer(refundRoute *transfer.RouteState, refundTransfer *mediatedtransfer.LockedTransferState, timeoutBlocks int, blockNumber int64) (events []transfer.Event) {
	/*
	   A refund transfer works like a special SendMediatedTransfer, so it must
	      follow the same rules and decrement reveal_timeout from the
	      payee_transfer.
	*/
	newLockTimeout := timeoutBlocks - refundRoute.RevealTimeout
	if newLockTimeout > 0 {
		newLockExpiration := int64(newLockTimeout) + blockNumber
		rtr2 := &mediatedtransfer.EventSendRefundTransfer{
			Identifier:   refundTransfer.Identifier,
			Token:        refundTransfer.Token,
			Amount:       new(big.Int).Set(refundTransfer.Amount),
			TargetAmount: refundTransfer.TargetAmount,
			Fee:          refundTransfer.Fee,
			HashLock:     refundTransfer.Hashlock,
			Initiator:    refundTransfer.Initiator,
			Target:       refundTransfer.Target,
			Expiration:   newLockExpiration,
			Receiver:     refundRoute.HopNode,
		}
		events = append(events, rtr2)
	}
	/*
	    Can not create a refund lock with a safe expiration, so don't do anything
	   and wait for the received lock to expire.
	*/
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
*/
func eventsForRevealSecret(transfersPair []*mediatedtransfer.MediationPairState, ourAddress common.Address) (events []transfer.Event) {
	for j := len(transfersPair) - 1; j >= 0; j-- {
		pair := transfersPair[j]
		isPayeeSecretKnown := stateSecretKnownMaps[pair.PayeeState]
		isPayerSecretKnown := stateSecretKnownMaps[pair.PayerState]
		if isPayeeSecretKnown && !isPayerSecretKnown { //todo 如果我在发送reveal secret过程中崩溃了，这个消息payer没有收到，payee也没有收到ack，payee会重发revealsecret消息，我的内部状态一定不能变。目前好像是这么做的。
			pair.PayerState = mediatedtransfer.StatePayerSecretRevealed
			tr := pair.PayerTransfer
			revealSecret := &mediatedtransfer.EventSendRevealSecret{
				Identifier: tr.Identifier,
				Secret:     tr.Secret,
				Token:      tr.Token,
				Receiver:   pair.PayerRoute.HopNode,
				Sender:     ourAddress,
			}
			events = append(events, revealSecret)
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
		payeeChannelOpen := pair.PayeeRoute.State == transfer.ChannelStateOpened

		/*
					  todo: All nodes must close the channel and withdraw on-chain if the
			         lock is nearing it's expiration block, what should be the strategy
			         for sending a balance proof to a node that knowns the secret but has
			         not gone on-chain while near the expiration? (The problem is how to
			         define the unsafe region, since that is a local configuration)
		*/
		lockValid := isLockValid(pair.PayeeTransfer, blockNumber)
		if payeeChannelOpen && payeeKnowsSecret && !payeePayed && lockValid {
			pair.PayeeState = mediatedtransfer.StatePayeeBalanceProof
			tr := pair.PayeeTransfer
			balanceProof := &mediatedtransfer.EventSendBalanceProof{
				Identifier:     tr.Identifier,
				ChannelAddress: pair.PayeeRoute.ChannelAddress,
				Token:          tr.Token,
				Receiver:       pair.PayeeRoute.HopNode,
				Secret:         tr.Secret,
			}
			unlockSuccess := &mediatedtransfer.EventUnlockSuccess{
				Identifier: pair.PayerTransfer.Identifier,
				Hashlock:   pair.PayerTransfer.Hashlock,
			}
			events = append(events, balanceProof, unlockSuccess)
		}
	}
	return
}

/*
Close the channels that are in the unsafe region prior to an on-chain
    withdraw
*/
func eventsForClose(transfersPair []*mediatedtransfer.MediationPairState, blockNumber int64) (events []transfer.Event) {
	pendings := getPendingTransferPairs(transfersPair)
	for j := len(pendings) - 1; j >= 0; j-- {
		pair := pendings[j]
		if isChannelCloseNeeded(pair, blockNumber) {
			pair.PayerState = mediatedtransfer.StatePayerWaitingClose
			channelClose := &mediatedtransfer.EventContractSendChannelClose{
				ChannelAddress: pair.PayerRoute.ChannelAddress,
				Token:          pair.PayerTransfer.Token,
			}
			events = append(events, channelClose)
		}
	}
	return
}

/*
Withdraw on any payer channel that is closed and the secret is known.

    If a channel is closed because of another task a balance proof will not be
    received, so there is no reason to wait for the unsafe region before
    calling close.

    This may break the reverse reveal order:

        Path: A -- B -- C -- B -- D

        B learned the secret from D and has revealed to C.
        C has not confirmed yet.
        channel(A, B).closed is True.
        B will withdraw on channel(A, B) before C's confirmation.
        A may learn the secret faster than other nodes.
*/
func eventsForWithdraw(transfersPair []*mediatedtransfer.MediationPairState, db channel.Db) (events []transfer.Event) {
	pendings := getPendingTransferPairs(transfersPair)
	for _, pair := range pendings {
		payerChannelOpen := pair.PayerRoute.State == transfer.ChannelStateOpened
		secretKnown := pair.PayerTransfer.Secret != utils.EmptyHash
		/*
				todo 这个消息应该会反复触发多次，加上限制，只触发一次就好了。。
				由于ContractReceiveWithdrawStateChange 并不派发，导致会不停的发送EventContractSendWithdraw，目前raiden的选择是忽略该事件，但是带来一个潜在安全性问题
			比如：
			A-B-C-B-E-D
			E拿到秘钥以后立即通知c(并把密钥告诉c），但是等一段时间再告诉B
			这时候C选择关闭B-C通道，然后在B-C通道上withdraw。
			这时候由于B在接到关闭B-C通道的时候，并不知道密钥，导致无法withdraw C给B的token，但是C拿到了B给C的token。
			由于忽略EventContractSendWithdraw事件，导致B会丢失C给B点token
		*/

		if !payerChannelOpen && secretKnown {
			//这个是临时补丁，后续要完整解决
			if db != nil {
				//避免无限制重发EventContractSendWithdraw
				if db.IsThisLockHasWithdraw(pair.PayerRoute.ChannelAddress, pair.PayerTransfer.Secret) {
					pair.PayerState = mediatedtransfer.StatePayerContractWithdraw
					continue
				}
			}
			pair.PayerState = mediatedtransfer.StatePayerWaitingWithdraw
			withdraw := &mediatedtransfer.EventContractSendWithdraw{
				Transfer:       pair.PayerTransfer,
				ChannelAddress: pair.PayerRoute.ChannelAddress,
			}
			events = append(events, withdraw)
		}
	}
	return events
}

/*
Set the state of the `payee_address` transfer, check the secret is
    being revealed backwards, and if necessary send out RevealSecret,
    SendBalanceProof, and Withdraws.
*/
func secretLearned(state *mediatedtransfer.MediatorState, secret common.Hash, payeeAddress common.Address, newPayeeState string) *transfer.TransitionResult {
	if !stateSecretKnownMaps[newPayeeState] {
		panic(fmt.Sprintf("%s not in STATE_SECRET_KNOWN", newPayeeState))
	}
	// TODO: if any of the transfers is in expired state, event for byzantine
	if state.Secret == utils.EmptyHash {
		state.SetSecret(secret)
	}
	var events []transfer.Event
	eventsWrongOrder := setPayeeStateAndCheckRevealOrder(state.TransfersPair, payeeAddress, newPayeeState)
	eventsSecretReveal := eventsForRevealSecret(state.TransfersPair, state.OurAddress)
	eventBalanceProof := eventsForBalanceProof(state.TransfersPair, state.BlockNumber)
	eventsWithdraw := eventsForWithdraw(state.TransfersPair, state.Db)
	events = append(events, eventsWrongOrder...)
	events = append(events, eventsSecretReveal...)
	events = append(events, eventBalanceProof...)
	events = append(events, eventsWithdraw...)
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

/*
update all routes state from db
*/
func updateAvaiableRoutesFromDb(db channel.Db, routes *transfer.RoutesState) {
	for _, r := range routes.AvailableRoutes {
		ch, err := db.GetChannelByAddress(r.ChannelAddress)
		if err != nil {
			log.Error(fmt.Sprintf("get channel %s from db err %s", r.ChannelAddress, err))
		} else {
			r.ClosedBlock = ch.ClosedBlock
			r.AvaibleBalance = ch.OurBalance.Sub(ch.OurBalance, ch.OurAmountLocked)
			r.State = ch.State
		}
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
func mediateTransfer(state *mediatedtransfer.MediatorState, payerRoute *transfer.RouteState, payerTransfer *mediatedtransfer.LockedTransferState) *transfer.TransitionResult {
	var transferPair *mediatedtransfer.MediationPairState
	var events []transfer.Event
	if params.TreatRefundTransferAsNormalMediatedTransfer {
		if state.HasRefunded {
			panic("has rufuned mediated node should never receive another transfer.")
		}
	}
	if state.Db != nil {
		updateAvaiableRoutesFromDb(state.Db, state.Routes)
	} else {
		log.Error(" db is nil can only be ignored when you are run testing...")
	}
	timeoutBlocks := int(getTimeoutBlocks(payerRoute, payerTransfer, state.BlockNumber))
	if timeoutBlocks > 0 {
		transferPair, events = nextTransferPair(payerRoute, payerTransfer, state.Routes, timeoutBlocks, state.BlockNumber)
	}
	if transferPair == nil {
		originalTransfer := payerTransfer
		originalRoute := payerRoute
		if len(state.TransfersPair) > 0 {
			//只有0不是refundtransfer，其他都已经是refundtransfer了，必须跳过比如A-B-C-B-E-F 这时E收到F的refund，E必须再次refund，但是refund的对象显然应该是
			originalTransfer = state.TransfersPair[0].PayerTransfer
			originalRoute = state.TransfersPair[0].PayerRoute
		}
		refundEvents := eventsForRefundTransfer(originalRoute, originalTransfer, timeoutBlocks, state.BlockNumber)
		if len(refundEvents) > 0 {
			if params.TreatRefundTransferAsNormalMediatedTransfer {
				rftr := refundEvents[0].(*mediatedtransfer.EventSendRefundTransfer)
				payeeLockedTransfer := &mediatedtransfer.LockedTransferState{
					Identifier:   rftr.Identifier,
					TargetAmount: rftr.TargetAmount,
					Amount:       new(big.Int).Set(rftr.Amount),
					Token:        rftr.Token,
					Initiator:    rftr.Initiator,
					Target:       rftr.Target,
					Expiration:   rftr.Expiration,
					Hashlock:     rftr.HashLock,
					Secret:       payerTransfer.Secret,
					Fee:          rftr.Fee, //refund transfer shouldnot charge.
				}
				payeeRoute := *originalRoute //route 信息是错误的，但是不影响，只要不继续route。
				transferPair = mediatedtransfer.NewMediationPairState(payerRoute, &payeeRoute, payerTransfer, payeeLockedTransfer)
				state.TransfersPair = append(state.TransfersPair, transferPair)
				state.HasRefunded = true
			}
		}
		return &transfer.TransitionResult{
			NewState: state,
			Events:   refundEvents,
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
After Raiden learns about a new block this function must be called to
    handle expiration of the hash time locks.

    Args:
        state (MediatorState): The current state.

    Return:
        TransitionResult: The resulting iteration
*/
func handleBlock(state *mediatedtransfer.MediatorState, st *transfer.BlockStateChange) *transfer.TransitionResult {
	blockNumber := state.BlockNumber
	if blockNumber < st.BlockNumber {
		blockNumber = st.BlockNumber
	}
	state.BlockNumber = blockNumber
	closeEvents := eventsForClose(state.TransfersPair, blockNumber)
	withdrawEvents := eventsForWithdraw(state.TransfersPair, state.Db)
	unlockfailEvents := setExpiredPairs(state.TransfersPair, blockNumber)
	var events []transfer.Event
	events = append(events, closeEvents...)
	events = append(events, withdrawEvents...)
	events = append(events, unlockfailEvents...)
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

/*
Validate and handle a ReceiveTransferRefund state change.

    A node might participate in mediated transfer more than once because of
    refund transfers, eg. A-B-C-B-D-T, B tried to mediate the transfer through
    C, which didn't have an available route to proceed and refunds B, at this
    point B is part of the path again and will try a new partner to proceed
    with the mediation through D, D finally reaches the target T.

    In the above scenario B has two pairs of payer and payee transfers:

        payer:A payee:C from the first SendMediatedTransfer
        payer:C payee:D from the following SendRefundTransfer

    Args:
        state (MediatorState): Current state.
        state_change (ReceiveTransferRefund): The state change.

    Returns:
        TransitionResult: The resulting iteration.
*/
func handleRefundTransfer(state *mediatedtransfer.MediatorState, st *mediatedtransfer.ReceiveTransferRefundStateChange) *transfer.TransitionResult {
	if state.Secret != utils.EmptyHash {
		panic("refunds are not allowed if the secret is revealed")
	}
	/*
			  The last sent transfer is the only one thay may be refunded, all the
		     previous ones are refunded already.
	*/
	l := len(state.TransfersPair)
	transferPair := state.TransfersPair[l-1]
	payeeTransfer := transferPair.PayeeTransfer
	if IsValidRefund(payeeTransfer, st.Transfer, st.Sender) {
		//什么时候会产生多个transferpair，就是发生refund的时候。 从这里也可以看出把refund当成普通mediatedtransfer来处理。
		return mediateTransfer(state, transferPair.PayeeRoute, st.Transfer)
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
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
	if utils.Sha3(secret[:]) == state.Hashlock {
		return secretLearned(state, secret, st.Sender, mediatedtransfer.StatePayeeSecretRevealed)
	}
	//  TODO: event for byzantine behavior,所有的reveal secret，是否与自己有关，都会到这里，是正常现象。
	return &transfer.TransitionResult{
		NewState: state,
		Events:   nil,
	}
}

//Handle a NettingChannelUnlock state change.
func handleContractWithDraw(state *mediatedtransfer.MediatorState, st *mediatedtransfer.ContractReceiveWithdrawStateChange) *transfer.TransitionResult {
	if utils.Sha3(state.Secret[:]) != state.Hashlock {
		panic("secret must be validated by the smart contract")
	}
	/*
			  For all but the last pair in transfer pair a refund transfer ocurred,
		     meaning the same channel was used twice, once when this node sent the
		     mediated transfer and once when the refund transfer was received. A
		     ContractReceiveWithdraw state change may be used for each.
	*/
	var events []transfer.Event
	//This node withdraw the refund
	if st.Receiver == state.OurAddress {
		for pos, pair := range state.TransfersPair {
			if pair.PayerRoute.ChannelAddress == st.ChannelAddress {
				/*
									  always set the contract_withdraw regardless of the previous
					                 state (even expired)
				*/
				pair.PayerState = mediatedtransfer.StatePayerContractWithdraw
				withdraw := &mediatedtransfer.EventWithdrawSuccess{
					Identifier: pair.PayerTransfer.Identifier,
					Hashlock:   pair.PayerTransfer.Hashlock,
				}
				events = append(events, withdraw)
				/*
									   if the current pair is backed by a refund set the sent
					                 mediated transfer to a 'secret known' state
				*/
				if pos > 0 {
					previousPair := state.TransfersPair[pos-1]
					if !stateTransferFinalMaps[previousPair.PayeeState] {
						previousPair.PayeeState = mediatedtransfer.StatePayeeRefundWithdraw
					}
				}
			}
		}
	} else {
		//A partner withdrew the mediated transfer
		for _, pair := range state.TransfersPair {
			if pair.PayerRoute.ChannelAddress == st.ChannelAddress {
				unlock := &mediatedtransfer.EventUnlockSuccess{
					Identifier: pair.PayeeTransfer.Identifier,
					Hashlock:   pair.PayeeTransfer.Hashlock,
				}
				events = append(events, unlock)
				pair.PayeeState = mediatedtransfer.StatePayeeContractWithdraw
			}
		}
	}
	tr := secretLearned(state, st.Secret, st.Receiver, mediatedtransfer.StatePayeeContractWithdraw)
	tr.Events = append(tr.Events, events...)
	return tr
}

// Handle a ReceiveBalanceProof state change.
func handleBalanceProof(state *mediatedtransfer.MediatorState, st *mediatedtransfer.ReceiveBalanceProofStateChange) *transfer.TransitionResult {
	var events []transfer.Event
	for _, pair := range state.TransfersPair {
		if pair.PayerRoute.HopNode == st.NodeAddress {
			withdraw := &mediatedtransfer.EventWithdrawSuccess{
				Identifier: pair.PayeeTransfer.Identifier,
				Hashlock:   pair.PayeeTransfer.Hashlock,
			}
			events = append(events, withdraw)
			pair.PayerState = mediatedtransfer.StatePayerBalanceProof
		}
	}
	return &transfer.TransitionResult{
		NewState: state,
		Events:   events,
	}
}

//Handle an ActionRouteChange state change
func handleRouteChange(state *mediatedtransfer.MediatorState, st *transfer.ActionRouteChangeStateChange) *transfer.TransitionResult {
	/*
	   TODO: `update_route` only changes the RoutesState, instead of moving the
	      routes to the MediationPairState use identifier to reference the routes
	*/
	newRoute := st.Route
	used := false
	/*
	   a route in use might be closed because of another task, update the pair
	      state in-place
	*/
	for _, pair := range state.TransfersPair {
		if pair.PayeeRoute.HopNode == newRoute.HopNode {
			pair.PayeeRoute = newRoute
			used = true
		}
		if pair.PayerRoute.HopNode == newRoute.HopNode {
			pair.PayerRoute = newRoute
			used = true
		}
	}
	if !used {
		mediatedtransfer.UpdateRoute(state.Routes, st)
	}
	//  a route might be closed by another task
	withdrawEvents := eventsForWithdraw(state.TransfersPair, state.Db)
	return &transfer.TransitionResult{
		NewState: state,
		Events:   withdrawEvents,
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
	//update all route state if possible
	if state != nil {
		if state.Db != nil {
			for _, pair := range state.TransfersPair {
				r := pair.PayerRoute
				ch, err := state.Db.GetChannelByAddress(r.ChannelAddress)
				if err != nil {
					log.Error(fmt.Sprintf("get channel %s from db err %s", utils.APex(r.ChannelAddress), err))
				} else {
					r.State = ch.State
					r.AvaibleBalance = ch.OurBalance.Sub(ch.OurBalance, ch.OurAmountLocked)
					r.ClosedBlock = ch.ClosedBlock
				}
				r = pair.PayeeRoute
				ch, err = state.Db.GetChannelByAddress(r.ChannelAddress)
				if err != nil {
					log.Error(fmt.Sprintf("get channel %s from db err %s", utils.APex(r.ChannelAddress), err))
				} else {
					r.State = ch.State
					r.AvaibleBalance = ch.OurBalance.Sub(ch.OurBalance, ch.OurAmountLocked)
					r.ClosedBlock = ch.ClosedBlock
				}
			}
		} else {
			log.Error(" db is nil can only be ignored when you are run testing...")
		}
	}
	if state == nil { //nil 判断小心了,有可能是空指针被当做非空了.
		if aim, ok := stateChange.(*mediatedtransfer.ActionInitMediatorStateChange); ok {
			state = &mediatedtransfer.MediatorState{
				OurAddress:  aim.OurAddress,
				Routes:      aim.Routes,
				BlockNumber: aim.BlockNumber,
				Hashlock:    aim.FromTranfer.Hashlock,
				Db:          aim.Db,
			}
			it = mediateTransfer(state, aim.FromRoute, aim.FromTranfer)
		}
	} else if state.Secret == utils.EmptyHash {
		switch st2 := stateChange.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
			//目前没用
		case *transfer.ActionRouteChangeStateChange:
			it = handleRouteChange(state, st2)
		case *mediatedtransfer.ReceiveTransferRefundStateChange:
			it = handleRefundTransfer(state, st2)
		case *mediatedtransfer.ReceiveSecretRevealStateChange:
			it = handleSecretReveal(state, st2)
		case *mediatedtransfer.ContractReceiveWithdrawStateChange:
			it = handleContractWithDraw(state, st2)
		default:
			log.Info(fmt.Sprintf("unknown statechange :%s", utils.StringInterface(st2, 3)))
		}
	} else {
		switch st2 := stateChange.(type) {
		case *transfer.BlockStateChange:
			it = handleBlock(state, st2)
			//目前没用
		case *transfer.ActionRouteChangeStateChange:
			it = handleRouteChange(state, st2)
		case *mediatedtransfer.ReceiveSecretRevealStateChange:
			it = handleSecretReveal(state, st2)
		case *mediatedtransfer.ReceiveBalanceProofStateChange:
			it = handleBalanceProof(state, st2)
		case *mediatedtransfer.ContractReceiveWithdrawStateChange:
			it = handleContractWithDraw(state, st2)
		default:
			log.Info(fmt.Sprintf("unknown statechange :%s", utils.StringInterface(st2, 3)))
		}
	}
	// this is the place for paranoia
	if it.NewState != nil {
		//iteration.new_state todo check this is equal
		sanityCheck(it.NewState.(*mediatedtransfer.MediatorState))
	}
	return clearIfFinalized(it)
}
