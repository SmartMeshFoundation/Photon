package channel

import (
	"errors"

	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

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
	ContractBalance     *big.Int                                //lock protect race codition with raidenapi
	Lock2PendingLocks   map[common.Hash]channeltype.PendingLock //the lock I have sent
	Lock2UnclaimedLocks map[common.Hash]channeltype.UnlockPartialProof
	Tree                *mtree.Merkletree
	BalanceProofState   *transfer.BalanceProofState //race codition with raidenapi
}

//NewChannelEndState create EndState
func NewChannelEndState(participantAddress common.Address, participantBalance *big.Int,
	balanceProof *transfer.BalanceProofState, tree *mtree.Merkletree) *EndState {
	c := &EndState{
		Address:             participantAddress,
		ContractBalance:     participantBalance,
		Tree:                tree,
		BalanceProofState:   balanceProof,
		Lock2PendingLocks:   make(map[common.Hash]channeltype.PendingLock),
		Lock2UnclaimedLocks: make(map[common.Hash]channeltype.UnlockPartialProof),
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

//SetContractTransferAmount update node's  transfer amount by contract event
func (node *EndState) SetContractTransferAmount(amount *big.Int) {
	if node.BalanceProofState != nil {
		if node.BalanceProofState.ContractTransferAmount.Cmp(node.BalanceProofState.TransferAmount) <= 0 {
			panic(fmt.Sprintf("ContractTransferAmount must be greater, ContractTransferAmount=%s,TransferAmount=%s",
				node.BalanceProofState.ContractTransferAmount,
				node.BalanceProofState.TransferAmount,
			))
		}
		node.BalanceProofState.ContractTransferAmount = new(big.Int).Set(amount)
	}
	return
}
func (node *EndState) contractTransferAmount() *big.Int {
	if node.BalanceProofState != nil {
		return node.BalanceProofState.ContractTransferAmount
	}
	return big.NewInt(0)
}

//SetContractLocksroot update node's locksroot by contract event
func (node *EndState) SetContractLocksroot(locksroot common.Hash) {
	if node.BalanceProofState != nil {
		node.BalanceProofState.ContractLocksRoot = locksroot
	}
}

//SetContractNonce update node's nonce by contract event
func (node *EndState) SetContractNonce(nonce int64) {
	if node.BalanceProofState != nil {
		node.BalanceProofState.Nonce = nonce
	}
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

//Balance returns the availabe tokens i have
func (node *EndState) Balance(counterpart *EndState) *big.Int {
	x := new(big.Int).Sub(node.ContractBalance, node.TransferAmount())
	x.Add(x, counterpart.TransferAmount())
	return x
}

//Distributable returns the availabe tokens i can send to partner. this equals `Balance`-`amountLocked`
func (node *EndState) Distributable(counterpart *EndState) *big.Int {
	return new(big.Int).Sub(node.Balance(counterpart), node.amountLocked())
}

//IsKnown returns True if the `hashlock` corresponds to a known lock.
func (node *EndState) IsKnown(lockSecretHash common.Hash) bool {
	_, ok := node.Lock2PendingLocks[lockSecretHash]
	if ok {
		return ok
	}
	_, ok = node.Lock2UnclaimedLocks[lockSecretHash]
	return ok
}

//GetSecret returns the secret corresponds to the lockSecretHash if found
func (node *EndState) GetSecret(lockSecretHash common.Hash) (secret common.Hash, found bool) {
	l, found := node.Lock2UnclaimedLocks[lockSecretHash]
	if found {
		secret = l.Secret
	}
	return
}

//IsLocked returns True if the `hashlock` is known and the correspoding secret is not.
func (node *EndState) IsLocked(hashlock common.Hash) bool {
	_, ok := node.Lock2PendingLocks[hashlock]
	return ok
}

/*
UpdateContractBalance returns Update the contract Balance, it must always increase.

return error If the `contract_balance` is smaller than the current
           Balance.
*/
func (node *EndState) UpdateContractBalance(balance *big.Int) error {
	if balance.Cmp(node.ContractBalance) < 0 {
		return errBalanceDecrease
	}
	node.ContractBalance = new(big.Int).Set(balance)
	return nil
}

//getLockByHashlock returns the hash corresponding Lock,nil if not found
func (node *EndState) getLockByHashlock(lockSecretHash common.Hash) *mtree.Lock {
	lock, ok := node.Lock2PendingLocks[lockSecretHash]
	if ok {
		return lock.Lock
	}
	plock, ok := node.Lock2UnclaimedLocks[lockSecretHash]
	if ok {
		return plock.Lock
	}
	return nil
}

/*
getSecretByLockSecretHash get secret by secret's hash
*/
func (node *EndState) getSecretByLockSecretHash(lockSecretHash common.Hash) (lock *mtree.Lock, secret common.Hash, err error) {
	plock, ok := node.Lock2UnclaimedLocks[lockSecretHash]
	if ok {
		return plock.Lock, plock.Secret, nil
	}
	return nil, utils.EmptyHash, errors.New("not found")
}

/*
computeMerkleRootWith Compute the resulting merkle root if the lock `include` is added in
       the tree.
*/
func (node *EndState) computeMerkleRootWith(include *mtree.Lock) (tree *mtree.Merkletree, hash common.Hash) {
	if !node.IsKnown(include.LockSecretHash) {
		tree := node.Tree.ComputeMerkleRootWith(include)
		return tree, tree.MerkleRoot()
	}
	return nil, node.Tree.MerkleRoot()
}

/*
 computeMerkleRootWithout Compute the resulting merkle root if the lock `without` is exclude from the tree
*/
func (node *EndState) computeMerkleRootWithout(without *mtree.Lock) (*mtree.Merkletree, common.Hash, error) {
	if !node.IsKnown(without.LockSecretHash) {
		return nil, utils.EmptyHash, errUnknownLock
	}
	newtree, err := node.Tree.ComputeMerkleRootWithout(without)
	if err != nil {
		return nil, utils.EmptyHash, err
	}
	return newtree, newtree.MerkleRoot(), nil
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
          lockedTransfer: The transfer to be added.

//Calculate the banlanceproof locksroot position before sending
*/
func (node *EndState) registerLockedTransfer(lockedTransfer encoding.EnvelopMessager) error {
	if lockedTransfer.Cmd() != encoding.MediatedTransferCmdID {
		return errors.New("not a locked lockedTransfer")
	}
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(lockedTransfer)
	mtranfer := encoding.GetMtrFromLockedTransfer(lockedTransfer)
	lock := mtranfer.GetLock()
	if node.IsKnown(lock.LockSecretHash) {
		return errors.New("hashlock is already registered")
	}
	newtree, locksroot := node.computeMerkleRootWith(lock)
	if balanceProof.LocksRoot != locksroot {
		return &InvalidLocksRootError{
			ExpectedLocksroot: locksroot,
			GotLocksroot:      balanceProof.LocksRoot,
		}
	}
	node.Lock2PendingLocks[lock.LockSecretHash] = channeltype.PendingLock{
		Lock:     lock,
		LockHash: lock.Hash(),
	}
	node.BalanceProofState = balanceProof
	node.Tree = newtree
	return nil
}

/*
registerDirectTransfer register a direct_transfer.

安全检查:
nonce,channel 由前置检查保证
transferAmount 必须增大,
locksroot 必须相等.
*/
func (node *EndState) registerDirectTransfer(directTransfer *encoding.DirectTransfer) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(directTransfer)
	if balanceProof.LocksRoot != node.Tree.MerkleRoot() {
		return &InvalidLocksRootError{node.Tree.MerkleRoot(), balanceProof.LocksRoot}
	}
	node.BalanceProofState = balanceProof
	return nil
}

/*
RegisterRemoveExpiredHashlockTransfer register a RemoveExpiredHashlockTransfer
this message may be sent out from this node or received from partner
*/
func (node *EndState) registerRemoveExpiredHashlockTransfer(removeExpiredHashlockTransfer *encoding.RemoveExpiredHashlockTransfer) error {
	return node.registerRemoveLock(removeExpiredHashlockTransfer, removeExpiredHashlockTransfer.LockSecretHash)
}

func (node *EndState) registerAnnounceDisdposedTransferResponse(response *encoding.AnnounceDisposedResponse) error {
	return node.registerRemoveLock(response, response.LockSecretHash)
}

func (node *EndState) registerRemoveLock(msg encoding.EnvelopMessager, lockSecretHash common.Hash) error {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(msg)
	node.BalanceProofState = balanceProof
	delete(node.Lock2PendingLocks, lockSecretHash)
	delete(node.Lock2UnclaimedLocks, lockSecretHash)
	return nil
}

/*
registerSecretMessage register a secret message
this message may be sent out from this node or received from partner
1.有这个锁
2.locksroot 要恰好等于去掉这个锁
3.transferAmount 要恰好等于这个锁的金额加上历史 transferAmount
*/
func (node *EndState) registerSecretMessage(unlock *encoding.UnLock) (err error) {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(unlock)
	lockSecretHash := utils.Sha3(unlock.LockSecret[:])
	lock := node.getLockByHashlock(lockSecretHash)
	if lock == nil {
		err = fmt.Errorf(" receive unlock message,but has no related lockSecretHash,msg=%s", utils.StringInterface(unlock, 3))
		log.Error(err.Error())
		return err
	}
	newtree, newLocksroot, err := node.computeMerkleRootWithout(lock)
	if err != nil {
		return err
	}
	if balanceProof.LocksRoot != newLocksroot {
		return &InvalidLocksRootError{newLocksroot, balanceProof.LocksRoot}
	}
	transferAmount := new(big.Int).Add(node.TransferAmount(), lock.Amount)
	/*
		金额只能是当前金额加上本次锁的金额,多了少了都是错的
	*/
	if unlock.TransferAmount.Cmp(transferAmount) != 0 {
		return fmt.Errorf("invalid transferred_amount, expected: %s got: %s",
			transferAmount, unlock.TransferAmount)
	}
	delete(node.Lock2PendingLocks, lock.LockSecretHash)
	delete(node.Lock2UnclaimedLocks, lock.LockSecretHash)
	/*
		确保所有的信息都是正确的,才能更新状态
	*/
	node.Tree = newtree
	node.BalanceProofState = balanceProof
	return nil
}

/*
registerMediatedMessage register a MediateTransfer message
this message may be sent out from this node or received from partner
1.这个锁一定要没有出现过
2.transferAmount 必须不变
3.locksroot 要恰好等于旧 locksroot 加上新锁
*/
func (node *EndState) registerMediatedMessage(mtr *encoding.MediatedTransfer) (err error) {
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(mtr)
	mtranfer := encoding.GetMtrFromLockedTransfer(mtr)
	lock := mtranfer.GetLock()
	if node.IsKnown(lock.LockSecretHash) {
		return errors.New("hashlock is already registered")
	}
	if node.getLockByHashlock(mtr.LockSecretHash) != nil {
		return fmt.Errorf("MediatedTransfer has duplicated lock, mtr=%s", mtr)
	}
	if balanceProof.TransferAmount.Cmp(node.TransferAmount()) < 0 {
		return fmt.Errorf("transfer amount decrease,now=%s, message=%s", node.TransferAmount(), mtr)
	}
	newtree, locksroot := node.computeMerkleRootWith(lock)
	lockhashed := utils.Sha3(lock.AsBytes())
	if balanceProof.LocksRoot != locksroot {
		return &InvalidLocksRootError{
			ExpectedLocksroot: locksroot,
			GotLocksroot:      balanceProof.LocksRoot,
		}
	}
	node.Lock2PendingLocks[lock.LockSecretHash] = channeltype.PendingLock{
		Lock:     lock,
		LockHash: lockhashed,
	}
	node.BalanceProofState = balanceProof
	node.Tree = newtree
	return nil
}

/*
TryRemoveHashLock try to remomve a expired hashlock
*/
func (node *EndState) TryRemoveHashLock(lockSecretHash common.Hash, blockNumber int64, mustExpired bool) (lock *mtree.Lock, newtree *mtree.Merkletree, newlocksroot common.Hash, err error) {
	lock = node.getLockByHashlock(lockSecretHash)
	if lock == nil {
		err = fmt.Errorf("%s donesn't know hashlock %s, cannot remove", utils.APex(node.Address), utils.HPex(lockSecretHash))
		return
	}
	if mustExpired && (lock.Expiration > blockNumber) {
		err = fmt.Errorf("try to remove a lock which is not expired, expired=%d,currentBlockNumber=%d", lock.Expiration, blockNumber)
		return
	}
	newtree, newlocksroot, err = node.computeMerkleRootWithout(lock)
	if err != nil {
		return
	}
	return
}

/*
RegisterSecret register a secret(not secret message) so that it can be used in a Balance proof.

            This methods needs to be called once a `Secret` message is received
*/
func (node *EndState) RegisterSecret(secret common.Hash) error {
	hashlock := utils.Sha3(secret[:])
	if !node.IsKnown(hashlock) {
		return errors.New("secret does not correspond to any hashlock")
	}
	if node.IsLocked(hashlock) {
		pendingLock := node.Lock2PendingLocks[hashlock]
		delete(node.Lock2PendingLocks, hashlock)
		node.Lock2UnclaimedLocks[hashlock] = channeltype.UnlockPartialProof{
			Lock:     pendingLock.Lock,
			LockHash: pendingLock.LockHash,
			Secret:   secret,
		}
	}
	return nil
}

/*
RegisterRevealedSecretHash a SecretReveal event on chain
*/
func (node *EndState) RegisterRevealedSecretHash(lockSecretHash common.Hash, blockNumber int64) error {
	if !node.IsKnown(lockSecretHash) {
		return errors.New("secret does not correspond to any lockSecretHash")
	}
	if node.IsLocked(lockSecretHash) {
		pendingLock := node.Lock2PendingLocks[lockSecretHash]
		if blockNumber > pendingLock.Lock.Expiration {
			return fmt.Errorf("secrethash %s  registerred on block chain,but already expired for me", utils.HPex(lockSecretHash))
		}
		delete(node.Lock2PendingLocks, lockSecretHash)
		node.Lock2UnclaimedLocks[lockSecretHash] = channeltype.UnlockPartialProof{
			Lock:     pendingLock.Lock,
			LockHash: pendingLock.LockHash,
			Secret:   utils.EmptyHash,
		}
	}
	return nil
}

//GetKnownUnlocks generate unlocking proofs for the known secrets
func (node *EndState) GetKnownUnlocks() []*channeltype.UnlockProof {
	tree := node.Tree
	var proofs []*channeltype.UnlockProof
	for _, v := range node.Lock2UnclaimedLocks {
		proof := ComputeProofForLock(v.Lock, tree)
		proofs = append(proofs, proof)
	}
	return proofs
}

//ComputeProofForLock returns unlockProof need by contracts
func ComputeProofForLock(lock *mtree.Lock, tree *mtree.Merkletree) *channeltype.UnlockProof {
	lockEncoded := lock.AsBytes()
	lockhash := utils.Sha3(lockEncoded)
	merkleProof := tree.MakeProof(lockhash)
	return &channeltype.UnlockProof{
		MerkleProof: merkleProof,
		Lock:        lock,
	}
}

//where to use?
func (node *EndState) equal(other *EndState) bool {
	return node.ContractBalance.Cmp(other.ContractBalance) == 0 && node.Address == other.Address
}
