package photon

import (
	"math/big"

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

var defaultFeePolicy = &models.FeePolicy{
	FeeConstant: big.NewInt(0),
	FeeRate:     big.NewInt(10000),
}

// FeeModule :
type FeeModule struct{}

func (c *FeeModule) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	var err error
	var feePolicy *models.FeePolicy
	feePolicy, err = models.ModelDB.GetFeePolicyByChannel(nodeAddress, tokenAddress)
	if err == nil {
		return calculateFee(feePolicy, amount)
	}
	feePolicy, err = models.ModelDB.GetFeePolicyByToken(tokenAddress)
	if err == nil {
		return calculateFee(feePolicy, amount)
	}
	feePolicy = models.ModelDB.GetFeePolicyByAccount()
	if feePolicy != nil {
		return calculateFee(feePolicy, amount)
	}
	return calculateFee(defaultFeePolicy, amount)
}

func calculateFee(feePolicy *models.FeePolicy, amount *big.Int) *big.Int {
	fee := big.NewInt(0)
	if feePolicy.FeeRate.Cmp(big.NewInt(0)) > 0 {
		fee = fee.Div(amount, feePolicy.FeeRate)
	}
	if feePolicy.FeeConstant.Cmp(big.NewInt(0)) > 0 {
		fee = fee.Add(fee, feePolicy.FeeConstant)
	}
	return fee
}
