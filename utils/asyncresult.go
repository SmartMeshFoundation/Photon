package utils

import "github.com/ethereum/go-ethereum/common"

/*
AsyncResult is designed for async notify
and Tag can be save anything by user.
*/
type AsyncResult struct {
	Result         chan error
	Tag            interface{}
	LockSecretHash common.Hash // only for /api/1/transfer use, return LockSecretHash to caller
}

//NewAsyncResult create a AsyncResult
func NewAsyncResult() *AsyncResult {
	return &AsyncResult{Result: make(chan error, 1)}
}

//NewAsyncResultWithError create AsyncResult with result
func NewAsyncResultWithError(err error) *AsyncResult {
	r := &AsyncResult{
		Result: make(chan error, 1),
	}
	r.Result <- err
	return r
}
