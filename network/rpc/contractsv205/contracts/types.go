package contracts

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

//ChannelIdentifier of contracts
type ChannelIdentifier common.Hash

//ChannelUniqueID unique id of a channel
type ChannelUniqueID struct {
	TokenNetworkAddress common.Address
	ChannelIdentifier   common.Hash
	OpenBlockNumber     int64
}

func (c ChannelUniqueID) String() string {
	return fmt.Sprintf("{tokenNetwork=%s,Channel=%s,OpenBlockNumber=%d}",
		utils.APex2(c.TokenNetworkAddress), utils.HPex(c.ChannelIdentifier), c.OpenBlockNumber)
}

const ChannelStateOpened = 1
const ChannelStateClosed = 2
const ChannelStateSettledOrNotExist = 0
