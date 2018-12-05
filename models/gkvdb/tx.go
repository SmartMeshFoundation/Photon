package gkvdb

import (
	"gitee.com/johng/gkvdb/gkvdb"
	"github.com/SmartMeshFoundation/Photon/models"
)

// GkvTX :
type GkvTX struct {
	tx *gkvdb.Transaction
}

// Set :
func (gtx *GkvTX) Set(table string, key interface{}, value interface{}) error {
	return gtx.tx.Set(gobEncode(key), gobEncode(value))
}

// Save :
func (gtx *GkvTX) Save(v models.KeyGetter) error {
	return gtx.tx.Set(gobEncode(v.GetKey()), gobEncode(v))
}

// Commit :
func (gtx *GkvTX) Commit() error {
	return gtx.tx.Commit()
}

// Rollback :
func (gtx *GkvTX) Rollback() error {
	gtx.tx.Rollback()
	return nil
}

//StartTx start a new tx of db
func (dao *GkvDB) StartTx(bucketName string) (tx models.TX) {
	var gtx *gkvdb.Transaction
	if bucketName == "" {
		gtx = dao.db.Begin()
	} else {
		gtx = dao.db.Begin(bucketName)
	}
	return &GkvTX{
		tx: gtx,
	}
}
