package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
)

//GetChainID :
func (dao *GkvDB) GetChainID() int64 {
	var chainID int64
	err := dao.getKeyValueToBucket(models.BucketChainID, models.KeyChainID, &chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models GetChainId err=%s", err))
	}
	return chainID
}

//SaveChainID :
func (dao *GkvDB) SaveChainID(chainID int64) {
	err := dao.saveKeyValueToBucket(models.BucketChainID, models.KeyChainID, chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models SaveChainId err=%s", err))
	}
}
