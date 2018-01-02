package mediator

import (
	"testing"

	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/SmartMeshFoundation/raiden-network/utils/utest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	assert2 "github.com/stretchr/testify/assert"
)

func assert(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert2.EqualValues(t, expected, actual, msgAndArgs...)
}
func TestTypConversion(t *testing.T) {
	var originalState transfer.State
	_, ok := originalState.(*mediated_transfer.MediatorState)
	if ok { //for a real nil pointer, it cannot convert
		t.Error("MediatorState StateTransition get type error ")
	}
}

func makeInitStateChange(fromTransfer *mediated_transfer.LockedTransferState, fromRoute *transfer.RouteState, routes []*transfer.RouteState, ourAddress common.Address) *mediated_transfer.ActionInitMediatorStateChange {
	return &mediated_transfer.ActionInitMediatorStateChange{
		OurAddress:  ourAddress,
		FromTranfer: fromTransfer,
		Routes:      transfer.NewRoutesState(routes),
		FromRoute:   fromRoute,
		BlockNumber: 1,
	}
}
func makeTransferPair(payer, payee, initiator, target common.Address, amount int64, expiration int64, secret common.Hash, revealTimeout int) *mediated_transfer.MediationPairState {
	payerExpiration := expiration
	payeeExpiration := expiration - int64(revealTimeout)
	return mediated_transfer.NewMediationPairState(
		utest.MakeRoute(payer, amount, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utest.ADDR),
		utest.MakeRoute(payee, amount, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utest.ADDR),
		utest.MakeTransfer(amount, initiator, target, payerExpiration, secret, utils.Sha3(secret[:]), utest.UNIT_IDENTIFIER, utest.UNIT_TOKEN_ADDRESS),
		utest.MakeTransfer(amount, initiator, target, payeeExpiration, secret, utils.Sha3(secret[:]), utest.UNIT_IDENTIFIER, utest.UNIT_TOKEN_ADDRESS),
	)
}
func makeTransfersPair(initiator common.Address, hops []common.Address, target common.Address, amount int64, secret common.Hash, initialExpiration int64, revealTimeout int) []*mediated_transfer.MediationPairState {
	if initialExpiration == 0 {
		initialExpiration = int64((2*len(hops) + 1) * revealTimeout)
	}
	var pairs []*mediated_transfer.MediationPairState
	nextExpiration := initialExpiration
	for i := 0; i < len(hops)-1; i++ {
		payer := hops[i]
		payee := hops[i+1]
		if nextExpiration <= 0 {
			log.Error("nextExpiration<=0")
		}
		pair := makeTransferPair(payer, payee, initiator, target, amount, nextExpiration, secret, revealTimeout)
		pairs = append(pairs, pair)
		/*
					    assumes that the node sending the refund will follow the protocol and
			         decrement the expiration for it's lock
		*/
		nextExpiration = pair.PayeeTransfer.Expiration - int64(revealTimeout)
	}
	return pairs
}
func makeMediatorState(fromTransfer *mediated_transfer.LockedTransferState, fromRoute *transfer.RouteState, routes []*transfer.RouteState, ourAddress common.Address) *mediated_transfer.MediatorState {
	statechange := makeInitStateChange(fromTransfer, fromRoute, routes, ourAddress)
	it := StateTransition(nil, statechange)
	return it.NewState.(*mediated_transfer.MediatorState)
}

//A hash time lock is valid up to the expiraiton block.
func TestIsLockValid(t *testing.T) {
	var amount int64 = 10
	var expiration int64 = 10
	initiator := utest.HOP1
	target := utest.HOP2
	tr := utest.MakeTransfer(amount, initiator, target, expiration, utils.EmptyHash, utils.EmptyHash, utest.UNIT_IDENTIFIER, utest.UNIT_TOKEN_ADDRESS)
	assert(t, IsLockValid(tr, 5), true)
	assert(t, IsLockValid(tr, 10), true)
	assert(t, IsLockValid(tr, 11), false)

}

/*
It's safe to wait for a secret while there are more than reveal timeout
    blocks until the lock expiration.
*/
func TestIsSafeToWait(t *testing.T) {
	var amount int64 = 10
	var expiration int64 = 40
	initiator := utest.HOP1
	target := utest.HOP2
	tr := utest.MakeTransfer(amount, initiator, target, expiration, utils.EmptyHash, utils.EmptyHash, utest.UNIT_IDENTIFIER, utest.UNIT_TOKEN_ADDRESS)
	//expiration is in 30 blocks, 19 blocks safe for waiting
	assert(t, IsSafeToWait(tr, 10, 1), true)
	assert(t, IsSafeToWait(tr, 10, 20), true)
	assert(t, IsSafeToWait(tr, 10, 29), true)

	assert(t, IsSafeToWait(tr, 10, 30), false)
	assert(t, IsSafeToWait(tr, 10, 50), false)

}

//Don't close the channel if the payee transfer is not paid.
func TestIsChannelCloseNeededUnpaid(t *testing.T) {
	var amount int64 = 10
	var expiration int64 = 10
	var revealTimeout int = 5

	/*
			    even if the secret is known by the payee, the transfer is paid only if a
		    withdraw on-chain happened or if the mediator has sent a balance proof
	*/
	unpaidStates := []string{mediated_transfer.STATE_PAYEE_PENDING, mediated_transfer.STATE_PAYEE_SECRET_REVEALED, mediated_transfer.STATE_PAYEE_REFUND_WITHDRAW}
	for _, unpaidState := range unpaidStates {
		unpaidPair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
		unpaidPair.PayeeState = unpaidState
		assert(t, unpaidPair.PayerRoute.State, transfer.CHANNEL_STATE_OPENED)
		safeBlock := expiration - int64(revealTimeout) - 1
		assert(t, isChannelCloseNeeded(unpaidPair, safeBlock), false)
		unsafeBlock := expiration - int64(revealTimeout)
		assert(t, isChannelCloseNeeded(unpaidPair, unsafeBlock), false)
	}
}

//Close the channel if the payee transfer is paid but the payer has not paid.
func TestIsChannelClosedNeededPaid(t *testing.T) {
	var amount int64 = 10
	var expiration int64 = 10
	var revealTimeout int = 5

	paidStates := []string{mediated_transfer.STATE_PAYEE_CONTRACT_WITHDRAW, mediated_transfer.STATE_PAYEE_BALANCE_PROOF}
	for _, paidState := range paidStates {
		paidPair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
		paidPair.PayeeState = paidState
		assert(t, paidPair.PayerRoute.State, transfer.CHANNEL_STATE_OPENED)
		safeBlock := expiration - int64(revealTimeout) - 1
		assert(t, isChannelCloseNeeded(paidPair, safeBlock), false)
		unsafeBlock := expiration - int64(revealTimeout)
		assert(t, isChannelCloseNeeded(paidPair, unsafeBlock), true)
	}
}

//If the channel is already closed the anser is always no.
func TestIsChannelCloseNeedChannelClosed(t *testing.T) {
	var amount int64 = 10
	var expiration int64 = 10
	var revealTimeout int = 5

	for state, _ := range mediated_transfer.ValidPayeeStateMap {
		pair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
		pair.PayeeState = state
		pair.PayerRoute.State = transfer.CHANNEL_STATE_CLOSED

		safeBlock := expiration - int64(revealTimeout) - 1
		assert(t, isChannelCloseNeeded(pair, safeBlock), false)
		unsafeBlock := expiration - int64(revealTimeout)
		assert(t, isChannelCloseNeeded(pair, unsafeBlock), false)
	}
}

func TestIsChannelCloseNeededClosed(t *testing.T) {
	var amount int64 = 10
	var expiration int64 = 10
	var revealTimeout int = 5

	pair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
	pair.PayeeState = mediated_transfer.STATE_PAYEE_BALANCE_PROOF

	assert(t, pair.PayerRoute.State, transfer.CHANNEL_STATE_OPENED)
	safeBlock := expiration - int64(revealTimeout) - 1
	assert(t, isChannelCloseNeeded(pair, safeBlock), false)
	unsafeBlock := expiration - int64(revealTimeout)
	assert(t, isChannelCloseNeeded(pair, unsafeBlock), true)

}
func TestIsValidRefund(t *testing.T) {
	initiator := utest.ADDR
	target := utest.HOP1
	validSender := utest.HOP2
	tr := &mediated_transfer.LockedTransferState{
		Identifier: 20,
		Amount:     30,
		Token:      utest.UNIT_TOKEN_ADDRESS,
		Initiator:  initiator,
		Target:     target,
		Expiration: 50,
		Hashlock:   utest.UNIT_HASHLOCK,
		Secret:     utils.EmptyHash,
	}
	refundLowerExpiration := &mediated_transfer.LockedTransferState{
		Identifier: 20,
		Amount:     30,
		Token:      utest.UNIT_TOKEN_ADDRESS,
		Initiator:  initiator,
		Target:     target,
		Expiration: 35,
		Hashlock:   utest.UNIT_HASHLOCK,
		Secret:     utils.EmptyHash,
	}
	assert(t, IsValidRefund(tr, refundLowerExpiration, validSender), true)
	assert(t, IsValidRefund(tr, refundLowerExpiration, target), false)

	refundSameExpiration := &mediated_transfer.LockedTransferState{
		Identifier: 20,
		Amount:     30,
		Token:      utest.UNIT_TOKEN_ADDRESS,
		Initiator:  initiator,
		Target:     target,
		Expiration: 50,
		Hashlock:   utest.UNIT_HASHLOCK,
		Secret:     utils.EmptyHash,
	}
	assert(t, IsValidRefund(tr, refundSameExpiration, validSender), false)
}

func TestGetTimeoutBlocks(t *testing.T) {
	var amount int64 = 10
	initiator := utest.HOP1
	nextHop := utest.HOP2
	settleTimeout := 30
	var blockNumber int64 = 5
	route := utest.MakeRoute(nextHop, amount, settleTimeout, 0, 0, utils.NewRandomAddress())
	var earlyExpire int64 = 10
	earlyTransfer := utest.MakeTransfer(amount, initiator, nextHop, earlyExpire, utils.EmptyHash, utils.EmptyHash, 0, utest.UNIT_TOKEN_ADDRESS)
	earlyBlock := getTimeoutBlocks(route, earlyTransfer, blockNumber)
	assert(t, earlyBlock, 5-TRANSIT_BLOCKS)

	var equalExpire int64 = 30
	equalTransfer := utest.MakeTransfer(amount, initiator, nextHop, equalExpire, utils.EmptyHash, utils.EmptyHash, 0, utest.UNIT_TOKEN_ADDRESS)
	equalBlock := getTimeoutBlocks(route, equalTransfer, blockNumber)
	assert(t, equalBlock, 25-TRANSIT_BLOCKS)

	var largeExpire int64 = 70
	largeTransfer := utest.MakeTransfer(amount, initiator, nextHop, largeExpire, utils.EmptyHash, utils.EmptyHash, 0, utest.UNIT_TOKEN_ADDRESS)
	largeBlock := getTimeoutBlocks(route, largeTransfer, blockNumber)
	assert(t, largeBlock, 30-TRANSIT_BLOCKS)

	closedRoute := utest.MakeRoute(nextHop, amount, settleTimeout, 0, 2, utils.NewRandomAddress())

	largeBlock = getTimeoutBlocks(closedRoute, largeTransfer, blockNumber)
	assert(t, largeBlock, 27-TRANSIT_BLOCKS)
	/*
		the computed timeout may be negative, in which case the calling code must /not/ use it
	*/
	negativeBlockNumber := largeExpire
	negativeBlock := getTimeoutBlocks(route, largeTransfer, negativeBlockNumber)
	assert(t, negativeBlock, -TRANSIT_BLOCKS)

}

//Routes that dont have enough available_balance must be ignored.
func TestNextRouteAmount(t *testing.T) {
	var amount int64 = 10
	revealTimeout := 30
	timeoutBlocks := 40
	routes := []*transfer.RouteState{
		utest.MakeRoute(utest.HOP2, amount*2, 0, revealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP1, amount+1, 0, revealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP3, amount/2, 0, revealTimeout, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP4, amount, 0, revealTimeout, 0, utils.NewRandomAddress()),
	}
	routesState := transfer.NewRoutesState(routes)
	route1 := nextRoute(routesState, timeoutBlocks, amount)
	assert(t, route1, routes[0])
	assert(t, routesState.AvailableRoutes, routes[1:])
	assert(t, len(routesState.IgnoredRoutes), 0)

	route2 := nextRoute(routesState, timeoutBlocks, amount)
	assert(t, route2, routes[1])
	assert(t, routesState.AvailableRoutes, routes[2:])
	assert(t, len(routesState.IgnoredRoutes), 0)

	route3 := nextRoute(routesState, timeoutBlocks, amount)
	assert(t, route3, routes[3])
	assert(t, len(routesState.AvailableRoutes), 0)
	assert(t, routesState.IgnoredRoutes, []*transfer.RouteState{routes[2]})

	assert(t, nextRoute(routesState, timeoutBlocks, amount) == nil, true)

}

//Routes with a larger reveal timeout than timeout_blocks must be ignored.
func TestNextRouteRevealTimeout(t *testing.T) {
	var amount int64 = 10
	var balance int64 = 20
	timeoutBlocks := 40

	routes := []*transfer.RouteState{
		utest.MakeRoute(utest.HOP2, balance, 0, timeoutBlocks*2, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP1, balance, 0, timeoutBlocks+1, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP3, balance, 0, timeoutBlocks/2, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP4, balance, 0, timeoutBlocks, 0, utils.NewRandomAddress()),
	}

	routesState := transfer.NewRoutesState(routes)
	route1 := nextRoute(routesState, timeoutBlocks, amount)
	assert(t, route1, routes[2])
	assert(t, routesState.AvailableRoutes, routes[3:])
	assert(t, routesState.IgnoredRoutes, routes[0:2])

	assert(t, nextRoute(routesState, timeoutBlocks, amount) == nil, true)
	assert(t, len(routesState.AvailableRoutes), 0)
	assert(t, routesState.IgnoredRoutes, append(routes[0:2], routes[3]))
}

func TestNextTransferPair(t *testing.T) {
	timeoutBlocks := 47
	var blockNumber int64 = 3
	var balance int64 = 10
	initiator := utest.HOP1
	target := utest.ADDR

	payerRoute := utest.MakeRoute(initiator, balance, 0, 0, 0, utils.NewRandomAddress())
	payerTransfer := utest.MakeTransfer(balance, initiator, target, 50, utils.EmptyHash, utils.EmptyHash, 1, utest.UNIT_TOKEN_ADDRESS)

	routes := []*transfer.RouteState{utest.MakeRoute(utest.HOP2, balance, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress())}
	routesState := transfer.NewRoutesState(routes)
	pair, events := nextTransferPair(payerRoute, payerTransfer, routesState, timeoutBlocks, blockNumber)

	assert(t, pair.PayerRoute, payerRoute)
	assert(t, pair.PayerTransfer, payerTransfer)
	assert(t, pair.PayeeRoute, routes[0])
	assert(t, pair.PayeeTransfer.Expiration < pair.PayerTransfer.Expiration, true)

	assert(t, len(events), 1)
	tr, ok := events[0].(*mediated_transfer.EventSendMediatedTransfer)
	assert(t, ok, true)
	assert(t, tr.Identifier, payerTransfer.Identifier)
	assert(t, tr.Token, payerTransfer.Token)
	assert(t, tr.Amount, payerTransfer.Amount)
	assert(t, tr.HashLock, payerTransfer.Hashlock)
	assert(t, tr.Initiator, payerTransfer.Initiator)
	assert(t, tr.Target, payerTransfer.Target)
	assert(t, tr.Expiration < payerTransfer.Expiration, true)
	assert(t, tr.Receiver, pair.PayeeRoute.HopNode)

	assert(t, len(routesState.AvailableRoutes), 0)
}

func TestSetPayee(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	assert(t, pairs[0].PayerState, mediated_transfer.STATE_PAYER_PENDING)
	assert(t, pairs[0].PayeeState, mediated_transfer.STATE_PAYEE_PENDING)
	assert(t, pairs[1].PayerState, mediated_transfer.STATE_PAYER_PENDING)
	assert(t, pairs[1].PayeeState, mediated_transfer.STATE_PAYEE_PENDING)
	SetPayeeStateAndCheckRevealOrder(pairs, utest.HOP2, mediated_transfer.STATE_PAYEE_SECRET_REVEALED)
	assert(t, pairs[0].PayerState, mediated_transfer.STATE_PAYER_PENDING)
	assert(t, pairs[0].PayeeState, mediated_transfer.STATE_PAYEE_PENDING)
	assert(t, pairs[1].PayerState, mediated_transfer.STATE_PAYER_PENDING)
	assert(t, pairs[1].PayeeState, mediated_transfer.STATE_PAYEE_PENDING)

	SetPayeeStateAndCheckRevealOrder(pairs, utest.HOP3, mediated_transfer.STATE_PAYEE_SECRET_REVEALED)
	assert(t, pairs[0].PayerState, mediated_transfer.STATE_PAYER_PENDING)
	assert(t, pairs[0].PayeeState, mediated_transfer.STATE_PAYEE_SECRET_REVEALED)
	assert(t, pairs[1].PayerState, mediated_transfer.STATE_PAYER_PENDING)
	assert(t, pairs[1].PayeeState, mediated_transfer.STATE_PAYEE_PENDING)

}

// The transfer pair must switch to expired at the right block.
func TestSetExpiredPairs(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	pair := pairs[0]
	// do not generate events if the secret is not known
	firstUnsafeBlock := pair.PayerTransfer.Expiration - int64(pair.PayerRoute.RevealTimeout)
	setExpiredPairs(pairs, firstUnsafeBlock)
	assert(t, pair.PayeeState, mediated_transfer.STATE_PAYEE_PENDING)
	assert(t, pair.PayerState, mediated_transfer.STATE_PAYER_PENDING)
	//payee lock expired
	payee_expiration_block := pair.PayeeTransfer.Expiration
	setExpiredPairs(pairs, payee_expiration_block)
	//dge case for the payee lock expiration
	assert(t, pair.PayeeState, mediated_transfer.STATE_PAYEE_PENDING)
	assert(t, pair.PayerState, mediated_transfer.STATE_PAYER_PENDING)
	setExpiredPairs(pairs, payee_expiration_block+1)
	assert(t, pair.PayeeState, mediated_transfer.STATE_PAYEE_EXPIRED)
	assert(t, pair.PayerState, mediated_transfer.STATE_PAYER_PENDING)
	// edge case for the payer lock expiration
	payer_expiration_block := pair.PayerTransfer.Expiration
	setExpiredPairs(pairs, payer_expiration_block)
	assert(t, pair.PayeeState, mediated_transfer.STATE_PAYEE_EXPIRED)
	assert(t, pair.PayerState, mediated_transfer.STATE_PAYER_PENDING)
	setExpiredPairs(pairs, payer_expiration_block+1)
	assert(t, pair.PayeeState, mediated_transfer.STATE_PAYEE_EXPIRED)
	assert(t, pair.PayerState, mediated_transfer.STATE_PAYER_EXPIRED)
}

func TestEventsForRefund(t *testing.T) {
	var amount int64 = 10
	var expiration int64 = 30
	var revealTimeout int = 17
	var timeoutBlocks int = int(expiration)
	var blockNumber int64 = 1
	initiator := utest.HOP1
	target := utest.HOP6
	refundRoute := utest.MakeRoute(initiator, amount, utest.UNIT_SETTLE_TIMEOUT, revealTimeout, 0, utils.NewRandomAddress())
	refundTransfer := utest.MakeTransfer(amount, initiator, target, expiration, utils.EmptyHash, utils.EmptyHash, 1, utest.UNIT_TOKEN_ADDRESS)
	smallTimeoutBlocks := revealTimeout
	smallRefundEvents := eventsForRefundTransfer(refundRoute, refundTransfer, smallTimeoutBlocks, blockNumber)
	assert(t, len(smallRefundEvents), 0)

	refundEvents := eventsForRefundTransfer(refundRoute, refundTransfer, timeoutBlocks, blockNumber)
	ev, ok := refundEvents[0].(*mediated_transfer.EventSendRefundTransfer)
	assert(t, ok, true)
	assert(t, ev.Expiration < blockNumber+int64(timeoutBlocks), true)
	assert(t, ev.Amount, amount)
	assert(t, ev.HashLock, refundTransfer.Hashlock)
}

/*
 The secret is revealed backwards to the payer once the payee sent the
    SecretReveal.
*/
func TestEventsForRevealSecret(t *testing.T) {
	secret := utest.UNIT_SECRET
	ourAddress := utest.ADDR
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	events := EventsForRevealSecret(pairs, ourAddress)
	/*
			   the secret is known by this node, but no other payee is at a secret known
		    state, do nothing
	*/
	assert(t, len(events), 0)
	firstPair := pairs[0]
	lastPair := pairs[1]

	lastPair.PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	events = EventsForRevealSecret(pairs, ourAddress)
	/*
			 # the last known hop sent a secret reveal message, this node learned the
		    # secret and now must reveal to the payer node from the transfer pair
	*/
	assert(t, len(events), 1)
	ev, ok := events[0].(*mediated_transfer.EventSendRevealSecret)
	assert(t, ok, true)
	assert(t, ev.Secret, secret)
	assert(t, ev.Receiver, lastPair.PayerRoute.HopNode)
	assert(t, lastPair.PayerState, mediated_transfer.STATE_PAYER_SECRET_REVEALED)

	events = EventsForRevealSecret(pairs, ourAddress)
	/*
			  # the payeee from the first_pair did not send a secret reveal message, do
		    # nothing
	*/
	assert(t, len(events), 0)
	firstPair.PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	events = EventsForRevealSecret(pairs, ourAddress)
	assert(t, len(events), 1)
	ev, ok = events[0].(*mediated_transfer.EventSendRevealSecret)
	assert(t, ok, true)
	assert(t, ev.Secret, secret)
	assert(t, ev.Receiver, firstPair.PayerRoute.HopNode)
	assert(t, firstPair.PayerState, mediated_transfer.STATE_PAYER_SECRET_REVEALED)
}

// When the secret is not know there is nothing to do.
func TestEventsForRevealSecretSecretUnkown(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utils.EmptyHash, 0, utest.UNIT_REVEAL_TIMEOUT)
	events := EventsForRevealSecret(pairs, utest.ADDR)
	//这个测试有什么意义呢?肯定是0啊,和上一个测试没有一点不一样.
	assert(t, len(events), 0)

}

// The secret must be revealed backwards to the payer if the payee knows
//the secret.
func TestEventsForRevealSecretAllStates(t *testing.T) {
	secret := utest.UNIT_SECRET
	ourAddress := utest.ADDR
	payee_secret_known := []string{mediated_transfer.STATE_PAYEE_SECRET_REVEALED, mediated_transfer.STATE_PAYEE_REFUND_WITHDRAW, mediated_transfer.STATE_PAYEE_CONTRACT_WITHDRAW, mediated_transfer.STATE_PAYEE_BALANCE_PROOF}
	for _, state := range payee_secret_known {
		pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, secret, 0, utest.UNIT_REVEAL_TIMEOUT)
		pair := pairs[0]
		pair.PayeeState = state
		events := EventsForRevealSecret(pairs, ourAddress)
		ev, ok := events[0].(*mediated_transfer.EventSendRevealSecret)
		assert(t, ok, true)
		assert(t, ev.Secret, secret)
		assert(t, ev.Receiver, utest.HOP2)
	}
}

/*
Test the simple case were the last hop has learned the secret and sent
    to the mediator node.
*/
func TestEventsForBanalceProof(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	lastpair := pairs[0]
	lastpair.PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	// the lock has not expired yet
	blockNumber := lastpair.PayeeTransfer.Expiration
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 2)
	var balanceProof *mediated_transfer.EventSendBalanceProof
	var unlockSuccess *mediated_transfer.EventUnlockSuccess

	for _, e := range events {
		switch e2 := e.(type) {
		case *mediated_transfer.EventSendBalanceProof:
			balanceProof = e2
		case *mediated_transfer.EventUnlockSuccess:
			unlockSuccess = e2
		}
	}
	assert(t, balanceProof.Receiver, lastpair.PayeeRoute.HopNode)
	assert(t, lastpair.PayeeState, mediated_transfer.STATE_PAYEE_BALANCE_PROOF)
	assert(t, unlockSuccess != nil, true)

}

/*
Balance proofs are useless if the channel is closed/settled, the payee
    needs to go on-chain and use the latest known balance proof which includes
    this lock in the locksroot.
*/
func TestEventsForBalanceProofChannelClosed(t *testing.T) {
	invalidStates := []string{transfer.CHANNEL_STATE_CLOSED, transfer.CHANNEL_STATE_SETTLED}
	for _, state := range invalidStates {
		pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
		var blockNumber int64 = 5
		lastpair := pairs[0]
		lastpair.PayeeRoute.State = state
		lastpair.PayeeRoute.ClosedBlock = blockNumber
		lastpair.PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
		events := eventsForBalanceProof(pairs, blockNumber)
		assert(t, len(events), 0)
	}
}

/*
Even though the secret should only propagate from the end of the chain
    to the front, if there is a payee node in the middle that knows the secret
    the Balance Proof is sent neverthless.

    This can be done safely because the secret is know to the mediator and
    there is reveal_timeout blocks to withdraw the lock on-chain with the payer.
*/
func TestEventsForBalanceProofMiddleSecret(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4, utest.HOP5}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	var blockNumber int64 = 1
	middlePair := pairs[1]
	middlePair.PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 2)
	var balanceProof *mediated_transfer.EventSendBalanceProof
	var unlockSuccess *mediated_transfer.EventUnlockSuccess

	for _, e := range events {
		switch e2 := e.(type) {
		case *mediated_transfer.EventSendBalanceProof:
			balanceProof = e2
		case *mediated_transfer.EventUnlockSuccess:
			unlockSuccess = e2
		}
	}
	assert(t, balanceProof.Receiver, middlePair.PayeeRoute.HopNode)
	assert(t, middlePair.PayeeState, mediated_transfer.STATE_PAYEE_BALANCE_PROOF)
	assert(t, unlockSuccess != nil, true)
}

//Nothing to do if the secret is not known.
func TestEventsFroBalanceProofSecretUnkown(t *testing.T) {
	var blockNumber int64 = 1
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utils.EmptyHash, 0, utest.UNIT_REVEAL_TIMEOUT)
	//pairs[1].PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 0)

	pairs = makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	/*
			# Even though the secret is set, there is not a single transfer pair with a
		    # 'secret known' state, so nothing should be done. This state is impossible
		    # to reach, in reality someone needs to reveal the secret to the mediator,
		    # so at least one other node knows the secret.
	*/
	events = eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 0)
}

//The balance proof should not be sent if the lock has expird.
func TestEventsForBalanceProofLockExpired(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4, utest.HOP5}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	lastpair := pairs[len(pairs)-1]
	lastpair.PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	var blockNumber int64 = lastpair.PayeeTransfer.Expiration + 1
	//the lock has expired, do not send a balance proof
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 0)
	middlePair := pairs[len(pairs)-2]
	middlePair.PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	/*
			    # Even though the last node did not receive the payment we should send the
		    # balance proof to the middle node to avoid unnecessarely closing the
		    # middle channel. This state should not be reached under normal operation,
		    # the last hop needs to choose a proper reveal_timeout and must go on-chain
		    # to withdraw the token before the lock expires.
	*/
	events = eventsForBalanceProof(pairs, blockNumber)
	var balanceProof *mediated_transfer.EventSendBalanceProof
	var unlockSuccess *mediated_transfer.EventUnlockSuccess
	for _, e := range events {
		switch e2 := e.(type) {
		case *mediated_transfer.EventSendBalanceProof:
			balanceProof = e2
		case *mediated_transfer.EventUnlockSuccess:
			unlockSuccess = e2
		}
	}
	assert(t, balanceProof.Receiver, middlePair.PayeeRoute.HopNode)
	assert(t, middlePair.PayeeState, mediated_transfer.STATE_PAYEE_BALANCE_PROOF)
	assert(t, unlockSuccess != nil, true)
}

// The node must close to unlock on-chain if the payee was paid.
func TestEventsForClose(t *testing.T) {
	var states = []string{mediated_transfer.STATE_PAYEE_BALANCE_PROOF, mediated_transfer.STATE_PAYEE_CONTRACT_WITHDRAW}
	for _, payeeState := range states {
		pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
		pair := pairs[0]
		pair.PayeeState = payeeState
		blockNumber := pair.PayerTransfer.Expiration - int64(pair.PayeeRoute.RevealTimeout)
		events := eventsForClose(pairs, blockNumber)
		assert(t, len(events), 1)
		ev := events[0].(*mediated_transfer.EventContractSendChannelClose)
		assert(t, pair.PayerState, mediated_transfer.STATE_PAYER_WAITING_CLOSE)
		assert(t, ev.ChannelAddress, pair.PayerRoute.ChannelAddress)
	}
}

/* 这个是怎么发生的呢?
If the secret is known but the payee transfer has not being paid the
   node must not settle on-chain, otherwise the payee can burn tokens to
   induce the mediator to close a channel.
*/
func TestEventsForCloseHoldForUnpaidPayee(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	pair := pairs[0]
	assert(t, pair.PayerTransfer.Secret, utest.UNIT_SECRET)
	assert(t, pair.PayeeTransfer.Secret, utest.UNIT_SECRET)
	assert(t, StateTransferPaidMap[pair.PayeeState], false)

	// do not generate events if the secret is known AND the payee is not paid
	firstUnSafeBlock := pair.PayerTransfer.Expiration - int64(pair.PayerRoute.RevealTimeout)
	events := eventsForClose(pairs, firstUnSafeBlock)
	assert(t, len(events), 0)
	assert(t, StateTransferPaidMap[pair.PayeeState], false)
	assert(t, StateTransferPaidMap[pair.PayerState], false)

	payerExpirationBlock := pair.PayerTransfer.Expiration
	events = eventsForClose(pairs, payerExpirationBlock)
	assert(t, len(events), 0)
	assert(t, StateTransferPaidMap[pair.PayeeState], false)
	assert(t, StateTransferPaidMap[pair.PayerState], false)

}

//The withdraw is done regardless of the current block.
func TestEventsForWithdrawChannelClosed(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	pair := pairs[0]
	pair.PayerRoute.State = transfer.CHANNEL_STATE_CLOSED

	// that's why this function doesn't receive the block_number
	events := eventsForWithdraw(pairs)
	assert(t, len(events), 1)
	ev, ok := events[0].(*mediated_transfer.EventContractSendWithdraw)
	assert(t, ok, true)
	assert(t, ev.ChannelAddress, pair.PayerRoute.ChannelAddress)
}
func TestEventsForWithdrawChannelOpen(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UNIT_SECRET, 0, utest.UNIT_REVEAL_TIMEOUT)
	// that's why this function doesn't receive the block_number
	events := eventsForWithdraw(pairs)
	assert(t, len(events), 0)
}

func TestSecretLearned(t *testing.T) {
	fromRoute, fromTransfer := utest.MakeFrom(utest.UNIT_TRANSFER_AMOUNT, utest.HOP5, int64(utest.HOP1_TIMEOUT), utils.NewRandomAddress(), utils.EmptyHash)
	routes := []*transfer.RouteState{
		utest.MakeRoute(utest.HOP2, utest.UNIT_TRANSFER_AMOUNT, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress()),
	}
	state := makeMediatorState(fromTransfer, fromRoute, routes, utils.NewRandomAddress())
	secret := utest.UNIT_SECRET
	payeeAddress := utest.HOP2
	it := secretLearned(state, secret, payeeAddress, mediated_transfer.STATE_PAYEE_SECRET_REVEALED)
	mstate := it.NewState.(*mediated_transfer.MediatorState)
	transferpair := mstate.TransfersPair[0]
	assert(t, fromTransfer.Expiration > transferpair.PayeeTransfer.Expiration, true)
	assert(t, transferpair.PayeeTransfer.AlmostEqual(fromTransfer), true)
	assert(t, transferpair.PayeeRoute, routes[0])

	assert(t, transferpair.PayerRoute, fromRoute)
	assert(t, transferpair.PayerTransfer, fromTransfer)
	assert(t, mstate.Secret, secret)
	assert(t, transferpair.PayeeTransfer.Secret, secret)
	assert(t, transferpair.PayerTransfer.Secret, secret)
	assert(t, transferpair.PayeeState, mediated_transfer.STATE_PAYEE_BALANCE_PROOF)
	assert(t, transferpair.PayerState, mediated_transfer.STATE_PAYER_SECRET_REVEALED)
	revealNumber := 0
	balanceNumber := 0
	for _, e := range it.Events {
		switch e.(type) {
		case *mediated_transfer.EventSendRevealSecret:
			revealNumber++
		case *mediated_transfer.EventSendBalanceProof:
			balanceNumber++
		}
	}
	assert(t, revealNumber, 1)
	assert(t, balanceNumber, 1)
}
func TestMediateTransfer(t *testing.T) {
	var amount int64 = 10
	var blockNumber int64 = 5
	var expiration int64 = 30
	var routes = []*transfer.RouteState{utest.MakeRoute(utest.HOP2, utest.UNIT_TRANSFER_AMOUNT, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress())}
	routesState := transfer.NewRoutesState(routes)
	state := &mediated_transfer.MediatorState{
		OurAddress:  utest.ADDR,
		Routes:      routesState,
		BlockNumber: blockNumber,
		Hashlock:    utest.UNIT_HASHLOCK,
	}
	payerroute, payertransfer := utest.MakeFrom(amount, utest.HOP6, expiration, utils.NewRandomAddress(), utils.EmptyHash)
	it := mediateTransfer(state, payerroute, payertransfer)
	var eventsMediated []*mediated_transfer.EventSendMediatedTransfer
	for _, e := range it.Events {
		switch e2 := e.(type) {
		case *mediated_transfer.EventSendMediatedTransfer:
			eventsMediated = append(eventsMediated, e2)
		}
	}
	assert(t, len(eventsMediated), 1)
	tr := eventsMediated[0]
	assert(t, tr.Identifier, payertransfer.Identifier)
	assert(t, tr.Token, payertransfer.Token)
	assert(t, tr.Amount, payertransfer.Amount)
	assert(t, tr.HashLock, payertransfer.Hashlock)
	assert(t, tr.Target, payertransfer.Target)
	assert(t, payertransfer.Expiration > tr.Expiration, true)
	assert(t, tr.Receiver, routes[0].HopNode)

}

func TestInitMediator(t *testing.T) {
	fromRoute, FromTransfer := utest.MakeFrom(utest.UNIT_TRANSFER_AMOUNT, utest.HOP2, int64(utest.HOP1_TIMEOUT), utils.NewRandomAddress(), utils.EmptyHash)
	var routes = []*transfer.RouteState{utest.MakeRoute(utest.HOP2, utest.UNIT_TRANSFER_AMOUNT, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress())}
	initStateChange := makeInitStateChange(FromTransfer, fromRoute, routes, utest.ADDR)
	sm := &transfer.StateManager{StateTransition, nil}
	assert(t, sm.CurrentState, nil)
	events := sm.Dispatch(initStateChange)
	mstate := sm.CurrentState.(*mediated_transfer.MediatorState)
	assert(t, mstate.OurAddress, utest.ADDR)
	assert(t, mstate.BlockNumber, initStateChange.BlockNumber)
	assert(t, mstate.TransfersPair[0].PayerTransfer, FromTransfer)
	assert(t, mstate.TransfersPair[0].PayerRoute, fromRoute)
	assert(t, len(events) > 0, true, "we have a valid route, the mediated transfer event must be emited")

	var mediatedTransfers []*mediated_transfer.EventSendMediatedTransfer
	for _, e := range events {
		e2, ok := e.(*mediated_transfer.EventSendMediatedTransfer)
		if ok {
			mediatedTransfers = append(mediatedTransfers, e2)
		}
	}
	assert(t, len(mediatedTransfers), 1, "mediated_transfer should /not/ split the transfer")
	mtr := mediatedTransfers[0]
	assert(t, mtr.Token, FromTransfer.Token)
	assert(t, mtr.Amount, FromTransfer.Amount)
	assert(t, mtr.Expiration < FromTransfer.Expiration, true)
	assert(t, mtr.HashLock, FromTransfer.Hashlock)
}

func TestNoValidRoutes(t *testing.T) {
	fromRoute, FromTransfer := utest.MakeFrom(utest.UNIT_TRANSFER_AMOUNT, utest.HOP2, int64(utest.HOP1_TIMEOUT), utils.NewRandomAddress(), utils.EmptyHash)
	var routes = []*transfer.RouteState{utest.MakeRoute(utest.HOP2, utest.UNIT_TRANSFER_AMOUNT-1, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress()),
		utest.MakeRoute(utest.HOP3, 1, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress()),
	}
	initStateChange := makeInitStateChange(FromTransfer, fromRoute, routes, utest.ADDR)
	sm := &transfer.StateManager{StateTransition, nil}
	assert(t, sm.CurrentState, nil)
	events := sm.Dispatch(initStateChange)
	assert(t, sm.CurrentState, nil)
	assert(t, len(events), 1)
	_, ok := events[0].(*mediated_transfer.EventSendRefundTransfer)
	assert(t, ok, true)
}

/*
@pytest.mark.xfail(reason='Not implemented')
def test_lock_timeout_lower_than_previous_channel_settlement_period():
    # For a path A-B-C, B cannot forward a mediated transfer to C with
    # a lock timeout larger than the settlement timeout of the A-B
    # channel.
    #
    # Consider that an attacker controls both nodes A and C:
    #
    # Channels A <-> B and B <-> C have a settlement=10 and B has a
    # reveal_timeout=5
    #
    # (block=1) A -> B [T1 expires=20]
    # (block=1) B -> C [T2 expires=20-5]
    # (block=1) A close channel A-B
    # (block=5) C close channel B-C (waited until lock_expiration=settle_timeout)
    # (block=11) A call settle on channel A-B (settle_timeout is over)
    # (block=12) C call unlock on channel B-C (lock is still valid)
    #
    # If B used min(lock.expiration, previous_channel.settlement)
    #
    # (block=1) A -> B [T1 expires=20]
    # (block=1) B -> C [T2 expires=min(20,10)-5]
    # (block=1) A close channel A-B
    # (block=4) C close channel B-C (waited all possible blocks)
    # (block=5) C call unlock on channel B-C (C is forced to unlock)
    # (block=6) B learns the secret
    # (block=7) B call unlock on channel A-B (settle_timeout is over)
    raise NotImplementedError()


@pytest.mark.xfail(reason='Not implemented. Issue: #382')
def test_do_not_withdraw_an_almost_expiring_lock_if_a_payment_didnt_occur():
    # For a path A1-B-C-A2, an attacker controlling A1 and A2 should not be
    # able to force B-C to close the channel by burning token.
    #
    # The attack would be as follows:
    #
    # - Attacker uses two nodes to open two really cheap channels A1 <-> B and
    #   node A2 <-> C
    # - Attacker sends a mediated message with the lowest possible token
    #   amount from A1 through B and C to A2
    # - Since the attacker controls A1 and A2 it knows the secret, she can choose
    #   when the secret is revealed
    # - The secret is hold back until the hash time lock B->C is almost expiring,
    #   then it's revealed (meaning that the attacker is losing token, that's why
    #   it's using the lowest possible amount)
    # - C wants the token from B, it will reveal the secret and close the channel
    #   (because it must assume the balance proof won't make in time and it needs
    #   to unlock on-chain)
    #
    # Mitigation:
    #
    # - C should only close the channel B-C if he has paid A2, since this may
    #   only happen if the lock for the transfer C-A2 has not yet expired then C
    #   has enough time to follow the protocol without closing the channel B-C.
    raise NotImplementedError()


@pytest.mark.xfail(reason='Not implemented. Issue: #382')
def mediate_transfer_payee_timeout_must_be_lower_than_settlement_and_payer_timeout():
    # Test:
    # - the current payer route/transfer is the reference, not the from_route / from_transfer
    # - the lowest value from blocks_until_settlement and lock expiration must be used
    raise NotImplementedError()


@pytest.mark.xfail(reason='Not implemented. Issue: #382')
def payee_timeout_must_be_lower_than_payer_timeout_minus_reveal_timeout():
    # The payee could reveal the secret on it's lock expiration block, the
    # mediator node will respond with a balance-proof to the payee since the
    # lock is valid and the mediator can safely get the token from the payer,
    # the secret is know and if there are no additional blocks the mediator
    # will be at risk of not being able to withdraw on-chain, so the channel
    # will be closed to safely withdraw.
    #
    # T2.expiration cannot be equal to T1.expiration - reveal_timeout:
    #
    # T1 |---|
    # T2     |---|
    #        ^- reveal the secret
    #        T1.expiration - reveal_timeout == current_block -> withdraw on chain
    #
    # If T2.expiration canot be equal to T1.expiration - reveal_timeout minus ONE:
    #
    # T1 |---|
    # T2      |---|
    #         ^- reveal the secret
    #
    # Race:
    #  1> Secret is learned
    #  2> balance-proof is sent to payee (payee transfer is paid)
    #  3! New block is mined and Raiden learns about it
    #  4> Now the secret is know, the payee is paid, and the current block is
    #     equal to the payer.expiration - reveal-timeout -> withdraw on chain
    #
    # The race is depending on the handling of 3 before 4.
    #
    raise NotImplementedError()

*/
