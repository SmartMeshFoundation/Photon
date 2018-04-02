package fee

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type FeeCharger interface {
	//returns how many tokens charge for transfer 'amount' tokens on token who's address is tokenAddress.
	GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int
}
