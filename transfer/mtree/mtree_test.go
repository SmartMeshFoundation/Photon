package mtree

import (
	"bytes"
	"testing"

	"errors"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func newTestLock(exp int) *Lock {
	return &Lock{
		Amount:         big.NewInt(0),
		Expiration:     int64(exp),
		LockSecretHash: utils.EmptyHash,
	}
}
func TestMerkleTreeEmpty(t *testing.T) {
	tree := NewMerkleTree(nil)
	emptyroot := tree.MerkleRoot()
	if !bytes.Equal(emptyroot[:], utils.EmptyHash[:]) {
		t.Error("empty merkle tree  error")
	}
}

func TestMerkleTreeSingle(t *testing.T) {
	/*
			def test_single():
		    hash_0 = keccak('x')
	*/
	lock0 := newTestLock(0)
	tree := NewMerkleTree([]*Lock{lock0})
	root := tree.MerkleRoot()
	assert.EqualValues(t, root, lock0.Hash())
	assert.Empty(t, tree.MakeProof(lock0.Hash()))
}

func TestMerkleTreeDuplicates(t *testing.T) {
	lock0 := newTestLock(0)
	lock1 := newTestLock(2)
	tree := NewMerkleTree([]*Lock{lock0, lock1})
	assert.NotEqual(t, tree.MerkleRoot(), utils.EmptyHash)
	defer func() {
		if err := recover(); err != nil {
			t.Logf("err=%s", err)
		} else {
			t.Error("should panic")
		}
	}()
	tree = NewMerkleTree([]*Lock{lock0, lock1, newTestLock(2)})
}

func TestMerkleTreeOne(t *testing.T) {
	lock0 := newTestLock(0)
	tree := NewMerkleTree([]*Lock{lock0})
	root := tree.MerkleRoot()
	proof := tree.MakeProof(lock0.Hash())
	if !checkProof(proof, root, lock0.Hash()) {
		t.Error("check proof error")
	}
}

func TestMerkleTreeTwo(t *testing.T) {
	lock0 := newTestLock(0)
	lock1 := newTestLock(2)
	leaves := []*Lock{lock0, lock1}
	tree := NewMerkleTree(leaves)
	root := tree.MerkleRoot()
	proof0 := tree.MakeProof(lock0.Hash())
	if !checkProof(proof0, root, lock0.Hash()) {
		t.Error(errors.New("proof0 error"))
		return
	}
	proof1 := tree.MakeProof(lock1.Hash())
	if !checkProof(proof1, root, lock1.Hash()) {
		t.Error(errors.New("proof1 error"))
	}
}

func TestMerkleTreeThree(t *testing.T) {
	lock0 := newTestLock(0)
	lock1 := newTestLock(3)
	lock2 := newTestLock(7)
	leaves := []*Lock{lock0, lock1, lock2}
	tree := NewMerkleTree(leaves)
	root := tree.MerkleRoot()
	t.Logf("root=%s", root.String())
	proof0 := tree.MakeProof(lock0.Hash())
	//spew.Dump("layers:", tree.Layers)
	//spew.Dump(proof0)
	if !checkProof(proof0, root, lock0.Hash()) {
		t.Error(errors.New("proof0 error"))
		return
	}
	proof1 := tree.MakeProof(lock1.Hash())
	if !checkProof(proof1, root, lock1.Hash()) {
		t.Error(errors.New("proof1 error"))
	}
	proof2 := tree.MakeProof(lock2.Hash())
	if !checkProof(proof2, root, lock2.Hash()) {
		t.Error(errors.New("proof2 error"))
	}
}

func TestMerkleTreeMany(t *testing.T) {
	var leaves []*Lock
	for i := 0; i < 35; i++ {
		leaves = append(leaves, newTestLock(i))
	}
	tree := NewMerkleTree(leaves)
	for _, l := range leaves {
		proof := tree.MakeProof(l.Hash())
		if !checkProof(proof, tree.MerkleRoot(), l.Hash()) {
			t.Error(errors.New("proof many error"))
		}
	}

}
