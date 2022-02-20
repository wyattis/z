package zhash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var DefaultBcryptConfig = BcryptConfig{
	Cost: bcrypt.DefaultCost,
}

func init() {
	available[AlgBcrypt] = &Bcrypt{
		config: DefaultBcryptConfig,
	}
}

type BcryptConfig struct {
	Cost int
}

type Bcrypt struct {
	config BcryptConfig
}

func (h *Bcrypt) Name() string {
	return "bcrypt"
}

func (h *Bcrypt) Configure(config interface{}) error {
	c, ok := config.(BcryptConfig)
	if !ok {
		return errors.New("expected BcryptConfig")
	}
	h.config = c
	return nil
}

func (h Bcrypt) Hash(pass []byte) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword(pass, h.config.Cost)
}

func (h Bcrypt) Compare(hash, pass []byte) (equal bool, err error) {
	err = bcrypt.CompareHashAndPassword(hash, pass)
	return err == nil, nil
}
