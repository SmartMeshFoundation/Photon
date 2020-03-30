package params

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/params"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/ethereum/go-ethereum/common"
)

/*
描述公链节点相关参数
*/
type chainConfig struct {
	Name string // 描述信息,无实际作用

	/*
		photon所连公链节点,使用者务必保证自己所链的节点是有效节点,
		如果是一个恶意节点,photon是无法检测以保证系统安全的,比如恶意通知photon没有在链上发生的事件.
		但是如果只是公链节点同步出了问题,那么photon能够检测出来,并阻止相关交易.
	*/
	GenesisBlockHash        common.Hash // 创世区块hash
	EthRPCEndPoint          string
	GasPrice                *big.Int
	BlockPeriod             time.Duration // 出块间隔
	EffectiveTimeoutSeconds int64         // 判断当前公链是否有效的依据,当前时间-最新块出块时间>该值,则认为公链无效

	/*
		事件模块轮询间隔,一般取出块间隔的一半
	*/
	PollPeriod time.Duration

	/*
		公链ChainID
	*/
	ChainID *big.Int

	/*
		photon主合约相关信息
	*/
	RegistryAddress       common.Address // 地址
	ContractVersionPrefix string         // 版本号前缀,做兼容性相关开发的时候使用
	EnableForkConfirm     bool           // 事件延迟确认开关
	ForkConfirmNumber     int64          // 延迟确认块数量,BlockNumber < 最新块-ForkConfirmNumber的事件被认为无分叉的风险
	SMTTokenName          string         // 合约支持的smt主币对应Token的名字

	/*
		photon向公链节点发起请求的超时时间,过短可能会造成因为网络延迟而把本应该成功的链接拒绝
		过长可能会造成错误迟迟无法被检测到.
	*/
	ChainRequestTimeout time.Duration

	/*
		合约调用超时时间
	*/
	TxTimeout time.Duration

	/*
		punish_block_number of contract,default is 257
	*/
	PunishBlockNumber int64
}

/*
描述节点通讯相关参数
*/
type rpcConfig struct {
	/*
		节点间UDP通信监听的Host:Port
	*/
	Host string
	Port int

	NetworkMode NetworkMode
	/*
		xmpp 相关参数
	*/
	XMPPServer string

	/*
		UDP通信相关参数
	*/
	UDPMaxMessageSize int64         // udp通信数据包最大长度
	EnableMDNS        bool          // 是否启用MDNS
	MDNSKeepalive     time.Duration // 默认mdns下20秒内检测不到在线,将该节点标志为下线
	MDNSQueryInterval time.Duration // 默认轮询间隔是1s,在测试代码中会更改他,以提高效率

	/*
		matrix相关参数
	*/
	DiscoveryServer string // discovery server
	AliasFragment   string // discovery AliasFragment
	NetworkName     string // Specify the network name of the Ethereum network to run Photon on
}

/*
photon运行节点相关参数
*/
type nodeConfig struct {
	/*
		该Photon节点的私钥,因为节点之间来往消息需要签名,因此必须保存该私钥在内存中
	*/
	PrivateKey *ecdsa.PrivateKey `json:"-"`
	MyAddress  common.Address
	/*
		最小余额,18 * GasPrice * GasLimit * 3, 当账户余额小于该值时,合约调用存在失败且tx消失的情况
	*/
	MinBalance *big.Int
}

/*
photon api 相关参数
*/
type apiConfig struct {
	RestAPIHost    string
	RestAPIPort    int
	HTTPUsername   string
	HTTPPassword   string
	RestAPITimeout time.Duration
}

/*
控制信息相关参数
*/
type controlConfig struct {
	DataBasePath      string
	DataDir           string
	Debug             bool
	DebugCrash        bool          //for test only,work with conditionQuit
	ConditionQuit     ConditionQuit //for test only
	LogFilePath       string
	EnableHealthCheck bool //send ping periodically?
}

type pmsConfig struct {
	PmsAddress common.Address
	PmsHost    string // pms server host
}

type pfsConfig struct {
	PfsHost string // pathfinder server host
}

type channelConfig struct {
	RevealTimeout             int
	SettleTimeout             int
	EnableMediationFee        bool // default false. which means no fee at all.
	IgnoreMediatedNodeRequest bool // true: this node will ignore any mediated transfer who's target is not me.
	MaxTransferDataLen        int  // 交易附带信息最大长度
	ChannelSettleTimeoutMin   int  // 支持的最小通道SettleTimeout值
}

/*
消息上报相关参数
*/
type reportConfig struct {
	ReceivedTransferReportURL string // 收到交易上报的url
}

/*
记录photon运行过程中用到的各种配置性全局变量的结构
*/
type config struct {
	/*
		定义photon的工作环境
	*/
	Env ENV
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
	IsMobile bool

	/*
		各模块相关配置及全局变量
	*/
	controlConfig
	nodeConfig
	chainConfig
	rpcConfig
	channelConfig
	pfsConfig
	pmsConfig
	apiConfig
	reportConfig
}

/*
Cfg 全局唯一,该变量应该在photon启动时初始化
初始化考虑到手机调用photon并退出然后再重启photon时,主进程不退出导致该变量无法释放的问题
*/
var Cfg *config

/*
InitDefaultCfg 刷新整个配置至各环境下的默认值
*/
func InitDefaultCfg(env ENV) {
	if Cfg != nil {
		log.Warn("params.cfg not nil when start,reset")
	}
	if env == Dev {
		Cfg = copy(&DefaultDevCfg)
	} else if env == TestNet {
		Cfg = copy(&DefaultTestNetCfg)
	} else {
		Cfg = copy(&DefaultMainNetCfg)
	}
	return
}

/*
InitForUnitTest 仅供单元测试使用
*/
func InitForUnitTest() {
	if Cfg != nil {
		return
	}
	Cfg = copy(&DefaultDevCfg)
	Cfg.ChainID = big.NewInt(7888)
}

/*
DefaultDevCfg 开发环境默认值
*/
var DefaultDevCfg = config{
	Env:      Dev,
	IsMobile: false,
	controlConfig: controlConfig{
		DataBasePath:      "", // 这个跟账户有关,不给默认值
		DataDir:           utils.DefaultDataDir(),
		Debug:             false,
		DebugCrash:        false,
		ConditionQuit:     ConditionQuit{},
		LogFilePath:       "",
		EnableHealthCheck: false,
	},
	nodeConfig: nodeConfig{
		PrivateKey: nil,
		MyAddress:  utils.EmptyAddress,
		MinBalance: big.NewInt(18 * 1e9 * 100000 * 3),
	},
	chainConfig: chainConfig{
		Name:                    "DevChain",
		GenesisBlockHash:        utils.EmptyHash,
		EthRPCEndPoint:          "",
		GasPrice:                big.NewInt(params.Shannon * 20),
		BlockPeriod:             time.Second, // 默认开发环境一秒一块
		EffectiveTimeoutSeconds: 180,
		PollPeriod:              time.Second / 2,
		ChainID:                 nil,
		RegistryAddress:         utils.EmptyAddress,
		ContractVersionPrefix:   "0.6",
		EnableForkConfirm:       false,
		ForkConfirmNumber:       17,
		SMTTokenName:            "SMTToken",
		ChainRequestTimeout:     3 * time.Second,
		TxTimeout:               5 * time.Minute,
		PunishBlockNumber:       257,
	},
	rpcConfig: rpcConfig{
		Host:        "0.0.0.0",
		Port:        40001,
		NetworkMode: MixUDPMatrix,
		XMPPServer:  "193.112.248.133:5222",

		UDPMaxMessageSize: 1200,
		EnableMDNS:        true,
		MDNSKeepalive:     20 * time.Second,
		MDNSQueryInterval: time.Second,

		DiscoveryServer: "transport01.smartmesh.cn",
		AliasFragment:   "discovery",
		NetworkName:     "ropsten",
	},
	channelConfig: channelConfig{
		RevealTimeout:             30,
		SettleTimeout:             600,
		EnableMediationFee:        true,
		IgnoreMediatedNodeRequest: false,
		MaxTransferDataLen:        256,
		ChannelSettleTimeoutMin:   60, // 开发环境默认60
	},
	pfsConfig: pfsConfig{
		PfsHost: "http://transport01.smartmesh.cn:7002",
		//PfsHostForXMPP:   "http://transport01.smartmesh.cn:7002", // matrix开发环境
		//PfsHostForMatrix: "http://transport01.smartmesh.cn:7012", // xmpp开发环境
	},
	pmsConfig: pmsConfig{
		PmsAddress: common.HexToAddress("0xa668da12fe5f5729cbce9ae697d56bac929766f4"),
		PmsHost:    "http://transport01.smartmesh.cn:7005",
	},
	apiConfig: apiConfig{
		RestAPIHost:    "127.0.0.1",
		RestAPIPort:    5001,
		HTTPUsername:   "",
		HTTPPassword:   "",
		RestAPITimeout: 20 * time.Minute,
	},
}

/*
DefaultTestNetCfg TestNet环境默认值
*/
var DefaultTestNetCfg = config{
	Env:      TestNet,
	IsMobile: false,
	controlConfig: controlConfig{
		DataBasePath:      "", // 这个跟账户有关,不给默认值
		DataDir:           utils.DefaultDataDir(),
		Debug:             false,
		DebugCrash:        false,
		ConditionQuit:     ConditionQuit{},
		LogFilePath:       "",
		EnableHealthCheck: false,
	},
	nodeConfig: nodeConfig{
		PrivateKey: nil,
		MyAddress:  utils.EmptyAddress,
		MinBalance: big.NewInt(18 * 1e9 * 100000 * 3),
	},
	chainConfig: chainConfig{
		Name:                    "SpectrumTestNet",
		GenesisBlockHash:        common.HexToHash("0xd011e2cc7f241996a074e2c48307df3971f5f1fe9e1f00cfa704791465d5efc3"),
		EthRPCEndPoint:          "",
		GasPrice:                big.NewInt(params.Shannon * 20),
		BlockPeriod:             14 * time.Second, // TestNet 14秒一块
		EffectiveTimeoutSeconds: 180,
		PollPeriod:              7 * time.Second,
		ChainID:                 big.NewInt(3), // Spectrum test net chain ID
		RegistryAddress:         common.HexToAddress("0x50839B01D28390048616C8f28dD1A21CF3CacbfF"),
		ContractVersionPrefix:   "0.6",
		EnableForkConfirm:       false,
		ForkConfirmNumber:       17,
		SMTTokenName:            "SMTToken",
		ChainRequestTimeout:     3 * time.Second,
		TxTimeout:               5 * time.Minute,
		PunishBlockNumber:       257,
	},
	rpcConfig: rpcConfig{
		Host:        "0.0.0.0",
		Port:        40001,
		NetworkMode: MixUDPMatrix,
		XMPPServer:  "193.112.248.133:5222",

		UDPMaxMessageSize: 1200,
		EnableMDNS:        true,
		MDNSKeepalive:     20 * time.Second,
		MDNSQueryInterval: time.Second,

		DiscoveryServer: "transport01.smartmesh.cn",
		AliasFragment:   "discovery",
		NetworkName:     "ropsten",
	},
	channelConfig: channelConfig{
		RevealTimeout:             30,
		SettleTimeout:             600,
		EnableMediationFee:        false,
		IgnoreMediatedNodeRequest: false,
		MaxTransferDataLen:        256,
		ChannelSettleTimeoutMin:   60, // 开发环境默认60
	},
	pfsConfig: pfsConfig{
		PfsHost: "http://transport01.smartmesh.cn:7001",
		//PfsHostForXMPP:   "http://transport01.smartmesh.cn:7001", // matrix testnet环境
		//PfsHostForMatrix: "http://transport01.smartmesh.cn:7011", // xmpp testnet环境
	},
	pmsConfig: pmsConfig{
		PmsAddress: common.HexToAddress("0xaed9188842c05e07bf5abdde2fb400432ae49d28"),
		PmsHost:    "http://transport01.smartmesh.cn:7004",
	},
	apiConfig: apiConfig{
		RestAPIHost:    "127.0.0.1",
		RestAPIPort:    5001,
		HTTPUsername:   "",
		HTTPPassword:   "",
		RestAPITimeout: 20 * time.Minute,
	},
}

/*
DefaultMainNetCfg MainNet环境默认值
*/
var DefaultMainNetCfg = config{
	Env:      MainNet,
	IsMobile: false,
	controlConfig: controlConfig{
		DataBasePath:      "", // 这个跟账户有关,不给默认值
		DataDir:           utils.DefaultDataDir(),
		Debug:             false,
		DebugCrash:        false,
		ConditionQuit:     ConditionQuit{},
		LogFilePath:       "",
		EnableHealthCheck: false,
	},
	nodeConfig: nodeConfig{
		PrivateKey: nil,
		MyAddress:  utils.EmptyAddress,
		MinBalance: big.NewInt(18 * 1e9 * 100000 * 3),
	},
	chainConfig: chainConfig{
		Name:                    "DevChain",
		GenesisBlockHash:        common.HexToHash("0x57e682b80257aad73c4f3ad98d20435b4e1644d8762ef1ea1ff2806c27a5fa3d"),
		EthRPCEndPoint:          "",
		GasPrice:                big.NewInt(params.Shannon * 20),
		BlockPeriod:             14 * time.Second, // 主链14秒一块
		EffectiveTimeoutSeconds: 180,
		PollPeriod:              7 * time.Second,
		ChainID:                 big.NewInt(1), // Spectrum 主链
		RegistryAddress:         common.HexToAddress("0x242e0de2B118279D1479545A131a90A8f67A2512"),
		ContractVersionPrefix:   "0.6",
		EnableForkConfirm:       false,
		ForkConfirmNumber:       17,
		SMTTokenName:            "SMTToken",
		ChainRequestTimeout:     3 * time.Second,
		TxTimeout:               5 * time.Minute,
		PunishBlockNumber:       257,
	},
	rpcConfig: rpcConfig{
		Host:        "0.0.0.0",
		Port:        40001,
		NetworkMode: MixUDPMatrix,
		XMPPServer:  "193.112.248.133:5222",

		UDPMaxMessageSize: 1200,
		EnableMDNS:        true,
		MDNSKeepalive:     20 * time.Second,
		MDNSQueryInterval: time.Second,

		DiscoveryServer: "transport01.smartmesh.cn",
		AliasFragment:   "discovery",
		NetworkName:     "ropsten",
	},
	channelConfig: channelConfig{
		RevealTimeout:             30,
		SettleTimeout:             600,
		EnableMediationFee:        true,
		IgnoreMediatedNodeRequest: false,
		MaxTransferDataLen:        256,
		ChannelSettleTimeoutMin:   40000,
	},
	pfsConfig: pfsConfig{
		PfsHost: "http://transport01.smartmesh.cn:7000",
		//PfsHostForXMPP:   "http://transport01.smartmesh.cn:7000", // matrix mainnet环境
		//PfsHostForMatrix: "http://transport01.smartmesh.cn:7010", // xmpp mainnet环境
	},
	pmsConfig: pmsConfig{
		PmsAddress: common.HexToAddress("0xa94399b93da31e25ab5612de8c64556694d5f2fd"),
		PmsHost:    "http://transport01.smartmesh.cn:7003",
	},
	apiConfig: apiConfig{
		RestAPIHost:    "127.0.0.1",
		RestAPIPort:    5001,
		HTTPUsername:   "",
		HTTPPassword:   "",
		RestAPITimeout: 20 * time.Minute,
	},
}

/*
Usable 该方法校验config是否已经初始化完毕,进入可用状态
*/
func (c *config) Usable() bool {
	/*
		这三个必须根据启动photon所使用的账户来进行初始化,否则无法启动
	*/
	if c.PrivateKey == nil {
		return false
	}
	if c.MyAddress == utils.EmptyAddress {
		return false
	}
	if c.DataBasePath == "" {
		return false
	}

	/*
		根据连接的公链节点来初始化
	*/
	if c.EthRPCEndPoint == "" {
		return false
	}
	if c.ChainID == nil {
		return false
	}
	return true
}

func copy(src *config) (dst *config) {
	if src == nil {
		return
	}
	dst = new(config)
	buf, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, dst)
	if err != nil {
		panic(err)
	}
	return
}
