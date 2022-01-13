package zstring

import "strings"

func Cut(val string, separator string) (left, right string, found bool) {
	if i := strings.Index(val, separator); i >= 0 {
		return val[:i], val[i+len(separator):], true
	}
	return val, "", false
}

func CutAny(val string, seperators ...string) (left, right string, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// func SplitMany(val string, seperators ...string) (parts []string) {
// 	parts = append(parts, val)
// 	for _, sep := range seperators {
// 		for i, part := range parts {
// 			vals := strings.Split(part, sep)
// 			if len(vals) > 1 {
// 				// Replace existing element with all vals
// 				parts = append(parts[:i], append(vals, parts[i+1:]...)...)
// 			}
// 		}
// 	}
// 	return
// }
