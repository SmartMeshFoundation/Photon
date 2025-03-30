package stormdb

import (
	"fmt"

	"sync"

	"time"

	"encoding/gob"

	"os"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/asdine/storm"
	gobcodec "github.com/asdine/storm/codec/gob"
	bolt "github.com/coreos/bbolt"
	"github.com/ethereum/go-ethereum/common"
)

//StormDB is thread safe
type StormDB struct {
	db                      *storm.DB
	lock                    sync.Mutex
	newTokenCallbacks       map[*cb.NewTokenCb]bool
	newChannelCallbacks     map[*cb.ChannelCb]bool
	channelDepositCallbacks map[*cb.ChannelCb]bool
	channelStateCallbacks   map[*cb.ChannelCb]bool
	channelSettledCallbacks map[*cb.ChannelCb]bool
	mlock                   sync.Mutex
	Name                    string
}

func newStormDB() (db *StormDB) {
	return &StormDB{
		newTokenCallbacks:       make(map[*cb.NewTokenCb]bool),
		newChannelCallbacks:     make(map[*cb.ChannelCb]bool),
		channelDepositCallbacks: make(map[*cb.ChannelCb]bool),
		channelStateCallbacks:   make(map[*cb.ChannelCb]bool),
		channelSettledCallbacks: make(map[*cb.ChannelCb]bool),
	}

}

//OpenDb open or create a bolt db at dbPath
func OpenDb(dbPath string) (model *StormDB, err error) {
	log.Trace(fmt.Sprintf("dbpath=%s", dbPath))
	model = newStormDB()
	needCreateDb := !common.FileExist(dbPath)
	var ver int
	model.db, err = storm.Open(dbPath, storm.BoltOptions(os.ModePerm, &bolt.Options{Timeout: 1 * time.Second}), storm.Codec(gobcodec.Codec))
	if err != nil {
		err = fmt.Errorf("cannot create or open db:%s,makesure you have write permission err:%v", dbPath, err)
		return
	}
	model.Name = dbPath
	if needCreateDb {
		err = model.db.Set(models.BucketMeta, models.KeyVersion, models.DbVersion)
		if err != nil {
			err = fmt.Errorf("unable to create db : %s", err.Error())
			return
		}
		err = model.db.Set(models.BucketToken, models.KeyToken, make(models.AddressMap))
		if err != nil {
			err = fmt.Errorf("unable to create db : %s", err.Error())
			return
		}
		model.initDb()
		model.MarkDbOpenedStatus()
	} else {
		err = model.db.Get(models.BucketMeta, models.KeyVersion, &ver)
		if err != nil {
			err = fmt.Errorf("wrong db file format : %s", err.Error())
			return
		}
		if ver != models.DbVersion {
			err = fmt.Errorf("db version not match : %s", err.Error())
			return
		}
		var closeFlag bool
		err = model.db.Get(models.BucketMeta, models.KeyCloseFlag, &closeFlag)
		if err != nil {
			err = fmt.Errorf("db meta data error : %s", err.Error())
			return
		}
		if closeFlag != true {
			log.Error("database not closed  last..., try to restore?")
		}
	}

	return
}

/*
MarkDbOpenedStatus First step   open the database
Second step detection for normal closure IsDbCrashedLastTime
Third step  recovers the data according to the second step
Fourth step mark the database for processing the data normally. MarkDbOpenedStatus
*/
func (model *StormDB) MarkDbOpenedStatus() {
	err := model.db.Set(models.BucketMeta, models.KeyCloseFlag, false)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//IsDbCrashedLastTime return true when quit but  db not closed
func (model *StormDB) IsDbCrashedLastTime() bool {
	var closeFlag bool
	err := model.db.Get(models.BucketMeta, models.KeyCloseFlag, &closeFlag)
	if err != nil {
		panic(fmt.Sprintf("db meta data error"))
	}
	return closeFlag != true
}

//CloseDB close db
func (model *StormDB) CloseDB() {
	model.lock.Lock()
	err := model.db.Set(models.BucketMeta, models.KeyCloseFlag, true)
	err = model.db.Close()
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	model.lock.Unlock()
}

//SaveContractStatus save registry address to db
func (model *StormDB) SaveContractStatus(contractStatus models.ContractStatus) {
	err := model.db.Set(models.BucketMeta, models.KeyRegistry, contractStatus)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//GetContractStatus returns registry address in db
func (model *StormDB) GetContractStatus() models.ContractStatus {
	var contractStatus models.ContractStatus
	err := model.db.Get(models.BucketMeta, models.KeyRegistry, &contractStatus)
	if err != nil && err != storm.ErrNotFound {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	return contractStatus
}

func init() {
	gob.Register(&StormDB{}) //cannot save and restore by gob,only avoid noise by gob
}

func (model *StormDB) initDb() {
	err := model.db.Init(&models.ReceivedTransfer{})
	err = model.db.Set(models.BucketBlockNumber, models.KeyBlockNumber, 0)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}
