package zcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var (
	ErrInvalidColumn = errors.New("no column with that name exists")
	ErrNoData        = errors.New("no data exists for that column")
)

func NewReader(source io.Reader, headers []string) (reader *CsvReader) {
	r := csv.NewReader(source)
	// r.ReuseRecord = true
	reader = &CsvReader{
		Reader: r,
	}
	if headers != nil && len(headers) > 0 {
		reader.headers = make(map[string]int)
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

func (r *CsvReader) Headers() (cols []string) {
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
	if r.headers == nil {
		r.headers = make(map[string]int)
		headers, err := r.Reader.Read()
		if err != nil {
			return line, err
		}
		for i := range headers {
			r.headers[headers[i]] = i
		}
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

func (r *CsvReader) scanLine(dest reflect.Value, line Line) (err error) {

	return
}

// Scan a row into this destination
func (r *CsvReader) Scan(dest *any) (err error) {
	t := reflect.TypeOf(dest)
	// v := reflect.ValueOf(dest)
	if t.Kind() != reflect.Pointer {
		return errors.New("CsvReader.Scan requires a pointer destination")
	}
	// v := reflect.Indirect(reflect.ValueOf(dest))
	return
}

// Scan all rows into a slice
func (r *CsvReader) ScanAll(dest any) (err error) {
	vPointer := reflect.ValueOf(dest)
	if reflect.Indirect(vPointer).Type().Kind() != reflect.Slice {
		return errors.New("CsvReader.ScanAll requires a pointer to a slice")
	}
	rDest := vPointer.Elem()
	for {
		line, err := r.Read()
		if err != nil {
			return err
		}
		rowVal := reflect.ValueOf(vPointer.Type().Elem())
		fmt.Println(rDest, vPointer.Type(), rowVal)
		if err = r.scanLine(rowVal, line); err != nil {
			return err
		}
		rDest.Set(reflect.Append(rDest, rowVal))
	}
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
