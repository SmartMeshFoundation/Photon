package restful

import (
	"encoding/json"
	"fmt"

	"errors"
	"math/big"
	"net/http"
	"time"

	"strings"

	"go.cryptoscope.co/ssb/restful/channel"
)

//var log kitlog.Logger

// PhotonNode a photon node
type PhotonNode struct {
	Host       string
	Address    string
	Name       string
	APIAddress string

	DebugCrash bool
	Running    bool
}

type TransferPayload struct {
	Amount   *big.Int `json:"amount"`
	IsDirect bool     `json:"is_direct"`
	Secret   string   `json:"secret"`
	Sync     bool     `json:"sync"`
	//RouteInfo []pfsproxy.FindPathResponse `json:"route_info,omitempty"`
	Data string `json:"data"`
}

// ChannelBigInt
type Channel struct {
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

	BlockNumberNow              int64 `json:"block_number_now"`
	BlockNumberChannelCanSettle int64 `json:"block_number_channel_can_settle,omitempty"`
}

// PhotonNodeRuntime
type PhotonNodeRuntime struct {
	MainChainBalance *big.Int // 主链货币余额
}

// GetChannels :
func (node *PhotonNode) GetChannels(tokenAddr string) ([]*Channel, error) {
	req := &Req{
		FullURL: node.Host + "/api/1/channels",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 30,
	}
	body, err := req.Invoke()
	if err != nil {
		return nil, err
	}
	//fmt.Println(fmt.Sprintf("[Pub]GetChannels returned=%s ", string(body)))
	var nodeChannels []*Channel
	err = json.Unmarshal(body, &nodeChannels)
	if err != nil {
		fmt.Println(fmt.Errorf("bodylen=%d,body=%s", len(body), string(body)))
		return nil, err
	}
	return nodeChannels, nil
}

// GetChannelWith :
func (node *PhotonNode) GetChannelWith(partnerNode *PhotonNode, tokenAddr string) (*Channel, error) {
	req := &Req{
		FullURL: node.Host + "/api/1/channels",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 30,
	}
	body, err := req.Invoke()
	if err != nil {
		return nil, err
	}
	var nodeChannels []Channel
	err = json.Unmarshal(body, &nodeChannels)
	if err != nil {
		fmt.Println(fmt.Sprintf("GetChannel Unmarshal err= %s", err))
		return nil, err
	}
	if len(nodeChannels) == 0 {
		return nil, nil
	}
	for _, channel := range nodeChannels {
		if channel.PartnerAddress == partnerNode.Address && channel.TokenAddress == tokenAddr {
			channel.SelfAddress = node.Address
			channel.Name = "CH-" + node.Name + "-" + partnerNode.Name
			return &channel, nil
		}
	}
	return nil, nil
}

// OpenChannel :
func (node *PhotonNode) OpenChannel(partnerAddress, tokenAddress string, balance *big.Int, settleTimeout int, waitSeconds ...int) error {
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
		fmt.Println(fmt.Sprintf("[Pub]OpenChannelApi err %s", err))
		return err
	}
	//fmt.Println(fmt.Sprintf("[Pub]OpenChannelApi returned %s", string(body)))
	ch := channel.ChannelDataDetail{}
	err = json.Unmarshal(body, &ch)
	if err != nil {
		fmt.Println(fmt.Sprintf("OpenChannel Unmarshal err= %s", err))
		return err
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
		//return errors.New("timeout")
		return errors.New("timeout")
	}
	return nil
}

func (node *PhotonNode) SpecifiedChannel(channelIdentifier string) (c channel.ChannelDataDetail, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/channels/%s", channelIdentifier),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		//fmt.Println(fmt.Sprintf("[Pub]SpecifiedChannel err %s", err))
		return
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return
	}
	return

}

func (node *PhotonNode) SendTrans(tokenAddress string, amount *big.Int, targetAddress string, isDirect bool, sync bool) error {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		IsDirect: isDirect,
		Sync:     sync,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 60,
	}
	_, err = req.Invoke()
	if err != nil {
		//fmt.Println(fmt.Sprintf("[Pub]SendTransApi err=%s,body=%s ", err, string(body)))
	}
	return err
}

func (node *PhotonNode) CheckChannelExist(partnerAddress, tokenAddress string) bool {
	//记录deposit之前的通道余额
	partners, err := node.TokenPartners(tokenAddress)
	if err != nil {
		fmt.Println(fmt.Sprintf("CheckChannelExist err :%s when get node-balance in this channel before deposit", err))
		return false
	}
	if len(partners) == 0 {
		fmt.Println(fmt.Sprintf("CheckChannelExist err :%s,no channel between %s and %s in token %s", err, node.Address, partnerAddress, tokenAddress))
		return false
	}
	channelInfo := ""
	for _, data := range partners {
		if data.PartnerAddress == partnerAddress {
			channelInfo = data.Channel
			break
		}
	}
	//"api/1/channles/0x9244a7c2bec98b59005656c5c98dba3ee394ccfd7710810a6af39929ca3d25a0"
	if strings.Count(channelInfo, "/") != 3 {
		return false
	}
	channelInfo = strings.Split(channelInfo, "/")[3]
	_, err = node.SpecifiedChannel(channelInfo)
	if err != nil {
		fmt.Println(fmt.Sprintf("CheckChannelExist err when get SpecifiedChannel:%s", err))
		return false
	}
	return true
}

func (node *PhotonNode) Deposit(partnerAddress, tokenAddress string, balance *big.Int, waitSeconds ...int) error {
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
	//记录deposit之前的通道余额
	partners, err := node.TokenPartners(tokenAddress)
	if err != nil {
		fmt.Println(fmt.Sprintf("DepositApi %s err :%s when get node-balance in this channel before deposit", req.FullURL, err))
		return err
	}
	if len(partners) == 0 {
		fmt.Println(fmt.Sprintf("DepositApi %s err :%s,no channel between %s and %s in token %s", req.FullURL, err, node.Address, partnerAddress, tokenAddress))
		return err
	}
	channelInfo := ""
	for _, data := range partners {
		if data.PartnerAddress == partnerAddress {
			channelInfo = data.Channel
			break
		}
	}
	//"api/1/channles/0x9244a7c2bec98b59005656c5c98dba3ee394ccfd7710810a6af39929ca3d25a0"
	if strings.Count(channelInfo, "/") != 3 {
		return errors.New("deposit error, check channel error")
	}
	channelInfo = strings.Split(channelInfo, "/")[3]
	c, err := node.SpecifiedChannel(channelInfo)
	if err != nil {
		fmt.Println(fmt.Sprintf("DepositApi %s err when get SpecifiedChannel:%s", req.FullURL, err))
		return err
	}
	nodeBalanceBeforeDeposit := c.Balance

	body, err := req.Invoke()
	if err != nil {
		//fmt.Println(fmt.Sprintf("[Pub]DepositApi err=%s ", err))
		return err
	}
	fmt.Println(fmt.Sprintf("[Pub]Deposit returned=%s ", string(body)))
	ch := channel.ChannelDataDetail{}
	err = json.Unmarshal(body, &ch)
	if err != nil {
		fmt.Println(fmt.Sprintf("Deposit Unmarshal err= %s", err))
		return err
	}
	var ws int
	if len(waitSeconds) > 0 {
		ws = waitSeconds[0]
	} else {
		ws = 45 //d等三块,应该会被打包进去的.
	}
	var i int
	for i = 0; i < ws; i++ {
		time.Sleep(time.Second * 1)
		_, err := node.SpecifiedChannel(ch.ChannelIdentifier)
		if err == nil {
			break
		}
	}
	if i == ws {
		return errors.New("timeout")
	}

	//UDP通信时间不定，deposit异步，需要验证通道余额
	var j int
	for j = 0; j < 90; j++ {
		time.Sleep(time.Second)
		cx, err := node.SpecifiedChannel(ch.ChannelIdentifier)
		/*fmt.Println(fmt.Sprintf("check (%d) partnerAddress=%s, balance of before\t:%v", j, partnerAddress, nodeBalanceBeforeDeposit))
		fmt.Println(fmt.Sprintf("check (%d) partnerAddress=%s, balance of balance\t:%v", j, partnerAddress, balance))
		fmt.Println(fmt.Sprintf("check (%d) partnerAddress=%s, balance of now\t:%v", j, partnerAddress, cx.Balance))*/
		if err == nil && cx.Balance.Cmp(new(big.Int).Add(nodeBalanceBeforeDeposit, balance)) == 0 {
			break
		}
	}
	if j == 90 {
		return errors.New("check result of Deposit timeout")
	}

	return nil
}

// Settle :
func (node *PhotonNode) Settle(channelIdentifier string, waitSeconds ...int) (err error) {
	type SettlePayload struct {
		State string `json:"state"`
		Force bool   `json:"force"`
	}
	p, err := json.Marshal(SettlePayload{
		State: "settled",
		Force: true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/channels/" + channelIdentifier,
		Method:  http.MethodPatch,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		return fmt.Errorf("SettleApi err :%s", err)
	}
	ch := channel.ChannelDataDetail{}
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
		time.Sleep(time.Second)
		_, err = node.SpecifiedChannel(ch.ChannelIdentifier)
		//找不到到这个channel了才返回
		if err != nil {
			break
		}
	}
	if i == ws {
		return errors.New("timeout")
	}
	return nil
}

// Close :
func (node *PhotonNode) Close(tokenAddress, partnerAddress string, waitSeconds ...int) (err error) {
	//
	partners, err := node.TokenPartners(tokenAddress)
	if err != nil {
		fmt.Println(fmt.Sprintf("Close Channel, Check TokenPartners err=%s, tokenAddress= %s", err, tokenAddress))
		return err
	}
	if len(partners) == 0 {
		fmt.Println(fmt.Sprintf("Close Channel, no channel between %s and %s in token %s", node.Address, partnerAddress, tokenAddress))
		return err
	}
	channelInfo := ""
	for _, data := range partners {
		if data.PartnerAddress == partnerAddress {
			channelInfo = data.Channel
			break
		}
	}
	//"api/1/channles/0x9244a7c2bec98b59005656c5c98dba3ee394ccfd7710810a6af39929ca3d25a0"
	if strings.Count(channelInfo, "/") != 3 {
		return errors.New("Close Channel, channel not exist")
	}
	channelInfo = strings.Split(channelInfo, "/")[3]
	//

	type ClosePayload struct {
		State string `json:"state"`
		Force bool   `json:"force"`
	}
	p, err := json.Marshal(ClosePayload{
		State: "closed",
		Force: true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/channels/" + channelInfo,
		Method:  http.MethodPatch,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		return fmt.Errorf("CloseApi err :%s", err)
	}
	fmt.Println(fmt.Sprintf("close channel returned=%s", string(body)))
	ch := channel.ChannelDataDetail{}
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
		time.Sleep(time.Second)
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

// PartnersDataResponse query by token
type PartnersDataResponse struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

// TokenPartners query token partners
func (node *PhotonNode) TokenPartners(token string) (partners []*PartnersDataResponse, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/tokens/%s/partners", token),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		fmt.Println(fmt.Sprintf("TokenPartners err :%s", err))
		return
	}
	err = json.Unmarshal(body, &partners)
	if err != nil {
		return
	}
	return
}

type repsNodeStatus struct {
	DeviceType string `json:"device_type"`
	IsOnline   bool   `json:"is_online"`
}

// GetNodeStatus query status of online, just for xmpp
func (node *PhotonNode) GetNodeStatus(nodeaddr string) (status *repsNodeStatus, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/node-status/%s", nodeaddr),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		fmt.Println(fmt.Sprintf("GetNodeStatus err :%s", err))
		return
	}
	err = json.Unmarshal(body, &status)
	if err != nil {
		return
	}
	return
}

// TokenTransfer
func (node *PhotonNode) TransferSMT(addr, value string) (err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/transfer-smt/%s/%s", addr, value),
		Method:  http.MethodPost,
		Timeout: time.Second * 180,
	}
	body, err := req.Invoke()
	if err != nil {
		fmt.Println(fmt.Sprintf("transfer-smt err :%s", err))
		return
	}
	err = json.Unmarshal(body, new(interface{}))
	if err != nil {
		return
	}
	return
}
