package models

import (
	"encoding/json"
	"fmt"
)

// Channel atmosphere chainnode
type Channel struct {
	Name                string `json:"name"`
	SelfAddress         string `json:"self_address"`
	ChannelIdentifier   string `json:"channel_identifier"`
	PartnerAddress      string `json:"partner_address"`
	Balance             int32  `json:"balance"`
	LockedAmount        int32  `json:"locked_amount"`
	PartnerBalance      int32  `json:"partner_balance"`
	PartnerLockedAmount int32  `json:"partner_locked_amount"`
	TokenAddress        string `json:"token_address"`
	State               int    `json:"state"`
	SettleTimeout       int32  `json:"settle_timeout"`
	RevealTimeout       int32  `json:"reveal_timeout"`
}

// PrintDataBeforeTransfer :
func (c *Channel) PrintDataBeforeTransfer() *Channel {
	header := fmt.Sprintf("Channel data before transfer %s-BeforeTransfer :", c.Name)
	return c.Println(header)
}

// PrintDataAfterTransfer :
func (c *Channel) PrintDataAfterTransfer() *Channel {
	header := fmt.Sprintf("Channel data after transfer %s-AfterTransfer :", c.Name)
	return c.Println(header)
}

// PrintDataAfterCrash :
func (c *Channel) PrintDataAfterCrash() *Channel {
	header := fmt.Sprintf("Channel data after crash %s-AfterCrash :", c.Name)
	return c.Println(header)
}

// PrintDataAfterRestart :b
func (c *Channel) PrintDataAfterRestart() *Channel {
	header := fmt.Sprintf("Channel data after restart %s-AfterRestart :", c.Name)
	return c.Println(header)
}

// Println print data to console
func (c *Channel) Println(header string) *Channel {
	Logger.Println(header)
	buf, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(buf))
	return c
}
