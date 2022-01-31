# zflag
Convert structs into flag.FlagSet arguments using reflection.

## Format
The flag tag uses the same order of arguments to define a flag using the library. Omitted values are ignored and trailing commas aren't necessary. Usage strings can also contain commas.

|                       |                                    |
|-----------------------|------------------------------------|
| Everything            | `flag:"<name>,<default>,<usage>"`  |
| Just name             | `flag:"<name>"`                    |
| Just name and default | `flag:"<name>,<default>"`          |
| Name and usage        | `flag:"<name>,,<usage>"`           |

## Example
```go
package main

import (
	"flag"
	"time"

	"github.com/wyattis/z/zflag"
)

type FileServerConfig struct {
	Addr    string        `flag:",:80,the address to use for this server"`
	Root    string        `flag:"dir,,which root directory to serve files from"`
	Timeout time.Duration `flag:",1s,how long to wait before timing out a request"`
	Tls     bool          `flag:",,start a TLS server"`
}

func main() {
	config := FileServerConfig{}
	set := flag.NewFlagSet("", flag.ExitOnError)
	if err := zflag.ReflectStruct(set, &config); err != nil {
		panic(err)
	}
	set.Usage()
}
```

#### Output
```
Usage:
  -addr string
        the address to use for this server (default ":80")
  -root string
        which root directory to serve files from
  -timeout duration
        how long to wait before timing out a request (default 1s)
  -tls
        start a TLS server
```