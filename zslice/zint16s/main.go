// Code generated by z. DO NOT EDIT.

package zint16s

import (
  "errors"
  "sort"
  
)

var (
	ErrInterfaceNotInt16 = errors.New("encountered non-int16 interface")
)

// Determine if two slices are equal to each other
func Equal(a []int16, b []int16) bool {
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

// Sort the slice
func Sort(s []int16) {
  sort.Slice(s, func (i, j int) bool {
    return s[i] > s[j]
  })
}

// Check if a slice ([]int16) contains a matching member
func Contains(haystack []int16, needle int16) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []int16 slice contains ANY of the provided members
func ContainsAny(haystack []int16, needles ...int16) bool {
	for _, needle := range needles {
		if Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Check if a []int16 slice contains ALL of the provided members
func ContainsAll(haystack []int16, needles ...int16) bool {
	for _, needle := range needles {
		if !Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// Convert a []interface{} slice into []int16 slice if possible
func As(slice []interface{}) (res []int16, err error) {
	res = make([]int16, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(int16); !ok {
			return res, ErrInterfaceNotInt16
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []int16 slice into a slice of interfaces ([]interface{})
func Interface(slice []int16) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't
// appear in the slice.
func Index(haystack []int16, needle int16) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []int16, separator int16) (left, right []int16, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []int16, seperators ...int16) (left, right []int16, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Split the slice into one or more slices separated by the separator
func Split(slice []int16, separator int16) (parts [][]int16) {
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
// func SplitMany(slice []int16, seperators ...int16) (parts [][]int16) {
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
// func Merge(a []int16, b []int16) (res []int16) {
//
// }

// // Remove the first occurrence of each value from the slice starting from the
// // supplied offset
// func Remove(slice []int16, values ...int16, offset int) (res []int16) {
//
// }

// // Replace the first occurrence of a value with the replacement value
// func Replace(slice []int16, val int16, replacement int16) (res []int16) {
//
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []int16, val int16, replacement int16) (res []int16) {
//
// }
//
