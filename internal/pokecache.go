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

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := CacheEntry{
		val: val,
		createdAt: time.Now(),
	}
	c.cacheMap[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var entry CacheEntry
	exists := false
	if entry, exists = c.cacheMap[key]; !exists {
		return nil, false
	}
	return entry.val, true
}
