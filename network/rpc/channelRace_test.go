package rpc

import (
	"testing"

	"sync"

	"os"

	"github.com/ethereum/go-ethereum/common"
)

func TestChannelConcurrentQuery(t *testing.T) {
	bcs := MakeTestBlockChainService()
	tn, err := bcs.TokenNetwork(common.HexToAddress(os.Getenv("TOKENNETWORK")))
	if err != nil {
		t.Error(err)
		return
	}
	_, p1 := TestGetParticipant1()
	_, p2 := TestGetParticipant2()
	_, s, _, _, _, err := tn.GetChannelInfo(p1, p2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("settile: %d", s)
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			_, s2, _, _, _, _ := tn.GetChannelInfo(p1, p2)
			if s != s2 {
				t.Error("not equal")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
