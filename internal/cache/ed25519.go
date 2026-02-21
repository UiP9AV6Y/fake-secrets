package cache

import (
	"crypto/ed25519"
	"fmt"
	"math/rand"
	"sync"
)

type ED25519Cache struct {
	store map[string]ed25519.PrivateKey
	lock  sync.RWMutex
	rnd   *rand.Rand
}

func NewED25519Cache(rnd *rand.Rand) *ED25519Cache {
	result := &ED25519Cache{
		rnd: rnd,
	}

	return result
}

func (c *ED25519Cache) Load(hostname string) (key ed25519.PrivateKey, err error) {
	var ok bool
	pk := fmt.Sprintf("ed25519-%s", hostname)

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

	key, err = c.generateED25519Key()
	if err != nil {
		return
	}

	if c.store == nil {
		c.store = map[string]ed25519.PrivateKey{
			pk: key,
		}
	} else {
		c.store[pk] = key
	}

	return
}

func (c *ED25519Cache) generateED25519Key() (ed25519.PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(c.rnd)

	return priv, err
}
