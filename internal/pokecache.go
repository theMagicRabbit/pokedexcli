package internal

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]CacheEntry
	mu	 sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val	  []byte
}
