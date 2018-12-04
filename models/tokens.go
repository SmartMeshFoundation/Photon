package models

import (
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
)

//AddressMap is token address to mananger address
type AddressMap map[common.Address]common.Address

func init() {
	gob.Register(common.Address{})
	gob.Register(make(AddressMap))
}
