package supernode

import (
	"encoding/json"
	"fmt"
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/kataras/go-errors"
	"math/big"
	"net/http"
	"time"
)

// PhotonNode a photon node
type SuperNode struct {
	Host          string
	Address       string
	Name          string
	APIAddress    string
	ListenAddress string
	ConditionQuit *params.ConditionQuit
	DebugCrash    bool
	Running       bool
	NoNetwork     bool
	DoPprof       bool
	Runtime       PhotonNodeRuntime
	PubApiHost    string
}

type TransferPayload struct {
	Amount    *big.Int                    `json:"amount"`
	IsDirect  bool                        `json:"is_direct"`
	Secret    string                      `json:"secret"`
	Sync      bool                        `json:"sync"`
	RouteInfo []pfsproxy.FindPathResponse `json:"route_info,omitempty"`
	Data      string                      `json:"data"`
}

// ChannelBigInt
type ChannelBigInt struct {
	Name                string   `json:"name"`
	SelfAddress         string   `json:"self_address"`
	ChannelIdentifier   string   `json:"channel_identifier"`
	PartnerAddress      string   `json:"partner_address"`
	Balance             *big.Int `json:"balance"`
	LockedAmount        *big.Int `json:"locked_amount"`
	PartnerBalance      *big.Int `json:"partner_balance"`
	PartnerLockedAmount *big.Int `json:"partner_locked_amount"`
	TokenAddress        string   `json:"token_address"`
	State               int      `json:"state"`
	SettleTimeout       *big.Int `json:"settle_timeout"`
	RevealTimeout       *big.Int `json:"reveal_timeout"`
}

// PhotonNodeRuntime
type PhotonNodeRuntime struct {
	MainChainBalance *big.Int // 主链货币余额
}

// LasterNumLikes
// VoteLink为被点赞的内容ID
// LasterAddVoteNum 为新增的点赞数量
// LasterAddVoteNum 收集到了总的点赞数量（如果发放奖励在先先，有取消点赞的，不收回奖励
type LasterNumLikes struct {
	ClientID         string `json:"client_id"`
	ClientAddress    string `json:"client_eth_address"`
	LasterAddVoteNum int64  `json:"laster_add_vote_num"`
	LasterVoteNum    int64  `json:"laster_vote_num"`
	//VoteLink         []string `json:"laster_add_vote_num"`
}

// GetChannelWithBigInt :
func (node *SuperNode) GetChannelWithBigInt(partnerNode *SuperNode, tokenAddr string) *ChannelBigInt {
	req := &Req{
		FullURL: node.Host + "/api/1/channels",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 30,
	}
	body, err := req.Invoke()
	if err != nil {
		panic(err)
	}
	var nodeChannels []ChannelBigInt
	err = json.Unmarshal(body, &nodeChannels)
	if err != nil {
		panic(err)
	}
	if len(nodeChannels) == 0 {
		return nil
	}
	for _, channel := range nodeChannels {
		if channel.PartnerAddress == partnerNode.Address && channel.TokenAddress == tokenAddr {
			channel.SelfAddress = node.Address
			channel.Name = "CH-" + node.Name + "-" + partnerNode.Name
			return &channel
		}
	}
	return nil
}

// OpenChannel :
func (node *SuperNode) OpenChannelBigInt(partnerAddress, tokenAddress string, balance *big.Int, settleTimeout int, waitSeconds ...int) error {
	type OpenChannelPayload struct {
		PartnerAddress string   `json:"partner_address"`
		TokenAddress   string   `json:"token_address"`
		Balance        *big.Int `json:"balance"`
		SettleTimeout  int      `json:"settle_timeout"`
		NewChannel     bool     `json:"new_channel"`
	}
	p, err := json.Marshal(OpenChannelPayload{
		PartnerAddress: partnerAddress,
		TokenAddress:   tokenAddress,
		Balance:        balance,
		SettleTimeout:  settleTimeout,
		NewChannel:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/deposit",
		Method:  http.MethodPut,
		Payload: string(p),
		Timeout: time.Second * 60,
	}
	body, err := req.Invoke()
	if err != nil {
		log.Error(fmt.Sprintf("[SuperNode]OpenChannelApi %s err :   body=%s err=%s", req.FullURL, string(body), err.Error()))
		return err
	}
	log.Info(fmt.Sprintf("[SuperNode]open channel returned=%s", string(body)))
	ch := channeltype.ChannelDataDetail{}
	err = json.Unmarshal(body, &ch)
	if err != nil {
		panic(err)
	}
	var ws int
	if len(waitSeconds) > 0 {
		ws = waitSeconds[0]
	} else {
		ws = 45 //d等三块,应该会被打包进去的.
	}
	var i int
	for i = 0; i < ws; i++ {
		time.Sleep(time.Second * 3)
		_, err = node.SpecifiedChannel(ch.ChannelIdentifier)
		//找到这个channel了才返回
		if err == nil {
			break
		}
	}
	if i == ws {
		return errors.New("timeout")
	}
	return nil
}

func (node *SuperNode) SpecifiedChannel(channelIdentifier string) (c channeltype.ChannelDataDetail, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/channels/%s", channelIdentifier),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		log.Error(fmt.Sprintf("[SuperNode]SpecifiedChannel err :%s", err))
		return
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return
	}
	return

}

//通过本节点查询其他节点的ssb账号、待付款
func (node *SuperNode) LatestNumberOfLikes() (lnum []*LasterNumLikes, err error) {
	req := &Req{
		FullURL: fmt.Sprintf("http://" + node.PubApiHost + "/ssb/api/likes"),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		log.Error(fmt.Sprintf("[SuperNode]getLatestNumberOfLikes err :%s", err))
		return
	}
	err = json.Unmarshal(body, &lnum)
	if err != nil {
		return
	}
	return
}

func (node *SuperNode) SendTransWithRouteInfo(tokenAddress string, amount *big.Int, targetAddress string, isDirect bool, routeInfo []pfsproxy.FindPathResponse) (err error) {
	if routeInfo == nil || len(routeInfo) == 0 {
		routeInfo, err = node.FindPath(targetAddress, tokenAddress, amount)
		if err != nil {
			return
		}
	}
	p, err := json.Marshal(TransferPayload{
		Amount:    amount,
		IsDirect:  false,
		Sync:      true,
		RouteInfo: routeInfo,
		Data:      "test",
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + tokenAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		log.Info(fmt.Sprintf("SendTransWithRouteInfo err :%s,body=%s", err, string(body)))
	}
	return
}

func (node *SuperNode) FindPath(target string, tokenAddress string, amount *big.Int) (path []pfsproxy.FindPathResponse, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/path/%s/%s/%v", target, tokenAddress, amount.String()),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		log.Error(fmt.Sprintf("FindPath err :%s", err))
		return
	}
	var resp []pfsproxy.FindPathResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Error(fmt.Sprintf("FindPath err :%s", err))
		return
	}
	path = resp
	log.Info(fmt.Sprintf("FindPath get RouteInfo from %s to %s on token %s :\n%s", node.Name, target, tokenAddress, MarshalIndent(path)))
	return
}

func MarshalIndent(v interface{}) string {
	buf, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (node *SuperNode) SendTrans(tokenAddress string, amount *big.Int, targetAddress string, isDirect bool) error {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		IsDirect: isDirect,
		Sync:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 60,
	}
	body, err := req.Invoke()
	if err != nil {
		log.Error(fmt.Sprintf("SendTransApi err :%s,body=%s", err, string(body)))
	}
	return err
}

func (node *SuperNode) Deposit(partnerAddress, tokenAddress string, balance *big.Int, waitSeconds ...int) error {
	type OpenChannelPayload struct {
		PartnerAddress string   `json:"partner_address"`
		TokenAddress   string   `json:"token_address"`
		Balance        *big.Int `json:"balance"`
		SettleTimeout  int64    `json:"settle_timeout"`
		NewChannel     bool     `json:"new_channel"`
	}
	p, err := json.Marshal(OpenChannelPayload{
		PartnerAddress: partnerAddress,
		TokenAddress:   tokenAddress,
		Balance:        balance,
		SettleTimeout:  0,
		NewChannel:     false,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/deposit",
		Method:  http.MethodPut,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		log.Error(fmt.Sprintf("[SuperNode]DepositApi %s err :%s", req.FullURL, err))
		return err
	}
	log.Info(fmt.Sprintf("[SuperNode]Deposit returned=%s", string(body)))
	ch := channeltype.ChannelDataDetail{}
	err = json.Unmarshal(body, &ch)
	if err != nil {
		panic(err)
	}
	var ws int
	if len(waitSeconds) > 0 {
		ws = waitSeconds[0]
	} else {
		ws = 45 //d等三块,应该会被打包进去的.
	}
	var i int
	for i = 0; i < ws; i++ {
		time.Sleep(time.Second * 3)
		_, err = node.SpecifiedChannel(ch.ChannelIdentifier)
		//找到这个channel了才返回
		if err == nil {
			break
		}
	}
	if i == ws {
		return errors.New("timeout")
	}
	return nil
}
