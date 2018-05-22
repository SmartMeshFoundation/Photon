package utest

import (
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
)

//UnitSettleTimeout for test only
var UnitSettleTimeout = 50

//UnitRevealTimeout for test
var UnitRevealTimeout = 5

//UnitTransferAmount for test
var UnitTransferAmount = big.NewInt(10)

//UnitBlockNumber for test
var UnitBlockNumber int64 = 1

//UnitIdentifier for test
var UnitIdentifier uint64 = 3

//UnitSecret for test
var UnitSecret = common.StringToHash("secretsecretsecretsecretsecretse")

//UnitHashLock for test
var UnitHashLock = utils.Sha3(UnitSecret[:])

//UnitTokenAddress for test
var UnitTokenAddress = utils.NewRandomAddress()

//ADDR for test
var ADDR = utils.NewRandomAddress()

//HOP1 for test
var HOP1 = common.HexToAddress("0x0101010101010101111111111111111111111111")

//HOP2 for test
var HOP2 = common.HexToAddress("0x0202020222222222222222222222222222222222")

//HOP3 for test
var HOP3 = common.HexToAddress("0x0303030303333333333333333333333333333333")

//HOP4 for test
var HOP4 = common.HexToAddress("0x0404040444444444444444444444444444444444")

//HOP5 for test
var HOP5 = common.HexToAddress("0x0505050505055555555555555555555555555555")

//HOP6 for test
var HOP6 = common.HexToAddress("0x060606060606666666666666666666666666666")

//Hop1Timeout for test
var Hop1Timeout = UnitSettleTimeout

//Hop2Timeout for test
var Hop2Timeout = Hop1Timeout - UnitRevealTimeout

//Hop3Timeout for test
var Hop3Timeout = Hop2Timeout - UnitRevealTimeout

/*
MakeRoute Helper for creating a route.

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

//MakeTransfer create test transfer
func MakeTransfer(amount *big.Int, initiator, target common.Address, expiration int64, secret common.Hash, hashlock common.Hash, identifier uint64, token /*UnitTokenAddress*/ common.Address) *mediatedtransfer.LockedTransferState {
	if secret != utils.EmptyHash {
		if utils.Sha3(secret[:]) != hashlock {
			log.Error("sha3(secret) != hashlock")
		}
	}
	if hashlock == utils.EmptyHash {
		hashlock = UnitHashLock
	}
	return &mediatedtransfer.LockedTransferState{
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

//MakeFrom create test from route and from transfer
func MakeFrom(amount *big.Int, target common.Address, fromExpiration int64, initiator /*HOP6*/ common.Address, secret common.Hash) (fromroute *transfer.RouteState, fromtransfer *mediatedtransfer.LockedTransferState) {
	fromroute = MakeRoute(initiator, amount, UnitSettleTimeout, UnitRevealTimeout, 0, utils.EmptyAddress)
	fromtransfer = MakeTransfer(amount, initiator, target, fromExpiration, secret, utils.EmptyHash, 0, UnitTokenAddress)
	return
}
