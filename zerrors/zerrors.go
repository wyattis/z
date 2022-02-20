package zerrors

import (
	"errors"
	"fmt"
)

// Returns the first non-nil error in a slice of errors.
func AnyError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute each function until an error is encountered. Returns nil if no error
// is returned.
func Exec(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// Returns a merged version of each error encountered in a slice of errors. Each
// error is a newline delimited message. Returns nil if no error is returned.
func Merge(errs ...error) error {
	res := ""
	for _, e := range errs {
		if e != nil {
			res += fmt.Sprintln(e)
		}
	}
	if res != "" {
		return errors.New(res)
	}
	return nil
}

// Execute each function in a slice of functions. Returns a merged version of
// each error encountered. Each error is a newline delimited message. Returns
// nil if no error is encountered.
func MergeExec(fns ...func() error) error {
	res := ""
	for _, fn := range fns {
		if err := fn(); err != nil {
			res += fmt.Sprintln(err)
		}
	}
	if res != "" {
		return errors.New(res)
	}
	return nil
}
