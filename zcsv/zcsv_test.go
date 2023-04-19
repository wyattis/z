package zcsv_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/wyattis/z/zcsv"
)

const CSV_ONE = "a,b,c\n1,2,3\n4,5,6"
const CSV_TWO = "username,password\naaaa,qwerty\nbbbb,asdf\n"

type StructOne struct {
	A, B, C string
}
type StructTwo struct {
	Username, Password string
}

var ONE_EXP = []StructOne{{"1", "2", "3"}, {"4", "5", "6"}}
var TWO_EXP = []StructTwo{{"aaaa", "qwerty"}, {"bbbb", "asdf"}}

func TestCsvReadSingleLines(t *testing.T) {
	s := bytes.NewBufferString(CSV_ONE)
	r := zcsv.NewReader(s, nil)

	resOne := []StructOne{}
	for {
		line, err := r.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				t.Error(err)
			}
			break
		}
		resOne = append(resOne, StructOne{
			A: line.MustGet("a"),
			B: line.MustGet("b"),
			C: line.MustGet("c"),
		})
	}
	if !reflect.DeepEqual(resOne, ONE_EXP) {
		t.Errorf("expected %+v, but got %+v", ONE_EXP, resOne)
	}

	s = bytes.NewBufferString(CSV_TWO)
	r = zcsv.NewReader(s, nil)

	resTwo := []StructTwo{}
	for {
		line, err := r.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				t.Error(err)
			}
			break
		}
		resTwo = append(resTwo, StructTwo{
			Username: line.MustGet("username"),
			Password: line.MustGet("password"),
		})
	}
	if !reflect.DeepEqual(resTwo, TWO_EXP) {
		t.Errorf("expected %+v, but got %+v", TWO_EXP, resTwo)
	}
}

func TestCsvReadAll(t *testing.T) {
	s := bytes.NewBufferString(CSV_ONE)
	r := zcsv.NewReader(s, nil)

	resOne := []StructOne{}
	lines, err := r.ReadAll()
	if err != nil {
		t.Error(err)
		return
	}
	for _, l := range lines {
		resOne = append(resOne, StructOne{
			A: l.MustGet("a"),
			B: l.MustGet("b"),
			C: l.MustGet("c"),
		})
	}
	if !reflect.DeepEqual(resOne, ONE_EXP) {
		t.Errorf("expected %+v, but got %+v", ONE_EXP, resOne)
	}

	s = bytes.NewBufferString(CSV_TWO)
	r = zcsv.NewReader(s, nil)

	resTwo := []StructTwo{}
	lines, err = r.ReadAll()
	if err != nil {
		t.Error(err)
		return
	}
	for _, l := range lines {
		resTwo = append(resTwo, StructTwo{
			Username: l.MustGet("username"),
			Password: l.MustGet("password"),
		})
	}
	if !reflect.DeepEqual(resTwo, TWO_EXP) {
		t.Errorf("expected %+v, but got %+v", TWO_EXP, resTwo)
	}
}

func TestScan(t *testing.T) {
	t.Fail()
}

func TestScanAll(t *testing.T) {
	s := bytes.NewBufferString(CSV_ONE)
	r := zcsv.NewReader(s, nil)
	resOne := []StructOne{}
	if err := r.ScanAll(&resOne); err != nil {
		t.Error(err)
		return
	}
}
