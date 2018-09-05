package utils

import (
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
