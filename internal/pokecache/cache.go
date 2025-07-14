package pokecache

import "time"

func NewCache(interval time.Duration) Cache {
	cache := Cache{}
	cache.entries = make(map[string]cacheEntry)
	cache.interval = interval
	go cache.reapLoop(interval)
	return cache
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.entries[key]
	if exists {
		return value.val, exists
	}
	return nil, exists

}

func (c Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		c.mu.Unlock()
	}

}
