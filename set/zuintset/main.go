// Code generated by z. DO NOT EDIT.

package zuintset

func New() *UintSet {
	return &UintSet{
		items: make(map[uint]bool),
	}
}

type UintSet struct {
	items map[uint]bool
}

func (s *UintSet) Add(items ...uint) {
	for _, item := range items {
		s.items[item] = true
	}
}

func (s UintSet) Contains(items ...uint) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

func (s *UintSet) Delete(items ...uint) {
	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *UintSet) Clear() {
	s.items = make(map[uint]bool)
}

func (s *UintSet) Size() int {
	return len(s.items)
}

func (s *UintSet) Items() (res []uint) {
	for key := range s.items {
		res = append(res, key)
	}
	return
}

func (s *UintSet) Union(others ...UintSet) {
	for _, b := range others {
		for key := range b.items {
			s.items[key] = true
		}
	}
}

func (s *UintSet) Complement(others ...UintSet) {
	for _, b := range others {
		for key := range b.items {
			delete(s.items, key)
		}
	}
}

func (s *UintSet) Clone() *UintSet {
	res := New()
	res.Add(s.Items()...)
	return res
}
