// Code generated by z. DO NOT EDIT.

package zuints

import (
  "errors"
  "sort"
  
)

var (
	ErrInterfaceNotUint = errors.New("encountered non-uint interface")
)

// Determine if two slices are equal to each other
func Equal(a []uint, b []uint) bool {
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

// Sort the slice in ascending order
func Sort(s []uint) {
  sort.Slice(s, func (i, j int) bool {
    return s[j] > s[i]
  })
}

// Sort the slice in descending order
func SortReverse(s []uint) {
  sort.Slice(s, func (i, j int) bool {
    return s[i] > s[j]
  })
}

// Resize a slice to the given length
func Resize(s []uint, n int) (res []uint) {
  diff := n - len(s)
  if diff > 0 {
    return append(s, make([]uint, diff)...)
  } else {
    return s[:n]
  }
}


// Check if a slice ([]uint) contains a matching member
func Contains(haystack []uint, needle uint) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []uint slice contains ANY of the provided members
func ContainsAny(haystack []uint, needles ...uint) bool {
	for _, needle := range needles {
		if Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Check if a []uint slice contains ALL of the provided members
func ContainsAll(haystack []uint, needles ...uint) bool {
	for _, needle := range needles {
		if !Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// Convert a []interface{} slice into []uint slice if possible
func As(slice []interface{}) (res []uint, err error) {
	res = make([]uint, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(uint); !ok {
			return res, ErrInterfaceNotUint
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []uint slice into a slice of interfaces ([]interface{})
func Interface(slice []uint) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't
// appear in the slice.
func Index(haystack []uint, needle uint) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []uint, separator uint) (left, right []uint, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []uint, seperators ...uint) (left, right []uint, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Split the slice into one or more slices separated by the separator
func Split(slice []uint, separator uint) (parts [][]uint) {
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
// func SplitMany(slice []uint, seperators ...uint) (parts [][]uint) {
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
// func Merge(a []uint, b []uint) (res []uint) {
//
// }

// // Remove the first occurrence of each value from the slice starting from the
// // supplied offset
// func Remove(slice []uint, values ...uint, offset int) (res []uint) {
//
// }

// // Replace the first occurrence of a value with the replacement value
// func Replace(slice []uint, val uint, replacement uint) (res []uint) {
//
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []uint, val uint, replacement uint) (res []uint) {
//
// }
//
