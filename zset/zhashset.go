package zset

type Hasher[T any, K comparable] func(T) K

func NewHashSet[T any, K comparable](hash Hasher[T, K], vals ...T) *HashSet[T, K] {
	s := &HashSet[T, K]{
		hasher: hash,
		items:  make(map[K]T),
	}
	s.Add(vals...)
	return s
}

// Create a new set that is the union of all provided sets
func NewHashSetUnion[T any, K comparable](sets ...HashSet[T, K]) (s *HashSet[T, K]) {
	s = &HashSet[T, K]{}
	s.Union(sets...)
	return
}

type HashSet[T any, K comparable] struct {
	hasher func(v T) K
	items  map[K]T
}

func (s *HashSet[T, K]) Add(items ...T) *HashSet[T, K] {
	for _, item := range items {
		k := s.hasher(item)
		s.items[k] = item
	}
	return s
}

func (s *HashSet[T, K]) Delete(items ...T) {
	for _, item := range items {
		k := s.hasher(item)
		delete(s.items, k)
	}
}

func (s *HashSet[T, K]) Clear() *HashSet[T, K] {
	s.items = make(map[K]T)
	return s
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

// Two sets are considered equal if they contain exactly the same elements.
func (s *HashSet[T, K]) Equal(other HashSet[T, K]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for key := range s.items {
		if _, exists := other.items[key]; !exists {
			return false
		}
	}
	return true
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

func (s *HashSet[T, K]) Union(others ...HashSet[T, K]) *HashSet[T, K] {
	for _, b := range others {
		for key, item := range b.items {
			s.items[key] = item
		}
	}
	return s
}

func (s *HashSet[T, K]) Complement(others ...HashSet[T, K]) *HashSet[T, K] {
	for _, b := range others {
		for key := range b.items {
			delete(s.items, key)
		}
	}
	return s
}

func (s *HashSet[T, K]) Clone() *HashSet[T, K] {
	return NewHashSet(s.hasher, s.Items()...)
}

// Intersection will reduce this set to the items that are present in this set and the other sets
func (s *HashSet[T, K]) Intersection(others ...HashSet[T, K]) *HashSet[T, K] {
	otherUnion := NewHashSetUnion(others...)
	for key := range s.items {
		if _, ok := otherUnion.items[key]; !ok {
			delete(s.items, key)
		}
	}
	return s
}
