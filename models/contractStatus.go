package models

import (
	"encoding/gob"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//ContractStatus 合约一旦部署就会确定下来的内容,也不会发生改变
type ContractStatus struct {
	RegistryAddress       common.Address
	ContractVersion       string
	SecretRegistryAddress common.Address
	PunishBlockNumber     int64
	ChainID               *big.Int
}

func init() {
	gob.Register(ContractStatus{})
}
