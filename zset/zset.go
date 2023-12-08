package zset

import "sync"

var _ ISet[int] = New[int]()

type ISet[T comparable] interface {
	Add(items ...T) *Set[T]
	Delete(items ...T) *Set[T]
	Clear() *Set[T]
	Contains(items ...T) bool
	ContainsAny(items ...T) bool
	Equal(other *Set[T]) bool
	Size() int
	Items() []T
	Union(others ...*Set[T]) *Set[T]
	Complement(others ...*Set[T]) *Set[T]
	Clone() *Set[T]
	Intersection(others ...*Set[T]) *Set[T]
	Filter(f func(T) bool) *Set[T]
	FilterItems(f func(T) bool) []T
}

func New[T comparable](vals ...T) *Set[T] {
	s := &Set[T]{
		items: make(map[T]bool),
	}
	s.Add(vals...)
	return s
}

// Create a new set that is the union of all provided sets
func NewUnion[T comparable](sets ...*Set[T]) (s *Set[T]) {
	s = &Set[T]{}
	s.Union(sets...)
	return
}

type Set[T comparable] struct {
	items map[T]bool
	lock  sync.RWMutex
}

// Efficiently set the internal map for the set. Must hold the lock before calling.
func (s *Set[T]) setItems(m map[T]bool) {
	s.items = make(map[T]bool, len(m))
	for k, v := range m {
		s.items[k] = v
	}
}

func (s *Set[T]) Add(items ...T) *Set[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, item := range items {
		s.items[item] = true
	}
	return s
}

func (s *Set[T]) Delete(items ...T) *Set[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, item := range items {
		delete(s.items, item)
	}
	return s
}

func (s *Set[T]) Clear() *Set[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[T]bool)
	return s
}

func (s *Set[T]) Contains(items ...T) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

func (s *Set[T]) ContainsAny(items ...T) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

// Two sets are considered equal if they contain exactly the same elements.
func (s *Set[T]) Equal(other *Set[T]) bool {
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

func (s *Set[T]) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

func (s *Set[T]) Items() (res []T) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for item := range s.items {
		res = append(res, item)
	}
	return
}

// Union will add all items from the other sets to this set. This mutates the set.
func (s *Set[T]) Union(others ...*Set[T]) *Set[T] {
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

// Complement will remove all items from this set that are present in the other sets. This mutates the set.
func (s *Set[T]) Complement(others ...*Set[T]) *Set[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, b := range others {
		b.lock.RLock()
		defer b.lock.RUnlock()
		for key := range b.items {
			delete(s.items, key)
		}
	}
	return s
}

// Clone will create a new set with the same items as this set. This
func (s *Set[T]) Clone() *Set[T] {
	c := &Set[T]{}
	c.setItems(s.items)
	return c
}

// Intersection will reduce this set to the items that are present in this set and the other sets. This mutates the set.
func (s *Set[T]) Intersection(others ...*Set[T]) *Set[T] {
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

// Filter will reduce this set to the items that match the provided filter function. This mutates the set.
func (s *Set[T]) Filter(f func(T) bool) *Set[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	for key := range s.items {
		if !f(key) {
			delete(s.items, key)
		}
	}
	return s
}

// FilterItems will return a slice of items that match the provided filter function.
func (s *Set[T]) FilterItems(f func(T) bool) (res []T) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for key := range s.items {
		if f(key) {
			res = append(res, key)
		}
	}
	return
}
