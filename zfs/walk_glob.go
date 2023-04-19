package zfs

import (
	"io/fs"

	"github.com/wyattis/z/zset/zstringset"
)

// func Match(path, pattern string) (bool, error) {
// 	patternParts := filepath.SplitList(pattern)
// 	pathParts := filepath.SplitList(path)
// 	j := len(pathParts) - 1
// 	for i := len(patternParts) - 1; i >= 0; i-- {
// 		patternPart := patternParts[i]
// 		pathPart := pathParts[j]
// 		if patternPart == "**" {
// 			j--
// 			continue
// 		} else if strings.Contains(patternPart, "*") {

// 		}
// 	}
// 	return filepath.Match(pattern, path)
// }

type WalkFS interface {
	fs.GlobFS
	fs.StatFS
	fs.ReadDirFS
}

type InfoDirEntry struct {
	fs.FileInfo
}

func (d InfoDirEntry) Type() fs.FileMode {
	return d.FileInfo.Mode().Type()
}

func (d InfoDirEntry) Info() (fs.FileInfo, error) {
	return d.FileInfo, nil
}

func walkMatches(fs WalkFS, matches []string, fn fs.WalkDirFunc) (err error) {
	queue := zstringset.New(matches...)
	for _, match := range matches {
		if !queue.Contains(match) {
			continue
		}
		queue.Delete(match)
		info, err := fs.Stat(match)
		if err != nil {
			return err
		}
		if err = fn(match, &InfoDirEntry{info}, nil); err != nil {
			return err
		}
		if info.IsDir() {
			// read the dir directly to make everything faster
			entries, err := fs.ReadDir(match)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if queue.Contains(entry.Name()) {
					queue.Delete(entry.Name())
					if err = fn(entry.Name(), entry, nil); err != nil {
						return err
					}
				}
			}
		}
	}
	return
}

// WalkDirGlob walks files matching the given pattern and calls fn for each file/directory
func WalkDirGlob(fs WalkFS, pattern string, fn fs.WalkDirFunc) (err error) {
	matches, err := fs.Glob(pattern)
	if err != nil {
		return
	}
	return walkMatches(fs, matches, fn)
}

// WalkDirGlobs walks files matching the given patterns and calls fn for each file/directory
func WalkDirGlobs(fs WalkFS, patterns []string, fn fs.WalkDirFunc) (err error) {
	matches := []string{}
	for _, pattern := range patterns {
		pMatches, err := fs.Glob(pattern)
		if err != nil {
			return err
		}
		matches = append(matches, pMatches...)
	}
	return walkMatches(fs, matches, fn)
}
