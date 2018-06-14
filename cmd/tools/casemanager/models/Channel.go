package models

import (
	"encoding/json"
)

// Channel raiden chainnode
type Channel struct {
	Name                string `json:"name"`
	SelfAddress         string `json:"self_address"`
	ChannelAddress      string `json:"channel_address"`
	PartnerAddress      string `json:"partner_address"`
	Balance             int32  `json:"balance"`
	LockedAmount        int32  `json:"locked_amount"`
	PartnerBalance      int32  `json:"partner_balance"`
	PartnerLockedAmount int32  `json:"partner_locked_amount"`
	TokenAddress        string `json:"token_address"`
	State               string `json:"state"`
	SettleTimeout       int32  `json:"settle_timeout"`
	RevealTimeout       int32  `json:"reveal_timeout"`
}

// Println print data to console
func (c *Channel) Println(header string) {
	Logger.Println(header)
	buf, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(buf))
}
