package photon

import (
	"math/big"

	"sync"

	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//NoFeePolicy charge no fee
type NoFeePolicy struct {
}

//GetNodeChargeFee always return 0
func (n *NoFeePolicy) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	return utils.BigInt0
}

//ConstantFeePolicy charge a constant fee
type ConstantFeePolicy struct {
}

var fixedFee = big.NewInt(3)

//GetNodeChargeFee returnx fixedFee
func (f *ConstantFeePolicy) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	return fixedFee
}

//CombinationFeePolicy should not used now
type CombinationFeePolicy struct {
}

//GetNodeChargeFee should not used now
func (c *CombinationFeePolicy) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	f := new(big.Int).Div(amount, big.NewInt(1000)) //fee rate: one in thousand.
	return f.Add(f, fixedFee)
}

// FeeModule :
type FeeModule struct {
	db        *models.ModelDB
	feePolicy *models.FeePolicy
	lock      sync.Mutex
}

// NewFeeModule :
func NewFeeModule(db *models.ModelDB) *FeeModule {
	if db == nil {
		panic("need init db first")
	}
	return &FeeModule{
		db:        db,
		feePolicy: db.GetFeePolicy(),
	}
}

// SetFeePolicy :
func (fm *FeeModule) SetFeePolicy() (err error) {
	return
}

//GetNodeChargeFee : impl of FeeCharge
func (fm *FeeModule) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	var feeSetting *models.FeeSetting
	var ok bool
	// 优先channel
	c, err := fm.db.GetChannel(tokenAddress, nodeAddress)
	if c != nil && err == nil {
		feeSetting, ok = fm.feePolicy.ChannelFeeMap[c.ChannelIdentifier.ChannelIdentifier]
		if ok {
			return calculateFee(feeSetting, amount)
		}
	}
	// 其次token
	feeSetting, ok = fm.feePolicy.TokenFeeMap[tokenAddress]
	if ok {
		return calculateFee(feeSetting, amount)
	}
	// 最后account
	return calculateFee(fm.feePolicy.AccountFee, amount)
}

func calculateFee(feeSetting *models.FeeSetting, amount *big.Int) *big.Int {
	fee := big.NewInt(0)
	if feeSetting.FeeRate.Cmp(big.NewInt(0)) > 0 {
		fee = fee.Div(amount, feeSetting.FeeRate)
	}
	if feeSetting.FeeConstant.Cmp(big.NewInt(0)) > 0 {
		fee = fee.Add(fee, feeSetting.FeeConstant)
	}
	return fee
}
