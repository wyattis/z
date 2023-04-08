package ziter

type TransformFunc[F any, T any] func(item F) (res T, skip bool, err error)

type Transformer[F any, T any] interface {
	Transform(item F) (res T, skip bool, err error)
}

type TransformerIterator[F any, T any] struct {
	iterator  Iterator[F]
	transform TransformFunc[F, T]
}

func (t *TransformerIterator[F, T]) Next() (item T, done bool, err error) {
	fitem, done, err := t.iterator.Next()
	if done || err != nil {
		return
	}
	item, skip, err := t.transform(fitem)
	if err != nil {
		return
	}
	if skip {
		return t.Next()
	}
	return
}

func TransformIterator[F any, T any](iterator Iterator[F], t TransformFunc[F, T]) Iterator[T] {
	return &TransformerIterator[F, T]{
		iterator:  iterator,
		transform: t,
	}
}
