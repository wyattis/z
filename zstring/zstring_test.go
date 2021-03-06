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
		{"", ""},
	}
	for _, c := range cases {
		res := CamelToSnake(c.In, "_", 0)
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
		{"", ""},
	}
	for _, c := range cases {
		res := SnakeToCamel(c.In, "_", 0)
		if c.Out != res {
			t.Errorf("Expected %s, but got %s", c.Out, res)
		}
	}
}

func TestCutAt(t *testing.T) {
	type cutCase struct {
		cutAt int
		in    string
		left  string
		right string
	}
	cases := []cutCase{
		{3, "one two", "one", " two"},
		{-2, "100mb", "100", "mb"},
	}
	for _, c := range cases {
		left, right := CutAt(c.in, c.cutAt)
		if left != c.left || right != c.right {
			t.Errorf("Expected %s:%s, but got %s:%s", c.left, c.right, left, right)
		}
	}
}
