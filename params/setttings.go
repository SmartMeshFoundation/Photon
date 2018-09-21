package params

import (
	"fmt"
	"math/big"

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
const DefaultRevealTimeout = 5

//DefaultSettleTimeout settle time of channel
const DefaultSettleTimeout = 600

//DefaultPollTimeout  request wait time
const DefaultPollTimeout = 180 * time.Second

//DefaultJoinableFundsTarget for connection api
const DefaultJoinableFundsTarget = 0.4

//DefaultInitialChannelTarget channels to create
const DefaultInitialChannelTarget = 3

//DefaultTxTimeout args
const DefaultTxTimeout = 5 * time.Minute //15seconds for one block,it may take sever minutes
//MaxRequestTimeout args
const MaxRequestTimeout = 20 * time.Minute //longest time for a request ,for example ,settle all channles?

var gasLimitHex string

//SpectrumTestNetRegistryAddress Registry contract address
var SpectrumTestNetRegistryAddress = common.HexToAddress("0x52d7167FAD53835a2356C7A872BfbC17C03aD758")

//ChannelSettleTimeoutMin min settle timeout
const ChannelSettleTimeoutMin = 6

/*
ChannelSettleTimeoutMax The maximum settle timeout is chosen as something above
 1 year with the assumption of very fast block times of 12 seconds.
 There is a maximum to avoidpotential overflows as described here:
 https://github.com/SmartRaiden/raiden/issues/1038
*/
const ChannelSettleTimeoutMax = 2700000

//UDPMaxMessageSize message size
const UDPMaxMessageSize = 1200

//DefaultXMPPServer xmpp server
const DefaultXMPPServer = "193.112.248.133:5222"

//DefaultTestXMPPServer xmpp server for test only
const DefaultTestXMPPServer = "193.112.248.133:5222" //"182.254.155.208:5222"
//ContractSignaturePrefix for EIP191 https://github.com/ethereum/EIPs/blob/master/EIPS/eip-191.md
var ContractSignaturePrefix = []byte("\x19Ethereum Signed Message:\n")

const (
	//ContractBalanceProofMessageLength balance proof  length
	ContractBalanceProofMessageLength = "176"
	//ContractBalanceProofDelegateMessageLength update balance proof delegate length
	ContractBalanceProofDelegateMessageLength = "144"
	//ContractCooperativeSettleMessageLength cooperative settle channel proof length
	ContractCooperativeSettleMessageLength = "176"
	//ContractDisposedProofMessageLength annouce disposed proof length
	ContractDisposedProofMessageLength = "136"
	//ContractWithdrawProofMessageLength withdraw proof length
	ContractWithdrawProofMessageLength = "156"
	//ContractUnlockDelegateProofMessageLength unlock delegate proof length
	ContractUnlockDelegateProofMessageLength = "188"
)

func init() {
	gasLimitHex = fmt.Sprintf("0x%x", GasLimit)
}

/*
MobileMode works on mobile device, 移动设备模式,这时候 smartraiden 并不是一个独立的进程,这时候很多工作模式要发生变化.
比如:
1.不能任意退出
2. 对于网络通信的处理要更谨慎
3. 对于资源的消耗如何控制?
*/
/*
 *	MobileMode : a boolean value to adapt with mobile modes.
 *
 *	Note : if true, then smartraiden is not an individual process, work mode is about to change.
 *		1. not support exit arbitrarily.
 *		2. handle internet communication more prudent.
 *		3. How to control amount of resource consumption.
 */
var MobileMode bool

/*
InTest are we test now?
*/
var InTest = true

//ChainID of this tokenNetwork
var ChainID = big.NewInt(1)

//MatrixServerConfig matrix server config
var MatrixServerConfig = [][]string{
	{"http://transport01.smartmesh.cn:8008", "transport01.smartmesh.cn"},
	{"http://transport02.smartmesh.cn:8008", "transport02.smartmesh.cn"},
	{"http://transport03.smartmesh.cn:8008", "transport03.smartmesh.cn"},
}

//MatrixDiscoveryRoomConfig a defalut matrix room config
var MatrixDiscoveryRoomConfig = [][]string{
	{"aliassegment", "discovery"},
	{"server", "transport01.smartmesh.cn"},
}

//NETWORKNAME Specify the network name of the Ethereum network to run SmartRaiden on
var NETWORKNAME = "ropsten"
