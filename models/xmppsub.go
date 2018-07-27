package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/ethereum/go-ethereum/common"
)

const bucketXMPP = "bucketxmpp"

//XMPPMarkAddrSubed mark `addr` subscribed
func (model *ModelDB) XMPPMarkAddrSubed(addr common.Address) {
	err := model.db.Set(bucketXMPP, addr[:], true)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//XMPPIsAddrSubed return true when `addr` already subscirbed
func (model *ModelDB) XMPPIsAddrSubed(addr common.Address) bool {
	var r bool
	err := model.db.Get(bucketXMPP, addr[:], &r)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	return r
}

//XMPPUnMarkAddr mark `addr` has been unsubscribed
func (model *ModelDB) XMPPUnMarkAddr(addr common.Address) {
	err := model.db.Set(bucketXMPP, addr[:], false)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}
