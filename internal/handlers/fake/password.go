package fake

import (
	"math/rand"
)

func generateRandomPool(upper, lower, numeric, special bool) []byte {
	s := make([]byte, 0, len(RandomPoolUpper)+len(RandomPoolLower)+len(RandomPoolNumeric)+len(RandomPoolSpecial))
	if upper {
		s = append(s, RandomPoolUpper...)
	}

	if lower {
		s = append(s, RandomPoolLower...)
	}

	if numeric {
		s = append(s, RandomPoolNumeric...)
	}

	if special {
		s = append(s, RandomPoolSpecial...)
	}

	return s
}

func generatePassword(rnd *rand.Rand, length int, pool []byte) []byte {
	b := make([]byte, length)

	j := len(pool)
	for i := range b {
		b[i] = pool[rnd.Intn(j)]
	}

	return b
}
