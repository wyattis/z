package zreflect

import (
	"fmt"
	"reflect"
	"testing"
)

type singleNestedStruct struct {
	A struct {
		B struct {
			C int
		}
	}
}

func TestFieldIteratorKeys(t *testing.T) {
	expected := []string{"A", "B", "C"}
	keys := []string{}
	val := singleNestedStruct{}
	fi := FieldIterator(&val)
	for fi.Next() {
		keys = append(keys, fi.Key())
	}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

type complexStruct struct {
	A int
	B string
	C []int
	D map[string]int
	E struct {
		F int
	}
}

func TestFieldIterator(t *testing.T) {
	expected := []string{"A", "B", "C", "D", "E", "F"}
	keys := []string{}
	val := complexStruct{}
	fi := FieldIterator(&val)
	for fi.Next() {
		keys = append(keys, fi.Key())
	}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

type mapTestStruct struct {
	G map[string]struct {
		H int
	}
}

func TestFieldIteratorMapStruct(t *testing.T) {
	expectedKeys := []string{"G", "I", "H", "J", "H"}
	keys := []string{}
	val := mapTestStruct{
		G: map[string]struct {
			H int
		}{
			"I": {H: 1},
			"J": {H: 2},
		},
	}
	fi := FieldIterator(&val)
	for fi.Next() {
		keys = append(keys, fi.Key())
	}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Expected %v, got %v", expectedKeys, keys)
	}
}

type mapTestInterface map[string]any

func TestFieldIteratorMapInterface(t *testing.T) {
	expectedKeys := []string{"G", "I", "H", "J", "H", "K"}
	val := mapTestInterface{
		// "G": map[string]any{
		// 	"I": map[string]any{
		// 		"H": 1,
		// 	},
		// 	"J": map[string]any{
		// 		"H": 2,
		// 	},
		// },
		"K": "none",
	}
	keys := []string{}
	fi := FieldIterator(&val)
	for fi.Next() {
		keys = append(keys, fi.Key())
		fmt.Println(keys, fi.container.Elem().CanSet(), fi.Value().CanSet())
	}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Expected %v, got %v", expectedKeys, keys)
	}
}

type testIntStruct struct {
	H int
}

func TestMapSet(t *testing.T) {
	val := map[string]any{
		"K": "none",
	}
	fi := FieldIterator(&val)
	for fi.Next() {
		fi.Set(reflect.ValueOf("many"))
	}
	if val["K"] != "many" {
		t.Errorf("Expected %v, got %v", "many", val["K"])
	}
	val = map[string]any{
		"K": testIntStruct{},
	}
	fi = FieldIterator(&val)
	for fi.Next() {
		if !fi.IsStructField() {
			continue
		}
		fi.Set(reflect.ValueOf(10))
	}
	expected := val["K"].(*testIntStruct).H
	if expected != 10 {
		t.Errorf("Expected %v, got %v", 10, expected)
	}
}
