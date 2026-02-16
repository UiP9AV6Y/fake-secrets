package cache

import (
	"crypto/rsa"
	"fmt"
	"math/rand"
	"sync"
)

type RSACache struct {
	store map[string]*rsa.PrivateKey
	lock  sync.RWMutex
	rnd   *rand.Rand
}

func NewRSACache(rnd *rand.Rand) *RSACache {
	result := &RSACache{
		rnd: rnd,
	}

	return result
}

func (c *RSACache) Load(hostname string, length int) (key *rsa.PrivateKey, err error) {
	var ok bool
	pk := fmt.Sprintf("rsa-%s-%d", hostname, length)

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

	key, err = c.generateRSAKey(length)
	if err != nil {
		return
	}

	if c.store == nil {
		c.store = map[string]*rsa.PrivateKey{
			pk: key,
		}
	} else {
		c.store[pk] = key
	}

	return
}

func (c *RSACache) generateRSAKey(length int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(c.rnd, length)
	if err != nil {
		return nil, err
	}

	if err := key.Validate(); err != nil {
		return nil, err
	}

	return key, nil
}
