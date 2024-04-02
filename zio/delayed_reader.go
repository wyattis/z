package zio

import (
	"io"
	"time"
)

func DelayedReader(r io.Reader, delay time.Duration) io.Reader {
	return &delayedReader{r: r, delay: delay}
}

type delayedReader struct {
	r     io.Reader
	delay time.Duration
}

func (r *delayedReader) Read(p []byte) (n int, err error) {
	time.Sleep(r.delay)
	return r.r.Read(p)
}
