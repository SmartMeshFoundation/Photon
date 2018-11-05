package models

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

// FeeSetting :
type FeeSetting struct {
	FeeConstant *big.Int `json:"fee_constant"`
	FeeRate     *big.Int `json:"fee_rate"`
}

// FeePolicy :
type FeePolicy struct {
	Key           string                         `storm:"id"`
	AccountFee    *FeeSetting                    `json:"account_fee"`
	TokenFeeMap   map[common.Address]*FeeSetting `json:"token_fee_map"`
	ChannelFeeMap map[common.Hash]*FeeSetting    `json:"channel_fee_map"`
}

const defaultKey string = "feePolicy"

// SaveFeePolicy :
func (model *ModelDB) SaveFeePolicy(fp *FeePolicy) (err error) {
	fp.Key = defaultKey
	err = model.db.Save(fp)
	return
}

// GetFeePolicy :
func (model *ModelDB) GetFeePolicy() (fp *FeePolicy) {
	fp = &FeePolicy{}
	err := model.db.One("Key", defaultKey, fp)
	if err == storm.ErrNotFound {
		return newDefaultFeePolicy()
	}
	if err != nil {
		log.Error(fmt.Sprintf("GetFeePolicy err %s, use default fee policy", err))
		return newDefaultFeePolicy()
	}
	return
}

// 默认手续费万分之一
func newDefaultFeePolicy() *FeePolicy {
	return &FeePolicy{
		AccountFee: &FeeSetting{
			FeeConstant: big.NewInt(0),
			FeeRate:     big.NewInt(10000),
		},
		TokenFeeMap:   make(map[common.Address]*FeeSetting),
		ChannelFeeMap: make(map[common.Hash]*FeeSetting),
	}
}
