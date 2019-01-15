package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon"

	"github.com/SmartMeshFoundation/Photon/log"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/models"
)

// GetChannelWith :
func (node *PhotonNode) GetChannelWith(partnerNode *PhotonNode, tokenAddr string) *Channel {
	req := &Req{
		FullURL: node.Host + "/api/1/channels",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 30,
	}
	_, body, err := req.Invoke()
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
	_, body, err := req.Invoke()
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
	statusCode, _, err := req.Invoke()
	if err != nil {
		return false
	}
	if statusCode != 200 {
		Logger.Printf("Exception response:%d\n", statusCode)
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
			n.UpdateMeshNetworkNodes(nodes...)
		}
	}
	return
}

// TransferPayload API  http body
type TransferPayload struct {
	Amount   int32  `json:"amount"`
	Fee      int64  `json:"fee"`
	IsDirect bool   `json:"is_direct"`
	Secret   string `json:"secret"`
	Sync     bool   `json:"sync"`
}

// Transfer send a transfer
func (node *PhotonNode) Transfer(tokenAddress string, amount int32, targetAddress string, isDirect bool) error {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		Fee:      0,
		IsDirect: isDirect,
		Sync:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 180,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TransferApi %s err :%s,body=%s", req.FullURL, err, string(body)))
		return err
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("TransferApi %s err : http status=%d,body=%s", req.FullURL, statusCode, string(body)))
		return fmt.Errorf("TransferApi err : http status=%d", statusCode)
	}
	return nil
}

// SendTrans send a transfer, should be instead of Transfer
func (node *PhotonNode) SendTrans(tokenAddress string, amount int32, targetAddress string, isDirect bool) {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		Fee:      0,
		IsDirect: isDirect,
		Sync:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 60,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransApi err :%s,body=%s", err, string(body)))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SendTransApi err : http status=%d,body=%s", statusCode, string(body)))
	}
}

// SendTransSyncWithFee send a transfer, should be instead of Transfer
func (node *PhotonNode) SendTransSyncWithFee(tokenAddress string, amount int32, targetAddress string, isDirect bool, fee int64) {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		Fee:      fee,
		IsDirect: isDirect,
		Sync:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 60,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SendTransApi err : http status=%d,body=%s", statusCode, string(body)))
	}
}

//SendTransWithSecret send a transfer
func (node *PhotonNode) SendTransWithSecret(tokenAddress string, amount int32, targetAddress string, secretSeed string) {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		Fee:      0,
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransWithSecretApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SendTransWithSecretApi err : http status=%d, body=%s", statusCode, string(body)))
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("WithdrawApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("WithdrawApi err : http status=%d,body=%s", statusCode, string(body)))
	}
}

// Close :
func (node *PhotonNode) Close(channelIdentifier string) (err error) {
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		return fmt.Errorf("CloseApi err :%s", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("CloseApi err : http status=%d,body=%s", statusCode, string(body))
	}
	return nil
}

// Settle :
func (node *PhotonNode) Settle(channelIdentifier string) (err error) {
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		return fmt.Errorf("SettleApi err :%s", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("SettleApi err : http status=%d,body=%s", statusCode, string(body))
	}
	return nil
}

// CooperateSettle :
func (node *PhotonNode) CooperateSettle(channelIdentifier string) (err error) {
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		return fmt.Errorf("CooperateSettle err :%s", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("CooperateSettle err : http status=%d,body=%s", statusCode, string(body))
	}
	return nil
}

// OpenChannel :
func (node *PhotonNode) OpenChannel(partnerAddress, tokenAddress string, balance, settleTimeout int64) error {
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("OpenChannelApi %s err : http status=%d body=%s err=%s", req.FullURL, statusCode, string(body), err.Error()))
		return err
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("OpenChannelApi %s err : http status=%d body=%s ", req.FullURL, statusCode, string(body)))
		return fmt.Errorf("http status=%d", statusCode)
	}
	return err
}

// Deposit :
func (node *PhotonNode) Deposit(partnerAddress, tokenAddress string, balance int64) error {
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("DepositApi %s err :%s", req.FullURL, err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("DepositApi %s err : http status=%d,body=%s", req.FullURL, statusCode, string(body)))
		return fmt.Errorf("http status=%d", statusCode)
	}
	return err
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("UpdateMeshNetworkNodes %s err :%s", req.FullURL, err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("UpdateMeshNetworkNodes %s err : http status=%d,body=%s", req.FullURL, statusCode, string(body)))
		return
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SetFeePolicy %s err :%s", req.FullURL, err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SetFeePolicy %s err : http status=%d,body=%s", req.FullURL, statusCode, string(body)))
		return fmt.Errorf("http status=%d", statusCode)
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("AllowSecret err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("AllowSecret err : http status=%d,body=%s", statusCode, string(body)))
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GenerateSecret err :%s", err))
		return
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("GenerateSecret err : http status=%d,body=%s", statusCode, string(body)))
		err = fmt.Errorf("errcode=%d", statusCode)
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

//GetSentTransfers query node's sent transfer
func (node *PhotonNode) GetSentTransfers() (trs []*models.SentTransfer, err error) {
	req := &Req{
		FullURL: node.Host + "/api/1/querysenttransfer",
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetSentTransfers err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("GetSentTransfers err : http status=%d", statusCode))
		err = fmt.Errorf("errcode=%d", statusCode)
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetReceivedTransfers err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("GetReceivedTransfers err : http status=%d", statusCode))
		err = fmt.Errorf("errcode=%d", statusCode)
		return
	}
	err = json.Unmarshal(body, &trs)
	if err != nil {
		return
	}
	return
}

//GetTransferStatus :
func (node *PhotonNode) GetTransferStatus(token, locksecrethash string) (status *models.TransferStatus, err error) {
	req := &Req{
		FullURL: fmt.Sprintf(node.Host+"/api/1/transferstatus/%s/%s", token, locksecrethash),
		Method:  http.MethodGet,
		Timeout: time.Second * 20,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetTransferStatus err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("GetTransferStatus err : http status=%d", statusCode))
		err = fmt.Errorf("errcode=%d", statusCode)
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("CancelTransfer err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("CancelTransfer err : http status=%d,body=%s", statusCode, string(body)))
		err = fmt.Errorf("errcode=%d,body=%s", statusCode, string(body))
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("GetUnfinishedReceivedTransfer err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("GetUnfinishedReceivedTransfer err : http status=%d", statusCode))
		err = fmt.Errorf("errcode=%d", statusCode)
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("Tokens err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("Tokens err : http status=%d", statusCode))
		err = fmt.Errorf("errcode=%d", statusCode)
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TokenPartners err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("TokenPartners err : http status=%d", statusCode))
		err = fmt.Errorf("errcode=%d", statusCode)
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("PrepareUpdate err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("PrepareUpdate err : http status=%d,body=%s", statusCode, string(body)))
		err = fmt.Errorf("errcode=%d", statusCode)
		return
	}
	return

}

// TokenSwap send a transfer
func (node *PhotonNode) TokenSwap(target, locksecrethash, sendingtoken, receivingtoken, role, secret string,
	sendingAmount, receivingAmount int) error {
	type TokenSwapPayload struct {
		Role            string `json:"role"`
		SendingAmount   int    `json:"sending_amount"`
		SendingToken    string `json:"sending_token"`
		ReceivingAmount int    `json:"receiving_amount"`
		ReceivingToken  string `json:"receiving_token"`
		Secret          string `json:"secret"` // taker无需填写,maker必填,且hash值需与url参数中的locksecrethash匹配,算法为SHA3
	}
	p, err := json.Marshal(TokenSwapPayload{
		Role:            role,
		SendingAmount:   sendingAmount,
		SendingToken:    sendingtoken,
		ReceivingAmount: receivingAmount,
		ReceivingToken:  receivingtoken,
		Secret:          secret,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/token_swaps/" + target + "/" + locksecrethash,
		Method:  http.MethodPut,
		Payload: string(p),
		Timeout: time.Second * 30,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TransferApi %s err :%s", req.FullURL, err))
		return err
	}
	if statusCode != http.StatusCreated {
		Logger.Println(fmt.Sprintf("TransferApi %s err : http status=%d", req.FullURL, statusCode))
		return fmt.Errorf("TransferApi err : http status=%d,body=%s", statusCode, string(body))
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
	statusCode, body, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SwitchNetwork err :%s", err))
		return
	}
	if statusCode < http.StatusOK || statusCode > http.StatusMultipleChoices {
		Logger.Println(fmt.Sprintf("SwitchNetwork err : http status=%d,body=%s", statusCode, string(body)))
		err = fmt.Errorf("errcode=%d", statusCode)
		return
	}
	return

}
