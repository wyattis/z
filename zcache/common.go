package zcache

import "time"

type Item struct {
	Id         string
	AccessedAt time.Time
	CreatedAt  time.Time
}

type gobIndex struct {
	Compressed bool
	Items      []Item
}

const indexPath = "zcache_index.gob"
