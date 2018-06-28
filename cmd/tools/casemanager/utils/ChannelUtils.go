package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
)

// CheckChannelNoLock :
func CheckChannelNoLock(self *models.RaidenNode, partner *models.RaidenNode, tokenAddr string, printMethod func(*models.Channel) *models.Channel) bool {
	c1 := GetChannelBetween(self, partner, tokenAddr)
	c2 := GetChannelBetween(partner, self, tokenAddr)
	if !IsEqualChannelData(c1, c2) {
		models.Logger.Printf("Check failed because channel %s not equal %s !!!\n", c1.Name, c2.Name)
		if printMethod != nil {
			printMethod(c1)
			printMethod(c2)
		}
		return false
	}
	if printMethod != nil {
		printMethod(c1)
	}
	if c1.LockedAmount != 0 || c1.PartnerLockedAmount != 0 {
		models.Logger.Printf("Check failed because channel %s has lock but expect no lock !!!\n", c1.Name)
		return false
	}
	return true
}

// CheckChannelLockPartner :
func CheckChannelLockPartner(self *models.RaidenNode, partner *models.RaidenNode, tokenAddr string, lockAmt int32, printMethod func(*models.Channel) *models.Channel) bool {
	c1 := GetChannelBetween(self, partner, tokenAddr)
	c2 := GetChannelBetween(partner, self, tokenAddr)
	if !IsEqualChannelData(c1, c2) {
		models.Logger.Printf("Check failed because channel %s not equal %s !!!\n", c1.Name, c2.Name)
		if printMethod != nil {
			printMethod(c1)
			printMethod(c2)
		}
		return false
	}
	if printMethod != nil {
		printMethod(c1)
	}
	if c1.PartnerLockedAmount != lockAmt {
		models.Logger.Printf("Check failed because channel %s PartnerLockedAmount=%d but expect PartnerLockedAmount=%d !!!\n", c1.Name, c1.PartnerLockedAmount, lockAmt)
		return false
	}
	return true
}

// CheckChannelLockBoth :
func CheckChannelLockBoth(self *models.RaidenNode, partner *models.RaidenNode, tokenAddr string, lockAmt int32, printMethod func(*models.Channel) *models.Channel) bool {
	c1 := GetChannelBetween(self, partner, tokenAddr)
	c2 := GetChannelBetween(partner, self, tokenAddr)
	if !IsEqualChannelData(c1, c2) {
		models.Logger.Printf("Check failed because channel %s not equal %s !!!\n", c1.Name, c2.Name)
		if printMethod != nil {
			printMethod(c1)
			printMethod(c2)
		}
		return false
	}
	if printMethod != nil {
		printMethod(c1)
	}
	if c1.PartnerLockedAmount != lockAmt && c1.LockedAmount != lockAmt {
		models.Logger.Printf("Check failed because channel %s LockedAmount,PartnerLockedAmount=%d but expect LockedAmount,PartnerLockedAmount=%d !!!\n", c1.Name, c1.PartnerLockedAmount, lockAmt)
		return false
	}
	return true
}

// CheckChannelPartnerBalance :
func CheckChannelPartnerBalance(self *models.RaidenNode, partner *models.RaidenNode, tokenAddr string, balance int32, printMethod func(*models.Channel) *models.Channel) bool {
	c1 := GetChannelBetween(self, partner, tokenAddr)
	c2 := GetChannelBetween(partner, self, tokenAddr)
	if !IsEqualChannelData(c1, c2) {
		models.Logger.Printf("Check failed because channel %s not equal %s !!!\n", c1.Name, c2.Name)
		if printMethod != nil {
			printMethod(c1)
			printMethod(c2)
		}
		return false
	}
	if printMethod != nil {
		printMethod(c1)
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

// PrintDataBeforeTransfer :
func PrintDataBeforeTransfer(c *models.Channel) *models.Channel {
	header := fmt.Sprintf("Channel data before transfer %s-BeforeTransfer :", c.Name)
	return c.Println(header)
}

// PrintDataAfterCrash :
func PrintDataAfterCrash(c *models.Channel) *models.Channel {
	header := fmt.Sprintf("Channel data after crash %s-AfterCrash :", c.Name)
	return c.Println(header)
}

// PrintDataAfterRestart :b
func PrintDataAfterRestart(c *models.Channel) *models.Channel {
	header := fmt.Sprintf("Channel data after restart %s-AfterRestart :", c.Name)
	return c.Println(header)
}
