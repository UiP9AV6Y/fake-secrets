package cache

import (
	"crypto/ed25519"
	"crypto/rand"
	"hash/maphash"
	"io"
)

type ED25519Loader struct {
	Hostname string
	Random   io.Reader
}

func (l *ED25519Loader) Hash(seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)
	_, _ = h.WriteString(l.Hostname)

	return h.Sum64()
}

func (l *ED25519Loader) Load() (key ed25519.PrivateKey, err error) {
	if l.Random == nil {
		_, key, err = ed25519.GenerateKey(rand.Reader)
	} else {
		_, key, err = ed25519.GenerateKey(l.Random)
	}

	if err != nil {
		return nil, err
	}

	return
}
