package zreflect

type StringSetErr interface {
	Set(string) (err error)
}

type StringSet interface {
	Set(string)
}

type IntSetErr interface {
	Set(int) (err error)
}

type IntSet interface {
	Set(int)
}

type FloatSetErr interface {
	Set(float64) (err error)
}

type FloatSet interface {
	Set(float64)
}

type BoolSetErr interface {
	Set(bool) (err error)
}

type BoolSet interface {
	Set(bool)
}
