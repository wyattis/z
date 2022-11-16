// Code generated by z. DO NOT EDIT.

package zintset

func New(items... int) (s *Set) {
	s = &Set{
		items: make(map[int]bool),
	}
  s.Add(items...)
  return
}

type Set struct {
	items map[int]bool
}

func (s *Set) Add(items ...int) {
	for _, item := range items {
		s.items[item] = true
	}
}

func (s Set) Contains(item int) bool {
  _, exists := s.items[item]
	return exists
}

func (s Set) ContainsAll(items ...int) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

func (s Set) ContainsAny(items ...int) bool {
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

func (s *Set) Delete(items ...int) {
	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *Set) Clear() {
	s.items = make(map[int]bool)
}

func (s *Set) Size() int {
	return len(s.items)
}

func (s *Set) Items() (res []int) {
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

func (s *Set) Intersection(others ...Set) *Set {
	res := s.Clone()
	res.Union(others...)
	for _, v := range res.Items() {
		for _, s := range others {
      if _, ok := s.items[v]; !ok {
        delete(res.items, v)
        break
      }
		}
	}
	return res
}
