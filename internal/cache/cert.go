package cache

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"hash/maphash"
	"math/big"
	"math/rand"
	"sync"
)

type CertCache struct {
	store  map[uint64][]byte
	lock   sync.RWMutex
	rnd    *rand.Rand
	seed   maphash.Seed
	parent *x509.Certificate
}

func NewCertCache(rnd *rand.Rand, parent *x509.Certificate) *CertCache {
	seed := maphash.MakeSeed()
	result := &CertCache{
		rnd:    rnd,
		seed:   seed,
		parent: parent,
	}

	return result
}

func (c *CertCache) Load(template *x509.Certificate, key interface{}) (cert []byte, err error) {
	var ok bool
	pk := certHash(template, c.seed)

	c.lock.RLock()
	cert, ok = c.store[pk]
	c.lock.RUnlock()
	if ok {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	cert, ok = c.store[pk]
	if ok {
		return
	}

	cert, err = c.generateCertificate(template, key)
	if err != nil {
		return
	}

	if c.store == nil {
		c.store = map[uint64][]byte{
			pk: cert,
		}
	} else {
		c.store[pk] = cert
	}

	return
}

func (c *CertCache) generateCertificate(template *x509.Certificate, key interface{}) ([]byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := crand.Int(c.rnd, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template.SerialNumber = serialNumber
	parent := c.parent
	if parent == nil {
		parent = template
	}

	switch k := key.(type) {
	case *rsa.PrivateKey:
		return x509.CreateCertificate(c.rnd, template, parent, &k.PublicKey, key)
	case *ecdsa.PrivateKey:
		return x509.CreateCertificate(c.rnd, template, parent, &k.PublicKey, key)
	case ed25519.PrivateKey:
		return x509.CreateCertificate(c.rnd, template, parent, k.Public().(ed25519.PublicKey), key)
	default:
		return nil, fmt.Errorf("unable to generate certificate for key type %T", key)
	}
}

func certHash(cert *x509.Certificate, seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)

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

	return h.Sum64()
}
