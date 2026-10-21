package nflverse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/cache"
)

func loadSampleIndex(t *testing.T) *GameIndex {
	t.Helper()
	f, err := os.Open("testdata/games_sample.csv")
	require.NoError(t, err)
	defer func() { _ = f.Close() }()
	ix, err := ParseGamesCSV(f)
	require.NoError(t, err)
	return ix
}

func TestParseGamesCSV_Lookups(t *testing.T) {
	ix := loadSampleIndex(t)

	ref, ok := ix.ByESPNID("401873008")
	require.True(t, ok)
	assert.Equal(t, "2026_06_WAS_SF", ref.GameID)
	assert.Equal(t, "SF", ref.HomeTeam)
	assert.Equal(t, "WAS", ref.AwayTeam)
	assert.Equal(t, "outdoors", ref.Roof)

	ref, ok = ix.ByGameID("2026_06_DAL_GB")
	require.True(t, ok)
	assert.Equal(t, "401873007", ref.ESPNID)
	assert.Equal(t, "Lambeau Field", ref.Stadium)
	assert.Equal(t, 2026, ref.Season)
	assert.Equal(t, 6, ref.Week)

	_, ok = ix.ByESPNID("")
	assert.False(t, ok, "a game with an empty espn column is not in the ESPN map")
	_, ok = ix.ByGameID("2026_07_NYG_NYJ")
	assert.True(t, ok, "but it is still found by game id")
}

func TestGameIndex_Weeks(t *testing.T) {
	ix := loadSampleIndex(t)

	assert.Contains(t, ix.Weeks(2020), 17)
	assert.Contains(t, ix.Weeks(2020), 18)
	assert.Contains(t, ix.Weeks(2026), 6)
	assert.Contains(t, ix.Weeks(2026), 22)
	assert.Equal(t, []int{2020, 2026}, ix.Seasons())
	assert.Empty(t, ix.Weeks(1990))
}

func TestGameIndex_ESPNWeek(t *testing.T) {
	ix := loadSampleIndex(t)

	cases := []struct {
		season, week         int
		seasonType, espnWeek int
		ok                   bool
	}{
		{2026, 6, 2, 6, true},
		{2020, 18, 3, 1, true},
		{2026, 22, 3, 5, true}, // the Super Bowl is ESPN week 5, not 4
		{2026, 40, 0, 0, false},
	}
	for _, tc := range cases {
		st, w, ok := ix.ESPNWeek(tc.season, tc.week)
		assert.Equal(t, tc.ok, ok, "%d/%d", tc.season, tc.week)
		assert.Equal(t, tc.seasonType, st, "%d/%d", tc.season, tc.week)
		assert.Equal(t, tc.espnWeek, w, "%d/%d", tc.season, tc.week)
	}
}

func TestClient_GameIndexCachedInMemory(t *testing.T) {
	csvData, err := os.ReadFile("testdata/games_sample.csv")
	require.NoError(t, err)
	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_, _ = w.Write(csvData)
	}))
	defer srv.Close()

	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	c := NewClient(cache.NewRedis("localhost:1", "", false), "", Options{
		GamesCSVURL: srv.URL + "/games.csv",
		Now:         func() time.Time { return now },
	})
	ctx := context.Background()

	_, err = c.GameIndex(ctx)
	require.NoError(t, err)
	_, err = c.GameIndex(ctx)
	require.NoError(t, err)
	assert.Equal(t, int32(1), requests.Load())

	now = now.Add(25 * time.Hour)
	_, err = c.GameIndex(ctx)
	require.NoError(t, err)
	assert.Equal(t, int32(2), requests.Load())
}
