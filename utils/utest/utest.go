package utest

import (
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
)

var UnitSettleTimeout = 50
var UnitRevealTimeout = 5
var UnitTransferAmount = big.NewInt(10)
var UnitBlockNumber int64 = 1
var UnitIdentifier uint64 = 3
var UnitSecret = common.StringToHash("secretsecretsecretsecretsecretse")
var UnitHashLock = utils.Sha3(UnitSecret[:])

var UnitTokenAddress = utils.NewRandomAddress()

var ADDR = utils.NewRandomAddress()
var HOP1 = common.HexToAddress("0x0101010101010101111111111111111111111111")
var HOP2 = common.HexToAddress("0x0202020222222222222222222222222222222222")
var HOP3 = common.HexToAddress("0x0303030303333333333333333333333333333333")
var HOP4 = common.HexToAddress("0x0404040444444444444444444444444444444444")
var HOP5 = common.HexToAddress("0x0505050505055555555555555555555555555555")
var HOP6 = common.HexToAddress("0x060606060606666666666666666666666666666")

var Hop1Timeout = UnitSettleTimeout
var Hop2Timeout = Hop1Timeout - UnitRevealTimeout
var Hop3Timeout = Hop2Timeout - UnitRevealTimeout

/*
Helper for creating a route.

    Args:
        node_address (address): The node address.
        available_balance (int): The available capacity of the route.
        settle_timeout (int): The settle_timeout of the route, as agreed in the netting contract.
        reveal_timeout (int): The configure reveal_timeout of the raiden node.
        channel_address (address): The correspoding channel address.
*/
func MakeRoute(nodeAddress common.Address, availableBalance *big.Int, settleTimeout /*UnitSettleTimeout*/ int, revealTimeout /*UnitRevealTimeout*/ int, closedBlock int64, channelAddress common.Address) *transfer.RouteState {
	return &transfer.RouteState{
		State:          transfer.ChannelStateOpened,
		HopNode:        nodeAddress,
		ChannelAddress: channelAddress,
		AvaibleBalance: new(big.Int).Set(availableBalance),
		SettleTimeout:  settleTimeout,
		RevealTimeout:  revealTimeout,
		ClosedBlock:    closedBlock,
	}
}

func MakeTransfer(amount *big.Int, initiator, target common.Address, expiration int64, secret common.Hash, hashlock common.Hash, identifier uint64, token /*UnitTokenAddress*/ common.Address) *mediated_transfer.LockedTransferState {
	if secret != utils.EmptyHash {
		if utils.Sha3(secret[:]) != hashlock {
			log.Error("sha3(secret) != hashlock")
		}
	}
	if hashlock == utils.EmptyHash {
		hashlock = UnitHashLock
	}
	return &mediated_transfer.LockedTransferState{
		Identifier:   identifier,
		TargetAmount: new(big.Int).Set(amount),
		Amount:       new(big.Int).Set(amount),
		Token:        token,
		Initiator:    initiator,
		Target:       target,
		Expiration:   expiration,
		Hashlock:     hashlock,
		Secret:       secret,
		Fee:          utils.BigInt0,
	}
}
func MakeFrom(amount *big.Int, target common.Address, fromExpiration int64, initiator /*HOP6*/ common.Address, secret common.Hash) (fromroute *transfer.RouteState, fromtransfer *mediated_transfer.LockedTransferState) {
	fromroute = MakeRoute(initiator, amount, UnitSettleTimeout, UnitRevealTimeout, 0, utils.EmptyAddress)
	fromtransfer = MakeTransfer(amount, initiator, target, fromExpiration, secret, utils.EmptyHash, 0, UnitTokenAddress)
	return
}
