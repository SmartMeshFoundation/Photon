package stormdb

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/log"
)

const bucketBlockNumber = "bucketBlockNumber"
const keyBlockNumber = "blocknumber"
const keyBlockTime = "blockTime"

//GetLatestBlockNumber lastest block number
func (model *StormDB) GetLatestBlockNumber() int64 {
	var number int64
	err := model.db.Get(bucketBlockNumber, keyBlockNumber, &number)
	if err != nil {
		log.Error(fmt.Sprintf("models GetLatestBlockNumber err=%s", err))
	}
	return number
}

//SaveLatestBlockNumber block numer has been processed
func (model *StormDB) SaveLatestBlockNumber(blockNumber int64) {
	err := model.db.Set(bucketBlockNumber, keyBlockNumber, blockNumber)
	if err != nil {
		log.Error(fmt.Sprintf("models SaveLatestBlockNumber err=%s", err))
	}
	err = model.db.Set(bucketBlockNumber, keyBlockTime, time.Now())
	if err != nil {
		log.Error(fmt.Sprintf("models SaveLatestBlockTime err=%s", err))
	}
}

//GetLastBlockNumberTime return when last block received
func (model *StormDB) GetLastBlockNumberTime() time.Time {
	var t time.Time
	err := model.db.Get(bucketBlockNumber, keyBlockTime, &t)
	if err != nil {
		log.Error(fmt.Sprintf("GetLastBlockNumberTime err %s", err))
	}
	return t
}
