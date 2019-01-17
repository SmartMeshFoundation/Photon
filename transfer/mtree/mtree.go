package mtree

import (
	"bytes"

	"errors"

	"fmt"

	"math/big"

	"io"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

var errorDuplicateElement = errors.New("Duplicated element")

// LayerLeaves is layer 0
const LayerLeaves = 0

// LayerMerkleRoot is top layer
const LayerMerkleRoot = -1

/*
Merkletree is hash tree
*/
type Merkletree struct {
	Layers [][]common.Hash
	Leaves []*Lock
}

// EmptyTree contains no locks
var EmptyTree = NewMerkleTree(nil)

/*
Lock is 	The Lock structure for Hashed TimeLock Contract.
In our messenger
在消息中 expiration 用的是 int64编码
而在合约中用的是 big.Int
*/
/*
 *	Lock : data structure for Hashed TimeLock Contract.
 *
 *	In our messenger, expiration should be int64, but in contract, it should be big.Int.
 */
type Lock struct {
	Expiration     int64 // expiration block number
	Amount         *big.Int
	LockSecretHash common.Hash
}

//AsBytes serialize Lock
func (l *Lock) AsBytes() []byte {
	buf := new(bytes.Buffer)
	_, err := buf.Write(utils.BigIntTo32Bytes(big.NewInt(l.Expiration)))
	_, err = buf.Write(utils.BigIntTo32Bytes(l.Amount))
	_, err = buf.Write(l.LockSecretHash[:])
	if err != nil {
		log.Crit(fmt.Sprintf("Lock AsBytes error %s", err))
	}
	return buf.Bytes()
}

//FromBytes deserialize Lock
func (l *Lock) FromBytes(locksencoded []byte) (err error) {
	buf := bytes.NewBuffer(locksencoded)
	return l.FromReader(buf)
}

//FromReader init lock from a reader
func (l *Lock) FromReader(r io.Reader) (err error) {
	l.Expiration = utils.ReadBigInt(r).Int64()
	l.Amount = utils.ReadBigInt(r)
	_, err = r.Read(l.LockSecretHash[:])
	return
}

//Hash of this lock
func (l *Lock) Hash() common.Hash {
	return utils.Sha3(l.AsBytes())
}

func (l *Lock) String() string {
	return fmt.Sprintf("{expiration=%d,amount=%s,secrethash=%s}", l.Expiration, l.Amount, utils.HPex(l.LockSecretHash))
}

//Equal return true when the two locks are exactly the same.
func (l *Lock) Equal(l2 *Lock) bool {
	if l2 == nil {
		return false
	}
	if l2.Expiration == l.Expiration && l2.LockSecretHash == l.LockSecretHash && l.Amount.Cmp(l2.Amount) == 0 {
		return true
	}
	return false
}

/*
NewMerkleTree create merkle tree from locks
保证不要包含重复的锁,否则会panic
*/
/*
 *	NewMerkleTree : function to create merkle tree from locks.
 *
 *	Note that do not contain repeated locks, otherwise panic will occur.
 */
func NewMerkleTree(leaves []*Lock) (m *Merkletree) {
	var err error
	elements := make([]common.Hash, len(leaves))
	for i := 0; i < len(elements); i++ {
		l := leaves[i]
		buf := new(bytes.Buffer)
		_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(l.Expiration)))
		_, err = buf.Write(utils.BigIntTo32Bytes(l.Amount))
		_, err = buf.Write(l.LockSecretHash[:])
		elements[i] = utils.Sha3(buf.Bytes())
	}
	if err != nil {
		log.Crit(fmt.Sprintf("NewMerkleTree err %s", err))
	}
	m = new(Merkletree)
	m.buildMerkleTreeLayers(elements)
	m.Leaves = leaves
	return m
}

func lenDiv2(l int) int {
	return l/2 + l%2
}

/*
   computes the layers of the merkletree. First layer is the list
   of elements and the last layer is a list with a single entry, the
   merkleroot
*/
func (m *Merkletree) buildMerkleTreeLayers(elements []common.Hash) {
	if len(elements) == 0 {
		m.Layers = append(m.Layers, []common.Hash(nil)) //make sure has one empty layer
		return
	}
	elementsMap := make(map[common.Hash]bool)
	for _, e := range elements {
		if elementsMap[e] {
			panic(fmt.Sprintf("elements %s duplicated", e.String()))
		}
		elementsMap[e] = true
	}
	//sort.Slice(elements, func(i, j int) bool {
	//	return bytes.Compare(elements[i][:], elements[j][:]) == -1
	//})
	//m.Layers = append(m.Layers, elements)
	prevLayer := elements
	for i := 0; ; i++ {
		m.Layers = append(m.Layers, prevLayer)
		if len(prevLayer) == 1 {
			break
		}
		curLayer := make([]common.Hash, lenDiv2(len(prevLayer)))
		for j := 0; j < len(prevLayer); j += 2 {
			if j == len(prevLayer)-1 {
				curLayer[j/2] = prevLayer[j]
			} else {
				curLayer[j/2] = HashPair(prevLayer[j], prevLayer[j+1])
			}
		}
		prevLayer = curLayer
	}
	return
}

/*
MerkleRoot Return the root element of the merkle tree.
*/
func (m *Merkletree) MerkleRoot() common.Hash {
	l := len(m.Layers)
	if l > 1 {
		return m.Layers[l-1][0]
	} else if l == 1 && len(m.Layers[0]) > 0 {
		//root layer may be emtpy
		return m.Layers[l-1][0]
	} else {
		return utils.EmptyHash
	}
}

/*
MakeProof contains all elements between `element` and `root`.
If on all of [element] + proof is recursively hash_pair applied one
gets the root.
*/
func (m *Merkletree) MakeProof(element common.Hash) []common.Hash {
	idx := 0
	for i := range m.Layers[0] {
		if bytes.Equal(element[:], m.Layers[0][i][:]) {
			idx = i
		}
	}
	var proof []common.Hash
	for _, layer := range m.Layers {
		pairidx := idx - 1
		if idx%2 == 0 {
			pairidx = idx + 1
		}
		if pairidx < len(layer) {
			proof = append(proof, layer[pairidx])
		}
		idx = idx / 2
	}
	return proof
}

//Leaves2Byets get bytes of locks
func (m *Merkletree) Leaves2Byets() []byte {
	var err error
	buf := new(bytes.Buffer)
	for i := 0; i < len(m.Leaves); i++ {
		l := m.Leaves[i]
		_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(l.Expiration)))
		_, err = buf.Write(utils.BigIntTo32Bytes(l.Amount))
		_, err = buf.Write(l.LockSecretHash[:])
	}
	if err != nil {
		log.Crit(fmt.Sprintf("Leaves2Byets err %s", err))
	}
	return buf.Bytes()
}

/*
ComputeMerkleRootWith 创建包含新 lock 的 merkleTree
保证 include 不在原来的锁里
*/
/*
 *	ComputeMerkleRootWith : function to create a merkleTree with a new lock contained in.
 *
 *	Note that make sure `include` is not contained original by the merkle tree.
 */
func (m *Merkletree) ComputeMerkleRootWith(include *Lock) (newm *Merkletree) {
	//我们并不会更改锁的内容,只会进行不同的排列组合.
	leaves := make([]*Lock, len(m.Leaves))
	copy(leaves, m.Leaves)
	leaves = append(leaves, include)
	newm = NewMerkleTree(leaves)
	return
}

/*
返回移除without 的数组,如果不包含 without 直接返回本身
*/
/*
 *	removeHash : function to remove array of locks that have removed `without`.
 *
 *	Note that if `without` does not include in array, then return the array.
 */
func removeHash(leaves []*Lock, without *Lock) ([]*Lock, bool) {
	var r bool
	i := -1
	for j := 0; j < len(leaves); j++ {
		if leaves[j].Expiration == without.Expiration &&
			leaves[j].Amount.Cmp(without.Amount) == 0 &&
			leaves[j].LockSecretHash == without.LockSecretHash {
			i = j
			r = true
			break
		}
	}
	if i >= 0 {
		leaves = append(leaves[:i], leaves[i+1:]...)
	}
	return leaves, r
}

/*
ComputeMerkleRootWithout Compute the resulting merkle root if the lock `without` is exclude from the tree
*/
func (m *Merkletree) ComputeMerkleRootWithout(without *Lock) (newm *Merkletree, err error) {
	leaves := make([]*Lock, len(m.Leaves))
	copy(leaves, m.Leaves)
	leaves, hasLock := removeHash(leaves, without)
	if !hasLock {
		err = fmt.Errorf("no such lock %s", utils.HPex(without.LockSecretHash))
		return
	}
	if len(leaves) > 0 {
		newm = NewMerkleTree(leaves)
		return
	}
	newm = NewMerkleTree(nil)
	return
}

func (m *Merkletree) String() string {
	return fmt.Sprintf("MerkleTreeState{root:%s,layer level:%d}", m.MerkleRoot(), len(m.Layers))
}

/*
HashPair makes first and secod ordered,then hash them
*/
func HashPair(first, second common.Hash) common.Hash {
	if first == utils.EmptyHash {
		return second
	}
	if second == utils.EmptyHash {
		return first
	}
	if bytes.Compare(first[:], second[:]) > 0 {
		return utils.Sha3(second[:], first[:])
	}
	return utils.Sha3(first[:], second[:])
}

func checkProof(proof []common.Hash, root, hash common.Hash) bool {
	for _, x := range proof {
		hash = HashPair(hash, x)
	}
	return hash == root
}

//Proof2Bytes convert proof to bytes
func Proof2Bytes(proof []common.Hash) []byte {
	buf := new(bytes.Buffer)
	for _, h := range proof {
		_, err := buf.Write(h[:])
		if err != nil {
			log.Trace(fmt.Sprintf("Proof2Bytes write err %s", err))
		}
	}
	return buf.Bytes()
}
