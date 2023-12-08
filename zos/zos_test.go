package zos

import (
	"testing"

	"github.com/wyattis/z/zslice/zstrings"
)

func TestReadLinesCRLF(t *testing.T) {
	lines, err := ReadLinesFromFile("../test/assets/lines.crlf.txt")
	if err != nil {
		t.Error(err)
	}
	if !zstrings.Equal(lines, []string{"one", "two", "three", "four"}) {
		t.Error("Failed to read lines")
	}
}

func TestReadLinesLF(t *testing.T) {
	lines, err := ReadLinesFromFile("../test/assets/lines.lf.txt")
	if err != nil {
		t.Error(err)
	}
	if !zstrings.Equal(lines, []string{"one", "two", "three", "four"}) {
		t.Error("Failed to read lines")
	}
}
