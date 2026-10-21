package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	_ = os.Setenv("REDIS_ADDR", "localhost:6379")
	_ = os.Setenv("PORT", "8080")
	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, "8080", cfg.Port)
}

func TestLoad_RedisTLS(t *testing.T) {
	t.Setenv("REDIS_TLS", "true")
	cfg, err := Load()
	assert.NoError(t, err)
	assert.True(t, cfg.RedisTLS)
}
