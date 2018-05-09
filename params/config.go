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
type NetworkMode int

const (
	/*
		节点不对外暴露网络接口,仅供测试使用
	*/
	NoNetwork = iota + 1
	/*
		通过udp ip 端口对外暴露服务,可以使用 stun,upnp 等方式,依赖节点发现合约或者直接告知其他节点 ip 端口
	*/
	UDPOnly
	/*
		通过信令服务器协助,建立连接.
	*/
	ICEOnly
	/*
		适应无网通信需要,将上面两种方式混合使用,有网时使用 ice 建立连接,无网时则使用 udp 直接暴露 ip 端口
	*/
	MixUDPICE
)

type Config struct {
	Host                      string
	Port                      int
	ExternIp                  string
	ExternPort                int
	PrivateKeyHex             string
	PrivateKey                *ecdsa.PrivateKey
	RevealTimeout             int
	SettleTimeout             int
	DataBasePath              string
	MsgTimeout                time.Duration
	Protocol                  protocolConfig
	UseRpc                    bool
	UseConsole                bool
	ApiHost                   string
	ApiPort                   int
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

var DefaultConfig = Config{
	Port:          INITIAL_PORT,
	ExternPort:    INITIAL_PORT,
	PrivateKeyHex: "",
	RevealTimeout: DEFAULT_REVEAL_TIMEOUT,
	SettleTimeout: DEFAULT_SETTLE_TIMEOUT,
	Protocol: protocolConfig{
		RetryInterval:        DEFAULT_PROTOCOL_RETRY_INTERVAL,
		RetriesBeforeBackoff: DEFAULT_PROTOCOL_RETRIES_BEFORE_BACKOFF,
		ThrottleCapacity:     DEFAULT_PROTOCOL_THROTTLE_CAPACITY,
		ThrottleFillRate:     DEFAULT_PROTOCOL_THROTTLE_FILL_RATE,
		NatInvitationTimeout: DEFAULT_NAT_INVITATION_TIMEOUT,
		NatKeepAliveRetries:  DEFAULT_NAT_KEEPALIVE_RETRIES,
		NatKeepAliveTimeout:  DEFAULT_NAT_KEEPALIVE_TIMEOUT,
	},
	UseRpc:           true,
	UseConsole:       false,
	RegistryAddress:  ROPSTEN_REGISTRY_ADDRESS,
	DiscoveryAddress: ROPSTEN_DISCOVERY_ADDRESS,
	MsgTimeout:       100 * time.Second,
	Ice: iceConfig{
		SignalServer: DefaultSignalServer,
	},
}

type ConditionQuit struct {
	QuitEvent  string //name match
	IsBefore   bool   //quit before event occur
	RandomQuit bool   //random exit
}

/*
When refund occurs in the intermediary node,is it treated as a common mediatedtransfer(that is to delete HandleSecret in raidenservice)?
*/
var TreatRefundTransferAsNormalMediatedTransfer = true

func init() {

}
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

func DefaultKeyStoreDir() string {
	return filepath.Join(node.DefaultDataDir(), "keystore")
}
