package models

import "github.com/SmartMeshFoundation/SmartRaiden/models/cb"

//RegisterNewTokenCallback register a new token callback
func (model *ModelDB) RegisterNewTokenCallback(f cb.NewTokenCb) {
	model.mlock.Lock()
	model.newTokenCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterNewChannellCallback register a new channel callback
func (model *ModelDB) RegisterNewChannellCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	model.newChannelCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterChannelDepositCallback register channel deposit callback
func (model *ModelDB) RegisterChannelDepositCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	model.channelDepositCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterChannelStateCallback notify when channel closed or settled
func (model *ModelDB) RegisterChannelStateCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	model.channelStateCallbacks[&f] = true
	model.mlock.Unlock()
}

/*
do we need remove a callback?
*/
func (model *ModelDB) unRegisterNewTokenCallback(f cb.NewTokenCb) {
	model.mlock.Lock()
	delete(model.newTokenCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) unRegisterNewChannellCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	delete(model.newChannelCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) unRegisterChannelDepositCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	delete(model.channelDepositCallbacks, &f)
	model.mlock.Unlock()
}
func (model *ModelDB) unRegisterChannelStateCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	delete(model.channelStateCallbacks, &f)
	model.mlock.Unlock()
}
