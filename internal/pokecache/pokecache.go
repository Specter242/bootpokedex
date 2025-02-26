package pokecache

import (
	"sync"
	"time"
)

var mux = &sync.RWMutex{}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache holds in-memory Pokemon cache data
type Cache struct {
	store map[string]cacheEntry
}

// NewCache creates and returns a new Cache instance
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		store: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

// Get returns a stored value identified by key
func (c *Cache) Get(key string) ([]byte, bool) {
	mux.RLock()
	defer mux.RUnlock()
	entry, exists := c.store[key]
	// if exists {
	//     fmt.Printf("Cache hit for key: %s\n", key)
	// } else {
	//     fmt.Printf("Cache miss for key: %s\n", key)
	// }
	return entry.val, exists
}

// Add stores a value in the cache with the given key
func (c *Cache) Add(key string, value []byte) {
	mux.Lock()
	defer mux.Unlock()
	c.store[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

// Delete removes the value associated with the key from the cache.
func (c *Cache) Delete(key string) {
	mux.Lock()
	defer mux.Unlock()
	delete(c.store, key)
}

func (c *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		c.clearExpired(interval)
	}
}

func (c *Cache) clearExpired(interval time.Duration) {
	mux.Lock()
	defer mux.Unlock()
	for key, entry := range c.store {
		if time.Since(entry.createdAt) > interval {
			// fmt.Printf("Removing expired cache entry for: %s\n", key)
			delete(c.store, key)
		}
	}
}
