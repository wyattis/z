package zpath

import "testing"

func TestIsChildOf(t *testing.T) {
	type Case struct {
		Dir   string
		Child string
		Res   bool
	}
	cases := []Case{
		{"temp", "temp/wow", true},
		{"temp", "temp", true},
		{"../zpath", "file.csv", true},
		{"root", "../another-file.json", false},
		{"../fake-dir", "fake.csv", false},
		{"C:\\windows-path/mixed\\slashes", "C:/windows-path/mixed/slashes/child.exe", true},
		{"C:\\windows-path/mixed\\slashes", "C:/windows-path/mixed/child.exe", false},
	}
	for _, c := range cases {
		isChild, err := IsChildOf(c.Dir, c.Child)
		if err != nil {
			t.Error(err)
		}
		if c.Res != isChild {
			t.Errorf("Expected %t, but got %t for IsChildOf(%s, %s)", c.Res, isChild, c.Dir, c.Child)
		}
	}
}
