package encoding

import (
	"os"

	"math/rand"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
)

//IsTest for test function only
var IsTest = false

func init() {
	IsTest = len(os.Getenv("ISTEST")) > 0
}

//TestChannelBlockNumberGetter only valid in test,if was used in production environment, always error
type TestChannelBlockNumberGetter struct {
}

//GetChannelOpenBlockNumber only works in
func (c TestChannelBlockNumberGetter) GetChannelOpenBlockNumber(chID *contracts.ChannelUniqueID) int64 {
	if IsTest {
		return 0
	}
	log.Warn(fmt.Sprintf("GetChannelOpenBlockNumber should only be called in test"))
	return rand.Int63()
}
