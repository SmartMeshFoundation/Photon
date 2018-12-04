package gkvdb

import (
	"errors"
	"fmt"

	"sync"

	"encoding/gob"

	"bytes"

	"gitee.com/johng/gkvdb/gkvdb"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/ethereum/go-ethereum/common"
)

// ErrorNotFound :
var ErrorNotFound = errors.New("not found")

//GkvDB is thread safe
type GkvDB struct {
	db                      *gkvdb.DB
	lock                    sync.Mutex
	newTokenCallbacks       map[*cb.NewTokenCb]bool
	newChannelCallbacks     map[*cb.ChannelCb]bool
	channelDepositCallbacks map[*cb.ChannelCb]bool
	channelStateCallbacks   map[*cb.ChannelCb]bool
	channelSettledCallbacks map[*cb.ChannelCb]bool
	mlock                   sync.Mutex
	Name                    string
}

func newGkvDB() (db *GkvDB) {
	return &GkvDB{
		newTokenCallbacks:       make(map[*cb.NewTokenCb]bool),
		newChannelCallbacks:     make(map[*cb.ChannelCb]bool),
		channelDepositCallbacks: make(map[*cb.ChannelCb]bool),
		channelStateCallbacks:   make(map[*cb.ChannelCb]bool),
		channelSettledCallbacks: make(map[*cb.ChannelCb]bool),
	}
}
func gobEncode(d interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func gobDecode(buf []byte, to interface{}) {
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(to)
	if err != nil {
		panic(err)
	}
}

func (dao *GkvDB) saveKeyValueToBucket(bucket string, key, value interface{}) error {
	tb, err := dao.db.Table(bucket)
	if err != nil {
		return err
	}
	err = tb.Set(gobEncode(key), gobEncode(value))
	if err != nil {
		return err
	}
	return nil
}

func (dao *GkvDB) getKeyValueToBucket(bucket string, key, to interface{}) error {
	tb, err := dao.db.Table(bucket)
	if err != nil {
		return err
	}
	buf := tb.Get(gobEncode(key))
	if buf == nil || len(buf) == 0 {
		return ErrorNotFound
	}
	gobDecode(buf, to)
	return nil
}

func (dao *GkvDB) removeKeyValueFromBucket(bucket string, key interface{}) error {
	tb, err := dao.db.Table(bucket)
	if err != nil {
		return err
	}
	return tb.Remove(gobEncode(key))
}

//OpenDb open or create a bolt db at dbPath
func OpenDb(dbPath string) (dao *GkvDB, err error) {
	log.Trace(fmt.Sprintf("dbpath=%s", dbPath))
	dao = newGkvDB()
	needCreateDb := !common.FileExist(dbPath)
	var ver int
	dao.db, err = gkvdb.New(dbPath)
	if err != nil {
		err = fmt.Errorf("cannot create or open db:%s,makesure you have write permission err:%v", dbPath, err)
		log.Crit(err.Error())
		return
	}
	dao.Name = dbPath
	if needCreateDb {
		err = dao.saveKeyValueToBucket(models.BucketMeta, models.KeyVersion, models.DbVersion)
		if err != nil {
			log.Error("save version err %s", err)
			return
		}
		err = dao.saveKeyValueToBucket(models.BucketToken, models.KeyToken, make(models.AddressMap))
		if err != nil {
			log.Error("init token table error %s", err)
			return
		}
		dao.initDb()
		dao.MarkDbOpenedStatus()
	} else {
		err = dao.getKeyValueToBucket(models.BucketMeta, models.KeyVersion, &ver)
		if err != nil {
			log.Error("get version error %s", err)
			return
		}
		if ver != models.DbVersion {
			log.Error("db version not match")
			return
		}
		var closeFlag bool
		err = dao.getKeyValueToBucket(models.BucketMeta, models.KeyCloseFlag, &closeFlag)
		if err != nil {
			log.Error("db meta data error %s", err)
			return
		}
		if closeFlag != true {
			log.Error("database not closed  last..., try to restore?")
		}
	}
	return
}

func init() {
	gob.Register(&GkvDB{}) //cannot save and restore by gob,only avoid noise by gob
}

func (dao *GkvDB) initDb() {
	err := dao.saveKeyValueToBucket(models.BucketBlockNumber, models.KeyBlockNumber, 0)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}
