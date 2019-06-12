package codefortest

import (
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/ethereum/go-ethereum/common"
)

//MockDb implement XMPPDb
type MockDb struct {
	channels []*channeltype.Serialization
}

//AddPartner implement XMPPDb
func (db *MockDb) AddPartner(address common.Address) {
	db.channels = append(db.channels, &channeltype.Serialization{
		PartnerAddressBytes: address[:],
	})
}

//XMPPIsAddrSubed implement XMPPDb
func (db *MockDb) XMPPIsAddrSubed(addr common.Address) bool {
	return true
}

//XMPPMarkAddrSubed implement XMPPDb
func (db *MockDb) XMPPMarkAddrSubed(addr common.Address) {
	return
}

//GetChannelList implement XMPPDb
func (db *MockDb) GetChannelList(token, partner common.Address) (cs []*channeltype.Serialization, err error) {
	return db.channels, nil
}

//RegisterNewChannelCallback implement XMPPDb
func (db *MockDb) RegisterNewChannelCallback(f cb.ChannelCb) {

}

//RegisterChannelStateCallback implement XMPPDb
func (db *MockDb) RegisterChannelStateCallback(f cb.ChannelCb) {

}

//XMPPUnMarkAddr implement XMPPDb
func (db *MockDb) XMPPUnMarkAddr(addr common.Address) {

}

//RegisterChannelSettleCallback implement XMPPDb
func (db *MockDb) RegisterChannelSettleCallback(f cb.ChannelCb) {

}
