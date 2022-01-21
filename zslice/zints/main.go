// Code generated by z. DO NOT EDIT.

package zints

import "errors"

var (
	ErrInterfaceNotInt = errors.New("encountered non-int interface")
)

func Equal(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Check if a slice ([]int) contains a matching member
func Contains(haystack []int, needle int) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []int slice contains ANY of the provided members
func ContainsAny(haystack []int, needles ...int) bool {
	for _, needle := range needles {
		if Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Check if a []int slice contains ALL of the provided members
func ContainsAll(haystack []int, needles ...int) bool {
	for _, needle := range needles {
		if !Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// Convert a []interface{} slice into []int slice if possible
func As(slice []interface{}) (res []int, err error) {
	res = make([]int, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(int); !ok {
			return res, ErrInterfaceNotInt
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []int slice into a slice of interfaces ([]interface{})
func Interface(slice []int) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't
// appear in the slice.
func Index(haystack []int, needle int) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []int, separator int) (left, right []int, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []int, seperators ...int) (left, right []int, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Split the slice into one or more slices separated by the separator
func Split(slice []int, separator int) (parts [][]int) {
	prevI := 0
	for i := range slice {
		if slice[i] == separator {
			if i > prevI {
				parts = append(parts, slice[prevI:i])
			}
			prevI = i + 1
		}
	}
	if len(slice) > prevI {
		parts = append(parts, slice[prevI:])
	}
	return
}

// // Split the slice into one or more slices using multiple separators
// func SplitMany(slice []int, seperators ...int) (parts [][]int) {
// 	parts = append(parts, slice)
// 	for _, sep := range seperators {
// 		for i, part := range parts {
// 			vals := Split(part, sep)
// 			if len(vals) > 1 {
// 				// Replace existing element with all vals
// 				parts = append(parts[:i], append(vals, parts[i+1:]...)...)
// 			}
// 		}
// 	}
// 	return
// }

// // Merge two slices together without repeating values
// func Merge(a []int, b []int) (res []int) {
//
// }

// // Remove the first occurrence of each value from the slice starting from the
// // supplied offset
// func Remove(slice []int, values ...int, offset int) (res []int) {
//
// }

// // Replace the first occurrence of a value with the replacement value
// func Replace(slice []int, val int, replacement int) (res []int) {
//
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []int, val int, replacement int) (res []int) {
//
// }
//