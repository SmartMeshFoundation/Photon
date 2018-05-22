package transfer

import (
	"sort"

	"bytes"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

var errorDuplicateElement = errors.New("Duplicated element")

//LayerLeaves is layer 0
const LayerLeaves = 0

//LayerMerkleRoot is top layer
const LayerMerkleRoot = -1

/*
Merkletree is hash tree
*/
type Merkletree struct {
	Layers [][]common.Hash
}

//NewMerkleTree create Merkletree
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
func (m *Merkletree) buildMerkleTreeLayers(elements []common.Hash) error {
	if len(elements) == 0 {
		m.Layers = append(m.Layers, []common.Hash(nil)) //make sure has one empty layer
		return nil
	}
	elementsMap := make(map[common.Hash]bool)
	for _, e := range elements {
		if elementsMap[e] {
			return errorDuplicateElement
		}
		elementsMap[e] = true
	}
	sort.Slice(elements, func(i, j int) bool {
		return bytes.Compare(elements[i][:], elements[j][:]) == -1
	})
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
	return nil
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
		buf.Write(h[:])
	}
	return buf.Bytes()
}
