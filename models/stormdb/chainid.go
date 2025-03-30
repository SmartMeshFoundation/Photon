package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
)

//GetChainID 查询ChainID
func (model *StormDB) GetChainID() int64 {
	var chainID int64
	err := model.db.Get(models.BucketChainID, models.KeyChainID, &chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models GetChainId err=%s", err))
	}
	return chainID
}

//SaveChainID 保存/更新数据库中的ChainID
func (model *StormDB) SaveChainID(chainID int64) {
	err := model.db.Set(models.BucketChainID, models.KeyChainID, chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models SaveChainId err=%s", err))
	}
}
