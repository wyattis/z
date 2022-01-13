{{- define "slice" -}}

package {{.PackageName}}

import "errors"

var (
	ErrInterfaceNot{{title .TypeName}} = errors.New("encountered non-{{.TypeName}} interface")
)

// Check if a slice ([]{{.Type}}) contains a matching member
func Contains(haystack []{{.Type}}, needle {{.Type}}) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Convert []interface{} into []{{.Type}} if possible
func As(slice []interface{}) (res []{{.Type}}, err error) {
	res = make([]{{.Type}}, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].({{.Type}}); !ok {
			return res, ErrInterfaceNot{{title .TypeName}}
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert []{{.Type}} into []interface{}
func Interface(slice []{{.Type}}) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// func Remove(slice *[]{{.Type}}, values ...{{.Type}}) (res []{{.Type}}) {
// 
// }
// 
// func Replace(slice *[]{{.Type}}, val {{.Type}}, replacement {{.Type}}) (res []{{.Type}}) {
// 
// }
// 
// func ReplaceAll(slice *[]{{.Type}}, val {{.Type}}, replacement {{.Type}}) (res []{{.Type}}) {
// 
// }
// 
// func IndexOf(haystack []{{.Type}}, needle {{.Type}}) int {
// 
// }

{{ end }}