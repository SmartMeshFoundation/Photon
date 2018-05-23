package smartraiden

import (
	"os"
	"testing"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func TestPing(t *testing.T) {
	reinit()
	r1, r2, r3 := makeTestRaidens()
	defer r1.Stop()
	defer r2.Stop()
	defer r3.Stop()
	ping := encoding.NewPing(32)
	ping.Sign(r1.PrivateKey, ping)
	err := r1.SendAndWait(r2.NodeAddress, ping, time.Second)
	if err != nil {
		t.Error(err)
	}
}
