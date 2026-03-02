package cache

import (
	"hash/maphash"
	"sync"
)

type CacheLoader[V any] interface {
	Hash(maphash.Seed) uint64
	Load() (V, error)
}

type Cacher[V any] interface {
	Load(CacheLoader[V]) (V, error)
}

type mapCacher[V any] struct {
	store map[uint64]V
	lock  sync.RWMutex
	seed  maphash.Seed
}

func NewCacher[V any]() Cacher[V] {
	seed := maphash.MakeSeed()
	result := &mapCacher[V]{
		seed: seed,
	}

	return result
}

func (c *mapCacher[V]) Load(r CacheLoader[V]) (value V, err error) {
	var ok bool
	pk := r.Hash(c.seed)

	c.lock.RLock()
	value, ok = c.store[pk]
	c.lock.RUnlock()
	if ok {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok = c.store[pk]
	if ok {
		return
	}

	value, err = r.Load()
	if err != nil {
		return
	}

	if c.store == nil {
		c.store = map[uint64]V{
			pk: value,
		}
	} else {
		c.store[pk] = value
	}

	return
}
