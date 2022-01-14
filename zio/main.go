package zio

import "io"

// Close multiple Closers returning an error if any of them produce one
func CloseAll(files ...io.Closer) error {
	for _, f := range files {
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}
