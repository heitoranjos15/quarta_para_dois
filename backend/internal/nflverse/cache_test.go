package nflverse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetGet(t *testing.T) {
	c := NewRedisCache("localhost:6379", "")
	c.Set("test", "value", time.Minute)
	val, ok := c.Get("test")
	assert.True(t, ok)
	assert.Equal(t, "value", val)
}

func TestCache_Delete(t *testing.T) {
	c := NewRedisCache("localhost:6379", "")
	c.Set("test", "value", time.Minute)
	c.Delete("test")
	_, ok := c.Get("test")
	assert.False(t, ok)
}

func TestCache_Miss(t *testing.T) {
	c := NewRedisCache("localhost:6379", "")
	_, ok := c.Get("nonexistent")
	assert.False(t, ok)
}