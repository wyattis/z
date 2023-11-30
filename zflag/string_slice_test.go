package zflag

import (
	"flag"
	"testing"
)

func TestStringSlice(t *testing.T) {
	s := StringSliceVar{}
	flag.Var(&s, "string", "idk")
	flag.Parse()
}
