package zreflect

import (
	"reflect"
	"testing"
	"time"

	"github.com/wyattis/z/ztime"
)

func TestSetTime(t *testing.T) {
	var ti time.Time
	if err := SetValue(reflect.ValueOf(&ti), "2015-01-01T00:00:00Z"); err != nil {
		t.Error(err)
	}
	expected := ztime.MustParse("2015-01-01T00:00:00Z", "2006-01-02T15:04:05Z")
	if !ti.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, t)
	}
}
