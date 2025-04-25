package zio

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestLimitReaderErr(t *testing.T) {
	reader := LimitReaderErr(bytes.NewReader([]byte("hello")), 5, errors.New("test error"))

	buf := make([]byte, 10)
	n, err := reader.Read(buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n != 5 {
		t.Fatalf("expected 5 bytes, got %d", n)
	}
}

func TestLimitReaderErr_Error(t *testing.T) {
	var ErrTest = errors.New("test error")
	reader := LimitReaderErr(bytes.NewReader([]byte("hello")), 4, ErrTest)

	buf := make([]byte, 10)
	_, err := reader.Read(buf)
	if err != ErrTest {
		t.Fatalf("expected %v, got %v", ErrTest, err)
	}
}

func TestLimitReaderErr_ReadAll(t *testing.T) {
	reader := LimitReaderErr(bytes.NewReader([]byte("hello")), 5, errors.New("test error"))

	all, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(all) != "hello" {
		t.Fatalf("expected %v, got %v", "hello", string(all))
	}
}
