// Package randomizer is used for generating random strings.
package randomizer

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// StringRunes generates a random string of runes of a specified length.
func StringRunes(length int) string {
	return generateRunes(length)
}

// Randomizer is used for generating random strings.
type Randomizer struct{}

// StringRunes generates a random string of runes of a specified length from a Randomizer struct.
func (r Randomizer) StringRunes(length int) string {
	return generateRunes(length)
}

func generateRunes(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
