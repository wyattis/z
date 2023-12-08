package ziter

type Iterator[T any] interface {
	Next() (stillOpen bool)
	Err() (err error)
	Item() (item T)
}

func SliceIterator[T any](items []T) *sliceIterator[T] {
	return &sliceIterator[T]{
		items: items,
		index: -1,
	}
}

type sliceIterator[T any] struct {
	items []T
	index int
}

func (s *sliceIterator[T]) Next() (hasMore bool) {
	s.index++
	return s.index < len(s.items)
}

func (s *sliceIterator[T]) Err() error {
	return nil
}

func (s *sliceIterator[T]) Item() T {
	return s.items[s.index]
}
