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

/*
PendingLock is lock of HTLC
*/
type PendingLock struct {
	Lock       *encoding.Lock
	LockHashed common.Hash
}

/*
UnlockPartialProof is the lock that I have known the secret ,but haven't receive the balance proof
*/
type UnlockPartialProof struct {
	Lock       *encoding.Lock
	LockHashed common.Hash
	Secret     common.Hash
}

/*
UnlockProof is the info needs withdraw on blockchain
*/
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

/*
InvalidLocksRootError is a wrong locks root error
*/
type InvalidLocksRootError struct {
	ExpectedLocksroot common.Hash
	GotLocksroot      common.Hash
}

//Error is err.Error interface
func (ilre *InvalidLocksRootError) Error() string {
	return fmt.Sprintf("Locksroot mismatch. Expected %s but get %s",
		utils.Pex(ilre.ExpectedLocksroot[:]), utils.Pex(ilre.GotLocksroot[:]))
}

var errBalanceDecrease = errors.New("contract_balance cannot decrease")
var errUnknownLock = errors.New("'unknown lock")
var errTransferAmountMismatch = errors.New("transfer amount mismatch")

/*
EndState Tracks the state of one of the participants in a channel
all the transfer (whenever lock or not ) I have sent.
*/
type EndState struct {
	Address             common.Address
	ContractBalance     *big.Int                    //lock protect race codition with raidenapi
	Lock2PendingLocks   map[common.Hash]PendingLock //the lock I have sent
	Lock2UnclaimedLocks map[common.Hash]UnlockPartialProof
	TreeState           *transfer.MerkleTreeState
	BalanceProofState   *transfer.BalanceProofState //race codition with raidenapi
}

//NewChannelEndState create EndState
func NewChannelEndState(participantAddress common.Address, participantBalance *big.Int,
	balanceProof *transfer.BalanceProofState, tree *transfer.MerkleTreeState) *EndState {
	c := &EndState{
		Address:             participantAddress,
		ContractBalance:     participantBalance,
		TreeState:           tree,
		BalanceProofState:   balanceProof,
		Lock2PendingLocks:   make(map[common.Hash]PendingLock),
		Lock2UnclaimedLocks: make(map[common.Hash]UnlockPartialProof),
	}
	return c
}

//TransferAmount is how many tokens I have sent to  partner.
func (node *EndState) TransferAmount() *big.Int {
	if node.BalanceProofState != nil {
		return node.BalanceProofState.TransferAmount
	}
	return big.NewInt(0)
}

//amountLocked is the tokens I have sent but partner doesn't have received the new blanceproof
func (node *EndState) amountLocked() *big.Int {
	sum := big.NewInt(0)
	for _, v := range node.Lock2PendingLocks {
		sum = sum.Add(sum, v.Lock.Amount)
	}
	for _, v := range node.Lock2UnclaimedLocks {
		sum = sum.Add(sum, v.Lock.Amount)
	}
	return sum
}

//nonce returns next nonce of this node.
func (node *EndState) nonce() int64 {
	if node.BalanceProofState != nil {
		return node.BalanceProofState.Nonce
	}
	return 0
}

//balance returns the availabe tokens i have
func (node *EndState) balance(counterpart *EndState) *big.Int {
	x := new(big.Int).Sub(node.ContractBalance, node.TransferAmount())
	x.Add(x, counterpart.TransferAmount())
	return x
}

//distributable returns the availabe tokens i can send to partner. this equals `balance`-`amountLocked`
func (node *EndState) distributable(counterpart *EndState) *big.Int {
	return new(big.Int).Sub(node.balance(counterpart), node.amountLocked())
}

//IsKnown returns True if the `hashlock` corresponds to a known lock.
func (node *EndState) IsKnown(hashlock common.Hash) bool {
	_, ok := node.Lock2PendingLocks[hashlock]
	if ok {
		return ok
	}
	_, ok = node.Lock2UnclaimedLocks[hashlock]
	return ok
}

//isLocked returns True if the `hashlock` is known and the correspoding secret is not.
func (node *EndState) isLocked(hashlock common.Hash) bool {
	_, ok := node.Lock2PendingLocks[hashlock]
	return ok
}

/*
UpdateContractBalance returns Update the contract balance, it must always increase.

return error If the `contract_balance` is smaller than the current
           balance.
*/
func (node *EndState) UpdateContractBalance(balance *big.Int) error {
	if balance.Cmp(node.ContractBalance) < 0 {
		return errBalanceDecrease
	}
	node.ContractBalance = new(big.Int).Set(balance)
	return nil
}

//getLockByHashlock returns the hash corresponding Lock,nil if not found
func (node *EndState) getLockByHashlock(hashlock common.Hash) *encoding.Lock {
	lock, ok := node.Lock2PendingLocks[hashlock]
	if ok {
		return lock.Lock
	}
	plock, ok := node.Lock2UnclaimedLocks[hashlock]
	if ok {
		return plock.Lock
	}
	return nil
}

/*
computeMerkleRootWith Compute the resulting merkle root if the lock `include` is added in
       the tree.
*/
func (node *EndState) computeMerkleRootWith(include *encoding.Lock) (tree *transfer.Merkletree, hash common.Hash) {
	if !node.IsKnown(include.HashLock) {
		leaves := make([]common.Hash, len(node.TreeState.Tree.Layers[transfer.LayerLeaves]))
		copy(leaves, node.TreeState.Tree.Layers[transfer.LayerLeaves])
		includeHash := utils.Sha3(include.AsBytes())
		leaves = append(leaves, includeHash)
		tree, err := transfer.NewMerkleTree(leaves)
		if err != nil {
			log.Error(fmt.Sprintf("NewMerkleTree err %s", err))
		}
		return tree, tree.MerkleRoot()
	}
	return nil, node.TreeState.Tree.MerkleRoot()
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
 computeMerkleRootWithout Compute the resulting merkle root if the lock `without` is exclude from the tree
*/
func (node *EndState) computeMerkleRootWithout(without *encoding.Lock) (*transfer.Merkletree, common.Hash, error) {
	if !node.IsKnown(without.HashLock) {
		return nil, utils.EmptyHash, errUnknownLock
	}
	leaves := make([]common.Hash, len(node.TreeState.Tree.Layers[transfer.LayerLeaves]))
	copy(leaves, node.TreeState.Tree.Layers[transfer.LayerLeaves])
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
  registerLockedTransfer  API design: using specialized methods to force the user to register the
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
func (node *EndState) registerLockedTransfer(lockedTransfer encoding.EnvelopMessager) error {
	if !encoding.IsLockedTransfer(lockedTransfer) {
		return errors.New("not a locked lockedTransfer")
	}
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(lockedTransfer)
	mtranfer := encoding.GetMtrFromLockedTransfer(lockedTransfer)
	lock := mtranfer.GetLock()
	if node.IsKnown(lock.HashLock) {
		return errors.New("hashlock is already registered")
	}
	newtree, locksroot := node.computeMerkleRootWith(lock)
	lockhashed := utils.Sha3(lock.AsBytes())
	if balanceProof.LocksRoot != locksroot {
		return &InvalidLocksRootError{locksroot, balanceProof.LocksRoot}
	}
	node.Lock2PendingLocks[lock.HashLock] = PendingLock{lock, lockhashed}
	node.BalanceProofState = balanceProof
	node.TreeState = transfer.NewMerkleTreeState(newtree)
	return nil
}

/*
registerDirectTransfer register a direct_transfer.

       Raises:
           InvalidLocksRoot: If the merkleroot of `direct_transfer` does not
           match the current value.
*/
func (node *EndState) registerDirectTransfer(directTransfer *encoding.DirectTransfer) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(directTransfer)
	if balanceProof.LocksRoot != node.TreeState.Tree.MerkleRoot() {
		return &InvalidLocksRootError{node.TreeState.Tree.MerkleRoot(), balanceProof.LocksRoot}
	}
	node.BalanceProofState = balanceProof
	return nil
}

/*
RegisterRemoveExpiredHashlockTransfer register a RemoveExpiredHashlockTransfer
this message may be sent out from this node or received from partner
*/
func (node *EndState) registerRemoveExpiredHashlockTransfer(removeExpiredHashlockTransfer *encoding.RemoveExpiredHashlockTransfer) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(removeExpiredHashlockTransfer)
	if balanceProof.TransferAmount.Cmp(node.TransferAmount()) != 0 {
		return errTransferAmountMismatch
	}
	node.BalanceProofState = balanceProof
	delete(node.Lock2PendingLocks, removeExpiredHashlockTransfer.HashLock)
	delete(node.Lock2UnclaimedLocks, removeExpiredHashlockTransfer.HashLock)
	return nil
}

/*
registerSecretMessage register a secret message
this message may be sent out from this node or received from partner
*/
func (node *EndState) registerSecretMessage(secret *encoding.Secret) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(secret)
	hashlock := utils.Sha3(secret.Secret[:])
	pendingLock, ok := node.Lock2PendingLocks[hashlock]
	var lock *encoding.Lock
	if ok {
		lock = pendingLock.Lock
	} else {
		unclaimedLock, ok := node.Lock2UnclaimedLocks[hashlock]
		if ok {
			lock = unclaimedLock.Lock
		}
	}
	//if !this.IsKnown(lock.HashLock) { // has corrupted because of lock is nil
	//	return errors.New("hashlock is not registered")
	//}
	newtree, newLocksroot, err := node.computeMerkleRootWithout(lock)
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
	delete(node.Lock2PendingLocks, lock.HashLock)
	delete(node.Lock2UnclaimedLocks, lock.HashLock)
	node.TreeState = transfer.NewMerkleTreeState(newtree)
	node.BalanceProofState = balanceProof
	return nil
}

/*
TryRemoveExpiredHashLock try to remomve a expired hashlock
*/
func (node *EndState) TryRemoveExpiredHashLock(hashlock common.Hash, blockNumber int64) (lock *encoding.Lock, newtree *transfer.Merkletree, newlocksroot common.Hash, err error) {
	if !node.IsKnown(hashlock) {
		err = fmt.Errorf("channel %s donesn't know hashlock %s, cannot remove", utils.APex(node.Address), utils.HPex(hashlock))
		return
	}
	pendingLock, ok := node.Lock2PendingLocks[hashlock]
	if ok {
		lock = pendingLock.Lock
	} else {
		unclaimedLock, ok := node.Lock2UnclaimedLocks[hashlock]
		if ok {
			lock = unclaimedLock.Lock
		}
	}
	if lock.Expiration > blockNumber {
		err = fmt.Errorf("try to remove a lock which is not expired, expired=%d,currentBlockNumber=%d", lock.Expiration, blockNumber)
		return
	}
	newtree, newlocksroot, err = node.computeMerkleRootWithout(lock)
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
RegisterSecret register a secret(not secret message) so that it can be used in a balance proof.

        Note:
            This methods needs to be called once a `Secret` message is received
            or a `SecretRevealed` event happens.

        Raises:
            ValueError: If the hashlock is not known.
        """
*/
func (node *EndState) RegisterSecret(secret common.Hash) error {
	hashlock := utils.Sha3(secret[:])
	if !node.IsKnown(hashlock) {
		return errors.New("secret does not correspond to any hashlock")
	}
	if node.isLocked(hashlock) {
		pendingLock := node.Lock2PendingLocks[hashlock]
		delete(node.Lock2PendingLocks, hashlock)
		node.Lock2UnclaimedLocks[hashlock] = UnlockPartialProof{
			pendingLock.Lock, pendingLock.LockHashed, secret}
	}
	return nil
}

//GetKnownUnlocks generate unlocking proofs for the known secrets
func (node *EndState) GetKnownUnlocks() []*UnlockProof {
	tree := node.TreeState.Tree
	var proofs []*UnlockProof
	for _, v := range node.Lock2UnclaimedLocks {
		proof := node.computeProofForLock(v.Secret, v.Lock, tree)
		proofs = append(proofs, proof)
	}
	return proofs
}

//computeProofForLock returns unlockProof need by contracts
func (node *EndState) computeProofForLock(secret common.Hash, lock *encoding.Lock, tree *transfer.Merkletree) *UnlockProof {
	if tree == nil {
		tree = node.TreeState.Tree
	}
	lockEncoded := lock.AsBytes()
	lockhash := utils.Sha3(lockEncoded)
	merkleProof := tree.MakeProof(lockhash)
	return &UnlockProof{merkleProof, lockEncoded, secret}
}

//where to use?
func (node *EndState) equal(other *EndState) bool {
	return node.ContractBalance.Cmp(other.ContractBalance) == 0 && node.Address == other.Address
}
