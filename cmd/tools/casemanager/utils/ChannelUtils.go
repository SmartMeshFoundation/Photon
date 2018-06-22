package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
)

// GetChannelBetween n1 and n2 on token
func GetChannelBetween(n1 *models.RaidenNode, n2 *models.RaidenNode, tokenAddr string) *models.Channel {
	req := &models.Req{
		FullURL: n1.Host + "/api/1/channels",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 30,
	}
	_, body, err := req.Invoke()
	if err != nil {
		panic(err)
	}
	var nodeChannels []models.Channel
	json.Unmarshal(body, &nodeChannels)
	if len(nodeChannels) == 0 {
		return nil
	}
	for _, channel := range nodeChannels {
		if channel.PartnerAddress == n2.Address && channel.TokenAddress == tokenAddr {
			channel.SelfAddress = n1.Address
			channel.Name = "CD-" + n1.Name + "-" + n2.Name
			return &channel
		}
	}
	return nil
}

// IsEqualChannelData compare two channel
func IsEqualChannelData(c1 *models.Channel, c2 *models.Channel) bool {
	if c1.TokenAddress != c2.TokenAddress {
		return false
	}
	if c1.SelfAddress != c2.SelfAddress {
		SwitchChannel(c1)
	}
	if c1.SelfAddress != c2.SelfAddress || c1.PartnerAddress != c2.PartnerAddress {
		return false
	}
	if c1.Balance == c2.Balance && c1.PartnerBalance == c2.PartnerBalance && c1.LockedAmount == c2.LockedAmount && c1.PartnerLockedAmount == c2.PartnerLockedAmount {
		return true
	}
	return false
}

//SwitchChannel switch channel
func SwitchChannel(c1 *models.Channel) {
	c1.SelfAddress, c1.PartnerAddress = c1.PartnerAddress, c1.SelfAddress
	c1.Balance, c1.PartnerBalance = c1.PartnerBalance, c1.Balance
	c1.LockedAmount, c1.PartnerLockedAmount = c1.PartnerLockedAmount, c1.LockedAmount
}
