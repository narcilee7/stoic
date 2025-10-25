package planner

import (
	"sync"
	"time"
)

type RuleCache struct {
	mu    sync.RWMutex
	cache map[string]cacheEntry
	ttl   time.Duration
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

func NewRuleCache(ttl time.Duration) *RuleCache {
	return &RuleCache{
		cache: make(map[string]cacheEntry),
		ttl:   ttl,
	}
}

func (c *RuleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.value, true
}

func (c *RuleCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}
