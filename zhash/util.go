package zhash

import "fmt"

func encodeWithSalt(hash, salt []byte) []byte {
	saltLen := uint8(len(salt))
	res := make([]byte, len(hash)+len(salt)+1)
	res[0] = saltLen
	copy(res[1:len(salt)+1], salt)
	copy(res[1+len(salt):], hash)
	return res
}

func splitHashAndSalt(raw []byte) (hash, salt []byte) {
	if len(raw) == 0 {
		fmt.Println(raw)
		panic("hash cannot have len(0)")
	}
	saltLen := raw[0]
	if len(raw) < int(saltLen+1) {
		panic("incorrectly encoded hash. len(hash + salt) < (expected salt length)")
	}
	salt = raw[1 : 1+saltLen]
	hash = raw[1+saltLen:]
	return
}

func encodeWithAlg(alg Alg, hash []byte) []byte {
	res := make([]byte, len(hash)+1)
	res[0] = alg
	copy(res[1:], hash)
	return res
}

func splitHashAndAlg(raw []byte) (hash []byte, alg Alg) {
	return hash[1:], hash[0]
}
