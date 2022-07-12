package cache

import (
	"post-loader-service/pkg/config"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type Cache interface {
	Set(key string, item interface{})
	Get(key string) (interface{}, bool)
}

type inMemoryCache struct {
	cache *cache.Cache
}

func NewCache(cfg config.Config) Cache {
	cacheExpiration := cfg.GetDuration("app.cache.expiration") * time.Minute

	c := cache.NewFrom(cacheExpiration, cache.NoExpiration, make(map[string]cache.Item, 20))

	return &inMemoryCache{cache: c}
}

func (c *inMemoryCache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *inMemoryCache) Set(key string, item interface{}) {
	c.cache.Set(key, item, 0)
}
