package zhash

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/scrypt"
)

func init() {
	available[AlgScrypt] = &Scrypt{
		Config: DefaultScryptConfig,
	}
}

var DefaultScryptConfig = ScryptConfig{
	N:      32768,
	R:      8,
	P:      1,
	KeyLen: 32,
	Salt:   DefaultSaltConfig,
}

type ScryptConfig struct {
	N      int
	R      int
	P      int
	KeyLen int
	Salt   SaltConfig
}

type Scrypt struct {
	Config ScryptConfig
}

func (h Scrypt) Name() string {
	return "scrypt"
}

func (h *Scrypt) Configure(config interface{}) error {
	c, ok := config.(ScryptConfig)
	if !ok {
		return errors.New("expected ScryptConfig")
	}
	h.Config = c
	return nil
}

func (h Scrypt) Hash(pass []byte) (hash []byte, err error) {
	salt, err := generateSalt(h.Config.Salt)
	if err != nil {
		return
	}
	return h.HashSalt(pass, salt)
}

func (h Scrypt) HashSalt(pass, salt []byte) (hash []byte, err error) {
	hash, err = scrypt.Key(pass, salt, h.Config.N, h.Config.R, h.Config.P, h.Config.KeyLen)
	if err != nil {
		return
	}
	hash = encodeWithSalt(hash, salt)
	return
}

func (h *Scrypt) Compare(hash, pass []byte) (equal bool, err error) {
	_, salt := splitHashAndSalt(hash)
	phash, err := h.HashSalt(pass, salt)
	if err != nil {
		return
	}
	return bytes.Equal(hash, phash), nil
}
