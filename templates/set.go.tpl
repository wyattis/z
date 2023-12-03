{{- define "set" -}}

package {{.PackageName}}

// Create a new set with the given items
func New(items... {{.Type}}) (s *Set) {
	s = &Set{
		items: make(map[{{.Type}}]bool),
	}
  s.Add(items...)
  return
}

// Create a new set that is the union of all provided sets
func NewUnion(sets ...Set) (s *Set) {
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
func NewIntersection(sets ...Set) *Set {
  s := &Set{}
  s.Intersection(sets...)
  return s
}

type Set struct {
	items map[{{.Type}}]bool
}

// Efficiently set the internal map for the set
func (s *Set) setItems(m map[{{.Type}}]bool) {
  s.items = make(map[{{.Type}}]bool, len(m))
  for k, v := range m {
    s.items[k] = v
  }
}

// Add items to the set
func (s *Set) Add(items ...{{ .Type }}) *Set {
	for _, item := range items {
		s.items[item] = true
	}
  return s
}

// Check if the set contains the given item
func (s Set) Contains(item {{ .Type }}) bool {
  _, exists := s.items[item]
	return exists
}

// Check if the set contains all the given items
func (s Set) ContainsAll(items ...{{ .Type }}) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

// Check if the set contains any of the given items
func (s Set) ContainsAny(items ...{{ .Type }}) bool {
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

// Delete items from the set
func (s *Set) Delete(items ...{{ .Type }}) *Set {
	for _, item := range items {
		delete(s.items, item)
	}
  return s
}

// Remove all items from the set
func (s *Set) Clear() *Set {
	s.items = make(map[{{.Type}}]bool)
  return s
}

// Size returns the size of the set
func (s *Set) Size() int {
	return len(s.items)
}

// Items returns a slice with the items of the set
func (s *Set) Items() (res []{{.Type}}) {
	for key := range s.items {
		res = append(res, key)
	}
	return
}

// Union adds all the items of the other sets to this set. This mutates the set.
func (s *Set) Union(others ...Set) *Set {
	for _, b := range others {
		for key := range b.items {
			s.items[key] = true
		}
	}
  return s
}

// Complement removes items that are not in the other sets. This mutates the set.
func (s *Set) Complement(others ...Set) *Set  {
	for _, b := range others {
		for key := range b.items {
			delete(s.items, key)
		}
	}
  return s
}

// Return a new set that contains the same items as the original
func (s *Set) Clone() *Set {
	res := New()
  res.setItems(s.items)
	return res
}

// Intersection will reduce this set to the items that are present in this set and the other sets. This mutates the set.
func (s *Set) Intersection(others ...Set) *Set  {
  otherUnion := NewUnion(others...)
	for key := range s.items {
    if _, ok := otherUnion.items[key]; !ok {
      delete(s.items, key)
    }
  }
  return s
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
{{end -}}