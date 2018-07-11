package testfilter

import (
	"testing"

	"os"

	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenNetworkRegistryTokenNetworkCreated struct {
	Token_address         common.Address
	Token_network_address common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

func TestStruct(t *testing.T) {
	tt := TokenNetworkRegistryTokenNetworkCreated{
		Token_address: utils.NewRandomAddress(),
	}
	t2 := tt
	t2.Token_address = utils.NewRandomAddress()
	t.Logf(fmt.Sprintf("tt=%s,t2=%s", tt.Token_address.String(), t2.Token_address.String()))
}
func TestFilter(t *testing.T) {
	client, err := ethclient.Dial(rpc.TestRPCEndpoint)
	if err != nil {
		t.Error(err)
		return
	}
	registryAddr := common.HexToAddress(os.Getenv("REGISTRY"))
	//registry, err := contracts.NewTokenNetworkRegistry(registryAddr, client)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	f, err := contracts.NewTokenNetworkRegistryFilterer(registryAddr, client)
	if err != nil {
		t.Error(err)
		return
	}
	it, err := f.FilterTokenNetworkCreated(nil, nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	for it.Next() {
		fmt.Printf("event=%s", utils.StringInterface(it.Event, 3))
	}
	ch := make(chan *contracts.TokenNetworkRegistryTokenNetworkCreated, 10)
	var start uint64
	sub, err := f.WatchTokenNetworkCreated(&bind.WatchOpts{
		Start: &start,
	}, ch, nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		select {
		case <-time.After(10 * time.Second):
			sub.Unsubscribe()
			return
		case e := <-ch:
			fmt.Printf("sub event=%s", utils.StringInterface(e, 3))
		}
	}
}
