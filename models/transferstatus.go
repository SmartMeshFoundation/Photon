package models

import (
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
)

/*
TransferStatusCode status of transfer
*/
type TransferStatusCode int

const (

	// TransferStatusInit init
	TransferStatusInit = iota

	// TransferStatusCanCancel transfer can cancel right now
	TransferStatusCanCancel

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
	Key            string             `storm:"id" json:"key,omitempty"`
	LockSecretHash common.Hash        `json:"lock_secret_hash"`
	TokenAddress   common.Address     `json:"token_address"`
	Status         TransferStatusCode `json:"status"`
	StatusMessage  string             `json:"status_message"`
}

func init() {
	gob.Register(&TransferStatus{})
}
