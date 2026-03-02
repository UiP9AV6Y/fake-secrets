package cache

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"hash/maphash"
	"io"
	"math/big"

	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
)

type CertLoader struct {
	Parent   *x509.Certificate
	Template *x509.Certificate
	Key      interface{}
	Random   io.Reader
}

func (l *CertLoader) Hash(seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)

	if l.Parent != nil {
		hashCertificate(&h, l.Parent)
	}

	hashCertificate(&h, l.Template)

	return h.Sum64()
}

func hashCertificate(h *maphash.Hash, cert *x509.Certificate) {
	_, _ = h.WriteString(cert.Subject.String())
	_, _ = h.WriteString(cert.NotBefore.String())
	_, _ = h.WriteString(cert.NotAfter.String())
	_, _ = h.WriteString(cert.KeyUsage.String())

	for _, u := range cert.ExtKeyUsage {
		_, _ = h.WriteString(u.String())
	}

	for _, n := range cert.DNSNames {
		_, _ = h.WriteString(n)
	}

	for _, a := range cert.EmailAddresses {
		_, _ = h.WriteString(a)
	}

	for _, a := range cert.IPAddresses {
		_, _ = h.Write([]byte(a))
	}

	for _, u := range cert.URIs {
		_, _ = h.WriteString(u.String())
	}
}

func (l *CertLoader) Load() (cert crypto.Certificate, err error) {
	var r io.Reader
	if l.Random == nil {
		r = rand.Reader
	} else {
		r = l.Random
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	l.Template.SerialNumber, err = rand.Int(r, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	parent := l.Parent
	if parent == nil {
		parent = l.Template
	}

	switch k := l.Key.(type) {
	case *rsa.PrivateKey:
		cert, err = x509.CreateCertificate(r, l.Template, parent, &k.PublicKey, l.Key)
	case *ecdsa.PrivateKey:
		cert, err = x509.CreateCertificate(r, l.Template, parent, &k.PublicKey, l.Key)
	case ed25519.PrivateKey:
		cert, err = x509.CreateCertificate(r, l.Template, parent, k.Public().(ed25519.PublicKey), l.Key)
	default:
		return nil, fmt.Errorf("unable to generate certificate for key type %T", l.Key)
	}

	if err != nil {
		return nil, err
	}

	return
}
