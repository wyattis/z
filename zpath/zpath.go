package zpath

import (
	"path/filepath"
	"strings"

	"github.com/wyattis/z/zstring"
)

// Determine if a path is child of the root path.
func IsChildOf(root, child string) (isChild bool, err error) {
	if !filepath.IsAbs(root) {
		root, err = filepath.Abs(root)
		if err != nil {
			return
		}
	}
	if !filepath.IsAbs(child) {
		child, err = filepath.Abs(child)
		if err != nil {
			return
		}
	}
	childParts := strings.Split(filepath.ToSlash(child), "/")
	rootParts := strings.Split(filepath.ToSlash(root), "/")
	if len(childParts) < len(rootParts) {
		return false, nil
	}
	for i := range rootParts {
		if rootParts[i] != childParts[i] {
			return false, nil
		}
	}
	return true, nil
}

// Replace the extension in a filepath with the new extension
func ReplaceExt(path, newExt string) string {
	ext := filepath.Ext(path)
	return strings.Replace(path, ext, newExt, len(ext)-1)
}

var escapeChars = []string{"/", "\\", ":", "?", "&", "="}

// Escape a filename to remove all potentially invalid characters. Reduces a
// path with multiple directories to a single result
func FileEscape(str string) string {
	str = zstring.TrimSuffixes(str, "/", "\\")
	str = zstring.TrimPrefixes(str, "/", "\\")
	return zstring.ReplaceManyWithOne(str, escapeChars, "_")
}
