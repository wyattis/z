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
	Bool    bool   `flag:",,make this available"`
	Dur     time.Duration
}

type nestedConf struct {
	Int    int
	Nested baseConfig
}

func TestReflectBase(t *testing.T) {
	var cases = []struct {
		args     string
		input    baseConfig
		expected baseConfig
	}{
		{"-int 1 -uint 20", baseConfig{}, baseConfig{Int: 1, Uint: 20}},
		{"-bool -float64 -10", baseConfig{}, baseConfig{Bool: true, Float64: -10}},
		{"-dur 10s", baseConfig{}, baseConfig{Dur: time.Second * 10}},
	}
	for _, c := range cases {
		set := flag.NewFlagSet("test", flag.ExitOnError)
		if err := Configure(set, &c.input, nil); err != nil {
			t.Error(err)
		}
		args := strings.Split(c.args, " ")
		if err := set.Parse(args); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(c.input, c.expected) {
			t.Errorf("Expected %+v, but got %+v\n", c.expected, c.input)
		}
	}
}

func TestNested(t *testing.T) {
	var cases = []struct {
		args     string
		input    nestedConf
		expected nestedConf
	}{
		{"-int 1 -nested-uint 20", nestedConf{}, nestedConf{Int: 1, Nested: baseConfig{Uint: 20}}},
	}
	for _, c := range cases {
		set := flag.NewFlagSet("test", flag.ExitOnError)
		if err := Configure(set, &c.input, nil); err != nil {
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
