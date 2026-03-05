package fake

import (
	"io"

	"github.com/google/uuid"
)

func generateRandomUUID(rnd io.Reader) ([]byte, error) {
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
