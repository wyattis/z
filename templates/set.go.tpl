{{- define "set" -}}

package {{.PackageName}}

func New() *{{title .TypeName}}Set {
  return &{{title .TypeName}}Set{
    items: make(map[{{.Type}}]bool),
  }
}

type {{title .TypeName}}Set struct {
  items map[{{.Type}}]bool
}

func (s *{{title .TypeName}}Set) Add(items ...{{ .Type }}) {
  for _, item := range items {
    s.items[item] = true
  }
}

func (s {{title .TypeName}}Set) Contains(items ...{{ .Type }}) bool {
  for _, item := range items {
    if _, exists := s.items[item]; !exists {
      return false
    }
  }
  return true
}

func (s *{{title .TypeName}}Set) Delete(items ...{{ .Type }}) {
  for _, item := range items {
      delete(s.items, item)
  }
}

func (s *{{title .TypeName}}Set) Clear() {
	s.items = make(map[{{.Type}}]bool)
}

func (s *{{title .TypeName}}Set) Size() int {
	return len(s.items)
}

func (s *{{title .TypeName}}Set) Items() (res []{{.Type}}) {
	for key := range s.items {
		res = append(res, key)
	}
	return
}

func (s *{{title .TypeName}}Set) Union(others ...{{title .TypeName}}Set) {
  for _, b := range others {
    for key := range b.items {
      s.items[key] = true
    }
  }
}

func (s *{{title .TypeName}}Set) Complement(others ...{{title .TypeName}}Set) {
  for _, b := range others {
    for key := range b.items {
      delete(s.items, key)
    }
  }
}

func (s *{{title .TypeName}}Set) Clone() *{{title .TypeName}}Set {
  res := New()
  res.Add(s.Items()...)
  return res
}

{{- end -}}