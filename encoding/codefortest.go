package encoding

import (
	"os"

	"math/rand"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
)

var isTest = false

func init() {
	isTest = len(os.Getenv("ISTEST")) > 0
}

//TestChannelBlockNumberGetter only valid in test,if was used in production environment, always error
type TestChannelBlockNumberGetter struct {
}

//GetChannelOpenBlockNumber only works in
func (c TestChannelBlockNumberGetter) GetChannelOpenBlockNumber(chID *contracts.ChannelUniqueID) int64 {
	if isTest {
		return 0
	}
	log.Warn(fmt.Sprintf("GetChannelOpenBlockNumber should only be called in test"))
	return rand.Int63()
}
