package zhash

import "crypto/rand"

type SaltConfig struct {
	Length int
}

var DefaultSaltConfig = SaltConfig{
	Length: 16,
}

func generateSalt(config SaltConfig) (salt []byte, err error) {
	salt = make([]byte, config.Length)
	_, err = rand.Read(salt)
	if err != nil {
		return
	}
	return
}
