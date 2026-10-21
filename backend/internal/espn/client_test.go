package espn

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/cache"
)

type fakeESPN struct {
	srv      *httptest.Server
	requests atomic.Int32
	lastURL  atomic.Value // string
}

// newFakeESPN serves body for every request and counts them.
func newFakeESPN(t *testing.T, body []byte) *fakeESPN {
	t.Helper()
	f := &fakeESPN{}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.requests.Add(1)
		f.lastURL.Store(r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeESPN) client(t *testing.T) *Client {
	t.Helper()
	mr := miniredis.RunT(t)
	return NewClient(cache.NewRedis(mr.Addr(), "", false), Options{BaseURL: f.srv.URL})
}

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	require.NoError(t, err)
	return b
}

func findEvent(t *testing.T, events []Event, id string) Event {
	t.Helper()
	for _, e := range events {
		if e.ID == id {
			return e
		}
	}
	t.Fatalf("event %s not found", id)
	return Event{}
}

func TestScoreboard_ParsesWeek6(t *testing.T) {
	f := newFakeESPN(t, fixture(t, "scoreboard_week6_final.json"))

	events, err := f.client(t).Scoreboard(context.Background(), 2026, 2, 6)
	require.NoError(t, err)
	require.Len(t, events, 14)

	gb := findEvent(t, events, "401873007")
	assert.Equal(t, "GB", gb.HomeAbbr)
	assert.Equal(t, "DAL", gb.AwayAbbr)
	require.NotNil(t, gb.HomeScore)
	require.NotNil(t, gb.AwayScore)
	assert.Equal(t, 23, *gb.HomeScore)
	assert.Equal(t, 13, *gb.AwayScore)
	assert.Equal(t, StatusFinal, gb.Status)
	assert.True(t, gb.Kickoff.Equal(time.Date(2026, 10, 19, 0, 20, 0, 0, time.UTC)), gb.Kickoff)
	assert.Equal(t, 6, gb.Week)
	assert.Equal(t, "Lambeau Field", gb.Venue)
	assert.Equal(t, "Green Bay", gb.City)
	assert.Equal(t, "WI", gb.State)

	assert.Equal(t, "WSH", findEvent(t, events, "401873008").AwayAbbr)

	url, _ := f.lastURL.Load().(string)
	assert.Contains(t, url, "seasontype=2")
	assert.Contains(t, url, "week=6")
	assert.Contains(t, url, "dates=2026")
}

func TestScoreboard_ScheduledHasNoScores(t *testing.T) {
	f := newFakeESPN(t, fixture(t, "scoreboard_week7_scheduled.json"))

	events, err := f.client(t).Scoreboard(context.Background(), 2026, 2, 7)
	require.NoError(t, err)

	for _, e := range events {
		assert.Equal(t, StatusScheduled, e.Status, e.ID)
		assert.Nil(t, e.HomeScore, e.ID)
		assert.Nil(t, e.AwayScore, e.ID)
	}
}

func TestScoreboardTTL(t *testing.T) {
	now := time.Date(2026, 10, 25, 16, 30, 0, 0, time.UTC)
	final := Event{Status: StatusFinal}
	scheduled := Event{Status: StatusScheduled, Kickoff: now.Add(3 * time.Hour)}
	live := Event{Status: StatusInProgress}

	ttl, ok := scoreboardTTL([]Event{final, final}, now)
	assert.True(t, ok)
	assert.Equal(t, 30*24*time.Hour, ttl)

	ttl, ok = scoreboardTTL([]Event{final, scheduled}, now)
	assert.True(t, ok)
	assert.Equal(t, time.Hour, ttl)

	_, ok = scoreboardTTL([]Event{final, scheduled, live}, now)
	assert.False(t, ok, "a game in progress means nothing is cached")
}

func TestScoreboardTTL_EmptyIsNotCached(t *testing.T) {
	_, ok := scoreboardTTL(nil, time.Now())
	assert.False(t, ok, "an empty answer is not 'all final': ESPN may not have published the week yet")
}

func TestScoreboardTTL_EndsAtNextKickoff(t *testing.T) {
	now := time.Date(2026, 10, 25, 16, 30, 0, 0, time.UTC)
	soon := Event{Status: StatusScheduled, Kickoff: now.Add(30 * time.Minute)}
	later := Event{Status: StatusScheduled, Kickoff: now.Add(4 * time.Hour)}

	ttl, ok := scoreboardTTL([]Event{later, soon}, now)

	assert.True(t, ok)
	assert.Equal(t, 30*time.Minute, ttl)
}

func TestScoreboardTTL_PastKickoffStillScheduled(t *testing.T) {
	now := time.Date(2026, 10, 25, 16, 30, 0, 0, time.UTC)
	delayed := Event{Status: StatusScheduled, Kickoff: now.Add(-5 * time.Minute)}

	ttl, ok := scoreboardTTL([]Event{delayed}, now)

	assert.True(t, ok)
	assert.Equal(t, time.Minute, ttl, "a game that should have started is rechecked soon")
}

// withLiveGame marks the first event of a scoreboard payload as in progress.
func withLiveGame(t *testing.T, body []byte) []byte {
	t.Helper()
	var doc map[string]any
	require.NoError(t, json.Unmarshal(body, &doc))
	ev := doc["events"].([]any)[0].(map[string]any)
	setLive := func(status map[string]any) {
		typ := status["type"].(map[string]any)
		typ["state"], typ["completed"], typ["name"] = "in", false, "STATUS_IN_PROGRESS"
	}
	setLive(ev["status"].(map[string]any))
	setLive(ev["competitions"].([]any)[0].(map[string]any)["status"].(map[string]any))
	out, err := json.Marshal(doc)
	require.NoError(t, err)
	return out
}

func TestScoreboard_InProgressNotCached(t *testing.T) {
	f := newFakeESPN(t, withLiveGame(t, fixture(t, "scoreboard_week6_final.json")))
	c := f.client(t)

	first, err := c.Scoreboard(context.Background(), 2026, 2, 6)
	require.NoError(t, err)
	assert.Equal(t, StatusInProgress, first[0].Status)
	_, err = c.Scoreboard(context.Background(), 2026, 2, 6)
	require.NoError(t, err)

	assert.Equal(t, int32(2), f.requests.Load())
}

func TestScoreboard_FinalWeekCached(t *testing.T) {
	f := newFakeESPN(t, fixture(t, "scoreboard_week6_final.json"))
	c := f.client(t)

	_, err := c.Scoreboard(context.Background(), 2026, 2, 6)
	require.NoError(t, err)
	_, err = c.Scoreboard(context.Background(), 2026, 2, 6)
	require.NoError(t, err)

	assert.Equal(t, int32(1), f.requests.Load())
}

func TestScoreboard_BadDate(t *testing.T) {
	var doc map[string]any
	require.NoError(t, json.Unmarshal(fixture(t, "scoreboard_week6_final.json"), &doc))
	ev := doc["events"].([]any)[0].(map[string]any)
	ev["date"] = ""
	body, err := json.Marshal(doc)
	require.NoError(t, err)
	f := newFakeESPN(t, body)

	events, err := f.client(t).Scoreboard(context.Background(), 2026, 2, 6)

	require.NoError(t, err)
	assert.True(t, events[0].Kickoff.IsZero())
}

func TestTeamSchedule_ParsesEvents(t *testing.T) {
	f := newFakeESPN(t, fixture(t, "team_schedule.json"))

	events, err := f.client(t).TeamSchedule(context.Background(), "WAS", 2026)
	require.NoError(t, err)

	first := events[0]
	assert.Equal(t, "401872929", first.ID)
	assert.Equal(t, "PHI", first.HomeAbbr)
	assert.Equal(t, "WSH", first.AwayAbbr)
	assert.Equal(t, StatusFinal, first.Status)
	require.NotNil(t, first.HomeScore)
	assert.Equal(t, 24, *first.HomeScore)
	url, _ := f.lastURL.Load().(string)
	assert.Contains(t, url, "/teams/wsh/schedule")
	assert.Contains(t, url, "season=2026")
}

func TestTeamSlug(t *testing.T) {
	assert.Equal(t, "wsh", teamSlug("WAS"))
	assert.Equal(t, "lar", teamSlug("LA"))
	assert.Equal(t, "kc", teamSlug("KC"))
}

func TestScoreboard_EmptyNotCached(t *testing.T) {
	f := newFakeESPN(t, []byte(`{"events":[]}`))
	c := f.client(t)

	for i := 0; i < 2; i++ {
		events, err := c.Scoreboard(context.Background(), 2026, 2, 18)
		require.NoError(t, err)
		assert.Empty(t, events)
	}

	assert.Equal(t, int32(2), f.requests.Load())
}

func TestTeamSchedule_IgnoresOldClientCacheEntry(t *testing.T) {
	f := newFakeESPN(t, fixture(t, "team_schedule.json"))
	mr := miniredis.RunT(t)
	// The value the previous ESPN client stored under this key (a []nflverse.Game).
	require.NoError(t, mr.Set("espn:team:wsh:2026", `[{"GameID":"2026_01_WSH_PHI","Season":2026,"HomeTeam":"PHI","AwayTeam":"WSH","HomeScore":24,"AwayScore":22}]`))
	c := NewClient(cache.NewRedis(mr.Addr(), "", false), Options{BaseURL: f.srv.URL})

	events, err := c.TeamSchedule(context.Background(), "WAS", 2026)

	require.NoError(t, err)
	assert.Equal(t, int32(1), f.requests.Load())
	assert.Equal(t, "401872929", events[0].ID)
}
