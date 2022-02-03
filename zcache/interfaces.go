package zcache

import "io"

type ReaderWriterCache interface {
	Has(key string) bool
	Open(key string) (r io.ReadCloser, err error)
	Create(key string) (r io.WriteCloser, err error)
	Remove(key string) error
}
