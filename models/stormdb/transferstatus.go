package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

// NewTransferStatus :
func (model *StormDB) NewTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash) {
	ts := &models.TransferStatus{
		Key:            utils.Sha3(tokenAddress[:], lockSecretHash[:]).String(),
		LockSecretHash: lockSecretHash,
		TokenAddress:   tokenAddress,
		Status:         models.TransferStatusInit,
		StatusMessage:  "",
	}
	err := model.db.Save(ts)
	if err != nil {
		log.Error(fmt.Sprintf("NewTransferStatus key=%s, err %s", ts.Key, err))
		return
	}
	log.Trace(fmt.Sprintf("NewTransferStatus key=%s lockSecertHash=%s", ts.Key, lockSecretHash.String()))
}

// UpdateTransferStatus :
func (model *StormDB) UpdateTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash, status models.TransferStatusCode, statusMessage string) {
	var ts models.TransferStatus
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := model.db.One("Key", key, &ts)
	if err == storm.ErrNotFound {
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
	ts.Status = status
	ts.StatusMessage = fmt.Sprintf("%s%s\n", ts.StatusMessage, statusMessage)
	err = model.db.Save(&ts)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
	log.Trace(fmt.Sprintf("UpdateTransferStatus key=%s lockSecretHash=%s %s", key, lockSecretHash.String(), statusMessage))
}

// UpdateTransferStatusMessage :
func (model *StormDB) UpdateTransferStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string) {
	var ts models.TransferStatus
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := model.db.One("Key", key, &ts)
	if err == storm.ErrNotFound {
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatusMessage err %s", err))
		return
	}
	ts.StatusMessage = fmt.Sprintf("%s%s\n", ts.StatusMessage, statusMessage)
	err = model.db.Save(&ts)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatusMessage err %s", err))
		return
	}
	log.Trace(fmt.Sprintf("UpdateTransferStatusMessage key=%s lockSecretHash=%s %s", key, lockSecretHash.String(), statusMessage))
}

// GetTransferStatus :
func (model *StormDB) GetTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash) (*models.TransferStatus, error) {
	var ts models.TransferStatus
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := model.db.One("Key", key, &ts)
	log.Trace(fmt.Sprintf("GetTransferStatus key=%s lockSecretHash=%s err=%s", key, lockSecretHash.String(), err))
	err = models.GeneratDBError(err)
	return &ts, err
}
