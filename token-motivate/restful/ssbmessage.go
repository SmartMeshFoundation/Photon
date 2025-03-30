package restful

import (
	"encoding/json"
	"time"

	"math/big"

	"sync"

	refs "go.mindeco.de/ssb-refs"
)

type MessageValue struct {
	Previous  *refs.MessageRef `json:"previous"`
	Sequence  int64            `json:"sequence"`
	Author    refs.FeedRef     `json:"author"`
	Timestamp float64          `json:"timestamp"`
	Hash      string           `json:"hash"`
	Content   json.RawMessage  `json:"content"`
	Signature string           `json:"signature"`
}

type DeserializedMessageStu struct {
	Key       string        `json:"key"`
	Value     *MessageValue `json:"value"`
	Timestamp float64       `json:"timestamp"`
}

// LikeDetail 存储一轮搜索到的被Like的消息ID
var LikeDetail []string

// LikeDetail 存储一轮搜索到的被Unlike的消息ID
var UnLikeDetail []string

// LikeCount for save message for search likes's author
var TempMsgMap map[string]*TempdMessage

// ClientID2Name for save message for search likes's author clientid---name
var ClientID2Name map[string]string

// 注册IP->time
var RegisterSourceMap map[string]int64

// PhotonNodeStatus only record the node status from backpay infos
type PhotonNodeStatusStu struct {
	lock            sync.RWMutex
	nodesstatusinfo map[string]*ClientOnlineStatus
}

//var PhotonNodeStatus map[string]*ClientOnlineStatus

// ClientOnlineStatus
type ClientOnlineStatus struct {
	ClientID         string
	ClientEthAddress string
	Online           bool
}

func PrintTime() string {
	//var LOC, _ = time.LoadLocation("Asia/Shanghai")
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "] "
}

// ContentVoteStru
type ContentVoteStru struct {
	Type string    `json:"type"`
	Vote *VoteStru `json:"vote"`
}

// ContentContactStru
type ContentContactStru struct {
	Type      string `json:"type"`
	Contact   string `json:"contact"`
	Following bool   `json:"following"`
	//Blocking  bool   `json:"blocking"`
	Pub bool `json:"pub"`
}

// UserFeedBackStu
type UserFeedBackStu struct {
	UserSsbId  string `json:"user_ssbid"`
	SubmitType string `json:"submit_type"`
	Content    string `json:"content"`
	SubmitTime int64  `json:"submit_time"`
	UserEmail  string `json:"user_email"`
	HadReplay  int    `json:"had_replay"`
}

// TippedOff reasons:"xxx|xxx|xxx"
type TippedOffStu struct {
	Plaintiff  string `json:"plaintiff"`
	Defendant  string `json:"defendant"`
	MessageKey string `json:"messagekey"`
	Reasons    string `json:"reasons"`
	DealTag    string `json:"dealtag"`
	Recordtime int64  `json:"recordtime"`
	Dealtime   int64  `json:"dealtime"`
	Dealreward string `json:"dealreward"`
}

// VoteStru
type VoteStru struct {
	Link       string `json:"link"`
	value      bool   `json:"value,omitempty"` //两种格式int，bool,2023-1-8
	Expression string `json:"expression"`
}

// ContentAboutStru
type ContentAboutStru struct {
	Type  string `json:"type"`
	About string `json:"about"`
	Name  string `json:"name"`
}

// ContentPostStru
type ContentPostStru struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Root     string          `json:"root"`
	Mentions json.RawMessage `json:"mentions"`
	Branch   string          `json:"branch"`
}

// LasterNumLikes
type LasterNumLikes struct {
	ClientID         string `json:"client_id"`
	LasterLikeNum    int    `json:"laster_like_num"`
	Name             string `json:"client_name"`
	ClientEthAddress string `json:"client_eth_address"`
	MessageFromPub   string `json:"message_from_pub"`
}

// TempdMessage 用于一次搜索的结果统计
type TempdMessage struct {
	Author string `json:"author"`
}

// Name2ProfileReponse
type Name2ProfileReponse struct {
	ID                 string `json:"client_id"`
	Name               string `json:"client_Name"`
	Alias              string `json:"client_alias"`
	Bio                string `json:"client_bio"`
	EthAddress         string `json:"client_eth_address"`
	AppVersion         string `json:"client_version_build"`
	RegisteTime        int64  `json:"registe_time"`
	LastactiveTime     int64  `json:"last_active_time"`
	PersonalInviteCode string `json:"personal_invite_code"`
}

type CloseChannelReq struct {
	PartnerAddress string `json:"partner_address"`
	Signature      string `json:"signature"`
}

// Whoami
type Whoami struct {
	Pub_Id          string `json:"pub_id"`
	Pub_Eth_Address string `json:"pub_eth_address"`
}

// SensitiveWords
var SensitiveWords []string

// EventSensitive
type EventSensitive struct {
	PubID           string `json:"pub_id"`
	MessageScanTime int64  `json:"message_scan_time"`
	MessageText     string `json:"message_text"`
	MessageKey      string `json:"message_key"`
	MessageAuthor   string `json:"message_author"`
	DealTag         string `json:"deal_tag"`
	DealTime        int64  `json:"deal_time"`
}

// UserTasks
type UserTasks struct {
	CollectFromPub string `json:"pub_id"`
	Author         string `json:"author"`
	MessageKey     string `json:"message_key"`
	MessageType    string `json:"message_type"`
	MessageRoot    string `json:"message_root"`
	MessageTime    int64  `json:"message_time"`

	NfttxHash    string `json:"nft_tx_hash"`
	NftTokenId   string `json:"nft_token_id"`
	NftStoredUrl string `json:"nft_store_url"`

	ClientEthAddress string `json:"client_eth_address"`
}

// ReqUserTask
type ReqUserTask struct {
	Author      string `json:"author"`
	MessageType string `json:"message_type"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
}

// ReqUserLoginApp
type ReqUserLoginApp struct {
	ClientID  string `json:"client_id"`
	LoginTime int64  `json:"login_time"`
}

// ReqCreatedNFT
type ReqCreatedNFT struct {
	ClientID       string `json:"client_id"`
	NftCreatedTime int64  `json:"nft_created_time"`
	NfttxHash      string `json:"nft_tx_hash"`
	NftTokenId     string `json:"nft_token_id"`
	NftStoredUrl   string `json:"nft_store_url"`
}

// RewardResult
type RewardResult struct {
	ClientID         string   `json:"client_id"`
	ClientEthAddress string   `json:"client_eth_address"`
	GrantSuccess     string   `json:"grant_success"`
	GrantTokenAmount *big.Int `json:"grant_token_amount"`
	RewardReason     string   `json:"reward_reason"`
	MessageKey       string   `json:"message_key"`
	MessageTime      int64    `json:"message_time"`
	RewardTime       int64    `json:"reward_time"`
	AppVersion       string   `json:"client_version_build"`
}

// RewardingReq
type RewardingReq struct {
	ClientID     string `json:"client_id"`
	GrantSuccess string `json:"grant_success"`
	//RewardReason string `json:"reward_reason"`
	TimeFrom int64 `json:"time_from"`
	TimeTo   int64 `json:"time_to"`
}

// RewardSum
type RewardSum struct {
	RewardReason      string   `json:"reward_reason"`
	GrantTokenAmounts *big.Int `json:"grant_token_amount_subtotals"`
}

// RewardSum
type RewardSummary struct {
	ClientID           string   `json:"client_id"`
	ClientEthAddress   string   `json:"client_eth_address"`
	RewardTokenSummary *big.Int `json:"reward_token_summary"`
}

// IPLoacation
type IPLoacation struct {
	PublicIp string `json:"public_ip"`
}

type PubInfoByIP struct {
	ReqPublicIP               string `json:"req_public_ip"`
	ContryShort               string `json:"country_short"`
	ContryLong                string `json:"country_long"`
	Region                    string `json:"region"`
	City                      string `json:"city"`
	FirstChoicePubHost        string `json:"first_choice_pub_host"`
	FirstChoicePubInviteCode  string `json:"first_choice_pub_invite_code"`
	SecondChoicePubHost       string `json:"second_choice_pub_host"`
	SecondChoicePubInviteCode string `json:"second_choice_pub_invite_code"`
}

type GameInfo struct {
	GameId           string `json:"game_id"`
	GameName         string `json:"game_name"`
	GameVersion      string `json:"game_version"`
	CoverPhoto       string `json:"game_cover_photo"`
	Thumbnail        string `json:"game_thumbnail"`
	GameType         string `json:"game_type"`
	Introduction     string `json:"game_introduction"`
	PlayContent      string `json:"game_play"`
	BannerPhotos     string `json:"game_banner_photos"`
	ResourceDownload string `json:"resource_download"`
	RegTime          int64  `json:"reg_time"`
}

type PlayHistory struct {
	EthAddress       string   `json:"wallet_address"`
	SubmitType       string   `json:"submit_type"`
	GameId           string   `json:"game_id"`
	PlayerId         string   `json:"player_id"`
	PlayerName       string   `json:"player_name"`
	PlayerMark       string   `json:"player_mark"`
	PhotoLink        string   `json:"player_play_photos"`
	SubmitTime       int64    `json:"submit_time"`
	Ssbid            string   `json:"ssb_id"`
	DealTag          string   `json:"deal_tag"`
	DealTime         int64    `json:"deal_time"`
	GrantTokenAmount *big.Int `json:"grant_token_amount"`
}
