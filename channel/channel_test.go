package channel

import (
	"testing"

	"github.com/SmartMeshFoundation/Photon/params"

	"math/big"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	params.InitForUnitTest()
}

var big10 = big.NewInt(10)
var x = big.NewInt(0)
var testOpenBlockNumber int64 = 3

func TestEndState(t *testing.T) {
	bcs := rpc.MakeTestBlockChainService()
	address1 := bcs.NodeAddress
	address2 := utils.NewRandomAddress()
	channelIdentifier := &contracts.ChannelUniqueID{
		ChannelIdentifier: utils.NewRandomHash(),
		OpenBlockNumber:   testOpenBlockNumber,
	}

	var balance1 = big.NewInt(70)
	var balance2 = big.NewInt(110)
	lockSecret := utils.ShaSecret([]byte("test_end_state"))
	var lockAmount = big.NewInt(30)
	var lockExpiration int64 = 10
	lockHashlock := utils.ShaSecret(lockSecret[:])
	state1 := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	state2 := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	assert.EqualValues(t, state1.ContractBalance, balance1)
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), balance1)
	assert.EqualValues(t, state2.Balance(state1), balance2)
	assert.Equal(t, state1.IsLocked(lockHashlock), false)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)

	assert.Equal(t, state1.Tree.MerkleRoot(), utils.EmptyHash)
	assert.Equal(t, state2.Tree.MerkleRoot(), utils.EmptyHash)
	assert.EqualValues(t, state1.nonce(), 0)
	assert.EqualValues(t, state2.nonce(), 0)
	lock := &mtree.Lock{
		Expiration:     lockExpiration,
		Amount:         lockAmount,
		LockSecretHash: lockHashlock,
	}
	lockHash := utils.Sha3(lock.AsBytes())
	var transferedAmount = utils.BigInt0
	_, locksroot := state2.computeMerkleRootWith(lock)
	/*
		ChannelIdentifier   common.Hash
			OpenBlockNumber     int64    //open blocknumber 和 channelIdentifier 一起作为通道的唯一标识
			TransferAmount      *big.Int //The number has been transferred to the other party
			Locksroot           common.Hash

	*/
	bp := &encoding.BalanceProof{
		Nonce:             1,
		ChannelIdentifier: channelIdentifier.ChannelIdentifier,
		OpenBlockNumber:   testOpenBlockNumber,
		TransferAmount:    transferedAmount,
		Locksroot:         locksroot,
	}
	mtr := encoding.NewMediatedTransfer(bp, lock, utils.NewRandomAddress(), utils.NewRandomAddress(), utils.BigInt0, []common.Address{utils.NewRandomAddress()})
	mtr.Sign(bcs.PrivKey, mtr)
	err := state1.registerMediatedMessage(mtr)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, state1.ContractBalance, balance1)
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), balance1)
	assert.EqualValues(t, state2.Balance(state1), balance2)

	assert.EqualValues(t, state1.Distributable(state2), new(big.Int).Sub(balance1, lockAmount))
	assert.EqualValues(t, state2.Distributable(state1), balance2)

	assert.EqualValues(t, state1.amountLocked(), lockAmount)
	assert.EqualValues(t, state2.amountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), true)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.Tree.MerkleRoot(), lockHash)
	assert.Equal(t, state2.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.nonce(), 1)
	assert.EqualValues(t, state2.nonce(), 0)
	if state1.UpdateContractBalance(new(big.Int).Sub(balance1, big10)) != rerr.ErrChannelBalanceDecrease {
		t.Error(rerr.ErrChannelBalanceDecrease)
		return
	}
	assert.Equal(t, state1.UpdateContractBalance(new(big.Int).Add(balance1, big10)), nil)
	assert.EqualValues(t, state1.ContractBalance, new(big.Int).Add(balance1, big10))
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), new(big.Int).Add(balance1, big10))
	assert.EqualValues(t, state2.Balance(state1), balance2)
	x = new(big.Int).Sub(balance1, lockAmount)
	assert.EqualValues(t, state1.Distributable(state2), x.Add(x, big10))
	assert.EqualValues(t, state1.amountLocked(), lockAmount)
	assert.EqualValues(t, state2.amountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), true)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.Tree.MerkleRoot(), lockHash)
	assert.Equal(t, state2.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.nonce(), 1)
	assert.EqualValues(t, state2.nonce(), 0)

	err = state1.RegisterSecret(lockSecret)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, state1.ContractBalance, x.Add(balance1, big10))
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), x.Add(balance1, big10))
	assert.EqualValues(t, state2.Balance(state1), balance2)

	assert.EqualValues(t, state1.Distributable(state2), x.Sub(balance1, lockAmount).Add(x, big10))
	assert.EqualValues(t, state1.amountLocked(), lockAmount)
	assert.EqualValues(t, state2.amountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), false)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.Tree.MerkleRoot(), lockHash)
	assert.Equal(t, state2.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.nonce(), 1)
	assert.EqualValues(t, state2.nonce(), 0)

	secretMessage := encoding.NewUnlock(encoding.NewBalanceProof(2, x.Add(transferedAmount, lockAmount), utils.EmptyHash, channelIdentifier), lockSecret)
	secretMessage.Sign(bcs.PrivKey, secretMessage)
	state1.registerSecretMessage(secretMessage)

	assert.EqualValues(t, state1.ContractBalance, x.Add(balance1, big10))
	assert.EqualValues(t, state2.ContractBalance, balance2)
	assert.EqualValues(t, state1.Balance(state2), x.Add(balance1, big10).Sub(x, lockAmount))
	assert.EqualValues(t, state2.Balance(state1), x.Add(balance2, lockAmount))

	assert.EqualValues(t, state1.Distributable(state2), x.Sub(balance1, lockAmount).Add(x, big10))
	assert.EqualValues(t, state2.Distributable(state1), x.Add(balance2, lockAmount))
	assert.EqualValues(t, state1.amountLocked(), utils.BigInt0)
	assert.EqualValues(t, state2.amountLocked(), utils.BigInt0)

	assert.Equal(t, state1.IsLocked(lockHashlock), false)
	assert.Equal(t, state2.IsLocked(lockHashlock), false)
	assert.Equal(t, state1.Tree.MerkleRoot(), utils.EmptyHash)
	assert.Equal(t, state2.Tree.MerkleRoot(), utils.EmptyHash)

	assert.EqualValues(t, state1.nonce(), 2)
	assert.EqualValues(t, state2.nonce(), 0)
}
func makeExternState() *ExternalState {
	bcs := newTestBlockChainService()
	ch := common.HexToHash(os.Getenv("TOKEN_NETWORK"))
	//must provide a valid netting channel address
	tokenNetwork, _ := bcs.TokenNetwork(common.HexToAddress(os.Getenv("TOKEN_NETWORK")))
	return NewChannelExternalState(testFuncRegisterChannelForHashlock,
		tokenNetwork,
		&contracts.ChannelUniqueID{
			ChannelIdentifier: ch,
			OpenBlockNumber:   testOpenBlockNumber,
		},
		bcs.PrivKey, bcs.Client,
		channeltype.NewMockChannelDb(),
		0,
		bcs.NodeAddress, utils.NewRandomAddress())
}
func TestSenderCannotOverSpend(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	privkey1, address1 := utils.MakePrivateKeyAddress()
	address2 := utils.NewRandomAddress()
	var balance1 = big.NewInt(70)
	var balance2 = big.NewInt(110)
	revealTimeout := 5
	settleTimeout := 15
	var blockNumber int64 = 10
	ourState := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, revealTimeout, settleTimeout)
	amount := balance1
	expiration := blockNumber + int64(settleTimeout)
	sentMediatedTransfer0, err := testChannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount, expiration, utils.ShaSecret([]byte("test_locked_amount_cannot_be_spent")), []common.Address{})
	if err != nil {
		t.Error(err)
		return
	}
	sentMediatedTransfer0.Sign(privkey1, sentMediatedTransfer0)
	testChannel.RegisterTransfer(blockNumber, sentMediatedTransfer0)
	lock2 := &mtree.Lock{
		Expiration:     expiration,
		Amount:         amount,
		LockSecretHash: utils.ShaSecret([]byte("test_locked_amount_cannot_be_spent2")),
	}
	leaves := []*mtree.Lock{sentMediatedTransfer0.GetLock(), lock2}
	tree2 := mtree.NewMerkleTree(leaves)
	locksroot2 := tree2.MerkleRoot()
	bp := &encoding.BalanceProof{
		Nonce:             sentMediatedTransfer0.Nonce + 1,
		ChannelIdentifier: testChannel.ChannelIdentifier.ChannelIdentifier,
		OpenBlockNumber:   testChannel.ChannelIdentifier.OpenBlockNumber,
		TransferAmount:    utils.BigInt0,
		Locksroot:         locksroot2,
	}
	sentMediatedTransfer1 := encoding.NewMediatedTransfer(bp, lock2, address2, address1, utils.BigInt0, []common.Address{utils.NewRandomAddress()})
	sentMediatedTransfer1.Sign(privkey1, sentMediatedTransfer1)
	err = testChannel.RegisterTransfer(blockNumber, sentMediatedTransfer1)
	if err != rerr.ErrInsufficientBalance {
		t.Error(err)
		return
	}
}
func TestReceiverCannotSpendLockedAmount(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	privkey1, address1 := utils.MakePrivateKeyAddress()
	privkey2, address2 := utils.MakePrivateKeyAddress()
	var balance1 = big.NewInt(33)
	var balance2 = big.NewInt(11)
	revealTimeout := 7
	settleTimeout := 11
	var blockNumber int64 = 7
	ourState := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, revealTimeout, settleTimeout)
	amount1 := balance2
	expiration := blockNumber + int64(settleTimeout)
	receiveMediatedTransfer0, _ := testChannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount1, expiration, utils.ShaSecret([]byte("test_locked_amount_cannot_be_spent")), []common.Address{})
	receiveMediatedTransfer0.Sign(privkey2, receiveMediatedTransfer0)
	err := testChannel.RegisterTransfer(blockNumber, receiveMediatedTransfer0)
	if err != nil {
		t.Error(err)
	}
	t.Log("after tr1 channel=", testChannel.String())
	amount2 := x.Add(balance1, big.NewInt(1))
	lock2 := &mtree.Lock{
		Expiration:     expiration,
		Amount:         amount2,
		LockSecretHash: utils.ShaSecret([]byte("lxllx")),
	}
	tree2 := mtree.NewMerkleTree([]*mtree.Lock{lock2})
	locksroot2 := tree2.MerkleRoot()
	bp := &encoding.BalanceProof{
		Nonce:             1,
		ChannelIdentifier: testChannel.ChannelIdentifier.ChannelIdentifier,
		OpenBlockNumber:   testChannel.ChannelIdentifier.OpenBlockNumber,
		TransferAmount:    utils.BigInt0,
		Locksroot:         locksroot2,
	}
	sendMediatedTransfer0 := encoding.NewMediatedTransfer(bp, lock2, address2, address1, utils.BigInt0, []common.Address{utils.NewRandomAddress()})
	sendMediatedTransfer0.Sign(privkey1, sendMediatedTransfer0)
	if testChannel.RegisterTransfer(blockNumber, sendMediatedTransfer0) != rerr.ErrInsufficientBalance {
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
	ourState := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	externState := makeExternState()
	_, err := NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, 50, 49)
	if err == nil {
		t.Error("should failed")
	}
	for _, invalidValue := range []int{-1, 0, 1} {
		_, err = NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, invalidValue, settleTimeout)
		assert.NotEqual(t, err, nil)
		_, err = NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, reavealTimeout, invalidValue)
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
	ourState := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	externState := makeExternState()
	testchannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, reavealTimeout, settleTimeout)
	_, err := testchannel.CreateDirectTransfer(big.NewInt(-10))
	assert.NotEqual(t, err, nil)
	_, err = testchannel.CreateDirectTransfer(x.Add(balance1, big10))
	assert.NotEqual(t, err, nil)
	var amount1 = big.NewInt(10)
	directTransfer, _ := testchannel.CreateDirectTransfer(amount1)
	directTransfer.Sign(privkey1, directTransfer)
	testchannel.RegisterTransfer(blockNumber, directTransfer)

	assert.EqualValues(t, testchannel.ContractBalance(), balance1)
	assert.EqualValues(t, testchannel.Balance(), x.Sub(balance1, amount1))
	assert.EqualValues(t, testchannel.TransferAmount(), amount1)
	assert.EqualValues(t, testchannel.Distributable(), x.Sub(balance1, amount1))
	assert.EqualValues(t, testchannel.Outstanding(), utils.BigInt0)
	assert.EqualValues(t, testchannel.Locked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.OurState.amountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.PartnerState.amountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.GetNextNonce(), 2)

	secret := utils.ShaSecret([]byte("test_channel"))
	hashlock := utils.ShaSecret(secret[:])
	var amount2 = big.NewInt(10)
	expiration := blockNumber + int64(settleTimeout) - 5
	mediatedTransfer, _ := testchannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount2, expiration, hashlock, []common.Address{})
	mediatedTransfer.Sign(privkey1, mediatedTransfer)
	testchannel.RegisterTransfer(blockNumber, mediatedTransfer)

	assert.EqualValues(t, testchannel.ContractBalance(), balance1)
	assert.EqualValues(t, testchannel.Balance(), x.Sub(balance1, amount1))
	assert.EqualValues(t, testchannel.TransferAmount(), amount1)
	assert.EqualValues(t, testchannel.Distributable(), x.Sub(balance1, amount1).Sub(x, amount2))
	assert.EqualValues(t, testchannel.Outstanding(), utils.BigInt0)
	assert.EqualValues(t, testchannel.Locked(), amount2)
	assert.EqualValues(t, testchannel.OurState.amountLocked(), amount2)
	assert.EqualValues(t, testchannel.PartnerState.amountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.GetNextNonce(), 3)

	err = testchannel.RegisterSecret(secret)
	if err != nil {
		t.Error(err)
		return
	}
	secretMessage, err := testchannel.CreateUnlock(utils.ShaSecret(secret[:]))
	if err != nil {
		t.Error(err)
		return
	}
	secretMessage.Sign(privkey1, secretMessage)
	log.Info(fmt.Sprintf("secret message=%s", utils.StringInterface(secretMessage, 4)))
	log.Info(fmt.Sprintf("bofore reg sec proof=%s", utils.StringInterface(testchannel.OurState.BalanceProofState, 2)))
	err = testchannel.RegisterTransfer(blockNumber, secretMessage)
	if err != nil {
		t.Error(err)
	}
	log.Info(fmt.Sprintf("after reg sec proof=%s", utils.StringInterface(testchannel.OurState.BalanceProofState, 2)))
	assert.EqualValues(t, testchannel.ContractBalance(), balance1)
	assert.EqualValues(t, testchannel.Balance(), x.Sub(balance1, amount1).Sub(x, amount2))
	assert.EqualValues(t, testchannel.TransferAmount(), x.Add(amount1, amount2))
	assert.EqualValues(t, testchannel.Distributable(), x.Sub(balance1, amount1).Sub(x, amount2))
	assert.EqualValues(t, testchannel.Outstanding(), utils.BigInt0)
	assert.EqualValues(t, testchannel.Locked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.OurState.amountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.OurState.amountLocked(), utils.BigInt0)
	assert.EqualValues(t, testchannel.PartnerState.amountLocked(), utils.BigInt0)
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
	ourState := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	externState := makeExternState()
	tch, _ := NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, reavealTimeout, settleTimeout)
	previousNonce := tch.GetNextNonce()
	previousTransfered := tch.TransferAmount()
	var amount = big.NewInt(7)
	for i := 0; i < 10; i++ {
		directTransfer, _ := tch.CreateDirectTransfer(amount)
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
	ourState := NewChannelEndState(externState1.MyAddress, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(externState2.MyAddress, balance2, nil, mtree.EmptyTree)

	testChannel, _ := NewChannel(ourState, partnerState, externState1, tokenAddress, &externState1.ChannelIdentifier, revealTimeout, settleTimeout)

	ourState = NewChannelEndState(externState1.MyAddress, balance1, nil, mtree.EmptyTree)
	partnerState = NewChannelEndState(externState2.MyAddress, balance2, nil, mtree.EmptyTree)
	testChannel2, _ := NewChannel(partnerState, ourState, externState2, tokenAddress, &externState2.ChannelIdentifier, revealTimeout, settleTimeout)
	return testChannel, testChannel2
}

/*
Assert that `channel0` has a correct `partner_state` to represent
    `channel1` and vice-versa.
*/
func assertMirror(ch0, ch1 *Channel, t *testing.T) {
	unclaimed0 := ch0.OurState.Tree.MerkleRoot()
	unclaimed1 := ch1.PartnerState.Tree.MerkleRoot()
	assert.EqualValues(t, unclaimed0, unclaimed1)

	assert.EqualValues(t, ch0.OurState.amountLocked(), ch1.PartnerState.amountLocked())
	assert.EqualValues(t, ch0.TransferAmount(), ch1.PartnerState.TransferAmount())
	balance0 := ch0.OurState.Balance(ch0.PartnerState)
	balance1 := ch1.PartnerState.Balance(ch1.OurState)
	assert.EqualValues(t, balance0, balance1)

	assert.EqualValues(t, ch0.Distributable(), ch0.OurState.Distributable(ch0.PartnerState))
	assert.EqualValues(t, ch0.Distributable(), ch1.PartnerState.Distributable(ch1.OurState))

	unclaimed0 = ch1.OurState.Tree.MerkleRoot()
	unclaimed1 = ch0.PartnerState.Tree.MerkleRoot()
	assert.EqualValues(t, unclaimed0, unclaimed1)

	assert.EqualValues(t, ch1.OurState.amountLocked(), ch0.PartnerState.amountLocked())
	assert.EqualValues(t, ch1.TransferAmount(), ch0.PartnerState.TransferAmount())
	balance0 = ch1.OurState.Balance(ch1.PartnerState)
	balance1 = ch0.PartnerState.Balance(ch0.OurState)
	assert.EqualValues(t, balance0, balance1)

	assert.EqualValues(t, ch1.Distributable(), ch1.OurState.Distributable(ch1.PartnerState))
	assert.EqualValues(t, ch1.Distributable(), ch0.PartnerState.Distributable(ch0.OurState))
}

//Assert the locks created from `from_channel`.
func assertLocked(ch *Channel, pendingLocks []*mtree.Lock, t *testing.T) {
	var root common.Hash
	if pendingLocks != nil {
		tree := mtree.NewMerkleTree(pendingLocks)
		root = tree.MerkleRoot()
	}
	assert.EqualValues(t, len(ch.OurState.Lock2PendingLocks), len(pendingLocks))
	assert.EqualValues(t, ch.OurState.Tree.MerkleRoot(), root)
	var sum = big.NewInt(0)
	for _, lock := range pendingLocks {
		sum.Add(sum, lock.Amount)
	}
	assert.EqualValues(t, ch.OurState.amountLocked(), sum)
	for _, lock := range pendingLocks {
		assert.Equal(t, ch.OurState.IsLocked(lock.LockSecretHash), true)
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
	assert.EqualValues(t, ch.PartnerState.amountLocked(), outstanding)
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
func assertSyncedChannels(ch0 *Channel, balance0 *big.Int, outstandingLocks0 []*mtree.Lock, ch1 *Channel, balance1 *big.Int, outstandingLocks1 []*mtree.Lock, t *testing.T) {
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
	var unclaimedLocks []*mtree.Lock
	var transfersList []*encoding.MediatedTransfer
	var transfersClaimed []bool
	var transfersAmount []*big.Int
	var transfersSecret []common.Hash
	for i := 1; i <= ArgNumberOfTransfers; i++ {
		transfersAmount = append(transfersAmount, big.NewInt(int64(i)))
		transfersSecret = append(transfersSecret, utils.ShaSecret(utils.Random(32)))
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
		var mtr *encoding.MediatedTransfer
		mtr, err = ch0.CreateMediatedTransfer(ch0.OurState.Address, ch1.OurState.Address, utils.BigInt0, amount, expiration, utils.ShaSecret(secret[:]), []common.Address{})
		assert.Equal(t, err, nil)
		mtr.Sign(ch0.ExternState.privKey, mtr)
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
			err = ch0.RegisterSecret(secret)
			if err != nil {
				t.Error(err)
				return
			}
			//synchronized claiming
			secretMessage, err := ch0.CreateUnlock(utils.ShaSecret(secret[:]))
			if err != nil {
				t.Error(err)
				return
			}
			secretMessage.Sign(ch0.ExternState.privKey, secretMessage)
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
	directTransfer, err := ch0.CreateDirectTransfer(amount)
	assert.Equal(t, err, nil)
	directTransfer.Sign(ch0.ExternState.privKey, directTransfer)
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
	secret := utils.ShaSecret([]byte("secret"))
	hashlock := utils.ShaSecret(secret[:])
	transfer1, err := ch0.CreateMediatedTransfer(ch0.OurState.Address, ch1.OurState.Address, utils.BigInt0, amount, expiration, hashlock, []common.Address{})
	assert.Equal(t, err, nil)
	transfer1.Sign(ch0.ExternState.privKey, transfer1)
	err = ch0.RegisterTransfer(blockNumber, transfer1)
	assert.Equal(t, err, nil)
	err = ch1.RegisterTransfer(blockNumber, transfer1)
	assert.Equal(t, err, nil)
	assertSyncedChannels(ch0, balance0, nil,
		ch1, balance1, []*mtree.Lock{transfer1.GetLock()}, t)
	// handcrafted transfer because channel.create_transfer won't create it
	transfer2 := encoding.NewDirectTransfer(encoding.NewBalanceProof(ch0.GetNextNonce(), x.Add(ch1.Balance(), balance0).Add(x, amount), ch0.PartnerState.Tree.MerkleRoot(), &ch0.ChannelIdentifier))
	transfer2.Sign(ch0.ExternState.privKey, transfer2)
	err = ch0.RegisterTransfer(blockNumber, transfer2)
	assert.Equal(t, err != nil, true)
	err = ch1.RegisterTransfer(blockNumber, transfer2)
	assert.Equal(t, err != nil, true)
	assertSyncedChannels(ch0, balance0, nil,
		ch1, balance1, []*mtree.Lock{transfer1.GetLock()}, t)
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
	_, address1 := utils.MakePrivateKeyAddress()
	privkey2, address2 := utils.MakePrivateKeyAddress()
	var balance1 = big.NewInt(33)
	var balance2 = big.NewInt(11)
	revealTimeout := 7
	settleTimeout := 11
	var blockNumber int64 = 7
	ourState := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, revealTimeout, settleTimeout)
	lock := &mtree.Lock{Expiration: blockNumber + int64(settleTimeout),
		Amount:         big.NewInt(1),
		LockSecretHash: utils.EmptyHash,
	}
	bp := &encoding.BalanceProof{
		Nonce:             testChannel.GetNextNonce(),
		ChannelIdentifier: testChannel.ChannelIdentifier.ChannelIdentifier,
		OpenBlockNumber:   testChannel.ChannelIdentifier.OpenBlockNumber,
		TransferAmount:    big.NewInt(1),
		Locksroot:         utils.Sha3(lock.AsBytes()),
	}
	transfer := encoding.NewMediatedTransfer(bp, lock, utils.EmptyAddress, utils.EmptyAddress, utils.BigInt0, []common.Address{utils.NewRandomAddress()})
	transfer.Sign(privkey2, transfer)
	err := testChannel.RegisterTransfer(blockNumber+int64(settleTimeout)+1, transfer)
	assert.Equal(t, err, nil)
}

func TestRemoveExpiredHashlock(t *testing.T) {
	tokenAddress := utils.NewRandomAddress()
	privkey1, address1 := utils.MakePrivateKeyAddress()
	privkey2, address2 := utils.MakePrivateKeyAddress()
	var balance1 = big.NewInt(33)
	var balance2 = big.NewInt(11)
	revealTimeout := 7
	settleTimeout := 11
	var blockNumber int64 = 7
	ourState := NewChannelEndState(address1, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(address2, balance2, nil, mtree.EmptyTree)
	externState := makeExternState()
	testChannel, _ := NewChannel(ourState, partnerState, externState, tokenAddress, &externState.ChannelIdentifier, revealTimeout, settleTimeout)
	amount1 := balance2
	expiration := blockNumber + int64(settleTimeout)
	//smtr: the mediated transfer i sent out
	smtr, _ := testChannel.CreateMediatedTransfer(address1, address2, utils.BigInt0, amount1, expiration, utils.ShaSecret([]byte("test_locked_amount_cannot_be_spent")), []common.Address{})
	smtr.Sign(privkey1, smtr)
	err := testChannel.RegisterTransfer(blockNumber, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("after tr1 channel=", testChannel.String())
	amount2 := balance2
	lock2 := &mtree.Lock{
		Expiration:     expiration,
		Amount:         amount2,
		LockSecretHash: utils.ShaSecret([]byte("lxllx")),
	}
	tree2 := mtree.NewMerkleTree([]*mtree.Lock{lock2})
	locksroot2 := tree2.MerkleRoot()
	//rmtr the mediatedtransfer i receive
	bp := &encoding.BalanceProof{
		Nonce:             1,
		ChannelIdentifier: testChannel.ChannelIdentifier.ChannelIdentifier,
		OpenBlockNumber:   testChannel.ChannelIdentifier.OpenBlockNumber,
		TransferAmount:    big.NewInt(0),
		Locksroot:         locksroot2,
	}
	rmtr := encoding.NewMediatedTransfer(bp, lock2, address1, address2, utils.BigInt0, []common.Address{utils.NewRandomAddress()})
	rmtr.Sign(privkey2, rmtr)
	err = testChannel.RegisterTransfer(blockNumber, rmtr)
	if err != nil {
		t.Error("RegisterTransfer error")
		return
	}
	t.Log("after tr2 channel=", testChannel.String())
	assert.Equal(t, testChannel.OurState.amountLocked(), amount1)
	assert.Equal(t, testChannel.PartnerState.amountLocked(), amount2)
	/*
		try to remove hashlock now
	*/

	//remove a not expired hashlock
	_, err = testChannel.CreateRemoveExpiredHashLockTransfer(smtr.LockSecretHash, blockNumber)
	if err == nil {
		t.Error("cannot remove a hashlock which is not expired.")
		return
	}
	_, _, _, err = testChannel.PartnerState.TryRemoveHashLock(rmtr.LockSecretHash, blockNumber, true)
	if err == nil {
		t.Error("cannot remove not expired hashlock")
		return
	}
	_, _, locksroot, err := testChannel.PartnerState.TryRemoveHashLock(rmtr.LockSecretHash, expiration+params.Cfg.ForkConfirmNumber+1, true)
	if err != nil {
		t.Errorf("can remove a expired hashlock err=%s", err)
		return
	}
	bp = &encoding.BalanceProof{
		Nonce:             rmtr.Nonce + 1,
		ChannelIdentifier: rmtr.ChannelIdentifier,
		OpenBlockNumber:   rmtr.OpenBlockNumber,
		TransferAmount:    rmtr.TransferAmount,
		Locksroot:         locksroot,
	}
	removeTransferFromPartner := encoding.NewRemoveExpiredHashlockTransfer(bp, rmtr.LockSecretHash)
	removeTransferFromPartner.Sign(privkey2, removeTransferFromPartner)
	err = testChannel.RegisterRemoveExpiredHashlockTransfer(removeTransferFromPartner, blockNumber)
	if err == nil {
		t.Error("can not register")
		return
	}
	err = testChannel.RegisterRemoveExpiredHashlockTransfer(removeTransferFromPartner, expiration+params.Cfg.ForkConfirmNumber)
	if err != nil {
		t.Error("must be  removed ", err)
		return
	}
	removeTransferFromMe, err := testChannel.CreateRemoveExpiredHashLockTransfer(smtr.LockSecretHash, expiration+params.Cfg.ForkConfirmNumber)
	if err != nil {
		t.Error("must be removed for a expired hashlock®")
		return
	}
	removeTransferFromMe.Sign(privkey1, removeTransferFromMe)
	err = testChannel.RegisterRemoveExpiredHashlockTransfer(removeTransferFromMe, expiration+params.Cfg.ForkConfirmNumber)
	if err != nil {
		t.Errorf(" err register mine remove transfer %s", err)
		return
	}
	assert.Equal(t, testChannel.OurState.BalanceProofState.LocksRoot, utils.EmptyHash)
	assert.Equal(t, testChannel.PartnerState.BalanceProofState.LocksRoot, utils.EmptyHash)
	assert.Equal(t, testChannel.OurState.BalanceProofState.IsBalanceProofValid(), true)
	assert.Equal(t, testChannel.PartnerState.BalanceProofState.IsBalanceProofValid(), true)
	assert.Equal(t, testChannel.OurState.amountLocked(), utils.BigInt0)
	assert.Equal(t, testChannel.PartnerState.amountLocked(), utils.BigInt0)
}

func TestChannel_RegisterAnnounceDisposedTransferResponse(t *testing.T) {
	var blockNumber int64 = 7
	ch0, ch1 := makePairChannel()
	expiration := blockNumber + int64(ch0.SettleTimeout)
	lockSecretHash := utils.ShaSecret([]byte("123"))
	smtr, _ := ch0.CreateMediatedTransfer(ch0.OurState.Address, ch0.PartnerState.Address, utils.BigInt0, big.NewInt(1), expiration, lockSecretHash, []common.Address{})
	err := smtr.Sign(ch0.ExternState.privKey, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	err = ch0.RegisterTransfer(blockNumber, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("after tr1 channel=", ch0.String())
	err = ch1.RegisterTransfer(blockNumber, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	assertMirror(ch0, ch1, t)
	req, err := ch1.CreateAnnouceDisposed(lockSecretHash, blockNumber, rerr.ErrNoAvailabeRoute)
	if err != nil {
		t.Error(err)
		return
	}
	err = req.Sign(ch1.ExternState.privKey, req)
	if err != nil {
		t.Error(err)
		return
	}
	err = ch0.RegisterAnnouceDisposed(req)
	if err != nil {
		t.Error(err)
		return
	}
	err = ch1.RegisterAnnouceDisposed(req)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := ch0.CreateAnnounceDisposedResponse(lockSecretHash, blockNumber)
	if err != nil {
		t.Error(err)
		return
	}
	err = res.Sign(ch0.ExternState.privKey, res)
	if err != nil {
		t.Error(err)
		return
	}
	err = ch0.RegisterAnnounceDisposedResponse(res, blockNumber)
	if err != nil {
		t.Error(err)
		return
	}
	err = ch1.RegisterAnnounceDisposedResponse(res, blockNumber)
	if err != nil {
		t.Error(err)
		return
	}
	assertMirror(ch0, ch1, t)
}

func TestChannel_RegisterWithdrawRequest(t *testing.T) {
	//var blockNumber int64 = 7
	//ch0, ch1 := makePairChannel()
	//expiration := blockNumber + int64(ch0.SettleTimeout)
	//secret := utils.ShaSecret([]byte("123"))
	//lockSecretHash := utils.ShaSecret(secret[:])
	//smtr, _ := ch0.CreateMediatedTransfer(ch0.OurState.Address, ch0.PartnerState.Address, utils.BigInt0, big.NewInt(1), expiration, lockSecretHash)
	//err := smtr.Sign(ch0.ExternState.privKey, smtr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//err = ch0.RegisterTransfer(blockNumber, smtr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Log("after tr1 channel=", ch0.String())
	//err = ch1.RegisterTransfer(blockNumber, smtr)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//_, err = ch0.CreateWithdrawRequest(big.NewInt(1))
	//if err == nil {
	//	t.Error("have lock doesn't allow withdraw")
	//	return
	//}
	//_, err = ch1.CreateWithdrawRequest(big.NewInt(1))
	//if err == nil {
	//	t.Error("have lock doesn't allow withdraw")
	//	return
	//}
	//err = ch0.RegisterSecret(secret)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//unlock, err := ch0.CreateUnlock(utils.ShaSecret(secret[:]))
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//unlock.Sign(ch0.ExternState.privKey, unlock)
	//err = ch0.RegisterTransfer(blockNumber, unlock)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//err = ch1.RegisterTransfer(blockNumber, unlock)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//assert.EqualValues(t, ch0.CanTransfer(), true)
	//assert.EqualValues(t, ch0.CanContinueTransfer(), true)
	//assert.EqualValues(t, ch1.CanTransfer(), true)
	//assert.EqualValues(t, ch1.CanContinueTransfer(), true)
	//
	//req, err := ch0.CreateWithdrawRequest(big.NewInt(1))
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//log.Trace(fmt.Sprintf("ch0=%s", utils.StringInterface(NewChannelSerialization(ch0), 3)))
	//log.Trace(fmt.Sprintf("req=%s", req))
	//req.Sign(ch1.ExternState.privKey, req)
	//err = ch0.RegisterWithdrawRequest(req)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//req.Sign(ch0.ExternState.privKey, req)
	//err = ch1.RegisterWithdrawRequest(req)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//assert.EqualValues(t, ch0.CanTransfer(), false)
	//assert.EqualValues(t, ch0.CanContinueTransfer(), false)
	//assert.EqualValues(t, ch1.CanTransfer(), false)
	//assert.EqualValues(t, ch1.CanContinueTransfer(), false)
	////目前 channel 并不验证自己是否发出了 withdrawRequest,这些请求应该保存在数据库中,由更高层验证.
	//res, err := ch1.CreateWithdrawResponse(req)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//res.Sign(ch1.ExternState.privKey, res)
	//err = ch0.RegisterWithdrawResponse(res)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//err = ch1.RegisterWithdrawResponse(res)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
}

func TestChannel_RegisterCooperativeSettleRequest(t *testing.T) {
	var blockNumber int64 = 7
	ch0, ch1 := makePairChannel()
	expiration := blockNumber + int64(ch0.SettleTimeout)
	secret := utils.ShaSecret([]byte("123"))
	lockSecretHash := utils.ShaSecret(secret[:])
	smtr, _ := ch0.CreateMediatedTransfer(ch0.OurState.Address, ch0.PartnerState.Address, utils.BigInt0, big.NewInt(1), expiration, lockSecretHash, []common.Address{})
	err := smtr.Sign(ch0.ExternState.privKey, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	err = ch0.RegisterTransfer(blockNumber, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("after tr1 channel=", ch0.String())
	err = ch1.RegisterTransfer(blockNumber, smtr)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = ch0.CreateCooperativeSettleRequest()
	if err == nil {
		t.Error("have lock doesn't allow withdraw")
		return
	}
	_, err = ch1.CreateCooperativeSettleRequest()
	if err == nil {
		t.Error("have lock doesn't allow withdraw")
		return
	}
	err = ch0.RegisterSecret(secret)
	if err != nil {
		t.Error(err)
		return
	}
	unlock, err := ch0.CreateUnlock(utils.ShaSecret(secret[:]))
	if err != nil {
		t.Error(err)
		return
	}
	unlock.Sign(ch0.ExternState.privKey, unlock)
	err = ch0.RegisterTransfer(blockNumber, unlock)
	if err != nil {
		t.Error(err)
		return
	}
	err = ch1.RegisterTransfer(blockNumber, unlock)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, ch0.CanTransfer(), true)
	assert.EqualValues(t, ch0.CanContinueTransfer(), true)
	assert.EqualValues(t, ch1.CanTransfer(), true)
	assert.EqualValues(t, ch1.CanContinueTransfer(), true)

	req, err := ch0.CreateCooperativeSettleRequest()
	if err != nil {
		t.Error(err)
		return
	}
	//log.Trace(fmt.Sprintf("ch0=%s", utils.StringInterface(NewChannelSerialization(ch0), 3)))
	log.Trace(fmt.Sprintf("req=%s", req))
	req.Sign(ch0.ExternState.privKey, req)
	//err = ch0.RegisterCooperativeSettleRequest(req)
	ch0.State = channeltype.StateCooprativeSettle
	if err != nil {
		t.Error(err)
		return
	}
	err = ch1.RegisterCooperativeSettleRequest(req)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, ch0.CanTransfer(), false)
	assert.EqualValues(t, ch0.CanContinueTransfer(), false)
	assert.EqualValues(t, ch1.CanTransfer(), false)
	assert.EqualValues(t, ch1.CanContinueTransfer(), false)
	//目前 channel 并不验证自己是否发出了 withdrawRequest,这些请求应该保存在数据库中,由更高层验证.
	// Currently, channel can't verify if he self sends out withdrawrequest,
	// these requests are backed up in local database, which needs to be verified by upper layer.
	res, err := ch1.CreateCooperativeSettleResponse(req)
	if err != nil {
		t.Error(err)
		return
	}
	res.Sign(ch1.ExternState.privKey, res)
	err = ch0.RegisterCooperativeSettleResponse(res)
	if err != nil {
		t.Error(err)
		return
	}
	// 目前在通道状态中区分了自己settle还是对方settle,所这里不能双方都注册request和response,只能一方注册request,一方注册response
	//err = ch1.RegisterCooperativeSettleResponse(res)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
}
