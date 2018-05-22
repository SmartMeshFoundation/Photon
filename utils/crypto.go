package utils

import (
	"crypto/ecdsa"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//EmptyHash all zero,invalid
var EmptyHash = common.Hash{}

//EmptyAddress all zero,invalid
var EmptyAddress = common.Address{}

//SignData sign with ethereum format
func SignData(privKey *ecdsa.PrivateKey, data []byte) (sig []byte, err error) {
	hash := Sha3(data)
	//why add 27 for the last byte?
	sig, err = crypto.Sign(hash[:], privKey)
	if err == nil {
		sig[len(sig)-1] += byte(27)
	}
	return
}

//Sha3 is short for Keccak256Hash
func Sha3(data ...[]byte) common.Hash {
	return crypto.Keccak256Hash(data...)
}

//Pex short string stands for data
func Pex(data []byte) string {
	return common.Bytes2Hex(data[:4])
}

//HPex pex for hash
func HPex(data common.Hash) string {
	return common.Bytes2Hex(data[:2])
}

//APex pex for address
func APex(data common.Address) string {
	return common.Bytes2Hex(data[:4])
}

//APex2 shorter than APex
func APex2(data common.Address) string {
	return common.Bytes2Hex(data[:2])
}

//PubkeyToAddress convert pubkey bin to address
func PubkeyToAddress(pubkey []byte) common.Address {
	return common.BytesToAddress(crypto.Keccak256(pubkey[1:])[12:])
}

//BigIntTo32Bytes convert a big int to bytes
func BigIntTo32Bytes(i *big.Int) []byte {
	data := i.Bytes()
	buf := make([]byte, 32)
	for i := 0; i < 32-len(data); i++ {
		buf[i] = 0
	}
	for i := 32 - len(data); i < 32; i++ {
		buf[i] = data[i-32+len(data)]
	}
	return buf
}
