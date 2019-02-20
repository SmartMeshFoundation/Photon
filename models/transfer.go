package models

import (
	"encoding/gob"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//ReceivedTransfer tokens I have received and where it comes from
type ReceivedTransfer struct {
	Key               string `storm:"id"`
	BlockNumber       int64  `json:"block_number" storm:"index"`
	OpenBlockNumber   int64
	ChannelIdentifier common.Hash    `json:"channel_identifier"`
	TokenAddress      common.Address `json:"token_address"`
	TokenAddressBytes []byte         `json:"-"`
	FromAddress       common.Address `json:"initiator_address"`
	Nonce             uint64         `json:"nonce"`
	Amount            *big.Int       `json:"amount"`
	Data              string         `json:"data"`
	TimeStamp         int64          `json:"time_stamp"`
}

func init() {
	gob.Register(&ReceivedTransfer{})
}
