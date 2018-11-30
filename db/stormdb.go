package db

import (
	"fmt"
	"os"
	"time"

	"github.com/asdine/storm"
	gobcodec "github.com/asdine/storm/codec/gob"
	"github.com/coreos/bbolt"
	"github.com/nkbai/log"
)

// StormDB :
type StormDB struct {
	db *storm.DB
}

// NewStormDB :
func NewStormDB(dbPath string) (*StormDB, error) {
	db, err := storm.Open(dbPath, storm.BoltOptions(os.ModePerm, &bolt.Options{Timeout: 1 * time.Second}), storm.Codec(gobcodec.Codec))
	if err != nil {
		err = fmt.Errorf("cannot create or open db:%s,makesure you have write permission err:%v", dbPath, err)
		log.Crit(err.Error())
		return nil, err
	}
	return &StormDB{
		db: db,
	}, nil
}

/*
	impl db
*/

// Set :
func (sdb *StormDB) Set(table string, key interface{}, value interface{}) error {
	return sdb.db.Set(table, key, value)
}

// Remove :
func (sdb *StormDB) Remove(table string, key interface{}) error {
	return sdb.db.Delete(table, key)
}

// Save :
func (sdb *StormDB) Save(v KeyGetter) error {
	return sdb.db.Save(v)
}

// Get :
func (sdb *StormDB) Get(table string, key interface{}, to interface{}) error {
	return sdb.db.Get(table, key, to)
}

// All :
func (sdb *StormDB) All(table string, to interface{}) error {
	return sdb.db.All(to)
}

// Find :
func (sdb *StormDB) Find(table string, fieldName string, value interface{}, to interface{}) error {
	return sdb.db.Find(fieldName, value, to)
}

// Range :
func (sdb *StormDB) Range(table string, fieldName string, min, max, to interface{}) error {
	return sdb.db.Range(fieldName, min, max, to)
}

// Begin :
func (sdb *StormDB) Begin(writable bool) (TX, error) {
	return newStormTx(sdb.db, writable)
}

// Close :
func (sdb *StormDB) Close() error {
	return sdb.db.Close()
}
