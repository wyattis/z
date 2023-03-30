package zio

import (
	"errors"
	"io"
)

var (
	ErrExceededMax = errors.New("Exceeded limit of maxReader")
)

func MaxReader(reader io.Reader, limit int) io.Reader {
	return &maxReader{Reader: reader, Limit: limit}
}

type maxReader struct {
	io.Reader
	bytesRead int
	Limit     int
}

func (r *maxReader) Read(d []byte) (n int, err error) {
	n, err = r.Reader.Read(d)
	if err != nil {
		return
	}
	r.bytesRead += n
	if r.bytesRead > r.Limit {
		err = ErrExceededMax
		return
	}
	return
}
