package stormdb

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
)

// SaveFeeChargeRecord :
func (model *StormDB) SaveFeeChargeRecord(r *models.FeeChargeRecord) (err error) {
	rs := r.ToSerialized()
	if rs.Key == nil || common.BytesToHash(rs.Key) == utils.EmptyHash {
		key := utils.NewRandomHash()
		rs.Key = key[:]
	}
	if rs.Timestamp <= 0 {
		rs.Timestamp = time.Now().Unix()
	}
	err = model.db.Save(rs)
	if err != nil {
		err = fmt.Errorf("SaveFeeChargeRecord err %s", err)
		err = models.GeneratDBError(err)
		return
	}
	log.Trace(fmt.Sprintf("charge for transfer:%s", r.ToString()))
	return
}

// GetAllFeeChargeRecord :
func (model *StormDB) GetAllFeeChargeRecord(tokenAddress common.Address, fromTime, toTime int64) (records []*models.FeeChargeRecord, err error) {
	var selectList []q.Matcher
	if tokenAddress != utils.EmptyAddress {
		selectList = append(selectList, q.Eq("TokenAddress", tokenAddress[:]))
	}
	if fromTime > 0 {
		selectList = append(selectList, q.Gte("Timestamp", fromTime))
	}
	if toTime > 0 {
		selectList = append(selectList, q.Lt("Timestamp", toTime))
	}
	var rs []*models.FeeChargerRecordSerialization
	if len(selectList) == 0 {
		err = model.db.All(&rs)
	} else {
		q := model.db.Select(selectList...)
		err = q.Find(&rs)
	}
	if err == storm.ErrNotFound {
		err = nil
	}
	for _, r := range rs {
		records = append(records, r.ToFeeChargeRecord())
	}
	return
}

// GetFeeChargeRecordByLockSecretHash :
func (model *StormDB) GetFeeChargeRecordByLockSecretHash(lockSecretHash common.Hash) (records []*models.FeeChargeRecord, err error) {
	var rs []*models.FeeChargerRecordSerialization
	err = model.db.Find("LockSecretHash", lockSecretHash[:], &rs)
	if err != nil {
		err = fmt.Errorf("GetAllFeeChargeRecordByLockSecretHash err %s", err)
		err = models.GeneratDBError(err)
		return
	}
	for _, r := range rs {
		records = append(records, r.ToFeeChargeRecord())
	}
	return
}
