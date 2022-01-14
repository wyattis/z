package zcache

import (
	"io"
	"os"
	"path/filepath"
)

type CacheFS interface {
	Open(path string) (io.ReadSeekCloser, error)
	Create(path string) (io.WriteCloser, error)
	Remove(path string) error
	RemoveAll(path string) error
	Rename(from, to string) error
	MkdirAll(dir string, perm os.FileMode) error
	ReadDir(dir string) ([]os.DirEntry, error)
}

func DirFS(root string) *osDirFS {
	return &osDirFS{root}
}

type osDirFS struct {
	root string
}

func (d *osDirFS) Open(path string) (io.ReadSeekCloser, error) {
	path = filepath.Join(d.root, path)
	return os.Open(path)
}

func (d *osDirFS) Create(path string) (io.WriteCloser, error) {
	path = filepath.Join(d.root, path)
	return os.Create(path)
}

func (d *osDirFS) Remove(path string) error {
	path = filepath.Join(d.root, path)
	return os.Remove(path)
}

func (d *osDirFS) RemoveAll(path string) error {
	path = filepath.Join(d.root, path)
	return os.RemoveAll(path)
}

func (d *osDirFS) MkdirAll(path string, perm os.FileMode) error {
	path = filepath.Join(d.root, path)
	return os.MkdirAll(path, perm)
}

func (d *osDirFS) Rename(from, to string) error {
	from, to = filepath.Join(d.root, from), filepath.Join(d.root, to)
	return os.Rename(from, to)
}

func (d *osDirFS) ReadDir(path string) ([]os.DirEntry, error) {
	path = filepath.Join(d.root, path)
	return os.ReadDir(path)
}
