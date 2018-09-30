package models

// Token : save raiden token info
type Token struct {
	Address      string    `json:"address"`       // token地址
	Channels     []Channel `json:"channels"`      // 该token下所有Channel信息
	IsRegistered bool      `json:"is_registered"` // 是否已注册
}

// judge a channel exist in this token
func (t *Token) hasChannel(channelIdentifier string) bool {
	for _, channel := range t.Channels {
		if channel.ChannelIdentifier == channelIdentifier {
			return true
		}
	}
	return false
}
