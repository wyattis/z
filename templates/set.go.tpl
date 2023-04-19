{{- define "set" -}}

package {{.PackageName}}

func New(items... {{.Type}}) (s *Set) {
	s = &Set{
		items: make(map[{{.Type}}]bool),
	}
  s.Add(items...)
  return
}

type Set struct {
	items map[{{.Type}}]bool
}

func (s *Set) Add(items ...{{ .Type }}) {
	for _, item := range items {
		s.items[item] = true
	}
}

func (s Set) Contains(item {{ .Type }}) bool {
  _, exists := s.items[item]
	return exists
}

func (s Set) ContainsAll(items ...{{ .Type }}) bool {
	for _, item := range items {
		if _, exists := s.items[item]; !exists {
			return false
		}
	}
	return true
}

func (s Set) ContainsAny(items ...{{ .Type }}) bool {
	for _, item := range items {
		if _, exists := s.items[item]; exists {
			return true
		}
	}
	return false
}

func (s *Set) Delete(items ...{{ .Type }}) {
	for _, item := range items {
		delete(s.items, item)
	}
}

func (s *Set) Clear() {
	s.items = make(map[{{.Type}}]bool)
}

func (s *Set) Size() int {
	return len(s.items)
}

func (s *Set) Items() (res []{{.Type}}) {
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
{{end -}}