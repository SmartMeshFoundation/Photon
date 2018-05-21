package rpc

import (
	"testing"

	"sync"

	"github.com/ethereum/go-ethereum/common"
)

func TestChannelConcurrentQuery(t *testing.T) {
	bcs := MakeTestBlockChainService()
	ch, err := bcs.NettingChannel(common.HexToAddress("0x2244d2509cbBFe7a8616ed86ab5c3ECe8FDC40Fe"))
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
