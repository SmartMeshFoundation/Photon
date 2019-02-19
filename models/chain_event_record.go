package models

import "encoding/gob"

// ChainEventID 一个链上事件的唯一ID, txHash+logIndex
type ChainEventID [25]byte

// ChainEventStatus :
type ChainEventStatus string

// #nosec
const (
	ChainEventStatusDelivered = "delivered" // 该状态标志事件已经投递到service层处理过
)

// ChainEventRecord 保存收到的链上事件
type ChainEventRecord struct {
	ID          ChainEventID     `json:"id" storm:"id"`
	BlockNumber uint64           `json:"block_number" storm:"index"`
	Status      ChainEventStatus `json:"status"`
}

func init() {
	gob.Register(&ChainEventRecord{})
}
