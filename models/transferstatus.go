package models

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
