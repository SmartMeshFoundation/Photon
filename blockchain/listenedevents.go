package blockchain

import (
	"strings"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

var secretRegistryAbi abi.ABI
var tokenNetworkRegistryAbi abi.ABI
var tokenNetworkAbi abi.ABI

func init() {
	var err error
	secretRegistryAbi, err = abi.JSON(strings.NewReader(contracts.SecretRegistryABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	tokenNetworkAbi, err = abi.JSON(strings.NewReader(contracts.TokenNetworkABI))
	if err != nil {
		panic(fmt.Sprintf("tokenNetworkAbi parse err %s", err))
	}
	tokenNetworkRegistryAbi, err = abi.JSON(strings.NewReader(contracts.TokenNetworkRegistryABI))
	if err != nil {
		panic(fmt.Sprintf("tokenNetworkRegistryAbi parse err %s", err))
	}
}

func newEventTokenNetworkCreated(el *types.Log) (event *contracts.TokenNetworkRegistryTokenNetworkCreated, err error) {
	event = &contracts.TokenNetworkRegistryTokenNetworkCreated{}
	err = UnpackLog(&tokenNetworkRegistryAbi, event, params.NameTokenNetworkCreated, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventEventChannelOpen(el *types.Log) (event *contracts.TokenNetworkChannelOpened, err error) {
	event = &contracts.TokenNetworkChannelOpened{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelOpened, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
func newEventChannelNewBalance(el *types.Log) (event *contracts.TokenNetworkChannelNewDeposit, err error) {
	event = &contracts.TokenNetworkChannelNewDeposit{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelNewDeposit, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelClosed(el *types.Log) (event *contracts.TokenNetworkChannelClosed, err error) {
	event = &contracts.TokenNetworkChannelClosed{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelClosed, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelWithdraw(el *types.Log) (event *contracts.TokenNetworkChannelWithdraw, err error) {
	event = &contracts.TokenNetworkChannelWithdraw{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelWithdraw, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventBalanceProofUpdated(el *types.Log) (event *contracts.TokenNetworkBalanceProofUpdated, err error) {
	event = &contracts.TokenNetworkBalanceProofUpdated{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameBalanceProofUpdated, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelSettled(el *types.Log) (event *contracts.TokenNetworkChannelSettled, err error) {
	event = &contracts.TokenNetworkChannelSettled{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelSettled, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
func newEventSecretRevealed(el *types.Log) (event *contracts.SecretRegistrySecretRevealed, err error) {
	event = &contracts.SecretRegistrySecretRevealed{}
	err = UnpackLog(&secretRegistryAbi, event, params.NameSecretRevealed, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelCooperativeSettled(el *types.Log) (event *contracts.TokenNetworkChannelCooperativeSettled, err error) {
	event = &contracts.TokenNetworkChannelCooperativeSettled{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelCooperativeSettled, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
