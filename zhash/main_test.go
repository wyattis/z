package zhash

import (
	"fmt"
	"testing"
)

var passwords = []string{"example", "two", "12039849012341234123412341234"}
var fails = []string{"one", "fiver", "q", ""}

func TestExample(t *testing.T) {
	pass := []byte("example")
	hash, err := Hash(pass)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
	ok, err := Compare(hash, pass)
	if err != nil {
		panic(err)
	}
	fmt.Println(ok)
}

func TestMainHashAndCompare(t *testing.T) {
	for _, pass := range passwords {
		hash, err := Hash([]byte(pass))
		if err != nil {
			t.Error(err)
		}
		equal, err := Compare(hash, []byte(pass))
		if err != nil {
			t.Error(err)
		}
		if !equal {
			t.Errorf("Expected '%s' to hash correctly", pass)
		}
	}
}
func TestAllHashAndCompare(t *testing.T) {
	for _, alg := range available {
		for _, pass := range passwords {
			hash, err := alg.Hash([]byte(pass))
			if err != nil {
				t.Error(err)
			}
			equal, err := alg.Compare(hash, []byte(pass))
			if err != nil {
				t.Error(err)
			}
			if !equal {
				t.Errorf("alg %s: Expected '%s' to hash correctly", alg.Name(), pass)
			}
		}
	}
}

func TestAllHashAndFail(t *testing.T) {
	for _, alg := range available {
		for _, pass := range passwords {
			hash, err := alg.Hash([]byte(pass))
			if err != nil {
				t.Error(err)
			}
			for _, other := range fails {
				equal, err := alg.Compare(hash, []byte(other))
				if err != nil {
					t.Error(err)
				}
				if equal {
					t.Errorf("alg %s: Expected %s compared with %s to fail", alg.Name(), other, pass)
				}
			}
		}
	}
}
