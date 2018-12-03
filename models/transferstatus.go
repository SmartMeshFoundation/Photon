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
	Key            string `storm:"id"`
	LockSecretHash common.Hash
	TokenAddress   common.Address
	Status         TransferStatusCode
	StatusMessage  string
}

func init() {
	gob.Register(&TransferStatus{})
}
