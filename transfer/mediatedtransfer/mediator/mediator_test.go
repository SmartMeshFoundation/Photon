package mediator

import (
	"testing"

	"math/big"

	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/route"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/utils/utest"
	"github.com/ethereum/go-ethereum/common"
	assert2 "github.com/stretchr/testify/assert"
)

var x = big.NewInt(0)
var big10 = big.NewInt(10)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func assert(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert2.EqualValues(t, expected, actual, msgAndArgs...)
}
func TestTypConversion(t *testing.T) {
	var originalState transfer.State
	_, ok := originalState.(*mediatedtransfer.MediatorState)
	if ok { //for a real nil pointer, it cannot convert
		t.Error("MediatorState StateTransition get type error ")
	}
}

func makeInitStateChange(fromTransfer *mediatedtransfer.LockedTransferState, fromRoute *route.State, routes []*route.State, ourAddress common.Address) *mediatedtransfer.ActionInitMediatorStateChange {
	return &mediatedtransfer.ActionInitMediatorStateChange{
		OurAddress:  ourAddress,
		FromTranfer: fromTransfer,
		Routes:      route.NewRoutesState(routes),
		FromRoute:   fromRoute,
		BlockNumber: 1,
	}
}
func makeTransferPair(payer, payee, initiator, target common.Address, amount *big.Int, expiration int64, secret common.Hash, revealTimeout int) *mediatedtransfer.MediationPairState {
	payerExpiration := expiration
	payeeExpiration := expiration - int64(revealTimeout)
	return mediatedtransfer.NewMediationPairState(
		utest.MakeRoute(payer, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utest.CHANNEL),
		utest.MakeRoute(payee, amount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utest.CHANNEL),
		utest.MakeTransfer(amount, initiator, target, payerExpiration, secret, utils.ShaSecret(secret[:]), utest.UnitTokenAddress),
		utest.MakeTransfer(amount, initiator, target, payeeExpiration, secret, utils.ShaSecret(secret[:]), utest.UnitTokenAddress),
	)
}
func makeTransfersPair(initiator common.Address, hops []common.Address, target common.Address, amount int64, secret common.Hash, initialExpiration int64, revealTimeout int) []*mediatedtransfer.MediationPairState {
	if initialExpiration == 0 {
		initialExpiration = int64((2*len(hops) + 1) * revealTimeout)
	}
	var pairs []*mediatedtransfer.MediationPairState
	nextExpiration := initialExpiration
	for i := 0; i < len(hops)-1; i++ {
		payer := hops[i]
		payee := hops[i+1]
		if nextExpiration <= 0 {
			log.Error("nextExpiration<=0")
		}
		pair := makeTransferPair(payer, payee, initiator, target, big.NewInt(amount), nextExpiration, secret, revealTimeout)
		pairs = append(pairs, pair)
		/*
					    assumes that the node sending the refund will follow the protocol and
			         decrement the expiration for it's lock
		*/
		nextExpiration = pair.PayeeTransfer.Expiration - int64(revealTimeout)
	}
	return pairs
}
func makeMediatorState(fromTransfer *mediatedtransfer.LockedTransferState, fromRoute *route.State, routes []*route.State, ourAddress common.Address) *mediatedtransfer.MediatorState {
	statechange := makeInitStateChange(fromTransfer, fromRoute, routes, ourAddress)
	it := StateTransition(nil, statechange)
	return it.NewState.(*mediatedtransfer.MediatorState)
}

//A hash time lock is valid up to the expiraiton block.
func TestIsLockValid(t *testing.T) {
	var amount = big.NewInt(10)
	var expiration int64 = 10
	initiator := utest.HOP1
	target := utest.HOP2
	tr := utest.MakeTransfer(amount, initiator, target, expiration, utils.EmptyHash, utils.EmptyHash, utest.UnitTokenAddress)
	assert(t, isLockValid(tr, 5), true)
	assert(t, isLockValid(tr, 10), true)
	assert(t, isLockValid(tr, 11), false)

}

/*
It's safe to wait for a secret while there are more than reveal timeout
    blocks until the lock expiration.
*/
func TestIsSafeToWait(t *testing.T) {
	var amount = big.NewInt(10)
	var expiration int64 = 40
	initiator := utest.HOP1
	target := utest.HOP2
	tr := utest.MakeTransfer(amount, initiator, target, expiration, utils.EmptyHash, utils.EmptyHash, utest.UnitTokenAddress)
	//expiration is in 30 blocks, 19 blocks safe for waiting
	assert(t, IsSafeToWait(tr, 10, 1), true)
	assert(t, IsSafeToWait(tr, 10, 20), true)
	assert(t, IsSafeToWait(tr, 10, 29), true)

	assert(t, IsSafeToWait(tr, 10, 30), false)
	assert(t, IsSafeToWait(tr, 10, 50), false)

}

//Don't close the channel if the payee transfer is not paid.
func TestIsRegisterSecretNeededUnpaid(t *testing.T) {
	var amount = big.NewInt(10)
	var expiration int64 = 10
	var revealTimeout = 5

	/*
			    even if the secret is known by the payee, the transfer is paid only if a
		    withdraw on-chain happened or if the mediator has sent a balance proof
	*/
	unpaidStates := []string{mediatedtransfer.StatePayeePending, mediatedtransfer.StatePayeeSecretRevealed}
	for _, unpaidState := range unpaidStates {
		unpaidPair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
		unpaidPair.PayeeState = unpaidState
		assert(t, unpaidPair.PayerRoute.State(), channeltype.StateOpened)
		safeBlock := expiration - int64(revealTimeout) - 1
		assert(t, isSecretRegisterNeeded(unpaidPair, safeBlock), false)
		unsafeBlock := expiration - int64(revealTimeout)
		assert(t, isSecretRegisterNeeded(unpaidPair, unsafeBlock), false)
	}
}

//Close the channel if the payee transfer is paid but the payer has not paid.
func TestIsRegisterSecretNeededPaid(t *testing.T) {
	var amount = big.NewInt(10)
	var expiration int64 = 10
	var revealTimeout = 5

	paidStates := []string{mediatedtransfer.StatePayeeBalanceProof}
	for _, paidState := range paidStates {
		paidPair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
		paidPair.PayeeState = paidState
		assert(t, paidPair.PayerRoute.State(), channeltype.StateOpened)
		safeBlock := expiration - int64(revealTimeout) - 1
		assert(t, isSecretRegisterNeeded(paidPair, safeBlock), false)
		unsafeBlock := expiration - int64(revealTimeout)
		assert(t, isSecretRegisterNeeded(paidPair, unsafeBlock), true)
	}
}

//If the channel is already closed the anser is always no.
func TestIsRegisterSecretNeedChannelClosed(t *testing.T) {
	var amount = big.NewInt(10)
	var expiration int64 = 10
	var revealTimeout = 5

	for state := range mediatedtransfer.ValidPayeeStateMap {
		pair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
		pair.PayeeState = state
		pair.PayerRoute.SetState(channeltype.StateClosed)

		safeBlock := expiration - int64(revealTimeout) - 1
		assert(t, isSecretRegisterNeeded(pair, safeBlock), false)
		unsafeBlock := expiration - int64(revealTimeout)
		assert(t, isSecretRegisterNeeded(pair, unsafeBlock), false)
	}
}

func TestIsRegisterSecretNeededClosed(t *testing.T) {
	var amount = big.NewInt(10)
	var expiration int64 = 10
	var revealTimeout = 5

	pair := makeTransferPair(utest.HOP1, utest.HOP2, utest.HOP3, utest.HOP4, amount, expiration, utils.EmptyHash, revealTimeout)
	pair.PayeeState = mediatedtransfer.StatePayeeBalanceProof

	assert(t, pair.PayerRoute.State(), channeltype.StateOpened)
	safeBlock := expiration - int64(revealTimeout) - 1
	assert(t, isSecretRegisterNeeded(pair, safeBlock), false)
	unsafeBlock := expiration - int64(revealTimeout)
	assert(t, isSecretRegisterNeeded(pair, unsafeBlock), true)

}
func TestIsValidRefund(t *testing.T) {
	initiator := utest.ADDR
	target := utest.HOP1
	validSender := utest.HOP2
	tr := &mediatedtransfer.LockedTransferState{
		Amount:         big.NewInt(30),
		TargetAmount:   big.NewInt(30),
		Fee:            utils.BigInt0,
		Token:          utest.UnitTokenAddress,
		Initiator:      initiator,
		Target:         target,
		Expiration:     50,
		LockSecretHash: utest.UnitHashLock,
		Secret:         utils.EmptyHash,
	}
	refundLowerExpiration := &mediatedtransfer.ReceiveAnnounceDisposedStateChange{
		Lock: &mtree.Lock{
			Expiration:     35,
			Amount:         big.NewInt(30),
			LockSecretHash: utest.UnitHashLock,
		},
		Sender: validSender,
		Token:  utest.UnitTokenAddress,
	}
	rt := utest.MakeRoute(utest.HOP2, big.NewInt(100), utest.UnitSettleTimeout, utest.UnitRevealTimeout, int64(0), utils.NewRandomHash())
	assert(t, IsValidRefund(tr, rt, refundLowerExpiration), false)
	refundLowerExpiration.Sender = target
	assert(t, IsValidRefund(tr, rt, refundLowerExpiration), false)
	refundSameExpiration := &mediatedtransfer.ReceiveAnnounceDisposedStateChange{
		Lock: &mtree.Lock{
			Expiration:     50,
			Amount:         big.NewInt(30),
			LockSecretHash: utest.UnitHashLock,
		},
		Sender: validSender,
		Token:  utest.UnitTokenAddress,
	}
	assert(t, IsValidRefund(tr, rt, refundSameExpiration), true)
}

func TestGetTimeoutBlocks(t *testing.T) {
	var amount = big.NewInt(10)
	initiator := utest.HOP1
	nextHop := utest.HOP2
	settleTimeout := 30
	var blockNumber int64 = 5
	route := utest.MakeRoute(nextHop, amount, settleTimeout, 0, 0, utils.NewRandomHash())
	var earlyExpire int64 = 10
	earlyTransfer := utest.MakeTransfer(amount, initiator, nextHop, earlyExpire, utils.EmptyHash, utils.EmptyHash, utest.UnitTokenAddress)
	earlyBlock := getTimeoutBlocks(route, earlyTransfer, blockNumber)
	assert(t, earlyBlock, 5)

	var equalExpire int64 = 30
	equalTransfer := utest.MakeTransfer(amount, initiator, nextHop, equalExpire, utils.EmptyHash, utils.EmptyHash, utest.UnitTokenAddress)
	equalBlock := getTimeoutBlocks(route, equalTransfer, blockNumber)
	assert(t, equalBlock, 25)

	var largeExpire int64 = 70
	largeTransfer := utest.MakeTransfer(amount, initiator, nextHop, largeExpire, utils.EmptyHash, utils.EmptyHash, utest.UnitTokenAddress)
	largeBlock := getTimeoutBlocks(route, largeTransfer, blockNumber)
	assert(t, largeBlock, 30)

	closedRoute := utest.MakeRoute(nextHop, amount, settleTimeout, 0, 2, utils.NewRandomHash())

	largeBlock = getTimeoutBlocks(closedRoute, largeTransfer, blockNumber)
	assert(t, largeBlock, 27)
	/*
		the computed timeout may be negative, in which case the calling code must /not/ use it
	*/
	negativeBlockNumber := largeExpire
	negativeBlock := getTimeoutBlocks(route, largeTransfer, negativeBlockNumber)
	assert(t, negativeBlock, 0)

}

/*
下一个route主要看，是否有足够的时间（reveal timeout），足够的钱。
*/
//Routes that dont have enough available_balance must be ignored.
func TestNextRouteAmount(t *testing.T) {
	var amount = big.NewInt(10)
	revealTimeout := 30
	timeoutBlocks := 40
	routes := []*route.State{
		utest.MakeRoute(utest.HOP2, x.Add(amount, amount), 0, revealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP1, x.Add(amount, big.NewInt(1)), 0, revealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP3, x.Div(amount, big.NewInt(2)), 0, revealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP4, amount, 0, revealTimeout, 0, utils.NewRandomHash()),
	}
	fromRoute := utest.MakeRoute(utest.HOP6, amount, 0, revealTimeout, 0, utils.NewRandomHash())
	routesState := route.NewRoutesState(routes)
	route1 := nextRoute(fromRoute, routesState, timeoutBlocks, amount, utils.BigInt0)
	assert(t, route1, routes[0])
	assert(t, routesState.AvailableRoutes, routes[1:])
	assert(t, len(routesState.IgnoredRoutes), 0)

	route2 := nextRoute(fromRoute, routesState, timeoutBlocks, amount, utils.BigInt0)
	assert(t, route2, routes[1])
	assert(t, routesState.AvailableRoutes, routes[2:])
	assert(t, len(routesState.IgnoredRoutes), 0)

	route3 := nextRoute(fromRoute, routesState, timeoutBlocks, amount, utils.BigInt0)
	assert(t, route3, routes[3])
	assert(t, len(routesState.AvailableRoutes), 0)
	assert(t, routesState.IgnoredRoutes, []*route.State{routes[2]})

	assert(t, nextRoute(fromRoute, routesState, timeoutBlocks, amount, utils.BigInt0) == nil, true)

}

//Routes with a larger reveal timeout than timeout_blocks must be ignored.
func TestNextRouteRevealTimeout(t *testing.T) {
	var amount = big.NewInt(10)
	var balance = big.NewInt(20)
	timeoutBlocks := 40

	routes := []*route.State{
		utest.MakeRoute(utest.HOP2, balance, 0, timeoutBlocks*2, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP1, balance, 0, timeoutBlocks+1, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP3, balance, 0, timeoutBlocks/2, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP4, balance, 0, timeoutBlocks, 0, utils.NewRandomHash()),
	}
	fromRoute := utest.MakeRoute(utest.HOP6, amount, 0, 10, 0, utils.NewRandomHash())
	routesState := route.NewRoutesState(routes)
	route1 := nextRoute(fromRoute, routesState, timeoutBlocks, amount, utils.BigInt0)
	assert(t, route1, routes[2])
	assert(t, routesState.AvailableRoutes, routes[3:])
	assert(t, routesState.IgnoredRoutes, routes[0:2])
	route2 := nextRoute(fromRoute, routesState, timeoutBlocks, amount, utils.BigInt0)
	assert(t, route2 == nil, true)
	assert(t, len(routesState.AvailableRoutes), 0)
	assert(t, routesState.IgnoredRoutes, append(routes[0:2], routes[3]))
}

func TestNextTransferPair(t *testing.T) {
	timeoutBlocks := 47
	var blockNumber int64 = 3
	var balance = big.NewInt(10)
	initiator := utest.HOP1
	target := utest.ADDR

	payerRoute := utest.MakeRoute(initiator, balance, 0, 0, 0, utils.NewRandomHash())
	payerTransfer := utest.MakeTransfer(balance, initiator, target, 50, utils.EmptyHash, utils.EmptyHash, utest.UnitTokenAddress)

	routes := []*route.State{utest.MakeRoute(utest.HOP2, balance, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash())}
	routesState := route.NewRoutesState(routes)
	pair, events := nextTransferPair(payerRoute, payerTransfer, routesState, timeoutBlocks, blockNumber)

	assert(t, pair.PayerRoute, payerRoute)
	assert(t, pair.PayerTransfer, payerTransfer)
	assert(t, pair.PayeeRoute, routes[0])
	assert(t, pair.PayeeTransfer.Expiration <= pair.PayerTransfer.Expiration, true)

	assert(t, len(events), 1)
	tr, ok := events[0].(*mediatedtransfer.EventSendMediatedTransfer)
	assert(t, ok, true)
	assert(t, tr.Token, payerTransfer.Token)
	assert(t, tr.Amount, payerTransfer.Amount)
	assert(t, tr.LockSecretHash, payerTransfer.LockSecretHash)
	assert(t, tr.Initiator, payerTransfer.Initiator)
	assert(t, tr.Target, payerTransfer.Target)
	assert(t, tr.Expiration <= payerTransfer.Expiration, true)
	assert(t, tr.Receiver, pair.PayeeRoute.HopNode())

	assert(t, len(routesState.AvailableRoutes), 0)
}

func TestSetPayee(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	assert(t, pairs[0].PayerState, mediatedtransfer.StatePayerPending)
	assert(t, pairs[0].PayeeState, mediatedtransfer.StatePayeePending)
	assert(t, pairs[1].PayerState, mediatedtransfer.StatePayerPending)
	assert(t, pairs[1].PayeeState, mediatedtransfer.StatePayeePending)
	setPayeeStateAndCheckRevealOrder(pairs, utest.HOP2, mediatedtransfer.StatePayeeSecretRevealed)
	assert(t, pairs[0].PayerState, mediatedtransfer.StatePayerPending)
	assert(t, pairs[0].PayeeState, mediatedtransfer.StatePayeePending)
	assert(t, pairs[1].PayerState, mediatedtransfer.StatePayerPending)
	assert(t, pairs[1].PayeeState, mediatedtransfer.StatePayeePending)

	setPayeeStateAndCheckRevealOrder(pairs, utest.HOP3, mediatedtransfer.StatePayeeSecretRevealed)
	assert(t, pairs[0].PayerState, mediatedtransfer.StatePayerPending)
	assert(t, pairs[0].PayeeState, mediatedtransfer.StatePayeeSecretRevealed)
	assert(t, pairs[1].PayerState, mediatedtransfer.StatePayerPending)
	assert(t, pairs[1].PayeeState, mediatedtransfer.StatePayeePending)

}

// The transfer pair must switch to expired at the right block.
func TestSetExpiredPairs(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	pair := pairs[0]
	// do not generate events if the secret is not known
	firstUnsafeBlock := pair.PayerTransfer.Expiration - int64(pair.PayerRoute.RevealTimeout())
	setExpiredPairs(pairs, firstUnsafeBlock)
	assert(t, pair.PayeeState, mediatedtransfer.StatePayeePending)
	assert(t, pair.PayerState, mediatedtransfer.StatePayerPending)
	//payee lock expired
	payeeExpirationBlock := pair.PayeeTransfer.Expiration
	setExpiredPairs(pairs, payeeExpirationBlock)
	//dge case for the payee lock expiration
	assert(t, pair.PayeeState, mediatedtransfer.StatePayeePending)
	assert(t, pair.PayerState, mediatedtransfer.StatePayerPending)
	setExpiredPairs(pairs, payeeExpirationBlock+1)
	assert(t, pair.PayeeState, mediatedtransfer.StatePayeeExpired)
	assert(t, pair.PayerState, mediatedtransfer.StatePayerPending)
	// edge case for the payer lock expiration
	payerExpirationBlock := pair.PayerTransfer.Expiration
	setExpiredPairs(pairs, payerExpirationBlock)
	assert(t, pair.PayeeState, mediatedtransfer.StatePayeeExpired)
	assert(t, pair.PayerState, mediatedtransfer.StatePayerPending)
	setExpiredPairs(pairs, payerExpirationBlock+1)
	assert(t, pair.PayeeState, mediatedtransfer.StatePayeeExpired)
	assert(t, pair.PayerState, mediatedtransfer.StatePayerExpired)
}

func TestEventsForRefund(t *testing.T) {
	var amount = big.NewInt(10)
	var expiration int64 = 30
	var revealTimeout = 17
	var timeoutBlocks = int(expiration)
	var blockNumber int64 = 1
	initiator := utest.HOP1
	target := utest.HOP6
	refundRoute := utest.MakeRoute(initiator, amount, utest.UnitSettleTimeout, revealTimeout, 0, utils.NewRandomHash())
	refundTransfer := utest.MakeTransfer(amount, initiator, target, expiration, utils.EmptyHash, utils.EmptyHash, utest.UnitTokenAddress)

	refundEvents := eventsForRefund(refundRoute, refundTransfer)
	ev, ok := refundEvents[0].(*mediatedtransfer.EventSendAnnounceDisposed)
	assert(t, ok, true)
	assert(t, ev.Expiration < blockNumber+int64(timeoutBlocks), true)
	assert(t, ev.Amount, amount)
	assert(t, ev.LockSecretHash, refundTransfer.LockSecretHash)
}

/*
 The secret is revealed backwards to the payer once the payee sent the
    SecretReveal.
Transfer Path 1...2->Me->3...7->Me->4...target
*/
func TestEventsForRevealSecret(t *testing.T) {
	secret := utest.UnitSecret
	ourAddress := utest.ADDR
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	events := eventsForRevealSecret(pairs, ourAddress)
	/*
			   the secret is known by this node, but no other payee is at a secret known
		    state, do nothing
	*/
	assert(t, len(events), 0)
	firstPair := pairs[0]
	lastPair := pairs[1]

	lastPair.PayeeState = mediatedtransfer.StatePayeeSecretRevealed
	events = eventsForRevealSecret(pairs, ourAddress)
	/*
			  the last known hop sent a secret reveal message, this node learned the
		     secret and now must reveal to the payer node from the transfer pair
	*/
	assert(t, len(events), 1)
	ev, ok := events[0].(*mediatedtransfer.EventSendRevealSecret)
	assert(t, ok, true)
	assert(t, ev.Secret, secret)
	assert(t, ev.Receiver, lastPair.PayerRoute.HopNode())
	assert(t, lastPair.PayerState, mediatedtransfer.StatePayerSecretRevealed)

	events = eventsForRevealSecret(pairs, ourAddress)
	/*
			   the payeee from the first_pair did not send a secret reveal message, do
		     nothing
	*/
	assert(t, len(events), 0)
	firstPair.PayeeState = mediatedtransfer.StatePayeeSecretRevealed
	events = eventsForRevealSecret(pairs, ourAddress)
	assert(t, len(events), 1)
	ev, ok = events[0].(*mediatedtransfer.EventSendRevealSecret)
	assert(t, ok, true)
	assert(t, ev.Secret, secret)
	assert(t, ev.Receiver, firstPair.PayerRoute.HopNode())
	assert(t, firstPair.PayerState, mediatedtransfer.StatePayerSecretRevealed)
}

// When the secret is not know there is nothing to do.
func TestEventsForRevealSecretSecretUnkown(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utils.EmptyHash, 0, utest.UnitRevealTimeout)
	events := eventsForRevealSecret(pairs, utest.ADDR)
	assert(t, len(events), 0)

}

// The secret must be revealed backwards to the payer if the payee knows
//the secret.
func TestEventsForRevealSecretAllStates(t *testing.T) {
	secret := utest.UnitSecret
	ourAddress := utest.ADDR
	payeeSecretKnowns := []string{mediatedtransfer.StatePayeeSecretRevealed, mediatedtransfer.StatePayeeBalanceProof}
	for _, state := range payeeSecretKnowns {
		pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, secret, 0, utest.UnitRevealTimeout)
		pair := pairs[0]
		pair.PayeeState = state
		events := eventsForRevealSecret(pairs, ourAddress)
		ev, ok := events[0].(*mediatedtransfer.EventSendRevealSecret)
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
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	lastpair := pairs[0]
	lastpair.PayeeState = mediatedtransfer.StatePayeeSecretRevealed
	// the lock has not expired yet
	blockNumber := lastpair.PayeeTransfer.Expiration
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 2)
	var balanceProof *mediatedtransfer.EventSendBalanceProof
	var unlockSuccess *mediatedtransfer.EventUnlockSuccess

	for _, e := range events {
		switch e2 := e.(type) {
		case *mediatedtransfer.EventSendBalanceProof:
			balanceProof = e2
		case *mediatedtransfer.EventUnlockSuccess:
			unlockSuccess = e2
		}
	}
	assert(t, balanceProof.Receiver, lastpair.PayeeRoute.HopNode())
	assert(t, lastpair.PayeeState, mediatedtransfer.StatePayeeBalanceProof)
	assert(t, unlockSuccess != nil, true)

}

/*
balance proofs are useless if the channel is closed/settled, the payee
    needs to go on-chain and use the latest known balance proof which includes
    this lock in the locksroot.
*/
func TestEventsForBalanceProofChannelClosed(t *testing.T) {
	invalidStates := []channeltype.State{channeltype.StateClosed, channeltype.StateSettled}
	for _, state := range invalidStates {
		pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
		var blockNumber int64 = 5
		lastpair := pairs[0]
		lastpair.PayeeRoute.SetState(state)
		lastpair.PayeeRoute.SetClosedBlock(blockNumber)
		lastpair.PayeeState = mediatedtransfer.StatePayeeSecretRevealed
		events := eventsForBalanceProof(pairs, blockNumber)
		assert(t, len(events), 0)
	}
}

/*
Even though the secret should only propagate from the end of the chain
    to the front, if there is a payee node in the middle that knows the secret
    the balance Proof is sent neverthless.

    This can be done safely because the secret is know to the mediator and
    there is reveal_timeout blocks to withdraw the lock on-chain with the payer.
*/
func TestEventsForBalanceProofMiddleSecret(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4, utest.HOP5}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	var blockNumber int64 = 1
	middlePair := pairs[1]
	middlePair.PayeeState = mediatedtransfer.StatePayeeSecretRevealed
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 2)
	var balanceProof *mediatedtransfer.EventSendBalanceProof
	var unlockSuccess *mediatedtransfer.EventUnlockSuccess

	for _, e := range events {
		switch e2 := e.(type) {
		case *mediatedtransfer.EventSendBalanceProof:
			balanceProof = e2
		case *mediatedtransfer.EventUnlockSuccess:
			unlockSuccess = e2
		}
	}
	assert(t, balanceProof.Receiver, middlePair.PayeeRoute.HopNode())
	assert(t, middlePair.PayeeState, mediatedtransfer.StatePayeeBalanceProof)
	assert(t, unlockSuccess != nil, true)
}

//Nothing to do if the secret is not known.
func TestEventsFroBalanceProofSecretUnkown(t *testing.T) {
	var blockNumber int64 = 1
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utils.EmptyHash, 0, utest.UnitRevealTimeout)
	//pairs[1].PayeeState = mediated_transfer.STATE_PAYEE_SECRET_REVEALED
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 0)

	pairs = makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	/*
			 Even though the secret is set, there is not a single transfer pair with a
		     'secret known' state, so nothing should be done. This state is impossible
		     to reach, in reality someone needs to reveal the secret to the mediator,
		     so at least one other node knows the secret.
	*/
	events = eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 0)
}

//The balance proof should not be sent if the lock has expird.
func TestEventsForBalanceProofLockExpired(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3, utest.HOP4, utest.HOP5}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	lastpair := pairs[len(pairs)-1]
	lastpair.PayeeState = mediatedtransfer.StatePayeeSecretRevealed
	var blockNumber = lastpair.PayeeTransfer.Expiration + 1
	//the lock has expired, do not send a balance proof
	events := eventsForBalanceProof(pairs, blockNumber)
	assert(t, len(events), 0)
	middlePair := pairs[len(pairs)-2]
	middlePair.PayeeState = mediatedtransfer.StatePayeeSecretRevealed
	/*
			     Even though the last node did not receive the payment we should send the
		     balance proof to the middle node to avoid unnecessarely closing the
		     middle channel. This state should not be reached under normal operation,
		     the last hop needs to choose a proper reveal_timeout and must go on-chain
		     to withdraw the token before the lock expires.
	*/
	events = eventsForBalanceProof(pairs, blockNumber)
	var balanceProof *mediatedtransfer.EventSendBalanceProof
	var unlockSuccess *mediatedtransfer.EventUnlockSuccess
	for _, e := range events {
		switch e2 := e.(type) {
		case *mediatedtransfer.EventSendBalanceProof:
			balanceProof = e2
		case *mediatedtransfer.EventUnlockSuccess:
			unlockSuccess = e2
		}
	}
	assert(t, balanceProof.Receiver, middlePair.PayeeRoute.HopNode())
	assert(t, middlePair.PayeeState, mediatedtransfer.StatePayeeBalanceProof)
	assert(t, unlockSuccess != nil, true)
}

// The node must close to unlock on-chain if the payee was paid.
func TestEventsForRegisterSecret(t *testing.T) {
	var states = []string{mediatedtransfer.StatePayeeBalanceProof}
	for _, payeeState := range states {
		pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
		pair := pairs[0]
		pair.PayeeState = payeeState
		blockNumber := pair.PayerTransfer.Expiration - int64(pair.PayeeRoute.RevealTimeout())
		events := eventsForRegisterSecret(pairs, blockNumber)
		assert(t, len(events), 1)
		ev := events[0].(*mediatedtransfer.EventContractSendRegisterSecret)
		assert(t, ev != nil, true)
		assert(t, pair.PayerState, mediatedtransfer.StatePayerWaitingRegisterSecret)
	}
}

/* 比如1..2->me->3...6
	如果6没收到消息，3就只能反复尝试，直到超时。 这个时候path上除了1知道密码以外，其他人都不知道，所以安全，不用关闭channel
If the secret is known but the payee transfer has not being paid the
   node must not settle on-chain, otherwise the payee can burn tokens to
   induce the mediator to close a channel.
*/
func TestEventsForCloseHoldForUnpaidPayee(t *testing.T) {
	pairs := makeTransfersPair(utest.HOP1, []common.Address{utest.HOP2, utest.HOP3}, utest.HOP6, 10, utest.UnitSecret, 0, utest.UnitRevealTimeout)
	pair := pairs[0]
	assert(t, pair.PayerTransfer.Secret, utest.UnitSecret)
	assert(t, pair.PayeeTransfer.Secret, utest.UnitSecret)
	assert(t, stateTransferPaidMaps[pair.PayeeState], false)

	// do not generate events if the secret is known AND the payee is not paid
	firstUnSafeBlock := pair.PayerTransfer.Expiration - int64(pair.PayerRoute.RevealTimeout())
	events := eventsForRegisterSecret(pairs, firstUnSafeBlock)
	assert(t, len(events), 0)
	assert(t, stateTransferPaidMaps[pair.PayeeState], false)
	assert(t, stateTransferPaidMaps[pair.PayerState], false)

	payerExpirationBlock := pair.PayerTransfer.Expiration
	events = eventsForRegisterSecret(pairs, payerExpirationBlock)
	assert(t, len(events), 0)
	assert(t, stateTransferPaidMaps[pair.PayeeState], false)
	assert(t, stateTransferPaidMaps[pair.PayerState], false)

}

func TestSecretLearned(t *testing.T) {
	fromRoute, fromTransfer := utest.MakeFrom(utest.UnitTransferAmount, utest.HOP5, int64(utest.Hop1Timeout), utils.NewRandomAddress(), utils.EmptyHash)
	routes := []*route.State{
		utest.MakeRoute(utest.HOP2, utest.UnitTransferAmount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	state := makeMediatorState(fromTransfer, fromRoute, routes, utils.NewRandomAddress())
	secret := utest.UnitSecret
	payeeAddress := utest.HOP2
	it := secretLearned(state, secret, payeeAddress, mediatedtransfer.StatePayeeSecretRevealed)
	mstate := it.NewState.(*mediatedtransfer.MediatorState)
	transferpair := mstate.TransfersPair[0]
	assert(t, fromTransfer.Expiration >= transferpair.PayeeTransfer.Expiration, true)
	assert(t, transferpair.PayeeTransfer.AlmostEqual(fromTransfer), true)
	assert(t, transferpair.PayeeRoute, routes[0])

	assert(t, transferpair.PayerRoute, fromRoute)
	assert(t, transferpair.PayerTransfer, fromTransfer)
	assert(t, mstate.Secret, secret)
	assert(t, transferpair.PayeeTransfer.Secret, secret)
	assert(t, transferpair.PayerTransfer.Secret, secret)
	assert(t, transferpair.PayeeState, mediatedtransfer.StatePayeeBalanceProof)
	assert(t, transferpair.PayerState, mediatedtransfer.StatePayerSecretRevealed)
	revealNumber := 0
	balanceNumber := 0
	for _, e := range it.Events {
		switch e.(type) {
		case *mediatedtransfer.EventSendRevealSecret:
			revealNumber++
		case *mediatedtransfer.EventSendBalanceProof:
			balanceNumber++
		}
	}
	assert(t, revealNumber, 1)
	assert(t, balanceNumber, 1)
}
func TestMediateTransfer(t *testing.T) {
	var amount = big.NewInt(10)
	var blockNumber int64 = 5
	var expiration int64 = 30
	var routes = []*route.State{utest.MakeRoute(utest.HOP2, utest.UnitTransferAmount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash())}
	routesState := route.NewRoutesState(routes)
	state := &mediatedtransfer.MediatorState{
		OurAddress:  utest.ADDR,
		Routes:      routesState,
		BlockNumber: blockNumber,
		Hashlock:    utest.UnitHashLock,
	}
	payerroute, payertransfer := utest.MakeFrom(amount, utest.HOP6, expiration, utils.NewRandomAddress(), utils.EmptyHash)
	it := mediateTransfer(state, payerroute, payertransfer)
	var eventsMediated []*mediatedtransfer.EventSendMediatedTransfer
	for _, e := range it.Events {
		switch e2 := e.(type) {
		case *mediatedtransfer.EventSendMediatedTransfer:
			eventsMediated = append(eventsMediated, e2)
		}
	}
	assert(t, len(eventsMediated), 1)
	tr := eventsMediated[0]
	assert(t, tr.Token, payertransfer.Token)
	assert(t, tr.Amount, payertransfer.Amount)
	assert(t, tr.LockSecretHash, payertransfer.LockSecretHash)
	assert(t, tr.Target, payertransfer.Target)
	assert(t, payertransfer.Expiration >= tr.Expiration, true)
	assert(t, tr.Receiver, routes[0].HopNode())

}

func TestInitMediator(t *testing.T) {
	fromRoute, FromTransfer := utest.MakeFrom(utest.UnitTransferAmount, utest.HOP2, int64(utest.Hop1Timeout), utils.NewRandomAddress(), utils.EmptyHash)
	var routes = []*route.State{utest.MakeRoute(utest.HOP2, utest.UnitTransferAmount, utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash())}
	initStateChange := makeInitStateChange(FromTransfer, fromRoute, routes, utest.ADDR)
	sm := transfer.NewStateManager(StateTransition, nil, "mediator", utils.ShaSecret([]byte("3")), utils.NewRandomAddress())
	assert(t, sm.CurrentState, nil)
	events := sm.Dispatch(initStateChange)
	mstate := sm.CurrentState.(*mediatedtransfer.MediatorState)
	assert(t, mstate.OurAddress, utest.ADDR)
	assert(t, mstate.BlockNumber, initStateChange.BlockNumber)
	assert(t, mstate.TransfersPair[0].PayerTransfer, FromTransfer)
	assert(t, mstate.TransfersPair[0].PayerRoute, fromRoute)
	assert(t, len(events) > 0, true, "we have a valid route, the mediated transfer event must be emited")

	var mediatedTransfers []*mediatedtransfer.EventSendMediatedTransfer
	for _, e := range events {
		e2, ok := e.(*mediatedtransfer.EventSendMediatedTransfer)
		if ok {
			mediatedTransfers = append(mediatedTransfers, e2)
		}
	}
	assert(t, len(mediatedTransfers), 1, "mediatedtransfer should /not/ split the transfer")
	mtr := mediatedTransfers[0]
	assert(t, mtr.Token, FromTransfer.Token)
	assert(t, mtr.Amount, FromTransfer.Amount)
	assert(t, mtr.Expiration <= FromTransfer.Expiration, true)
	assert(t, mtr.LockSecretHash, FromTransfer.LockSecretHash)
}

func TestNoValidRoutes(t *testing.T) {
	fromRoute, FromTransfer := utest.MakeFrom(utest.UnitTransferAmount, utest.HOP2, int64(utest.Hop1Timeout), utils.NewRandomAddress(), utils.EmptyHash)
	var routes = []*route.State{utest.MakeRoute(utest.HOP2, x.Sub(utest.UnitTransferAmount, big.NewInt(1)), utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
		utest.MakeRoute(utest.HOP3, big.NewInt(1), utest.UnitSettleTimeout, utest.UnitRevealTimeout, 0, utils.NewRandomHash()),
	}
	initStateChange := makeInitStateChange(FromTransfer, fromRoute, routes, utest.ADDR)
	sm := transfer.NewStateManager(StateTransition, nil, "mediator", utils.ShaSecret([]byte("3")), utils.NewRandomAddress())
	assert(t, sm.CurrentState, nil)
	events := sm.Dispatch(initStateChange)
	//assert(t, sm.CurrentState, nil)
	assert(t, len(events), 2)
	_, ok := events[0].(*mediatedtransfer.EventSendAnnounceDisposed)
	assert(t, ok, true)
}
