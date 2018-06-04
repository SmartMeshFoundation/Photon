package rpc

import (
	"testing"

	"sync"

	"os"

	"github.com/ethereum/go-ethereum/common"
)

func TestChannelConcurrentQuery(t *testing.T) {
	bcs := MakeTestBlockChainService()
	ch, err := bcs.NettingChannel(common.HexToAddress(os.Getenv("CHANNEL")))
	if err != nil {
		t.Error(err)
		return
	}
	s, _ := ch.SettleTimeout()
	t.Log("settile:", s)
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			s2, _ := ch.SettleTimeout()
			if s != s2 {
				t.Error("not equal")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
