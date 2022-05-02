package zos

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type FileHandler = func(f *os.File) error

// Determine if a path exists or not
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Returns true if a path exists and is a directory. It still returns false even
// if the path doesn't exist.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
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

// Read all lines from a text file into a slice
func ReadLines(path string) (lines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return ReadLinesFromFile(f)
}

// Read all lines from an open text file into a slice
func ReadLinesFromFile(f *os.File) (lines []string, err error) {
	r := bufio.NewScanner(f)
	r.Split(bufio.ScanLines)
	for r.Scan() {
		if err = r.Err(); err != nil {
			return
		}
		lines = append(lines, r.Text())
	}
	return
}

// Open the first existing path in a list of paths
func OpenFirst(paths ...string) (f *os.File, err error) {
	for _, p := range paths {
		if Exists(p) {
			return os.Open(p)
		}
	}
	err = os.ErrNotExist
	return
}

// Copy a file from one location to another
func Copy(from, to string) (err error) {
	in, err := os.Open(from)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(to)
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return
}

// Copy a file from one location to another using a temp file
func CopyTemp(from, to string) (err error) {
	in, err := os.Open(from)
	if err != nil {
		return
	}
	defer in.Close()
	tmp, err := os.CreateTemp(filepath.Dir(to), filepath.Base(to))
	if err != nil {
		return
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	_, err = io.Copy(tmp, in)
	if err != nil {
		return
	}
	if err = tmp.Close(); err != nil {
		return
	}
	return os.Rename(tmp.Name(), to)
}
