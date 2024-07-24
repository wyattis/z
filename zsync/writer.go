package zsync

import (
	"io"
	"sync"
)

func Writer(w io.Writer) io.Writer {
	return &writer{
		mut: &sync.Mutex{},
		w:   w,
	}
}

type writer struct {
	mut *sync.Mutex
	w   io.Writer
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.w.Write(p)
}
