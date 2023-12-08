package zset

import "testing"

func TestSetAll(t *testing.T) {
	SetAddTest[int](t, 1, 2, 3)
	SetContainsTest[int](t, 1, 2, 3)
	SetDeleteTest[int](t, 1, 2, 3)
	SetClearTest[int](t, 1, 2, 3)
	SetSizeTest[int](t, 1, 2, 3)
	SetItemsTest[int](t, 1, 2, 3)
	SetEqualTest[int](t, 1, 2, 3)
	SetContainsAnyTest[int](t, 1, 2, 3)
	SetContainsAllTest[int](t, 1, 2, 3)
	SetCloneTest[int](t, 1, 2, 3)
	SetUnionTest[int](t, 1, 2, 3)
	SetIntersectionTest[int](t, 1, 2, 3)
	SetComplementTest[int](t, 1, 2, 3)
	SetFilterTest[int](t, 1, 2, 3)
	SetFilterSetItemsTest[int](t, 1, 2, 3)

	SetAddTest[string](t, "hello", "world", "!")
	SetContainsTest[string](t, "hello", "world", "!")
	SetDeleteTest[string](t, "hello", "world", "!")
	SetClearTest[string](t, "hello", "world", "!")
	SetSizeTest[string](t, "hello", "world", "!")
	SetItemsTest[string](t, "hello", "world", "!")
	SetEqualTest[string](t, "hello", "world", "!")
	SetContainsAnyTest[string](t, "hello", "world", "!")
	SetContainsAllTest[string](t, "hello", "world", "!")
	SetCloneTest[string](t, "hello", "world", "!")
	SetUnionTest[string](t, "hello", "world", "!")
	SetIntersectionTest[string](t, "hello", "world", "!")
	SetComplementTest[string](t, "hello", "world", "!")
	SetFilterTest[string](t, "hello", "world", "!")
	SetFilterSetItemsTest[string](t, "hello", "world", "!")
}

func SetAddTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	if s.Size() != 0 {
		t.Errorf("Expected size to be 0, got %d", s.Size())
	}
	s.Add(v)
	if s.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", s.Size())
	}
}

func SetContainsTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v)
	if !s.Contains(v) {
		t.Errorf("Expected set to contain \"%v\"", v)
	}
}

func SetDeleteTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v)
	if !s.Contains(v) {
		t.Errorf("Expected set to contain \"%v\"", v)
	}
	s.Delete(v)
	if s.Contains(v) {
		t.Errorf("Expected set to not contain \"%v\"", v)
	}
}

func SetClearTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
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

func SetSizeTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v, v, v, v)
	if s.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", s.Size())
	}
	s.Add(v2, v2, v2, v2)
	if s.Size() != 2 {
		t.Errorf("Expected size to be 2, got %d", s.Size())
	}
}

func SetItemsTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
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

func SetEqualTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s2 := New[T]()
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

func SetContainsAnyTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v, v2)
	if !s.ContainsAny(v, v2, v3) {
		t.Errorf("Expected set to contain any of \"%v\", \"%v\", \"%v\"", v, v2, v3)
	}
	if s.ContainsAny(v3) {
		t.Errorf("Expected set to not contain \"%v\"", v3)
	}
}

func SetContainsAllTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v, v2)
	if !s.Contains(v, v2) {
		t.Errorf("Expected set to contain all of \"%v\", \"%v\"", v, v2)
	}
	if s.Contains(v, v2, v3) {
		t.Errorf("Expected set to not contain \"%v\"", v3)
	}
}

func SetCloneTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v, v2)
	s2 := s.Clone()
	if s == s2 {
		t.Errorf("Expected sets to be different instances")
	}
	if !s2.Equal(s) {
		t.Errorf("Expected sets to be equal")
	}
}

func SetUnionTest[T comparable](t *testing.T, v, v2, v3 T) {
	s, s2 := New[T](), New[T]()
	s.Add(v)
	s2.Add(v2)
	s3 := s.Clone().Union(s2)
	if !s3.Contains(v, v2) {
		t.Errorf("Expected set to contain \"%v\", \"%v\"", v, v2)
	}
}

func SetIntersectionTest[T comparable](t *testing.T, v, v2, v3 T) {
	s, s2 := New[T](), New[T]()
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

func SetComplementTest[T comparable](t *testing.T, v, v2, v3 T) {
	s, s2 := New[T](), New[T]()
	s.Add(v, v2)
	s2.Add(v2, v3)
	s3 := s.Clone().Complement(s2)
	if !s3.Contains(v) {
		t.Errorf("Expected set to contain \"%v\"", v)
	}
	if s3.ContainsAny(v2, v3) {
		t.Errorf("Expected set to not contain \"%v\", \"%v\"", v2, v3)
	}
}

func SetFilterTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v, v2, v3)
	s2 := s.Clone().Filter(func(item T) bool {
		return item == v || item == v2
	})
	if !s2.Contains(v, v2) {
		t.Errorf("Expected set to contain \"%v\", \"%v\"", v, v2)
	}
	if s2.Contains(v3) {
		t.Errorf("Expected set to not contain \"%v\"", v3)
	}
}

func SetFilterSetItemsTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := New[T]()
	s.Add(v, v2, v3)
	s2 := s.Clone().FilterItems(func(item T) bool {
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
