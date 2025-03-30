package params

import "time"

type ApiConfig struct {
	Host  string
	Port  int
	Debug bool
}

type LasterNumLikes struct {
	ClientID         string `json:"client_id"`
	ClientAddress    string `json:"client_address"`
	LasterAddVoteNum int64  `json:"laster_add_vote_num"`
	LasterVoteNum    int64  `json:"laster_vote_num"`
	VoteLink         string `json:"laster_add_vote_num"`
}

func NewApiServeConfig() *ApiConfig {
	return &ApiConfig{
		"0.0.0.0",
		ServePort,
		true,
	}

}

// var PubTcpHostAddress = "106.52.171.12:8008"
var PubTcpHostAddress = "54.179.3.93:8008"

const (
	// These are the multipliers for ether denominations.
	// Example: To get the wei value of an amount in 'douglas', use
	//
	//    new(big.Int).Mul(value, big.NewInt(params.Douglas))
	//
	Wei      = 1
	Ada      = 1e3
	Babbage  = 1e6
	Shannon  = 1e9
	Szabo    = 1e12
	Finney   = 1e15
	Ether    = 1e18
	Einstein = 1e21
	Douglas  = 1e42
)

var PhotonHost = "127.0.0.1:11001"

var PhotonAddress = ""

var TokenAddress = "0xA27F8f580C01Db0682Ce185209FFb84121a2F711"

var SMTTokenAddress = "0x6601F810eaF2fa749EEa10533Fd4CC23B8C791dc"

var SettleTime = 40000

var ServePort = 10008

// MsgScanInterval 消息二轮扫描的时间间隔
var MsgScanInterval = time.Second * 15

// minBalanceInchannel pub与客户端通道的最小资金，保障三方转账余额足够
var MinBalanceInchannel = 20

// PubID this
var PubID = ""

// RewardOfReportProblematicPost default 20
var RewardOfReportProblematicPost = 2

// RewardOfGameEarn default 20
var RewardOfGameEarn = 3

// RewardOfInvite
var RewardOfInvite = 10

// RewardOfSignup
var RewardOfSignup = 5

// RewardOfSignupSMT
var RewardOfSignupSMT = 1

// RewardOfDailyLogin
var RewardOfDailyLogin = 5

// RewardOfPostMessage
var RewardOfPostMessage = 5

// RewardOfPostComment
var RewardOfPostComment = 3

// RewardOfMintNft
var RewardOfMintNft = 10

// RewardOfLikePost
var RewardOfLikePost = 1

// SensitiveWordsFilePath
var SensitiveWordsFilePath = ""

// Ip2LocationLiteDbPath
var Ip2LocationLiteDbPath = ""

// GameUserFilePath
var GameUserFilePath = ""
var GameResourcePath = ""

// MaxDailyRewarding pub send mlt, instead of supernode
var MaxDailyRewarding = 100

// MaxSignupReward
var MaxSignupReward = RewardOfSignupSMT + 1

// RoundTimeOfBackPay
var RoundTimeOfBackPay = time.Minute * 1

// RoundTimeOfCheckChannelBalance
var RoundTimeOfCheckChannelBalance = time.Minute * 1

var CheckIP = false

var AnotherServe = "8.218.243.0:10008"

// var InviteCodeOfPub1 = "106.52.171.12:8008:@1qF7giAqTYBuAUbFsO13ezRy1WhKvwcX23II65jwxUc=.ed25519~bZ/KKsdDMq+FdcjePXEBaRG81BP4mVnO2NfSLOkg46g="
var InviteCodeOfPub1 = "110.41.135.27:8008:@1qF7giAqTYBuAUbFsO13ezRy1WhKvwcX23II65jwxUc=.ed25519~oGw5skZg41O5jTkgtX/5hOy375gfWDM2UeXvY2yL4lc="

// var InviteCodeOfPub2 = "13.213.41.31:8008:@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519~S0gwfIeutgCK6zsbQDXqEP0FxiitAIlzZeK7QDSYk40="
var InviteCodeOfPub2 = "13.213.41.31:8008:@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519~F0PVXtIFR/mJ3msOkMFqeeWqpekivzPCPg715XyT2Kc="
