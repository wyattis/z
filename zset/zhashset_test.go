package zset

import "testing"

func TestHashSetAll(t *testing.T) {
	HashSetAddTest[int](t, 1, 2, 3)
	HashSetContainsTest[int](t, 1, 2, 3)
	HashSetDeleteTest[int](t, 1, 2, 3)
	HashSetClearTest[int](t, 1, 2, 3)
	HashSetSizeTest[int](t, 1, 2, 3)
	HashSetItemsTest[int](t, 1, 2, 3)
	HashSetEqualTest[int](t, 1, 2, 3)
	HashSetContainsAnyTest[int](t, 1, 2, 3)
	HashSetContainsAllTest[int](t, 1, 2, 3)
	HashSetCloneTest[int](t, 1, 2, 3)
	HashSetUnionTest[int](t, 1, 2, 3)
	HashSetIntersectionTest[int](t, 1, 2, 3)
	HashSetComplementTest[int](t, 1, 2, 3)
	HashSetFilterTest[int](t, 1, 2, 3)
	HashSetFilterItemsTest[int](t, 1, 2, 3)

	HashSetAddTest[string](t, "hello", "world", "!")
	HashSetContainsTest[string](t, "hello", "world", "!")
	HashSetDeleteTest[string](t, "hello", "world", "!")
	HashSetClearTest[string](t, "hello", "world", "!")
	HashSetSizeTest[string](t, "hello", "world", "!")
	HashSetItemsTest[string](t, "hello", "world", "!")
	HashSetEqualTest[string](t, "hello", "world", "!")
	HashSetContainsAnyTest[string](t, "hello", "world", "!")
	HashSetContainsAllTest[string](t, "hello", "world", "!")
	HashSetCloneTest[string](t, "hello", "world", "!")
	HashSetUnionTest[string](t, "hello", "world", "!")
	HashSetIntersectionTest[string](t, "hello", "world", "!")
	HashSetComplementTest[string](t, "hello", "world", "!")
	HashSetFilterTest[string](t, "hello", "world", "!")
	HashSetFilterItemsTest[string](t, "hello", "world", "!")
}

func HashSetAddTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	if s.Size() != 0 {
		t.Errorf("Expected size to be 0, got %d", s.Size())
	}
	s.Add(v)
	if s.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", s.Size())
	}
}

func HashSetContainsTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	s.Add(v)
	if !s.Contains(v) {
		t.Errorf("Expected set to contain \"%v\"", v)
	}
}

func HashSetDeleteTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	s.Add(v)
	if !s.Contains(v) {
		t.Errorf("Expected set to contain \"%v\"", v)
	}
	s.Delete(v)
	if s.Contains(v) {
		t.Errorf("Expected set to not contain \"%v\"", v)
	}
}

func HashSetClearTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
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

func HashSetSizeTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	s.Add(v, v, v, v)
	if s.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", s.Size())
	}
	s.Add(v2, v2, v2, v2)
	if s.Size() != 2 {
		t.Errorf("Expected size to be 2, got %d", s.Size())
	}
}

func HashSetItemsTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
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

func HashSetEqualTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	s2 := NewHashSet(func(item T) T { return item })
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

func HashSetContainsAnyTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	s.Add(v, v2)
	if !s.ContainsAny(v, v2, v3) {
		t.Errorf("Expected set to contain any of \"%v\", \"%v\", \"%v\"", v, v2, v3)
	}
	if s.ContainsAny(v3) {
		t.Errorf("Expected set to not contain \"%v\"", v3)
	}
}

func HashSetContainsAllTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	s.Add(v, v2)
	if !s.Contains(v, v2) {
		t.Errorf("Expected set to contain all of \"%v\", \"%v\"", v, v2)
	}
	if s.Contains(v, v2, v3) {
		t.Errorf("Expected set to not contain \"%v\"", v3)
	}
}

func HashSetCloneTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
	s.Add(v, v2)
	s2 := s.Clone()
	if s == s2 {
		t.Errorf("Expected sets to be different instances")
	}
	if !s2.Equal(s) {
		t.Errorf("Expected sets to be equal")
	}
}

func HashSetUnionTest[T comparable](t *testing.T, v, v2, v3 T) {
	s, s2 := NewHashSet(func(item T) T { return item }), NewHashSet(func(item T) T { return item })
	s.Add(v)
	s2.Add(v2)
	s3 := s.Clone().Union(s2)
	if !s3.Contains(v, v2) {
		t.Errorf("Expected set to contain \"%v\", \"%v\"", v, v2)
	}
}

func HashSetIntersectionTest[T comparable](t *testing.T, v, v2, v3 T) {
	s, s2 := NewHashSet(func(item T) T { return item }), NewHashSet(func(item T) T { return item })
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

func HashSetComplementTest[T comparable](t *testing.T, v, v2, v3 T) {
	s, s2 := NewHashSet(func(item T) T { return item }), NewHashSet(func(item T) T { return item })
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

func HashSetFilterTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
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

func HashSetFilterItemsTest[T comparable](t *testing.T, v, v2, v3 T) {
	s := NewHashSet(func(item T) T { return item })
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
