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
	"github.com/coreos/bbolt"
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

var bucketMeta = "meta"

const dbVersion = 1

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
		log.Crit(err.Error())
		return
	}
	model.Name = dbPath
	if needCreateDb {
		err = model.db.Set(bucketMeta, "version", dbVersion)
		if err != nil {
			log.Crit(fmt.Sprintf("unable to create db "))
			return
		}
		err = model.db.Set(bucketToken, keyToken, make(models.AddressMap))
		if err != nil {
			log.Crit(fmt.Sprintf("unable to create db "))
			return
		}
		model.initDb()
		model.MarkDbOpenedStatus()
	} else {
		err = model.db.Get(bucketMeta, "version", &ver)
		if err != nil {
			log.Crit(fmt.Sprintf("wrong db file format "))
			return
		}
		if ver != dbVersion {
			log.Crit("db version not match")
		}
		var closeFlag bool
		err = model.db.Get(bucketMeta, "close", &closeFlag)
		if err != nil {
			log.Crit(fmt.Sprintf("db meta data error"))
		}
		if closeFlag != true {
			log.Error("database not closed  last..., try to restore?")
		}
	}

	return
}

//StartTx start a new tx of db
func (model *StormDB) StartTx() (tx models.TX) {
	var err error
	tx, err = model.db.Begin(true)
	if err != nil {
		panic(fmt.Sprintf("start transaction error %s", err))
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
	err := model.db.Set(bucketMeta, "close", false)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//IsDbCrashedLastTime return true when quit but  db not closed
func (model *StormDB) IsDbCrashedLastTime() bool {
	var closeFlag bool
	err := model.db.Get(bucketMeta, "close", &closeFlag)
	if err != nil {
		log.Crit(fmt.Sprintf("db meta data error"))
	}
	return closeFlag != true
}

//CloseDB close db
func (model *StormDB) CloseDB() {
	model.lock.Lock()
	err := model.db.Set(bucketMeta, "close", true)
	err = model.db.Close()
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	model.lock.Unlock()
}

//SaveRegistryAddress save registry address to db
func (model *StormDB) SaveRegistryAddress(registryAddress common.Address) {
	err := model.db.Set(bucketMeta, "registry", registryAddress)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//GetRegistryAddress returns registry address in db
func (model *StormDB) GetRegistryAddress() common.Address {
	var registry common.Address
	err := model.db.Get(bucketMeta, "registry", &registry)
	if err != nil && err != storm.ErrNotFound {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	return registry
}

//SaveSecretRegistryAddress save secret registry contract address to db
func (model *StormDB) SaveSecretRegistryAddress(secretRegistryAddress common.Address) {
	err := model.db.Set(bucketMeta, "secretregistry", secretRegistryAddress)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//GetSecretRegistryAddress return secret registry contract address
func (model *StormDB) GetSecretRegistryAddress() common.Address {
	var secretRegistry common.Address
	err := model.db.Get(bucketMeta, "secretregistry", &secretRegistry)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	return secretRegistry
}
func init() {
	gob.Register(&StormDB{}) //cannot save and restore by gob,only avoid noise by gob
}

func (model *StormDB) initDb() {
	err := model.db.Init(&models.SentTransfer{})
	err = model.db.Init(&models.ReceivedTransfer{})
	err = model.db.Set(bucketBlockNumber, keyBlockNumber, 0)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}
