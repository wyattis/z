package zconf

import (
	"errors"
	"reflect"
	"time"

	"github.com/wyattis/z/zreflect"
)

type defaultConfigurer struct{}

func (d *defaultConfigurer) Init(val interface{}) (err error) {
	return
}

// Iterate over the value and set any fields that have a default tag
func (d *defaultConfigurer) Apply(val interface{}, args ...string) (err error) {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	it := zreflect.FieldIterator(val)
	for it.Next() {
		v := it.Value()
		if !it.IsStructField() {
			continue
		}
		field := it.Field()
		defaultTag := field.Tag.Get("default")
		if defaultTag == "" {
			continue
		}
		if v.CanInterface() {
			f := v
			if f.CanAddr() {
				f = f.Addr()
			}
			if set, ok := f.Interface().(ConfigSettable); ok {
				if err = set.SetConfig(defaultTag); err != nil {
					return
				}
				it.DontDescend()
			}
		}
		k := v.Kind()
		newVal, err := zreflect.GetValueFromString(v, defaultTag, zreflect.WithStructTags(field.Tag))
		if err != nil {
			if err == zreflect.ErrCantConvertFromString {
				err = nil
				newVal = reflect.ValueOf(defaultTag)
			} else {
				return err
			}
		}
		// Handle pointers by creating a new instance if it's nil
		isPointer := k == reflect.Ptr
		if isPointer {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
			k = v.Kind()
		}
		it.Set(newVal)

		if field.Type == reflect.TypeOf(time.Time{}) {
			it.DontDescend()
		}
	}
	return
}
