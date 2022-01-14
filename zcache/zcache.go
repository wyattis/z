package zcache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dsnet/compress/xflate"
	lru "github.com/hashicorp/golang-lru"
	"github.com/wyattis/z/zio"
	"github.com/wyattis/z/zpath"
)

type Item struct {
	Id         string
	AccessedAt time.Time
	CreatedAt  time.Time
}

func New(dir CacheFS, maxSize int) (c *Cache, err error) {
	c = &Cache{
		fs:    dir,
		mutex: NewMutexMap(),
	}
	c.index, err = lru.NewWithEvict(maxSize, c.onEvict)
	return
}

type Cache struct {
	fs           CacheFS
	mutex        *MutexMap
	index        *lru.Cache
	prevOldestId string
	prevLength   int
	Compressed   bool
}

type gobIndex struct {
	Compressed bool
	Items      []Item
}

const indexPath = "zcache_index.gob"

func (c *Cache) onEvict(key interface{}, value interface{}) {
	id := key.(string)
	c.mutex.Lock(id)
	defer c.mutex.Unlock(id)
	c.fs.Remove(id)
}

func (c *Cache) Init() (err error) {
	f, err := c.fs.Open(indexPath)
	if err == nil {
		index := gobIndex{}
		defer f.Close()
		dec := gob.NewDecoder(f)
		if err = dec.Decode(&index); err != nil {
			return
		}
		if index.Compressed != c.Compressed {
			return errors.New("compression state changed. cache is invalid")
		}
		for _, item := range index.Items {
			c.index.Add(item.Id, item)
		}
		oldestKey, _, exists := c.index.GetOldest()
		if exists {
			c.prevOldestId = oldestKey.(string)
			c.prevLength = c.index.Len()
		}
		fmt.Printf("loaded %d existing items\n", c.index.Len())
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	go c.loop()
	return
}

func (c *Cache) loop() {
	t := time.NewTicker(time.Second * 10)
	for range t.C {
		if err := c.sync(); err != nil {
			fmt.Println(err)
		}
	}
}

func (c *Cache) sync() (err error) {
	oldestKey, _, exists := c.index.GetOldest()
	oldestId := ""
	if exists {
		oldestId = oldestKey.(string)
	}
	l := c.index.Len()
	if exists && (c.prevOldestId != oldestId || c.prevLength != l) {
		index := gobIndex{Compressed: c.Compressed}
		for _, key := range c.index.Keys() {
			v, exists := c.index.Get(key)
			if exists {
				index.Items = append(index.Items, v.(Item))
			}
		}
		fmt.Printf("syncing %d items\n", len(index.Items))
		f, err := c.fs.Create(indexPath)
		if err != nil {
			return err
		}
		defer f.Close()
		enc := gob.NewEncoder(f)
		c.prevOldestId = oldestId
		c.prevLength = l
		return enc.Encode(index)
	}
	return
}

func (c *Cache) Contains(id string) (exists bool) {
	c.mutex.Lock(id)
	defer c.mutex.Unlock(id)
	return c.index.Contains(id)
}

func (c *Cache) Get(id string) (item Item, exists bool) {
	c.mutex.Lock(id)
	defer c.mutex.Unlock(id)
	v, exists := c.index.Get(id)
	item = v.(Item)
	item.AccessedAt = time.Now()
	return
}

func (c *Cache) slugId(id string) string {
	id = zpath.FileEscape(id)
	if c.Compressed {
		id += ".xflate"
	}
	return id
}

func (c *Cache) open(id string) (file io.ReadSeekCloser, item Item, err error) {
	if v, exists := c.index.Get(id); exists {
		id = c.slugId(id)
		file, err = c.fs.Open(id)
		if err != nil {
			return
		}
		if c.Compressed {
			file, err = xflate.NewReader(file, nil)
			if err != nil {
				return
			}
		}
		item := v.(Item)
		return file, item, err
	}
	err = os.ErrNotExist
	return
}

func (c *Cache) Open(id string) (file io.ReadSeekCloser, item Item, err error) {
	c.mutex.Lock(id)
	defer c.mutex.Unlock(id)
	return c.open(id)
}

func (c *Cache) OpenOrCreate(id string, create func(w io.Writer) error) (res io.ReadSeekCloser, item Item, err error) {
	c.mutex.Lock(id)
	defer c.mutex.Unlock(id)
	fileId := c.slugId(id)
	if v, exists := c.index.Get(id); exists {
		item.AccessedAt = time.Now()
		res, err = c.fs.Open(fileId)
		item = v.(Item)
		return res, item, err
	}
	tmpId := fileId + ".tmp"
	file, source, err := c.create(tmpId)
	if err != nil {
		return
	}
	if err = create(file); err != nil {
		zio.CloseAll(file, source)
		c.fs.Remove(tmpId)
		return
	}
	if err = zio.CloseAll(file, source); err != nil {
		return
	}
	if err = c.fs.Rename(tmpId, fileId); err != nil {
		return
	}
	item.Id = id
	item.CreatedAt = time.Now()
	item.AccessedAt = time.Now()
	c.index.Add(id, item)
	return c.open(id)
}

func (c *Cache) create(id string) (file io.WriteCloser, source io.Closer, err error) {
	if err = c.fs.MkdirAll(filepath.Dir(id), os.ModePerm); err != nil {
		return
	}
	file, err = c.fs.Create(id)
	if err != nil {
		return
	}
	if c.Compressed {
		source = file
		file, err = xflate.NewWriter(file, &xflate.WriterConfig{
			Level: xflate.BestSpeed,
		})
	}
	return
}

func (c *Cache) Remove(id string) error {
	c.mutex.Lock(id)
	defer c.mutex.Unlock(id)
	c.index.Remove(id)
	return c.fs.Remove(id)
}
