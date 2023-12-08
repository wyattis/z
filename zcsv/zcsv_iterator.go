package zcsv

import (
	"io"

	"github.com/wyattis/z/ziter"
)

func NewIterator(reader *CsvReader) ziter.Iterator[Line] {
	return &CsvIterator{
		reader: reader,
	}
}

type CsvIterator struct {
	reader *CsvReader
	item   Line
	err    error
}

func (c *CsvIterator) Next() (hasMore bool) {
	c.item, c.err = c.reader.Read()
	hasMore = c.err == io.EOF
	return
}

func (c *CsvIterator) Err() error {
	return c.err
}

func (c *CsvIterator) Item() Line {
	return c.item
}
