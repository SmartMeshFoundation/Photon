package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
)

// CheckChannelNoLock :
func CheckChannelNoLock(env *models.TestEnv, c1 *models.Channel) bool {
	if !checkByTwoSize(env, c1) {
		return false
	}
	if c1.LockedAmount != 0 || c1.PartnerLockedAmount != 0 {
		models.Logger.Printf("Check failed because channel %s has lock but expect no lock !!!\n", c1.Name)
		return false
	}
	return true
}

// CheckChannelLockPartner :
func CheckChannelLockPartner(env *models.TestEnv, c1 *models.Channel, lockAmt int32) bool {
	if !checkByTwoSize(env, c1) {
		return false
	}
	if c1.PartnerLockedAmount != lockAmt {
		models.Logger.Printf("Check failed because channel %s PartnerLockedAmount=%d but expect PartnerLockedAmount=%d !!!\n", c1.Name, c1.PartnerLockedAmount, lockAmt)
		return false
	}
	return true
}

// CheckChannelLockBoth :
func CheckChannelLockBoth(env *models.TestEnv, c1 *models.Channel, lockAmt int32) bool {
	if !checkByTwoSize(env, c1) {
		return false
	}
	if c1.PartnerLockedAmount != lockAmt && c1.LockedAmount != lockAmt {
		models.Logger.Printf("Check failed because channel %s LockedAmount,PartnerLockedAmount=%d but expect LockedAmount,PartnerLockedAmount=%d !!!\n", c1.Name, c1.PartnerLockedAmount, lockAmt)
		return false
	}
	return true
}

// CheckChannelPartnerBalance :
func CheckChannelPartnerBalance(env *models.TestEnv, c1 *models.Channel, balance int32) bool {
	if !checkByTwoSize(env, c1) {
		return false
	}
	if c1.PartnerBalance != balance {
		models.Logger.Printf("Check failed because channel %s PartnerBalance=%d but expect PartnerBalance=%d !!!\n", c1.Name, c1.PartnerBalance, balance)
		return false
	}
	return true
}

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

func checkByTwoSize(env *models.TestEnv, c1 *models.Channel) bool {
	c2 := GetChannelBetween(env.GetNodeByAddress(c1.PartnerAddress), env.GetNodeByAddress(c1.SelfAddress), c1.TokenAddress)
	if !IsEqualChannelData(c1, c2) {
		models.Logger.Printf("Check failed because channel %s not equal %s !!!\n", c1.Name, c2.Name)
		header := fmt.Sprintf("Channel data after case fail %s-CaseFail :", c1.Name)
		c1.Println(header)
		c2.Println(header)
		return false
	}
	return true
}
