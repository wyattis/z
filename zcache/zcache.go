package zcache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/wyattis/z/zsignal"
)

type indexItem struct {
	id   interface{}
	data interface{}
}

func New(dir CacheFS, maxSize int) (c *Cache, err error) {
	c = &Cache{
		fs: dir,
	}
	c.Cache, err = lru.New(maxSize)
	return
}

type Cache struct {
	*lru.Cache
	sig           *zsignal.Signal
	Debounce      time.Duration
	isInitialized bool
	mut           *sync.Mutex
	cond          *sync.Cond
	fs            CacheFS
	prevOldestId  string
	prevLength    int
}

func (c *Cache) init() (err error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	if c.isInitialized {
		return
	}
	c.isInitialized = true
	f, err := c.fs.Open(indexPath)
	if err == nil {
		index := []indexItem{}
		defer f.Close()
		dec := gob.NewDecoder(f)
		if err = dec.Decode(&index); err != nil {
			return
		}
		for _, item := range index {
			c.Add(item.id, item.data)
		}
		oldestKey, _, exists := c.GetOldest()
		if exists {
			c.prevOldestId = oldestKey.(string)
			c.prevLength = c.Len()
		}
		fmt.Printf("loaded %d existing items\n", c.Len())
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	c.sig = zsignal.New()
	if c.Debounce == 0 {
		c.Debounce = time.Second * 5
	}
	go c.sig.Debounce(c.Debounce, c.sync)
	return
}

func (c *Cache) sync() (err error) {
	oldestKey, _, exists := c.GetOldest()
	oldestId := ""
	if exists {
		oldestId = oldestKey.(string)
	}
	l := c.Len()
	if exists && (c.prevOldestId != oldestId || c.prevLength != l) {
		items := []indexItem{}
		for _, key := range c.Keys() {
			v, exists := c.Get(key)
			if exists {
				items = append(items, indexItem{id: key, data: v})
			}
		}
		fmt.Printf("syncing %d items\n", len(items))
		f, err := c.fs.Create(indexPath)
		if err != nil {
			return err
		}
		defer f.Close()
		enc := gob.NewEncoder(f)
		c.prevOldestId = oldestId
		c.prevLength = l
		return enc.Encode(items)
	}
	return
}

func (c *Cache) Add(key, val interface{}) (evicted bool) {
	c.sig.Notify()
	return c.Cache.Add(key, val)
}

func (c *Cache) ContainsOrAdd(key, val interface{}) (existed, evicted bool) {
	existed, evicted = c.Cache.ContainsOrAdd(key, val)
	if !existed || evicted {
		c.sig.Notify()
	}
	return
}

func (c *Cache) Purge() {
	c.sig.Notify()
	c.Cache.Purge()
}

func (c *Cache) PeekOrAdd(key, val interface{}) (prev interface{}, existed, evicted bool) {
	prev, existed, evicted = c.Cache.PeekOrAdd(key, val)
	if !existed || evicted {
		c.sig.Notify()
	}
	return
}

func (c *Cache) Remove(key interface{}) (present bool) {
	present = c.Cache.Remove(key)
	if present {
		c.sig.Notify()
	}
	return
}

func (c *Cache) RemoveAll() (key, val interface{}, removed bool) {
	key, val, removed = c.Cache.RemoveOldest()
	if removed {
		c.sig.Notify()
	}
	return
}
