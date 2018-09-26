package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

/*
TransferStatusCode status of transfer
*/
type TransferStatusCode int

const (

	// TransferStatusCanCancel transfer can cancel right now
	TransferStatusCanCancel = iota

	// TransferStatusCanNotCancel transfer can not cancel
	TransferStatusCanNotCancel

	// TransferStatusSuccess transfer already success
	TransferStatusSuccess

	// TransferStatusCanceled transfer cancel by user request
	TransferStatusCanceled

	// TransferStatusFailed transfer already failed
	TransferStatusFailed
)

/*
TransferStatus :
	save status of transfer for api, most time for debug
*/
type TransferStatus struct {
	Key            common.Hash `storm:"id"`
	LockSecretHash common.Hash
	TokenAddress   common.Address
	Status         TransferStatusCode
	StatusMessage  string
}

// UpdateTransferStatus :
func (model *ModelDB) UpdateTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash, status TransferStatusCode, statusMessage string) {
	var ts TransferStatus
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:])
	err := model.db.One("Key", key, &ts)
	if err == storm.ErrNotFound {
		ts = TransferStatus{
			Key:            key,
			LockSecretHash: lockSecretHash,
			TokenAddress:   tokenAddress,
		}
		err = nil
	}
	if err != nil {
		log.Warn(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
	ts.Status = status
	ts.StatusMessage = fmt.Sprintf("%s%s\n", ts.StatusMessage, statusMessage)
	err = model.db.Save(&ts)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
}

// UpdateTransferStatusMessage :
func (model *ModelDB) UpdateTransferStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string) {
	var ts TransferStatus
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:])
	err := model.db.One("Key", key, &ts)
	if err != nil {
		log.Warn(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
	ts.StatusMessage = fmt.Sprintf("%s%s\n", ts.StatusMessage, statusMessage)
	err = model.db.Save(&ts)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
}

// GetTransferStatus :
func (model *ModelDB) GetTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash) (*TransferStatus, error) {
	var ts TransferStatus
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:])
	err := model.db.One("Key", key, &ts)
	return &ts, err
}
