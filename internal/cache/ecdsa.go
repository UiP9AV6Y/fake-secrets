package cache

import (
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"sync"

	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
)

type ECDSACache struct {
	store map[string]*ecdsa.PrivateKey
	lock  sync.RWMutex
	rnd   *rand.Rand
}

func NewECDSACache(rnd *rand.Rand) *ECDSACache {
	result := &ECDSACache{
		rnd: rnd,
	}

	return result
}

func (c *ECDSACache) Load(hostname string, curve crypto.ECDSACurve) (key *ecdsa.PrivateKey, err error) {
	var ok bool
	pk := fmt.Sprintf("ecdsa-%s-%s", hostname, curve)

	c.lock.RLock()
	key, ok = c.store[pk]
	c.lock.RUnlock()
	if ok {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	key, ok = c.store[pk]
	if ok {
		return
	}

	key, err = c.generateECDSAKey(curve)
	if err != nil {
		return
	}

	if c.store == nil {
		c.store = map[string]*ecdsa.PrivateKey{
			pk: key,
		}
	} else {
		c.store[pk] = key
	}

	return
}

func (c *ECDSACache) generateECDSAKey(curve crypto.ECDSACurve) (*ecdsa.PrivateKey, error) {
	ell := curve.Curve()
	if ell == nil {
		return nil, fmt.Errorf("no curve size available for %q", curve)
	}

	return ecdsa.GenerateKey(ell, c.rnd)
}
