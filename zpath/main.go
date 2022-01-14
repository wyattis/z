package zpath

import (
	"path/filepath"
	"strings"
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
		child, err = filepath.Abs(root)
		if err != nil {
			return
		}
	}
	childParts := strings.Split(child, string(filepath.Separator))
	rootParts := strings.Split(child, string(filepath.Separator))
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
