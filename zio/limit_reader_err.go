package zio

import (
	"io"
)

// LimitReaderErr is the samw as io.LimitReader except you can specify an error to return when the limit is reached.
func LimitReaderErr(r io.Reader, limit int64, err error) io.Reader {
	return &limitReaderErr{
		Reader: r,
		err:    err,
		limit:  limit,
	}
}

type limitReaderErr struct {
	io.Reader
	err   error
	limit int64
	read  int64
}

func (l *limitReaderErr) Read(p []byte) (n int, err error) {
	n, err = l.Reader.Read(p)
	l.read += int64(n)
	if l.read > l.limit {
		if l.err == nil {
			l.err = io.EOF
		}
		return n, l.err
	}
	return n, err
}
