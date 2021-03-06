// Code generated by z. DO NOT EDIT.

package zbools

import (
  "errors"
  "sort"
  
)

var (
	ErrInterfaceNotBool = errors.New("encountered non-bool interface")
)

// Determine if two slices are equal to each other
func Equal(a []bool, b []bool) bool {
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
func Sort(s []bool) {
  sort.Slice(s, func (i, j int) bool {
    return s[j] && !s[i]
  })
}

// Sort the slice in descending order
func SortReverse(s []bool) {
  sort.Slice(s, func (i, j int) bool {
    return s[i] && !s[j]
  })
}

// Resize a slice to the given length
func Resize(s []bool, n int) (res []bool) {
  diff := n - len(s)
  if diff > 0 {
    return append(s, make([]bool, diff)...)
  } else {
    return s[:n]
  }
}


// Check if a slice ([]bool) contains a matching member
func Contains(haystack []bool, needle bool) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []bool slice contains ANY of the provided members
func ContainsAny(haystack []bool, needles ...bool) bool {
	for _, needle := range needles {
		if Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Check if a []bool slice contains ALL of the provided members
func ContainsAll(haystack []bool, needles ...bool) bool {
	for _, needle := range needles {
		if !Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// Convert a []interface{} slice into []bool slice if possible
func As(slice []interface{}) (res []bool, err error) {
	res = make([]bool, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(bool); !ok {
			return res, ErrInterfaceNotBool
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []bool slice into a slice of interfaces ([]interface{})
func Interface(slice []bool) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't
// appear in the slice.
func Index(haystack []bool, needle bool) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []bool, separator bool) (left, right []bool, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []bool, seperators ...bool) (left, right []bool, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Cut the slice at the specified index
func CutAt(val []bool, index int) (left, right []bool) {
  if index < 0 {
    index += len(val)
  }
  if len(val) > index {
    return val[:index], val[index:]
  }
  return val, nil
}

// Split the slice into one or more slices separated by the separator
func Split(slice []bool, separator bool) (parts [][]bool) {
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
func Fill (slice []bool, val bool) []bool {
  for i := range slice {
    slice[i] = val
  }
  return slice
}

// // Split the slice into one or more slices using multiple separators
// func SplitMany(slice []bool, seperators ...bool) (parts [][]bool) {
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
// func Merge(a []bool, b []bool) (res []bool) {
//
// }

// Remove the first occurrence of each value from the slice starting from the
// supplied offset
func Remove(slice []bool, offset int, values ...bool) (res []bool) {
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
func RemoveAt(slice []bool, index int) (res []bool) {
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
// func Replace(slice []bool, val bool, replacement bool) (res []bool) {
//
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []bool, val bool, replacement bool) (res []bool) {
//
// }
//
