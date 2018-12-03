package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/asdine/storm"
)

const defaultKey string = "feePolicy"

// SaveFeePolicy :
func (model *StormDB) SaveFeePolicy(fp *models.FeePolicy) (err error) {
	fp.Key = defaultKey
	err = model.db.Save(fp)
	return
}

// GetFeePolicy :
func (model *StormDB) GetFeePolicy() (fp *models.FeePolicy) {
	fp = &models.FeePolicy{}
	err := model.db.One("Key", defaultKey, fp)
	if err == storm.ErrNotFound {
		return models.NewDefaultFeePolicy()
	}
	if err != nil {
		log.Error(fmt.Sprintf("GetFeePolicy err %s, use default fee policy", err))
		return models.NewDefaultFeePolicy()
	}
	return
}
