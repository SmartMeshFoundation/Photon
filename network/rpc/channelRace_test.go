package rpc

import (
	"testing"

	"sync"

	"github.com/ethereum/go-ethereum/common"
)

func TestChannelConcurrentQuery(t *testing.T) {
	bcs := MakeTestBlockChainService()
	ch, err := bcs.NettingChannel(common.HexToAddress("0xf029c3ec22b5dc7194dfa2650ae701e57068781e"))
	if err != nil {
		t.Error(err)
		return
	}
	s, _ := ch.SettleTimeout()
	t.Log("settile:", s)
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
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
