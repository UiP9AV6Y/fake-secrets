package fake

import (
	"math/rand"

	"github.com/google/uuid"
)

func generateRandomUUID(rnd *rand.Rand) ([]byte, error) {
	result, err := uuid.NewRandomFromReader(rnd)
	if err != nil {
		return nil, err
	}

	return result.MarshalText()
}

func generateSeededUUID(seed string) ([]byte, error) {
	result := uuid.NewMD5(uuid.NameSpaceX500, []byte(seed))

	return result.MarshalText()
}
