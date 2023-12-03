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

// Efficiently set the internal map for the set
func (s *Set[T]) setItems(m map[T]bool) {
	s.items = make(map[T]bool, len(m))
	for k, v := range m {
		s.items[k] = v
	}
}

func (s *Set[T]) Add(items ...T) *Set[T] {
	for _, item := range items {
		s.items[item] = true
	}
	return s
}

func (s *Set[T]) Delete(items ...T) *Set[T] {
	for _, item := range items {
		delete(s.items, item)
	}
	return s
}

func (s *Set[T]) Clear() *Set[T] {
	s.items = make(map[T]bool)
	return s
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

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) Items() (res []T) {
	for item := range s.items {
		res = append(res, item)
	}
	return
}

// Union will add all items from the other sets to this set. This mutates the set.
func (s *Set[T]) Union(others ...Set[T]) *Set[T] {
	for _, b := range others {
		for key, item := range b.items {
			s.items[key] = item
		}
	}
	return s
}

// Complement will remove all items from this set that are present in the other sets. This mutates the set.
func (s *Set[T]) Complement(others ...Set[T]) *Set[T] {
	for _, b := range others {
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
func (s *Set[T]) Intersection(others ...Set[T]) *Set[T] {
	otherUnion := NewUnion(others...)
	for key := range s.items {
		if _, ok := otherUnion.items[key]; !ok {
			delete(s.items, key)
		}
	}
	return s
}
