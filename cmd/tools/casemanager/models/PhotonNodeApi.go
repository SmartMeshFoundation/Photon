package models

import (
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
	models2 "github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// GetChannelWith :
func (node *PhotonNode) GetChannelWith(partnerNode *PhotonNode, tokenAddr string) *Channel {
	req := &models.Req{
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
	req := &models.Req{
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
	for _, n := range env.Nodes {
		if n.Running {
			n.UpdateMeshNetworkNodes(env.Nodes...)
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
	Data     string `json:"data"`
}

// Transfer send a transfer
func (node *PhotonNode) Transfer(tokenAddress string, amount int32, targetAddress string, isDirect bool) error {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		Fee:      0,
		IsDirect: isDirect,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 180,
	}
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("TransferApi %s err :%s", req.FullURL, err))
		return err
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("TransferApi %s err : http status=%d", req.FullURL, statusCode))
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
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 60,
	}
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SendTransApi err : http status=%d", statusCode))
	}
}

// SendTransSync send a transfer, should be instead of Transfer
func (node *PhotonNode) SendTransSync(tokenAddress string, amount int32, targetAddress string, isDirect bool) {
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
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SendTransApi err : http status=%d", statusCode))
	}
}

// SendTransWithData send a transfer, should be instead of Transfer
func (node *PhotonNode) SendTransWithData(tokenAddress string, amount int32, targetAddress string, isDirect bool, data string) {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		Fee:      0,
		IsDirect: isDirect,
		Data:     data,
		Sync:     true,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/transfers/" + tokenAddress + "/" + targetAddress,
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 300,
	}
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SendTransApi err : http status=%d", statusCode))
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
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SendTransApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SendTransApi err : http status=%d", statusCode))
	}
}

//SendTransWithSecret send a transfer
func (node *PhotonNode) SendTransWithSecret(tokenAddress string, amount int32, targetAddress string, secretSeed string) {
	p, err := json.Marshal(TransferPayload{
		Amount:   amount,
		Fee:      0,
		IsDirect: false,
		Secret:   utils.Sha3([]byte(secretSeed)).String(),
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
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("WithdrawApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("WithdrawApi err : http status=%d", statusCode))
	}
}

// Close :
func (node *PhotonNode) Close(channelIdentifier string) {
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
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("CloseApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("CloseApi err : http status=%d", statusCode))
	}
}

// Settle :
func (node *PhotonNode) Settle(channelIdentifier string) {
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
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SettleApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SettleApi err : http status=%d", statusCode))
	}
}

// CooperateSettle :
func (node *PhotonNode) CooperateSettle(channelIdentifier string) {
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
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("CloseApi err :%s", err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("CloseApi err : http status=%d", statusCode))
	}
}

// OpenChannel :
func (node *PhotonNode) OpenChannel(partnerAddress, tokenAddress string, balance, settleTimeout int64) error {
	type OpenChannelPayload struct {
		PartnerAddress string `json:"partner_address"`
		TokenAddress   string `json:"token_address"`
		Balance        int64  `json:"balance"`
		SettleTimeout  int64  `json:"settle_timeout"`
	}
	p, err := json.Marshal(OpenChannelPayload{
		PartnerAddress: partnerAddress,
		TokenAddress:   tokenAddress,
		Balance:        balance,
		SettleTimeout:  settleTimeout,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/channels",
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
func (node *PhotonNode) Deposit(channelIdentifier string, balance int64) error {
	type DepostitPayload struct {
		Balance int64 `json:"balance"`
	}
	p, err := json.Marshal(DepostitPayload{
		Balance: balance,
	})
	req := &Req{
		FullURL: node.Host + "/api/1/channels/" + channelIdentifier,
		Method:  http.MethodPatch,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("DepositApi %s err :%s", req.FullURL, err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("DepositApi %s err : http status=%d", req.FullURL, statusCode))
		return fmt.Errorf("http status=%d", statusCode)
	}
	return err
}

// UpdateMeshNetworkNodes :
func (node *PhotonNode) UpdateMeshNetworkNodes(nodes ...*PhotonNode) {
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
		if n.Running {
			payloads = append(payloads, UpdateMeshNetworkNodesPayload{
				Address: n.Address,
				IPPort:  n.Host[7:] + "0",
			})
		}
	}
	p, err := json.Marshal(payloads)
	req := &Req{
		FullURL: node.Host + "/api/1/updatenodes",
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 5,
	}
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("UpdateMeshNetworkNodes %s err :%s", req.FullURL, err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("UpdateMeshNetworkNodes %s err : http status=%d", req.FullURL, statusCode))
		return
	}
	return
}

// SetFeePolicy :
func (node *PhotonNode) SetFeePolicy(fp *models2.FeePolicy) error {
	req := &Req{
		FullURL: node.Host + "/api/1/fee_policy",
		Method:  http.MethodPost,
		Payload: marshal(fp),
		Timeout: time.Second * 20,
	}
	statusCode, _, err := req.Invoke()
	if err != nil {
		Logger.Println(fmt.Sprintf("SetFeePolicy %s err :%s", req.FullURL, err))
	}
	if statusCode != 200 {
		Logger.Println(fmt.Sprintf("SetFeePolicy %s err : http status=%d", req.FullURL, statusCode))
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
