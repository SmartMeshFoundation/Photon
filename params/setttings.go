package params

import (
	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

const InitialPort = 40001

const CacheTTL = 60
const EstimatedBlockTime = 7
const GasLimit = 3141592 //den's gasLimit.

const GasPrice = params.Shannon * 20

const DefaultProtocolRetiesBeforeBackoff = 5
const DefaultProtocolRhrottleCapacity = 10.
const DefaultProtocolThrottleFillRate = 10.
const DefaultprotocolRetryInterval = 1.

const DefaultRevealTimeout = 3
const DefaultSettleTimeout = DefaultRevealTimeout * 9
const DefaultPollTimeout = 180 * time.Second
const DefaultJoinableFundsTarget = 0.4
const DefaultInitialChannelTarget = 3
const DefaultWaitForSettle = true

const DefaultKeepAliveReties = 5
const DefaultNATKeepAliveTimeout = 500
const DefaultNATInvitationTimeout = 15000
const DefaultTxTimeout = 5 * time.Minute   //15seconds for one block,it may take sever minutes
const MaxRequestTimeout = 20 * time.Minute //longest time for a request ,for example ,settle all channles?

var GasLimitHex string

var RopstenRegistryAddress = common.HexToAddress("0xd01Ca23F2B84AF393550271bFCC2A8b48d6f65b8")
var RopstenDiscoveryAddress = common.HexToAddress("0x4CDAF98516490d42E1E6F050bcfBD143dCb58CcD")

const NettingChannelSettleTimeoutMin = 6

/*
 The maximum settle timeout is chosen as something above
 1 year with the assumption of very fast block times of 12 seconds.
 There is a maximum to avoidpotential overflows as described here:
 https://github.com/SmartRaiden/raiden/issues/1038
*/
const NettingChannelSettleTimeoutMax = 2700000

const UDPMaxMessageSize = 1200
const DefaultSignalServer = "139.199.6.114:5222"
const DefaultTurnServer = "182.254.155.208:3478"

func init() {
	GasLimitHex = fmt.Sprintf("0x%x", GasLimit)
}
