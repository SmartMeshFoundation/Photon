package utest

import (
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
)

var UNIT_SETTLE_TIMEOUT = 50
var UNIT_REVEAL_TIMEOUT = 5
var UNIT_TRANSFER_AMOUNT int64 = 10
var UNIT_BLOCK_NUMBER int64 = 1
var UNIT_IDENTIFIER uint64 = 3
var UNIT_SECRET = common.StringToHash("secretsecretsecretsecretsecretse")
var UNIT_HASHLOCK = utils.Sha3(UNIT_SECRET[:])

var UNIT_TOKEN_ADDRESS = utils.NewRandomAddress()

var ADDR = utils.NewRandomAddress()
var HOP1 = common.HexToAddress("0x0101010101010101111111111111111111111111")
var HOP2 = common.HexToAddress("0x0202020222222222222222222222222222222222")
var HOP3 = common.HexToAddress("0x0303030303333333333333333333333333333333")
var HOP4 = common.HexToAddress("0x0404040444444444444444444444444444444444")
var HOP5 = common.HexToAddress("0x0505050505055555555555555555555555555555")
var HOP6 = common.HexToAddress("0x060606060606666666666666666666666666666")

var HOP1_TIMEOUT = UNIT_SETTLE_TIMEOUT
var HOP2_TIMEOUT = HOP1_TIMEOUT - UNIT_REVEAL_TIMEOUT
var HOP3_TIMEOUT = HOP2_TIMEOUT - UNIT_REVEAL_TIMEOUT

/*
Helper for creating a route.

    Args:
        node_address (address): The node address.
        available_balance (int): The available capacity of the route.
        settle_timeout (int): The settle_timeout of the route, as agreed in the netting contract.
        reveal_timeout (int): The configure reveal_timeout of the raiden node.
        channel_address (address): The correspoding channel address.
*/
func MakeRoute(nodeAddress common.Address, availableBalance int64, settleTimeout /*UNIT_SETTLE_TIMEOUT*/ int, revealTimeout /*UNIT_REVEAL_TIMEOUT*/ int, closedBlock int64, channelAddress common.Address) *transfer.RouteState {
	return &transfer.RouteState{
		State:          transfer.CHANNEL_STATE_OPENED,
		HopNode:        nodeAddress,
		ChannelAddress: channelAddress,
		AvaibleBalance: availableBalance,
		SettleTimeout:  settleTimeout,
		RevealTimeout:  revealTimeout,
		ClosedBlock:    closedBlock,
	}
}

func MakeTransfer(amount int64, initiator, target common.Address, expiration int64, secret common.Hash, hashlock common.Hash, identifier uint64, token /*UNIT_TOKEN_ADDRESS*/ common.Address) *mediated_transfer.LockedTransferState {
	if secret != utils.EmptyHash {
		if utils.Sha3(secret[:]) != hashlock {
			log.Error("sha3(secret) != hashlock")
		}
	}
	if hashlock == utils.EmptyHash {
		hashlock = UNIT_HASHLOCK
	}
	return &mediated_transfer.LockedTransferState{
		Identifier: identifier,
		Amount:     amount,
		Token:      token,
		Initiator:  initiator,
		Target:     target,
		Expiration: expiration,
		Hashlock:   hashlock,
		Secret:     secret,
	}
}
func MakeFrom(amount int64, target common.Address, fromExpiration int64, initiator /*HOP6*/ common.Address, secret common.Hash) (fromroute *transfer.RouteState, fromtransfer *mediated_transfer.LockedTransferState) {
	fromroute = MakeRoute(initiator, amount, UNIT_SETTLE_TIMEOUT, UNIT_REVEAL_TIMEOUT, 0, utils.EmptyAddress)
	fromtransfer = MakeTransfer(amount, initiator, target, fromExpiration, secret, utils.EmptyHash, 0, UNIT_TOKEN_ADDRESS)
	return
}
