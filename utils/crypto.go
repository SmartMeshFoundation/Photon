package utils

import (
	"crypto/ecdsa"
	"io"

	"math/big"

	"encoding/hex"
	"fmt"

	"crypto/sha256"

	"github.com/SmartMeshFoundation/Photon/log"
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

//Ecrecover is a wrapper for crypto.Ecrecover
func Ecrecover(hash common.Hash, signature []byte) (addr common.Address, err error) {
	if len(signature) != 65 {
		err = fmt.Errorf("signature errr, len=%d,signature=%s", len(signature), hex.EncodeToString(signature))
		return
	}
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		signature[len(signature)-1] += 27
		return
	}
	addr = PubkeyToAddress(pubkey)
	signature[len(signature)-1] += 27
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

//ShaSecret is short for sha256
func ShaSecret(data []byte) common.Hash {
	//	return crypto.Keccak256Hash(data...)
	return sha256.Sum256(data)
}

//HPex pex for hash
func HPex(data common.Hash) string {
	return common.Bytes2Hex(data[:2])
}

//BPex bytes to string
func BPex(data []byte) string {
	return common.Bytes2Hex(data)
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

//ReadBigInt read big.Int from buffer
func ReadBigInt(reader io.Reader) *big.Int {
	bi := new(big.Int)
	tmpbuf := make([]byte, 32)
	_, err := reader.Read(tmpbuf)
	if err != nil {
		log.Error(fmt.Sprintf("read BigInt error %s", err))
	}
	bi.SetBytes(tmpbuf)
	return bi
}
