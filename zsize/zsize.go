package zsize

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wyattis/z/zstring"
)

type Size int64

const (
	Byte     Size = 1
	Kilobyte Size = 1000
	Kibibyte Size = 1024
	Megabyte Size = Kilobyte * Kilobyte
	Mebibyte Size = Kibibyte * Kibibyte
	Gigabyte Size = Megabyte * Kilobyte
	GibiByte Size = Mebibyte * Kibibyte
	Terabyte Size = Gigabyte * Kilobyte
	Tebibyte Size = GibiByte * Kibibyte
	Petabyte Size = Terabyte * Kilobyte
	Pebibyte Size = Tebibyte * Kibibyte
	Exabyte  Size = Petabyte * Kilobyte
)

var unitMap = map[string]Size{
	"KB":  Kilobyte,
	"KiB": Kibibyte,
	"MB":  Megabyte,
	"MiB": Mebibyte,
	"GB":  Gigabyte,
	"GiB": GibiByte,
	"TB":  Terabyte,
	"TiB": Tebibyte,
	"PB":  Petabyte,
	"PiB": Pebibyte,
	"EB":  Exabyte,
}

// Convert a human readable string value into a Size
func Parse(val string) (res Size, err error) {
	val = strings.TrimSpace(val)
	val, units := zstring.CutAt(val, -2)
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return
	}
	if units == "" {
		res = Size(v)
		return
	}
	conv, ok := unitMap[units]
	if !ok {
		err = fmt.Errorf("unknown unit: %s", units)
		return
	}
	v *= int64(conv)
	res = Size(v)
	return
}

// Set the value of Size by parsing a string value
func (s *Size) Set(val string) (err error) {
	v, err := Parse(val)
	if err != nil {
		return
	}
	*s = v
	return
}

// Return the underlying integer value as a string
func (s Size) String() string {
	return fmt.Sprint(int64(s))
}

// Return the formatted string
func (s Size) Formatted() string {
	p := int64(s)
	units := "b"
	// TODO
	for u, size := range unitMap {
		r := int64(s % size)
		fmt.Println(u, p, r)
		if r == 0 {
			units = u
			break
		}
		p = r
	}
	return fmt.Sprintf("%d%s", p, units)
}
