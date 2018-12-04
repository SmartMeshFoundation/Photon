package gkvdb

import (
	"encoding/hex"
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
IsThisLockRemoved return true when  a expired hashlock has been removed from channel status.
*/
func (dao *GkvDB) IsThisLockRemoved(channel common.Hash, sender common.Address, lockHash common.Hash) bool {
	var result bool
	key := utils.Sha3(channel[:], lockHash[:], sender[:])
	err := dao.getKeyValueToBucket(models.BucketExpiredHashlock, key.Bytes(), &result)
	if err != nil {
		return false
	}
	if result != true {
		panic("expiredHashlock cannot be set to false")
	}
	return result
}

/*
RemoveLock remember this lock has been removed from channel status.
*/
func (dao *GkvDB) RemoveLock(channel common.Hash, sender common.Address, lockHash common.Hash) {
	key := utils.Sha3(channel[:], lockHash[:], sender[:])
	err := dao.saveKeyValueToBucket(models.BucketExpiredHashlock, key.Bytes(), true)
	if err != nil {
		log.Error(fmt.Sprintf("UnlockThisLock write %s to db err %s", hex.EncodeToString(key.Bytes()), err))
	}
}
