package main

import (
	"embed"
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
}

//go:embed templates
var templates embed.FS

var root = ""
var funcMap = template.FuncMap{
	"title": strings.Title,
}

var sliceTypes = []PrimitiveType{
	{"string", "string", "zstrings"},
	{"int", "int", "zints"},
	{"int16", "int16", "zint16s"},
	{"int32", "int32", "zint32s"},
	{"int64", "int64", "zint64s"},
	{"uint", "uint", "zuints"},
	{"uint32", "uint32", "zuint32s"},
	{"uint64", "uint64", "zuint64s"},
	{"float32", "float32", "zfloat32s"},
	{"float64", "float64", "zfloat64s"},
}

var setTypes = []PrimitiveType{
	{"string", "string", "zstringset"},
	{"int", "int", "zintset"},
	{"int16", "int16", "zint16set"},
	{"int32", "int32", "zint32set"},
	{"int64", "int64", "zint64set"},
	{"uint", "uint", "zuintset"},
	{"uint32", "uint32", "zuint32set"},
	{"uint64", "uint64", "zuint64set"},
	{"float32", "float32", "zfloat32set"},
	{"float64", "float64", "zfloat64set"},
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
	if err = os.RemoveAll(dir); err != nil {
		return
	}
	for _, t := range types {
		outputPath := filepath.Join(dir, t.PackageName, "main.go")
		if err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return
		}
		if err = executeTemplate(outputPath, tmpl, tempName, t); err != nil {
			return
		}
	}

	return
}

func generate() (err error) {
	if err = makeTypes(filepath.Join(root, "slice"), "slice", sliceTypes); err != nil {
		return
	}
	return makeTypes(filepath.Join(root, "set"), "set", setTypes)
}

func main() {
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}
