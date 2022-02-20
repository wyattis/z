package zhash

import "testing"

var passwords = []string{"example", "two", "12039849012341234123412341234"}
var fails = []string{"one", "fiver", "q", ""}

func TestHashAndCompare(t *testing.T) {
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

func TestHashAndFail(t *testing.T) {
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
