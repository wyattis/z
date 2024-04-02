package zconfig

import (
	"flag"
	"os"

	"github.com/wyattis/z/zdefaults"
	"github.com/wyattis/z/zenv"
	"github.com/wyattis/z/zflag"
	"github.com/wyattis/z/zos"
	"gopkg.in/yaml.v3"
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
func (c *Configurer) Apply(val interface{}) (err error) {
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
	opts := []configOption{Defaults(), Env(".env"), Flag(os.Args)}
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
		return zenv.Set(val, &zenv.EnvOptions{Overwrite: false}, paths...)
	}
}

// Configure a struct using values from .env files
func EnvFiles(paths ...string) configOption {
	return func(val interface{}) error {
		return zenv.SetFiles(val, &zenv.EnvOptions{Overwrite: false}, paths...)
	}
}

// Configure a struct using command line flags
func Flag(args []string) configOption {
	return func(val interface{}) error {
		set := flag.NewFlagSet("", flag.ContinueOnError)
		if err := zflag.Configure(set, val, &zflag.FlagOptions{Overwrite: false}); err != nil {
			return err
		}
		return set.Parse(args)
	}
}

// Configure an interface{} using a yaml file
func Yaml(paths ...string) configOption {
	return func(val interface{}) (err error) {
		for _, loc := range paths {
			if zos.Exists(loc) {
				f, err := os.Open(loc)
				if err != nil {
					return err
				}
				defer f.Close()
				dec := yaml.NewDecoder(f)
				return dec.Decode(val)
			}
		}
		return nil
	}
}

// Configure using default tags configured on the struct
func Defaults() configOption {
	return func(val interface{}) error {
		return zdefaults.Set(val)
	}
}
