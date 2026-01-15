package fatek

import (
	"bytes"
	"sync"
)

// bufferPool is a pool of bytes.Buffer to reduce memory allocation pressure.
var bufferPool = sync.Pool{
	New: func() interface{} {
		// Pre-allocate a buffer of 64 bytes (sufficient for most Fatek ASCII frames)
		// It will grow automatically if needed.
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
}

// GetBuffer retrieves a buffer from the pool and resets it.
// The caller must return the buffer using PutBuffer.
func GetBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a buffer to the pool.
// It performs a safety check to avoid keeping excessively large buffers in memory.
func PutBuffer(buf *bytes.Buffer) {
	// If the buffer grew too large (e.g., > 4KB), drop it to let GC collect it.
	// This prevents memory leaks if a huge anomalous packet was processed.
	if buf.Cap() > 4096 {
		return
	}
	bufferPool.Put(buf)
}
