package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/asdine/storm"
)

// StormTx proxy of storm.Node
type StormTx struct {
	tx storm.Node
}

// Set proxy of storm.Node
func (stx *StormTx) Set(table string, key interface{}, value interface{}) error {
	err := stx.tx.Set(table, key, value)
	return err
}

// Save proxy of storm.Node
func (stx *StormTx) Save(v models.KeyGetter) error {
	err := stx.tx.Save(v)
	return err
}

// Commit proxy of storm.Node
func (stx *StormTx) Commit() error {
	return stx.tx.Commit()
}

// Rollback proxy of storm.Node
func (stx *StormTx) Rollback() error {
	return stx.tx.Rollback()
}

//StartTx start a new tx of db
func (model *StormDB) StartTx() (tx models.TX) {
	stx, err := model.db.Begin(true)
	if err != nil {
		panic(fmt.Sprintf("start transaction error %s", err))
	}
	return &StormTx{
		tx: stx,
	}
}
