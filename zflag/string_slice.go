package zflag

import (
	"fmt"
	"strings"
)

func StringSlice(defaults ...string) *StringSliceVar {
	return &StringSliceVar{nil, defaults}
}

type StringSliceVar struct {
	val        []string
	defaultVal []string
}

func (s *StringSliceVar) SetDefault(val []string) {
	s.defaultVal = val
}

func (s StringSliceVar) Val() []string {
	if s.val != nil {
		return s.val
	} else {
		return s.defaultVal
	}
}

func (s StringSliceVar) Len() int {
	return len(s.Val())
}

func (s *StringSliceVar) Set(val string) error {
	parts := strings.Split(val, ",")
	for _, p := range parts {
		s.val = append(s.val, strings.TrimSpace(p))
	}
	return nil
}

func (s StringSliceVar) String() string {
	if s.val != nil {
		return fmt.Sprint(s.val)
	} else {
		return fmt.Sprint(s.defaultVal)
	}
}
