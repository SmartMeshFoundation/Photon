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
ActionCancelTransferStateChange The user requests the transfer to be cancelled.

    This state change can fail, it depends on the node's role and the current
    state of the transfer.
*/
type ActionCancelTransferStateChange struct {
	LockSecretHash common.Hash
}

//ActionTransferDirectStateChange send a direct transfer
type ActionTransferDirectStateChange struct {
	Amount       *big.Int
	TokenAddress common.Address
	NodeAddress  common.Address
}

//ReceiveTransferDirectStateChange receive a direct transfer
type ReceiveTransferDirectStateChange struct {
	Amount       *big.Int
	TokenAddress common.Address
	Sender       common.Address
	Message      *encoding.DirectTransfer
}

//CooperativeSettleStateChange user request to cooperative settle
type CooperativeSettleStateChange struct {
	Message *encoding.SettleRequest
}

//WithdrawRequestStateChange user request to withdraw on chain
type WithdrawRequestStateChange struct {
	Message *encoding.WithdrawRequest
}

/*
StopTransferRightNowStateChange 收到了 WithdrawRequest 或者 CooperativeSettleRequest, 应该理解停止进行中的交易.
*/
type StopTransferRightNowStateChange struct {
	TokenAddress      common.Address
	ChannelIdentifier common.Hash
}

func init() {
	gob.Register(&BlockStateChange{})
	gob.Register(&ActionCancelTransferStateChange{})
	gob.Register(&ActionTransferDirectStateChange{})
	gob.Register(&ReceiveTransferDirectStateChange{})
}
