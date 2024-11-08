package pokecache

import "time"

type Cache struct {
    data map[string]cacheEntry
    
}

type cacheEntry struct {
    createdAt time.Time
    val []byte
}
