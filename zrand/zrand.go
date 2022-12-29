package zrand

import (
	"errors"
	"math/rand"
	"time"
)

const ALPHA = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const NUMERIC = "0123456789"
const ALPHANUMERIC = ALPHA + NUMERIC

var (
	ErrLengthExceedsAvailable = errors.New("requested length exceeds length of set")
)

func SubsetFromSource[T any](source rand.Source, set []T, length int) (res []T) {
	res = make([]T, length)
	r := rand.New(source)
	for i := range res {
		res[i] = set[r.Intn(len(set))]
	}
	return res
}

func Subset[T any](set []T, length int) (res []T) {
	res = make([]T, length)
	rand.Seed(time.Now().UnixNano())
	for i := range res {
		res[i] = set[rand.Intn(len(set))]
	}
	return res
}

func UniqueSubset[T any](set []T, length int) (res []T) {
	if len(set) < length {
		panic(ErrLengthExceedsAvailable)
	}
	res = make([]T, length)
	for i := range res {
		j := rand.Intn(len(set))
		set[0], set[j] = set[j], set[0]
		res[i], set = set[j], set[1:]
	}
	return res
}

func AlphaWord(length int) string {
	return string(Subset([]rune(ALPHA), length))
}

func AlphaWordFromSource(source rand.Source, length int) string {
	return string(SubsetFromSource(source, []rune(ALPHA), length))
}

func AlphaNumericWord(length int) string {
	return string(Subset([]rune(ALPHANUMERIC), length))
}

func AlphaNumericWordFromSource(source rand.Source, length int) string {
	return string(SubsetFromSource(source, []rune(ALPHANUMERIC), length))
}
