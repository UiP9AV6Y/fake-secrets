package cache

import (
	"crypto/rand"
	"hash/maphash"
	"io"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/UiP9AV6Y/fake-secrets/internal/hash"
)

type TOTPLoader struct {
	Issuer      string
	AccountName string
	Period      uint
	SecretSize  uint
	Algorithm   hash.Algorithm
	Random      io.Reader
}

func (l *TOTPLoader) Hash(seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)
	_, _ = h.WriteString(l.Issuer)
	_, _ = h.WriteString(l.AccountName)
	_, _ = h.WriteString(l.Algorithm.String())
	_, _ = h.Write(Uint64Bytes(uint64(l.SecretSize)))
	_, _ = h.Write(Uint64Bytes(uint64(l.Period)))

	return h.Sum64()
}

func (l *TOTPLoader) Load() (key *otp.Key, err error) {
	opts := totp.GenerateOpts{
		Issuer:      l.Issuer,
		AccountName: l.AccountName,
		SecretSize:  l.SecretSize,
		Algorithm:   l.Algorithm.OTPAlgorithm(),
		Period:      l.Period,
	}

	if l.Random == nil {
		opts.Rand = rand.Reader
	} else {
		opts.Rand = l.Random
	}

	key, err = totp.Generate(opts)
	if err != nil {
		return nil, err
	}

	return
}
