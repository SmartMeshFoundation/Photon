package gkvdb

import (
	"fmt"

	"time"

	"gitee.com/johng/gkvdb/gkvdb"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// SaveFeeChargeRecord :
func (dao *GkvDB) SaveFeeChargeRecord(r *models.FeeChargeRecord) (err error) {
	if r.Key == utils.EmptyHash {
		r.Key = utils.NewRandomHash()
	}
	if r.Timestamp <= 0 {
		r.Timestamp = time.Now().Unix()
	}
	err = dao.saveKeyValueToBucket(models.BucketFeeChargeRecord, r.Key, r)
	if err != nil {
		err = fmt.Errorf("SaveFeeChargeRecord err %s", err)
		return
	}
	log.Trace(fmt.Sprintf("charge for transfer:%s", r.ToString()))
	return
}

// GetAllFeeChargeRecord :
func (dao *GkvDB) GetAllFeeChargeRecord() (records []*models.FeeChargeRecord, err error) {
	var tb *gkvdb.Table
	tb, err = dao.db.Table(models.BucketFeeChargeRecord)
	if err != nil {
		return
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var r models.FeeChargeRecord
		gobDecode(v, &r)
		records = append(records, &r)
	}
	return
}

// GetFeeChargeRecordByLockSecretHash :
func (dao *GkvDB) GetFeeChargeRecordByLockSecretHash(lockSecretHash common.Hash) (records []*models.FeeChargeRecord, err error) {
	var rs []*models.FeeChargeRecord
	rs, err = dao.GetAllFeeChargeRecord()
	if err != nil {
		return
	}
	for _, r := range rs {
		if r.LockSecretHash == lockSecretHash {
			records = append(records, r)
		}
	}
	return
}
