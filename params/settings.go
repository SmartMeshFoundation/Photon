package params

import (
	"fmt"
	"time"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

//InitialPort listening port for communication bewtween nodes
const InitialPort = 40001

//DefaultGasLimit max gas usage for photon tx
const DefaultGasLimit = 3141592 //den's gasLimit.
//DefaultGasPrice from ethereum
const DefaultGasPrice = params.Shannon * 20

//defaultProtocolRetiesBeforeBackoff
const defaultProtocolRetiesBeforeBackoff = 5
const defaultProtocolRhrottleCapacity = 10.
const defaultProtocolThrottleFillRate = 10.
const defaultprotocolRetryInterval = 1.

//DefaultRevealTimeout blocks needs to update transfer
//this time is used for a participant to register secret on chain
// and unlock the lock if need.
var DefaultRevealTimeout = 30

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

//ChannelSettleTimeoutMin min settle timeout
const ChannelSettleTimeoutMin = 6

/*
ChannelSettleTimeoutMax The maximum settle timeout is chosen as something above
 1 year with the assumption of very fast block times of 12 seconds.
 There is a maximum to avoidpotential overflows as described here:
 https://github.com/Photon/photon/issues/1038
*/
const ChannelSettleTimeoutMax = 2700000

//UDPMaxMessageSize message size
const UDPMaxMessageSize = 1200

//DefaultXMPPServer xmpp server
const DefaultXMPPServer = "193.112.248.133:5222"

//TestLogServer only for test, enabled if --debug flag is set
var TestLogServer = "http://transport01.smartmesh.cn:8008"

//var TestLogServer = "http://127.0.0.1:5000"

//DefaultTestXMPPServer xmpp server for test only
const DefaultTestXMPPServer = "193.112.248.133:5222" //"182.254.155.208:5222"
//ContractSignaturePrefix for EIP191 https://github.com/ethereum/EIPs/blob/master/EIPS/eip-191.md
var ContractSignaturePrefix = []byte("\x19Spectrum Signed Message:\n")

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
	gasLimitHex = fmt.Sprintf("0x%x", DefaultGasLimit)
}

/*
MobileMode works on mobile device, 移动设备模式,这时候 photon 并不是一个独立的进程,这时候很多工作模式要发生变化.
比如:
1.不能任意退出
2. 对于网络通信的处理要更谨慎
3. 对于资源的消耗如何控制?
*/
/*
 *	MobileMode : a boolean value to adapt with mobile modes.
 *
 *	Note : if true, then photon is not an individual process, work mode is about to change.
 *		1. not support exit arbitrarily.
 *		2. handle internet communication more prudent.
 *		3. How to control amount of resource consumption.
 */
var MobileMode bool

/*
InTest are we test now?
*/
var InTest = true

// DefaultChainID :
var DefaultChainID = big.NewInt(0)

//ChainID of this tokenNetwork
var ChainID = DefaultChainID

//MatrixServerConfig matrix server config
var MatrixServerConfig = map[string]string{
	"transport01.smartmesh.cn": "http://transport01.smartmesh.cn:8008",
	//"transport02.smartmesh.cn": "http://transport02.smartmesh.cn:8008",
	//"transport03.smartmesh.cn": "http://transport03.smartmesh.cn:8008",
}

//AliasFragment  is discovery AliasFragment
const AliasFragment = "discovery"

//DiscoveryServer is discovery server
const DiscoveryServer = "transport01.smartmesh.cn"

//NETWORKNAME Specify the network name of the Ethereum network to run Photon on
var NETWORKNAME = "ropsten"

//GenesisBlockHashToDefaultRegistryAddress :
var GenesisBlockHashToDefaultRegistryAddress = map[common.Hash]common.Address{
	// spectrum
	common.HexToHash("0x57e682b80257aad73c4f3ad98d20435b4e1644d8762ef1ea1ff2806c27a5fa3d"): common.HexToAddress("0x08b7d79ec4ebd53e5b89c7c062cc64bb09d063e3"),
	// spectrum test net
	common.HexToHash("0xd011e2cc7f241996a074e2c48307df3971f5f1fe9e1f00cfa704791465d5efc3"): common.HexToAddress("0xc479184abeb8c508ee96e4c093ee47af2256cbbf"),
	// ethereum
	common.HexToHash("0x88e96d4537bea4d9c05d12549907b32561d3bf31f45aae734cdc119f13406cb6"): utils.EmptyAddress,
	// ethereum test net
	common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"): utils.EmptyAddress,
	// ethereum private
	common.HexToHash("0x38a88a9ddffe522df5c07585a7953f8c011c94327a494188bd0cc2410dc40a1a"): common.HexToAddress("0x2907b8bf0fF92dA818E2905fB5218b1A8323Ffb4"),
}

//GenesisBlockHashToPFS : default pfs provider
var GenesisBlockHashToPFS = map[common.Hash]string{
	// spectrum
	common.HexToHash("0x57e682b80257aad73c4f3ad98d20435b4e1644d8762ef1ea1ff2806c27a5fa3d"): "http://transport01.smartmesh.cn:7000",
	// spectrum test net
	common.HexToHash("0xd011e2cc7f241996a074e2c48307df3971f5f1fe9e1f00cfa704791465d5efc3"): "http://transport01.smartmesh.cn:7001",
	// ethereum
	common.HexToHash("0x88e96d4537bea4d9c05d12549907b32561d3bf31f45aae734cdc119f13406cb6"): "",
	// ethereum test net
	common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"): "",
	// ethereum private
	common.HexToHash("0x38a88a9ddffe522df5c07585a7953f8c011c94327a494188bd0cc2410dc40a1a"): "http://transport01.smartmesh.cn:7002",
}

// DefaultEthRPCPollPeriodForTest :
var DefaultEthRPCPollPeriodForTest = 500 * time.Millisecond

// DefaultEthRPCPollPeriod :
var DefaultEthRPCPollPeriod = 7500 * time.Millisecond

// TestPrivateChainID :
var TestPrivateChainID int64 = 8888

// TestPrivateChainID2 : for travis fast test
var TestPrivateChainID2 int64 = 7888

// EthRPCTimeout :
var EthRPCTimeout = 3 * time.Second

// ContractVersionPrefix :
var ContractVersionPrefix = "0.6"

// EnableForkConfirm : 事件延迟确认开关
var EnableForkConfirm = false

// ForkConfirmNumber : 分叉确认块数量,BlockNumber < 最新块-ForkConfirmNumber的事件被认为无分叉的风险
var ForkConfirmNumber int64 = 17

// MaxTransferDataLen : 交易附件信息最大长度
var MaxTransferDataLen = 256
