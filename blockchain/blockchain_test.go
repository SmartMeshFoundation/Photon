package blockchain

import (
	"math/big"
	"os"

	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var client *helper.SafeEthClient
var auth *bind.TransactOpts
var ChainID *big.Int
var secretRegistryAddress common.Address
var TokenNetworkRegistryAddress common.Address
var be *Events

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	setup()
}

func setup() {
	var err error
	client, err = helper.NewSafeClient(rpc.TestRPCEndpoint)
	auth = bind.NewKeyedTransactor(rpc.TestPrivKey)
	if err != nil {
		panic(err)
	}
	TokenNetworkRegistryAddress = rpc.TestGetTokenNetworkRegistryAddress()
	tokenNetworkRegistry, err := contracts.NewTokenNetworkRegistry(TokenNetworkRegistryAddress, client)
	if err != nil {
		panic(err)
	}
	secretRegistryAddress, err = tokenNetworkRegistry.Secret_registry_address(nil)
	if err != nil {
		panic(err)
	}
	be = NewBlockChainEvents(client, TokenNetworkRegistryAddress, secretRegistryAddress)
}

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
	channels, err := be.GetAllChannels(0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("channels=%#v", channels)
}

func TestEvents_GetAllChannelClosed(t *testing.T) {
	events, err := be.GetAllChannelClosed(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events=\n%s", utils.StringInterface(events, 3))
}

func TestEvents_GetAllChannelSettled(t *testing.T) {
	events, err := be.GetAllChannelSettled(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events=\n%s", utils.StringInterface(events, 3))
}

func TestEvents_GetAllSecretRevealed(t *testing.T) {
	events, err := be.GetAllSecretRevealed(0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events=\n%s", utils.StringInterface(events, 3))
}
