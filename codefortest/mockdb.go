package codefortest

import (
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/ethereum/go-ethereum/common"
)

type MockDb struct {
	channels []*channeltype.Serialization
}

func (db *MockDb) AddPartner(address common.Address) {
	db.channels = append(db.channels, &channeltype.Serialization{
		PartnerAddressBytes: address[:],
	})
}
func (db *MockDb) XMPPIsAddrSubed(addr common.Address) bool {
	return true
}
func (db *MockDb) XMPPMarkAddrSubed(addr common.Address) {
	return
}
func (db *MockDb) GetChannelList(token, partner common.Address) (cs []*channeltype.Serialization, err error) {
	return db.channels, nil
}
func (db *MockDb) RegisterNewChannelCallback(f cb.ChannelCb) {

}
func (db *MockDb) RegisterChannelStateCallback(f cb.ChannelCb) {

}
func (db *MockDb) XMPPUnMarkAddr(addr common.Address) {

}
func (db *MockDb) RegisterChannelSettleCallback(f cb.ChannelCb) {

}
