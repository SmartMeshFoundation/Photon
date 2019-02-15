package models

import (
	"encoding/gob"

	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

// TXInfoStatus tx的状态
type TXInfoStatus string

/* #nosec */
const (
	TXInfoStatusPending = "pending"
	TXInfoStatusSuccess = "success"
	TXInfoStatusFailed  = "failed"
)

// TXInfoType 类型
type TXInfoType string

/* #nosec */
const (
	TXInfoTypeDeposit            = "ChannelDeposit"
	TXInfoTypeClose              = "ChannelClose"
	TXInfoTypeSettle             = "ChannelSettle"
	TXInfoTypeCooperateSettle    = "CooperateSettle"
	TXInfoTypeUpdateBalanceProof = "UpdateBalanceProof"
	TXInfoTypeUnlock             = "Unlock"
	TXInfoTypePunish             = "Punish"
	TXInfoTypeWithdraw           = "Withdraw"
	TXInfoTypeApprove            = "Approve"
	TXInfoTypeRegiterSecret      = "RegisterSecret"
)

// TXInfo 记录已经提交到公链节点的tx信息
type TXInfo struct {
	TXHash             common.Hash   `json:"tx_hash"`
	ChannelIdentifier  common.Hash   `json:"channel_identifier"` // 结合OpenBlockNumber唯一确定一个通道
	OpenBlockNumber    int64         `json:"open_block_number"`
	Type               TXInfoType    `json:"type"`
	IsSelfCall         bool          `json:"is_self_call"` // 是否自己发起的
	TXParams           string        `json:"tx_params"`    // 保存调用tx的参数信息,json格式,内容根据TXType不同而不同,仅自己发起的部分tx里面会带此参数
	Status             TXInfoStatus  `json:"tx_status"`
	Events             []interface{} `json:"events"`               // 保存这个tx成功之后对应的所有事件
	PendingBlockNumber int64         `json:"pending_block_number"` // tx最终所在的块号
}

// String :
func (ti *TXInfo) String() string {
	buf, err := json.MarshalIndent(ti, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// ToTXInfoSerialization :
func (ti *TXInfo) ToTXInfoSerialization() *TXInfoSerialization {
	return &TXInfoSerialization{
		TXHash:             ti.TXHash[:],
		ChannelIdentifier:  ti.ChannelIdentifier[:],
		OpenBlockNumber:    ti.OpenBlockNumber,
		Type:               string(ti.Type),
		IsSelfCall:         ti.IsSelfCall,
		TXParams:           ti.TXParams,
		Status:             string(ti.Status),
		Events:             ti.Events,
		PendingBlockNumber: ti.PendingBlockNumber,
	}
}

// TXInfoSerialization :
type TXInfoSerialization struct {
	TXHash             []byte        `storm:"id"`
	ChannelIdentifier  []byte        `storm:"index"` // 结合OpenBlockNumber唯一确定一个通道
	OpenBlockNumber    int64         `storm:"index"`
	Type               string        `storm:"index"`
	IsSelfCall         bool          `storm:"index"` // 是否自己发起的
	TXParams           string        // 保存调用tx的参数信息,json格式,内容根据TXType不同而不同,仅自己发起的部分tx里面会带此参数
	Status             string        `storm:"index"`
	Events             []interface{} // 保存这个tx成功之后对应的所有事件
	PendingBlockNumber int64         `storm:"index"` // tx最终所在的块号
}

// ToTXInfo :
func (tis *TXInfoSerialization) ToTXInfo() *TXInfo {
	return &TXInfo{
		TXHash:             common.BytesToHash(tis.TXHash),
		ChannelIdentifier:  common.BytesToHash(tis.ChannelIdentifier),
		OpenBlockNumber:    tis.OpenBlockNumber,
		Type:               TXInfoType(tis.Type),
		IsSelfCall:         tis.IsSelfCall,
		TXParams:           tis.TXParams,
		Status:             TXInfoStatus(tis.Status),
		Events:             tis.Events,
		PendingBlockNumber: tis.PendingBlockNumber,
	}
}

// TXParams tx的参数,自己发起的tx会带上
type TXParams interface{}

func init() {
	gob.Register(&TXInfoSerialization{})
}
