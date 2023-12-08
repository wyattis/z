package ziter

import (
	"io"
	"math"
)

type ChunkHandler[T any] func(chunk []T) error

// Run the handler on each chunk with a maximum size of chunkSize. The last chunk may be smaller than the others.
func WithChunksOfSize[T any](items []T, chunkSize int, handler ChunkHandler[T]) (err error) {
	numChunks := int(math.Ceil(float64(len(items)) / float64(chunkSize)))
	return withChunks(items, chunkSize, numChunks, handler)
}

// Run the handler N times with chunks of approximately equal size. The last chunk may be smaller than the others.
func WithNChunks[T any](items []T, numChunks int, handler ChunkHandler[T]) (err error) {
	chunkSize := int(math.Ceil(float64(len(items)) / float64(numChunks)))
	return withChunks(items, chunkSize, numChunks, handler)
}

// Break a slice into multiple slices of the provided size. The last chunk may be smaller than the others.
func SliceChunksOfSize[T any](items []T, chunkSize int) (chunks [][]T) {
	WithChunksOfSize(items, chunkSize, func(chunk []T) error {
		chunks = append(chunks, chunk)
		return nil
	})
	return
}

// Break a slice into N slices of approximately equal size. The last chunk may be smaller than the others.
func SliceNChunks[T any](items []T, numChunks int) (chunks [][]T) {
	WithNChunks(items, numChunks, func(chunk []T) error {
		chunks = append(chunks, chunk)
		return nil
	})
	return
}

func withChunks[T any](items []T, chunkSize, numChunks int, handler ChunkHandler[T]) (err error) {
	for i := 0; i < numChunks; i++ {
		end := (i + 1) * chunkSize
		start := i * chunkSize
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

type chunkIterator[T any] struct {
	iterator Iterator[T]
	size     int
	chunk    []T
	err      error
}

func (c *chunkIterator[T]) Next() (hasMore bool) {
	c.chunk = make([]T, 0, c.size)
	for c.iterator.Next() && len(c.chunk) < c.size {
		if c.iterator.Err() != nil {
			c.err = c.iterator.Err()
			return c.err != io.EOF
		}
		c.chunk = append(c.chunk, c.iterator.Item())
	}
	return len(c.chunk) > 0
}

func (c *chunkIterator[T]) Err() error {
	return c.err
}

func (c *chunkIterator[T]) Item() []T {
	return c.chunk
}

// Transforms an iterator of T into an iterator of []T, where each chunk has a max size of chunkSize.
func Chunkify[T any](iterator Iterator[T], chunkSize int) Iterator[[]T] {
	return &chunkIterator[T]{
		iterator: iterator,
		size:     chunkSize,
	}
}
