package transfer

import (
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
)

//Transition used when a new block is mined.
type BlockStateChange struct {
	BlockNumber int64
}

/*
A route change.

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
The user requests the transfer to be cancelled.

    This state change can fail, it depends on the node's role and the current
    state of the transfer.
*/
type ActionCancelTransferStateChange struct {
	Identifier uint64
}
type ActionTransferDirectStateChange struct {
	Identifier   uint64
	Amount       int64
	TokenAddress common.Address
	NodeAddress  common.Address
}

type ReceiveTransferDirectStateChange struct {
	Identifier   uint64
	Amount       int64
	TokenAddress common.Address
	Sender       common.Address
}

func init() {
	gob.Register(&BlockStateChange{})
	gob.Register(&ActionRouteChangeStateChange{})
	gob.Register(&ActionCancelTransferStateChange{})
	gob.Register(&ActionTransferDirectStateChange{})
	gob.Register(&ReceiveTransferDirectStateChange{})
}
