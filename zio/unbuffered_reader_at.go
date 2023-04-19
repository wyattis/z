package zio

import (
	"errors"
	"io"
	"io/ioutil"
)

func NewUnbufferedReaderAt(r io.Reader) io.ReaderAt {
	return &unbufferedReaderAt{R: r}
}

type unbufferedReaderAt struct {
	R io.Reader
	N int64
}

func (u *unbufferedReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if s, ok := u.R.(io.ReaderAt); ok {
		return s.ReadAt(p, off)
	}
	if s, ok := u.R.(io.Seeker); ok {
		if _, err = s.Seek(off, io.SeekStart); err != nil {
			return
		}
		return u.R.Read(p)
	}
	if off < u.N {
		return 0, errors.New("invalid offset")
	}
	diff := off - u.N
	written, err := io.CopyN(ioutil.Discard, u.R, diff)
	u.N += written
	if err != nil {
		return 0, err
	}

	n, err = u.R.Read(p)
	u.N += int64(n)
	return
}
