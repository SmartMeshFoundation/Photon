package channel

import "math/big"

// ChannelDataDetail for user api
type ChannelDataDetail struct {
	ChannelIdentifier   string   `json:"channel_identifier"`
	OpenBlockNumber     int64    `json:"open_block_number"`
	PartnerAddress      string   `json:"partner_address"`
	Balance             *big.Int `json:"balance"`
	PartnerBalance      *big.Int `json:"partner_balance"`
	LockedAmount        *big.Int `json:"locked_amount"`
	PartnerLockedAmount *big.Int `json:"partner_locked_amount"`
	TokenAddress        string   `json:"token_address"`
	State               State    `json:"state"`
	StateString         string   `json:"state_string"`
	SettleTimeout       int      `json:"settle_timeout"`
	RevealTimeout       int      `json:"reveal_timeout"`
}
