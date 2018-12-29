package blockchain

import (
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/ethereum/go-ethereum/core/types"
)

func newEventTokenNetworkCreated(el *types.Log) (event *contracts.TokensNetworkTokenNetworkCreated, err error) {
	event = &contracts.TokensNetworkTokenNetworkCreated{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameTokenNetworkCreated, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelOpenAndDeposit(el *types.Log) (event *contracts.TokensNetworkChannelOpenedAndDeposit, err error) {
	event = &contracts.TokensNetworkChannelOpenedAndDeposit{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelOpenedAndDeposit, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
func newEventChannelNewDeposit(el *types.Log) (event *contracts.TokensNetworkChannelNewDeposit, err error) {
	event = &contracts.TokensNetworkChannelNewDeposit{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelNewDeposit, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelClosed(el *types.Log) (event *contracts.TokensNetworkChannelClosed, err error) {
	event = &contracts.TokensNetworkChannelClosed{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelClosed, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelWithdraw(el *types.Log) (event *contracts.TokensNetworkChannelWithdraw, err error) {
	event = &contracts.TokensNetworkChannelWithdraw{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelWithdraw, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
func newEventChannelUnlocked(el *types.Log) (event *contracts.TokensNetworkChannelUnlocked, err error) {
	event = &contracts.TokensNetworkChannelUnlocked{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelUnlocked, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
func newEventBalanceProofUpdated(el *types.Log) (event *contracts.TokensNetworkBalanceProofUpdated, err error) {
	event = &contracts.TokensNetworkBalanceProofUpdated{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameBalanceProofUpdated, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelSettled(el *types.Log) (event *contracts.TokensNetworkChannelSettled, err error) {
	event = &contracts.TokensNetworkChannelSettled{}
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

func newEventChannelCooperativeSettled(el *types.Log) (event *contracts.TokensNetworkChannelCooperativeSettled, err error) {
	event = &contracts.TokensNetworkChannelCooperativeSettled{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelCooperativeSettled, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelPunished(el *types.Log) (event *contracts.TokensNetworkChannelPunished, err error) {
	event = &contracts.TokensNetworkChannelPunished{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelPunished, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
