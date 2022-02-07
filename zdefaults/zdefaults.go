package zdefaults

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/wyattis/z/zreflect"
	"github.com/wyattis/z/ztime"
)

func SetDefaults(val interface{}) (err error) {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	return setDefaultsRecursive(reflect.Indirect(reflect.ValueOf(val)))
}

func setDefaultsRecursive(val reflect.Value) (err error) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		t := val.Type().Field(i)
		k := field.Kind()
		defaultValue := t.Tag.Get("default")
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t, err := parseTime(defaultValue, t.Tag.Get("time-format"))
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(t))
		} else if k == reflect.Struct {
			if err = setDefaultsRecursive(field); err != nil {
				return
			}
		} else if !field.IsZero() {
			continue
		} else if defaultValue != "" {
			conv, exists := zreflect.ConvMap[k]
			if exists {
				val, err := conv(defaultValue, field, zreflect.ConvMap)
				if err != nil {
					return err
				}
				field.Set(val)
			} else {
				fmt.Println("Skipping", t, k)
				// TODO: check if there is a custom setter for the type
			}
		}
	}
	return
}

func parseTime(val string, format string) (t time.Time, err error) {
	formats := strings.Split(format, ";")
	return ztime.Parse(val, formats...)
}
