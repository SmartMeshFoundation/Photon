package fee

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

/*
Charger defines how to charge fee for a mediated node
it relates which token and how much to transfer.
*/
type Charger interface {
	//GetNodeChargeFee returns how many tokens charge for transfer 'amount' tokens on token who's address is tokenAddress.
	GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int
}
