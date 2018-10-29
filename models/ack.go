package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

const bucketAck = "ack"

//GetAck get message related ack message
func (model *ModelDB) GetAck(echohash common.Hash) []byte {
	var data []byte
	err := model.db.Get(bucketAck, echohash[:], &data)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAck err %s", err))
	}
	log.Trace(fmt.Sprintf("get ack %s from db,result=%d", utils.HPex(echohash), len(data)))
	return data
}

//SaveAck save a new ack to db
func (model *ModelDB) SaveAck(echohash common.Hash, ack []byte, tx storm.Node) {
	log.Trace(fmt.Sprintf("save ack %s to db", utils.HPex(echohash)))
	err := tx.Set(bucketAck, echohash[:], ack)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//SaveAckNoTx save a ack to db
func (model *ModelDB) SaveAckNoTx(echohash common.Hash, ack []byte) {
	err := model.db.Set(bucketAck, echohash[:], ack)
	if err != nil {
		log.Error(fmt.Sprintf("save ack to db err %s", err))
	}
}
