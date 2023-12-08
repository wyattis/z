package ziter

import (
	"reflect"
	"testing"
)

func TestSliceChunk(t *testing.T) {
	cases := [][3]int{{100, 10, 10}, {99, 10, 10}, {100, 9, 12}, {0, 100, 0}}
	for _, c := range cases {
		callCount := 0
		total := 0
		vals := make([]int, c[0])
		err := SliceChunks(vals, c[1], func(chunk []int) error {
			callCount++
			total += len(chunk)
			return nil
		})
		if err != nil {
			t.Error(err)
		} else if callCount != c[2] {
			t.Errorf("Expected handler to be called %d times, but got %d", c[2], callCount)
		} else if total != len(vals) {
			t.Errorf("Expected handler to be called with %d items, but got %d", len(vals), total)
		}
	}
}

func TestIteratorChunk(t *testing.T) {
	cases := [][]int8{{1, 2, 3, 4, 5}, {1}, {2, 3, 4, 5}}
	for _, c := range cases {
		res := make([]int8, 0, len(c))
		iterator := NewSliceIterator(c)
		err := IteratorChunk[int8](iterator, 3, func(chunk []int8) error {
			t.Log(chunk)
			res = append(res, chunk...)
			return nil
		})
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(c, res) {
			t.Errorf("expected %v to equal %v", c, res)
		}
	}
}
