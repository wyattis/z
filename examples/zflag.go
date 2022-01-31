package main

import (
	"flag"
	"time"

	"github.com/wyattis/z/zflag"
)

type FileServerConfig struct {
	Addr    string        `flag:",:80,the address to use for this server"`
	Root    string        `flag:"root,,which root directory to serve files from"`
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
