package main

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func NewTarDirReader(dir string) *tarDirReader {
	return &tarDirReader{
		fs: os.DirFS(dir),
	}
}

type tarDirReader struct {
	fs          fs.FS
	queue       []fs.DirEntry
	reader      *io.PipeReader
	writer      *io.PipeWriter
	archive     *tar.Writer
	initialized bool
}

func (t *tarDirReader) init() (err error) {
	if t.initialized {
		return
	}
	t.initialized = true
	t.reader, t.writer = io.Pipe()
	t.archive = tar.NewWriter(t.writer)
	err = fs.WalkDir(t.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// TODO: some kind of filtering using .ignore files
		t.queue = append(t.queue, d)
		return nil
	})
	if err != nil {
		return
	}
	go t.writeQueue()
	return
}

func (t *tarDirReader) writeFile(f fs.DirEntry) (err error) {
	file, err := t.fs.Open(f.Name())
	if err != nil {
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = f.Name()
	if err = t.archive.WriteHeader(header); err != nil {
		return
	}
	if _, err = io.Copy(t.archive, file); err != nil {
		return
	}
	return t.archive.Flush()
}

func (t *tarDirReader) writeQueue() {
	for _, f := range t.queue {
		if f.IsDir() {
			continue
		}
		fmt.Println("writing file to tar", f.Name())
		if err := t.writeFile(f); err != nil {
			panic(err)
		}
	}
	if err := t.archive.Close(); err != nil {
		panic(err)
	}
	if err := t.writer.Close(); err != nil {
		panic(err)
	}
	fmt.Println("writing complete")
}

func (t *tarDirReader) Read(p []byte) (n int, err error) {
	if err = t.init(); err != nil {
		return
	}
	n, err = t.reader.Read(p)
	return
}
