package stormdb

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
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
func (model *StormDB) GetAllFeeChargeRecord() (records []*models.FeeChargeRecord, err error) {
	var rs []*models.FeeChargerRecordSerialization
	err = model.db.All(&rs)
	if err != nil {
		err = fmt.Errorf("GetAllFeeChargeRecord err %s", err)
		err = models.GeneratDBError(err)
		return
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
