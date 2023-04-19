package zcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/wyattis/z/ztime"
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

func (l Line) GetFloat32(key string) (val float32, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	v64, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return
	}
	val = float32(v64)
	return
}

func (l Line) GetFloat64(key string) (val float64, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	return strconv.ParseFloat(str, 64)
}

func (l Line) GetInt(key string) (val int, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	return strconv.Atoi(str)
}

func (l Line) GetInt32(key string) (val int32, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	v64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return
	}
	val = int32(v64)
	return
}

func (l Line) GetInt64(key string) (val int64, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	return strconv.ParseInt(str, 10, 64)
}

func (l Line) GetBool(key string) (val bool, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	return strconv.ParseBool(str)
}

func (l Line) GetUint(key string) (val uint, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	v64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return
	}
	val = uint(v64)
	return
}

func (l Line) GetUint32(key string) (val uint32, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	v64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return
	}
	val = uint32(v64)
	return
}

func (l Line) GetUint64(key string) (val uint64, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	return strconv.ParseUint(str, 10, 64)
}

func (l Line) GetTime(key string, layouts ...string) (val time.Time, err error) {
	str, err := l.Get(key)
	if err != nil {
		return
	}
	return ztime.Parse(str, layouts...)
}
