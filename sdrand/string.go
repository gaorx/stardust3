package sdrand

import (
	"math/rand"
)

var (
	LowerLetters       = "abcdefghijklmnopqrstuvwxyz"
	UpperLetters       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Letters            = LowerLetters + UpperLetters
	Numbers            = "0123456789"
	LettersNumbers     = Letters + Numbers
	LowerLetterNumbers = LowerLetters + Numbers
	UpperLetterNumbers = UpperLetters + Numbers
)

func String(n int, set string) string {
	set1 := []rune(set)
	nSet := len(set1)
	r := make([]rune, n)
	for i := range r {
		r[i] = set1[rand.Intn(nSet)]
	}
	return string(r)
}
