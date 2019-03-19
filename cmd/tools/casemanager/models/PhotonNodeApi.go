package models

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kataras/go-errors"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	photon "github.com/SmartMeshFoundation/Photon"

	"github.com/SmartMeshFoundation/Photon/log"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/pfsproxy"
)

// GetChannelWith :
func (node *PhotonNode) GetChannelWith(partnerNode *PhotonNode, tokenAddr string) *Channel {
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
	var nodeChannels []Channel
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
			channel.Name = "CD-" + node.Name + "-" + partnerNode.Name
			return &channel
		}
	}
	return nil
}

// GetChannels :
func (node *PhotonNode) GetChannels(tokenAddr string) []*Channel {
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
	var nodeChannels []Channel
	err = json.Unmarshal(body, &nodeChannels)
	if err != nil {
		log.Info(fmt.Sprintf("bodylen=%d,body=%s", len(body), string(body)))
		panic(err)
	}
	var channels []*Channel
	for _, channel := range nodeChannels {
		if channel.TokenAddress == tokenAddr {
			channel.SelfAddress = node.Address
			channel.Name = "CD-" + node.Name + "-"
			channels = append(channels, &channel)
		}
	}
	return channels
}

// IsRunning check by api address
func (node *PhotonNode) IsRunning() bool {
	req := &Req{
		FullURL: node.Host + "/api/1/address",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 3,
	}
	_, err := req.Invoke()
	if err != nil {
		return false
	}
	return true
}

// Shutdown check by api address
func (node *PhotonNode) Shutdown(env *TestEnv) {
	req := &Req{
		FullURL: node.Host + "/api/1/debug/shutdown",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 3,
	}
	go req.Invoke()
	time.Sleep(10 * time.Second)
	node.Running = false
	var nodes []*PhotonNode
	for _, n := range env.Nodes {
		if n.Running {
			nodes = append(nodes, n)
		}
	}
	for _, n := range env.Nodes {
		if n.Running {
			//n.UpdateMeshNetworkNodes(nodes...)
		}
	}
	time.Sleep(time.Second * 2)
	return
}

// TransferPayload API  http body
type TransferPayload struct {
	Amount    int32                       `json:"amount"`
	IsDirect  bool                        `json:"is_direct"`
	Secret    string                      `json:"secret"`
	Sync      bool                        `json:"sync"`
	RouteInfo []pfsproxy.FindPathResponse `json:"route_info,omitempty"`
	Data      string                      `json:"data"`
}

// Transfer send a transfer
func (node *PhotonNode) Transfer(tokenAddress string, amount int32, targetAddress string, isDirect bool) error {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		IsDirect: isDirect,
		Sync:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 180,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TransferApi %s err :%s,body=%s", req.FullURL, err, string(body)))
		return err
	}
	return nil
}

// SendTransWithRouteInfo send a transfer with route info from pfs
func (node *PhotonNode) SendTransWithRouteInfo(target *PhotonNode, tokenAddress string, amount int32, routeInfo []pfsproxy.FindPathResponse) {
	if routeInfo == nil || len(routeInfo) == 0 {
		routeInfo = node.FindPath(target, common.HexToAddress(tokenAddress), amount)
	}
	p, err := json.Marshal(TransferPayload{
		Amount:    amount,
		IsDirect:  false,
		Sync:      true,
		RouteInfo: routeInfo,
		Data:      "test",
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + target.Address,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 60,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransWithRouteInfo err :%s,body=%s", err, string(body)))
	}
}

// SendTrans send a transfer, should be instead of Transfer
func (node *PhotonNode) SendTrans(tokenAddress string, amount int32, targetAddress string, isDirect bool) {
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
		Logger.Println(fmt.Sprintf("SendTransApi err :%s,body=%s", err, string(body)))
	}
}

//SendTransWithSecret send a transfer
func (node *PhotonNode) SendTransWithSecret(tokenAddress string, amount int32, targetAddress string, secretSeed string) {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		IsDirect: false,
		Secret:   secretSeed,
		Sync:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransWithSecretApi err :%s", err))
	}

}

// Withdraw :
func (node *PhotonNode) Withdraw(channelIdentifier string, withdrawAmount int32) {
	type WithdrawPayload struct {
		Amount int32
		Op     string
	}
	p, err := json.Marshal(WithdrawPayload{
		Amount: withdrawAmount,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/withdraw/" + channelIdentifier,
		Method:  http.MethodPut,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("WithdrawApi err :%s", err))
	}

}

// Close :
func (node *PhotonNode) Close(channelIdentifier string, waitSeconds ...int) (err error) {
	type ClosePayload struct {
		State string `json:"state"`
		Force bool   `json:"force"`
	}
	p, err := json.Marshal(ClosePayload{
		State: "closed",
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
		return fmt.Errorf("CloseApi err :%s", err)
	}
	Logger.Println(fmt.Sprintf("close channel returned=%s", string(body)))
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

// Settle :
func (node *PhotonNode) Settle(channelIdentifier string, waitSeconds ...int) (err error) {
	type SettlePayload struct {
		State string `json:"state"`
	}
	p, err := json.Marshal(SettlePayload{
		State: "settled",
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
	Logger.Println(fmt.Sprintf("settle channel returned=%s", string(body)))
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

// CooperateSettle : 由于CooperateSettle,close,settle,withdraw都是异步调用,因此必须再次封装
func (node *PhotonNode) CooperateSettle(channelIdentifier string, waitSeconds ...int) (err error) {
	type ClosePayload struct {
		State string `json:"state"`
	}
	p, err := json.Marshal(ClosePayload{
		State: "closed",
	})
	req := &Req{
		FullURL: node.Host + "/api/1/channels/" + channelIdentifier,
		Method:  http.MethodPatch,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	if err != nil {
		return fmt.Errorf("CooperateSettle err :%s", err)
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
		_, err = node.SpecifiedChannel(channelIdentifier)
		//找不到这个channel了才返回
		if err != nil {
			break
		}
	}
	if i == ws {
		return errors.New("timeout")
	}
	return nil
}

// OpenChannel :
func (node *PhotonNode) OpenChannel(partnerAddress, tokenAddress string, balance, settleTimeout int64, waitSeconds ...int) error {
	type OpenChannelPayload struct {
		PartnerAddress string `json:"partner_address"`
		TokenAddress   string `json:"token_address"`
		Balance        int64  `json:"balance"`
		SettleTimeout  int64  `json:"settle_timeout"`
		NewChannel     bool   `json:"new_channel"`
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
		Logger.Println(fmt.Sprintf("OpenChannelApi %s err :   body=%s err=%s", req.FullURL, string(body), err.Error()))
		return err
	}
	Logger.Println(fmt.Sprintf("open channel returned=%s", string(body)))
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

// Deposit :
func (node *PhotonNode) Deposit(partnerAddress, tokenAddress string, balance int64, waitSeconds ...int) error {
	type OpenChannelPayload struct {
		PartnerAddress string `json:"partner_address"`
		TokenAddress   string `json:"token_address"`
		Balance        int64  `json:"balance"`
		SettleTimeout  int64  `json:"settle_timeout"`
		NewChannel     bool   `json:"new_channel"`
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
		Logger.Println(fmt.Sprintf("DepositApi %s err :%s", req.FullURL, err))
		return err
	}
	Logger.Println(fmt.Sprintf("Deposit returned=%s", string(body)))
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

// UpdateMeshNetworkNodes :
func (node *PhotonNode) UpdateMeshNetworkNodes(nodes ...*PhotonNode) {
	Logger.Printf("UpdateMeshNetworkNodes for %s", node.Name)
	for _, n := range nodes {
		Logger.Printf("%s is online", n.Name)
	}
	type UpdateMeshNetworkNodesPayload struct {
		Address    string `json:"address"`
		IPPort     string `json:"ip_port"`
		DeviceType string `json:"device_type"` // must be mobile?
	}
	var payloads []UpdateMeshNetworkNodesPayload
	if len(nodes) == 0 {
		return
	}
	for _, n := range nodes {
		payloads = append(payloads, UpdateMeshNetworkNodesPayload{
			Address: n.Address,
			IPPort:  n.Host[7:] + "0",
		})
	}
	p, err := json.Marshal(payloads)
	req := &Req{
		FullURL: node.Host + "/api/1/updatenodes",
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 5,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("UpdateMeshNetworkNodes %s err :%s", req.FullURL, err))
	}
	return
}

// SetFeePolicy :
func (node *PhotonNode) SetFeePolicy(fp *models.FeePolicy) error {
	req := &Req{
		FullURL: node.Host + "/api/1/fee_policy",
		Method:  http.MethodPost,
		Payload: marshal(fp),
		Timeout: time.Second * 20,
	}
	_, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SetFeePolicy %s err :%s", req.FullURL, err))
	}

	return err
}

func marshal(v interface{}) string {
	p, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(p)
}

// AllowSecret :
func (node *PhotonNode) AllowSecret(secretHash, token string) {
	type AllowRevealSecretPayload struct {
		LockSecretHash string `json:"lock_secret_hash"`
		TokenAddress   string `json:"token_address"`
	}
	p, err := json.Marshal(AllowRevealSecretPayload{
		LockSecretHash: secretHash,
		TokenAddress:   token,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/allowrevealsecret",
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("AllowSecret err :%s", err))
	}
}

// GenerateSecret :
func (node *PhotonNode) GenerateSecret() (secret, secretHash string, err error) {
	type SecretPair struct {
		LockSecretHash string `json:"lock_secret_hash"`
		Secret         string `json:"secret"`
	}
	rs := SecretPair{}
	req := &Req{
		FullURL: node.Host + "/api/1/secret",
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GenerateSecret err :%s", err))
		return
	}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return
	}
	secret = rs.Secret
	secretHash = rs.LockSecretHash
	return
}

//GetSentTransferDetails query node's sent transfer
func (node *PhotonNode) GetSentTransferDetails() (trs []*models.SentTransferDetail, err error) {
	req := &Req{
		FullURL: node.Host + "/api/1/querysenttransfer",
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetSentTransferDetails err :%s", err))
		return
	}
	err = json.Unmarshal(body, &trs)
	if err != nil {
		return
	}
	return
}

//GetReceivedTransfers query node's received transfer
func (node *PhotonNode) GetReceivedTransfers() (trs []*models.ReceivedTransfer, err error) {
	req := &Req{
		FullURL: node.Host + "/api/1/queryreceivedtransfer",
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetReceivedTransfers err :%s", err))
		return
	}
	err = json.Unmarshal(body, &trs)
	if err != nil {
		return
	}
	return
}

//GetSentTransferDetail :
func (node *PhotonNode) GetSentTransferDetail(token, locksecrethash string) (status *models.SentTransferDetail, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/transferstatus/%s/%s", token, locksecrethash),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetSentTransferDetail err :%s", err))
		return
	}
	err = json.Unmarshal(body, &status)
	if err != nil {
		return
	}
	return
}

//CancelTransfer cancel a on transfer which secret is not revealed
func (node *PhotonNode) CancelTransfer(token, locksecrethash string) (err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/transfercancel/%s/%s", token, locksecrethash),
		Method:  http.MethodPost,
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("CancelTransfer err :%s", err))
		return
	}
	return nil
}

//GetUnfinishedReceivedTransfer query unfinished received transfers
func (node *PhotonNode) GetUnfinishedReceivedTransfer(token, locksecrethash string) (resp *photon.TransferDataResponse, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/getunfinishedreceivedtransfer/%s/%s", token, locksecrethash),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetUnfinishedReceivedTransfer err :%s", err))
		return
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	return
}

//Tokens : query registered tokens
func (node *PhotonNode) Tokens() (tokens []string, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host + "/api/1/tokens"),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("Tokens err :%s", err))
		return
	}
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		return
	}
	return

}

//PartnersDataResponse query by token
type PartnersDataResponse struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

//TokenPartners query token partners
func (node *PhotonNode) TokenPartners(token string) (partners []*PartnersDataResponse, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/tokens/%s/partners", token),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TokenPartners err :%s", err))
		return
	}
	err = json.Unmarshal(body, &partners)
	if err != nil {
		return
	}
	return

}

//PrepareUpdate query token partners
func (node *PhotonNode) PrepareUpdate() (err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host + "/api/1/prepare-update"),
		Method:  http.MethodPost,
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("PrepareUpdate err :%s", err))
		return
	}
	return

}

// TokenSwap send a transfer
func (node *PhotonNode) TokenSwap(target, locksecrethash, sendingtoken, receivingtoken, role, secret string,
	sendingAmount, receivingAmount int, routeInfo []pfsproxy.FindPathResponse) error {
	type TokenSwapPayload struct {
		Role            string                      `json:"role"`
		SendingAmount   int                         `json:"sending_amount"`
		SendingToken    string                      `json:"sending_token"`
		ReceivingAmount int                         `json:"receiving_amount"`
		ReceivingToken  string                      `json:"receiving_token"`
		Secret          string                      `json:"secret"` // taker无需填写,maker必填,且hash值需与url参数中的locksecrethash匹配,算法为SHA3
		RouteInfo       []pfsproxy.FindPathResponse `json:"route_info"`
	}
	p, err := json.Marshal(TokenSwapPayload{
		Role:            role,
		SendingAmount:   sendingAmount,
		SendingToken:    sendingtoken,
		ReceivingAmount: receivingAmount,
		ReceivingToken:  receivingtoken,
		Secret:          secret,
		RouteInfo:       routeInfo,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/token_swaps/" + target + "/" + locksecrethash,
		Method:  http.MethodPut,
		Payload: string(p),
		Timeout: time.Second * 30,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TransferApi %s err :%s", req.FullURL, err))
		return err
	}

	return nil
}

//SwitchNetwork disable mediated transfer
func (node *PhotonNode) SwitchNetwork(tomesh string) (err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/switch/%s", tomesh),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	_, err = req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SwitchNetwork err :%s", err))
		return
	}
	return

}

//TokenBalance query this account's balance of this token
func (node *PhotonNode) TokenBalance(token string) (v int, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/debug/balance/%s/%s", token, node.Address),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TokenBalance err :%s", err))
		return
	}
	log.Trace(string(body))
	b := new(big.Int)
	b.SetString(string(body), 0)
	v = int(b.Int64())
	return
}

//SpecifiedChannel query channel's detail
func (node *PhotonNode) SpecifiedChannel(channelIdentifier string) (c channeltype.ChannelDataDetail, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/channels/%s", channelIdentifier),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TokenPartners err :%s", err))
		return
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return
	}
	return

}

//ForceUnlock  unlock a unlock whenever i send annoucedisposed or not
func (node *PhotonNode) ForceUnlock(channelIdentifier string, secret string) (err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/debug/force-unlock/%s/%s", channelIdentifier, secret),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TokenPartners err :%s", err))
		return
	}
	Logger.Printf("forunlock body=%s", string(body))
	return

}

//RegisterSecret  register a secret to contract
func (node *PhotonNode) RegisterSecret(secret string) (err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/debug/register-secret-onchain/%s", secret),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("RegisterSecret err :%s", err))
		return
	}
	Logger.Printf("RegisterSecret body=%s", string(body))
	return

}

// FindPath :
func (node *PhotonNode) FindPath(target *PhotonNode, tokenAddress common.Address, amount int32) (path []pfsproxy.FindPathResponse) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/path/%s/%s/%d", target.Address, tokenAddress.String(), amount),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("FindPath err :%s", err))
		return
	}
	var resp []pfsproxy.FindPathResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		Logger.Println(fmt.Sprintf("FindPath err :%s", err))
		return
	}
	path = resp
	Logger.Printf("FindPath get RouteInfo from %s to %s on token %s :\n%s", node.Name, target.Name, tokenAddress.String(), MarshalIndent(path))
	return
}
