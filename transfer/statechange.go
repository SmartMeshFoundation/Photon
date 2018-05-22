package transfer

import (
	"encoding/gob"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/ethereum/go-ethereum/common"
)

//BlockStateChange used when a new block is mined.
type BlockStateChange struct {
	BlockNumber int64
}

/*
ActionRouteChangeStateChange is A route change.

    State change used for:
        - when a new route is added.
        - when the counter party is unresponsive (fails the healthcheck) and the
          route cannot be used.
        - when a different transfer uses the channel, changing the available
          balance.
*/
type ActionRouteChangeStateChange struct {
	Route      *RouteState
	Identifier uint64
}

/*
ActionCancelTransferStateChange The user requests the transfer to be cancelled.

    This state change can fail, it depends on the node's role and the current
    state of the transfer.
*/
type ActionCancelTransferStateChange struct {
	Identifier uint64
}

//ActionTransferDirectStateChange send a direct transfer
type ActionTransferDirectStateChange struct {
	Identifier   uint64
	Amount       *big.Int
	TokenAddress common.Address
	NodeAddress  common.Address
}

//ReceiveTransferDirectStateChange receive a direct transfer
type ReceiveTransferDirectStateChange struct {
	Identifier   uint64
	Amount       *big.Int
	TokenAddress common.Address
	Sender       common.Address
	Message      *encoding.DirectTransfer
}

func init() {
	gob.Register(&BlockStateChange{})
	gob.Register(&ActionRouteChangeStateChange{})
	gob.Register(&ActionCancelTransferStateChange{})
	gob.Register(&ActionTransferDirectStateChange{})
	gob.Register(&ReceiveTransferDirectStateChange{})
}
