package zset

package zset

type Hasher[T any, K comparable] func(T) K

func New[T comparable](vals ...T) *Set[T] {
	s := &Set[T]{
		items: make(map[T]bool),
	}
	s.Add(vals...)
	return s
}

type Set[T comparable] struct {
	items map[T]bool
}

func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.items[item] = true
	}
}

func (s Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

func (s Set[T]) ContainsAll(items ...T) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

func (s Set[T]) ContainsAny(items ...T) bool {
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

func (s *Set[T]) Delete(items ...T) {
	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *Set[T]) Clear() {
	s.items = make(map[T]bool)
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) Items() (res []T) {
	for item := range s.items {
		res = append(res, item)
	}
	return
}

func (s *Set[T]) Union(others ...Set[T]) {
	for _, b := range others {
		for key, item := range b.items {
			s.items[key] = item
		}
	}
}

func (s *Set[T]) Complement(others ...Set[T]) {
	for _, b := range others {
		for key := range b.items {
			delete(s.items, key)
		}
	}
}

func (s *Set[T]) Clone() *Set[T] {
	return New(s.Items()...)
}

func (s *Set[T]) Intersection(others ...Set[T]) *Set[T] {
	res := s.Clone()
	res.Union(others...)
	for _, v := range res.Items() {
		for _, s := range others {
			if _, ok := s.items[v]; !ok {
				delete(res.items, v)
				break
			}
		}
	}
	return res
}