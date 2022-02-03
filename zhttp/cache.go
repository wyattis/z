package zhttp

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wyattis/z/zpath"
	"github.com/wyattis/z/zsync"
)

type ReadUnlocker struct {
	io.ReadCloser
	locker *sync.RWMutex
}

func (u *ReadUnlocker) Close() error {
	defer u.locker.RUnlock()
	return u.ReadCloser.Close()
}

type WriteUnlocker struct {
	io.WriteCloser
	locker sync.Locker
}

func (u *WriteUnlocker) Close() error {
	defer u.locker.Unlock()
	return u.WriteCloser.Close()
}

func NewDirCache(dir string, duration time.Duration) *DirCache {
	return &DirCache{
		dir:           dir,
		locker:        *zsync.NewRWLockerMap(10),
		index:         make(map[string]bool),
		CacheDuration: duration,
	}
}

type DirCache struct {
	dir           string
	CacheDuration time.Duration
	locker        zsync.RWLockerMap
	index         map[string]bool
}

func (c *DirCache) Has(key string) bool {
	key = c.sluggify(key)
	return c.index[key]
}

func (c *DirCache) sluggify(key string) string {
	return zpath.FileEscape(key)
}

func (c *DirCache) Open(key string) (r io.ReadCloser, err error) {
	key = c.sluggify(key)
	mut := c.locker.RLock(key)
	r, err = os.Open(filepath.Join(c.dir, key))
	if err != nil {
		return
	}
	r = &ReadUnlocker{r, mut}
	return
}

func (c *DirCache) Create(key string) (w io.WriteCloser, err error) {
	key = c.sluggify(key)
	mut := c.locker.Lock(key)
	w, err = os.Create(filepath.Join(c.dir, key))
	if err != nil {
		return
	}
	w = &WriteUnlocker{w, mut}
	return
}

func (c *DirCache) Remove(key string) (err error) {
	key = c.sluggify(key)
	mut := c.locker.Lock(key)
	defer mut.Unlock()
	err = os.Remove(filepath.Join(c.dir, key))
	return
}
