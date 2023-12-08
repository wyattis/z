package ziter

import "testing"

var _ Iterator[[]int] = &chunkIterator[int]{} // Ensure interface is implemented

func TestSliceNChunks(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7}
	chunks := SliceNChunks(slice, 3)
	if len(chunks) != 3 {
		t.Error("Expected 3 chunks")
	}
	if len(chunks[0]) != 3 {
		t.Error("Expected 3 items in first chunk")
	}
	if len(chunks[1]) != 3 {
		t.Error("Expected 3 items in second chunk")
	}
	if len(chunks[2]) != 1 {
		t.Error("Expected 1 item in third chunk")
	}
}

func TestSliceChunksOfSize(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7}
	chunks := SliceChunksOfSize(slice, 3)
	if len(chunks) != 3 {
		t.Error("Expected 3 chunks")
	}
	if len(chunks[0]) != 3 {
		t.Error("Expected 3 items in first chunk")
	}
	if len(chunks[1]) != 3 {
		t.Error("Expected 3 items in second chunk")
	}
	if len(chunks[2]) != 1 {
		t.Error("Expected 1 item in third chunk")
	}
}
