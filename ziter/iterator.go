package ziter

type Iterator[T any] interface {
	Next() (item T, done bool, err error)
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

func (s *sliceIterator[T]) Next() (item T, done bool, err error) {
	done = s.index >= len(s.Slice)
	if done {
		return
	}
	item = s.Slice[s.index]
	s.index++
	return
}
