package zio

import (
	"errors"
	"io"
)

func NewReaderToSeeker(reader io.Reader) *ReaderToSeeker {
	return &ReaderToSeeker{
		reader:      reader,
		buf:         make([]byte, 0),
		IsRecording: true,
	}
}

// Add seeking capability to a reader by storing the read data in memory
type ReaderToSeeker struct {
	reader      io.Reader
	buf         []byte
	cursor      int64
	IsRecording bool
}

func (r *ReaderToSeeker) Read(d []byte) (n int, err error) {
	if _, ok := r.reader.(io.Seeker); ok {
		return r.reader.Read(d)
	}
	// Handle periods where the cursor is inside the buffer
	if int(r.cursor) < len(r.buf) {
		n = copy(d, r.buf[r.cursor:])
		r.cursor += int64(n)
		// we could also handle reads past len(r.buf) by reading more data from the reader, but that adds some complexity
		return
	}

	//  Handle copying data from the underlying reader into the buffer
	n, err = r.reader.Read(d)
	if err != nil && err != io.EOF {
		return
	}
	r.buf = append(r.buf, d[:n]...)
	r.cursor += int64(n)
	return
}

func (r *ReaderToSeeker) Seek(offset int64, whence int) (newOffset int64, err error) {
	if s, ok := r.reader.(io.Seeker); ok {
		return s.Seek(offset, whence)
	}
	switch whence {
	case io.SeekCurrent:
		offset += r.cursor
	case io.SeekEnd:
		err = errors.New("cannot seek from the end as the end is not known")
		return
	}

	// Handle seeking ahead of buffer position
	if int(offset) > len(r.buf) {
		b := make([]byte, int(offset)-len(r.buf))
		n, err := r.reader.Read(b)
		if err != nil && err != io.EOF {
			return int64(n), err
		}
		r.buf = append(r.buf, b[:n]...)
		offset = int64(len(r.buf))
	}
	r.cursor = offset
	return r.cursor, err
}

func (r *ReaderToSeeker) Close() (err error) {
	if c, ok := r.reader.(io.Closer); ok {
		return c.Close()
	}
	return
}
