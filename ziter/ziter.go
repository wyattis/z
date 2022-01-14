package ziter

import (
	"errors"
	"math"
	"reflect"
)

var (
	ErrValueNotSlice = errors.New("value is not a slice")
)

type ChunkHandler func(chunk []interface{}) error

// Convert interface{} into []interface{} if the input is a slice, otherwise,
// return an error.
func InterfaceToSlice(input interface{}) (slice []interface{}, err error) {
	rv := reflect.ValueOf(input)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			slice = append(slice, rv.Index(i).Interface())
		}
	} else {
		err = ErrValueNotSlice
	}
	return
}

// Break a slice into multiple slices of the provided size
func WithChunks(items interface{}, size int, handler ChunkHandler) error {
	vals, err := InterfaceToSlice(items)
	if err != nil {
		return err
	}
	numChunks := int(math.Ceil(float64(len(vals)) / float64(size)))
	for i := 0; i < numChunks; i++ {
		end := (i + 1) * size
		start := i * size
		if end > len(vals) {
			end = len(vals)
		}
		chunk := vals[start:end]
		if err := handler(chunk); err != nil {
			return err
		}
	}
	return nil
}
