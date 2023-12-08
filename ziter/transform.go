package ziter

import "io"

type IFilter[T any] interface {
	Filter(item T) (keep bool, err error)
}

type ITransform[F any, T any] interface {
	Transform(item F) (res T, err error)
}

type filteringIterator[T any] struct {
	iterator Iterator[T]
	filter   IFilter[T]

	item T
	err  error
}

func (f *filteringIterator[T]) Next() (hasItem bool) {
	for f.iterator.Next() {
		if f.err = f.iterator.Err(); f.err != nil {
			return f.err != io.EOF
		}
		f.item = f.iterator.Item()
		keep, err := f.filter.Filter(f.item)
		if err != nil {
			f.err = err
			return true
		}
		if !keep {
			continue
		}
		return true
	}
	return false
}

func (f *filteringIterator[T]) Err() error {
	return f.err
}

func (f *filteringIterator[T]) Item() T {
	return f.item
}

type transformingIterator[F any, T any] struct {
	iterator    Iterator[F]
	transformer ITransform[F, T]

	item T
	err  error
}

// Iterate through the underlying iterator items until an unfiltered item is found, and transforming the items that
// do match the filter.
func (t *transformingIterator[F, T]) Next() (hasItem bool) {
	if !t.iterator.Next() {
		return false
	}
	if t.err = t.iterator.Err(); t.err != nil {
		return true
	}
	t.item, t.err = t.transformer.Transform(t.iterator.Item())
	return true
}

func (t *transformingIterator[F, T]) Err() error {
	return t.err
}

func (t *transformingIterator[F, T]) Item() T {
	return t.item
}

func Transform[F any, T any](iterator Iterator[F], transformer ITransform[F, T]) Iterator[T] {
	return &transformingIterator[F, T]{
		iterator:    iterator,
		transformer: transformer,
	}
}

func Filter[T any](iterator Iterator[T], filter IFilter[T]) Iterator[T] {
	return &filteringIterator[T]{
		iterator: iterator,
		filter:   filter,
	}
}
