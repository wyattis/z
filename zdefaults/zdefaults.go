package zdefaults

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type typeConv func(val string, field reflect.Value, m conversionMap) (reflect.Value, error)
type conversionMap map[reflect.Kind]typeConv

func SetDefaults(val interface{}) (err error) {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	v := reflect.Indirect(reflect.ValueOf(val))
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.IsZero() {
			continue
		}
		t := v.Type().Field(i)
		k := field.Kind()
		defaultValue := t.Tag.Get("default")
		if defaultValue != "" {
			conv, exists := convMap[k]
			if exists {
				val, err := conv(defaultValue, field, convMap)
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

var convMap = map[reflect.Kind]typeConv{
	reflect.String: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		v = reflect.ValueOf(val)
		return
	},
	reflect.Bool: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return
		}
		v = reflect.ValueOf(b)
		return
	},
	reflect.Int: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.Atoi(val)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Int8: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return
		}
		v = reflect.ValueOf(int8(i))
		return
	},
	reflect.Int16: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return
		}
		v = reflect.ValueOf(int16(i))
		return
	},
	reflect.Int32: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return
		}
		v = reflect.ValueOf(int32(i))
		return
	},
	reflect.Int64: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Uint: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint(i))
		return
	},
	reflect.Uint8: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint8(i))
		return
	},
	reflect.Uint16: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint16(i))
		return
	},
	reflect.Uint32: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return
		}
		v = reflect.ValueOf(uint32(i))
		return
	},
	reflect.Uint64: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Float32: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseFloat(val, 10)
		if err != nil {
			return
		}
		v = reflect.ValueOf(float32(i))
		return
	},
	reflect.Float64: func(val string, _ reflect.Value, _ conversionMap) (v reflect.Value, err error) {
		i, err := strconv.ParseFloat(val, 10)
		if err != nil {
			return
		}
		v = reflect.ValueOf(i)
		return
	},
	reflect.Slice: func(val string, field reflect.Value, convMap conversionMap) (res reflect.Value, err error) {
		parts := strings.Split(val, ",")
		elem := field.Type().Elem()
		res = reflect.MakeSlice(reflect.SliceOf(elem), 0, len(parts))
		fmt.Println("reflect.Slice", parts, field, elem, res)
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
