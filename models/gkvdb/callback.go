package gkvdb

import "github.com/SmartMeshFoundation/Photon/models/cb"

// RegisterNewTokenCallback register a new token callback
func (dao *GkvDB) RegisterNewTokenCallback(f cb.NewTokenCb) {
	dao.mlock.Lock()
	dao.newTokenCallbacks[&f] = true
	dao.mlock.Unlock()
}

// RegisterNewChannelCallback register a new channel callback
func (dao *GkvDB) RegisterNewChannelCallback(f cb.ChannelCb) {
	dao.mlock.Lock()
	dao.newChannelCallbacks[&f] = true
	dao.mlock.Unlock()
}

//RegisterChannelDepositCallback register channel deposit callback
func (dao *GkvDB) RegisterChannelDepositCallback(f cb.ChannelCb) {
	dao.mlock.Lock()
	dao.channelDepositCallbacks[&f] = true
	dao.mlock.Unlock()
}

//RegisterChannelStateCallback notify when channel closed
func (dao *GkvDB) RegisterChannelStateCallback(f cb.ChannelCb) {
	dao.mlock.Lock()
	dao.channelStateCallbacks[&f] = true
	dao.mlock.Unlock()
}

//RegisterChannelSettleCallback notify when channel settled
func (dao *GkvDB) RegisterChannelSettleCallback(f cb.ChannelCb) {
	dao.mlock.Lock()
	dao.channelSettledCallbacks[&f] = true
	dao.mlock.Unlock()
}

/*
do we need remove a callback?
*/
func (dao *GkvDB) unRegisterNewTokenCallback(f cb.NewTokenCb) {
	dao.mlock.Lock()
	delete(dao.newTokenCallbacks, &f)
	dao.mlock.Unlock()
}
func (dao *GkvDB) unRegisterNewChannelCallback(f cb.ChannelCb) {
	dao.mlock.Lock()
	delete(dao.newChannelCallbacks, &f)
	dao.mlock.Unlock()
}
func (dao *GkvDB) unRegisterChannelDepositCallback(f cb.ChannelCb) {
	dao.mlock.Lock()
	delete(dao.channelDepositCallbacks, &f)
	dao.mlock.Unlock()
}
func (dao *GkvDB) unRegisterChannelStateCallback(f cb.ChannelCb) {
	dao.mlock.Lock()
	delete(dao.channelStateCallbacks, &f)
	dao.mlock.Unlock()
}
