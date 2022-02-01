package zio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Close multiple Closers returning an error if any of them produce one
func CloseAll(files ...io.Closer) error {
	for _, f := range files {
		if f != nil {
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Determine if two readers have the same data. Returns an error if they don't
func ReadersMatch(r1 io.Reader, r2 io.Reader, bufSize int) error {
	if bufSize == 0 {
		bufSize = 32 * 1024
	}
	buf1 := make([]byte, bufSize)
	buf2 := make([]byte, bufSize)
	var err1, err2 error
	var n1, n2 int
	c := 0
	for err1 != io.EOF && err2 != io.EOF {
		n1, err1 = r1.Read(buf1)
		if err1 != nil && err1 != io.EOF {
			return err1
		}
		n2, err2 = r2.Read(buf2)
		if err2 != nil && err2 != io.EOF {
			return err2
		}
		c += n2
		if n1 != n2 {
			return errors.New(fmt.Sprintf("Expected n1 == n2, but got %d != %d", n1, n2))
		}
		if !bytes.Equal(buf1, buf2) {
			return errors.New(fmt.Sprint("didn't read the same bytes from the readers around", c, "\n", buf1, "\n", buf2))
		}
	}
	return nil
}