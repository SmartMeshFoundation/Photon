package contracts

import "github.com/ethereum/go-ethereum/common"

//ChannelIdentifier of contracts
type ChannelIdentifier common.Hash

//ChannelUniqueID unique id of a channel
type ChannelUniqueID struct {
	TokenNetworkAddress common.Address
	ChannelIdentifier   ChannelIdentifier
}

const ChannelStateOpened = 1
const ChannelStateClosed = 2
const ChannelStateSettledOrNotExist = 0
