package zio

import (
	"bytes"
	"fmt"
	"io"
)

func CutReader(r io.Reader, needle []byte) (left []byte, remaining io.Reader, err error) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil {
			return nil, nil, err
		}
		if n == 0 {
			break
		}
		// TODO: This can skip through a needle if it's split across two reads
		for i := 0; i < n; i++ {
			if buf[i] == needle[0] {
				if n-i < len(needle) {
					break
				}
				if bytes.Equal(buf[i:i+len(needle)], needle) {
					return buf[:i], io.MultiReader(bytes.NewBuffer(buf[i:]), r), nil
				}
			}
		}
	}
	return nil, r, fmt.Errorf("needle not found")
}
