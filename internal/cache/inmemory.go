package cache

import (
	"context"
	"sync"
	"time"
)

type InMemoryCache[T any] struct {
	mu      sync.RWMutex
	data    map[string]*entry[T]
	stats   Stats
	maxSize int
}

type entry[T any] struct {
	Data      T
	ExpiresAt time.Time
}

func NewInMemoryCache[T any](maxSize int) *InMemoryCache[T] {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &InMemoryCache[T]{
		data:    make(map[string]*entry[T]),
		maxSize: maxSize,
	}
}

func (c *InMemoryCache[T]) Get(ctx context.Context, key string) (T, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zero T
	entry, exists := c.data[key]
	if !exists {
		c.stats.Misses++
		return zero, false, nil
	}

	if time.Now().After(entry.ExpiresAt) {
		// Entry expired
		c.stats.Misses++
		// Clean up expired entry (defer cleanup to avoid holding read lock)
		go c.cleanupExpired(key)
		return zero, false, nil
	}

	c.stats.Hits++
	return entry.Data, true, nil
}

func (c *InMemoryCache[T]) Set(ctx context.Context, key string, data T, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict old entries
	if len(c.data) >= c.maxSize {
		c.evictOldest()
	}

	c.data[key] = &entry[T]{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.stats.Size = int64(len(c.data))
	return nil
}

func (c *InMemoryCache[T]) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*entry[T])
	c.stats.Size = 0
	return nil
}

func (c *InMemoryCache[T]) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}

func (c *InMemoryCache[T]) cleanupExpired(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	c.stats.Size = int64(len(c.data))
}

func (c *InMemoryCache[T]) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range c.data {
		if first || entry.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.ExpiresAt
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)
	}
}
