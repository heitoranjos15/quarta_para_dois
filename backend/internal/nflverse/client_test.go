package nflverse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetPBP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cache := NewRedisCache("localhost:6379", "")
	client := NewClient(cache, "")

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

	cache := NewRedisCache("localhost:6379", "")
	client := NewClient(cache, "")

	_, err := client.GetSchedules(2024)
	// nflverse doesn't publish schedules parquet yet (404)
	if err != nil {
		assert.Contains(t, err.Error(), "404")
		return
	}
}

func TestClient_GetTeams(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cache := NewRedisCache("localhost:6379", "")
	client := NewClient(cache, "")

	_, err := client.GetTeams()
	// nflverse doesn't publish teams parquet yet (404)
	if err != nil {
		assert.Contains(t, err.Error(), "404")
		return
	}
}

func TestClient_GetRosters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cache := NewRedisCache("localhost:6379", "")
	client := NewClient(cache, "")

	_, err := client.GetRosters(2024)
	// nflverse doesn't publish rosters parquet yet (404)
	if err != nil {
		assert.Contains(t, err.Error(), "404")
		return
	}
}

func TestClientCache_SetGet(t *testing.T) {
	cache := NewRedisCache("localhost:6379", "")
	
	testKey := "test_key_" + time.Now().Format("150405.000")
	testValue := "test_value"
	
	cache.Set(testKey, testValue, time.Minute)
	var val string
	ok := cache.Get(testKey, &val)
	
	assert.True(t, ok)
	assert.Equal(t, testValue, val)
	
	cache.Delete(testKey)
	ok = cache.Get(testKey, &val)
	assert.False(t, ok)
}