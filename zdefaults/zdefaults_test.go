package zdefaults

import (
	"reflect"
	"testing"
	"time"
)

type simpleconf struct {
	Int        int           `default:"1"`
	Int8       int8          `default:"2"`
	Int16      int16         `default:"3"`
	Int32      int32         `default:"4"`
	Int64      int64         `default:"5"`
	Uint       uint          `default:"6"`
	Uint8      uint8         `default:"7"`
	Uint16     uint16        `default:"8"`
	Uint32     uint32        `default:"9"`
	Uint64     uint64        `default:"10"`
	Float32    float32       `default:"11"`
	Float64    float64       `default:"12"`
	Bool1      bool          `default:"true"`
	Bool2      bool          `default:"1"`
	String     string        `default:"hello"`
	Dur        time.Duration `default:"60s"`
	Time       time.Time     `default:"2000-01-01"`
	TimeCustom time.Time     `default:"Feb 4, 2008" time-format:"Jan _2, 2006"`
}

type sliceconf struct {
	Strings  []string  `default:"hello,world"`
	Ints     []int     `default:"1,2,3"`
	Int8s    []int8    `default:"1,2,3"`
	Float32s []float32 `default:"10.1,1234.1234,56"`
}

type structconf struct {
	SliceConf  sliceconf
	SimpleConf simpleconf
}

var t1, _ = time.Parse("2006-01-02", "2000-01-01")
var t2, _ = time.Parse("2006-01-02", "2008-02-04")

var simpleCases = [][2]simpleconf{
	{
		simpleconf{},
		simpleconf{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, true, true, "hello", time.Minute, t1, t2},
	},
	{
		simpleconf{Int: -1, Int8: -1},
		simpleconf{-1, -1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, true, true, "hello", time.Minute, t1, t2},
	},
}

var sliceCases = [][2]sliceconf{
	{
		{},
		{[]string{"hello", "world"}, []int{1, 2, 3}, []int8{1, 2, 3}, []float32{10.1, 1234.1234, 56}},
	},
	{
		{Strings: []string{}},
		{[]string{}, []int{1, 2, 3}, []int8{1, 2, 3}, []float32{10.1, 1234.1234, 56}},
	},
}

var structCases = [][2]structconf{
	{
		{},
		{
			SliceConf:  sliceconf{[]string{"hello", "world"}, []int{1, 2, 3}, []int8{1, 2, 3}, []float32{10.1, 1234.1234, 56}},
			SimpleConf: simpleconf{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, true, true, "hello", time.Minute, t1, t2},
		},
	},
}

func TestSimple(t *testing.T) {
	for i, c := range simpleCases {
		in := c[0]
		if err := Set(&in); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(in, c[1]) {
			t.Errorf("Expected %+v, but got %+v", c[1], c[0])
		} else {
			t.Logf("passed simple %d", i)
		}
	}
}

func TestArray(t *testing.T) {
	for i, c := range sliceCases {
		in := c[0]
		if err := Set(&in); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(in, c[1]) {
			t.Errorf("Expected %+v, but got %+v", c[1], c[0])
		} else {
			t.Logf("passed complex %d", i)
		}
	}
}

func TestStruct(t *testing.T) {
	for i, c := range structCases {
		in := c[0]
		if err := Set(&in); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(in, c[1]) {
			t.Errorf("Expected %+v, but got %+v", c[1], c[0])
		} else {
			t.Logf("passed struct %d", i)
		}
	}
}
