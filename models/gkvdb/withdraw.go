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
IsThisLockHasUnlocked return ture when  lockhash has unlocked on channel?
*/
func (dao *GkvDB) IsThisLockHasUnlocked(channel common.Hash, lockHash common.Hash) bool {
	var result bool
	key := utils.Sha3(channel[:], lockHash[:])
	err := dao.getKeyValueToBucket(models.BucketWithDraw, key.Bytes(), &result)
	if err != nil {
		return false
	}
	if result != true {
		panic("withdraw cannot be set to false")
	}
	return result
}

/*
UnlockThisLock marks that I have withdrawed this secret on channel.
*/
func (dao *GkvDB) UnlockThisLock(channel common.Hash, lockHash common.Hash) {
	key := utils.Sha3(channel[:], lockHash[:])
	err := dao.saveKeyValueToBucket(models.BucketWithDraw, key.Bytes(), true)
	if err != nil {
		log.Error(fmt.Sprintf("UnlockThisLock write %s to db err %s", hex.EncodeToString(key.Bytes()), err))
	}
}
