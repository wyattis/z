// Code generated by z. DO NOT EDIT.

package zstrings

import (
  "errors"
  "sort"
  "strings"
)

var (
	ErrInterfaceNotString = errors.New("encountered non-string interface")
)

// Determine if two slices are equal to each other
func Equal(a []string, b []string) bool {
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
func Sort(s []string) {
  sort.Slice(s, func (i, j int) bool {
    return strings.Compare(s[j], s[i]) > 0
  })
}

// Sort the slice in descending order
func SortReverse(s []string) {
  sort.Slice(s, func (i, j int) bool {
    return strings.Compare(s[i], s[j]) > 0
  })
}

// Resize a slice to the given length
func Resize(s []string, n int) (res []string) {
  diff := n - len(s)
  if diff > 0 {
    return append(s, make([]string, diff)...)
  } else {
    return s[:n]
  }
}


// Check if a slice ([]string) contains a matching member
func Contains(haystack []string, needle string) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []string slice contains ANY of the provided members
func ContainsAny(haystack []string, needles ...string) bool {
	for _, needle := range needles {
		if Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Check if a []string slice contains ALL of the provided members
func ContainsAll(haystack []string, needles ...string) bool {
	for _, needle := range needles {
		if !Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// Convert a []interface{} slice into []string slice if possible
func As(slice []interface{}) (res []string, err error) {
	res = make([]string, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(string); !ok {
			return res, ErrInterfaceNotString
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []string slice into a slice of interfaces ([]interface{})
func Interface(slice []string) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't
// appear in the slice.
func Index(haystack []string, needle string) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []string, separator string) (left, right []string, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []string, seperators ...string) (left, right []string, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Split the slice into one or more slices separated by the separator
func Split(slice []string, separator string) (parts [][]string) {
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

// Fill a slice with the given value
func Fill (slice []string, val string) []string {
  for i := range slice {
    slice[i] = val
  }
  return slice
}

// // Split the slice into one or more slices using multiple separators
// func SplitMany(slice []string, seperators ...string) (parts [][]string) {
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
// func Merge(a []string, b []string) (res []string) {
//
// }

// Remove the first occurrence of each value from the slice starting from the
// supplied offset
func Remove(slice []string, offset int, values ...string) (res []string) {
  for i := offset; i < len(slice); i++ {
    if !Contains(values, slice[i]) {
      res = append(res, slice[i])
    }
  }
  return
}

// Remove the item at the index while preserving the order of the items
// TODO: maybe change this to take multiple indices and to remove them at the 
// same time
func RemoveAt(slice []string, index int) (res []string) {
  if index >= len(slice) {
    panic("index cannot be greater than len(slice)")
  } else if index == len(slice) - 1 {
    res = slice[:index]
  } else {
    res = append(slice[:index], slice[index+1:]...)
  }
  return
}

// // Replace the first occurrence of a value with the replacement value
// func Replace(slice []string, val string, replacement string) (res []string) {
//
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []string, val string, replacement string) (res []string) {
//
// }
//
