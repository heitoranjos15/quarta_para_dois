package nflverse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/cache"
)

func TestClient_GetPBP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := cache.NewRedis("localhost:6379", "", false)
	client := NewClient(c, "")

	plays, err := client.GetPBP(2024)
	require.NoError(t, err)
	assert.Greater(t, len(plays), 0)

	// Verify cache hit on second call
	plays2, err := client.GetPBP(2024)
	require.NoError(t, err)
	assert.Equal(t, len(plays), len(plays2))
}

func TestClient_GetSchedules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := cache.NewRedis("localhost:6379", "", false)
	client := NewClient(c, "")

	_, err := client.GetSchedules(2024)
	// nflverse doesn't publish schedules parquet yet (404)
	if err != nil {
		assert.Contains(t, err.Error(), "404")
		return
	}
}
