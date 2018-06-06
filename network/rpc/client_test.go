package rpc

import (
	"testing"

	"context"

	"time"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func init() {
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = true
	spew.Config.MaxDepth = 7
}
func TestToken(t *testing.T) {
	bcs := MakeTestBlockChainService()
	reg := bcs.Registry(bcs.RegistryAddress)
	address, err := reg.TokenAddresses()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("address:%v", address)
}

func TestAddToken(t *testing.T) {
	bcs := MakeTestBlockChainService()
	reg := bcs.Registry(bcs.RegistryAddress)
	tokenAddress := utils.EmptyAddress
	_, err := reg.AddToken(tokenAddress)
	if err == nil {
		t.Errorf("should fail for invalid token")
		return
	}
}

func TestGetAddTokenLog(t *testing.T) {
	bcs := MakeTestBlockChainService()
	logs, err := EventGetInternal(context.Background(), bcs.RegistryAddress, rpc.EarliestBlockNumber,
		rpc.LatestBlockNumber, "TokenAdded", contracts.RegistryABI, bcs.Client)
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump(logs)
}
func TestEventSubscribe(t *testing.T) {
	bcs := MakeTestBlockChainService()
	ch := make(chan types.Log, 1)
	t.Log("wait for tokenadded event")
	sub, err := EventSubscribeInternal(context.Background(), bcs.RegistryAddress, rpc.EarliestBlockNumber,
		rpc.LatestBlockNumber, "TokenAdded", contracts.RegistryABI, bcs.Client.Client, ch)
	if err != nil {
		t.Error(err)
		return
	}
	//select {
	//case log := <-ch:
	//	spew.Dump(log)
	//	break
	//case err = <-sub.Err():
	//	t.Error(err)
	//	break
	//}
	sub.Unsubscribe()
}

func TestEventGetChannelNew(t *testing.T) {
	bcs := MakeTestBlockChainService()
	oneChannelManagerAddress := common.HexToAddress("0x2a00314c128855512ce77c16c839c7f263bbe99")
	logs, err := EventGetInternal(context.Background(), oneChannelManagerAddress, rpc.EarliestBlockNumber,
		rpc.LatestBlockNumber, params.NameChannelNew, contracts.ChannelManagerContractABI, bcs.Client)
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump(logs)
}

func TestCodeAt(t *testing.T) {
	bcs := MakeTestBlockChainService()
	addrNotExist := common.HexToAddress("0x0000000000000000000000000000000000000000")
	addrHasContract := common.HexToAddress(os.Getenv("REGISTRY"))

	code, err := bcs.Client.CodeAt(context.Background(), addrNotExist, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(code) != 0 {
		t.Error("not exist account's code shoule be empty")
		return
	}
	code, err = bcs.Client.CodeAt(context.Background(), addrHasContract, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(code) <= 0 {
		t.Error("this conctract 's code should not be empty")
		return
	}
}

func TestNewHead(t *testing.T) {
	fmt.Println("start...")
	bcs := MakeTestBlockChainService()
	ch := make(chan *types.Header, 1)
	sub, err := bcs.Client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		t.Error(err)
		return
	}
	timeoutCh := time.After(time.Second * 10)
	for {
		select {
		case h := <-ch:
			fmt.Printf("receive header:%d\n", h.Number.Int64())
		case <-timeoutCh:
			sub.Unsubscribe()
			return
		}
	}
}
