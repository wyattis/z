// Code generated by z. DO NOT EDIT.

package zstringset

func New(items... string) (s *Set) {
	s = &Set{
		items: make(map[string]bool),
	}
  s.Add(items...)
  return
}

type Set struct {
	items map[string]bool
}

func (s *Set) Add(items ...string) {
	for _, item := range items {
		s.items[item] = true
	}
}

func (s Set) Contains(item string) bool {
  _, exists := s.items[item]
	return exists
}

func (s Set) ContainsAll(items ...string) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

func (s Set) ContainsAny(items ...string) bool {
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

func (s *Set) Delete(items ...string) {
	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *Set) Clear() {
	s.items = make(map[string]bool)
}

func (s *Set) Size() int {
	return len(s.items)
}

func (s *Set) Items() (res []string) {
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
