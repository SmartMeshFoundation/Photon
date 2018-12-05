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
	return gtx.tx.SetTo(gobEncode(key), gobEncode(value), table)
}

// Save :
func (gtx *GkvTX) Save(v models.KeyGetter) error {
	panic("should not use this")
}

// Commit :
func (gtx *GkvTX) Commit() error {
	return gtx.tx.Commit(true)
}

// Rollback :
func (gtx *GkvTX) Rollback() error {
	gtx.tx.Rollback()
	return nil
}

//StartTx start a new tx of db
func (dao *GkvDB) StartTx() (tx models.TX) {
	gtx := dao.db.Begin()
	return &GkvTX{
		tx: gtx,
	}
}
