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
	//NoNetwork 节点不对外暴露网络接口,仅供测试使用
	// NoNetwork : Node does not expose interface, just for case.
	NoNetwork NetworkMode = iota + 1
	//UDPOnly 通过udp ip 端口对外暴露服务,可以使用 stun,upnp 等方式,依赖节点发现合约或者直接告知其他节点 ip 端口
	// UDPOnly : expose service via udp ip, we can use stun, upnp and node to find contracts or tell other nodes with ip port.
	UDPOnly
	//XMPPOnly 通过XMPP服务器进行通信
	// XMPPOnly : communicate via XMPP server.
	XMPPOnly
	//MixUDPXMPP 适应无网通信需要,将上面两种方式混合使用,有网时使用 ice 建立连接,无网时则使用 udp 直接暴露 ip 端口
	// MixUDPXMPP : used for Internet-free network, combining UDPOnly and XMPPOnly.
	// While Internet, it use ice to create connection; while Internet-free, it use udp to expose ip port.
	MixUDPXMPP
	//MixUDPMatrix Matrix and UDP at the same time
	MixUDPMatrix
)

//Config is configuration for Photon,
type Config struct {
	EthRPCEndPoint            string
	Host                      string
	Port                      int
	PrivateKey                *ecdsa.PrivateKey
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
	IsMeshNetwork             bool   //is mesh now?
	PfsHost                   string // pathfinder server host
	HTTPUsername              string
	HTTPPassword              string
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
