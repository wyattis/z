package zos

import (
	"errors"
	"os"
)

type FileHandler = func(f *os.File) error

// Determine if a path exists or not
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Create a file and use it inside a closure. Will automatically close the file
// after the closure completes and remove the corrupt file if an error is thrown.
func CreateWith(path string, handler FileHandler) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = handler(f)
	if err != nil {
		f.Close()
		os.Remove(path)
	}
	return
}

// Like CreateWith, but will create a temporary file and rename it to the final
// path once the file has been written successfully. If an error occurs, the
// temp file will be removed.
func CreateWithTemp(path string, handler FileHandler) error {
	tmpPath := path + ".tmp"
	if err := CreateWith(tmpPath, handler); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
