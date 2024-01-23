{{- define "set" -}}

package {{.PackageName}}

import "sync"

// Create a new set with the given items
func New(items... {{.Type}}) (s *Set) {
	s = &Set{
		items: make(map[{{.Type}}]bool),
	}
  s.Add(items...)
  return
}

// Create a new set that is the union of all provided sets
func NewUnion(sets ...*Set) (s *Set) {
  size := 0
  for _, set := range sets {
    size += set.Size()
  }
  s = &Set{
    items: make(map[{{.Type}}]bool, size),
  }
  s.Union(sets...)
  return
}

// Create a new set that is the intersection of all provided sets
func NewIntersection(sets ...*Set) *Set {
  s := &Set{}
  s.Intersection(sets...)
  return s
}

type Set struct {
	items  map[{{.Type}}]bool
  lock   sync.RWMutex
}

// Efficiently set the internal map for the set. Must hold the lock before calling this.
func (s *Set) setItems(m map[{{.Type}}]bool) {
  s.items = make(map[{{.Type}}]bool, len(m))
  for k, v := range m {
    s.items[k] = v
  }
}

// Add items to the set
func (s *Set) Add(items ...{{ .Type }}) *Set {
  s.lock.Lock()
  defer s.lock.Unlock()
	for _, item := range items {
		s.items[item] = true
	}
  return s
}

// Check if the set contains all of the given items
func (s *Set) Contains(items ...{{ .Type }}) bool {
  s.lock.RLock()
  defer s.lock.RUnlock()
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

// Check if the set contains any of the given items
func (s *Set) ContainsAny(items ...{{ .Type }}) bool {
  s.lock.RLock()
  defer s.lock.RUnlock()
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

// Delete items from the set
func (s *Set) Delete(items ...{{ .Type }}) *Set {
  s.lock.Lock()
  defer s.lock.Unlock()
	for _, item := range items {
		delete(s.items, item)
	}
  return s
}

// Remove all items from the set
func (s *Set) Clear() *Set {
  s.lock.Lock()
  defer s.lock.Unlock()
	s.items = make(map[{{.Type}}]bool)
  return s
}

// Size returns the size of the set
func (s *Set) Size() int {
  s.lock.RLock()
  defer s.lock.RUnlock()
	return len(s.items)
}

// Items returns a slice with the items of the set
func (s *Set) Items() (res []{{.Type}}) {
  s.lock.RLock()
  defer s.lock.RUnlock()
	for key := range s.items {
		res = append(res, key)
	}
	return
}

// Union adds all the items of the other sets to this set. This mutates the set.
func (s *Set) Union(others ...*Set) *Set {
  s.lock.RLock()
  defer s.lock.RUnlock()
	for _, b := range others {
    b.lock.RLock()
    defer b.lock.RUnlock()
		for key := range b.items {
			s.items[key] = true
		}
	}
  return s
}

// Difference removes items that are in the other sets. This mutates the set.
func (s *Set) Difference(others ...*Set) *Set  {
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

// Return a new set that contains the same items as the original
func (s *Set) Clone() *Set {
  s.lock.RLock()
  defer s.lock.RUnlock()
	res := New()
  res.setItems(s.items)
	return res
}

// Intersection will reduce this set to the items that are present in this set and the other sets. This mutates the set.
func (s *Set) Intersection(others ...*Set) *Set  {
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

// Equal returns if boths sets contain the same items
func (s *Set) Equal(other *Set) bool {
  s.lock.RLock()
  defer s.lock.RUnlock()
  if len(s.items) != len(other.items) {
    return false
  }
  other.lock.RLock()
  defer other.lock.RUnlock()
  for key := range s.items {
    if _, exists := other.items[key]; !exists {
      return false
    }
  }
  return true
}

// Filter removes all items from the set that do not match the given filter function. This mutates the set.
func (s *Set) Filter(f func({{.Type}}) bool) *Set {
  s.lock.Lock()
  defer s.lock.Unlock()
  for key := range s.items {
    if !f(key) {
      delete(s.items, key)
    }
  }
  return s
}

// FilterItems returns a slice of items from the set that match the given filter function.
func (s *Set) FilterItems(f func({{.Type}}) bool) (res []{{.Type}}) {
  s.lock.RLock()
  defer s.lock.RUnlock()
  for key := range s.items {
    if f(key) {
      res = append(res, key)
    }
  }
  return
}
{{end -}}
