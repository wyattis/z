package zenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/wyattis/z/zos"
	"github.com/wyattis/z/zreflect"
	"github.com/wyattis/z/zstring"
	"github.com/wyattis/z/ztime"
)

type EnvMap = map[string]string

type EnvOptions struct {
	Overwrite bool
}

var defaultConfig = EnvOptions{
	Overwrite: false,
}

// Set non-zero values on a struct using reflection to parse values from
// os.Environ() followed by parsing values from any provided .env files.
// Existing values will not be overwritten if the same key is encountered later.
// If a file doesn't exist it will be skipped.
func Set(val interface{}, config *EnvOptions, paths ...string) error {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	if config == nil {
		config = &defaultConfig
	}
	env := make(EnvMap)
	for _, pair := range os.Environ() {
		key, val, _ := zstring.Cut(pair, "=")
		env[key] = val
	}
	filesEnv, err := ParseFiles(paths...)
	if err != nil {
		return err
	}
	mergeEnv(&env, filesEnv)
	return SetMap(reflect.Indirect(reflect.ValueOf(val)), *config, env, "")
}

// Set non-zero values on a struct using reflect to parse values from .env files.
// This method has the same behavior as Set, but doesn't use os.Environ() at all.
func SetFiles(val interface{}, config *EnvOptions, paths ...string) error {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	env, err := ParseFiles(paths...)
	if err != nil {
		return err
	}
	fmt.Println(env)
	return SetMap(reflect.Indirect(reflect.ValueOf(val)), *config, env, "")
}

// Parse multiple env files and merge the results into a single map. Existing
// values are skipped when merging. This defines the precedence as left to right.
func ParseFiles(paths ...string) (env EnvMap, err error) {
	env = make(EnvMap)
	for _, p := range paths {
		fmt.Println("checking", p)
		if zos.Exists(p) {
			fileEnv, err := ParseEnvFile(p)
			if err != nil {
				return env, err
			}
			mergeEnv(&env, fileEnv)
		}
	}
	return
}

// Parse a single env file using godotenv.Parse
func ParseEnvFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return godotenv.Parse(f)
}

// Set non-zero values on a struct using reflect to get values from a map
func SetMap(val reflect.Value, config EnvOptions, env EnvMap, prefix string) (err error) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		t := val.Type().Field(i)
		k := field.Kind()
		name := t.Tag.Get("env")
		if name == "" {
			if !zstring.IsUppercase(t.Name) {
				name = strings.ToUpper(zstring.CamelToSnake(t.Name, "_", 1))
			} else {
				name = t.Name
			}
		}
		if prefix != "" {
			name = strings.Join([]string{prefix, name}, "_")
		}
		fmt.Println("zenv.SetMap", t.Name, name, field.IsZero())
		if k == reflect.Struct {
			if err = SetMap(field, config, env, name); err != nil {
				return
			}
		} else if field.IsZero() || config.Overwrite {
			fmt.Println("zenv.SetMap", t.Name, name)
			if val, exists := env[name]; exists {
				if field.Type() == reflect.TypeOf(time.Time{}) {
					t, err := ztime.Parse(val, t.Tag.Get("time-format"))
					if err != nil {
						return err
					}
					field.Set(reflect.ValueOf(t))
				} else if converter, exists := zreflect.ConvMap[k]; exists {
					fmt.Println("converting", k, t.Name, val)
					rVal, err := converter(val, field, zreflect.ConvMap)
					if err != nil {
						return err
					}
					fmt.Println("setting", k, t.Name, rVal)
					field.Set(rVal)
				} else {
					return fmt.Errorf("unknown type %s %s", name, k)
				}
			}
		}
	}
	return
}

func mergeEnv(dest *EnvMap, source EnvMap) {
	for key, val := range source {
		if v, exists := (*dest)[key]; !exists || v == "" {
			(*dest)[key] = val
		}
	}
}
