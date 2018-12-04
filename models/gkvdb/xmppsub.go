package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/ethereum/go-ethereum/common"
)

//XMPPMarkAddrSubed mark `addr` subscribed
func (dao *GkvDB) XMPPMarkAddrSubed(addr common.Address) {
	err := dao.saveKeyValueToBucket(models.BucketXMPP, addr[:], true)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//XMPPIsAddrSubed return true when `addr` already subscirbed
func (dao *GkvDB) XMPPIsAddrSubed(addr common.Address) bool {
	var r bool
	err := dao.getKeyValueToBucket(models.BucketXMPP, addr[:], &r)
	if err != nil {
		log.Trace(fmt.Sprintf("db err %s", err))
	}
	return r
}

//XMPPUnMarkAddr mark `addr` has been unsubscribed
func (dao *GkvDB) XMPPUnMarkAddr(addr common.Address) {
	err := dao.saveKeyValueToBucket(models.BucketXMPP, addr[:], false)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}
