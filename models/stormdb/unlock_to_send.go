package stormdb

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

//NewUnlockToSend save a UnlockToSend to db
func (model *StormDB) NewUnlockToSend(lockSecretHash common.Hash, tokenAddress, receiver common.Address, blockNumber int64) *models.UnlockToSend {

	key := utils.Sha3(lockSecretHash[:], tokenAddress[:], receiver[:]).Bytes()
	us := &models.UnlockToSend{
		Key:              key,
		LockSecretHash:   lockSecretHash[:],
		TokenAddress:     tokenAddress[:],
		ReceiverAddress:  receiver[:],
		SavedTimestamp:   time.Now().Unix(),
		SavedBlockNumber: blockNumber,
	}
	err := model.db.Save(us)
	if err != nil {
		log.Error(fmt.Sprintf("NewUnlockToSend err %s", err))
	}
	return us
}

// GetAllUnlockToSend query all
func (model *StormDB) GetAllUnlockToSend() (list []*models.UnlockToSend) {
	err := model.db.All(&list)
	if err == storm.ErrNotFound {
		err = nil
	}
	if err != nil {
		log.Error(fmt.Sprintf("GetAllUnlockToSend err %s", err))
	}
	return
}

// RemoveUnlockToSend remove by primary key
func (model *StormDB) RemoveUnlockToSend(key []byte) {
	err := model.db.DeleteStruct(&models.UnlockToSend{
		Key: key,
	})
	if err != nil {
		log.Error(fmt.Sprintf("RemoveUnlockToSend err %s", err))
	}
}
