package zhash

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/argon2"
)

func init() {
	available[AlgArgon2i] = &Argon2i{
		config: DefaultArgon2Config,
	}
}

type Argon2Config struct {
	Salt    SaltConfig
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

var DefaultArgon2Config = Argon2Config{
	Salt:    DefaultSaltConfig,
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

type Argon2i struct {
	config Argon2Config
}

func (h *Argon2i) Name() string {
	return "argon2i"
}

func (h *Argon2i) Configure(config interface{}) error {
	c, ok := config.(Argon2Config)
	if !ok {
		return errors.New("invalid config. expected Argon2Config")
	}
	h.config = c
	return nil
}

func (h Argon2i) Hash(pass []byte) (hash []byte, err error) {
	salt, err := generateSalt(h.config.Salt)
	if err != nil {
		return
	}
	return h.HashSalt(pass, salt)
}

func (h Argon2i) HashSalt(pass, salt []byte) (hash []byte, err error) {
	hash = argon2.Key(pass, salt, h.config.Time, h.config.Memory, h.config.Threads, h.config.KeyLen)
	hash = encodeWithSalt(hash, salt)
	return
}

func (h Argon2i) Compare(hash, pass []byte) (equal bool, err error) {
	// fmt.Println("input", hash, pass)
	_, salt := splitHashAndSalt(hash)
	phash, err := h.HashSalt(pass, salt)
	if err != nil {
		return
	}
	// fmt.Println("salt", salt)
	// fmt.Println("hash", hash)
	// fmt.Println("phash", phash)
	equal = bytes.Equal(hash, phash)
	return
}
