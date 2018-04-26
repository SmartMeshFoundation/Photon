package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
)

const bucketBlockNumber = "bucketBlockNumber"
const keyBlockNumber = "blocknumber"

func (model *ModelDB) GetLatestBlockNumber() int64 {
	var number int64
	err := model.db.Get(bucketBlockNumber, keyBlockNumber, &number)
	if err != nil {
		log.Error(fmt.Sprintf("models GetLatestBlockNumber err=%s", err))
	}
	return number
}

func (model *ModelDB) SaveLatestBlockNumber(blockNumber int64) {
	err := model.db.Set(bucketBlockNumber, keyBlockNumber, blockNumber)
	if err != nil {
		log.Error(fmt.Sprintf("models SaveLatestBlockNumber err=%s", err))
	}
}
