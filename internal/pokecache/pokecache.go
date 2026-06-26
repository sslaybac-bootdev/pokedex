package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex
}

// Adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
}

// Gets an entry from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// Deletes an entry from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Identifies expired entries in the cache
func (c *Cache) getExpiredKeys(d time.Duration) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.entries))
	for k, e := range c.entries {
		if time.Since(e.createdAt) > d {
			keys = append(keys, k)
		}
	}
	return keys
}

// Each interval, identifies expired entries and deletes them
func (c *Cache) reapLoop(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for range ticker.C {
		keys := c.getExpiredKeys(d)
		for _, k := range keys {
			c.Delete(k)
		}
	}
}

// Creates and returns a new cache. Also sets up the cache clearing loop
func NewCache(d time.Duration) *Cache {
	new_cache := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go new_cache.reapLoop(d)
	return new_cache
}
