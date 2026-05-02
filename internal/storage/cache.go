package storage

import (
	"github.com/MontFerret/ferret/v2"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/MontFerret/worker/pkg/caching"
)

type InMemoryCache struct {
	store *lru.Cache[string, *ferret.Plan]
}

func NewCache(opts ...caching.Option) (caching.Cache[*ferret.Plan], error) {
	options := caching.NewOptions(opts...)

	store, err := lru.New[string, *ferret.Plan](int(options.Size))

	if err != nil {
		return nil, err
	}

	return &InMemoryCache{store}, nil
}

func (cache *InMemoryCache) Set(key string, value *ferret.Plan) {
	cache.store.Add(key, value)
}

func (cache *InMemoryCache) Get(key string) (*ferret.Plan, bool) {
	return cache.store.Get(key)
}

func (cache *InMemoryCache) Contains(key string) bool {
	return cache.store.Contains(key)
}
