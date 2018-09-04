package rpc

import (
	"errors"
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//RegistryProxy proxy for registry contract
type RegistryProxy struct {
	Address  common.Address //contract address
	bcs      *BlockChainService
	registry *contracts.TokenNetworkRegistry
}

// TokenNetworkByToken Get the ChannelManager address for a specific token
// @param token_address The address of the given token
// @return Address of tokenNetwork
func (r *RegistryProxy) TokenNetworkByToken(tokenAddress common.Address) (tokenNetworkAddress common.Address, err error) {
	tokenNetworkAddress, err = r.registry.TokenToTokenNetworks(r.bcs.getQueryOpts(), tokenAddress)
	if tokenNetworkAddress == utils.EmptyAddress {
		err = rerr.ErrNoTokenManager
	}
	return
}

//GetContract return contract itself
func (r *RegistryProxy) GetContract() *contracts.TokenNetworkRegistry {
	return r.registry
}

//AddToken register a new token,this token must be a valid erc20
func (r *RegistryProxy) AddToken(tokenAddress common.Address) (tokenNetworkAddress common.Address, err error) {
	tx, err := r.registry.CreateERC20TokenNetwork(r.bcs.Auth, tokenAddress)
	if err != nil {
		return
	}
	receipt, err := bind.WaitMined(GetCallContext(), r.bcs.Client, tx)
	if err != nil {
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("AddToken failed %s,receipt=%s", utils.APex(r.Address), receipt))
		err = errors.New("AddToken tx execution failed")
		return
	}
	log.Info(fmt.Sprintf("AddToken success %s,token=%s, gasused=%d\n", utils.APex(r.Address), tokenAddress.String(), receipt.GasUsed))
	return r.registry.TokenToTokenNetworks(r.bcs.getQueryOpts(), tokenAddress)
}
