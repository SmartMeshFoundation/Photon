package contracts

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//ChannelIdentifier of contracts
type ChannelIdentifier common.Hash

//ChannelUniqueID unique id of a channel
type ChannelUniqueID struct {
	ChannelIdentifier common.Hash `json:"channel_identifier"`
	OpenBlockNumber   int64       `json:"open_block_number"`
}

func (c *ChannelUniqueID) String() string {
	return fmt.Sprintf("{ch=%s,OpenBlockNumber=%d}",
		utils.HPex(c.ChannelIdentifier), c.OpenBlockNumber)
}

const ChannelStateOpened = 1
const ChannelStateClosed = 2
const ChannelStateSettledOrNotExist = 0
