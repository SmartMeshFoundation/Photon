package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewDeliveredChainEvent save one
func (model *StormDB) NewDeliveredChainEvent(id models.ChainEventID, blockNumber uint64) {
	e := &models.ChainEventRecord{
		ID:          id,
		BlockNumber: blockNumber,
		Status:      models.ChainEventStatusDelivered,
	}
	err := model.db.Save(e)
	if err != nil {
		log.Error(fmt.Sprintf("models NewDeliveredChainEvent err=%s", err))
	}
	log.Trace(fmt.Sprintf("NewDeliveredChainEvent id=%s blockNumber=%d", e.ID, e.BlockNumber))
}

// CheckChainEventDelivered check one ChainEvent is delivered or not
func (model *StormDB) CheckChainEventDelivered(id models.ChainEventID) (blockNumber uint64, delivered bool) {
	e := &models.ChainEventRecord{}
	err := model.db.One("ID", id, e)
	if err == storm.ErrNotFound {
		delivered = false
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("models CheckChainEventDelivered err=%s", err))
		delivered = false
		return
	}
	if e.Status != models.ChainEventStatusDelivered {
		delivered = false
		return
	}
	delivered = true
	blockNumber = e.BlockNumber
	return
}

// ClearOldChainEventRecord delete records which blockNumber <= blockNumber in param
func (model *StormDB) ClearOldChainEventRecord(blockNumber uint64) {
	var list []*models.ChainEventRecord
	err := model.db.Range("BlockNumber", 0, blockNumber, &list)
	if err == storm.ErrNotFound {
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("models ClearOldChainEventRecord err=%s", err))
		return
	}
	for _, r := range list {
		err2 := model.db.DeleteStruct(r)
		if err2 != nil {
			log.Error(fmt.Sprintf("models ClearOldChainEventRecord DeleteStruct id=%s blockNumber=%d status=%s err=%s", r.ID, r.BlockNumber, r.Status, err2.Error()))
		}
	}
	log.Trace(fmt.Sprintf("ClearOldChainEventRecord remove %d events witch blockNumber < %d", len(list), blockNumber))
}

// MakeChainEventID 根据log构造一个ChainEventID
func (model *StormDB) MakeChainEventID(l *types.Log) models.ChainEventID {
	var t [25]byte
	copy(t[:], l.TxHash[:])
	t[24] = byte(l.Index)
	return models.ChainEventID(common.Bytes2Hex(t[:]))
}
