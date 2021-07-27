package storage

import (
	"github.com/hashicorp/golang-lru"

	"github.com/MontFerret/worker/pkg/caching"
)

type InMemoryCache struct {
	store *lru.Cache
}

func NewCache(opts ...caching.Option) (caching.Cache, error) {
	options := caching.NewOptions(opts...)

	store, err := lru.New(int(options.Size))

	if err != nil {
		return nil, err
	}

	return &InMemoryCache{store}, nil
}

func (cache *InMemoryCache) Set(key, value interface{}) {
	cache.store.Add(key, value)
}

func (cache *InMemoryCache) Get(key interface{}) (interface{}, bool) {
	return cache.store.Get(key)
}

func (cache *InMemoryCache) Contains(key interface{}) bool {
	return cache.store.Contains(key)
}
