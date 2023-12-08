package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type PrimitiveType struct {
	Type        string
	TypeName    string
	PackageName string
	IsString    bool
	IsBytes     bool
	IsBool      bool
}

//go:embed templates
var templates embed.FS

var root = "../"
var funcMap = template.FuncMap{
	"title": strings.Title,
}

var sliceTypes = []PrimitiveType{
	{"string", "string", "zstrings", true, false, false},
	{"byte", "byte", "zbytes", false, true, false},
	{"bool", "bool", "zbools", false, false, true},
	{"int", "int", "zints", false, false, false},
	{"int16", "int16", "zint16s", false, false, false},
	{"int32", "int32", "zint32s", false, false, false},
	{"int64", "int64", "zint64s", false, false, false},
	{"uint", "uint", "zuints", false, false, false},
	{"uint32", "uint32", "zuint32s", false, false, false},
	{"uint64", "uint64", "zuint64s", false, false, false},
	{"float32", "float32", "zfloat32s", false, false, false},
	{"float64", "float64", "zfloat64s", false, false, false},
}

var setTypes = []PrimitiveType{
	{"string", "string", "zstringset", true, false, false},
	{"byte", "byte", "zbytes", false, true, false},
	{"int", "int", "zintset", false, false, false},
	{"int16", "int16", "zint16set", false, false, false},
	{"int32", "int32", "zint32set", false, false, false},
	{"int64", "int64", "zint64set", false, false, false},
	{"uint", "uint", "zuintset", false, false, false},
	{"uint32", "uint32", "zuint32set", false, false, false},
	{"uint64", "uint64", "zuint64set", false, false, false},
	{"float32", "float32", "zfloat32set", false, false, false},
	{"float64", "float64", "zfloat64set", false, false, false},
}

func executeTemplate(out string, tmpl *template.Template, name string, data interface{}) (err error) {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = tmpl.ExecuteTemplate(f, "generated", nil); err != nil {
		return err
	}
	if err = tmpl.ExecuteTemplate(f, name, data); err != nil {
		return err
	}
	return
}

func makeTypes(dir, tempName string, types []PrimitiveType) (err error) {
	tmpl := template.New("").Funcs(funcMap)
	tmpl, err = tmpl.ParseFS(templates, "**/*.tpl")
	if err != nil {
		return
	}
	for _, t := range types {
		outputPath, err := filepath.Abs(filepath.Join(dir, t.PackageName, t.PackageName+".go"))
		if err != nil {
			return err
		}
		fmt.Println("generating into", outputPath)
		if err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}
		if err = executeTemplate(outputPath, tmpl, tempName, t); err != nil {
			return err
		}
	}

	return
}

func generate() (err error) {
	if err = makeTypes(filepath.Join(root, "zslice"), "slice", sliceTypes); err != nil {
		return
	}
	return makeTypes(filepath.Join(root, "zset"), "set", setTypes)
}

func main() {
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}
