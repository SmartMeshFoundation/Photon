package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// FeePolicy :
type FeePolicy struct {
	Key         common.Hash `storm:"id"`
	FeeConstant *big.Int
	FeeRate     *big.Int
}

func (model *ModelDB) GetFeePolicyByAccount() *FeePolicy {
}
