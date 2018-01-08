package rpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestBrokenClient(t *testing.T) {
	return
	bcs := MakeTestBlockChainService()
	_, err := bcs.Client.BalanceAt(context.Background(), bcs.NodeAddress, nil)
	if err != nil {
		t.Error(err)
		return

	}
	fmt.Println("shutdown geth now...")
	time.Sleep(5 * time.Second)
	fmt.Println("try  operation on broken connection")
	_, err = bcs.Client.BalanceAt(context.Background(), bcs.NodeAddress, nil)
	spew.Dump(err)
	t.Error(err)
}
