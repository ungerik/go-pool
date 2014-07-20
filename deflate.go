package pool

import (
	"compress/flate"
	"io"
	"sync"
)

var Deflate DeflatePool

// DeflatePool manages a pool of flate.Writer
// flate.NewWriter allocates a lot of memory, so if flate.Writer
// are needed frequently, it's more efficient to use a pool of them.
// The pool uses sync.Pool internally.
// 
// There is no pool for flate readers, because they implement io.ReadCloser
// and can't be recycled after closing.
type DeflatePool struct {
	writers sync.Pool
}

// GetWriter returns flate.Writer from the pool, or creates a new one
// with flate.BestCompression if the pool is empty.
func (pool *DeflatePool) GetWriter(dst io.Writer) (writer *flate.Writer) {
	if w := pool.writers.Get(); w != nil {
		writer = w.(*flate.Writer)
		writer.Reset(dst)
	} else {
		writer, _ = flate.NewWriter(dst, flate.BestCompression)
	}
	return writer
}

// PutWriter returns a flate.Writer to the pool that can
// late be reused via GetWriter.
// Don't close the writer, Flush will be called before returning
// it to the pool.
func (pool *DeflatePool) PutWriter(writer *flate.Writer) {
	writer.Close()
	pool.writers.Put(writer)
}
