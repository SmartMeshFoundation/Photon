package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewDeliveredChainEvent save one
func (dao *GkvDB) NewDeliveredChainEvent(id models.ChainEventID, blockNumber uint64) {
	e := &models.ChainEventRecord{
		ID:          id,
		BlockNumber: blockNumber,
		Status:      models.ChainEventStatusDelivered,
	}
	err := dao.saveKeyValueToBucket(models.BucketChainEventRecord, e.ID, e)
	if err != nil {
		log.Error(fmt.Sprintf("models NewDeliveredChainEvent err=%s", err))
	}
}

// CheckChainEventDelivered check one ChainEvent is delivered or not
func (dao *GkvDB) CheckChainEventDelivered(id models.ChainEventID) (blockNumber uint64, delivered bool) {
	e := &models.ChainEventRecord{}
	err := dao.getKeyValueToBucket(models.BucketChainEventRecord, e.ID, e)
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
func (dao *GkvDB) ClearOldChainEventRecord(blockNumber uint64) {
	tb, err := dao.db.Table(models.BucketChainEventRecord)
	if err != nil {
		err = models.GeneratDBError(err)
		return
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var r models.ChainEventRecord
		gobDecode(v, &r)
		if r.BlockNumber <= blockNumber {
			err2 := dao.removeKeyValueFromBucket(models.BucketChainEventRecord, r.ID)
			if err2 != nil {
				log.Error(fmt.Sprintf("models ClearOldChainEventRecord err=%s", err))
			}
		}
	}
}

// MakeChainEventID :
func (dao *GkvDB) MakeChainEventID(l *types.Log) models.ChainEventID {
	var t [25]byte
	copy(t[:], l.TxHash[:])
	t[24] = byte(l.Index)
	return models.ChainEventID(common.Bytes2Hex(t[:]))
}
