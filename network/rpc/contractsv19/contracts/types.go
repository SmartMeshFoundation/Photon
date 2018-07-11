package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//ChannelIdentifier of contracts
type ChannelIdentifier *big.Int

//ChannelUniqueID unique id of a channel
type ChannelUniqueID struct {
	TokenNetworkAddress common.Address
	ChannelIdentifier   ChannelIdentifier
}

const ChannelStateOpened = 1
const ChannelStateClosed = 2
const ChannelStateSettledOrNotExist = 0
