package random

import "math/rand"

var lettersForRndString = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewRandomString(stringLength int) string {
	rndString := make([]rune, stringLength)

	for i := range rndString {
		rndString[i] = lettersForRndString[rand.Intn(len(lettersForRndString))]
	}

	return string(rndString)
}
