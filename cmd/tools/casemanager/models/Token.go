package models

import (
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/ethereum/go-ethereum/common"
)

// Token name and address
type Token struct {
	Name                string
	Token               *contracts.Token
	TokenAddress        common.Address
	TokenNetwork        *contracts.TokenNetwork
	TokenNetworkAddress common.Address
}
