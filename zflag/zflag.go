package zflag

import (
	"flag"

	"github.com/wyattis/z/ztime"
)

func Parse() {
	flag.Parse()
}

func DurationVar(v *ztime.Duration, name string, defaultVal ztime.Duration, description string) {
	if defaultVal != 0 {
		if err := v.Set(defaultVal.String()); err != nil {
			panic(err)
		}
	}
	flag.Var(v, name, description)
}

func Duration(name string, defaultVal ztime.Duration, description string) *ztime.Duration {
	d := new(ztime.Duration)
	if defaultVal != 0 {
		if err := d.Set(defaultVal.String()); err != nil {
			panic(err)
		}
	}
	flag.Var(d, name, description)
	return d
}
