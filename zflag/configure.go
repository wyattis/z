package zflag

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/wyattis/z/zreflect"
	"github.com/wyattis/z/zstring"
	"github.com/wyattis/z/ztime"
)

type convSetter struct {
	value *reflect.Value
	field *reflect.StructField
}

func (f *convSetter) Set(val string) error {
	k := f.value.Kind()
	if f.value.Type() == reflect.TypeOf(time.Time{}) {
		t, err := ztime.Parse(val, f.field.Tag.Get("time-format"))
		if err != nil {
			return err
		}
		f.value.Set(reflect.ValueOf(t))
		return nil
	} else if convert, exists := zreflect.ConvMap[k]; exists {
		val, err := convert(val, *f.value, zreflect.ConvMap)
		if err != nil {
			return err
		}
		f.value.Set(val)
		return nil
	} else {
		return errors.New("uknown type")
	}
}

func (f *convSetter) String() string {
	if f.value == nil {
		return "<nil>"
	}
	return (*f.value).String()
}

type FlagOptions struct {
	Overwrite bool
}

// Use reflection infer options for a flag.FlagSet based on the types and tags
// defined on a struct
func Configure(set *flag.FlagSet, config interface{}, opts *FlagOptions) error {
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer")
	}
	if opts == nil {
		opts = &FlagOptions{}
	}
	return recursiveSetFlags(set, reflect.Indirect(reflect.ValueOf(config)), "", *opts)
}

func recursiveSetFlags(set *flag.FlagSet, v reflect.Value, prefix string, opts FlagOptions) (err error) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		t := v.Type().Field(i)
		k := field.Kind()
		var found bool
		var name, usage, defaultVal string
		if tag := t.Tag.Get("flag"); tag != "" {
			// using cut here to allow commas in the usage string at the end
			name, defaultVal, found = zstring.Cut(tag, ",")
			if found && defaultVal != "" {
				defaultVal, usage, found = zstring.Cut(defaultVal, ",")
			}
		}
		if defaultVal == "" {
			defaultVal = t.Tag.Get("default")
		}
		if name == "" {
			name = zstring.CamelToSnake(t.Name, "-", 1)
		}
		if prefix != "" {
			name = strings.Join([]string{prefix, name}, "-")
			// name = strings.Join(append([]string{prefix}, strings.Split(name, "-")...), "-")
		}
		fmt.Println(t.Name, name)
		if usage == "" {
			usage = fmt.Sprintf("%s is a %s", name, field.Kind().String())
		}
		if field.Type() == reflect.TypeOf(time.Time{}) {
			res, err := ztime.Parse(defaultVal, t.Tag.Get("time-format"))
			if err != nil {
				return err
			}
			if field.IsZero() || opts.Overwrite {
				field.Set(reflect.ValueOf(res))
			}
			set.Var(&convSetter{value: &field, field: &t}, name, usage)
		} else if k == reflect.Struct {
			if err = recursiveSetFlags(set, field, name, opts); err != nil {
				return
			}
		} else if k == reflect.Bool {
			defVal := false
			if defaultVal != "" {
				if defVal, err = strconv.ParseBool(defaultVal); err != nil {
					return
				}
			}
			p := (*bool)(unsafe.Pointer(field.Addr().Pointer()))
			set.BoolVar(p, name, defVal, usage)
		} else if conv, exists := zreflect.ConvMap[k]; exists {
			// set the default value if one exists
			if defaultVal != "" {
				val, err := conv(defaultVal, field, zreflect.ConvMap)
				if err != nil {
					fmt.Println(err)
					return err
				}
				if field.IsZero() || opts.Overwrite {
					field.Set(val)
				}
			}
			set.Var(&convSetter{value: &field, field: &t}, name, usage)
		} else {
			fmt.Println("skipping invalid type", field, k, v)
		}
	}
	return nil
}
