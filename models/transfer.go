package models

import (
	"encoding/gob"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//SentTransfer transfer's I have sent and success.
type SentTransfer struct {
	Key               string         `storm:"id"`
	BlockNumber       int64          `json:"block_number" storm:"index"`
	OpenBlockNumber   int64          `json:"open_block_number"`
	ChannelIdentifier common.Hash    `json:"channel_identifier"`
	ToAddress         common.Address `json:"to_address"`
	TokenAddress      common.Address `json:"token_address"`
	Nonce             uint64         `json:"nonce"`
	Amount            *big.Int       `json:"amount"`
	Data              string         `json:"data"`
	TimeStamp         string         `json:"time_stamp"`
}

//ReceivedTransfer tokens I have received and where it comes from
type ReceivedTransfer struct {
	Key               string `storm:"id"`
	BlockNumber       int64  `json:"block_number" storm:"index"`
	OpenBlockNumber   int64
	ChannelIdentifier common.Hash    `json:"channel_identifier"`
	TokenAddress      common.Address `json:"token_address"`
	FromAddress       common.Address `json:"from_address"`
	Nonce             uint64         `json:"nonce"`
	Amount            *big.Int       `json:"amount"`
	Data              string         `json:"data"`
	TimeStamp         string         `json:"time_stamp"`
}

func init() {
	gob.Register(&SentTransfer{})
	gob.Register(&ReceivedTransfer{})
}
