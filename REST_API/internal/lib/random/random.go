package random

import (
	"math/rand/v2"
)

func NewRandomString(size int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rand.IntN(len(chars))]
	}

	return string(b)
}
