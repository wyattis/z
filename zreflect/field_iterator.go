package zreflect

import (
	"reflect"
)

type containerKind int

const (
	Struct containerKind = iota
	Map
)

func FieldIterator(v any) *fieldIterator {
	return newFieldIterator(reflect.ValueOf(v))
}

func newFieldIterator(v reflect.Value) *fieldIterator {
	container := reflect.Indirect(v)
	container, kind, ok := kindFromValue(container)
	if !ok {
		panic("FieldIterator only works on structs and maps")
	}
	if kind == Map {
		container = container.Addr()
	}
	it := &fieldIterator{
		container: container,
		kind:      kind,
	}
	it.init()
	return it
}

type fieldIterator struct {
	container     reflect.Value
	kind          containerKind
	path          []string
	fieldIndex    int
	numFields     int
	keys          []reflect.Value
	shouldDescend bool // if we should descend into the field or not
	child         *fieldIterator
}

func (f *fieldIterator) DontDescend() {
	if f.child != nil {
		f.child.DontDescend()
		return
	}
	f.shouldDescend = false
}

func (f *fieldIterator) init() {
	if f.kind == Struct {
		f.numFields = f.container.NumField()
	} else {
		f.keys = f.container.Elem().MapKeys()
	}
	f.fieldIndex = -1
}

func (f *fieldIterator) remainingFields() int {
	if f.kind == Struct {
		return f.numFields - f.fieldIndex - 1
	} else {
		return len(f.keys) - f.fieldIndex - 1
	}
}

func (f *fieldIterator) Next() bool {
	if f.child != nil {
		if f.child.Next() {
			return true
		}
		f.child = nil
		f.shouldDescend = false
	}

	if f.shouldDescend {
		f.child = newFieldIterator(f.Value())
		f.child.path = append(f.path, f.containerKey())
		return f.Next()
	} else if f.remainingFields() > 0 {
		f.fieldIndex++
		currentVal, k, ok := kindFromValue(f.Value())
		f.shouldDescend = ok
		if k == Map {
			// if key is not a string, we can't descend
			if currentVal.Type().Key().Kind() != reflect.String {
				f.shouldDescend = false
			}
		}
		return true
	} else {
		return false
	}
}

func (f *fieldIterator) IsStructField() bool {
	if f.child != nil {
		return f.child.IsStructField()
	}
	return f.kind == Struct
}

func (f *fieldIterator) IsMap() bool {
	if f.child != nil {
		return f.child.IsMap()
	}
	return f.kind == Map
}

func (f *fieldIterator) Type() reflect.Type {
	return f.Value().Type()
}

func (f *fieldIterator) Value() reflect.Value {
	if f.child != nil {
		return f.child.Value()
	} else if f.kind == Struct {
		return f.container.Field(f.fieldIndex)
	} else {
		v := f.container.Elem().MapIndex(f.keys[f.fieldIndex])
		if v.Kind() == reflect.Interface {
			return v.Elem()
		}
		return v
	}
}

func (f *fieldIterator) Set(val reflect.Value) (err error) {
	if f.child != nil {
		f.child.Set(val)
	} else if f.kind == Struct {
		field := f.container.Field(f.fieldIndex)
		isPointer := field.Kind() == reflect.Ptr
		if isPointer {
			if field.IsNil() {
				// create a new instance and set it
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
			// if val.Kind() != reflect.Ptr {
			// 	if val.CanAddr() {
			// 		val = val.Addr()
			// 	} else {
			// 		return fmt.Errorf("can't address value: %s", val)
			// 	}
			// }
		}
		// fmt.Println("Set", field.CanAddr(), field.CanSet(), field.CanInterface(), isPointer)
		if field.CanInterface() {
			v := field
			if field.CanAddr() {
				v = field.Addr()
			}
			if f.assignViaInterface(v, val) {
				return
			}
		}
		field.Set(val)
	} else {
		f.container.Elem().SetMapIndex(f.keys[f.fieldIndex], val)
	}
	return nil
}

func (f *fieldIterator) assignViaInterface(field, val reflect.Value) (assigned bool) {
	switch k := field.Interface().(type) {
	case StringSet:
		if val.Kind() == reflect.String {
			k.Set(val.String())
			assigned = true
		}
	case StringSetErr:
		if val.Kind() == reflect.String {
			if err := k.Set(val.String()); err != nil {
				panic(err)
			}
			assigned = true
		}
	case IntSet:
		if val.Kind() == reflect.Int {
			k.Set(int(val.Int()))
			assigned = true
		}
	case IntSetErr:
		if val.Kind() == reflect.Int {
			if err := k.Set(int(val.Int())); err != nil {
				panic(err)
			}
			assigned = true
		}
	case FloatSet:
		if val.Kind() == reflect.Float64 {
			k.Set(val.Float())
			assigned = true
		}
	case FloatSetErr:
		if val.Kind() == reflect.Float64 {
			if err := k.Set(val.Float()); err != nil {
				panic(err)
			}
			assigned = true
		}
	case BoolSet:
		if val.Kind() == reflect.Bool {
			k.Set(val.Bool())
			assigned = true
		}
	case BoolSetErr:
		if val.Kind() == reflect.Bool {
			if err := k.Set(val.Bool()); err != nil {
				panic(err)
			}
			assigned = true
		}
	}
	return
}

func (f *fieldIterator) Field() reflect.StructField {
	if f.child != nil {
		return f.child.Field()
	}
	return f.container.Type().Field(f.fieldIndex)
}

func (f *fieldIterator) containerKey() string {
	if f.kind == Map {
		return f.keys[f.fieldIndex].String()
	} else {
		return f.container.Type().Field(f.fieldIndex).Name
	}
}

func (f *fieldIterator) Key() string {
	if f.child != nil {
		return f.child.Key()
	}
	return f.containerKey()
}

func (f *fieldIterator) Path() []string {
	if f.child != nil {
		return f.child.Path()
	}
	return f.path
}

func kindFromValue(v reflect.Value) (newV reflect.Value, kind containerKind, ok bool) {
	switch v.Kind() {
	case reflect.Struct:
		return v, Struct, true
	case reflect.Map:
		return v, Map, true
	case reflect.Interface:
		return kindFromValue(v.Elem())
	default:
		return v, 0, false
	}
}
