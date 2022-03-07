package zsize

import "testing"

func TestSet(t *testing.T) {
	type setCase struct {
		in  string
		out Size
	}
	cases := []setCase{
		{"1MB", Megabyte},
		{"1KB", Kilobyte},
		{"10MB", 10 * Megabyte},
	}
	for _, c := range cases {
		s, err := Parse(c.in)
		if err != nil {
			t.Error(err)
		} else if s != c.out {
			t.Errorf("Expected %d, but got %d for %s", c.out, s, c.in)
		}
	}
}
