// Code generated by z. DO NOT EDIT.

package zfloat64s

import "errors"

var (
	ErrInterfaceNotFloat64 = errors.New("encountered non-float64 interface")
)

func Equal(a []float64, b []float64) bool {
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

// Check if a slice ([]float64) contains a matching member
func Contains(haystack []float64, needle float64) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []float64 slice contains ANY of the provided members
func ContainsAny(haystack []float64, needles ...float64) bool {
	for _, needle := range needles {
		if Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Check if a []float64 slice contains ALL of the provided members
func ContainsAll(haystack []float64, needles ...float64) bool {
	for _, needle := range needles {
		if !Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// Convert a []interface{} slice into []float64 slice if possible
func As(slice []interface{}) (res []float64, err error) {
	res = make([]float64, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(float64); !ok {
			return res, ErrInterfaceNotFloat64
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []float64 slice into a slice of interfaces ([]interface{})
func Interface(slice []float64) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't
// appear in the slice.
func Index(haystack []float64, needle float64) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []float64, separator float64) (left, right []float64, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []float64, seperators ...float64) (left, right []float64, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Split the slice into one or more slices separated by the separator
func Split(slice []float64, separator float64) (parts [][]float64) {
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
// func SplitMany(slice []float64, seperators ...float64) (parts [][]float64) {
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
// func Merge(a []float64, b []float64) (res []float64) {
//
// }

// // Remove the first occurrence of each value from the slice starting from the
// // supplied offset
// func Remove(slice []float64, values ...float64, offset int) (res []float64) {
//
// }

// // Replace the first occurrence of a value with the replacement value
// func Replace(slice []float64, val float64, replacement float64) (res []float64) {
//
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []float64, val float64, replacement float64) (res []float64) {
//
// }
//
