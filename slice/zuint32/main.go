// Code generated by z. DO NOT EDIT.

package zuint32

import "errors"

var (
	ErrInterfaceNotUint32 = errors.New("encountered non-uint32 interface")
)

// Check if a slice ([]uint32) contains a matching member
func Contains(haystack []uint32, needle uint32) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

// Convert []interface{} into []uint32 if possible
func As(slice []interface{}) (res []uint32, err error) {
	res = make([]uint32, len(slice))
	for i := range slice {
		if strVal, ok := slice[i].(uint32); !ok {
			return res, ErrInterfaceNotUint32
		} else {
			res[i] = strVal
		}
	}
	return
}

// Convert []uint32 into []interface{}
func Interface(slice []uint32) (res []interface{}, err error) {
	res = make([]interface{}, len(slice))
	for i := range slice {
		res[i] = slice[i]
	}
	return
}

// func Remove(slice *[]uint32, values ...uint32) (res []uint32) {
// 
// }
// 
// func Replace(slice *[]uint32, val uint32, replacement uint32) (res []uint32) {
// 
// }
// 
// func ReplaceAll(slice *[]uint32, val uint32, replacement uint32) (res []uint32) {
// 
// }
// 
// func IndexOf(haystack []uint32, needle uint32) int {
// 
// }

