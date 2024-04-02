package zconf

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/wyattis/z/zreflect"
	"github.com/wyattis/z/zstring"
)

type flagOpt = func(c *FlagConfigurer) error

// Provide a custom flag set. This is useful if you want to parse the flags yourself or are using an external library.
func UseFlagSet(fs *flag.FlagSet) flagOpt {
	return func(c *FlagConfigurer) error {
		c.flagSet = fs
		return nil
	}
}

// Tell the Flag configurer not to parse the flag set. This is useful if you want to parse the flags yourself or are
// using an external library.
func SkipParse() flagOpt {
	return func(c *FlagConfigurer) error {
		c.dontParseFlagSet = true
		return nil
	}
}

type flagVal struct {
	rVal reflect.Value
	tag  reflect.StructTag
	val  any
}

type FlagConfigurer struct {
	flagSet          *flag.FlagSet
	prefix           string
	vals             map[string]flagVal
	dontParseFlagSet bool
}

func (f *FlagConfigurer) SetPrefix(prefix string) {
	f.prefix = prefix
}

func (f *FlagConfigurer) SetFlagSet(flagSet *flag.FlagSet) {
	f.flagSet = flagSet
}

func (f *FlagConfigurer) Init(val interface{}) (err error) {
	// TODO: create the flags
	if f.flagSet == nil {
		f.flagSet = flag.CommandLine
	}
	f.vals = make(map[string]flagVal)
	it := zreflect.FieldIterator(val)
	for it.Next() {
		v := it.Value()
		if !it.IsStructField() {
			continue
		}
		field := it.Field()
		path := it.Path()
		key := it.Key()
		flagTag := field.Tag.Get("flag")
		name, usage, _ := strings.Cut(flagTag, ",")
		if name == "" {
			parts := make([]string, len(path)+1)
			copy(parts, path)
			if f.prefix != "" {
				parts = append([]string{f.prefix}, parts...)
			}
			parts[len(parts)-1] = key
			for i := range parts {
				parts[i] = zstring.CamelToSnake(parts[i], "-", 2)
			}
			name = strings.ToLower(strings.Join(parts, "-"))
		} else {
			name = strings.ToLower(name)
		}
		if usage == "" {
			usage = fmt.Sprintf("Set the %s value", name)
		}
		k := v.Kind()
		t := v.Type()
		if t == reflect.TypeOf(time.Time{}) || t == reflect.TypeOf(time.Duration(0)) {
			val := ""
			f.vals[name] = flagVal{rVal: v, val: &val, tag: field.Tag}
			f.flagSet.StringVar(&val, name, val, usage)
			it.DontDescend()
			continue
		} else if k == reflect.Struct {
			continue
		}
		switch k {
		case reflect.Bool:
			val := false
			f.vals[name] = flagVal{rVal: v, val: &val, tag: field.Tag}
			f.flagSet.BoolVar(&val, name, v.Bool(), usage)
		case reflect.String:
			val := ""
			f.vals[name] = flagVal{rVal: v, val: &val, tag: field.Tag}
			f.flagSet.StringVar(&val, name, v.String(), usage)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := 0
			f.vals[name] = flagVal{rVal: v, val: &val, tag: field.Tag}
			f.flagSet.IntVar(&val, name, int(v.Int()), usage)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val := uint(0)
			f.vals[name] = flagVal{rVal: v, val: &val, tag: field.Tag}
			f.flagSet.UintVar(&val, name, uint(v.Uint()), usage)
		case reflect.Float32, reflect.Float64:
			val := float64(0)
			f.vals[name] = flagVal{rVal: v, val: &val, tag: field.Tag}
			f.flagSet.Float64Var(&val, name, float64(v.Float()), usage)
		default:
			return fmt.Errorf("unsupported type %s", v.Type())
		}
	}
	return
}

func (f *FlagConfigurer) Apply(val interface{}, args ...string) (err error) {
	if !f.dontParseFlagSet {
		if err = f.flagSet.Parse(args); err != nil {
			return
		}
	}
	for _, v := range f.vals {
		t := v.rVal.Type()
		if err = zreflect.SetValue(v.rVal, v.val, zreflect.WithStructTags(v.tag)); err != nil {
			return
		}
		if t == reflect.TypeOf(time.Time{}) {

			if err = setTimeVal(v.rVal, v.tag, reflect.ValueOf(v.val).Elem().String()); err != nil {
				return
			}
			continue
		} else if t == reflect.TypeOf(time.Duration(0)) {
			if err = setDurationVal(v.rVal, reflect.ValueOf(v.val).Elem().String()); err != nil {
				return
			}
			continue
		}

		switch v.rVal.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if reflect.ValueOf(v.val).Elem().Int() == 0 {
				continue
			}
			v.rVal.SetInt(reflect.ValueOf(v.val).Elem().Int())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			if reflect.ValueOf(v.val).Elem().Uint() == 0 {
				continue
			}
			v.rVal.SetUint(reflect.ValueOf(v.val).Elem().Uint())
		case reflect.Float32, reflect.Float64:
			if reflect.ValueOf(v.val).Elem().Float() == 0 {
				continue
			}
			v.rVal.SetFloat(reflect.ValueOf(v.val).Elem().Float())
		case reflect.String:
			if reflect.ValueOf(v.val).Elem().String() == "" {
				continue
			}
			v.rVal.SetString(reflect.ValueOf(v.val).Elem().String())
		case reflect.Bool:
			if !reflect.ValueOf(v.val).Elem().Bool() {
				continue
			}
			v.rVal.SetBool(reflect.ValueOf(v.val).Elem().Bool())
		default:
			v.rVal.Set(reflect.ValueOf(v.val).Elem())
		}
	}
	return
}
