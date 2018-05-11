package rpc

import (
	"testing"

	"context"

	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/params"
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
	tokenAddress := common.HexToAddress("0xa9b61a3cc7cc1810e133174caa7ead7ef909d701")
	_, err := reg.AddToken(tokenAddress)
	if err != nil {
		t.Error(err)
		return
	}
	//no way to know the result of transaction ,is failure or success?
	manager, err := reg.ChannelManagerByToken(tokenAddress)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("manager address=%s", common.Bytes2Hex(manager[:]))
}

func TestGetAddTokenLog(t *testing.T) {
	bcs := MakeTestBlockChainService()
	logs, err := EventGetInternal(context.Background(), bcs.RegistryAddress, rpc.EarliestBlockNumber,
		rpc.LatestBlockNumber, "TokenAdded", RegistryABI, bcs.Client)
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
		rpc.LatestBlockNumber, "TokenAdded", RegistryABI, bcs.Client.Client, ch)
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case log := <-ch:
		spew.Dump(log)
		break
	case err = <-sub.Err():
		t.Error(err)
		break
	}
	sub.Unsubscribe()
}

func TestEventGetChannelNew(t *testing.T) {
	bcs := MakeTestBlockChainService()
	oneChannelManagerAddress := common.HexToAddress("0x2a00314c128855512ce77c16c839c7f263bbe99")
	logs, err := EventGetInternal(context.Background(), oneChannelManagerAddress, rpc.EarliestBlockNumber,
		rpc.LatestBlockNumber, params.NameChannelNew, ChannelManagerContractABI, bcs.Client)
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump(logs)
}

func TestEventAddressRegistered(t *testing.T) {
	bcs := MakeTestBlockChainService()
	logs, err := EventGetInternal(context.Background(), params.RopstenDiscoveryAddress, rpc.EarliestBlockNumber,
		rpc.LatestBlockNumber, params.NameAddressRegistered, EndpointRegistryABI, bcs.Client)
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump(logs)
}

func TestCodeAt(t *testing.T) {
	bcs := MakeTestBlockChainService()
	addrKilled := common.HexToAddress("0xad65d5b1210a80e8664aa58185bcd492184a43fa")
	addrNotExist := common.HexToAddress("0x0000000000000000000000000000000000000000")
	addrHasContract := params.RopstenRegistryAddress
	//
	code, err := bcs.Client.CodeAt(context.Background(), addrKilled, nil)
	if err != nil {
		t.Error(err)
	}
	if len(code) != 0 {
		t.Error("selfdestruct contract's code should be empty")
	}
	code, err = bcs.Client.CodeAt(context.Background(), addrNotExist, nil)
	if err != nil {
		t.Error(err)
	}
	if len(code) != 0 {
		t.Error("not exist account's code shoule be empty")
	}
	code, err = bcs.Client.CodeAt(context.Background(), addrHasContract, nil)
	if err != nil {
		t.Error(err)
	}
	if len(code) <= 0 {
		t.Error("this conctract 's code should not be empty")
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
	timeoutCh := time.After(time.Minute * 1)
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
