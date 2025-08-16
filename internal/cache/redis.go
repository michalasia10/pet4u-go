package cache

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache[T any] struct {
	client *redis.Client
	prefix string
	hits   int64
	misses int64
}

func NewRedisCache[T any](client *redis.Client, prefix string) *RedisCache[T] {
	if prefix == "" {
		prefix = "cache"
	}
	return &RedisCache[T]{
		client: client,
		prefix: prefix,
	}
}

func (c *RedisCache[T]) Get(ctx context.Context, key string) (T, bool, error) {
	var zero T
	fullKey := c.prefix + ":" + key

	data, err := c.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			atomic.AddInt64(&c.misses, 1)
			return zero, false, nil
		}
		return zero, false, err
	}

	var result T
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		c.client.Del(ctx, fullKey)
		atomic.AddInt64(&c.misses, 1)
		return zero, false, nil
	}

	atomic.AddInt64(&c.hits, 1)
	return result, true, nil
}

func (c *RedisCache[T]) Set(ctx context.Context, key string, data T, ttl time.Duration) error {
	fullKey := c.prefix + ":" + key

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, fullKey, jsonData, ttl).Err()
}

func (c *RedisCache[T]) Clear(ctx context.Context) error {
	pattern := c.prefix + ":*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

func (c *RedisCache[T]) Stats() Stats {
	return Stats{
		Hits:   atomic.LoadInt64(&c.hits),
		Misses: atomic.LoadInt64(&c.misses),
		Size:   -1,
	}
}
