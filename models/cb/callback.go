package cb

import (
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/ethereum/go-ethereum/common"
)

//NewTokenCb notify when new token registered
//return true to remove this callback, all the callback should never block.
type NewTokenCb func(token common.Address) (remove bool)

//ChannelCb notify when channel status changed
//return true to remove this callback, all the callback should never block.
type ChannelCb func(c *channeltype.Serialization) (remove bool)
