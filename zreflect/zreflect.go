package zreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/wyattis/z/ztime"
)

// Get the names of each field in a struct
func FieldNames(val interface{}, tags ...string) (res []string) {
	v := reflect.ValueOf(val)
	for i := 0; i < v.NumField(); i++ {
		t := v.Type().Field(i)
		name := t.Name
		for _, tag := range tags {
			candidate := t.Tag.Get(tag)
			if candidate != "" {
				name = candidate
				break
			}
		}
		res = append(res, name)
	}
	return
}

func WithStructTags(tags reflect.StructTag) opt {
	return func(c *setConfig) {
		c.tags = tags
	}
}

type setConfig struct {
	tags reflect.StructTag
}

func (c *setConfig) Load(opts []opt) *setConfig {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type opt func(*setConfig)

func SetValue(val reflect.Value, toVal any, opts ...opt) (err error) {
	conf := (&setConfig{}).Load(opts)
	val = reflect.Indirect(val)
	if val.CanInterface() {
		wasSet, err := setViaInterface(val, toVal)
		if err != nil || wasSet {
			return err
		}
	}
	if !val.CanSet() {
		err = fmt.Errorf("val is not settable")
		return
	}
	if val.Type() == reflect.TypeOf(time.Time{}) {
		if t, ok := toVal.(time.Time); ok {
			val.Set(reflect.ValueOf(t))
		} else {
			s, ok := toVal.(string)
			if !ok {
				if sp, ok := toVal.(*string); ok {
					s = *sp
				}
			}
			t, err = getTimeFromString(val, s, conf)
			if err != nil {
				return
			}
			val.Set(reflect.ValueOf(t))
		}
		return
	}

	if val.Type() == reflect.TypeOf(time.Second) {
		if d, ok := toVal.(time.Duration); ok {
			val.Set(reflect.ValueOf(d))
		} else if d, ok := toVal.(int64); ok {
			dur := time.Duration(d)
			val.Set(reflect.ValueOf(dur))
		} else {
			s, ok := toVal.(string)
			if !ok {
				if sp, ok := toVal.(*string); ok {
					s = *sp
				} else {
					err = fmt.Errorf("cannot set duration from %T", toVal)
					return
				}
			}
			d, err := getDurationFromString(val, s, conf)
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(d))
		}
		return
	}

	k := val.Kind()
	switch k {
	case reflect.String:
		if s, ok := toVal.(string); ok {
			val.SetString(s)
		} else if s, ok := toVal.(*string); ok {
			val.SetString(*s)
		} else if s, ok := toVal.(fmt.Stringer); ok {
			val.SetString(s.String())
		} else {
			err = fmt.Errorf("cannot set string from %T", toVal)
		}
	case reflect.Bool:
		if b, ok := toVal.(bool); ok {
			val.SetBool(b)
		} else if b, ok := toVal.(*bool); ok {
			val.SetBool(*b)
		} else {
			err = fmt.Errorf("cannot set bool from %T", toVal)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := getInt(val, toVal.(int64))
		if err != nil {
			return err
		}
		val.Set(i)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := getUint(val, toVal.(uint64))
		if err != nil {
			return err
		}
		val.Set(i)
		return nil
	case reflect.Float32, reflect.Float64:
		i, err := getFloat(val, toVal.(float64))
		if err != nil {
			return err
		}
		val.Set(i)
		return nil
	}
	return
}

// Parse the given string and use it to set the given value
func SetValueFromString(val reflect.Value, toVal string, opts ...opt) (err error) {
	val = reflect.Indirect(val)
	if val.CanInterface() {
		vi := val.Interface()
		switch t := vi.(type) {
		case StringSet:
			t.Set(toVal)
			return
		case StringSetErr:
			return t.Set(toVal)
		}
	}
	res, err := GetValueFromString(val, toVal, opts...)
	if err != nil {
		return
	}
	val.Set(res)
	return
}

var ErrCantConvertFromString = errors.New("unable to convert string into value")

func GetValueFromString(field reflect.Value, strVal string, opts ...opt) (res reflect.Value, err error) {
	conf := (&setConfig{}).Load(opts)
	isPointer := field.Kind() == reflect.Ptr

	val := reflect.Indirect(field)
	if isPointer && field.IsNil() {
		val = reflect.New(field.Type().Elem()).Elem()
	}

	res, err = getValFromString(val, strVal, conf)
	if err != nil {
		return
	}
	return
}

func getValFromString(val reflect.Value, strVal string, conf *setConfig) (res reflect.Value, err error) {
	switch val.Type() {
	case reflect.TypeOf(time.Time{}):
		t, err := getTimeFromString(val, strVal, conf)
		if err != nil {
			return res, err
		}
		return reflect.ValueOf(t), err
	case reflect.TypeOf(time.Second):
		d, err := getDurationFromString(val, strVal, conf)
		if err != nil {
			return res, err
		}
		return reflect.ValueOf(d), err
	}

	k := val.Kind()
	switch k {
	case reflect.String:
		return reflect.ValueOf(strVal), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(strVal)
		if err != nil {
			return res, err
		}
		return reflect.ValueOf(b), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		size := 64
		if k == reflect.Int8 {
			size = 8
		} else if k == reflect.Int16 {
			size = 16
		} else if k == reflect.Int32 {
			size = 32
		}
		i, err := strconv.ParseInt(strVal, 10, size)
		if err != nil {
			return res, err
		}
		return getInt(val, i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		size := 64
		if k == reflect.Uint8 {
			size = 8
		} else if k == reflect.Uint16 {
			size = 16
		} else if k == reflect.Uint32 {
			size = 32
		}
		i, err := strconv.ParseUint(strVal, 10, size)
		if err != nil {
			return res, err
		}
		return getUint(val, i)
	case reflect.Float32, reflect.Float64:
		size := 64
		if k == reflect.Float32 {
			size = 32
		}
		f, err := strconv.ParseFloat(strVal, size)
		if err != nil {
			return res, err
		}
		return getFloat(val, f)
	case reflect.Slice:
		// parts := strings.Split(toVal, ",")
		// elem := val.Type().Elem()
		// res := reflect.MakeSlice(reflect.SliceOf(elem), 0, len(parts))
		return res, errors.New("slice is not implemented")
	}

	return res, ErrCantConvertFromString
}

func getTimeFromString(val reflect.Value, strVal string, conf *setConfig) (t time.Time, err error) {
	format, ok := conf.tags.Lookup("time_format")
	if ok {
		t, err = time.Parse(format, strVal)
	} else {
		t, err = ztime.Parse(strVal)
	}
	if err != nil {
		return
	}
	return
}

func getDurationFromString(val reflect.Value, strVal string, conf *setConfig) (d time.Duration, err error) {
	return time.ParseDuration(strVal)
}

func setViaInterface(val reflect.Value, toVal any) (wasSet bool, err error) {
	vi := val.Interface()
	switch k := toVal.(type) {
	case string:
		if s, ok := vi.(StringSet); ok {
			s.Set(k)
			wasSet = true
		}
		if s, ok := vi.(StringSetErr); ok {
			if err = s.Set(k); err != nil {
				return
			}
			wasSet = true
		}
	case int:
		if s, ok := vi.(IntSet); ok {
			s.Set(k)
			wasSet = true
		}
		if s, ok := vi.(IntSetErr); ok {
			if err = s.Set(k); err != nil {
				return
			}
			wasSet = true
		}
	case float64:
		if s, ok := vi.(FloatSet); ok {
			s.Set(k)
			wasSet = true
		}
		if s, ok := vi.(FloatSetErr); ok {
			if err = s.Set(k); err != nil {
				return
			}
			wasSet = true
		}
	case bool:
		if s, ok := vi.(BoolSet); ok {
			s.Set(k)
			wasSet = true
		}
		if s, ok := vi.(BoolSetErr); ok {
			if err = s.Set(k); err != nil {
				return
			}
			wasSet = true
		}
	}
	return
}

func getInt(val reflect.Value, toVal int64) (res reflect.Value, err error) {
	switch val.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(toVal)), nil
	case reflect.Int8:
		return reflect.ValueOf(int8(toVal)), nil
	case reflect.Int16:
		return reflect.ValueOf(int16(toVal)), nil
	case reflect.Int32:
		return reflect.ValueOf(int32(toVal)), nil
	case reflect.Int64:
		return reflect.ValueOf(toVal), nil
	}
	err = fmt.Errorf("cannot set %T from int64", val)
	return
}

func getUint(val reflect.Value, toVal uint64) (res reflect.Value, err error) {
	switch val.Kind() {
	case reflect.Uint:
		return reflect.ValueOf(uint(toVal)), nil
	case reflect.Uint8:
		return reflect.ValueOf(uint8(toVal)), nil
	case reflect.Uint16:
		return reflect.ValueOf(uint16(toVal)), nil
	case reflect.Uint32:
		return reflect.ValueOf(uint32(toVal)), nil
	case reflect.Uint64:
		return reflect.ValueOf(toVal), nil
	}
	err = fmt.Errorf("cannot set %T from uint64", val)
	return
}

func getFloat(val reflect.Value, toVal float64) (res reflect.Value, err error) {
	switch val.Kind() {
	case reflect.Float64:
		return reflect.ValueOf(toVal), nil
	case reflect.Float32:
		return reflect.ValueOf(float32(toVal)), nil
	}
	err = fmt.Errorf("cannot set %T from float64", val)
	return
}

type TypeConv func(val string, field reflect.Value, m ConversionMap) (reflect.Value, error)
type ConversionMap map[reflect.Kind]TypeConv

var ConvMap = map[reflect.Kind]TypeConv{
	reflect.String: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		v = reflect.ValueOf(val)
		return
	},
	reflect.Bool: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return
		}
		v = reflect.ValueOf(b)
		return
	},
	reflect.Int: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.Atoi(val)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Int8: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return
		}
		v = reflect.ValueOf(int8(i))
		return
	},
	reflect.Int16: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return
		}
		v = reflect.ValueOf(int16(i))
		return
	},
	reflect.Int32: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return
		}
		v = reflect.ValueOf(int32(i))
		return
	},
	reflect.Int64: func(val string, field reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		if field.Type() == reflect.TypeOf(time.Second) {
			dur, err := time.ParseDuration(val)
			if err != nil {
				return v, err
			}
			v = reflect.ValueOf(dur)
			return v, err
		}
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Uint: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint(i))
		return
	},
	reflect.Uint8: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint8(i))
		return
	},
	reflect.Uint16: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint16(i))
		return
	},
	reflect.Uint32: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint32(i))
		return
	},
	reflect.Uint64: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Float32: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseFloat(val, 10)
		if err != nil {
			return
		}
		v = reflect.ValueOf(float32(i))
		return
	},
	reflect.Float64: func(val string, _ reflect.Value, _ ConversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseFloat(val, 10)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Slice: func(val string, field reflect.Value, convMap ConversionMap) (res reflect.Value, err error) {
		parts := strings.Split(val, ",")
		elem := field.Type().Elem()
		res = reflect.MakeSlice(reflect.SliceOf(elem), 0, len(parts))
		for _, p := range parts {
			conv, exists := convMap[elem.Kind()]
			if exists {
				val, err := conv(p, field, convMap)
				if err != nil {
					return res, err
				}
				res = reflect.Append(res, val)
			} else {
				err = fmt.Errorf("invalid type of %s found in slice", elem.Kind())
				return
			}
		}
		return
	},
}
