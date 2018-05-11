package channel

import (
	"errors"

	"fmt"

	"encoding/gob"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

type PendingLock struct {
	Lock       *encoding.Lock
	LockHashed common.Hash
}
type UnlockPartialProof struct {
	Lock       *encoding.Lock
	LockHashed common.Hash
	Secret     common.Hash
}
type UnlockProof struct {
	MerkleProof []common.Hash
	LockEncoded []byte
	Secret      common.Hash
}

func init() {
	gob.Register(&PendingLock{})
	gob.Register(&UnlockPartialProof{})
	//make sure don't save this data
	//gob.Register(&UnlockProof{})
}

type InvalidLocksRootError struct {
	ExpectedLocksroot common.Hash
	GotLocksroot      common.Hash
}

func (this *InvalidLocksRootError) Error() string {
	return fmt.Sprintf("Locksroot mismatch. Expected %s but get %s",
		utils.Pex(this.ExpectedLocksroot[:]), utils.Pex(this.GotLocksroot[:]))
}

var errBalanceDecrease = errors.New("contract_balance cannot decrease")
var errUnknownLock = errors.New("'unknown lock")
var errTransferAmountMismatch = errors.New("transfer amount mismatch")

/*
Tracks the state of one of the participants in a channel
*/
type ChannelEndState struct {
	Address             common.Address
	ContractBalance     *big.Int                    //lock protect race codition with raidenapi
	Lock2PendingLocks   map[common.Hash]PendingLock //the lock I have sent
	Lock2UnclaimedLocks map[common.Hash]UnlockPartialProof
	TreeState           *transfer.MerkleTreeState
	BalanceProofState   *transfer.BalanceProofState //race codition with raidenapi
}

func NewChannelEndState(participantAddress common.Address, participantBalance *big.Int,
	balanceProof *transfer.BalanceProofState, tree *transfer.MerkleTreeState) *ChannelEndState {
	c := &ChannelEndState{
		Address:             participantAddress,
		ContractBalance:     participantBalance,
		TreeState:           tree,
		BalanceProofState:   balanceProof,
		Lock2PendingLocks:   make(map[common.Hash]PendingLock),
		Lock2UnclaimedLocks: make(map[common.Hash]UnlockPartialProof),
	}
	return c
}

//how many tokens I have sent to  partner.
func (this *ChannelEndState) TransferAmount() *big.Int {
	if this.BalanceProofState != nil {
		return this.BalanceProofState.TransferAmount
	}
	return big.NewInt(0)
}
func (this *ChannelEndState) AmountLocked() *big.Int {
	sum := big.NewInt(0)
	for _, v := range this.Lock2PendingLocks {
		sum = sum.Add(sum, v.Lock.Amount)
	}
	for _, v := range this.Lock2UnclaimedLocks {
		sum = sum.Add(sum, v.Lock.Amount)
	}
	return sum
}

func (this *ChannelEndState) Nonce() int64 {
	if this.BalanceProofState != nil {
		return this.BalanceProofState.Nonce
	}
	return 0
}

func (this *ChannelEndState) Balance(counterpart *ChannelEndState) *big.Int {
	x := new(big.Int).Sub(this.ContractBalance, this.TransferAmount())
	x.Add(x, counterpart.TransferAmount())
	return x
}

func (this *ChannelEndState) Distributable(counterpart *ChannelEndState) *big.Int {
	return new(big.Int).Sub(this.Balance(counterpart), this.AmountLocked())
}

//True if the `hashlock` corresponds to a known lock.
func (this *ChannelEndState) IsKnown(hashlock common.Hash) bool {
	_, ok := this.Lock2PendingLocks[hashlock]
	if ok {
		return ok
	}
	_, ok = this.Lock2UnclaimedLocks[hashlock]
	return ok
}

//True if the `hashlock` is known and the correspoding secret is not.
func (this *ChannelEndState) IsLocked(hashlock common.Hash) bool {
	_, ok := this.Lock2PendingLocks[hashlock]
	return ok
}

/*
Update the contract balance, it must always increase.

return error If the `contract_balance` is smaller than the current
           balance.
*/
func (this *ChannelEndState) UpdateContractBalance(balance *big.Int) error {
	if balance.Cmp(this.ContractBalance) < 0 {
		return errBalanceDecrease
	}
	this.ContractBalance = new(big.Int).Set(balance)
	return nil
}

func (this *ChannelEndState) GetLockByHashlock(hashlock common.Hash) *encoding.Lock {
	lock, ok := this.Lock2PendingLocks[hashlock]
	if ok {
		return lock.Lock
	}
	plock, ok := this.Lock2UnclaimedLocks[hashlock]
	if ok {
		return plock.Lock
	}
	return nil
}

/*
Compute the resulting merkle root if the lock `include` is added in
       the tree.
*/
func (this *ChannelEndState) ComputeMerkleRootWith(include *encoding.Lock) (tree *transfer.Merkletree, hash common.Hash) {
	if !this.IsKnown(include.HashLock) {
		leaves := make([]common.Hash, len(this.TreeState.Tree.Layers[transfer.LayerLeaves]))
		copy(leaves, this.TreeState.Tree.Layers[transfer.LayerLeaves])
		includeHash := utils.Sha3(include.AsBytes())
		leaves = append(leaves, includeHash)
		tree, err := transfer.NewMerkleTree(leaves)
		if err != nil {
			log.Error(fmt.Sprintf("NewMerkleTree err %s", err))
		}
		return tree, tree.MerkleRoot()
	}
	return nil, this.TreeState.Tree.MerkleRoot()
}

func removeHash(leaves []common.Hash, hash common.Hash) []common.Hash {
	i := -1
	for j := 0; j < len(leaves); j++ {
		if leaves[j] == hash {
			i = j
			break
		}
	}
	if i >= 0 {
		leaves = append(leaves[:i], leaves[i+1:]...)
	}
	return leaves
}

/*
 Compute the resulting merkle root if the lock `without` is exclude from the tree
*/
func (this *ChannelEndState) ComputeMerkleRootWithout(without *encoding.Lock) (*transfer.Merkletree, common.Hash, error) {
	if !this.IsKnown(without.HashLock) {
		return nil, utils.EmptyHash, errUnknownLock
	}
	leaves := make([]common.Hash, len(this.TreeState.Tree.Layers[transfer.LayerLeaves]))
	copy(leaves, this.TreeState.Tree.Layers[transfer.LayerLeaves])
	withoutHash := utils.Sha3(without.AsBytes())
	leaves = removeHash(leaves, withoutHash)
	if len(leaves) > 0 {
		tree, err := transfer.NewMerkleTree(leaves)
		if err != nil {
			return nil, utils.EmptyHash, err
		}
		return tree, tree.MerkleRoot(), nil
	}
	return nil, utils.EmptyHash, nil
}

/*
    Api design: using specialized methods to force the user to register the
    transfer and the lock in a single step
	Register the latest known transfer.

       The sender needs to use this method before sending a locked transfer,
       otherwise the calculate locksroot of the transfer message will be
       invalid and the transfer will be rejected by the partner. Since the
       sender wants the transfer to be accepted by the receiver otherwise the
       transfer won't proceed and the sender won't receive their fee.

       The receiver needs to use this method to update the container with a
       _valid_ transfer, otherwise the locksroot will not contain the pending
       transfer. The receiver needs to ensure that the merkle root has the
       hashlock included, otherwise it won't be able to claim it.

       Args:
           transfer (LockedTransfer): The transfer to be added.

       Raises:
           InvalidLocksRoot: If the merkleroot of `locked_transfer` does not
           match with the expected value.

           ValueError: If the transfer contains a lock that was registered
           previously.
//Calculate the banlanceproof locksroot position before sending
*/
func (this *ChannelEndState) RegisterLockedTransfer(lockedTransfer encoding.EnvelopMessager) error {
	if !encoding.IsLockedTransfer(lockedTransfer) {
		return errors.New("not a locked lockedTransfer")
	}
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(lockedTransfer)
	mtranfer := encoding.GetMtrFromLockedTransfer(lockedTransfer)
	lock := mtranfer.GetLock()
	if this.IsKnown(lock.HashLock) {
		return errors.New("hashlock is already registered")
	}
	newtree, locksroot := this.ComputeMerkleRootWith(lock)
	lockhashed := utils.Sha3(lock.AsBytes())
	if balanceProof.LocksRoot != locksroot {
		return &InvalidLocksRootError{locksroot, balanceProof.LocksRoot}
	}
	this.Lock2PendingLocks[lock.HashLock] = PendingLock{lock, lockhashed}
	this.BalanceProofState = balanceProof
	this.TreeState = transfer.NewMerkleTreeState(newtree)
	return nil
}

/*
Register a direct_transfer.

       Raises:
           InvalidLocksRoot: If the merkleroot of `direct_transfer` does not
           match the current value.
*/
func (this *ChannelEndState) RegisterDirectTransfer(directTransfer *encoding.DirectTransfer) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(directTransfer)
	if balanceProof.LocksRoot != this.TreeState.Tree.MerkleRoot() {
		return &InvalidLocksRootError{this.TreeState.Tree.MerkleRoot(), balanceProof.LocksRoot}
	}
	this.BalanceProofState = balanceProof
	return nil
}
func (this *ChannelEndState) RegisterRemoveExpiredHashlockTransfer(removeExpiredHashlockTransfer *encoding.RemoveExpiredHashlockTransfer) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(removeExpiredHashlockTransfer)
	if balanceProof.TransferAmount.Cmp(this.TransferAmount()) != 0 {
		return errTransferAmountMismatch
	}
	this.BalanceProofState = balanceProof
	delete(this.Lock2PendingLocks, removeExpiredHashlockTransfer.HashLock)
	delete(this.Lock2UnclaimedLocks, removeExpiredHashlockTransfer.HashLock)
	return nil
}
func (this *ChannelEndState) RegisterSecretMessage(secret *encoding.Secret) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(secret)
	hashlock := utils.Sha3(secret.Secret[:])
	pendingLock, ok := this.Lock2PendingLocks[hashlock]
	var lock *encoding.Lock
	if ok {
		lock = pendingLock.Lock
	} else {
		unclaimedLock, ok := this.Lock2UnclaimedLocks[hashlock]
		if ok {
			lock = unclaimedLock.Lock
		}
	}
	//if !this.IsKnown(lock.HashLock) { // has corrupted because of lock is nil
	//	return errors.New("hashlock is not registered")
	//}
	newtree, newLocksroot, err := this.ComputeMerkleRootWithout(lock)
	if err != nil {
		return err
	}
	if newtree == nil {
		newtree, err = transfer.NewMerkleTree(nil)
		if err != nil {
			return err
		}
	}
	if balanceProof.LocksRoot != newLocksroot {
		return &InvalidLocksRootError{newLocksroot, balanceProof.LocksRoot}
	}
	delete(this.Lock2PendingLocks, lock.HashLock)
	delete(this.Lock2UnclaimedLocks, lock.HashLock)
	this.TreeState = transfer.NewMerkleTreeState(newtree)
	this.BalanceProofState = balanceProof
	return nil
}

/*
try to remomve a expired hashlock
*/
func (this *ChannelEndState) TryRemoveExpiredHashLock(hashlock common.Hash, blockNumber int64) (lock *encoding.Lock, newtree *transfer.Merkletree, newlocksroot common.Hash, err error) {
	if !this.IsKnown(hashlock) {
		err = fmt.Errorf("channel %s donesn't know hashlock %s, cannot remove", utils.APex(this.Address), utils.HPex(hashlock))
		return
	}
	pendingLock, ok := this.Lock2PendingLocks[hashlock]
	if ok {
		lock = pendingLock.Lock
	} else {
		unclaimedLock, ok := this.Lock2UnclaimedLocks[hashlock]
		if ok {
			lock = unclaimedLock.Lock
		}
	}
	if lock.Expiration > blockNumber {
		err = fmt.Errorf("try to remove a lock which is not expired, expired=%d,currentBlockNumber=%d", lock.Expiration, blockNumber)
		return
	}
	newtree, newlocksroot, err = this.ComputeMerkleRootWithout(lock)
	if err != nil {
		return
	}
	if newtree == nil {
		newtree, err = transfer.NewMerkleTree(nil)
		if err != nil {
			return
		}
	}
	return
}

/*
Register a secret so that it can be used in a balance proof.

        Note:
            This methods needs to be called once a `Secret` message is received
            or a `SecretRevealed` event happens.

        Raises:
            ValueError: If the hashlock is not known.
        """
*/
func (this *ChannelEndState) RegisterSecret(secret common.Hash) error {
	hashlock := utils.Sha3(secret[:])
	if !this.IsKnown(hashlock) {
		return errors.New("secret does not correspond to any hashlock")
	}
	if this.IsLocked(hashlock) {
		pendingLock := this.Lock2PendingLocks[hashlock]
		delete(this.Lock2PendingLocks, hashlock)
		this.Lock2UnclaimedLocks[hashlock] = UnlockPartialProof{
			pendingLock.Lock, pendingLock.LockHashed, secret}
	}
	return nil
}

//Generate unlocking proofs for the known secrets
func (this *ChannelEndState) GetKnownUnlocks() []*UnlockProof {
	tree := this.TreeState.Tree
	var proofs []*UnlockProof
	for _, v := range this.Lock2UnclaimedLocks {
		proof := this.ComputeProofForLock(v.Secret, v.Lock, tree)
		proofs = append(proofs, proof)
	}
	return proofs
}

func (this *ChannelEndState) ComputeProofForLock(secret common.Hash, lock *encoding.Lock, tree *transfer.Merkletree) *UnlockProof {
	if tree == nil {
		tree = this.TreeState.Tree
	}
	lockEncoded := lock.AsBytes()
	lockhash := utils.Sha3(lockEncoded)
	merkleProof := tree.MakeProof(lockhash)
	return &UnlockProof{merkleProof, lockEncoded, secret}
}

//where to use?
func (this *ChannelEndState) Equal(other *ChannelEndState) bool {
	return this.ContractBalance.Cmp(other.ContractBalance) == 0 && this.Address == other.Address
}
