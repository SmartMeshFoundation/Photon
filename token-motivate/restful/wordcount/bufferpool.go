package wordcount

import (
	"bytes"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// GetBuffer returns a buffer from the pool.
func GetBuffer() (buf *bytes.Buffer) {
	return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool.
// The buffer is reset before it is put back into circulation.
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}
