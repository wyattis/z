package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/wyattis/z/zcache"
)

type PrimitiveType struct {
	Type        string
	TypeName    string
	PackageName string
	IsString    bool
}

//go:embed templates
var templates embed.FS

var root = ""
var funcMap = template.FuncMap{
	"title": strings.Title,
}

var sliceTypes = []PrimitiveType{
	{"string", "string", "zstrings", true},
	{"int", "int", "zints", false},
	{"int16", "int16", "zint16s", false},
	{"int32", "int32", "zint32s", false},
	{"int64", "int64", "zint64s", false},
	{"uint", "uint", "zuints", false},
	{"uint32", "uint32", "zuint32s", false},
	{"uint64", "uint64", "zuint64s", false},
	{"float32", "float32", "zfloat32s", false},
	{"float64", "float64", "zfloat64s", false},
}

var setTypes = []PrimitiveType{
	{"string", "string", "zstringset", true},
	{"int", "int", "zintset", false},
	{"int16", "int16", "zint16set", false},
	{"int32", "int32", "zint32set", false},
	{"int64", "int64", "zint64set", false},
	{"uint", "uint", "zuintset", false},
	{"uint32", "uint32", "zuint32set", false},
	{"uint64", "uint64", "zuint64set", false},
	{"float32", "float32", "zfloat32set", false},
	{"float64", "float64", "zfloat64set", false},
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
