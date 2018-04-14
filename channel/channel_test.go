package channel

import (
	"testing"

	"math/big"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

var big10 = big.NewInt(10)
var x = big.NewInt(0)

func TestEndState(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	bcs := rpc.MakeTestBlockChainService()
	address1 := bcs.NodeAddress
	address2 := utils.NewRandomAddress()
	channelAddress := utils.NewRandomAddress()

	var balance1 = big.NewInt(70)
	var balance2 = big.NewInt(110)
	lockSecret := utils.Sha3([]byte("test_end_state"))
	var lockAmount = big.NewInt(30)
	var lockExpiration int64 = 10
	lockHashlock := utils.Sha3(lockSecret[:])
	state1 := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	state2 := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	assert.EqualValues(t, state1.ContractBalance, balance1)
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), balance1)
	assert.EqualValues(t, state2.Balance(state1), balance2)
	assert.Equal(t, state1.IsLocked(lockHashlock), false)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)

	assert.Equal(t, state1.TreeState.Tree.MerkleRoot(), utils.EmptyHash)
	assert.Equal(t, state2.TreeState.Tree.MerkleRoot(), utils.EmptyHash)
	assert.EqualValues(t, state1.Nonce(), 0)
	assert.EqualValues(t, state2.Nonce(), 0)
	lock := &encoding.Lock{lockExpiration, lockAmount, lockHashlock}
	lockHash := utils.Sha3(lock.AsBytes())
	var transferedAmount = utils.BigInt0
	_, locksroot := state2.ComputeMerkleRootWith(lock)
	/*
		identifier uint64, nonce int64, token common.Address,
		channel common.Address, transferAmount *big.Int,
		recipient common.Address, locksroot common.Hash, lock *Lock,
		target common.Address, initiator common.Address, fee int64
	*/
	mediated_transfer := encoding.NewMediatedTransfer(1, 1, tokenAddress, channelAddress, transferedAmount, state2.Address, locksroot,
		lock, utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0)
	mediated_transfer.Sign(bcs.PrivKey, mediated_transfer)
	state1.RegisterLockedTransfer(mediated_transfer)
	assert.EqualValues(t, state1.ContractBalance, balance1)
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), balance1)
	assert.EqualValues(t, state2.Balance(state1), balance2)

	assert.EqualValues(t, state1.Distributable(state2), new(big.Int).Sub(balance1, lockAmount))
	assert.EqualValues(t, state2.Distributable(state1), balance2)

	assert.EqualValues(t, state1.AmountLocked(), lockAmount)
	assert.EqualValues(t, state2.AmountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), true)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.TreeState.Tree.MerkleRoot(), lockHash)
	assert.Equal(t, state2.TreeState.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.Nonce(), 1)
	assert.EqualValues(t, state2.Nonce(), 0)
	if state1.UpdateContractBalance(new(big.Int).Sub(balance1, big10)) != errBalanceDecrease {
		t.Error(errBalanceDecrease)
	}
	assert.Equal(t, state1.UpdateContractBalance(new(big.Int).Add(balance1, big10)), nil)
	assert.EqualValues(t, state1.ContractBalance, new(big.Int).Add(balance1, big10))
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), new(big.Int).Add(balance1, big10))
	assert.EqualValues(t, state2.Balance(state1), balance2)
	x := new(big.Int).Sub(balance1, lockAmount)
	assert.EqualValues(t, state1.Distributable(state2), x.Add(x, big10))
	assert.EqualValues(t, state1.AmountLocked(), lockAmount)
	assert.EqualValues(t, state2.AmountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), true)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.TreeState.Tree.MerkleRoot(), lockHash)
	assert.Equal(t, state2.TreeState.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.Nonce(), 1)
	assert.EqualValues(t, state2.Nonce(), 0)

	state1.RegisterSecret(lockSecret)
	assert.EqualValues(t, state1.ContractBalance, x.Add(balance1, big10))
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), x.Add(balance1, big10))
	assert.EqualValues(t, state2.Balance(state1), balance2)

	assert.EqualValues(t, state1.Distributable(state2), x.Sub(balance1, lockAmount).Add(x, big10))
	assert.EqualValues(t, state1.AmountLocked(), lockAmount)
	assert.EqualValues(t, state2.AmountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), false)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.TreeState.Tree.MerkleRoot(), lockHash)
	assert.Equal(t, state2.TreeState.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.Nonce(), 1)
	assert.EqualValues(t, state2.Nonce(), 0)
	secretMessage := encoding.NewSecret(1, 2, channelAddress, x.Add(transferedAmount, lockAmount), utils.EmptyHash, lockSecret)
	secretMessage.Sign(bcs.PrivKey, secretMessage)
	state1.RegisterSecretMessage(secretMessage)

	assert.EqualValues(t, state1.ContractBalance, x.Add(balance1, big10))
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), x.Add(balance1, big10).Sub(x, lockAmount))
	assert.EqualValues(t, state2.Balance(state1), x.Add(balance2, lockAmount))

	assert.EqualValues(t, state1.Distributable(state2), x.Sub(balance1, lockAmount).Add(x, big10))
	assert.EqualValues(t, state2.Distributable(state1), x.Add(balance2, lockAmount))
	assert.EqualValues(t, state1.AmountLocked(), utils.BigInt0)
	assert.EqualValues(t, state2.AmountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), false)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.TreeState.Tree.MerkleRoot(), utils.EmptyHash)
	assert.Equal(t, state2.TreeState.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.Nonce(), 2)
	assert.EqualValues(t, state2.Nonce(), 0)
}
func makeExternState() *ChannelExternalState {
	bcs := newTestBlockChainService()
	//must provide a valid netting channel address
	nettingChannel, _ := bcs.NettingChannel(common.HexToAddress("0x93b84FF17268b6a2636D94Ecc58949527BB4ac9d"))
	return NewChannelExternalState(func(channel *Channel, hashlock common.Hash) {}, nettingChannel, nettingChannel.Address, bcs, NewMockChannelDb())
}
func TestSenderCannotOverSpend(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	bcs := rpc.MakeTestBlockChainService()
	address1 := bcs.NodeAddress
	privkey1 := bcs.PrivKey
	address2 := utils.NewRandomAddress()
	var balance1 = big.NewInt(70)
	var balance2 = big.NewInt(110)
	revealTimeout := 5
	settleTimeout := 15
	var blockNumber int64 = 10
	ourState := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, externState.ChannelAddress, bcs, revealTimeout, settleTimeout)
	amount := balance1
	expiration := blockNumber + int64(settleTimeout)
	sent_mediated_transfer0, _ := testChannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount, 1, expiration, utils.Sha3([]byte("test_locked_amount_cannot_be_spent")))
	sent_mediated_transfer0.Sign(privkey1, sent_mediated_transfer0)
	testChannel.RegisterTransfer(blockNumber, sent_mediated_transfer0)
	lock2 := &encoding.Lock{expiration, amount, utils.Sha3([]byte("test_locked_amount_cannot_be_spent2"))}
	leaves := []common.Hash{utils.Sha3(sent_mediated_transfer0.GetLock().AsBytes()), utils.Sha3(lock2.AsBytes())}
	tree2, _ := transfer.NewMerkleTree(leaves)
	locksroot2 := tree2.MerkleRoot()
	sent_mediated_transfer1 := encoding.NewMediatedTransfer(2, sent_mediated_transfer0.Nonce+1, tokenAddress, testChannel.MyAddress, big.NewInt(0), address2, locksroot2, lock2, address2, address1, utils.BigInt0)
	sent_mediated_transfer1.Sign(privkey1, sent_mediated_transfer1)
	err := testChannel.RegisterTransfer(blockNumber, sent_mediated_transfer1)
	if err != rerr.InsufficientBalance {
		t.Error(err)
	}
}
func TestReceiverCannotSpendLockedAmount(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	bcs := rpc.MakeTestBlockChainService()
	privkey1, address1 := utils.MakePrivateKeyAddress()
	privkey2, address2 := utils.MakePrivateKeyAddress()
	var balance1 = big.NewInt(33)
	var balance2 = big.NewInt(11)
	revealTimeout := 7
	settleTimeout := 11
	var blockNumber int64 = 7
	ourState := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, externState.ChannelAddress, bcs, revealTimeout, settleTimeout)
	amount1 := balance2
	expiration := blockNumber + int64(settleTimeout)
	receiveMediatedTransfer0, _ := testChannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount1, 1, expiration, utils.Sha3([]byte("test_locked_amount_cannot_be_spent")))
	receiveMediatedTransfer0.Sign(privkey2, receiveMediatedTransfer0)
	err := testChannel.RegisterTransfer(blockNumber, receiveMediatedTransfer0)
	if err != nil {
		t.Error(err)
	}
	t.Log("after tr1 channel=", testChannel.String())
	amount2 := x.Add(balance1, big.NewInt(1))
	lock2 := &encoding.Lock{expiration, amount2, utils.Sha3([]byte("lxllx"))}
	tree2, _ := transfer.NewMerkleTree([]common.Hash{utils.Sha3(lock2.AsBytes())})
	locksroot2 := tree2.MerkleRoot()
	sendMediatedTransfer0 := encoding.NewMediatedTransfer(1, 1, tokenAddress, testChannel.MyAddress, big.NewInt(0), address2, locksroot2, lock2, address2, address1, utils.BigInt0)
	sendMediatedTransfer0.Sign(privkey1, sendMediatedTransfer0)
	if testChannel.RegisterTransfer(blockNumber, sendMediatedTransfer0) != rerr.InsufficientBalance {
		t.Error("RegisterTransfer should be failed ")
	}
	t.Log("after tr2 channel=", testChannel.String())
}

func TestInvalidTimeouts(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	reavealTimeout := 5
	settleTimeout := 15
	address1 := utils.NewRandomAddress()
	address2 := utils.NewRandomAddress()
	var balance1 = big.NewInt(10)
	var balance2 = big.NewInt(10)
	ourState := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	externState := makeExternState()
	bcs := rpc.MakeTestBlockChainService()
	_, err := NewChannel(ourState, partnerState, externState, externState.ChannelAddress, tokenAddress, bcs, 50, 49)
	if err == nil {
		t.Error("should failed")
	}
	for _, invalidValue := range []int{-1, 0, 1} {
		_, err = NewChannel(ourState, partnerState, externState, externState.ChannelAddress, tokenAddress, bcs, invalidValue, settleTimeout)
		assert.NotEqual(t, err, nil)
		_, err = NewChannel(ourState, partnerState, externState, externState.ChannelAddress, tokenAddress, bcs, reavealTimeout, invalidValue)
		assert.NotEqual(t, err, nil)
	}
}
func TestPythonChannel(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	reavealTimeout := 5
	settleTimeout := 15
	privkey1, address1 := utils.MakePrivateKeyAddress()
	address2 := utils.NewRandomAddress()
	var balance1 = big.NewInt(70)
	var balance2 = big.NewInt(110)
	var blockNumber int64 = 10
	ourState := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	externState := makeExternState()
	bcs := rpc.MakeTestBlockChainService()
	testchannel, _ := NewChannel(ourState, partnerState, externState, externState.ChannelAddress, tokenAddress, bcs, reavealTimeout, settleTimeout)
	_, err := testchannel.CreateDirectTransfer(big.NewInt(-10), 1)
	assert.NotEqual(t, err, nil)
	_, err = testchannel.CreateDirectTransfer(x.Add(balance1, big10), 1)
	assert.NotEqual(t, err, nil)
	var amount1 = big.NewInt(10)
	directTransfer, _ := testchannel.CreateDirectTransfer(amount1, 1)
	directTransfer.Sign(privkey1, directTransfer)
	testchannel.RegisterTransfer(blockNumber, directTransfer)

	assert.EqualValues(t, testchannel.ContractBalance(), balance1)
	assert.EqualValues(t, testchannel.Balance(), x.Sub(balance1, amount1))
	assert.EqualValues(t, testchannel.TransferAmount(), amount1)
	assert.EqualValues(t, testchannel.Distributable(), x.Sub(balance1, amount1))
	assert.EqualValues(t, testchannel.Outstanding(), utils.BigInt0)
	assert.EqualValues(t, testchannel.Locked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.OurState.AmountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.PartnerState.AmountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.GetNextNonce(), 2)

	secret := utils.Sha3([]byte("test_channel"))
	hashlock := utils.Sha3(secret[:])
	var amount2 = big.NewInt(10)
	expiration := blockNumber + int64(settleTimeout) - 5
	var identifier uint64 = 2
	mediatedTransfer, _ := testchannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount2, identifier, expiration, hashlock)
	mediatedTransfer.Sign(privkey1, mediatedTransfer)
	testchannel.RegisterTransfer(blockNumber, mediatedTransfer)

	assert.EqualValues(t, testchannel.ContractBalance(), balance1)
	assert.EqualValues(t, testchannel.Balance(), x.Sub(balance1, amount1))
	assert.EqualValues(t, testchannel.TransferAmount(), amount1)
	assert.EqualValues(t, testchannel.Distributable(), x.Sub(balance1, amount1).Sub(x, amount2))
	assert.EqualValues(t, testchannel.Outstanding(), utils.BigInt0)
	assert.EqualValues(t, testchannel.Locked(), amount2)
	assert.EqualValues(t, testchannel.OurState.AmountLocked(), amount2)
	assert.EqualValues(t, testchannel.PartnerState.AmountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.GetNextNonce(), 3)

	secretMessage, _ := testchannel.CreateSecret(identifier, secret)
	secretMessage.Sign(privkey1, secretMessage)
	log.Info(fmt.Sprintf("secret message=%s", utils.StringInterface(secretMessage, 4)))
	log.Info("bofore reg sec proof=%s", utils.StringInterface(testchannel.OurState.BalanceProofState, 2))
	err = testchannel.RegisterTransfer(blockNumber, secretMessage)
	if err != nil {
		t.Error(err)
	}
	log.Info("after reg sec proof=%s", utils.StringInterface(testchannel.OurState.BalanceProofState, 2))
	assert.EqualValues(t, testchannel.ContractBalance(), balance1)
	assert.EqualValues(t, testchannel.Balance(), x.Sub(balance1, amount1).Sub(x, amount2))
	assert.EqualValues(t, testchannel.TransferAmount(), x.Add(amount1, amount2))
	assert.EqualValues(t, testchannel.Distributable(), x.Sub(balance1, amount1).Sub(x, amount2))
	assert.EqualValues(t, testchannel.Outstanding(), utils.BigInt0)
	assert.EqualValues(t, testchannel.Locked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.OurState.AmountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.OurState.AmountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.PartnerState.AmountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.GetNextNonce(), 4)

}

//The nonce must increase with each new transfer.
func TestChannelIncreaseNonceAndTransferedAmount(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	reavealTimeout := 5
	settleTimeout := 15
	privkey1, address1 := utils.MakePrivateKeyAddress()
	address2 := utils.NewRandomAddress()
	var balance1 = big.NewInt(70)
	var balance2 = big.NewInt(110)
	var blockNumber int64 = 1
	ourState := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	externState := makeExternState()
	bcs := rpc.MakeTestBlockChainService()
	tch, _ := NewChannel(ourState, partnerState, externState, externState.ChannelAddress, tokenAddress, bcs, reavealTimeout, settleTimeout)
	previousNonce := tch.GetNextNonce()
	previousTransfered := tch.TransferAmount()
	var amount = big.NewInt(7)
	for i := 0; i < 10; i++ {
		directTransfer, _ := tch.CreateDirectTransfer(amount, 1)
		directTransfer.Sign(privkey1, directTransfer)
		tch.RegisterTransfer(blockNumber, directTransfer)
		newNonce := tch.GetNextNonce()
		newTransfered := tch.TransferAmount()
		assert.EqualValues(t, newNonce, previousNonce+1)
		assert.EqualValues(t, newTransfered, x.Add(previousTransfered, amount))
		previousNonce = tch.GetNextNonce()
		previousTransfered = tch.TransferAmount()
	}
}
func makePairChannel() (*Channel, *Channel) {
	tokenAddress := utils.NewRandomAddress()
	externState1 := makeExternState()
	externState2 := makeExternState()
	var balance1 = big.NewInt(330)
	var balance2 = big.NewInt(110)
	revealTimeout := 7
	settleTimeout := 30
	ourState := NewChannelEndState(externState1.bcs.NodeAddress, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(externState2.bcs.NodeAddress, balance2, nil, transfer.EmptyMerkleTreeState)

	testChannel, _ := NewChannel(ourState, partnerState, externState1, tokenAddress, externState1.ChannelAddress, externState1.bcs, revealTimeout, settleTimeout)

	ourState = NewChannelEndState(externState1.bcs.NodeAddress, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState = NewChannelEndState(externState2.bcs.NodeAddress, balance2, nil, transfer.EmptyMerkleTreeState)
	testChannel2, _ := NewChannel(partnerState, ourState, externState2, tokenAddress, externState2.ChannelAddress, externState2.bcs, revealTimeout, settleTimeout)
	return testChannel, testChannel2
}

/*
Assert that `channel0` has a correct `partner_state` to represent
    `channel1` and vice-versa.
*/
func assertMirror(ch0, ch1 *Channel, t *testing.T) {
	unclaimed0 := ch0.OurState.TreeState.Tree.MerkleRoot()
	unclaimed1 := ch1.PartnerState.TreeState.Tree.MerkleRoot()
	assert.EqualValues(t, unclaimed0, unclaimed1)

	assert.EqualValues(t, ch0.OurState.AmountLocked(), ch1.PartnerState.AmountLocked())
	assert.EqualValues(t, ch0.TransferAmount(), ch1.PartnerState.TransferAmount())
	balance0 := ch0.OurState.Balance(ch0.PartnerState)
	balance1 := ch1.PartnerState.Balance(ch1.OurState)
	assert.EqualValues(t, balance0, balance1)

	assert.EqualValues(t, ch0.Distributable(), ch0.OurState.Distributable(ch0.PartnerState))
	assert.EqualValues(t, ch0.Distributable(), ch1.PartnerState.Distributable(ch1.OurState))

	unclaimed0 = ch1.OurState.TreeState.Tree.MerkleRoot()
	unclaimed1 = ch0.PartnerState.TreeState.Tree.MerkleRoot()
	assert.EqualValues(t, unclaimed0, unclaimed1)

	assert.EqualValues(t, ch1.OurState.AmountLocked(), ch0.PartnerState.AmountLocked())
	assert.EqualValues(t, ch1.TransferAmount(), ch0.PartnerState.TransferAmount())
	balance0 = ch1.OurState.Balance(ch1.PartnerState)
	balance1 = ch0.PartnerState.Balance(ch0.OurState)
	assert.EqualValues(t, balance0, balance1)

	assert.EqualValues(t, ch1.Distributable(), ch1.OurState.Distributable(ch1.PartnerState))
	assert.EqualValues(t, ch1.Distributable(), ch0.PartnerState.Distributable(ch0.OurState))
}

//Assert the locks created from `from_channel`.
func assertLocked(ch *Channel, pendingLocks []*encoding.Lock, t *testing.T) {
	var root common.Hash
	if pendingLocks != nil {
		var leaves []common.Hash
		for _, lock := range pendingLocks {
			leaves = append(leaves, utils.Sha3(lock.AsBytes()))
		}
		tree, _ := transfer.NewMerkleTree(leaves)
		root = tree.MerkleRoot()
	}
	assert.EqualValues(t, len(ch.OurState.Lock2PendingLocks), len(pendingLocks))
	assert.EqualValues(t, ch.OurState.TreeState.Tree.MerkleRoot(), root)
	var sum = big.NewInt(0)
	for _, lock := range pendingLocks {
		sum.Add(sum, lock.Amount)
	}
	assert.EqualValues(t, ch.OurState.AmountLocked(), sum)
	for _, lock := range pendingLocks {
		assert.Equal(t, ch.OurState.IsLocked(lock.HashLock), true)
	}
}

//Assert the from_channel overall token values.
func assertBalance(ch *Channel, balance, outstanding, distributable *big.Int, t *testing.T) {
	assert.EqualValues(t, ch.Balance(), balance)
	assert.EqualValues(t, ch.Distributable(), distributable)
	assert.EqualValues(t, ch.Outstanding(), outstanding)
	/*
			     the amount of token locked in the partner end of the from_channel is equal to how much
		     we have outstanding
	*/
	assert.EqualValues(t, ch.PartnerState.AmountLocked(), outstanding)
	assert.EqualValues(t, ch.Balance(), ch.OurState.Balance(ch.PartnerState))
	assert.EqualValues(t, ch.Balance().Cmp(utils.BigInt0) >= 0, true)
	assert.EqualValues(t, ch.Distributable().Cmp(utils.BigInt0) >= 0, true)
	assert.EqualValues(t, ch.Locked().Cmp(utils.BigInt0) >= 0, true)
	assert.EqualValues(t, ch.Balance(), x.Add(ch.Locked(), ch.Distributable()))
}

/*
Assert the values of two synched channels.

    Note:
        This assert does not work if for a intermediate state, were one message
        hasn't being delivered yet or has been completely lost.
*/
func assertSyncedChannels(ch0 *Channel, balance0 *big.Int, outstandingLocks0 []*encoding.Lock, ch1 *Channel, balance1 *big.Int, outstandingLocks1 []*encoding.Lock, t *testing.T) {
	totalToken := new(big.Int).Set(x.Add(ch0.ContractBalance(), ch1.ContractBalance()))
	assert.EqualValues(t, totalToken, x.Add(ch0.Balance(), ch1.Balance()))

	var lockedAmount0 = big.NewInt(0)
	for _, lock := range outstandingLocks0 {
		lockedAmount0.Add(lockedAmount0, lock.Amount)
	}
	var lockedAmount1 = big.NewInt(0)
	for _, lock := range outstandingLocks1 {
		lockedAmount1.Add(lockedAmount1, lock.Amount)
	}
	assertBalance(ch0, balance0, lockedAmount0, x.Sub(ch0.Balance(), lockedAmount1), t)
	assertBalance(ch1, balance1, lockedAmount1, x.Sub(ch1.Balance(), lockedAmount0), t)
	assertLocked(ch0, outstandingLocks1, t)
	assertLocked(ch1, outstandingLocks0, t)
	assertMirror(ch0, ch1, t)
}
func TestSetup(t *testing.T) {
	ch0, ch1 := makePairChannel()
	assertSyncedChannels(ch0, ch0.Balance(), nil, ch1, ch1.Balance(), nil, t)
}

/*
Can keep doing transactions even if not all secrets have been released.
*/
func TestInterwovenTransfers(t *testing.T) {
	var err error
	ArgNumberOfTransfers := 10 //To make sure if there is  money can be transferred
	ch0, ch1 := makePairChannel()
	contractBalance0 := ch0.ContractBalance()
	contractBalance1 := ch1.ContractBalance()
	var unclaimedLocks []*encoding.Lock
	var transfersList []*encoding.MediatedTransfer
	var transfersClaimed []bool
	var transfersAmount []*big.Int
	var transfersSecret []common.Hash
	for i := 1; i <= ArgNumberOfTransfers; i++ {
		transfersAmount = append(transfersAmount, big.NewInt(int64(i)))
		transfersSecret = append(transfersSecret, utils.Sha3(utils.Random(32)))
	}
	var claimedAmount = big.NewInt(0)
	var distributedAmount = big.NewInt(0)
	var blockNumber int64 = 7
	var settleTimeout int64 = 30
	logState := func() {

	}
	for i := 0; i < len(transfersAmount); i++ {
		amount := transfersAmount[i]
		secret := transfersSecret[i]
		expiration := blockNumber + settleTimeout - 1
		identifier := uint64(i)
		var mtr *encoding.MediatedTransfer
		mtr, err = ch0.CreateMediatedTransfer(ch0.OurState.Address, ch1.OurState.Address, utils.BigInt0, amount, identifier, expiration, utils.Sha3(secret[:]))
		assert.Equal(t, err, nil)
		mtr.Sign(ch0.ExternState.bcs.PrivKey, mtr)
		err = ch0.RegisterTransfer(blockNumber, mtr)
		assert.Equal(t, err, nil)
		err = ch1.RegisterTransfer(blockNumber, mtr)
		assert.Equal(t, err, nil)
		distributedAmount.Add(distributedAmount, amount)
		transfersClaimed = append(transfersClaimed, false)
		transfersList = append(transfersList, mtr)
		unclaimedLocks = append(unclaimedLocks, mtr.GetLock())
		logState()
		assertSyncedChannels(ch0, new(big.Int).Sub(contractBalance0, claimedAmount), nil, ch1, new(big.Int).Add(contractBalance1, claimedAmount), unclaimedLocks, t)
		assert.EqualValues(t, ch0.Distributable(), x.Sub(contractBalance0, distributedAmount))
		/*
					 claim a transaction at every other iteration, leaving the current one
			        in place
		*/
		if i > 0 && i%2 == 0 {
			transfer := transfersList[i-1]
			secret := transfersSecret[i-1]
			//synchronized claiming
			secretMessage, _ := ch0.CreateSecret(identifier, secret)
			secretMessage.Sign(ch0.ExternState.bcs.PrivKey, secretMessage)
			err = ch0.RegisterTransfer(blockNumber, secretMessage)
			assert.Equal(t, err, nil)
			err = ch1.RegisterTransfer(blockNumber, secretMessage)
			assert.Equal(t, err, nil)
			//update test state
			claimedAmount.Add(claimedAmount, transfer.GetLock().Amount)
			transfersClaimed[i-1] = true
			unclaimedLocks = nil
			for i := 0; i < len(transfersList); i++ {
				if !transfersClaimed[i] {
					unclaimedLocks = append(unclaimedLocks, transfersList[i].GetLock())
				}
			}
			logState()
			//test the state of the channels after the claim
			assertSyncedChannels(ch0, new(big.Int).Sub(contractBalance0, claimedAmount), nil,
				ch1, new(big.Int).Add(contractBalance1, claimedAmount), unclaimedLocks, t)
			assert.EqualValues(t, ch0.Distributable(), x.Sub(contractBalance0, distributedAmount))
		}
	}
}

func TestTransfer(t *testing.T) {
	ch0, ch1 := makePairChannel()
	var amount = big.NewInt(10)
	directTransfer, err := ch0.CreateDirectTransfer(amount, 1)
	assert.Equal(t, err, nil)
	directTransfer.Sign(ch0.ExternState.bcs.PrivKey, directTransfer)
	err = ch0.RegisterTransfer(10, directTransfer)
	assert.Equal(t, err, nil)
	err = ch1.RegisterTransfer(10, directTransfer)
	assert.Equal(t, err, nil)
	assertSyncedChannels(ch0, x.Sub(ch0.ContractBalance(), amount), nil,
		ch1, x.Add(ch1.ContractBalance(), amount), nil, t)
}

/*
Regression test for registration of invalid transfer.

    The bug occurred if a transfer with an invalid allowance but a valid secret
    was registered, when the local end registered the transfer it would
    "unlock" the partners' token, but the transfer wouldn't be sent because the
    allowance check failed, leaving the channel in an inconsistent state.
*/
func TestRegisterInvalidTransfer(t *testing.T) {
	settleTimeout := 30
	ch0, ch1 := makePairChannel()
	balance0 := ch0.Balance()
	balance1 := ch1.Balance()
	var amount = big.NewInt(10)
	var blockNumber int64 = 10
	expiration := blockNumber + int64(settleTimeout) - 1
	secret := utils.Sha3([]byte("secret"))
	hashlock := utils.Sha3(secret[:])
	transfer1, err := ch0.CreateMediatedTransfer(ch0.OurState.Address, ch1.OurState.Address, utils.BigInt0, amount, 1, expiration, hashlock)
	assert.Equal(t, err, nil)
	transfer1.Sign(ch0.ExternState.bcs.PrivKey, transfer1)
	err = ch0.RegisterTransfer(blockNumber, transfer1)
	assert.Equal(t, err, nil)
	err = ch1.RegisterTransfer(blockNumber, transfer1)
	assert.Equal(t, err, nil)
	assertSyncedChannels(ch0, balance0, nil,
		ch1, balance1, []*encoding.Lock{transfer1.GetLock()}, t)
	// handcrafted transfer because channel.create_transfer won't create it
	transfer2 := encoding.NewDirectTransfer(1, ch0.GetNextNonce(), ch0.TokenAddress, ch0.MyAddress, x.Add(ch1.Balance(), balance0).Add(x, amount), ch0.PartnerState.Address, ch0.PartnerState.TreeState.Tree.MerkleRoot())
	transfer2.Sign(ch0.ExternState.bcs.PrivKey, transfer2)
	err = ch0.RegisterTransfer(blockNumber, transfer2)
	assert.Equal(t, err != nil, true)
	err = ch1.RegisterTransfer(blockNumber, transfer2)
	assert.Equal(t, err != nil, true)
	assertSyncedChannels(ch0, balance0, nil,
		ch1, balance1, []*encoding.Lock{transfer1.GetLock()}, t)
}

/*
A node may go offline for an undetermined period of time, and when it
    comes back online it must accept the messages that are waiting, otherwise
    the partner node won't make progress with its queue.

    If a N node goes offline for a number B of blocks, and the partner does not
    close the channel, when N comes back online some of the messages from its
    partner may become expired. Neverthless these messages are ordered and must
    be accepted for the partner to make progress with its queue.

    Note: Accepting a message with an expired lock does *not* imply the token
    transfer happened, and the receiver node must *not* forward the transfer,
    only accept the message allowing the partner to progress with its message
    queue.
*/
func TestChannelMustAcceptExpiredLocks(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	bcs := rpc.MakeTestBlockChainService()
	_, address1 := utils.MakePrivateKeyAddress()
	privkey2, address2 := utils.MakePrivateKeyAddress()
	var balance1 = big.NewInt(33)
	var balance2 = big.NewInt(11)
	revealTimeout := 7
	settleTimeout := 11
	var blockNumber int64 = 7
	ourState := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, externState.ChannelAddress, bcs, revealTimeout, settleTimeout)
	lock := &encoding.Lock{blockNumber + int64(settleTimeout), big.NewInt(1), utils.EmptyHash}
	transfer := encoding.NewMediatedTransfer(1, testChannel.GetNextNonce(), testChannel.TokenAddress, testChannel.MyAddress, big.NewInt(1), address1, utils.Sha3(lock.AsBytes()), lock, utils.EmptyAddress, utils.EmptyAddress, utils.BigInt0)
	transfer.Sign(privkey2, transfer)
	err := testChannel.RegisterTransfer(blockNumber+int64(settleTimeout)+1, transfer)
	assert.Equal(t, err, nil)
}

func TestRemoveExpiredHashlock(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	bcs := rpc.MakeTestBlockChainService()
	privkey1, address1 := utils.MakePrivateKeyAddress()
	privkey2, address2 := utils.MakePrivateKeyAddress()
	var balance1 = big.NewInt(33)
	var balance2 = big.NewInt(11)
	revealTimeout := 7
	settleTimeout := 11
	var blockNumber int64 = 7
	ourState := NewChannelEndState(address1, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(address2, balance2, nil, transfer.EmptyMerkleTreeState)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, externState.ChannelAddress, bcs, revealTimeout, settleTimeout)
	amount1 := balance2
	expiration := blockNumber + int64(settleTimeout)
	//smtr: the mediated transfer i sent out
	smtr, _ := testChannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount1, 1, expiration, utils.Sha3([]byte("test_locked_amount_cannot_be_spent")))
	smtr.Sign(privkey1, smtr)
	err := testChannel.RegisterTransfer(blockNumber, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("after tr1 channel=", testChannel.String())
	amount2 := balance2
	lock2 := &encoding.Lock{expiration, amount2, utils.Sha3([]byte("lxllx"))}
	tree2, _ := transfer.NewMerkleTree([]common.Hash{utils.Sha3(lock2.AsBytes())})
	locksroot2 := tree2.MerkleRoot()
	//rmtr the mediatedtransfer i receive
	rmtr := encoding.NewMediatedTransfer(1, 1, tokenAddress, testChannel.MyAddress, big.NewInt(0), address1, locksroot2, lock2, address1, address2, utils.BigInt0)
	rmtr.Sign(privkey2, rmtr)
	err=testChannel.RegisterTransfer(blockNumber, rmtr)
	if err!= nil {
		t.Error("RegisterTransfer error")
		return
	}
	t.Log("after tr2 channel=", testChannel.String())
	assert.Equal(t,testChannel.OurState.AmountLocked(),amount1)
	assert.Equal(t,testChannel.PartnerState.AmountLocked(),amount2)
	/*
	try to remove hashlock now
	 */

	 //remove a not expired hashlock
	_,err=testChannel.CreateRemoveExpiredHashLockTransfer(smtr.HashLock,blockNumber)
	if err==nil{
		t.Error("cannot remove a hashlock which is not expired.")
		return
	}
	_,_,_,err=testChannel.PartnerState.TryRemoveExpiredHashLock(rmtr.HashLock,blockNumber)
	if err==nil{
		t.Error("cannot remove not expired hashlock")
		return
	}
	_,_,locksroot,err:= testChannel.PartnerState.TryRemoveExpiredHashLock(rmtr.HashLock,expiration)
	if err!=nil{
		t.Error("can remove a expired hashlock")
		return
	}
	removeTransferFromPartner:=encoding.NewRemoveExpiredHashlockTransfer(0, rmtr.Nonce+1, rmtr.Channel, rmtr.TransferAmount,locksroot,rmtr.HashLock)
	removeTransferFromPartner.Sign(privkey2,removeTransferFromPartner)
	err=testChannel.RegisterRemoveExpiredHashlockTransfer(removeTransferFromPartner,blockNumber)
	if err==nil{
		t.Error("can not register")
		return
	}
	err=testChannel.RegisterRemoveExpiredHashlockTransfer(removeTransferFromPartner,expiration)
	if err!=nil{
		t.Error("must be  removed ",err)
		return
	}
	removeTransferFromMe,err:=testChannel.CreateRemoveExpiredHashLockTransfer(smtr.HashLock,expiration)
	if err!=nil{
		t.Error("must be removed for a expired hashlock®")
		return
	}
	removeTransferFromMe.Sign(privkey1,removeTransferFromMe)
	err=testChannel.RegisterRemoveExpiredHashlockTransfer(removeTransferFromMe,expiration)
	if err!=nil{
		t.Errorf(" err register mine remove transfer ",err)
		return
	}
	assert.Equal(t,testChannel.OurState.BalanceProofState.LocksRoot,utils.EmptyHash)
	assert.Equal(t,testChannel.PartnerState.BalanceProofState.LocksRoot,utils.EmptyHash)
	assert.Equal(t,testChannel.OurState.BalanceProofState.IsBalanceProofValid(),true)
	assert.Equal(t,testChannel.PartnerState.BalanceProofState.IsBalanceProofValid(),true)
	assert.Equal(t,testChannel.OurState.AmountLocked(),utils.BigInt0)
	assert.Equal(t,testChannel.PartnerState.AmountLocked(),utils.BigInt0)
}