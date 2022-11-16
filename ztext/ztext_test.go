package ztext

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

var truth = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

func TestUTF8Reader(t *testing.T) {
	filepath.WalkDir("./test", func(path string, d fs.DirEntry, err error) (err2 error) {
		if d.IsDir() {
			return
		}
		f, err := os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		reader := NewUTF8Reader(f)
		data, err := io.ReadAll(reader)
		if err != nil {
			return
		}
		if string(truth) != string(data) {
			t.Errorf("%s\nExpected\n%s\nreceived\n%s\n", path, string(truth), string(data))
		}
		return
	})
}
