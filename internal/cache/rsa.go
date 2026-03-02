package cache

import (
	"crypto/rand"
	"crypto/rsa"
	"hash/maphash"
	"io"
)

type RSALoader struct {
	Hostname string
	Length   int
	Random   io.Reader
}

func (l *RSALoader) Hash(seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)
	_, _ = h.WriteString(l.Hostname)
	_, _ = h.Write(Uint32Bytes(uint32(l.Length)))

	return h.Sum64()
}

func (l *RSALoader) Load() (key *rsa.PrivateKey, err error) {
	if l.Random == nil {
		key, err = rsa.GenerateKey(rand.Reader, l.Length)
	} else {
		key, err = rsa.GenerateKey(l.Random, l.Length)
	}

	if err != nil {
		return nil, err
	}

	if err := key.Validate(); err != nil {
		return nil, err
	}

	return
}
