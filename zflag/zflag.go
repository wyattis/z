package zflag

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/wyattis/z/zstring"
)

// Use reflection infer options for a flag.FlagSet based on the types and tags
// defined on a struct
func ReflectStruct(set *flag.FlagSet, config interface{}) (err error) {
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer")
	}
	v := reflect.Indirect(reflect.ValueOf(config))
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		t := v.Type().Field(i)
		name := zstring.CamelToSnake(t.Name)
		if tag := t.Tag.Get("flag"); tag != "" {
			name = tag
		}
		usage := t.Tag.Get("usage")
		if usage == "" {
			usage = fmt.Sprintf("%s is a %s", name, field.Kind().String())
		}
		defaultVal := t.Tag.Get("default")
		switch k := field.Kind(); k {
		case reflect.Bool:
			err = setBool(set, field, defaultVal, name, usage)
		case reflect.Int:
			err = setInt(set, field, defaultVal, name, usage)
		case reflect.Int64:
			err = setInt64(set, field, defaultVal, name, usage)
		case reflect.Uint:
			err = setUint(set, field, defaultVal, name, usage)
		case reflect.Uint64:
			err = setUint64(set, field, defaultVal, name, usage)
		case reflect.Float64:
			err = setFloat64(set, field, defaultVal, name, usage)
		case reflect.String:
			err = setString(set, field, defaultVal, name, usage)
		default:
			return fmt.Errorf("type of %s is not implemented", k)
		}
		if err != nil {
			return
		}
	}
	return
}

func setBool(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	defVal := false
	if defaultVal != "" {
		if defVal, err = strconv.ParseBool(defaultVal); err != nil {
			return
		}
	}
	p := (*bool)(unsafe.Pointer(field.Addr().Pointer()))
	set.BoolVar(p, name, defVal, usage)
	return
}

func setInt(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	defVal := 0
	if defaultVal != "" {
		if i, err := strconv.ParseInt(defaultVal, 10, 64); err != nil {
			return err
		} else {
			defVal = int(i)
		}
	}
	p := (*int)(unsafe.Pointer(field.Addr().Pointer()))
	set.IntVar(p, name, defVal, usage)

	return
}

func setUint(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	var defVal uint
	if defaultVal != "" {
		if i, err := strconv.ParseUint(defaultVal, 10, 64); err != nil {
			return err
		} else {
			defVal = uint(i)
		}
	}
	p := (*uint)(unsafe.Pointer(field.Addr().Pointer()))
	set.UintVar(p, name, defVal, usage)
	return
}

func setInt64(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	if field.Type().String() == "time.Duration" {
		return setDur(set, field, defaultVal, name, usage)
	}
	var defVal int64
	if defaultVal != "" {
		if i, err := strconv.ParseInt(defaultVal, 10, 64); err != nil {
			return err
		} else {
			defVal = i
		}
	}
	p := (*int64)(unsafe.Pointer(field.Addr().Pointer()))
	set.Int64Var(p, name, defVal, usage)
	return
}

func setDur(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	var defVal time.Duration
	if defaultVal != "" {
		if defVal, err = time.ParseDuration(defaultVal); err != nil {
			return err
		}
	}
	p := (*time.Duration)(unsafe.Pointer(field.Addr().Pointer()))
	set.DurationVar(p, name, defVal, usage)
	return
}

func setUint64(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	var defVal uint64
	if defaultVal != "" {
		if i, err := strconv.ParseUint(defaultVal, 10, 64); err != nil {
			return err
		} else {
			defVal = i
		}
	}
	p := (*uint64)(unsafe.Pointer(field.Addr().Pointer()))
	set.Uint64Var(p, name, defVal, usage)
	return
}

func setString(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	p := (*string)(unsafe.Pointer(field.Addr().Pointer()))
	set.StringVar(p, name, defaultVal, usage)
	return
}

func setFloat64(set *flag.FlagSet, field reflect.Value, defaultVal, name, usage string) (err error) {
	var defVal float64
	if defaultVal != "" {
		if i, err := strconv.ParseFloat(defaultVal, 10); err != nil {
			return err
		} else {
			defVal = i
		}
	}
	p := (*float64)(unsafe.Pointer(field.Addr().Pointer()))
	set.Float64Var(p, name, defVal, usage)
	return
}
