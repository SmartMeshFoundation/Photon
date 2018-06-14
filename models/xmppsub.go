package models

import "github.com/ethereum/go-ethereum/common"

const bucketXMPP = "bucketxmpp"

//XMPPMarkAddrSubed mark `addr` subscribed
func (model *ModelDB) XMPPMarkAddrSubed(addr common.Address) {
	model.db.Set(bucketXMPP, addr, true)
}

//XMPPIsAddrSubed return true when `addr` already subscirbed
func (model *ModelDB) XMPPIsAddrSubed(addr common.Address) bool {
	var r bool
	model.db.Get(bucketXMPP, addr, &r)
	return r
}

//XMPPUnMarkAddr mark `addr` has been unsubscribed
func (model *ModelDB) XMPPUnMarkAddr(addr common.Address) {
	model.db.Set(bucketXMPP, addr, false)
}
