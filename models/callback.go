package models

import (
	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/ethereum/go-ethereum/common"
)

// new token, return true to remove this callback, all the callback should never block.
type NewTokenCb func(token common.Address) (remove bool)
type ChannelCb func(c *channel.ChannelSerialization) (remove bool)

func (model *ModelDB) RegisterNewTokenCallback(f NewTokenCb) {
	model.mlock.Lock()
	model.newTokenCallbacks[&f] = true
	model.mlock.Unlock()
}
func (model *ModelDB) RegisterNewChannellCallback(f ChannelCb) {
	model.mlock.Lock()
	model.newChannelCallbacks[&f] = true
	model.mlock.Unlock()
}
func (model *ModelDB) RegisterChannelDepositCallback(f ChannelCb) {
	model.mlock.Lock()
	model.channelDepositCallbacks[&f] = true
	model.mlock.Unlock()
}
func (model *ModelDB) RegisterChannelStateCallback(f ChannelCb) {
	model.mlock.Lock()
	model.channelStateCallbacks[&f] = true
	model.mlock.Unlock()
}

func (model *ModelDB) UnRegisterNewTokenCallback(f NewTokenCb) {
	model.mlock.Lock()
	delete(model.newTokenCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) UnRegisterNewChannellCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.newChannelCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) UnRegisterChannelDepositCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.channelDepositCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) UnRegisterChannelStateCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.channelStateCallbacks, &f)
	model.mlock.Unlock()
}
