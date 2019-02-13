package utils

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

/*
HexToAddress convert hex encoded address string to common.Address,
and verify it is in accordance with EIP55
*/
func HexToAddress(addr string) (address common.Address, err error) {
	address = common.HexToAddress(addr)
	if address.String() != addr {
		err = fmt.Errorf("address checksum error expect=%s,got=%s", address.String(), addr)
	}
	return
}

//HexToAddressWithoutValidation disable EIP55 validation
func HexToAddressWithoutValidation(addr string) (address common.Address, err error) {
	return common.HexToAddress(addr), nil
}

//CalcChannelID 计算ChannelID的方式,注意与合约上计算方式保持完全一致.
func CalcChannelID(token, tokensNetwork, p1, p2 common.Address) common.Hash {
	var channelID common.Hash
	//log.Trace(fmt.Sprintf("p1=%s,p2=%s,tokennetwork=%s", p1.String(), p2.String(), tokenNetwork.String()))
	if bytes.Compare(p1[:], p2[:]) < 0 {
		channelID = Sha3(p1[:], p2[:], token[:], tokensNetwork[:])
	} else {
		channelID = Sha3(p2[:], p1[:], token[:], tokensNetwork[:])
	}
	return channelID
}
