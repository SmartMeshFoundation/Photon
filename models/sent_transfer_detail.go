package models

import (
	"encoding/gob"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// SentTransferDetail :
type SentTransferDetail struct {
	Key               string         `storm:"id"`
	BlockNumber       int64          `json:"block_number" storm:"index"`
	TokenAddressBytes []byte         `json:"-"  storm:"index"`
	TokenAddress      common.Address `json:"token_address"`
	LockSecretHash    common.Hash
	TargetAddress     common.Address     `json:"target_address"`
	Amount            *big.Int           `json:"amount"`
	Data              string             `json:"data"`
	IsDirect          bool               `json:"is_direct"`
	SendingTime       int64              `json:"sending_time" storm:"index"`
	FinishTime        int64              `json:"finish_time" storm:"index"`
	Status            TransferStatusCode `json:"status"`
	StatusMessage     string             `json:"status_message"`

	/*
		通道相关信息,如果为MediatorTransfer, 保存的是我与第一个mediator节点的通道上的信息,这部分信息仅交易成功才会有
	*/
	ChannelIdentifier common.Hash `json:"channel_identifier"`
	OpenBlockNumber   int64       `json:"open_block_number"`
}

func init() {
	gob.Register(&SentTransferDetail{})
}
