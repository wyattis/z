package zconfig

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestEnv(t *testing.T) {
	os.Setenv("GOOS", runtime.GOOS)
	os.Setenv("GOARCH", runtime.GOARCH)
	os.Setenv("GOOT", "wow")
	type config struct {
		GOOS   string
		GOARCH string
		Goot   string
	}
	conf := config{}
	c := New(Env())
	if err := c.Parse(&conf); err != nil {
		t.Error(err)
	}
	expected := config{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH, Goot: "wow"}
	if !reflect.DeepEqual(conf, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, conf)
	}
}

func TestDotEnv(t *testing.T) {
	os.Clearenv()
	data := []byte(fmt.Sprintf("GOOS=%s\nGOARCH=%s\n", runtime.GOOS, runtime.GOARCH))
	if err := os.WriteFile(".env.test", data, os.ModePerm); err != nil {
		t.Error(err)
	}
	defer os.Remove(".env.test")
	type config struct {
		GOOS   string
		GOARCH string
	}
	conf := config{}
	c := New(EnvFile(".env.test"))
	if err := c.Parse(&conf); err != nil {
		t.Error(err)
	}
	expected := config{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}
	if !reflect.DeepEqual(conf, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, conf)
	}
}
