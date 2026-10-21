package nflverse

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quarta-para-dois/backend/internal/cache"
)

func TestClient_GetSchedules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := cache.NewRedis("localhost:6379", "", false)
	client := NewClient(c, "", Options{})

	_, err := client.GetSchedules(2024)
	// nflverse doesn't publish schedules parquet yet (404)
	if err != nil {
		assert.Contains(t, err.Error(), "404")
		return
	}
}
