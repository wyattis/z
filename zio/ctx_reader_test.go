package zio

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"strings"
	"testing"
	"time"
)

func TestCtxReaderSingleRead(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	r := strings.NewReader("hello")
	cr := CtxReader(ctx, r)
	buf := make([]byte, 5)
	n, err := cr.Read(buf)
	if n != 5 || err != nil {
		t.Errorf("Read() = %d, %v, want 5, nil", n, err)
	}
}

func TestCtxReaderIO(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	nBytes := int64(10000)
	r := io.LimitReader(rand.Reader, nBytes)
	truth := bytes.Buffer{}
	r = io.TeeReader(r, &truth)
	cr := CtxReader(ctx, r)
	crRes := bytes.Buffer{}
	n, err := io.Copy(&crRes, cr)
	if err != nil {
		t.Errorf("Copy() = %v, want nil", err)
	}
	if n != nBytes {
		t.Errorf("Copy() = %d, want %d", n, nBytes)
	}
}

func TestCtxReaderCancels(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	r := DelayedReader(strings.NewReader("hello"), 2*time.Second)
	cr := CtxReader(ctx, r)
	buf := make([]byte, 5)
	n, err := cr.Read(buf)
	if n != 0 || err == nil {
		t.Errorf("Read() = %d, %v, want 0, non-nil", n, err)
	}
}
