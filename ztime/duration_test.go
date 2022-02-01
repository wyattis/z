package ztime

import (
	"testing"
)

type parseCase struct {
	in  string
	out Duration
}

func TestParseDuration(t *testing.T) {
	cases := []parseCase{
		{"1h", Hour},
		{"1d", Day},
		{"1w", Week},
		{"1M", Month},
		{"1Y", Year},
		{"1w2d5h", Week + 2*Day + 5*Hour},
		{"3Y6d", 3*Year + 6*Day},
		{"0Y", 0},
	}

	for _, c := range cases {
		res, err := ParseDuration(c.in)
		if err != nil {
			t.Error(err)
		}
		if res != c.out {
			t.Errorf("Expected %s, but got %s", c.out, res)
		}
	}
}
