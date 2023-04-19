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
}

func (i *CsvIterator) Next() (item Line, done bool, err error) {
	item, err = i.reader.Read()
	done = err == io.EOF
	if done {
		err = nil
		return
	}

	return
}
