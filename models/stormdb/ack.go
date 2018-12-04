package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

//GetAck get message related ack message
func (model *StormDB) GetAck(echoHash common.Hash) []byte {
	var data []byte
	err := model.db.Get(models.BucketAck, echoHash[:], &data)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAck err %s", err))
	}
	log.Trace(fmt.Sprintf("get ack %s from db,result=%d", utils.HPex(echoHash), len(data)))
	return data
}

//SaveAck save a new ack to db
func (model *StormDB) SaveAck(echoHash common.Hash, ack []byte, tx models.TX) {
	log.Trace(fmt.Sprintf("save ack %s to db", utils.HPex(echoHash)))
	err := tx.Set(models.BucketAck, echoHash[:], ack)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//SaveAckNoTx save a ack to db
func (model *StormDB) SaveAckNoTx(echoHash common.Hash, ack []byte) {
	err := model.db.Set(models.BucketAck, echoHash[:], ack)
	if err != nil {
		log.Error(fmt.Sprintf("save ack to db err %s", err))
	}
}
