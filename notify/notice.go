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
