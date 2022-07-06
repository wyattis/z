package zstringset

import "testing"

func TestComplement(t *testing.T) {
	a := New()
	a.Add("one", "two", "three")
	b := New()
	b.Add("one", "three")
	a.Complement(*b)
	if !a.Contains("two") || a.Size() != 1 {
		t.Errorf("expected the complement to only contain 'two', but instead got %s", a.Items())
	}
}
