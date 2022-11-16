package storage

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/hashicorp/golang-lru/v2"

	"github.com/MontFerret/worker/pkg/caching"
)

type InMemoryCache struct {
	store *lru.Cache[string, *runtime.Program]
}

func NewCache(opts ...caching.Option) (caching.Cache[*runtime.Program], error) {
	options := caching.NewOptions(opts...)

	store, err := lru.New[string, *runtime.Program](int(options.Size))

	if err != nil {
		return nil, err
	}

	return &InMemoryCache{store}, nil
}

func (cache *InMemoryCache) Set(key string, value *runtime.Program) {
	cache.store.Add(key, value)
}

func (cache *InMemoryCache) Get(key string) (*runtime.Program, bool) {
	return cache.store.Get(key)
}

func (cache *InMemoryCache) Contains(key string) bool {
	return cache.store.Contains(key)
}
