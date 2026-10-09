package nflverse

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, val any, ttl time.Duration)
	Delete(key string)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisCache{client: rdb}
}

func (c *RedisCache) Get(key string) (any, bool) {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false
		}
		return nil, false
	}

	var val any
	if err := json.Unmarshal(data, &val); err != nil {
		return nil, false
	}
	return val, true
}

func (c *RedisCache) Set(key string, val any, ttl time.Duration) {
	ctx := context.Background()
	data, err := json.Marshal(val)
	if err != nil {
		return
	}
	c.client.Set(ctx, key, data, ttl)
}

func (c *RedisCache) Delete(key string) {
	ctx := context.Background()
	c.client.Del(ctx, key)
}