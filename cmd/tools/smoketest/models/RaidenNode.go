package models

// RaidenNode description
type RaidenNode struct {
	AccountAddress string `json:"account_address"` // 账户地址
	Host           string `json:"api_address"`     // api服务Host
}
