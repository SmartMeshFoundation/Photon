package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
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
)

/*
TransferStatus :
	save status of transfer for api, most time for debug
*/
type TransferStatus struct {
	LockSecretHash common.Hash `storm:"id"`
	Status         TransferStatusCode
	StatusMessage  string
}

// UpdateTransferStatus :
func (model *ModelDB) UpdateTransferStatus(lockSecretHash common.Hash, status TransferStatusCode, statusMessage string) {
	var ts TransferStatus
	err := model.db.One("LockSecretHash", lockSecretHash, &ts)
	if err == storm.ErrNotFound {
		ts = TransferStatus{}
		err = nil
	}
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
	ts.LockSecretHash = lockSecretHash
	ts.Status = status
	ts.StatusMessage = fmt.Sprintf("%s%s\n", ts.StatusMessage, statusMessage)
	err = model.db.Save(&ts)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatus err %s", err))
		return
	}
}

// UpdateTransferStatusMessage :
func (model *ModelDB) UpdateTransferStatusMessage(lockSecretHash common.Hash, statusMessage string) {
	var ts TransferStatus
	err := model.db.One("LockSecretHash", lockSecretHash, &ts)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTransferStatus err %s", err))
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
func (model *ModelDB) GetTransferStatus(lockSecretHash common.Hash) (*TransferStatus, error) {
	var ts TransferStatus
	err := model.db.One("LockSecretHash", lockSecretHash, &ts)
	return &ts, err
}
