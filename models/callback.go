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
	model.NewTokenCallbacks[&f] = true
	model.mlock.Unlock()
}
func (model *ModelDB) RegisterNewChannellCallback(f ChannelCb) {
	model.mlock.Lock()
	model.NewChannelCallbacks[&f] = true
	model.mlock.Unlock()
}
func (model *ModelDB) RegisterChannelDepositCallback(f ChannelCb) {
	model.mlock.Lock()
	model.ChannelDepositCallbacks[&f] = true
	model.mlock.Unlock()
}
func (model *ModelDB) RegisterChannelStateCallback(f ChannelCb) {
	model.mlock.Lock()
	model.ChannelStateCallbacks[&f] = true
	model.mlock.Unlock()
}

func (model *ModelDB) UnRegisterNewTokenCallback(f NewTokenCb) {
	model.mlock.Lock()
	delete(model.NewTokenCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) UnRegisterNewChannellCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.NewChannelCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) UnRegisterChannelDepositCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.ChannelDepositCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) UnRegisterChannelStateCallback(f ChannelCb) {
	model.mlock.Lock()
	delete(model.ChannelStateCallbacks, &f)
	model.mlock.Unlock()
}
