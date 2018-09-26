package models

import (
	"encoding/json"
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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

/*
newNotice :
*/
func newNotice(level Level, info interface{}) *Notice {
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

// Notify : 通知上层,不让阻塞,以免影响正常业务
func (model *ModelDB) Notify(level Level, info interface{}) {
	select {
	case model.NoticeChan <- newNotice(level, info):
	default:
		// never block
	}
}

// NotifyReceiveMediatedTransfer :
func (model *ModelDB) NotifyReceiveMediatedTransfer(msg *encoding.MediatedTransfer, ch *channel.Channel) {
	if msg == nil {
		return
	}
	info := fmt.Sprintf("收到token=%s,amount=%d,locksecrethash=%s的交易",
		utils.APex2(ch.TokenAddress), msg.PaymentAmount, utils.HPex(msg.LockSecretHash))
	select {
	case model.NoticeChan <- newNotice(LevelInfo, info):
	default:
		// never block
	}
}
