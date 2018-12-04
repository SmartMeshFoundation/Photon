package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/ethereum/go-ethereum/common"
)

/*
MarkDbOpenedStatus First step   open the database
Second step detection for normal closure IsDbCrashedLastTime
Third step  recovers the data according to the second step
Fourth step mark the database for processing the data normally. MarkDbOpenedStatus
*/
func (dao *GkvDB) MarkDbOpenedStatus() {
	err := dao.saveKeyValueToBucket(models.BucketMeta, models.KeyCloseFlag, false)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//IsDbCrashedLastTime return true when quit but  db not closed
func (dao *GkvDB) IsDbCrashedLastTime() bool {
	var closeFlag bool
	err := dao.getKeyValueToBucket(models.BucketMeta, models.KeyCloseFlag, &closeFlag)
	if err != nil {
		log.Crit(fmt.Sprintf("db meta data error"))
	}
	return closeFlag != true
}

//CloseDB close db
func (dao *GkvDB) CloseDB() {
	dao.lock.Lock()
	err := dao.saveKeyValueToBucket(models.BucketMeta, models.KeyCloseFlag, true)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	dao.db.Close()
	dao.lock.Unlock()
}

//SaveRegistryAddress save registry address to db
func (dao *GkvDB) SaveRegistryAddress(registryAddress common.Address) {
	err := dao.saveKeyValueToBucket(models.BucketMeta, models.KeyRegistry, registryAddress)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//GetRegistryAddress returns registry address in db
func (dao *GkvDB) GetRegistryAddress() common.Address {
	var registry common.Address
	err := dao.getKeyValueToBucket(models.BucketMeta, models.KeyRegistry, &registry)
	if err != nil && err != ErrorNotFound {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	return registry
}

//SaveSecretRegistryAddress save secret registry contract address to db
func (dao *GkvDB) SaveSecretRegistryAddress(secretRegistryAddress common.Address) {
	err := dao.saveKeyValueToBucket(models.BucketMeta, models.KeySecretRegistry, secretRegistryAddress)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//GetSecretRegistryAddress return secret registry contract address
func (dao *GkvDB) GetSecretRegistryAddress() common.Address {
	var secretRegistry common.Address
	err := dao.getKeyValueToBucket(models.BucketMeta, models.KeySecretRegistry, &secretRegistry)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	return secretRegistry
}
