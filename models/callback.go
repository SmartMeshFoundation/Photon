package models

import (
	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/ethereum/go-ethereum/common"
)

//NewTokenCb notify when new token registered
//return true to remove this callback, all the callback should never block.
type NewTokenCb func(token common.Address) (remove bool)

//ChannelCb notify when channel status changed
//return true to remove this callback, all the callback should never block.
type ChannelCb func(c *channel.Serialization) (remove bool)

//RegisterNewTokenCallback register a new token callback
func (model *ModelDB) RegisterNewTokenCallback(f NewTokenCb) {
	model.mlock.Lock()
	model.newTokenCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterNewChannellCallback register a new channel callback
func (model *ModelDB) RegisterNewChannellCallback(f ChannelCb) {
	model.mlock.Lock()
	model.newChannelCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterChannelDepositCallback register channel deposit callback
func (model *ModelDB) RegisterChannelDepositCallback(f ChannelCb) {
	model.mlock.Lock()
	model.channelDepositCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterChannelStateCallback notify when channel closed or settled
func (model *ModelDB) RegisterChannelStateCallback(f ChannelCb) {
	model.mlock.Lock()
	model.channelStateCallbacks[&f] = true
	model.mlock.Unlock()
}

/*
do we need remove a callback?
*/
func (model *ModelDB) unRegisterNewTokenCallback(f NewTokenCb) {
	model.mlock.Lock()
	delete(model.newTokenCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) unRegisterNewChannellCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.newChannelCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) unRegisterChannelDepositCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.channelDepositCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) unRegisterChannelStateCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.channelStateCallbacks, &f)
	model.mlock.Unlock()
}
