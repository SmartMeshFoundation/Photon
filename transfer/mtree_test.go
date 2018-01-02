package transfer

import (
	"bytes"
	"testing"

	"errors"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestMerkleTreeEmpty(t *testing.T) {
	tree, err := NewMerkleTree(nil)
	if err != nil {
		t.Error(err)
		return
	}
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
	hash0 := utils.Sha3([]byte{'x'})
	tree, err := NewMerkleTree([]common.Hash{hash0})
	if err != nil {
		t.Error(err)
		return
	}
	root := tree.MerkleRoot()
	if !bytes.Equal(root[:], hash0[:]) {
		t.Error("single merkle tree error")
	}
}

func TestMerkleTreeDuplicates(t *testing.T) {
	hash0 := utils.Sha3([]byte{'x'})
	hash1 := utils.Sha3([]byte{'y'})
	_, err := NewMerkleTree([]common.Hash{hash0, hash1})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = NewMerkleTree([]common.Hash{hash0, hash1, hash0})
	if err != ErrorDuplicateElement {
		t.Error(errors.New("duplicate error not found"))
		return
	}

}

func TestMerkleTreeOne(t *testing.T) {
	hash0 := utils.Sha3([]byte{'a', 'b'})
	tree, _ := NewMerkleTree([]common.Hash{hash0})
	root := tree.MerkleRoot()
	proof := tree.MakeProof(hash0)
	if !CheckProof(proof, root, hash0) {
		t.Error("check proof error")
	}
}

func TestMerkleTreeTwo(t *testing.T) {
	hash0 := utils.Sha3([]byte{'a'})
	hash1 := utils.Sha3([]byte{'b'})
	leaves := []common.Hash{hash0, hash1}
	tree, _ := NewMerkleTree(leaves)
	root := tree.MerkleRoot()
	proof0 := tree.MakeProof(hash0)
	if !CheckProof(proof0, root, hash0) {
		t.Error(errors.New("proof0 error"))
		return
	}
	proof1 := tree.MakeProof(hash1)
	if !CheckProof(proof1, root, hash1) {
		t.Error(errors.New("proof1 error"))
	}
}

/*
def test_three():
    def sort_join(first, second):
        return ''.join(sorted([first, second]))

    hash_0 = 'a' * 32
    hash_1 = 'b' * 32
    hash_2 = 'c' * 32

    leaves = [hash_0, hash_1, hash_2]
    tree = Merkletree(leaves)
    merkle_root = tree.merkleroot

    hash_01 = (
        b'me\xef\x9c\xa9=5\x16\xa4\xd3\x8a\xb7\xd9\x89\xc2\xb5\x00'
        b'\xe2\xfc\x89\xcc\xdc\xf8x\xf9\xc4m\xaa\xf6\xad\r['
    )
    assert keccak(hash_0 + hash_1) == hash_01
    calculated_root = keccak(hash_2 + hash_01)

    merkle_proof0 = tree.make_proof(hash_0)
    assert merkle_proof0 == [hash_1, hash_2]
    assert merkle_root == calculated_root
    assert check_proof(merkle_proof0, merkle_root, hash_0)

    merkle_proof1 = tree.make_proof(hash_1)
    assert merkle_proof1 == [hash_0, hash_2]
    assert merkle_root == calculated_root
    assert check_proof(merkle_proof1, merkle_root, hash_1)

    # with an odd number of values, the last value wont appear by itself in the
    # proof since it isn't hashed with another value
    merkle_proof2 = tree.make_proof(hash_2)
    assert merkle_proof2 == [keccak(hash_0 + hash_1)]
    assert merkle_root == calculated_root
    assert check_proof(merkle_proof2, merkle_root, hash_2)
*/

func TestMerkleTreeThree(t *testing.T) {
	hash0 := common.StringToHash("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	hash1 := common.StringToHash("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	hash2 := common.StringToHash("cccccccccccccccccccccccccccccccc")
	leaves := []common.Hash{hash0, hash1, hash2}
	hash01 := utils.Sha3(hash0[:], hash1[:])
	t.Logf("hash0=%s", hash0.String())
	t.Logf("hash1=%s", hash1.String())
	t.Logf("hash01=%s", hash01.String())
	tree, _ := NewMerkleTree(leaves)
	root := tree.MerkleRoot()
	t.Logf("root=%s", root.String())
	proof0 := tree.MakeProof(hash0)
	//spew.Dump("layers:", tree.Layers)
	//spew.Dump(proof0)
	if !CheckProof(proof0, root, hash0) {
		t.Error(errors.New("proof0 error"))
		return
	}
	proof1 := tree.MakeProof(hash1)
	if !CheckProof(proof1, root, hash1) {
		t.Error(errors.New("proof1 error"))
	}
	proof2 := tree.MakeProof(hash2)
	if !CheckProof(proof2, root, hash2) {
		t.Error(errors.New("proof2 error"))
	}
}

func TestMerkleTreeMany(t *testing.T) {
	var leaves []common.Hash
	for i := 0; i < 35; i++ {
		leaves = append(leaves, utils.Sha3([]byte{byte(i)}))
	}
	tree, _ := NewMerkleTree(leaves)
	for _, l := range leaves {
		proof := tree.MakeProof(l)
		if !CheckProof(proof, tree.MerkleRoot(), l) {
			t.Error(errors.New("proof many error"))
		}
	}

}
