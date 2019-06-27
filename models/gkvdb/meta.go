package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
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
		panic(fmt.Sprintf("db meta data error"))
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

//SaveContractStatus save registry address to db
func (dao *GkvDB) SaveContractStatus(contractStatus models.ContractStatus) {
	err := dao.saveKeyValueToBucket(models.BucketMeta, models.KeyRegistry, contractStatus)
	if err != nil {
		log.Error(fmt.Sprintf("db err %s", err))
	}
}

//GetContractStatus returns registry address in db
func (dao *GkvDB) GetContractStatus() models.ContractStatus {
	var contractStatus models.ContractStatus
	err := dao.getKeyValueToBucket(models.BucketMeta, models.KeyRegistry, &contractStatus)
	if err != nil && err != ErrorNotFound {
		log.Error(fmt.Sprintf("db err %s", err))
	}
	return contractStatus
}
