package zset

func NewHashSet[T any, K comparable](hash Hasher[T, K], vals ...T) *HashSet[T, K] {
	s := &HashSet[T, K]{
		hasher: hash,
		items:  make(map[K]T),
	}
	s.Add(vals...)
	return s
}

type HashSet[T any, K comparable] struct {
	hasher func(v T) K
	items  map[K]T
}

func (s *HashSet[T, K]) Add(items ...T) {
	for _, item := range items {
		k := s.hasher(item)
		s.items[k] = item
	}
}

func (s HashSet[T, K]) Contains(item T) bool {
	k := s.hasher(item)
	_, exists := s.items[k]
	return exists
}

func (s HashSet[T, K]) ContainsAll(items ...T) bool {
	for _, item := range items {
		k := s.hasher(item)
		if _, exists := s.items[k]; !exists {
			return false
		}
	}
	return true
}

func (s HashSet[T, K]) ContainsAny(items ...T) bool {
	for _, item := range items {
		k := s.hasher(item)
		if _, exists := s.items[k]; exists {
			return true
		}
	}
	return false
}

func (s *HashSet[T, K]) Delete(items ...T) {
	for _, item := range items {
		k := s.hasher(item)
		delete(s.items, k)
	}
}

func (s *HashSet[T, K]) Clear() {
	s.items = make(map[K]T)
}

func (s *HashSet[T, K]) Size() int {
	return len(s.items)
}

func (s *HashSet[T, K]) Items() (res []T) {
	for _, item := range s.items {
		res = append(res, item)
	}
	return
}

func (s *HashSet[T, K]) Union(others ...HashSet[T, K]) {
	for _, b := range others {
		for key, item := range b.items {
			s.items[key] = item
		}
	}
}

func (s *HashSet[T, K]) Complement(others ...HashSet[T, K]) {
	for _, b := range others {
		for key := range b.items {
			delete(s.items, key)
		}
	}
}

func (s *HashSet[T, K]) Clone() *HashSet[T, K] {
	return NewHashSet(s.hasher, s.Items()...)
}

func (s *HashSet[T, K]) Intersection(others ...HashSet[T, K]) *HashSet[T, K] {
	res := s.Clone()
	res.Union(others...)
	for _, v := range res.Items() {
		for _, s := range others {
			k := s.hasher(v)
			if _, ok := s.items[k]; !ok {
				delete(res.items, k)
				break
			}
		}
	}
	return res
}
