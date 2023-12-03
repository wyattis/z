package zset

func New[T comparable](vals ...T) *Set[T] {
	s := &Set[T]{
		items: make(map[T]bool),
	}
	s.Add(vals...)
	return s
}

// Create a new set that is the union of all provided sets
func NewUnion[T comparable](sets ...Set[T]) (s *Set[T]) {
	s = &Set[T]{}
	s.Union(sets...)
	return
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

// Intersection will reduce this set to the items that are present in this set and the other sets
func (s *Set[T]) Intersection(others ...Set[T]) {
	otherUnion := NewUnion(others...)
	for key := range s.items {
		if _, ok := otherUnion.items[key]; !ok {
			delete(s.items, key)
		}
	}
}

// Two sets are considered equal if they contain exactly the same elements.
func (s *Set[T]) Equal(other Set[T]) bool {
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
