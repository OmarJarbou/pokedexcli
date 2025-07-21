package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	CacheMap map[string]cacheEntry
	mutex    sync.RWMutex
	duration time.Duration
}

func NewCache(duration time.Duration) *Cache {
	newCacheMap := make(map[string]cacheEntry)
	var newMutex sync.RWMutex
	newCache := &Cache{
		CacheMap: newCacheMap,
		mutex: newMutex,
		duration: duration,
	}
	go newCache.reapLoop()
	return newCache
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newCacheEntry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.CacheMap[key] = newCacheEntry
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	fetchedCacheEntry, ok := c.CacheMap[key]
	if !ok {
		return nil, ok
	}
	return fetchedCacheEntry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.duration) // ticks every "duration" of time
	defer ticker.Stop()
	for range ticker.C { // when a tick recieved on the channel
			c.mutex.Lock()
			for key, entry := range c.CacheMap {
				livingTime := time.Now().Sub(entry.createdAt) // amount of time the entry have been in the cache
				if livingTime >= c.duration {
					delete(c.CacheMap, key)
				}
			}
			c.mutex.Unlock()
	}
}