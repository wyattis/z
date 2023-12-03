// Code generated by z. DO NOT EDIT.

package zuintset

// Create a new set with the given items
func New(items... uint) (s *Set) {
	s = &Set{
		items: make(map[uint]bool),
	}
  s.Add(items...)
  return
}

// Create a new set that is the union of all provided sets
func NewUnion(sets ...Set) (s *Set) {
  s = &Set{}
  s.Union(sets...)
  return
}

// Create a new set that is the intersection of all provided sets
func NewIntersection(sets ...Set) *Set {
  s := &Set{}
  s.Intersection(sets...)
  return s
}

type Set struct {
	items map[uint]bool
}

// Efficiently set the internal map for the set
func (s *Set) setItems(m map[uint]bool) {
  s.items = make(map[uint]bool, len(m))
  for k, v := range m {
    s.items[k] = v
  }
}

// Add items to the set
func (s *Set) Add(items ...uint) {
	for _, item := range items {
		s.items[item] = true
	}
}

// Check if the set contains the given item
func (s Set) Contains(item uint) bool {
  _, exists := s.items[item]
	return exists
}

// Check if the set contains all the given items
func (s Set) ContainsAll(items ...uint) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

// Check if the set contains any of the given items
func (s Set) ContainsAny(items ...uint) bool {
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

// Delete items from the set
func (s *Set) Delete(items ...uint) {
	for _, item := range items {
		delete(s.items, item)
	}
}

// Remove all items from the set
func (s *Set) Clear() {
	s.items = make(map[uint]bool)
}

// Size returns the size of the set
func (s *Set) Size() int {
	return len(s.items)
}

// Items returns a slice with the items of the set
func (s *Set) Items() (res []uint) {
	for key := range s.items {
		res = append(res, key)
	}
	return
}

// Union adds all the items of the other sets to this set
func (s *Set) Union(others ...Set) {
	for _, b := range others {
		for key := range b.items {
			s.items[key] = true
		}
	}
}

// Complement removes items that are not in the other sets
func (s *Set) Complement(others ...Set)  {
	for _, b := range others {
		for key := range b.items {
			delete(s.items, key)
		}
	}
}

// Return a new set that contains the same items as the original
func (s *Set) Clone() *Set {
	res := New()
  res.setItems(s.items)
	return res
}

// Intersection will reduce this set to the items that are present in this set and the other sets
func (s *Set) Intersection(others ...Set)  {
  otherUnion := NewUnion(others...)
	for key := range s.items {
    if _, ok := otherUnion.items[key]; !ok {
      delete(s.items, key)
    }
  }
}

// Equal returns if boths sets contain the same items
func (s Set) Equal(other Set) bool {
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
