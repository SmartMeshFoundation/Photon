package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func (model *ModelDB) StartTx() (tx storm.Node) {
	var err error
	tx, err = model.db.Begin(true)
	if err != nil {
		panic(fmt.Sprintf("start transaction error %s", err))
	}
	return
}
func (model *ModelDB) AddStateManager(mgr *transfer.StateManager) error {
	err := model.db.Save(mgr)
	if err != nil {
		log.Error(fmt.Sprintf(" AddStateManager err=%s", err))
	}
	return err
}
func (model *ModelDB) UpdateStateManaer(mgr *transfer.StateManager, tx storm.Node) error {
	//log.Trace(fmt.Sprintf("UpdateStateManaer %s\n", utils.StringInterface(mgr, 7)))
	err := tx.Save(mgr)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStateManaer err=%s", err))
	}
	return err
}
func (model *ModelDB) GetAllStateManager() []*transfer.StateManager {
	var mgrs []*transfer.StateManager
	//err := model.db.Find("IsFinished", false, &mgrs)
	err := model.db.All(&mgrs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllUnfinishedStateManager err %s", err))
	}
	return mgrs
}
func (model *ModelDB) GetAck(echohash common.Hash) []byte {
	var data []byte
	log.Trace(fmt.Sprintf("quer ack %s from db", utils.HPex(echohash)))
	err := model.db.Get("ack", echohash.String(), &data)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAck err %s", err))
	}
	return data
}

func (model *ModelDB) SaveAck(echohash common.Hash, ack []byte, tx storm.Node) {
	log.Trace(fmt.Sprintf("save ack %s to db", utils.HPex(echohash)))
	tx.Set("ack", echohash.String(), ack)
}
