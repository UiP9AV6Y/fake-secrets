package cache

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"hash/maphash"
	"io"

	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
)

type ECDSALoader struct {
	Hostname string
	Curve    crypto.ECDSACurve
	Random   io.Reader
}

func (l *ECDSALoader) Hash(seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)
	_, _ = h.WriteString(l.Hostname)
	_, _ = h.WriteString(l.Curve.String())

	return h.Sum64()
}

func (l *ECDSALoader) Load() (key *ecdsa.PrivateKey, err error) {
	ell := l.Curve.Curve()
	if ell == nil {
		return nil, fmt.Errorf("no curve size available for %q", l.Curve)
	}

	if l.Random == nil {
		key, err = ecdsa.GenerateKey(ell, rand.Reader)
	} else {
		key, err = ecdsa.GenerateKey(ell, l.Random)
	}

	if err != nil {
		return nil, err
	}

	return
}
