package ztemplate

import (
	"bytes"
	"text/template"
)

/*
 * Quickly execute a single template string with "text/template"
 */
func ExecString(tmpl string, data any) (res string, err error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return
	}
	buf := bytes.Buffer{}
	if err = t.Execute(&buf, data); err != nil {
		return
	}
	res = buf.String()
	return
}
