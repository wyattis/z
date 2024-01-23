{{ define "set_test" }}
package {{ .PackageName }}

import "testing"

func TestAdd(t *testing.T) {
  s := New()
  if s.Size() != 0 {
    t.Errorf("Expected size to be 0, got %d", s.Size())
  }
  {{ if .IsString -}}
  s.Add("hello")
  {{- else -}}
  s.Add(1)
  {{- end }}
  if s.Size() != 1 {
    t.Errorf("Expected size to be 1, got %d", s.Size())
  }
}

func TestContains (t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  {{- else -}}
  var v {{.TypeName}} = 10
  {{- end }}
  s.Add(v)
  if !s.Contains(v) {
    t.Errorf("Expected set to contain \"%v\"", v)
  }
}

func TestDelete (t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  {{- else -}}
  var v {{.TypeName}} = 10
  {{- end }}
  s.Add(v)
  if !s.Contains(v) {
    t.Errorf("Expected set to contain \"%v\"", v)
  }
  s.Delete(v)
  if s.Contains(v) {
    t.Errorf("Expected set to not contain \"%v\"", v)
  }
}

func TestClear (t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  {{- else -}}
  var v {{.TypeName}} = 10
  {{- end }}
  s.Add(v)
  if !s.Contains(v) {
    t.Errorf("Expected set to contain \"%v\"", v)
  } else if s.Size() != 1 {
    t.Errorf("Expected size to be 1, got %d", s.Size())
  }
  s.Clear()
  if s.Contains(v) {
    t.Errorf("Expected set to not contain \"%v\"", v)
  }
  if s.Size() != 0 {
    t.Errorf("Expected size to be 0, got %d", s.Size())
  }
}

func TestSize (t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  {{- end }}
  s.Add(v, v, v, v)
  if s.Size() != 1 {
    t.Errorf("Expected size to be 1, got %d", s.Size())
  }
  s.Add(v2, v2, v2, v2)
  if s.Size() != 2 {
    t.Errorf("Expected size to be 2, got %d", s.Size())
  }
}

func TestItems(t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  {{- end }}
  s.Add(v, v2)
  items := s.Items()
  if len(items) != 2 {
    t.Errorf("Expected items to have length 2, got %d", len(items))
  }
  if items[0] != v && items[0] != v2 {
    t.Errorf("Expected items to contain \"%v\" or \"%v\", got \"%v\"", v, v2, items[0])
  }
  if items[1] != v && items[1] != v2 {
    t.Errorf("Expected items to contain \"%v\" or \"%v\", got \"%v\"", v, v2, items[1])
  }
}

func TestEqual (t *testing.T) {
  s := New()
  s2 := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  {{- end }}
  s.Add(v, v2)
  s2.Add(v)
  if s.Equal(s2) {
    t.Errorf("Expected sets to not be equal")
  }
  s2.Add(v2)
  if !s.Equal(s2) {
    t.Errorf("Expected sets to be equal")
  }
}

func TestContainsAny(t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  v3 := "foo"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  var v3 {{.TypeName}} = 30
  {{- end }}
  s.Add(v, v2)
  if !s.ContainsAny(v, v2, v3) {
    t.Errorf("Expected set to contain any of \"%v\", \"%v\", \"%v\"", v, v2, v3)
  }
  if s.ContainsAny(v3) {
    t.Errorf("Expected set to not contain \"%v\"", v3)
  }
}

func TestContainsAll(t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  v3 := "foo"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  var v3 {{.TypeName}} = 30
  {{- end }}
  s.Add(v, v2)
  if !s.Contains(v, v2) {
    t.Errorf("Expected set to contain all of \"%v\", \"%v\"", v, v2)
  }
  if s.Contains(v, v2, v3) {
    t.Errorf("Expected set to not contain \"%v\"", v3)
  }
}

func TestClone (t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  {{- end }}
  s.Add(v, v2)
  s2 := s.Clone()
  if s == s2 {
    t.Errorf("Expected sets to be different instances")
  }
  if !s2.Equal(s) {
    t.Errorf("Expected sets to be equal")
  }
}

func TestUnion(t *testing.T) {
  s, s2 := New(), New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  {{- end }}
  s.Add(v)
  s2.Add(v2)
  s3 := s.Clone().Union(s2)
  if !s3.Contains(v, v2) {
    t.Errorf("Expected set to contain \"%v\", \"%v\"", v, v2)
  }
}

func TestIntersection(t *testing.T) {
  s, s2 := New(), New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  v3 := "foo"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  var v3 {{.TypeName}} = 30
  {{- end }}
  s.Add(v, v2)
  s2.Add(v2, v3)
  s3 := s.Clone().Intersection(s2)
  if !s3.Contains(v2) {
    t.Errorf("Expected set to contain \"%v\"", v2)
  }
  if s3.ContainsAny(v, v3) {
    t.Errorf("Expected set to not contain \"%v\", \"%v\"", v, v3)
  }
}

func TestDifference(t *testing.T) {
  s, s2 := New(), New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  v3 := "foo"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  var v3 {{.TypeName}} = 30
  {{- end }}
  s.Add(v, v2)
  s2.Add(v2, v3)
  s3 := s.Clone().Difference(s2)
  if !s3.Contains(v) {
    t.Errorf("Expected set to contain \"%v\"", v)
  }
  if s3.ContainsAny(v2, v3) {
    t.Errorf("Expected set to not contain \"%v\", \"%v\"", v2, v3)
  }
}

func TestFilter(t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  v3 := "foo"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  var v3 {{.TypeName}} = 30
  {{- end }}
  s.Add(v, v2, v3)
  s2 := s.Clone().Filter(func (item {{.TypeName}}) bool {
    return item == v || item == v2
  })
  if !s2.Contains(v, v2) {
    t.Errorf("Expected set to contain \"%v\", \"%v\"", v, v2)
  }
  if s2.Contains(v3) {
    t.Errorf("Expected set to not contain \"%v\"", v3)
  }
}

func TestFilterItems(t *testing.T) {
  s := New()
  {{ if .IsString -}}
  v := "hello"
  v2 := "world"
  v3 := "foo"
  {{- else -}}
  var v {{.TypeName}} = 10
  var v2 {{.TypeName}} = 20
  var v3 {{.TypeName}} = 30
  {{- end }}
  s.Add(v, v2, v3)
  s2 := s.Clone().FilterItems(func (item {{.TypeName}}) bool {
    return item == v || item == v2
  })
  if len(s2) != 2 {
    t.Errorf("Expected set to have length 2, got %d", len(s2))
  }
  if s2[0] != v && s2[0] != v2 {
    t.Errorf("Expected set to contain \"%v\" or \"%v\", got \"%v\"", v, v2, s2[0])
  }
  if s2[1] != v && s2[1] != v2 {
    t.Errorf("Expected set to contain \"%v\" or \"%v\", got \"%v\"", v, v2, s2[1])
  }
}

{{ end }}