package internal

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]CacheEntry
	mu	 sync.Mutex
	interval time.Duration
}

type CacheEntry struct {
	createdAt time.Time
	val	  []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		interval: interval,
	}
	return &cache
}
