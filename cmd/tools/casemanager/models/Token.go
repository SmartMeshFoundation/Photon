package models

import "github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"

// Token name and address
type Token struct {
	Name    string
	Address string
	Manager *contracts.ChannelManagerContract
	Token   *contracts.Token
}
