package zenv

import (
	"errors"
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

// Set non-zero values on a struct using reflection to parse values from
// os.Environ() followed by parsing values from any provided .env files.
// Existing values will not be overwritten if the same key is encountered later.
// If a file doesn't exist it will be skipped.
func Set(val interface{}, paths ...string) error {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
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
	return SetMap(reflect.Indirect(reflect.ValueOf(val)), env)
}

// Set non-zero values on a struct using reflect to parse values from .env files.
// This method has the same behavior as Set, but doesn't use os.Environ() at all.
func SetFiles(val interface{}, paths ...string) error {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	env, err := ParseFiles(paths...)
	if err != nil {
		return err
	}
	return SetMap(reflect.Indirect(reflect.ValueOf(val)), env)
}

// Parse multiple env files and merge the results into a single map. Existing
// values are skipped when merging. This defines the precedence as left to right.
func ParseFiles(paths ...string) (env EnvMap, err error) {
	env = make(EnvMap)
	for _, p := range paths {
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
func SetMap(val reflect.Value, env EnvMap) (err error) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.IsZero() {
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
			if val, exists := env[name]; exists {
				if field.Type() == reflect.TypeOf(time.Time{}) {
					t, err := ztime.Parse(val, t.Tag.Get("time-format"))
					if err != nil {
						return err
					}
					field.Set(reflect.ValueOf(t))
				} else if k == reflect.Struct {
					if err = SetMap(field, env); err != nil {
						return
					}
				} else if converter, exists := zreflect.ConvMap[k]; exists {
					rVal, err := converter(val, field, zreflect.ConvMap)
					if err != nil {
						return err
					}
					field.Set(rVal)
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
