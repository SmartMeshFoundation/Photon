package gkvdb

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
)

//GetLatestBlockNumber lastest block number
func (dao *GkvDB) GetLatestBlockNumber() int64 {
	var number int64
	err := dao.getKeyValueToBucket(models.BucketBlockNumber, models.KeyBlockNumber, &number)
	if err != nil {
		log.Error(fmt.Sprintf("models GetLatestBlockNumber err=%s", err))
	}
	return number
}

//SaveLatestBlockNumber block numer has been processed
func (dao *GkvDB) SaveLatestBlockNumber(blockNumber int64) {
	err := dao.saveKeyValueToBucket(models.BucketBlockNumber, models.KeyBlockNumber, blockNumber)
	if err != nil {
		log.Error(fmt.Sprintf("models SaveLatestBlockNumber err=%s", err))
	}
	err = dao.saveKeyValueToBucket(models.BucketBlockNumber, models.KeyBlockNumberTime, time.Now())
	if err != nil {
		log.Error(fmt.Sprintf("models SaveLatestBlockTime err=%s", err))
	}
}

//GetLastBlockNumberTime return when last block received
func (dao *GkvDB) GetLastBlockNumberTime() time.Time {
	var t time.Time
	err := dao.getKeyValueToBucket(models.BucketBlockNumber, models.KeyBlockNumberTime, &t)
	if err != nil {
		log.Error(fmt.Sprintf("GetLastBlockNumberTime err %s", err))
	}
	return t
}
