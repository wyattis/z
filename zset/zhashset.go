package zset

import "sync"

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
func NewHashSetUnion[T any, K comparable](sets ...*HashSet[T, K]) (s *HashSet[T, K]) {
	s = &HashSet[T, K]{}
	s.Union(sets...)
	return
}

type HashSet[T any, K comparable] struct {
	hasher func(v T) K
	items  map[K]T
	lock   sync.RWMutex
}

func (s *HashSet[T, K]) Add(items ...T) *HashSet[T, K] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, item := range items {
		k := s.hasher(item)
		s.items[k] = item
	}
	return s
}

func (s *HashSet[T, K]) Delete(items ...T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, item := range items {
		k := s.hasher(item)
		delete(s.items, k)
	}
}

func (s *HashSet[T, K]) Clear() *HashSet[T, K] {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[K]T)
	return s
}

func (s HashSet[T, K]) Contains(items ...T) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, item := range items {
		k := s.hasher(item)
		if _, exists := s.items[k]; !exists {
			return false
		}
	}
	return true
}

func (s HashSet[T, K]) ContainsAny(items ...T) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, item := range items {
		k := s.hasher(item)
		if _, exists := s.items[k]; exists {
			return true
		}
	}
	return false
}

// Two sets are considered equal if they contain exactly the same elements.
func (s *HashSet[T, K]) Equal(other *HashSet[T, K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	other.lock.RLock()
	defer other.lock.RUnlock()
	if len(s.items) != len(other.items) {
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
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

func (s *HashSet[T, K]) Items() (res []T) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, item := range s.items {
		res = append(res, item)
	}
	return
}

func (s *HashSet[T, K]) Union(others ...*HashSet[T, K]) *HashSet[T, K] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, b := range others {
		b.lock.RLock()
		defer b.lock.RUnlock()
		for key, item := range b.items {
			s.items[key] = item
		}
	}
	return s
}

func (s *HashSet[T, K]) Clone() *HashSet[T, K] {
	return NewHashSet(s.hasher, s.Items()...)
}

// Intersection will reduce this set to the items that are present in this set and the other sets
func (s *HashSet[T, K]) Intersection(others ...*HashSet[T, K]) *HashSet[T, K] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, o := range others {
		o.lock.RLock()
		defer o.lock.RUnlock()
	}
	for key := range s.items {
		for _, other := range others {
			if _, ok := other.items[key]; !ok {
				delete(s.items, key)
				break
			}
		}
	}
	return s
}

// Difference will reduce this set to the items that are present in this set but not in the other sets. This mutates the set.
func (s *HashSet[T, K]) Difference(others ...*HashSet[T, K]) *HashSet[T, K] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, o := range others {
		o.lock.RLock()
		defer o.lock.RUnlock()
	}
	for key := range s.items {
		for _, other := range others {
			if _, ok := other.items[key]; ok {
				delete(s.items, key)
				break
			}
		}
	}
	return s
}

// Filter will reduce this set to the items that match the provided filter function. This mutates the set.
func (s *HashSet[T, K]) Filter(keep func(item T) bool) *HashSet[T, K] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for key, item := range s.items {
		if !keep(item) {
			delete(s.items, key)
		}
	}
	return s
}

// FilterItems will reduce this set to the items that match the provided filter function. This mutates the set.
func (s *HashSet[T, K]) FilterItems(keep func(item T) bool) (res []T) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, item := range s.items {
		if keep(item) {
			res = append(res, item)
		}
	}
	return
}
