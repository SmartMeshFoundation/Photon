package blockchain

import (
	"fmt"
	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

var client *helper.SafeEthClient
var secretRegistryAddress common.Address
var TokenNetworkRegistryAddress common.Address
var be *Events
var at *AlarmTask

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	setup()
}

func setup() {
	var err error
	client, err = helper.NewSafeClient(rpc.TestRPCEndpoint)
	if err != nil {
		panic(err)
	}
	TokenNetworkRegistryAddress = rpc.TestGetTokenNetworkRegistryAddress()
	tokenNetworkRegistry, err := contracts.NewTokenNetworkRegistry(TokenNetworkRegistryAddress, client)
	if err != nil {
		panic(err)
	}
	secretRegistryAddress, err = tokenNetworkRegistry.SecretRegistryAddress(nil)
	if err != nil {
		panic(err)
	}
	be = NewBlockChainEvents(client, &fakaRPCModule{
		RegistryAddress:       TokenNetworkRegistryAddress,
		SecretRegistryAddress: secretRegistryAddress,
	}, nil)
	tokens, err := be.GetAllTokenNetworks(0)
	if err != nil {
		panic(fmt.Sprintf("cannot get all token networks err %s", err))
	}
	if len(tokens) == 0 {
		panic(fmt.Sprintf("empty registyr network"))
	}
	at = NewAlarmTask(client)
}

type fakaRPCModule struct {
	RegistryAddress       common.Address
	SecretRegistryAddress common.Address
}

func (r *fakaRPCModule) GetRegistryAddress() common.Address {
	return r.RegistryAddress
}

func (r *fakaRPCModule) GetSecretRegistryAddress() common.Address {
	return r.SecretRegistryAddress
}
