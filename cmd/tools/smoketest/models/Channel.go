package models

// Channel : raiden channel, add SelfChannel for test use
type Channel struct {
	ChannelIdentifier   string `json:"channel_identifier"`
	OpenBlockNumber     uint64 `json:"open_block_number"`
	PartnerAddress      string `json:"partner_address"`
	Balance             int32  `json:"balance"`
	PartnerBalance      int32  `json:"partner_balance"`
	LockedAmount        int32  `json:"locked_amount"`
	PartnerLockedAmount int32  `json:"partner_locked_amount"`
	TokenAddress        string `json:"token_address"`
	State               int    `json:"state"`
	SettleTimeout       int32  `json:"settle_timeout"`
	RevealTimeout       int32  `json:"reveal_timeout"`
	SelfAddress         string `json:"self_address"`
}
