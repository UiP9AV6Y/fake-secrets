package cache

import (
	"crypto/rand"
	"hash/maphash"
	"io"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"

	"github.com/UiP9AV6Y/fake-secrets/internal/hash"
)

type HOTPLoader struct {
	Issuer      string
	AccountName string
	SecretSize  uint
	Algorithm   hash.Algorithm
	Random      io.Reader
}

func (l *HOTPLoader) Hash(seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)
	_, _ = h.WriteString(l.Issuer)
	_, _ = h.WriteString(l.AccountName)
	_, _ = h.WriteString(l.Algorithm.String())
	_, _ = h.Write(Uint64Bytes(uint64(l.SecretSize)))

	return h.Sum64()
}

func (l *HOTPLoader) Load() (key *otp.Key, err error) {
	opts := hotp.GenerateOpts{
		Issuer:      l.Issuer,
		AccountName: l.AccountName,
		SecretSize:  l.SecretSize,
		Algorithm:   l.Algorithm.OTPAlgorithm(),
	}

	if l.Random == nil {
		opts.Rand = rand.Reader
	} else {
		opts.Rand = l.Random
	}

	key, err = hotp.Generate(opts)
	if err != nil {
		return nil, err
	}

	return
}
