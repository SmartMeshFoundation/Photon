package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//GetAck get message related ack message
func (dao *GkvDB) GetAck(echoHash common.Hash) []byte {
	var data []byte
	err := dao.getKeyValueToBucket(models.BucketAck, echoHash[:], &data)
	if err != nil && err != ErrorNotFound {
		panic(fmt.Sprintf("GetAck err %s", err))
	}
	log.Trace(fmt.Sprintf("get ack %s from db,result=%d", utils.HPex(echoHash), len(data)))
	return data
}

//SaveAck save a new ack to db
func (dao *GkvDB) SaveAck(echoHash common.Hash, ack []byte, tx models.TX) {
	log.Trace(fmt.Sprintf("save ack %s to db", utils.HPex(echoHash)))
	err := tx.Set(models.BucketAck, echoHash[:], ack)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//SaveAckNoTx save a ack to db
func (dao *GkvDB) SaveAckNoTx(echoHash common.Hash, ack []byte) {
	err := dao.saveKeyValueToBucket(models.BucketAck, echoHash[:], ack)
	if err != nil {
		log.Error(fmt.Sprintf("save ack to db err %s", err))
	}
}
