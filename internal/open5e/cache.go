package open5e

import (
	"sync"
	"time"
)

type cacheEntry struct {
	data []byte
	expiresAt time.Time
}

type Cache struct {
	mu sync.Mutex
	entries map[string]cacheEntry
	ttl time.Duration
	stopCh chan struct{}
}

func NewCache(ttl, reapInterval time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]cacheEntry),
		ttl: ttl,
		stopCh: make(chan struct{}),
	}

	go cache.reapLoop(reapInterval)
	return cache
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry, ok := cache.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		delete(cache.entries, key)
		return nil, false
	}
	return entry.data, true
}

func (cache *Cache) Set(key string, data []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.entries[key] = cacheEntry{
		data: data,
		expiresAt: time.Now().Add(cache.ttl),
	}
}

func (cache *Cache) purgeExpired() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	now := time.Now()
	for key, entry := range cache.entries {
		if now.After(entry.expiresAt) {
			delete(cache.entries, key)
		}
	}
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cache.purgeExpired()
		case <-cache.stopCh:
			return
		}
	}
}

func (cache *Cache) Stop() {
	close(cache.stopCh)
}

var apiCache = NewCache(5*time.Minute, 1*time.Minute)

