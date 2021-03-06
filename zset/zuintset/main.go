// Code generated by z. DO NOT EDIT.

package zuintset

func New(items... uint) (s *Set) {
	s = &Set{
		items: make(map[uint]bool),
	}
  s.Add(items...)
  return
}

type Set struct {
	items map[uint]bool
}

func (s *Set) Add(items ...uint) {
	for _, item := range items {
		s.items[item] = true
	}
}

func (s Set) Contains(item uint) bool {
  _, exists := s.items[item]
	return exists
}

func (s Set) ContainsAll(items ...uint) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

func (s Set) ContainsAny(items ...uint) bool {
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

func (s *Set) Delete(items ...uint) {
	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *Set) Clear() {
	s.items = make(map[uint]bool)
}

func (s *Set) Size() int {
	return len(s.items)
}

func (s *Set) Items() (res []uint) {
	for key := range s.items {
		res = append(res, key)
	}
	return
}

func (s *Set) Union(others ...Set) {
	for _, b := range others {
		for key := range b.items {
			s.items[key] = true
		}
	}
}

func (s *Set) Complement(others ...Set) {
	for _, b := range others {
		for key := range b.items {
			delete(s.items, key)
		}
	}
}

func (s *Set) Clone() *Set {
	res := New()
	res.Add(s.Items()...)
	return res
}
