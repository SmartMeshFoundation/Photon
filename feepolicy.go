package smartraiden

import (
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
