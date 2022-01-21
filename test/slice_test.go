package test

import (
	"testing"

	"github.com/wyattis/z/zslice/zstrings"
)

func TestSplit(t *testing.T) {
	cases := [][2][][]string{
		{zstrings.Split([]string{"one", "two", "three"}, "two"), [][]string{{"one"}, {"three"}}},
		{zstrings.Split([]string{"one", "two", "three"}, "three"), [][]string{{"one", "two"}}},
		{zstrings.Split([]string{"one", "two", "three"}, "one"), [][]string{{"two", "three"}}},
		{zstrings.Split([]string{"one", "two", "three"}, "four"), [][]string{{"one", "two", "three"}}},
		{zstrings.Split([]string{"one"}, "one"), [][]string{}},
	}
	for n, c := range cases {
		if len(c[0]) != len(c[1]) {
			t.Errorf("Expected %s, but got %s len", c[1], c[0])
		}
		for i := range c[1] {
			t.Log(n, c)
			if !zstrings.Equal(c[0][i], c[1][i]) {
				t.Errorf("Expected %s, but got %s", c[1], c[0])
			}
		}
	}
}

// func TestSplitMany(t *testing.T) {
// 	cases := [][2][][]string{
// 		{zstrings.SplitMany([]string{"one", "two", "three"}, "two", "three"), [][]string{{"one"}}},
// 		{zstrings.SplitMany([]string{"one", "two", "three"}, "three"), [][]string{{"one", "two"}}},
// 		{zstrings.SplitMany([]string{"one", "two", "three"}, "one"), [][]string{{"two", "three"}}},
// 		{zstrings.SplitMany([]string{"one", "two", "three"}, "one", "three", "two"), [][]string{}},
// 	}
// 	for _, c := range cases {
// 		for i := range c[1] {
// 			if !zstrings.Equal(c[0][i], c[1][i]) {
// 				t.Errorf("Expected %s, but got %s", c[1], c[0])
// 			}
// 		}
// 	}
// }
