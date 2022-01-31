package zstring

import "testing"

type tcase struct {
	In  string
	Out string
}

func TestCamelToSnake(t *testing.T) {
	cases := []tcase{
		{"HelloWorld", "hello_world"},
		{"onceUponATime", "once_upon_a_time"},
		{"do_nothing", "do_nothing"},
		{"AMix_of_types", "a_mix_of_types"},
	}
	for _, c := range cases {
		res := CamelToSnake(c.In, "_")
		if c.Out != res {
			t.Errorf("Expected %s, but got %s", c.Out, res)
		}
	}
}

func TestSnakeToCamel(t *testing.T) {
	cases := []tcase{
		{"hello_world", "HelloWorld"},
		{"once_upon_a_time", "OnceUponATime"},
		{"DoNothing", "DoNothing"},
		{"AMix_of_types", "AMixOfTypes"},
		{"a_test___whew", "ATestWhew"},
		{"_ALeadingWord_test", "ALeadingWordTest"},
	}
	for _, c := range cases {
		res := SnakeToCamel(c.In, "_")
		if c.Out != res {
			t.Errorf("Expected %s, but got %s", c.Out, res)
		}
	}
}
