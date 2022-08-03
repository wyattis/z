package zflag

import (
	"fmt"
	"strings"
)

func StringSlice(defaults ...string) *stringSlice {
	return &stringSlice{nil, defaults}
}

type stringSlice struct {
	val        []string
	defaultVal []string
}

func (s *stringSlice) SetDefault(val []string) {
	s.defaultVal = val
}

func (s stringSlice) Val() []string {
	if s.val != nil {
		return s.val
	} else {
		return s.defaultVal
	}
}

func (s stringSlice) Len() int {
	return len(s.Val())
}

func (s *stringSlice) Set(val string) error {
	parts := strings.Split(val, ",")
	for _, p := range parts {
		s.val = append(s.val, strings.TrimSpace(p))
	}
	return nil
}

func (s stringSlice) String() string {
	if s.val != nil {
		return fmt.Sprint(s.val)
	} else {
		return fmt.Sprint(s.defaultVal)
	}
}
