package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
)

const bucketChainID = "bucketChainID"
const keyChainID = "chainID"

//GetChainID :
func (model *StormDB) GetChainID() int64 {
	var chainID int64
	err := model.db.Get(bucketChainID, keyChainID, &chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models GetChainId err=%s", err))
	}
	return chainID
}

//SaveChainID :
func (model *StormDB) SaveChainID(chainID int64) {
	err := model.db.Set(bucketChainID, keyChainID, chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models SaveChainId err=%s", err))
	}
}
