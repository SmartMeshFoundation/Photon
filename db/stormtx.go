package db

import (
	"github.com/asdine/storm"
)

// StormTx :
type StormTx struct {
	tx storm.Node
}

// newStormTx :
func newStormTx(db *storm.DB, writable bool) (*StormTx, error) {
	tx, err := db.Begin(writable)
	if err != nil {
		return nil, err
	}
	return &StormTx{
		tx: tx,
	}, nil
}

// Set :
func (stx *StormTx) Set(table string, key interface{}, value interface{}) error {
	return stx.tx.Set(table, key, value)
}

// Remove :
func (stx *StormTx) Remove(table string, key interface{}) error {
	return stx.tx.Delete(table, key)
}

// Save :
func (stx *StormTx) Save(v KeyGetter) error {
	return stx.tx.Save(v)
}

// Get :
func (stx *StormTx) Get(table string, key interface{}, to interface{}) error {
	return stx.tx.Get(table, key, to)
}

// All :
func (stx *StormTx) All(table string, to interface{}) error {
	return stx.tx.All(to)
}

// Find :
func (stx *StormTx) Find(table string, fieldName string, value interface{}, to interface{}) error {
	return stx.tx.Find(fieldName, value, to)
}

// Range :
func (stx *StormTx) Range(table string, fieldName string, min, max, to interface{}) error {
	return stx.tx.Range(fieldName, min, max, to)
}

// Commit :
func (stx *StormTx) Commit() error {
	return stx.tx.Commit()
}

// Rollback :
func (stx *StormTx) Rollback() error {
	return stx.tx.Rollback()
}
