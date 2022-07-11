package zcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidColumn = errors.New("no column with that name exists")
	ErrNoData        = errors.New("no data exists for that column")
)

func NewReader(source io.Reader, headers []string) (reader *CsvReader) {
	r := csv.NewReader(source)
	reader = &CsvReader{
		Reader: r,
	}
	if headers != nil && len(headers) > 0 {
		reader.headers = make(map[string]int, len(headers))
		for i := range headers {
			reader.headers[headers[i]] = i
		}
	}
	return reader
}

type CsvReader struct {
	*csv.Reader
	headers map[string]int
}

func (r *CsvReader) init() (err error) {
	if r.headers == nil {
		headers, err := r.Reader.Read()
		if err != nil {
			return err
		}
		r.headers = make(map[string]int, len(headers))
		for i := range headers {
			r.headers[headers[i]] = i
		}
	}
	return
}

func (r *CsvReader) Headers() (cols []string, err error) {
	if err = r.init(); err != nil {
		return
	}
	for col := range r.headers {
		cols = append(cols, col)
	}
	return
}

func (r *CsvReader) ReadAll() (lines []Line, err error) {
	for {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return lines, err
		}
		lines = append(lines, line)
	}
}

func (r *CsvReader) Read() (line Line, err error) {
	if err = r.init(); err != nil {
		return
	}
	row, err := r.Reader.Read()
	if err != nil {
		return
	}
	line.data = row
	line.headers = &r.headers
	line.errOnInvalid = r.FieldsPerRecord >= 0
	return
}

type Line struct {
	data         []string
	errOnInvalid bool
	headers      *map[string]int
}

func (l Line) Data() []string {
	return l.data
}

func (l Line) Get(key string) (val string, err error) {
	i, ok := (*l.headers)[key]
	if !ok {
		err = fmt.Errorf("no column exists with name '%s'", key)
		return
	}
	if i < len(l.data) {
		val = l.data[i]
	} else if l.errOnInvalid {
		err = ErrNoData
	}
	return
}

func (l Line) MustGet(key string) string {
	val, err := l.Get(key)
	if err != nil {
		panic(err)
	}
	return val
}
