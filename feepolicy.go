package raiden_network

import (
	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
)

type NoFeePolicy struct {
}

func (n *NoFeePolicy) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	return utils.BigInt0
}

type ConstantFeePolicy struct {
}

var fixedFee = big.NewInt(3)

func (f *ConstantFeePolicy) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	return fixedFee
}

type CombinationFeePolicy struct {
}

func (c *CombinationFeePolicy) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	f := new(big.Int).Div(amount, big.NewInt(1000)) //fee rate: one in thousand.
	return f.Add(f, fixedFee)
}
