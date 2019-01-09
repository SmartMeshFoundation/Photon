package rpc

import (
	"testing"

	"context"

	"time"

	"fmt"

	"os"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func init() {
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = true
	spew.Config.MaxDepth = 7
}

func TestCodeAt(t *testing.T) {
	bcs := MakeTestBlockChainService()
	addrNotExist := common.HexToAddress("0x0000000000000000000000000000000000000000")
	addrHasContract := common.HexToAddress(os.Getenv("TOKEN_NETWORK"))
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
	if true {
		return //不再强制要求必须是ws连接,
	}
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
	now := time.Now()
	for {
		if time.Since(now) > 20*time.Second {
			break
		}
		pendingNonce, _ := client.PendingNonceAt(context.Background(), account)
		nonce, _ := client.NonceAt(context.Background(), account, nil)
		fmt.Println("pendingNonce", pendingNonce)
		fmt.Println("nonce", nonce)
		fmt.Println("=============")
		time.Sleep(3 * time.Second)
	}
}
