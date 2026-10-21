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

func TestLoad_UpstreamURLs(t *testing.T) {
	t.Setenv("PBP_URL_FORMAT", "http://pbp.test/play_by_play_%d.parquet")
	t.Setenv("GAMES_CSV_URL", "http://games.test/games.csv")
	t.Setenv("ESPN_BASE_URL", "http://espn.test")
	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "http://pbp.test/play_by_play_%d.parquet", cfg.PBPURLFormat)
	assert.Equal(t, "http://games.test/games.csv", cfg.GamesCSVURL)
	assert.Equal(t, "http://espn.test", cfg.ESPNBaseURL)
}

func TestLoad_UpstreamURLsDefaultEmpty(t *testing.T) {
	cfg, err := Load()
	assert.NoError(t, err)
	assert.Empty(t, cfg.PBPURLFormat, "empty means the client's built-in default")
	assert.Empty(t, cfg.GamesCSVURL)
	assert.Empty(t, cfg.ESPNBaseURL)
}
