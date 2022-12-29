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

// Cut the string at the specified index
func CutAt(val string, index int) (left, right string) {
	if index < 0 {
		index += len(val)
	}
	if len(val) > index {
		return val[:index], val[index:]
	}
	return val, ""
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

// Determines if all characters in a string are uppercase
func IsUppercase(val string) bool {
	return strings.ToUpper(val) == val
}

// Determines if all characters in a string are lowercase
func IsLowercase(val string) bool {
	return strings.ToLower(val) == val
}

// Convert a camel string to snake case. Already snake cased strings aren't
// modified.
func CamelToSnake(val, separator string, minWordSize int) string {
	if separator == "" {
		separator = "_"
	}
	changes := []int{}
	wordSize := 0
	for i := range val {
		wordSize++
		if i == 0 || (val[i:i+1] != strings.ToLower(val[i:i+1]) && wordSize > minWordSize) {
			changes = append(changes, i)
			wordSize = 0
		}
	}
	changes = append(changes, len(val))
	parts := []string{}
	for i := 0; i < len(changes)-1; i++ {
		part := val[changes[i]:changes[i+1]]
		parts = append(parts, strings.ToLower(part))
	}
	return strings.Join(parts, separator)
}

// Convert a snake cased string to camel case. Ignores strings that already have
// camel casing.
func SnakeToCamel(val, separator string, minWordSize int) (res string) {
	if separator == "" {
		separator = "_"
	}
	val = CamelToSnake(val, separator, minWordSize)
	parts := strings.Split(val, separator)
	if len(parts) == 1 {
		return val
	}
	for _, p := range parts {
		if len(p) > 0 {
			res += strings.ToUpper(p[:1])
			if len(p) > 0 {
				res += strings.ToLower(p[1:])
			}
		}
	}
	return
}

// Split a string on space characters and trim leading/trailing space from each part
func SplitSpace(val string) (parts []string) {
	pParts := strings.Split(val, " ")
	for _, p := range pParts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			parts = append(parts, trimmed)
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
