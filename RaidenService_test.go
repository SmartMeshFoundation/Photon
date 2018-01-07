package raiden_network

import (
	"os"
	"testing"

	"time"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/ethereum/go-ethereum/log"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}
func TestPing(t *testing.T) {
	r1, r2, _ := makeTestRaidens()
	ping := encoding.NewPing(32)
	ping.Sign(r1.PrivateKey, ping)
	err := r1.SendAndWait(r2.NodeAddress, ping, time.Second)
	if err != nil {
		t.Error(err)
	}
}
