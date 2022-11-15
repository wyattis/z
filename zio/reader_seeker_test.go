package zio

import (
	"bytes"
	"io"
	"testing"
)

func TestSimpleSeekStart(t *testing.T) {
	r := NewReaderToSeeker(bytes.NewBuffer(make([]byte, 100)))
	offset, err := r.Seek(100, io.SeekStart)
	if err != nil {
		t.Error(err)
	}
	if offset != 100 {
		t.Errorf("expected offset to be %d, but got %d instead", 100, offset)
	}
}

func FuzzReaderSeekerSeekStart(f *testing.F) {
	f.Add([]byte("7456rftgo12345rolkasjhdgv0opc89zx7i\"|||vbas;kl;][;]';'.,"), uint(4), 50)
	f.Add([]byte{1, 45, 255, 3, 65}, uint(2), 0)
	f.Fuzz(func(t *testing.T, data []byte, initRead uint, seek int) {
		// Validate the input
		if seek < 0 || seek > len(data) || initRead > uint(len(data)) {
			t.SkipNow()
		}
		r := NewReaderToSeeker(bytes.NewBuffer(data))
		d := make([]byte, initRead)
		n, err := io.ReadFull(r, d)
		if err != nil {
			t.Error(err)
		}
		if uint(n) != initRead {
			t.Errorf("expected to read %d, but read %d instead", initRead, n)
		}
		if !bytes.Equal(data[:initRead], d) {
			t.Errorf("expected %b, but got %b instead", data[:initRead], d)
		}
		offset, err := r.Seek(int64(seek), io.SeekStart)
		if err != nil {
			t.Error(err)
		}
		if int(offset) != seek {
			t.Errorf("expected offset to be %d, but got %d instead", seek, offset)
		}
		remaining := len(data) - seek
		d = make([]byte, remaining)
		n, err = io.ReadFull(r, d)
		if err != nil {
			t.Error(err)
		}
		if n != remaining {
			t.Errorf("expected to read remaining %d, but read %d instead", remaining, n)
		}
		if !bytes.Equal(d, data[seek:]) {
			t.Errorf("expected to read %b, but got %b instead", data[seek:], d)
		}
	})
}
