package channel

import (
	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/params"

	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
EndState Tracks the state of one of the participants in a channel
all the transfer (whenever lock or not ) I have sent.
*/
type EndState struct {
	Address             common.Address
	ContractBalance     *big.Int                                //lock protect race codition with Photonapi
	Lock2PendingLocks   map[common.Hash]channeltype.PendingLock //the lock I have sent
	Lock2UnclaimedLocks map[common.Hash]channeltype.UnlockPartialProof
	Tree                *mtree.Merkletree
	BalanceProofState   *transfer.BalanceProofState //race codition with Photonapi
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
	if c.BalanceProofState == nil {
		c.BalanceProofState = transfer.NewEmptyBalanceProofState()
	}
	return c
}

//TransferAmount is how many tokens I have sent to  partner.
func (node *EndState) TransferAmount() *big.Int {
	return node.BalanceProofState.TransferAmount
}

func (node *EndState) locksRoot() common.Hash {
	return node.BalanceProofState.LocksRoot
}

//SetContractTransferAmount update node's  transfer amount by contract event
func (node *EndState) SetContractTransferAmount(amount *big.Int) {
	// amount 为0,只有一种情况就是发生了 punish 事件
	// amount is 0, which means punish event occurs.
	if amount.Cmp(utils.BigInt0) != 0 && amount.Cmp(node.BalanceProofState.TransferAmount) < 0 {
		log.Error(fmt.Sprintf("ContractTransferAmount must be greater, ContractTransferAmount=%s,TransferAmount=%s",
			amount,
			node.BalanceProofState.TransferAmount,
		))
	}
	node.BalanceProofState.ContractTransferAmount = new(big.Int).Set(amount)
}
func (node *EndState) contractLocksRoot() common.Hash {
	return node.BalanceProofState.ContractLocksRoot
}
func (node *EndState) contractTransferAmount() *big.Int {
	return node.BalanceProofState.ContractTransferAmount
}

//SetContractLocksroot update node's locksroot by contract event
func (node *EndState) SetContractLocksroot(locksroot common.Hash) {
	node.BalanceProofState.ContractLocksRoot = locksroot
}

//SetContractNonce update node's nonce by contract event
func (node *EndState) SetContractNonce(nonce uint64) {
	node.BalanceProofState.ContractNonce = nonce
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
func (node *EndState) nonce() uint64 {
	return node.BalanceProofState.Nonce
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
		return rerr.ErrChannelBalanceDecrease
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

//GetUnkownSecretLockByHashlock returns the hash corresponding Lock,nil if not found
func (node *EndState) GetUnkownSecretLockByHashlock(lockSecretHash common.Hash) *mtree.Lock {
	lock, ok := node.Lock2PendingLocks[lockSecretHash]
	if ok {
		return lock.Lock
	}
	plock, ok := node.Lock2UnclaimedLocks[lockSecretHash]
	if ok && !plock.IsRegisteredOnChain {
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
	return nil, utils.EmptyHash, rerr.ErrChannelEndStateNoSuchLock
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
		return nil, utils.EmptyHash, rerr.ErrChannelEndStateNoSuchLock
	}
	newtree, err := node.Tree.ComputeMerkleRootWithout(without)
	if err != nil {
		return nil, utils.EmptyHash, err
	}
	return newtree, newtree.MerkleRoot(), nil
}
func (node *EndState) balanceProofRegisteredOnChain() bool {
	b := node.BalanceProofState
	//只要transferamount或者locksroot其一不为空,就表示已经上链更新过了.
	return b != nil &&
		((b.ContractTransferAmount != nil && b.ContractTransferAmount.Cmp(utils.BigInt0) != 0) ||
			b.ContractLocksRoot != utils.EmptyHash)
}

/*
registerDirectTransfer register a direct_transfer.

安全检查:
nonce,channel 由前置检查保证
transferAmount 必须增大,
locksroot 必须相等.
*/
/*
 *	registerDirectTransfer : function to register message of direct transfer.
 *
 *	Security Check :
 *		1. nonce, channel must be ensured by pre-set check.
 *		2. transferAmount must increase.
 *		3. locksroot must not get changed.
 */
func (node *EndState) registerDirectTransfer(directTransfer *encoding.DirectTransfer) error {
	if node.balanceProofRegisteredOnChain() {
		return rerr.ErrChannelBalanceProofAlreadyRegisteredOnChain
	}
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(directTransfer)
	if balanceProof.LocksRoot != node.Tree.MerkleRoot() {
		return rerr.InvalidLocksRoot(node.Tree.MerkleRoot(), balanceProof.LocksRoot)
	}
	node.BalanceProofState = balanceProof
	return nil
}

/*
RegisterRemoveExpiredHashlockTransfer register a RemoveExpiredHashlockTransfer
this message may be sent out from this node or received from partner
*/
//func (node *EndState) registerRemoveExpiredHashlockTransfer(removeExpiredHashlockTransfer *encoding.RemoveExpiredHashlockTransfer) error {
//	return node.registerRemoveLock(removeExpiredHashlockTransfer, removeExpiredHashlockTransfer.LockSecretHash)
//}

//func (node *EndState) registerAnnounceDisdposedTransferResponse(response *encoding.AnnounceDisposedResponse) error {
//	return node.registerRemoveLock(response, response.LockSecretHash)
//}

func (node *EndState) registerRemoveLock(msg encoding.EnvelopMessager, lockSecretHash common.Hash) error {
	if node.balanceProofRegisteredOnChain() {
		return rerr.ErrChannelBalanceProofAlreadyRegisteredOnChain
	}
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
/*
 *	registerSecretMessage : function to register message of secret.
 *
 *	Note that this message may be sent out from this node or received from his partner.
 *		1. this transfer lock exists.
 *		2. locksroot must remove lock of this transfer.
 *		3. transferAmount must equal to the value of previous transferAmount plus the amount of this transfer.
 */
func (node *EndState) registerSecretMessage(unlock *encoding.UnLock) (err error) {
	if node.balanceProofRegisteredOnChain() {
		return rerr.ErrChannelBalanceProofAlreadyRegisteredOnChain
	}
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(unlock)
	lockSecretHash := utils.ShaSecret(unlock.LockSecret[:])
	lock := node.getLockByHashlock(lockSecretHash)
	if lock == nil {
		err = rerr.ErrChannelLockSecretHashNotFound.Errorf(" receive unlock message,but has no related lockSecretHash,msg=%s", utils.StringInterface(unlock, 3))
		log.Error(err.Error())
		return err
	}
	newtree, newLocksroot, err := node.computeMerkleRootWithout(lock)
	if err != nil {
		return err
	}
	if balanceProof.LocksRoot != newLocksroot {
		return rerr.InvalidLocksRoot(newLocksroot, balanceProof.LocksRoot)
	}
	transferAmount := new(big.Int).Add(node.TransferAmount(), lock.Amount)
	/*
		金额只能是当前金额加上本次锁的金额,多了少了都是错的
	*/
	// transferAmount = previous transferAmount + token amount in this lock
	if unlock.TransferAmount.Cmp(transferAmount) != 0 {
		return rerr.ErrChannelTransferAmountMismatch.Errorf("invalid transferred_amount, expected: %s got: %s",
			transferAmount, unlock.TransferAmount)
	}
	delete(node.Lock2PendingLocks, lock.LockSecretHash)
	delete(node.Lock2UnclaimedLocks, lock.LockSecretHash)
	/*
		确保所有的信息都是正确的,才能更新状态
	*/
	// Verify messages are correct then update channel state.
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
/*
 *	registerMediatedMessage : function to register messages of MediatedTransfer.
 *
 *	Note that this message may be sent out from this node or received by his partner.
 *		1. the lock must not be existent.
 *		2. transferAmount must not have any change.
 * 		3. locksroot must be equal to value of previous locksroot plus the amount of this new amount.
 */
func (node *EndState) registerMediatedMessage(mtr *encoding.MediatedTransfer) (err error) {
	if node.balanceProofRegisteredOnChain() {
		return rerr.ErrChannelBalanceProofAlreadyRegisteredOnChain
	}
	balanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(mtr)
	mtranfer := encoding.GetMtrFromLockedTransfer(mtr)
	lock := mtranfer.GetLock()
	if node.IsKnown(lock.LockSecretHash) {
		return rerr.ErrChannelDuplicateLock
	}
	if balanceProof.TransferAmount.Cmp(node.TransferAmount()) < 0 {
		return rerr.ErrChannelTransferAmountDecrease.Errorf("transfer amount decrease,now=%s, message=%s", node.TransferAmount(), mtr)
	}
	newtree, locksroot := node.computeMerkleRootWith(lock)
	lockhashed := utils.Sha3(lock.AsBytes())
	if balanceProof.LocksRoot != locksroot {
		return rerr.InvalidLocksRoot(locksroot, balanceProof.LocksRoot)
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
	//链上已经注册密码的锁是一定不能移除的,无论什么原因都不能移除此锁.除非是unlock消息
	lock = node.GetUnkownSecretLockByHashlock(lockSecretHash)
	if lock == nil {
		err = rerr.ErrChannelEndStateNoSuchLock.Errorf("%s donesn't know hashlock %s, cannot remove", utils.APex(node.Address), utils.HPex(lockSecretHash))
		return
	}
	if mustExpired && (lock.Expiration > blockNumber-params.Cfg.ForkConfirmNumber) {
		err = rerr.ErrRemoveNotExpiredLock.Errorf("try to remove a lock which is not expired, expired=%d,currentBlockNumber=%d", lock.Expiration, blockNumber)
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
	hashlock := utils.ShaSecret(secret[:])
	if !node.IsKnown(hashlock) {
		return rerr.ErrChannelEndStateNoSuchLock
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
func (node *EndState) RegisterRevealedSecretHash(lockSecretHash, secret common.Hash, blockNumber int64) error {
	lock := node.getLockByHashlock(lockSecretHash)
	if lock == nil {
		return rerr.ErrChannelEndStateNoSuchLock
	}
	if blockNumber > lock.Expiration {
		return rerr.ErrChannelEndStateNoSuchLock.Errorf("secrethash %s  registerred on block chain,but already expired for me", utils.HPex(lockSecretHash))
	}
	//有可能这个lock已经在Lock2UnclaimedLocks,不过无所谓
	delete(node.Lock2PendingLocks, lockSecretHash)
	node.Lock2UnclaimedLocks[lockSecretHash] = channeltype.UnlockPartialProof{
		Lock:                lock,
		LockHash:            lock.Hash(),
		Secret:              secret,
		IsRegisteredOnChain: true,
	}
	return nil
}

//GetCanUnlockOnChainLocks generate unlocking proofs for the known secrets
func (node *EndState) GetCanUnlockOnChainLocks() []*channeltype.UnlockProof {
	tree := node.Tree
	var proofs []*channeltype.UnlockProof
	for _, v := range node.Lock2UnclaimedLocks {
		if !v.IsRegisteredOnChain {
			continue
		}
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
