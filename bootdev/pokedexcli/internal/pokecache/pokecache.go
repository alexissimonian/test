package pokecache

import "time"

func NewCache() Cache {
	return Cache{
		cache: make(map[string]cacheEntry),
    }
}

func (c *Cache) Add(key string, val []byte) {
	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now().UTC(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
    cacheE, ok := c.cache[key]
    return cacheE.val, ok
}
