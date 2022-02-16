package helpers

import (
	"math/rand"
)

// The random runes to use when generating strings.
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// Generates a random alphanumeric string of
// the given length.
//
// The included characters are a-z, A-Z, 0-9.
func RandomAlphaString(length int) string {
	b := make([]rune, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
