package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type Cache struct {
	mut          sync.Mutex
	reapInterval time.Duration
	ticker       *time.Ticker
	data         map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{val: val, createdAt: time.Now()}

	c.mut.Lock()
	defer c.mut.Unlock()
	c.data[key] = entry
	fmt.Printf("CACHE add: %v\n", key)

}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()

	entry, ok := c.data[key]
	if ok {
		fmt.Printf("CACHE hit: %v\n", key)
		return entry.val, true
	}
	fmt.Printf("CACHE miss: %v\n", key)

	return nil, false
}

func (c *Cache) reapLoop() {
	for range c.ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	oldestAllowed := time.Now().Add(-c.reapInterval)

	c.mut.Lock()
	defer c.mut.Unlock()

	for key, entry := range c.data {
		if entry.createdAt.Before(oldestAllowed) {
			fmt.Printf("CACHE clearing: %v\n", key)
			delete(c.data, key)
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{reapInterval: interval, ticker: time.NewTicker(interval), data: make(map[string]cacheEntry)}
	go c.reapLoop()
	return c
}
