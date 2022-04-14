package zflag

import (
	"fmt"
	"strings"
)

type StringSlice struct {
	val        []string
	defaultVal []string
}

func (s *StringSlice) SetDefault(val []string) {
	s.defaultVal = val
}

func (s *StringSlice) Val() []string {
	if s.val != nil {
		return s.val
	} else {
		return s.defaultVal
	}
}

func (s *StringSlice) Set(val string) error {
	parts := strings.Split(val, ",")
	for _, p := range parts {
		s.val = append(s.val, strings.TrimSpace(p))
	}
	return nil
}

func (s *StringSlice) String() string {
	if s.val != nil {
		return fmt.Sprint(s.val)
	} else {
		return fmt.Sprint(s.defaultVal)
	}
}
