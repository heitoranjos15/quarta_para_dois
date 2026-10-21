package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type payload struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func newTestRedis(t *testing.T) (*Redis, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	return NewRedis(mr.Addr(), "", false), mr
}

func TestRedis_SetGet(t *testing.T) {
	c, _ := newTestRedis(t)
	ctx := context.Background()

	require.NoError(t, c.Set(ctx, "k", payload{Name: "a", Count: 2}, time.Minute))

	var got payload
	ok, err := c.Get(ctx, "k", &got)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, payload{Name: "a", Count: 2}, got)
}

func TestRedis_Miss(t *testing.T) {
	c, _ := newTestRedis(t)

	var got payload
	ok, err := c.Get(context.Background(), "missing", &got)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestRedis_TTL(t *testing.T) {
	c, mr := newTestRedis(t)
	ctx := context.Background()
	ttl := time.Minute
	require.NoError(t, c.Set(ctx, "k", payload{Name: "a"}, ttl))

	mr.FastForward(ttl + time.Second)

	var got payload
	ok, err := c.Get(ctx, "k", &got)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestRedis_Down(t *testing.T) {
	c, mr := newTestRedis(t)
	mr.Close()
	ctx := context.Background()

	var got payload
	ok, err := c.Get(ctx, "k", &got)
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Error(t, c.Set(ctx, "k", payload{}, time.Minute))
}
