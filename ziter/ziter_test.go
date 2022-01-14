package ziter

import "testing"

func TestChunk(t *testing.T) {
	cases := [][3]int{{100, 10, 10}, {99, 10, 10}, {100, 9, 12}, {0, 100, 0}}
	for _, c := range cases {
		callCount := 0
		total := 0
		vals := make([]int, c[0])
		err := WithChunks(vals, c[1], func(chunk []interface{}) error {
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
