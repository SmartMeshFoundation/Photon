package params

import (
	"crypto/ecdsa"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/node"
)

type protocolConfig struct {
	RetryInterval        int
	RetriesBeforeBackoff int
	ThrottleCapacity     float64
	ThrottleFillRate     float64
}

//NetworkMode is transport status
type NetworkMode int

const (
	//NoNetwork 不与其他节点之间进行通信,仅供测试使用
	// NoNetwork : Node does not  communicates with other nodes, just for test.
	NoNetwork NetworkMode = iota + 1
	//UDPOnly 通过udp ip 端口对外暴露服务,通过使用MDNS进行节点服务发现
	// UDPOnly : expose service via udp ip,
	UDPOnly
	//XMPPOnly 通过XMPP服务器进行通信,目前暂时不使用
	// XMPPOnly : communicate via XMPP server.
	XMPPOnly
	//MixUDPXMPP 适应无网通信需要,将上面两种方式混合使用,有网时使用XMPP建立连接,无网时则使用 udp 直接暴露 ip 端口
	// MixUDPXMPP : used for Internet-free network, combining UDPOnly and XMPPOnly.
	// While Internet, it use XMPP to create connection; while Internet-free, it use udp to expose ip port.
	MixUDPXMPP
	//MixUDPMatrix Matrix and UDP at the same time
	MixUDPMatrix
)

//Config is configuration for Photon,
type Config struct {
	/*
		photon所连公链节点,使用者务必保证自己所链的节点是有效节点,
		如果是一个恶意节点,photon是无法检测以保证系统安全的,比如恶意通知photon没有在链上发生的事件.
		但是如果只是公链节点同步出了问题,那么photon能够检测出来,并阻止相关交易.
	*/
	EthRPCEndPoint string
	/*
		Host 节点间UDP通信监听的Host
	*/
	Host string
	/*
		Port 节点间UDP通信监听的Port
	*/
	Port int
	/*
		该Photon节点的私钥,因为节点之间来往消息需要签名,因此必须保存该私钥在内存中
	*/
	PrivateKey *ecdsa.PrivateKey
	/*
		专门留给节点进行链上unlock的时间,
	*/
	RevealTimeout             int
	SettleTimeout             int
	DataBasePath              string
	MsgTimeout                time.Duration
	Protocol                  protocolConfig
	UseRPC                    bool
	UseConsole                bool
	APIHost                   string
	APIPort                   int
	RegistryAddress           common.Address
	DataDir                   string
	MyAddress                 common.Address
	Debug                     bool
	DebugCrash                bool          //for test only,work with conditionQuit
	ConditionQuit             ConditionQuit //for test only
	NetworkMode               NetworkMode
	EnableMediationFee        bool //default false. which means no fee at all.
	IgnoreMediatedNodeRequest bool // true: this node will ignore any mediated transfer who's target is not me.
	EnableHealthCheck         bool //send ping periodically?
	XMPPServer                string
	PfsHost                   string // pathfinder server host
	HTTPUsername              string
	HTTPPassword              string
	PmsHost                   string // pms server host
	PmsAddress                common.Address
}

//DefaultConfig default config
var DefaultConfig = Config{
	Port:          InitialPort,
	RevealTimeout: DefaultRevealTimeout,
	SettleTimeout: DefaultSettleTimeout,
	Protocol: protocolConfig{
		RetryInterval:        defaultprotocolRetryInterval,
		RetriesBeforeBackoff: defaultProtocolRetiesBeforeBackoff,
		ThrottleCapacity:     defaultProtocolRhrottleCapacity,
		ThrottleFillRate:     defaultProtocolThrottleFillRate,
	},
	UseRPC:            true,
	UseConsole:        false,
	MsgTimeout:        100 * time.Second,
	EnableHealthCheck: false,
	XMPPServer:        DefaultXMPPServer,
}

//ConditionQuit is for test
type ConditionQuit struct {
	QuitEvent  string //name match
	IsBefore   bool   //quit before event occur
	RandomQuit bool   //random exit
}

//DefaultDataDir default work directory
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "photon")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "photon")
		} else {
			return filepath.Join(home, ".photon")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

//DefaultKeyStoreDir keystore path of ethereum
func DefaultKeyStoreDir() string {
	return filepath.Join(node.DefaultDataDir(), "keystore")
}
