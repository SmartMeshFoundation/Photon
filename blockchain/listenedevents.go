package blockchain

import (
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/ethereum/go-ethereum/core/types"
)

func newEventTokenNetworkCreated(el *types.Log) (event *contracts.TokenNetworkRegistryTokenNetworkCreated, err error) {
	event = &contracts.TokenNetworkRegistryTokenNetworkCreated{}
	err = UnpackLog(&tokenNetworkRegistryAbi, event, params.NameTokenNetworkCreated, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}

func newEventChannelOpen(el *types.Log) (event *contracts.TokenNetworkChannelOpened, err error) {
	event = &contracts.TokenNetworkChannelOpened{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelOpened, el)
	if err != nil {
		return
	}
	event.Raw = *el
	//log.Trace(fmt.Sprintf("newEventChannelOpen el=%s, event=%s", utils.StringInterface(el, 3), utils.StringInterface(event, 3)))
	return
}
func newEventChannelOpenAndDeposit(el *types.Log) (event *contracts.TokenNetworkChannelOpenedAndDeposit, err error) {
	event = &contracts.TokenNetworkChannelOpenedAndDeposit{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelOpenedAndDeposit, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
func newEventChannelNewDeposit(el *types.Log) (event *contracts.TokenNetworkChannelNewDeposit, err error) {
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
func newEventChannelUnlocked(el *types.Log) (event *contracts.TokenNetworkChannelUnlocked, err error) {
	event = &contracts.TokenNetworkChannelUnlocked{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelUnlocked, el)
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

func newEventChannelPunished(el *types.Log) (event *contracts.TokenNetworkChannelPunished, err error) {
	event = &contracts.TokenNetworkChannelPunished{}
	err = UnpackLog(&tokenNetworkAbi, event, params.NameChannelPunished, el)
	if err != nil {
		return
	}
	event.Raw = *el
	return
}
