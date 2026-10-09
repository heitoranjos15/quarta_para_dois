package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("PORT", "8080")
	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, "8080", cfg.Port)
}