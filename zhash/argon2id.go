package zhash

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/argon2"
)

func init() {
	available[AlgArgon2id] = &Argon2id{
		config: DefaultArgon2Config,
	}
}

type Argon2id struct {
	config Argon2Config
}

func (h *Argon2id) Name() string {
	return "argon2id"
}

func (h *Argon2id) Configure(config interface{}) error {
	c, ok := config.(Argon2Config)
	if !ok {
		return errors.New("invalid config. expected Argon2Config")
	}
	h.config = c
	return nil
}

func (h Argon2id) Hash(pass []byte) (hash []byte, err error) {
	salt, err := generateSalt(h.config.Salt)
	if err != nil {
		return
	}
	return h.HashSalt(pass, salt)
}

func (h Argon2id) HashSalt(pass, salt []byte) (hash []byte, err error) {
	hash = argon2.IDKey(pass, salt, h.config.Time, h.config.Memory, h.config.Threads, h.config.KeyLen)
	hash = encodeWithSalt(hash, salt)
	return
}

func (h Argon2id) Compare(hash, pass []byte) (equal bool, err error) {
	_, salt := splitHashAndSalt(hash)
	phash, err := h.HashSalt(pass, salt)
	if err != nil {
		return
	}
	equal = bytes.Equal(hash, phash)
	return
}
