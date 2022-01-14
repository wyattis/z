# zcache
zcache is a naïve LRU file caching solution w/ optional support for file 
compression.

```go
  cache, err := zcache.New(zcache.DirFs("data"), 5000)
  if err != nil {
    panic(err)
  }
  cache.Compressed = true // set before init
  if err = cache.Init(); err != nil {
    panic(err)
  }
  reader, item, _ := cache.OpenOrCreate("my-key", func (w io.Writer) error {
    _, err := w.Write([]byte("helloski\n"))
    return err
  })
  msg, _ := io.ReadAll(reader)
  fmt.Println(string(msg), item)
```