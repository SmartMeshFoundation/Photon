package contracttest

import "github.com/ethereum/go-ethereum/common"

// TestSettleTimeoutMin :
var TestSettleTimeoutMin uint64 = 6

// TestSettleTimeoutMax :
var TestSettleTimeoutMax uint64 = 2700000

// FakeAccountAddress :
var FakeAccountAddress = common.HexToAddress("0x03432")

// EmptyAccountAddress :
var EmptyAccountAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// ChannelStateOpened :
const ChannelStateOpened uint8 = 1

// ChannelStateClosed :
const ChannelStateClosed uint8 = 2

// ChannelStateSettledOrNotExist :
const ChannelStateSettledOrNotExist uint8 = 0

// EmptyBalanceHash :
const EmptyBalanceHash = "000000000000000000000000000000000000000000000000"
