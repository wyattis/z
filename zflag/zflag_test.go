package zflag

import (
	"flag"
	"reflect"
	"strings"
	"testing"
	"time"
)

type baseConfig struct {
	Int     int
	Int64   int64
	Uint    uint
	Uint64  uint64
	Float64 float64
	String  string `flag:"name"`
	Bool    bool   `usage:"make this available"`
	Dur     time.Duration
}

var cases = []struct {
	args     string
	input    baseConfig
	expected baseConfig
}{
	{"-int 1 -uint 20", baseConfig{}, baseConfig{Int: 1, Uint: 20}},
	{"-bool -float64 -10", baseConfig{}, baseConfig{Bool: true, Float64: -10}},
	{"-dur 10s", baseConfig{}, baseConfig{Dur: time.Second * 10}},
}

func TestReflectStruct(t *testing.T) {
	for _, c := range cases {
		set := flag.NewFlagSet("test", flag.ExitOnError)
		if err := ReflectStruct(set, &c.input); err != nil {
			t.Error(err)
		}
		if err := set.Parse(strings.Split(c.args, " ")); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(c.input, c.expected) {
			t.Errorf("Expected %+v, but got %+v\n", c.expected, c.input)
		}
	}
}
