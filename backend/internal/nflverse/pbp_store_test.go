package nflverse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/cache"
)

const (
	gameA = "2026_06_DAL_GB"
	gameB = "2026_06_WAS_SF"
)

type pbpFixture struct {
	srv       *httptest.Server
	body      atomic.Pointer[[]byte] // the parquet file the server returns
	downloads atomic.Int32
	mr        *miniredis.Miniredis
	addr      string // kept so a client can still be built after mr.Close()
	now       time.Time
}

func newPBPFixture(t *testing.T) *pbpFixture {
	t.Helper()
	data := writePlaysParquet(t, []map[string]any{
		{"game_id": gameA, "play_id": 1, "fixed_drive": 1},
		{"game_id": gameA, "play_id": 2, "fixed_drive": 1},
		{"game_id": gameB, "play_id": 1, "fixed_drive": 1},
	})

	mr := miniredis.RunT(t)
	f := &pbpFixture{mr: mr, addr: mr.Addr(), now: time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)}
	f.body.Store(&data)
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.downloads.Add(1)
		_, _ = w.Write(*f.body.Load())
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *pbpFixture) client(opts Options) *Client {
	opts.PBPURLFormat = f.srv.URL + "/play_by_play_%d.parquet"
	if opts.Now == nil {
		opts.Now = func() time.Time { return f.now }
	}
	return NewClient(cache.NewRedis(f.addr, "", false), "", opts)
}

func TestGamePlays_DownloadsOnceAndSplits(t *testing.T) {
	f := newPBPFixture(t)
	c := f.client(Options{})
	ctx := context.Background()

	a1, err := c.GamePlays(ctx, gameA)
	require.NoError(t, err)
	a2, err := c.GamePlays(ctx, gameA)
	require.NoError(t, err)
	b, err := c.GamePlays(ctx, gameB)
	require.NoError(t, err)

	assert.Equal(t, int32(1), f.downloads.Load())
	assert.Len(t, a1, 2)
	assert.Len(t, a2, 2)
	require.Len(t, b, 1)
	assert.Equal(t, gameB, b[0].GameID)

	for _, id := range []string{gameA, gameB} {
		assert.True(t, f.mr.Exists("pbp:game:"+id), id)
		assert.Equal(t, 30*24*time.Hour, f.mr.TTL("pbp:game:"+id))
	}
}

func TestGamePlays_FromRedis(t *testing.T) {
	f := newPBPFixture(t)
	_, err := f.client(Options{}).GamePlays(context.Background(), gameA)
	require.NoError(t, err)
	require.Equal(t, int32(1), f.downloads.Load())

	plays, err := f.client(Options{}).GamePlays(context.Background(), gameA)
	require.NoError(t, err)

	assert.Len(t, plays, 2)
	assert.Equal(t, int32(1), f.downloads.Load())
}

func TestGamePlays_UnknownGameThrottled(t *testing.T) {
	f := newPBPFixture(t)
	now := f.now
	c := f.client(Options{SeasonRefetch: time.Hour, Now: func() time.Time { return now }})
	ctx := context.Background()
	const unknown = "2026_06_NYG_NYJ"

	_, err := c.GamePlays(ctx, unknown)
	assert.ErrorIs(t, err, ErrGameNotFound)
	_, err = c.GamePlays(ctx, unknown)
	assert.ErrorIs(t, err, ErrGameNotFound)
	assert.Equal(t, int32(1), f.downloads.Load())

	now = now.Add(time.Hour + time.Second)
	_, err = c.GamePlays(ctx, unknown)
	assert.ErrorIs(t, err, ErrGameNotFound)
	assert.Equal(t, int32(2), f.downloads.Load())
}

func TestGamePlays_RedisDown(t *testing.T) {
	f := newPBPFixture(t)
	f.mr.Close()

	plays, err := f.client(Options{}).GamePlays(context.Background(), gameA)

	require.NoError(t, err)
	assert.Len(t, plays, 2)
}

func TestGamePlays_MemoryBounded(t *testing.T) {
	f := newPBPFixture(t)
	c := f.client(Options{MemoryGames: 1})
	ctx := context.Background()

	_, err := c.GamePlays(ctx, gameA)
	require.NoError(t, err)
	_, err = c.GamePlays(ctx, gameB)
	require.NoError(t, err)

	f.srv.Close()
	f.mr.Close()

	_, err = c.GamePlays(ctx, gameA)
	assert.Error(t, err, "A was evicted from memory")
	b, err := c.GamePlays(ctx, gameB)
	require.NoError(t, err)
	assert.Len(t, b, 1)
}

func TestGamePlays_RedisDownOtherGamesStillServed(t *testing.T) {
	f := newPBPFixture(t)
	f.mr.Close()
	c := f.client(Options{SeasonRefetch: time.Hour})
	ctx := context.Background()

	_, err := c.GamePlays(ctx, gameA)
	require.NoError(t, err)
	b, err := c.GamePlays(ctx, gameB)

	require.NoError(t, err, "a game that was in the file is not throttled just because Redis lost it")
	assert.Len(t, b, 1)
}

func TestGamePlays_RefetchSkipsUnchangedGames(t *testing.T) {
	f := newPBPFixture(t)
	now := f.now
	c := f.client(Options{SeasonRefetch: time.Hour, Now: func() time.Time { return now }})
	ctx := context.Background()
	_, err := c.GamePlays(ctx, gameA)
	require.NoError(t, err)
	f.mr.SetTTL(gamePlaysKey(gameA), time.Hour)

	now = now.Add(2 * time.Hour)
	_, err = c.GamePlays(ctx, "2026_06_NYG_NYJ")
	require.ErrorIs(t, err, ErrGameNotFound)

	require.Equal(t, int32(2), f.downloads.Load())
	assert.Equal(t, time.Hour, f.mr.TTL(gamePlaysKey(gameA)), "game A did not change, so it was not written again")
}

func TestGamePlays_RefetchRewritesChangedGames(t *testing.T) {
	f := newPBPFixture(t)
	now := f.now
	c := f.client(Options{SeasonRefetch: time.Hour, Now: func() time.Time { return now }})
	ctx := context.Background()
	_, err := c.GamePlays(ctx, gameA)
	require.NoError(t, err)

	grown := writePlaysParquet(t, []map[string]any{
		{"game_id": gameA, "play_id": 1, "fixed_drive": 1},
		{"game_id": gameA, "play_id": 2, "fixed_drive": 1},
		{"game_id": gameB, "play_id": 1, "fixed_drive": 1},
		{"game_id": gameB, "play_id": 2, "fixed_drive": 1},
	})
	f.body.Store(&grown)
	now = now.Add(2 * time.Hour)
	_, _ = c.GamePlays(ctx, "2026_06_NYG_NYJ")

	b, err := f.client(Options{}).GamePlays(ctx, gameB) // fresh client: read from Redis
	require.NoError(t, err)
	assert.Len(t, b, 2)
}
