package params

import (
	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

//InitialPort listening port for communication bewtween nodes
const InitialPort = 40001

//GasLimit max gas usage for raiden tx
const GasLimit = 3141592 //den's gasLimit.
//GasPrice from ethereum
const GasPrice = params.Shannon * 20

//defaultProtocolRetiesBeforeBackoff
const defaultProtocolRetiesBeforeBackoff = 5
const defaultProtocolRhrottleCapacity = 10.
const defaultProtocolThrottleFillRate = 10.
const defaultprotocolRetryInterval = 1.

//DefaultRevealTimeout blocks needs to update transfer
const DefaultRevealTimeout = 3

//DefaultSettleTimeout settle time of channel
const DefaultSettleTimeout = DefaultRevealTimeout * 9

//DefaultPollTimeout  request wait time
const DefaultPollTimeout = 180 * time.Second

//DefaultJoinableFundsTarget for connection api
const DefaultJoinableFundsTarget = 0.4

//DefaultInitialChannelTarget channels to create
const DefaultInitialChannelTarget = 3

//DefaultKeepAliveReties args i don't know
const DefaultKeepAliveReties = 5

//DefaultNATKeepAliveTimeout args
const DefaultNATKeepAliveTimeout = 500

//DefaultNATInvitationTimeout args
const DefaultNATInvitationTimeout = 15000

//DefaultTxTimeout args
const DefaultTxTimeout = 5 * time.Minute //15seconds for one block,it may take sever minutes
//MaxRequestTimeout args
const MaxRequestTimeout = 20 * time.Minute //longest time for a request ,for example ,settle all channles?

var gasLimitHex string

//RopstenRegistryAddress Registry contract address
var RopstenRegistryAddress = common.HexToAddress("0xd01Ca23F2B84AF393550271bFCC2A8b48d6f65b8")

//RopstenDiscoveryAddress discovery contract address
var RopstenDiscoveryAddress = common.HexToAddress("0x4CDAF98516490d42E1E6F050bcfBD143dCb58CcD")

//NettingChannelSettleTimeoutMin min settle timeout
const NettingChannelSettleTimeoutMin = 6

/*
NettingChannelSettleTimeoutMax The maximum settle timeout is chosen as something above
 1 year with the assumption of very fast block times of 12 seconds.
 There is a maximum to avoidpotential overflows as described here:
 https://github.com/SmartRaiden/raiden/issues/1038
*/
const NettingChannelSettleTimeoutMax = 2700000

//UDPMaxMessageSize message size
const UDPMaxMessageSize = 1200

//DefaultSignalServer signal server for ice
const DefaultSignalServer = "193.112.248.133:5222"

//DefaultTurnServer  turn server
const DefaultTurnServer = "193.112.248.133:3478"

//DefaultTurnUserName turn user
const DefaultTurnUserName = "smartraiden"

//DefaultTurnPassword turn password
const DefaultTurnPassword = "smartraiden"

func init() {
	gasLimitHex = fmt.Sprintf("0x%x", GasLimit)
}
