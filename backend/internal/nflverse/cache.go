package nflverse

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(key string, dest any) bool
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

func (c *RedisCache) Get(key string, dest any) bool {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false
		}
		return false
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false
	}
	return true
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