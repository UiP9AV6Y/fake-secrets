package fake

import (
	"crypto/sha256"
	"hash/crc32"
	"math/rand"

	"github.com/jxskiss/base62"
)

var APIKeyEntropySize = 30

func generateRandomAPIKey(rnd *rand.Rand, label []byte) ([]byte, error) {
	pool := generateRandomPool(true, true, true, false)
	entropy := generatePassword(rnd, APIKeyEntropySize, pool)

	return concatAPIKey(entropy, label), nil
}

func generateSeededAPIKey(seed, label []byte) ([]byte, error) {
	h := sha256.New()
	_, _ = h.Write(seed)
	entropy := base62.Encode(h.Sum(nil))[:APIKeyEntropySize]

	return concatAPIKey(entropy, label), nil
}

func concatAPIKey(entropy, label []byte) []byte {
	sum := crc32.ChecksumIEEE(entropy)
	check := base62.FormatUint(uint64(sum))

	result := make([]byte, 0, len(label)+1+len(entropy)+6)
	result = append(result, label...)
	result = append(result, '_')
	result = append(result, entropy...)

	if left := 6 - len(check); left > 0 {
		pad := []byte{'0', '0', '0', '0', '0', '0'}[:left]
		result = append(result, pad...)
	}

	result = append(result, check...)

	return result
}
