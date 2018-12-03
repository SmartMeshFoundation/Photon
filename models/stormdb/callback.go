package stormdb

import "github.com/SmartMeshFoundation/Photon/models/cb"

// RegisterNewTokenCallback register a new token callback
func (model *StormDB) RegisterNewTokenCallback(f cb.NewTokenCb) {
	model.mlock.Lock()
	model.newTokenCallbacks[&f] = true
	model.mlock.Unlock()
}

// RegisterNewChannelCallback register a new channel callback
func (model *StormDB) RegisterNewChannelCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	model.newChannelCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterChannelDepositCallback register channel deposit callback
func (model *StormDB) RegisterChannelDepositCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	model.channelDepositCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterChannelStateCallback notify when channel closed
func (model *StormDB) RegisterChannelStateCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	model.channelStateCallbacks[&f] = true
	model.mlock.Unlock()
}

//RegisterChannelSettleCallback notify when channel settled
func (model *StormDB) RegisterChannelSettleCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	model.channelSettledCallbacks[&f] = true
	model.mlock.Unlock()
}

/*
do we need remove a callback?
*/
func (model *StormDB) unRegisterNewTokenCallback(f cb.NewTokenCb) {
	model.mlock.Lock()
	delete(model.newTokenCallbacks, &f)
	model.mlock.Unlock()
}
func (model *StormDB) unRegisterNewChannelCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	delete(model.newChannelCallbacks, &f)
	model.mlock.Unlock()
}
func (model *StormDB) unRegisterChannelDepositCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	delete(model.channelDepositCallbacks, &f)
	model.mlock.Unlock()
}
func (model *StormDB) unRegisterChannelStateCallback(f cb.ChannelCb) {
	model.mlock.Lock()
	delete(model.channelStateCallbacks, &f)
	model.mlock.Unlock()
}
