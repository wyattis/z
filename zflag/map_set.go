package zflag

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wyattis/z/zstring"
)

type MapSet map[string][]string

func (s *MapSet) Set(val string) error {
	key, val, found := zstring.Cut(val, "|")
	if !found {
		return errors.New("not separator found in map")
	}
	parts := strings.Split(val, ",")
	_, exists := (*s)[key]
	if !exists {
		(*s)[key] = parts
	} else {
		(*s)[key] = append((*s)[key], parts...)
	}
	return nil
}

func (s *MapSet) String() string {
	return fmt.Sprint(*s)
}
