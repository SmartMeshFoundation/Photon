package transfer

import (
	"sort"

	"bytes"

	"errors"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
)

var ErrorDuplicateElement = errors.New("Duplicated element")

const LayerLeaves = 0
const LayerMerkleRoot = -1

type Merkletree struct {
	Layers [][]common.Hash
}

func NewMerkleTree(elements []common.Hash) (m *Merkletree, err error) {
	m = new(Merkletree)
	err = m.buildMerkleTreeLayers(elements)
	return m, err
}
func lenDiv2(l int) int {
	return l/2 + l%2
}

/*
   computes the layers of the merkletree. First layer is the list
   of elements and the last layer is a list with a single entry, the
   merkleroot
*/
func (this *Merkletree) buildMerkleTreeLayers(elements []common.Hash) error {
	if len(elements) == 0 {
		this.Layers = append(this.Layers, []common.Hash(nil)) //make sure has one empty layer
		return nil
	}
	elementsMap := make(map[common.Hash]bool)
	for _, e := range elements {
		if elementsMap[e] {
			return ErrorDuplicateElement
		}
		elementsMap[e] = true
	}
	sort.Slice(elements, func(i, j int) bool {
		return bytes.Compare(elements[i][:], elements[j][:]) == -1
	})
	//this.Layers = append(this.Layers, elements)
	prevLayer := elements
	for i := 0; ; i++ {
		this.Layers = append(this.Layers, prevLayer)
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
	return nil
}

/*
""" Return the root element of the merkle tree. """
*/
func (this *Merkletree) MerkleRoot() common.Hash {
	l := len(this.Layers)
	if l > 1 {
		return this.Layers[l-1][0]
	} else if l == 1 && len(this.Layers[0]) > 0 {
		//root layer may be emtpy
		return this.Layers[l-1][0]
	} else {
		return utils.EmptyHash
	}
}

/*
		The proof contains all elements between `element` and `root`.
            If on all of [element] + proof is recursively hash_pair applied one
            gets the root.
*/
func (this *Merkletree) MakeProof(element common.Hash) []common.Hash {
	idx := 0
	for i, _ := range this.Layers[0] {
		if bytes.Equal(element[:], this.Layers[0][i][:]) {
			idx = i
		}
	}
	var proof []common.Hash
	for _, layer := range this.Layers {
		pair_idx := idx - 1
		if idx%2 == 0 {
			pair_idx = idx + 1
		}
		if pair_idx < len(layer) {
			proof = append(proof, layer[pair_idx])
		}
		idx = idx / 2
	}
	return proof
}

/*
def hash_pair(first, second):
    if second is None:
        return first
    if first is None:
        return second
    if first > second:
        return keccak(second + first)
    return keccak(first + second)
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
	} else {
		return utils.Sha3(first[:], second[:])
	}
}

func CheckProof(proof []common.Hash, root, hash common.Hash) bool {
	for _, x := range proof {
		hash = HashPair(hash, x)
	}
	return hash == root
}

func Proof2Bytes(proof []common.Hash) []byte {
	buf := new(bytes.Buffer)
	for _, h := range proof {
		buf.Write(h[:])
	}
	return buf.Bytes()
}
