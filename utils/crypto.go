package utils

import (
	"crypto/ecdsa"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var EmptyHash = common.Hash{}
var EmptyAddress = common.Address{}

func SignData(privKey *ecdsa.PrivateKey, data []byte) (sig []byte, err error) {
	hash := Sha3(data)
	//why add 27 for the last byte?
	sig, err = crypto.Sign(hash[:], privKey)
	if err == nil {
		sig[len(sig)-1] += byte(27)
	}
	return
}
func Sha3(data ...[]byte) common.Hash {
	return crypto.Keccak256Hash(data...)
}
func Pex(data ...[]byte) string {
	hash := Sha3(data...)
	return common.Bytes2Hex(hash[:8])
}

//pex for hash
func HPex(data common.Hash) string {
	hash := Sha3(data[:])
	return common.Bytes2Hex(hash[:8])
}
func APex(data common.Address) string {
	hash := Sha3(data[:])
	return common.Bytes2Hex(hash[:8])
}
func PubkeyToAddress(pubkey []byte) common.Address {
	return common.BytesToAddress(crypto.Keccak256(pubkey[1:])[12:])
}

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
