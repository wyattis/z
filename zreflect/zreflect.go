package zreflect

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

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
