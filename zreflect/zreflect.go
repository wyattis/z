package zreflect

import (
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

func SetValue(val reflect.Value, toVal any) (err error) {
	// isPointer := val.Kind() == reflect.Ptr
	// fmt.Println("isPointer", isPointer)
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
		} else if s, ok := toVal.(string); ok {
			t, err := ztime.Parse(s)
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(t))
		} else if s, ok := toVal.(*string); ok {
			t, err := ztime.Parse(*s)
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(t))
		} else {
			err = fmt.Errorf("cannot set time from %T", toVal)
		}
		return
	}

	if val.Type() == reflect.TypeOf(time.Second) {
		if d, ok := toVal.(time.Duration); ok {
			val.Set(reflect.ValueOf(d))
		} else if d, ok := toVal.(int64); ok {
			dur := time.Duration(d)
			val.Set(reflect.ValueOf(dur))
		} else if s, ok := toVal.(string); ok {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(d))
		} else if s, ok := toVal.(*string); ok {
			d, err := time.ParseDuration(*s)
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(d))
		} else {
			err = fmt.Errorf("cannot set duration from %T", toVal)
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
		return setInt(val, toVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUint(val, toVal)
	case reflect.Float32, reflect.Float64:
		return setFloat(val, toVal)
	}
	return
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

func setInt(val reflect.Value, toVal any) (err error) {
	i, ok := toVal.(int)
	if !ok {
		if ip, ok := toVal.(*int); ok {
			i = *ip
		} else if i64, ok := toVal.(int64); ok {
			i = int(i64)
		} else if i64p, ok := toVal.(*int64); ok {
			i = int(*i64p)
		} else if i32, ok := toVal.(int32); ok {
			i = int(i32)
		} else if i32p, ok := toVal.(*int32); ok {
			i = int(*i32p)
		} else if i16, ok := toVal.(int16); ok {
			i = int(i16)
		} else if i16p, ok := toVal.(*int16); ok {
			i = int(*i16p)
		} else if i8, ok := toVal.(int8); ok {
			i = int(i8)
		} else if i8p, ok := toVal.(*int8); ok {
			i = int(*i8p)
		} else {
			err = fmt.Errorf("cannot set int from %T", toVal)
			return
		}
	}
	switch val.Kind() {
	case reflect.Int:
		val.SetInt(int64(i))
	case reflect.Int8:
		val.SetInt(int64(int8(i)))
	case reflect.Int16:
		val.SetInt(int64(int16(i)))
	case reflect.Int32:
		val.SetInt(int64(int32(i)))
	case reflect.Int64:
		val.SetInt(int64(i))
	}
	return
}

func setUint(val reflect.Value, toVal any) (err error) {
	i, ok := toVal.(uint)
	if !ok {
		if ip, ok := toVal.(*uint); ok {
			i = *ip
		} else if i64, ok := toVal.(uint64); ok {
			i = uint(i64)
		} else if i64p, ok := toVal.(*uint64); ok {
			i = uint(*i64p)
		} else if i32, ok := toVal.(uint32); ok {
			i = uint(i32)
		} else if i32p, ok := toVal.(*uint32); ok {
			i = uint(*i32p)
		} else if i16, ok := toVal.(uint16); ok {
			i = uint(i16)
		} else if i16p, ok := toVal.(*uint16); ok {
			i = uint(*i16p)
		} else if i8, ok := toVal.(uint8); ok {
			i = uint(i8)
		} else if i8p, ok := toVal.(*uint8); ok {
			i = uint(*i8p)
		} else {
			err = fmt.Errorf("cannot set uint from %T", toVal)
			return
		}
	}

	switch val.Kind() {
	case reflect.Uint:
		val.SetUint(uint64(i))
	case reflect.Uint8:
		val.SetUint(uint64(uint8(i)))
	case reflect.Uint16:
		val.SetUint(uint64(uint16(i)))
	case reflect.Uint32:
		val.SetUint(uint64(uint32(i)))
	case reflect.Uint64:
		val.SetUint(uint64(i))
	}
	return
}

func setFloat(val reflect.Value, toVal any) (err error) {
	var f float64
	if i, ok := toVal.(float64); ok {
		f = i
	} else if ip, ok := toVal.(*float64); ok {
		f = *ip
	} else if i32, ok := toVal.(float32); ok {
		f = float64(i32)
	} else if i32p, ok := toVal.(*float32); ok {
		f = float64(*i32p)
	} else {
		err = fmt.Errorf("cannot set float from %T", toVal)
		return
	}

	switch val.Kind() {
	case reflect.Float32:
		val.SetFloat(float64(float32(f)))
	case reflect.Float64:
		val.SetFloat(float64(f))
	}
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
