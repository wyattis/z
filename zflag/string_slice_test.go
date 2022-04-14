package zflag

import (
	"flag"
	"testing"
)

func TestStringSlice(t *testing.T) {
	s := StringSlice{}
	flag.Var(&s, "string", "idk")
	flag.Parse()
}
