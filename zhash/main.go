package zhash

import "errors"

type Alg = uint8

const (
	AlgBcrypt   Alg = 1
	AlgScrypt       = 2
	AlgArgon2i      = 3
	AlgArgon2id     = 4
)

type Algorithm interface {
	Name() string
	Configure(config interface{}) error
	Hash(password []byte) (hash []byte, err error)
	Compare(hash []byte, password []byte) (same bool, err error)
}

var defaultAlgPrecedence = []Alg{AlgArgon2id, AlgArgon2i, AlgScrypt, AlgBcrypt}

var algPrecendence = defaultAlgPrecedence
var available = map[Alg]Algorithm{}

// Hash takes a plain text password and hashes it using the best available
// algorithm. The algorithm used is included at the beginning of the hash result.
func Hash(password []byte) (hash []byte, err error) {
	return HashAlg(password, algPrecendence...)
}

// HashAlg takes a plain text password and hashes using the first available
// algorithm type provided. The algorithm used is included at the beginning of
// the hash result
func HashAlg(password []byte, algs ...Alg) (hash []byte, err error) {
	for _, alg := range algs {
		hasher, ok := available[alg]
		if ok {
			hash, err = hasher.Hash(password)
			if err != nil {
				return
			}
			hash = encodeWithAlg(alg, hash)
			return
		}
	}
	err = errors.New("no hash algorithms available")
	return
}

// Compare takes the hashed password and the plain text password a validates
// they are the same. It assumes the hash was generated using Hash or HashAlg
// which include the hashing algorithm used.
func Compare(hash, password []byte) (equal bool, err error) {
	hash, alg := splitHashAndAlg(hash)
	hasher, ok := available[alg]
	if !ok {
		err = errors.New("invalid hashing function")
		return
	}
	return hasher.Compare(hash, password)
}

// Set the precedence of which algorithms to use for encoding. A single hashing
// algorith can be used if desired
func SetPrecedence(algs ...Alg) {
	algPrecendence = algPrecendence[:len(algs)]
	copy(algPrecendence, algs)
}
