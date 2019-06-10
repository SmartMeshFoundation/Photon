package models

import (
	"fmt"
	"math/big"

	"crypto/ecdsa"

	"bytes"
	"encoding/binary"
	"encoding/gob"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// FeeSetting :
// 其中FeeConstant为固定费率,比如5代表手续费固定部分为5个token,设置为0即不收费
// FeePercent为比例费率,计算方式为 交易金额/FeePercent,比如交易金额50000,FeePercent=10000,那么手续费比例部分=50000/10000=5,设置为0即不收费
// 最终为手续费为固定收费+比例收费
type FeeSetting struct {
	FeeConstant *big.Int `json:"fee_constant"`
	FeePercent  int64    `json:"fee_percent"`
	Signature   []byte   `json:"signature"` // used when set fee policy to pfs
}

func (fs *FeeSetting) sign(key *ecdsa.PrivateKey) []byte {
	var err error
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, fs.FeePercent)
	_, err = buf.Write(utils.BigIntTo32Bytes(fs.FeeConstant))
	if err != nil {
		log.Error(fmt.Sprintf("signData err %s", err))
	}
	fs.Signature, err = utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("signDataFor FeeSetting err %s", err))
	}
	return fs.Signature
}

// FeePolicy 定义个账号的收费策略
type FeePolicy struct {
	Key           string                         `storm:"id"`
	AccountFee    *FeeSetting                    `json:"account_fee"`
	TokenFeeMap   map[common.Address]*FeeSetting `json:"token_fee_map"`
	ChannelFeeMap map[common.Hash]*FeeSetting    `json:"channel_fee_map"`
}

// Sign for pfs
func (fp *FeePolicy) Sign(key *ecdsa.PrivateKey) {
	fp.AccountFee.sign(key)
	for _, fs := range fp.TokenFeeMap {
		fs.sign(key)
	}
	for _, fs := range fp.ChannelFeeMap {
		fs.sign(key)
	}
}

const defaultKey string = "feePolicy"

// NewDefaultFeePolicy : 默认手续费万分之一
func NewDefaultFeePolicy() *FeePolicy {
	return &FeePolicy{
		AccountFee: &FeeSetting{
			FeeConstant: big.NewInt(0),
			FeePercent:  10000,
		},
		TokenFeeMap:   make(map[common.Address]*FeeSetting),
		ChannelFeeMap: make(map[common.Hash]*FeeSetting),
	}
}

func init() {
	gob.Register(&FeePolicy{})
}
