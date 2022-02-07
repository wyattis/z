package zconfig

import (
	"flag"

	"github.com/wyattis/z/zdefaults"
	"github.com/wyattis/z/zenv"
	"github.com/wyattis/z/zflag"
)

type configOption = func(val interface{}) error

type Configurer struct {
	options []configOption
}

// Create a configurer
func New(options ...configOption) *Configurer {
	return &Configurer{
		options: options,
	}
}

// Actually update the struct given the values provided
func (c *Configurer) Parse(val interface{}) (err error) {
	options := c.options
	if len(options) == 0 {
		options = append(options, Auto())
	}
	for _, opt := range c.options {
		if err = opt(val); err != nil {
			return
		}
	}
	return
}

// Configure a struct using any command line flags, environment variables,
// variables in .env and finally, values set using default tags
func Auto() configOption {
	opts := []configOption{Flag(), Env(".env"), Defaults()}
	return func(val interface{}) error {
		for _, opt := range opts {
			if err := opt(val); err != nil {
				return err
			}
		}
		return nil
	}
}

// Configure a struct using values defined in the environment and .env files
func Env(paths ...string) configOption {
	return func(val interface{}) error {
		return zenv.Set(val, paths...)
	}
}

// Configure a struct using values from .env files
func EnvFile(paths ...string) configOption {
	return func(val interface{}) error {
		return zenv.SetFiles(val, paths...)
	}
}

// Configure a struct using command line flags
func Flag(args []string) configOption {
	return func(val interface{}) error {
		set := flag.NewFlagSet("", flag.ExitOnError)
		if err := zflag.Configure(set, val); err != nil {
			return err
		}
		return set.Parse(args)
	}
}

// Configure using default tags configured on the struct
func Defaults() configOption {
	return func(val interface{}) error {
		return zdefaults.Set(val)
	}
}
