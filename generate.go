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
		{"string", "string", "zstring"},
		{"int", "int", "zint"},
		{"int16", "int16", "zint16"},
		{"int32", "int32", "zint32"},
		{"int64", "int64", "zint64"},
		{"uint", "uint", "zuint"},
		{"uint32", "uint32", "zuint32"},
		{"uint64", "uint64", "zuint64"},
		{"float32", "float32", "zfloat32"},
		{"float64", "float64", "zfloat64"},
	}
	tmpl := template.New("").Funcs(funcMap)
	tmpl, err = tmpl.ParseFS(templates, "**/*.tpl")
	if err != nil {
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
