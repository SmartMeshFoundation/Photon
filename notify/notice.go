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
	infoTypeString = iota
	infoTypeTransferStatus
)

//InfoStruct for notify to mobile
type InfoStruct struct {
	Type    int         `json:"info"` //infoTypeString 表示Message是一个string,InfoTypeTransferStatus表示Message是TransferStatus
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
