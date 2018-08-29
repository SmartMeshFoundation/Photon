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
		Timeout: time.Second * 600,
	}
	statusCode, _, err := req.Invoke()
	if err != nil {
		panic(err)
	}
	if statusCode != 200 {
		panic(err)
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
	statusCode, _, err := req.Invoke()
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(statusCode, err)
	//if statusCode != 200 {
	//	panic(err)
	//}
}
