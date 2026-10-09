package nflverse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetGet(t *testing.T) {
	c := NewRedisCache("localhost:6379", "")
	c.Set("test", "value", time.Minute)
	var val string
	ok := c.Get("test", &val)
	assert.True(t, ok)
	assert.Equal(t, "value", val)
}

func TestCache_Delete(t *testing.T) {
	c := NewRedisCache("localhost:6379", "")
	c.Set("test", "value", time.Minute)
	c.Delete("test")
	var val string
	ok := c.Get("test", &val)
	assert.False(t, ok)
}

func TestCache_Miss(t *testing.T) {
	c := NewRedisCache("localhost:6379", "")
	var val string
	ok := c.Get("nonexistent", &val)
	assert.False(t, ok)
}