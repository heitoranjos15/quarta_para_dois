package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache stores JSON-encoded values with a TTL. Callers log a non-nil error
// and treat it as a miss.
type Cache interface {
	Get(ctx context.Context, key string, dest any) (bool, error) // (false, nil) on miss
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
}

type Redis struct {
	client *redis.Client
}

func NewRedis(addr, password string, useTLS bool) *Redis {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}
	if useTLS {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	return &Redis{client: redis.NewClient(opts)}
}

func (c *Redis) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Redis) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}
