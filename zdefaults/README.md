# zdefaults
Set default values for a struct using tags. Will ignore any values that aren't already at their zero value. Supports parsing time.Duration and time.Time as well as custom implementations of the `Settable` interface.

## Example
Set defaults for an example file server configuration.
```go
type Config struct {
  Addr  string    `default:":80"`
  Roots []string  `default:"/var/www,/home/www"`
  Ints  []int     // ignored
}

func main () {
  conf := Config{}
  err = zdefaults.SetDefaults(&conf)
}
```

### time.Time
```go
type Config struct {
  Time          time.Time     `default:"2009-01-20"`
  CustomTime    time.Time     `default:"random-20-05-21" time-format:"random-06-01-02"`
  CacheDuration time.Duration `default:"10s"`
}

func main () {
  conf := FileServerConfig{}
  err = zdefaults.SetDefaults(&conf)
}
```