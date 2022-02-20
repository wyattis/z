# zhash
Abstracts away four password hashing algorithms for simple hashing and comparing of passwords. 
Future-proof design will make adding future hashing algorithms easy and shouldn't change existing 
APIs.

## Algorithm precedence
1. argon2id
2. argon2i
3. bcrypt
4. scrypt


## Install
```
go get github.com/wyattis/z/zhash
```

## Usage
```go
import "github.com/wyattis/z/zhash"

func main () {
  pass := []byte("example")
  hash, err := zhash.Hash(pass)
  if err != nil {
    panic(err)
  }
  fmt.Println(hash)
  ok, err := zhash.Compare(hash, pass)
  if err != nil {
    panic(err)
  }
  fmt.Println(ok)
}

// Output:
// [4 16 254 147 73 103 170 201 27 248 107 164 137 117 218 220 187 122 53 119 23 231 106 80 240 181 106 151 31 133 63 241 111 48 146 99 102 42 16 70 77 18 1 91 113 246 15 55 46 21]
// true
```
