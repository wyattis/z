package ziter

type Iterator[T any] interface {
	Next() bool
	Item() T
	Err() error
}

func NewSliceIterator[T any](slice []T) *sliceIterator[T] {
	return &sliceIterator[T]{
		Slice: slice,
		index: -1,
	}
}

type sliceIterator[T any] struct {
	Slice []T
	index int
}

func (s *sliceIterator[T]) Next() bool {
	s.index++
	return s.index < len(s.Slice)
}

func (s *sliceIterator[T]) Item() T {
	return s.Slice[s.index]
}

func (s sliceIterator[T]) Err() error {
	return nil
}
