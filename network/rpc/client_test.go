package rpc

import (
	"testing"

	"context"

	"time"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/codefortest"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func init() {
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = true
	spew.Config.MaxDepth = 7
}

func TestAddToken(t *testing.T) {
	bcs := MakeTestBlockChainService()
	reg := bcs.Registry(bcs.RegistryProxy.Address, true)
	tokenAddress := utils.EmptyAddress
	_, err := reg.AddToken(tokenAddress)
	if err == nil {
		t.Errorf("should fail for invalid token")
		return
	}
}

func TestCodeAt(t *testing.T) {
	bcs := MakeTestBlockChainService()
	addrNotExist := common.HexToAddress("0x0000000000000000000000000000000000000000")
	addrHasContract := common.HexToAddress(os.Getenv("REGISTRY"))
	t.Logf("token network registry=%s\n", addrHasContract.String())
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

func TestPendingNonceAt(t *testing.T) {
	accounts, err := codefortest.GetAccounts()
	if err != nil {
		panic(err)
	}
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	account := accounts[0].Address
	for {
		pendingNonce, _ := client.PendingNonceAt(context.Background(), account)
		nonce, _ := client.NonceAt(context.Background(), account, nil)
		fmt.Println("pendingNonce", pendingNonce)
		fmt.Println("nonce", nonce)
		fmt.Println("=============")
		time.Sleep(3 * time.Second)
	}
}
