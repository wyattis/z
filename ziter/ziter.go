package ziter

import (
	"errors"
	"math"
	"reflect"
)

var (
	ErrValueNotSlice = errors.New("value is not a slice")
)

type ChunkHandler[T any] func(chunk []T) error

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
func SliceChunks[T any](items []T, size int, handler ChunkHandler[T]) (err error) {
	numChunks := int(math.Ceil(float64(len(items)) / float64(size)))
	for i := 0; i < numChunks; i++ {
		end := (i + 1) * size
		start := i * size
		if end > len(items) {
			end = len(items)
		}
		chunk := items[start:end]
		if err = handler(chunk); err != nil {
			return
		}
	}
	return
}

func IteratorChunk[T any](iterator Iterator[T], size int, handler ChunkHandler[T]) (err error) {
	chunk := make([]T, 0, size)
	for iterator.Next() {
		if err = iterator.Err(); err != nil {
			return
		}
		chunk = append(chunk, iterator.Item())
		if len(chunk) >= size {
			if err = handler(chunk); err != nil {
				return
			}
			chunk = chunk[:0]
		}
	}
	if len(chunk) > 0 {
		err = handler(chunk)
	}
	return
}
