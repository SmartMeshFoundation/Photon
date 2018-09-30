package models

import (
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// GetChannelWith :
func (node *RaidenNode) GetChannelWith(partnerNode *RaidenNode, tokenAddr string) *Channel {
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

// IsRunning check by api address
func (node *RaidenNode) IsRunning() bool {
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
		panic("Exception response")
	}
	return true
}

// TransferPayload API  http body
type TransferPayload struct {
	Amount   int32  `json:"amount"`
	Fee      int64  `json:"fee"`
	IsDirect bool   `json:"is_direct"`
	Secret   string `json:"secret"`
}

// SendTrans send a transfer
func (node *RaidenNode) SendTrans(tokenAddress string, amount int32, targetAddress string, isDirect bool) {
	p, _ := json.Marshal(TransferPayload{
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

//SendTransWithSecret send a transfer
func (node *RaidenNode) SendTransWithSecret(tokenAddress string, amount int32, targetAddress string, secretSeed string) {
	p, _ := json.Marshal(TransferPayload{
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
func (node *RaidenNode) Withdraw(channelIdentifier string, withdrawAmount int32) {
	type WithdrawPayload struct {
		Amount int32
		Op     string
	}
	p, _ := json.Marshal(WithdrawPayload{
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
func (node *RaidenNode) Close(channelIdentifier string) {
	type ClosePayload struct {
		State string `json:"state"`
		Force bool   `json:"force"`
	}
	p, _ := json.Marshal(ClosePayload{
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
func (node *RaidenNode) Settle(channelIdentifier string) {
	type SettlePayload struct {
		State string `json:"state"`
	}
	p, _ := json.Marshal(SettlePayload{
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
func (node *RaidenNode) CooperateSettle(channelIdentifier string) {
	type ClosePayload struct {
		State string `json:"state"`
	}
	p, _ := json.Marshal(ClosePayload{
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
