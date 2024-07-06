package zio

import (
	"bufio"
	"io"
)

// A utility function that uses a bufio.Scanner to read lines from an io.Reader using bufio.ScanLines
func ReadLines(r io.Reader) (lines []string, err error) {
	return ReadLinesSplit(r, bufio.ScanLines)
}

// A utility function that uses a bufio.Scanner to read lines from an io.Reader using the provided split function
func ReadLinesSplit(r io.Reader, split bufio.SplitFunc) (lines []string, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(split)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}
