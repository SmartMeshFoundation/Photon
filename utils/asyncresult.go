package utils

/*
AsyncResult is designed for async notify
and Tag can be save anything by user.
*/
type AsyncResult struct {
	Result chan error
	Tag    interface{}
}

//NewAsyncResult create a AsyncResult
func NewAsyncResult() *AsyncResult {
	return &AsyncResult{Result: make(chan error, 1)}
}
