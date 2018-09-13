package models

import (
	"fmt"
)

// CheckNoLock :
func (c1 *Channel) CheckNoLock() bool {
	if c1.LockedAmount != 0 || c1.PartnerLockedAmount != 0 {
		Logger.Printf("Check failed because channel %s has lock but expect no lock !!!\n", c1.Name)
		return false
	}
	return true
}

// CheckLockSelf :
func (c1 *Channel) CheckLockSelf(lockAmt int32) bool {
	if c1.LockedAmount != lockAmt {
		Logger.Printf("Check failed because channel %s LockedAmount=%d but expect LockedAmount=%d !!!\n",
			c1.Name, c1.LockedAmount, lockAmt)
		return false
	}
	return true
}

// CheckLockPartner :
func (c1 *Channel) CheckLockPartner(lockAmt int32) bool {
	if c1.PartnerLockedAmount != lockAmt {
		Logger.Printf("Check failed because channel %s PartnerLockedAmount=%d but expect PartnerLockedAmount=%d !!!\n",
			c1.Name, c1.PartnerLockedAmount, lockAmt)
		return false
	}
	return true
}

// CheckLockBoth :
func (c1 *Channel) CheckLockBoth(lockAmt int32) bool {
	if c1.PartnerLockedAmount != lockAmt || c1.LockedAmount != lockAmt {
		Logger.Printf("Check failed because channel %s LockedAmount=%d,PartnerLockedAmount=%d but expect LockedAmount,PartnerLockedAmount=%d !!!\n",
			c1.Name, c1.LockedAmount, c1.PartnerLockedAmount, lockAmt)
		return false
	}
	return true
}

// CheckSelfBalance :
func (c1 *Channel) CheckSelfBalance(balance int32) bool {
	if c1.Balance != balance {
		Logger.Printf("Check failed because channel %s Balance=%d but expect Balance=%d !!!\n", c1.Name, c1.Balance, balance)
		return false
	}
	return true
}

// CheckPartnerBalance :
func (c1 *Channel) CheckPartnerBalance(balance int32) bool {
	if c1.PartnerBalance != balance {
		Logger.Printf("Check failed because channel %s PartnerBalance=%d but expect PartnerBalance=%d !!!\n", c1.Name, c1.PartnerBalance, balance)
		return false
	}
	return true
}

// CheckEqualByPartnerNode :
func (c1 *Channel) CheckEqualByPartnerNode(env *TestEnv) bool {
	n1 := env.GetNodeByAddress(c1.SelfAddress)
	n2 := env.GetNodeByAddress(c1.PartnerAddress)
	c2 := n2.GetChannelWith(n1, c1.TokenAddress)
	if !c1.isEqualChannelData(c2) {
		if c1.SelfAddress == c2.SelfAddress {
			c2.switchChannel()
		}
		Logger.Printf("Check failed because channel %s not equal %s !!!\n", c1.Name, c2.Name)
		header := fmt.Sprintf("Channel data after CheckByTwoSize fail :")
		c1.Println(header)
		c2.Println(header)
		return false
	}
	return true
}

// CheckState :
func (c1 *Channel) CheckState(state int) bool {
	if state == c1.State {
		return true
	}
	return false
}

// isEqualChannelData compare two channel
func (c1 *Channel) isEqualChannelData(c2 *Channel) bool {
	if c1.TokenAddress != c2.TokenAddress {
		return false
	}
	if c1.SelfAddress != c2.SelfAddress {
		c2.switchChannel()
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
func (c1 *Channel) switchChannel() {
	c1.SelfAddress, c1.PartnerAddress = c1.PartnerAddress, c1.SelfAddress
	c1.Balance, c1.PartnerBalance = c1.PartnerBalance, c1.Balance
	c1.LockedAmount, c1.PartnerLockedAmount = c1.PartnerLockedAmount, c1.LockedAmount
}
