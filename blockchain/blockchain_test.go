package blockchain

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestGetTokenNetworkCreated(t *testing.T) {
	//NewBlockChainEvents create BlockChainEvents
	tokens, err := be.GetAllTokenNetworks(0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("tokens=%#v", tokens)
}

func TestEvents_GetAllChannels(t *testing.T) {
	channels, err := be.GetChannelNew(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("channels num =%d\n", len(channels))
}

func TestEvents_GetAllChannelClosed(t *testing.T) {
	events, err := be.GetChannelClosed(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events num =%d\n", len(events))
}

func TestEvents_GetAllChannelSettled(t *testing.T) {
	events, err := be.GetChannelSettled(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events num =%d\n", len(events))
}

func TestEvents_GetAllSecretRevealed(t *testing.T) {
	events, err := be.GetAllSecretRevealed(0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events num =%d\n", len(events))
}

func TestEvents_GetChannelNewAndDeposit(t *testing.T) {
	events, err := be.GetChannelNewAndDeposit(0, utils.EmptyAddress)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events num =%d\n", len(events))
}
