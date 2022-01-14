package zstring

import (
	"errors"
	"strings"
)

// Cut a string where a separator occurs
func Cut(val string, separator string) (left, right string, found bool) {
	if i := strings.Index(val, separator); i >= 0 {
		return val[:i], val[i+len(separator):], true
	}
	return val, "", false
}

// Cut the string where the first occurence of a separator occurs
func CutAny(val string, seperators ...string) (left, right string, found bool) {
	for _, sep := range seperators {
		left, right, found = Cut(val, sep)
		if found {
			return
		}
	}
	return
}

// Extract the part of a string surrounded by two separators
func CutOut(val, leftBound, rightBound string) (remaining string, found bool) {
	_, remaining, found = Cut(val, leftBound)
	if !found {
		return remaining, false
	}
	remaining, _, found = Cut(remaining, rightBound)
	return
}

// Extract the part of a string surrounded by multiple sets of separators
func CutOutAny(val string, leftBounds []string, rightBounds []string) (remaining string, found bool) {
	if len(leftBounds) != len(rightBounds) {
		panic(errors.New("left and right bounds must be the same length"))
	}
	for i := range leftBounds {
		remaining, found = CutOut(val, leftBounds[i], rightBounds[i])
		if found {
			return
		}
	}
	return
}

// Like ReplaceAll, but for replacing several sequences at once
func ReplaceManyWithOne(val string, needles []string, replacement string) string {
	for _, needle := range needles {
		val = strings.ReplaceAll(val, needle, replacement)
	}
	return val
}

// Like TrimSuffix, but for removing several suffixes at the same time
func TrimSuffixes(val string, suffixes ...string) string {
	for _, suffix := range suffixes {
		val = strings.TrimSuffix(val, suffix)
	}
	return val
}

// Like TrimPrefix, but for removing several prefixes at the same time
func TrimPrefixes(val string, prefixes ...string) string {
	for _, prefix := range prefixes {
		val = strings.TrimPrefix(val, prefix)
	}
	return val
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
