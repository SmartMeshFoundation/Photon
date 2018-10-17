package blockchain

import (
	"testing"

	"fmt"

	"time"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/codefortest"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

type fakeRPCModule struct {
	RegistryAddress       common.Address
	SecretRegistryAddress common.Address
}

func (r *fakeRPCModule) GetRegistryAddress() common.Address {
	return r.RegistryAddress
}

func (r *fakeRPCModule) GetSecretRegistryAddress() common.Address {
	return r.SecretRegistryAddress
}

func TestNewBlockChainEvents(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	be := NewBlockChainEvents(client, &fakeRPCModule{}, nil)
	if be == nil {
		t.Error("NewBlockChainEvents failed")
	}
}

func TestEvents_Start(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	be := NewBlockChainEvents(client, &fakeRPCModule{
		RegistryAddress: common.HexToAddress("0x7B25494cF297D63eA2AF72d43Fc133408674c43a"),
	}, nil)
	if be == nil {
		t.Error("NewBlockChainEvents failed")
	}
	params.ChainID = big.NewInt(8888)
	be.Start(-1)
	begin := time.Now()
	for {
		if time.Since(begin) > 10*time.Second {
			be.Stop()
			time.Sleep(5 * time.Second)
			return
		}
		select {
		case sc := <-be.StateChangeChannel:
			fmt.Println(utils.StringInterface(sc, 0))
			//BlockStateChange, ok := sc.(transfer.BlockStateChange)
			//if ok {
			//	fmt.Println(BlockStateChange.BlockNumber)
			//}
		}
	}
}
