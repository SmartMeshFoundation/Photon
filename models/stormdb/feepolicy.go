package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/asdine/storm"
)

// SaveFeePolicy 创建/更新记录
func (model *StormDB) SaveFeePolicy(fp *models.FeePolicy) (err error) {
	fp.Key = models.KeyFeePolicy
	err = model.db.Save(fp)
	err = models.GeneratDBError(err)
	return
}

// GetFeePolicy 查询
func (model *StormDB) GetFeePolicy() (fp *models.FeePolicy) {
	fp = &models.FeePolicy{}
	err := model.db.One("Key", models.KeyFeePolicy, fp)
	if err == storm.ErrNotFound {
		return models.NewDefaultFeePolicy()
	}
	if err != nil {
		log.Error(fmt.Sprintf("GetFeePolicy err %s, use default fee policy", err))
		return models.NewDefaultFeePolicy()
	}
	return
}
