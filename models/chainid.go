package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
)

const bucketChainID = "bucketChainID"
const keyChainID = "chainID"

//GetChainID :
func (model *ModelDB) GetChainID() int64 {
	var chainID int64
	err := model.db.Get(bucketChainID, keyChainID, &chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models GetChainId err=%s", err))
	}
	return chainID
}

//SaveChainID :
func (model *ModelDB) SaveChainID(chainID int64) {
	err := model.db.Set(bucketChainID, keyChainID, chainID)
	if err != nil {
		log.Error(fmt.Sprintf("models SaveChainId err=%s", err))
	}
}
