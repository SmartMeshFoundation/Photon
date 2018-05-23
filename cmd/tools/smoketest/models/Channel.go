package models

// Channel : raiden channel, add SelfChannel for test use
type Channel struct {
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
