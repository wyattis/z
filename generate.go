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

func generate() (err error) {
	sliceTypes := []PrimitiveType{
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
	tmpl := template.New("").Funcs(funcMap)
	tmpl, err = tmpl.ParseFS(templates, "**/*.tpl")
	if err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join(root, "slice/")); err != nil {
		return
	}
	for _, t := range sliceTypes {
		outputPath := filepath.Join(root, "slice/", t.PackageName, "main.go")
		if err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return
		}
		if err = executeTemplate(outputPath, tmpl, "slice", t); err != nil {
			return
		}
	}

	return
}

func main() {
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}
