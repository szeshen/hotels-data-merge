package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheInterface interface {
	Set(key string, value interface{}, d time.Duration)
	Get(key string) (interface{}, bool)
}

type GoCacheWrapper struct {
	cache *cache.Cache
}

func NewGoCacheWrapper(defaultExpiration, cleanupInterval time.Duration) *GoCacheWrapper {
	return &GoCacheWrapper{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Set stores a value in the cache
func (c *GoCacheWrapper) Set(key string, value interface{}, d time.Duration) {
	c.cache.Set(key, value, d)
}

// Get retrieves a value from the cache
func (c *GoCacheWrapper) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}
