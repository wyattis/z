package zoc

import (
	"errors"
	"os"
)

// Determine if a path exists or not
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
