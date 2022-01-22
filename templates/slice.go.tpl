{{- define "slice" -}}

package {{.PackageName}}

import (
  "errors"
  "sort"
  {{ if .IsString -}}
  "strings"
  {{- end }}
)

var (
	ErrInterfaceNot{{title .TypeName}} = errors.New("encountered non-{{.TypeName}} interface")
)

// Determine if two slices are equal to each other
func Equal(a []{{.Type}}, b []{{.Type}}) bool {
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
{{- if .IsString}}
func Sort(s []{{.Type}}) {
  sort.Slice(s, func (i, j int) bool {
    return strings.Compare(s[j], s[i]) > 0
  })
}
{{- else if .IsBool }}
func Sort(s []{{.Type}}) {
  sort.Slice(s, func (i, j int) bool {
    return s[j] && !s[i]
  })
}
{{- else}}
func Sort(s []{{.Type}}) {
  sort.Slice(s, func (i, j int) bool {
    return s[j] > s[i]
  })
}
{{- end }}

// Sort the slice in descending order
{{- if .IsString }}
func SortReverse(s []{{.Type}}) {
  sort.Slice(s, func (i, j int) bool {
    return strings.Compare(s[i], s[j]) > 0
  })
}
{{- else if .IsBool }}
func SortReverse(s []{{.Type}}) {
  sort.Slice(s, func (i, j int) bool {
    return s[i] && !s[j]
  })
}
{{- else }}
func SortReverse(s []{{.Type}}) {
  sort.Slice(s, func (i, j int) bool {
    return s[i] > s[j]
  })
}
{{- end }}

// Resize a slice to the given length
func Resize(s []{{.Type}}, n int) (res []{{.Type}}) {
  diff := n - len(s)
  if diff > 0 {
    return append(s, make([]{{.Type}}, diff)...)
  } else {
    return s[:n]
  }
}


// Check if a slice ([]{{.Type}}) contains a matching member
func Contains(haystack []{{.Type}}, needle {{.Type}}) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Check if a []{{.Type}} slice contains ANY of the provided members
func ContainsAny(haystack []{{.Type}}, needles ...{{.Type}}) bool {
	for _, needle := range needles {
		if Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// Check if a []{{.Type}} slice contains ALL of the provided members
func ContainsAll(haystack []{{.Type}}, needles ...{{.Type}}) bool {
	for _, needle := range needles {
		if !Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// Convert a []interface{} slice into []{{.Type}} slice if possible
func As(slice []interface{}) (res []{{.Type}}, err error) {
	res = make([]{{.Type}}, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].({{.Type}}); !ok {
			return res, ErrInterfaceNot{{title .TypeName}}
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert a []{{.Type}} slice into a slice of interfaces ([]interface{})
func Interface(slice []{{.Type}}) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// Find the index where the needle appears. Returns -1 if the needle doesn't
// appear in the slice.
func Index(haystack []{{.Type}}, needle {{.Type}}) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Cut the slice into two slices separated by the separator
func Cut(slice []{{.Type}}, separator {{.Type}}) (left, right []{{.Type}}, found bool) {
	if i := Index(slice, separator); i >= 0 {
		return slice[:i], slice[i+1:], true
	}
	return slice, nil, false
}

// Cut the slice into two slices separated by the first separator
func CutAny(val []{{.Type}}, seperators ...{{.Type}}) (left, right []{{.Type}}, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Split the slice into one or more slices separated by the separator
func Split(slice []{{.Type}}, separator {{.Type}}) (parts [][]{{.Type}}) {
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
// func SplitMany(slice []{{.Type}}, seperators ...{{.Type}}) (parts [][]{{.Type}}) {
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
// func Merge(a []{{.Type}}, b []{{.Type}}) (res []{{.Type}}) {
//
// }

// Remove the first occurrence of each value from the slice starting from the
// supplied offset
func Remove(slice []{{.Type}}, offset int, values ...{{.Type}}) (res []{{.Type}}) {
  for i := offset; i < len(slice); i++ {
    for _, val := range values {
      if slice[i] == val {
        break
      }
      res = append(res, slice[i])
    }
  }
  return
}

// Remove the item at the index while preserving the order of the items
func RemoveAt(slice []{{.Type}}, index int) (res []{{.Type}}) {
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
// func Replace(slice []{{.Type}}, val {{.Type}}, replacement {{.Type}}) (res []{{.Type}}) {
//
// }

// // Replace all occurrences of a value with the replacement value
// func ReplaceAll(slice []{{.Type}}, val {{.Type}}, replacement {{.Type}}) (res []{{.Type}}) {
//
// }
//
{{ end -}}