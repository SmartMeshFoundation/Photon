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
	NatInvitationTimeout int
	NatKeepAliveRetries  int
	NatKeepAliveTimeout  int64
}

//NetworkMode is transport status
type NetworkMode int

const (
	//NoNetwork 节点不对外暴露网络接口,仅供测试使用
	NoNetwork NetworkMode = iota + 1
	//UDPOnly 通过udp ip 端口对外暴露服务,可以使用 stun,upnp 等方式,依赖节点发现合约或者直接告知其他节点 ip 端口
	UDPOnly
	//ICEOnly 通过信令服务器协助,建立连接.
	ICEOnly
	//MixUDPICE 适应无网通信需要,将上面两种方式混合使用,有网时使用 ice 建立连接,无网时则使用 udp 直接暴露 ip 端口
	MixUDPICE
)

//Config is configuration for Raiden,
type Config struct {
	Host                      string
	Port                      int
	ExternIP                  string
	ExternPort                int
	PrivateKeyHex             string
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
	DiscoveryAddress          common.Address
	DataDir                   string
	MyAddress                 common.Address
	DebugCrash                bool
	ConditionQuit             ConditionQuit
	Ice                       iceConfig
	NetworkMode               NetworkMode
	EnableMediationFee        bool //default false. which means no fee at all.
	IgnoreMediatedNodeRequest bool // true: this node will ignore any mediated transfer who's target is not me.
	EnableHealthCheck         bool //send ping periodically?
}
type iceConfig struct {
	/*
		signal server url for ice
	*/
	SignalServer string
	/*
		must be xmpp
	*/
	SignalEngine string
	/*
		turn server for ice
	*/
	TurnServer   string
	StunServer   string
	TurnUser     string
	TurnPassword string
}

//DefaultConfig default config
var DefaultConfig = Config{
	Port:          InitialPort,
	ExternPort:    InitialPort,
	PrivateKeyHex: "",
	RevealTimeout: DefaultRevealTimeout,
	SettleTimeout: DefaultSettleTimeout,
	Protocol: protocolConfig{
		RetryInterval:        defaultprotocolRetryInterval,
		RetriesBeforeBackoff: defaultProtocolRetiesBeforeBackoff,
		ThrottleCapacity:     defaultProtocolRhrottleCapacity,
		ThrottleFillRate:     defaultProtocolThrottleFillRate,
		NatInvitationTimeout: DefaultNATInvitationTimeout,
		NatKeepAliveRetries:  DefaultKeepAliveReties,
		NatKeepAliveTimeout:  DefaultNATKeepAliveTimeout,
	},
	UseRPC:           true,
	UseConsole:       false,
	RegistryAddress:  RopstenRegistryAddress,
	DiscoveryAddress: RopstenDiscoveryAddress,
	MsgTimeout:       100 * time.Second,
	Ice: iceConfig{
		SignalServer: DefaultSignalServer,
	},
	EnableHealthCheck: false,
}

//conditionQuit is for test
type ConditionQuit struct {
	QuitEvent  string //name match
	IsBefore   bool   //quit before event occur
	RandomQuit bool   //random exit
}

/*
TreatRefundTransferAsNormalMediatedTransfer When refund occurs in the intermediary node,is it treated as a common mediatedtransfer(that is to delete handleSecret in raidenservice)?
todo remove?
*/
var TreatRefundTransferAsNormalMediatedTransfer = true

func init() {

}

//DefaultDataDir default work directory
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "smartraiden")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "smartraiden")
		} else {
			return filepath.Join(home, ".smartraiden")
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
