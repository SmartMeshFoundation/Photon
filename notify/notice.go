package notify

import (
	"encoding/json"
)

/*
Level :
*/
type Level int

const (
	// LevelInfo :
	LevelInfo = iota
	// LevelWarn :
	LevelWarn
	// LevelError :
	LevelError
)

/*
Notice for mobile or app
*/
type Notice struct {
	Level Level  `json:"level"`
	Info  string `json:"info"`
}

const (
	//InfoTypeString 0 简单的string通知
	InfoTypeString = iota
	//InfoTypeTransferStatus 1 发起的交易状态发生了变化
	InfoTypeTransferStatus
	//InfoTypeChannelCallID 2 关于通道的操作,有了结果
	InfoTypeChannelCallID
	//InfoTypeChannelStatus 3 通道状态发生了变化,包括但不限于
	//balance
	//patner_balance
	//locked_amount
	//partner_locked_amount
	//state
	InfoTypeChannelStatus
)

//InfoStruct for notify to mobile
type InfoStruct struct {
	Type    int         `json:"type"` //InfoTypeString 表示Message是一个string,InfoTypeTransferStatus表示Message是TransferStatus
	Message interface{} `json:"message"`
}

/*
newNotice :
*/
func newNotice(level Level, info *InfoStruct) *Notice {
	n := &Notice{
		Level: level,
	}
	buf, err := json.Marshal(info)
	if err != nil {
		n.Info = "unknown info"
	} else {
		n.Info = string(buf)
	}
	return n
}
