package target

import (
	"testing"

	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/SmartMeshFoundation/raiden-network/utils/utest"
	"github.com/ethereum/go-ethereum/common"
	assert2 "github.com/stretchr/testify/assert"
)

func assert(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert2.EqualValues(t, expected, actual, msgAndArgs...)
}
func makeInitStateChange(ourAddress common.Address, amount int64, blocknumber int64, initiator common.Address, expire int64) *mediated_transfer.ActionInitTargetStateChange {
	if expire == 0 {
		expire = blocknumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	}
	fromRoute, fromTransfer := utest.MakeFrom(amount, ourAddress, expire, initiator, utils.EmptyHash)
	init := &mediated_transfer.ActionInitTargetStateChange{
		OurAddress:  ourAddress,
		FromRoute:   fromRoute,
		FromTranfer: fromTransfer,
		BlockNumber: blocknumber,
	}
	return init
}

func makeTargetState(ouraddress common.Address, amount, blocknumber int64, initiator common.Address, expire int64) *mediated_transfer.TargetState {
	if expire == 0 {
		expire = blocknumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	}
	fromRoute, fromTransfer := utest.MakeFrom(amount, ouraddress, expire, initiator, utils.EmptyHash)
	state := &mediated_transfer.TargetState{
		OurAddress:   ouraddress,
		FromRoute:    fromRoute,
		FromTransfer: fromTransfer,
		BlockNumber:  blocknumber,
	}
	return state
}

//" Channel must be closed when the unsafe region is reached and the secret is known.
func TestEventsForClose(t *testing.T) {
	var amount int64 = 3
	var expire int64 = 10
	initiator := utest.HOP1
	ourAddress := utest.ADDR
	secret := utest.UNIT_SECRET
	fromRoute, fromTransfer := utest.MakeFrom(amount, ourAddress, expire, initiator, secret)
	safeToWait := expire - int64(fromRoute.RevealTimeout) - 1
	unsafeToWait := expire - int64(fromRoute.RevealTimeout)

	state := &mediated_transfer.TargetState{
		OurAddress:   ourAddress,
		FromRoute:    fromRoute,
		FromTransfer: fromTransfer,
		BlockNumber:  safeToWait,
	}
	events := eventsForClose(state)
	assert(t, len(events), 0)
	state.BlockNumber = unsafeToWait
	events = eventsForClose(state)
	assert(t, len(events) > 0, true)
	ev, ok := events[0].(*mediated_transfer.EventContractSendChannelClose)
	assert(t, ok, true)
	assert(t, fromTransfer.Secret != utils.EmptyHash, true)
	assert(t, ev.ChannelAddress, fromRoute.ChannelAddress)
}

/*
Channel must not be closed when the unsafe region is reached and the
    secret is not known.
*/
func TestEventsForCloseSecretUnkown(t *testing.T) {
	var amount int64 = 3
	var expire int64 = 10
	initiator := utest.HOP1
	ourAddress := utest.ADDR

	fromRoute, fromTransfer := utest.MakeFrom(amount, ourAddress, expire, initiator, utils.EmptyHash)

	state := &mediated_transfer.TargetState{
		OurAddress:   ourAddress,
		FromRoute:    fromRoute,
		FromTransfer: fromTransfer,
		BlockNumber:  expire,
	}
	events := eventsForClose(state)
	assert(t, len(events), 0)
	events = eventsForClose(state)
	assert(t, len(events), 0)
	assert(t, fromTransfer.Secret, utils.EmptyHash)
}

/*
On-chain withdraw must be done if the channel is closed, regardless of
    the unsafe region.
*/
func TestEventsForWithDraw(t *testing.T) {
	var amount int64 = 3
	var expire int64 = 10
	initiator := utest.HOP1
	tr := utest.MakeTransfer(amount, initiator, utest.ADDR, expire, utest.UNIT_SECRET, utils.Sha3(utest.UNIT_SECRET[:]), 1, utest.UNIT_TOKEN_ADDRESS)
	route := utest.MakeRoute(initiator, amount, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress())
	events := eventsForWithdraw(tr, route)
	assert(t, len(events), 0)
	route.State = transfer.CHANNEL_STATE_CLOSED
	events = eventsForWithdraw(tr, route)
	assert(t, len(events) > 0, true)
	ev, ok := events[0].(*mediated_transfer.EventContractSendWithdraw)
	assert(t, ok, true)
	assert(t, ev.ChannelAddress, route.ChannelAddress)
}

/*
Init transfer must send a secret request if the expiration is valid.
*/
func TestHandleInitTarget(t *testing.T) {
	var blockNumber int64 = 1
	var amount int64 = 1
	var expire int64 = int64(utest.UNIT_REVEAL_TIMEOUT) + blockNumber + 1
	initiator := utest.HOP1

	//fromroute,fromtransfer:=utest.MakeFrom(amount,utest.ADDR,expire,initiator,utils.EmptyHash)
	st := makeInitStateChange(utest.ADDR, amount, blockNumber, initiator, expire)
	fromTransfer := st.FromTranfer
	it := handleInitTraget(st)
	assert(t, len(it.Events) > 0, true)
	ev := it.Events[0].(*mediated_transfer.EventSendSecretRequest)

	assert(t, ev.Identifer, fromTransfer.Identifier)
	assert(t, ev.Amount, fromTransfer.Amount)
	assert(t, ev.Hashlock, fromTransfer.Hashlock)
	assert(t, ev.Receiver, initiator)
}

// Init transfer must do nothing if the expiration is bad.
func TestHandleInitTargetBadExpiration(t *testing.T) {
	var blockNumber int64 = 1
	var amount int64 = 1
	var expire int64 = int64(utest.UNIT_REVEAL_TIMEOUT) + blockNumber
	initiator := utest.HOP1

	//fromroute,fromtransfer:=utest.MakeFrom(amount,utest.ADDR,expire,initiator,utils.EmptyHash)
	st := makeInitStateChange(utest.ADDR, amount, blockNumber, initiator, expire)
	it := handleInitTraget(st)
	assert(t, len(it.Events), 0)
}

/*
The target node needs to inform the secret to the previous node to
    receive an updated balance proof.
*/
func TestHandleSecretReveal(t *testing.T) {
	var blockNumber int64 = 1
	var amount int64 = 1
	var expire int64 = int64(utest.UNIT_REVEAL_TIMEOUT) + blockNumber
	initiator := utest.HOP1
	ourAddress := utest.ADDR
	secret := utest.UNIT_SECRET
	state := makeTargetState(ourAddress, amount, blockNumber, initiator, expire)
	stateChange := &mediated_transfer.ReceiveSecretRevealStateChange{
		Secret: secret,
		Sender: initiator,
	}
	//use mediatedTransfer to implement direct transfer
	it := handleSecretReveal(state, stateChange)
	assert(t, len(it.Events), 0)
	//real mediatedTransfere, have a hopnode
	state.FromRoute = utest.MakeRoute(utest.HOP2, amount, utest.UNIT_SETTLE_TIMEOUT, utest.UNIT_REVEAL_TIMEOUT, 0, utils.NewRandomAddress())
	it = handleSecretReveal(state, stateChange)

	assert(t, len(it.Events), 1)
	ev := it.Events[0].(*mediated_transfer.EventSendRevealSecret)
	assert(t, state.State, mediated_transfer.STATE_REVEAL_SECRET)
	assert(t, ev.Identifier, state.FromTransfer.Identifier)
	assert(t, ev.Secret, secret)
	assert(t, ev.Receiver, state.FromRoute.HopNode)
	assert(t, ev.Sender, ourAddress)

}

func TestHandleBlock(t *testing.T) {
	initiator := utest.HOP6
	ourAddress := utest.ADDR
	var amount int64 = 3
	var blockNumber int64 = 1
	expire := blockNumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	state := makeTargetState(ourAddress, amount, blockNumber, initiator, expire)
	newBlock := &transfer.BlockStateChange{blockNumber + 1}
	StateTransiton(state, newBlock)
	assert(t, state.BlockNumber, blockNumber+1)
}

func TestHandleBlockEqualBlockNumber(t *testing.T) {
	initiator := utest.HOP6
	ourAddress := utest.ADDR
	var amount int64 = 3
	var blockNumber int64 = 1
	expire := blockNumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	state := makeTargetState(ourAddress, amount, blockNumber, initiator, expire)
	newBlock := &transfer.BlockStateChange{blockNumber}
	StateTransiton(state, newBlock)
	assert(t, state.BlockNumber, blockNumber)
}
func TestHandleBlockLowerBlockNumber(t *testing.T) {
	initiator := utest.HOP6
	ourAddress := utest.ADDR
	var amount int64 = 3
	var blockNumber int64 = 1
	expire := blockNumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	state := makeTargetState(ourAddress, amount, blockNumber, initiator, expire)
	newBlock := &transfer.BlockStateChange{blockNumber - 1}
	StateTransiton(state, newBlock)
	assert(t, state.BlockNumber, blockNumber)
}

//Clear if the transfer is paid with a proof.
func TestClearIfFinalizedPayed(t *testing.T) {
	initiator := utest.HOP6
	ourAddress := utest.ADDR
	var amount int64 = 3
	var blockNumber int64 = 1
	expire := blockNumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	state := makeTargetState(ourAddress, amount, blockNumber, initiator, expire)
	state.State = mediated_transfer.STATE_BALANCE_PROOF
	it := &transfer.TransitionResult{state, nil}
	it = clearIfFinalized(it)
	assert(t, it.NewState, nil)
}

// Clear expired locks that we don't know the secret for.
func TestClearIfFinalizedExpired(t *testing.T) {
	initiator := utest.HOP6
	ourAddress := utest.ADDR
	var amount int64 = 3
	var blockNumber int64 = 1
	expire := blockNumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	beforestate := makeTargetState(ourAddress, amount, expire, initiator, expire)
	beforeIt := &transfer.TransitionResult{beforestate, nil}
	beforeIt = clearIfFinalized(beforeIt)
	assert(t, beforestate.FromTransfer.Secret, utils.EmptyHash)
	assert(t, beforeIt.NewState != nil, true)

	expiredState := &mediated_transfer.TargetState{
		OurAddress:   ourAddress,
		FromRoute:    beforestate.FromRoute,
		FromTransfer: beforestate.FromTransfer,
		BlockNumber:  expire + 1,
	}
	expireIt := &transfer.TransitionResult{expiredState, nil}
	expireIt = clearIfFinalized(expireIt)
	assert(t, expireIt.NewState == nil, true)
}

func TestStateTransition(t *testing.T) {
	initiator := utest.HOP6
	var amount int64 = 3
	var blockNumber int64 = 1
	expire := blockNumber + int64(utest.UNIT_REVEAL_TIMEOUT)
	fromRoute, fromTransfer := utest.MakeFrom(amount, utest.ADDR, expire, initiator, utils.EmptyHash)
	init := &mediated_transfer.ActionInitTargetStateChange{
		OurAddress:  utest.ADDR,
		FromRoute:   fromRoute,
		FromTranfer: fromTransfer,
		BlockNumber: blockNumber,
	}
	initIt := StateTransiton(nil, init)
	assert(t, initIt.NewState != nil, true)
	newstate := initIt.NewState.(*mediated_transfer.TargetState)
	assert(t, newstate.FromRoute, fromRoute)
	assert(t, newstate.FromTransfer, fromTransfer)

	firstNewBlock := &transfer.BlockStateChange{blockNumber + 1}
	StateTransiton(newstate, firstNewBlock)
	assert(t, newstate.BlockNumber, blockNumber+1)

}

/*
@pytest.mark.xfail(reason='Not implemented #522')
def test_transfer_succesful_after_secret_learned():
    # TransferCompleted event must be used only after the secret is learned and
    # there is enough time to unlock the lock on chain.
    #
    # A mediated transfer might be received during the settlement period of the
    # current channel, the secret request is sent to the initiator and at time
    # the secret is revealed there might not be enough time to safely unlock
    # the token on-chain.
    raise NotImplementedError()

*/
