// Code generated by z. DO NOT EDIT.

package zint32s

import "errors"

var (
	ErrInterfaceNotInt32 = errors.New("encountered non-int32 interface")
)

func Equal(a []int32, b []int32) bool {
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

// Check if a slice ([]int32) contains a matching member
func Contains(haystack []int32, needle int32) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []int32 slice contains ANY of the provided members
func ContainsAny(haystack []int32, needles ...int32) bool {
  for _, needle := range needles {
    if Contains(haystack, needle) {
      return true
    }
  }
  return false
}

// Check if a []int32 slice contains ALL of the provided members
func ContainsAll(haystack []int32, needles ...int32) bool {
  for _, needle := range needles {
    if !Contains(haystack, needle) {
      return false
    }
  }
  return true
}

// Convert a []interface{} slice into []int32 slice if possible
func As(slice []interface{}) (res []int32, err error) {
	res = make([]int32, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(int32); !ok {
			return res, ErrInterfaceNotInt32
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []int32 slice into a slice of interfaces ([]interface{})
func Interface(slice []int32) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't 
// appear in the slice.
func Index(haystack []int32, needle int32) int {
  for i := range haystack {
    if haystack[i] == needle {
      return i
    }
  }
  return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []int32, separator int32) (left, right []int32, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []int32, seperators ...int32) (left, right []int32, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Split the slice into one or more slices separated by the separator
func Split(slice []int32, separator int32) (parts [][]int32) {
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
// func SplitMany(slice []int32, seperators ...int32) (parts [][]int32) {
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
// func Merge(a []int32, b []int32) (res []int32) {
//   
// }

// // Remove the first occurrence of each value from the slice starting from the 
// // supplied offset
// func Remove(slice []int32, values ...int32, offset int) (res []int32) {
// 
// }

// // Replace the first occurrence of a value with the replacement value
// func Replace(slice []int32, val int32, replacement int32) (res []int32) {
// 
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []int32, val int32, replacement int32) (res []int32) {
// 
// }
// 

