package games

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/cache"
	"quarta-para-dois/backend/internal/espn"
	"quarta-para-dois/backend/internal/nflverse"
)

type fakeIndex struct{ ix *nflverse.GameIndex }

func (f fakeIndex) GameIndex(context.Context) (*nflverse.GameIndex, error) { return f.ix, nil }

// fakeESPN serves events from the espn package's fixtures and records every call.
type fakeESPN struct {
	t           *testing.T
	scoreboards map[[3]int][]espn.Event
	team        []espn.Event
	calls       [][3]int
}

func (f *fakeESPN) Scoreboard(_ context.Context, season, seasonType, week int) ([]espn.Event, error) {
	f.calls = append(f.calls, [3]int{season, seasonType, week})
	return f.scoreboards[[3]int{season, seasonType, week}], nil
}

func (f *fakeESPN) TeamSchedule(context.Context, string, int) ([]espn.Event, error) {
	return f.team, nil
}

func sampleIndex(t *testing.T) fakeIndex {
	t.Helper()
	fh, err := os.Open("../nflverse/testdata/games_sample.csv")
	require.NoError(t, err)
	defer func() { _ = fh.Close() }()
	ix, err := nflverse.ParseGamesCSV(fh)
	require.NoError(t, err)
	return fakeIndex{ix}
}

// eventsFromFixture runs an ESPN payload through a real espn.Client, so the test data
// has exactly the shape the service sees in production.
func eventsFromFixture(t *testing.T, file string) []espn.Event {
	t.Helper()
	body, err := os.ReadFile("../espn/testdata/" + file)
	require.NoError(t, err)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	c := espn.NewClient(cache.NewRedis(miniredis.RunT(t).Addr(), "", false), espn.Options{BaseURL: srv.URL})
	events, err := c.Scoreboard(context.Background(), 2026, 2, 1) // both payload shapes parse the same way
	require.NoError(t, err)
	return events
}

func newFake(t *testing.T) *fakeESPN {
	return &fakeESPN{
		t: t,
		scoreboards: map[[3]int][]espn.Event{
			{2026, 2, 6}: eventsFromFixture(t, "scoreboard_week6_final.json"),
		},
		team: eventsFromFixture(t, "team_schedule.json"),
	}
}

func ids(gs []Game) []string {
	out := make([]string, len(gs))
	for i, g := range gs {
		out[i] = g.GameID
	}
	return out
}

func TestService_Week_UsesNflverseIDs(t *testing.T) {
	svc := NewService(sampleIndex(t), newFake(t))

	gs, err := svc.Week(context.Background(), 2026, 6)
	require.NoError(t, err)

	assert.Contains(t, ids(gs), "2026_06_WAS_SF")
	assert.Contains(t, ids(gs), "2026_06_ARI_LA")
	for _, g := range gs {
		switch g.GameID {
		case "2026_06_WAS_SF":
			assert.Equal(t, "WAS", g.AwayTeam)
			assert.Equal(t, "SF", g.HomeTeam)
		case "2026_06_ARI_LA":
			assert.Equal(t, "LA", g.HomeTeam)
		}
	}
}

func TestService_Week_UnmatchedEventStillListed(t *testing.T) {
	f := newFake(t)
	f.scoreboards[[3]int{2026, 2, 6}] = []espn.Event{{
		ID: "999", Week: 6, Status: espn.StatusScheduled, HomeAbbr: "LAR", AwayAbbr: "WSH",
		Kickoff: time.Date(2026, 10, 25, 17, 0, 0, 0, time.UTC),
	}}
	svc := NewService(sampleIndex(t), f)

	gs, err := svc.Week(context.Background(), 2026, 6)

	require.NoError(t, err)
	require.Len(t, gs, 1)
	assert.Equal(t, "2026_06_WAS_LA", gs[0].GameID)
	assert.Equal(t, "LA", gs[0].HomeTeam)
	assert.Equal(t, "WAS", gs[0].AwayTeam)
	assert.Equal(t, "2026-10-25T17:00:00Z", gs[0].Kickoff)
}

func TestService_Game_UnmatchedEventPageLoads(t *testing.T) {
	f := newFake(t)
	f.scoreboards[[3]int{2026, 2, 6}] = []espn.Event{{
		ID: "999", Week: 6, Status: espn.StatusScheduled, HomeAbbr: "LAR", AwayAbbr: "WSH",
	}}
	svc := NewService(sampleIndex(t), f)

	g, err := svc.Game(context.Background(), "2026_06_WAS_LA")

	require.NoError(t, err)
	assert.Equal(t, "999", g.ESPNID)
}

func TestService_Game(t *testing.T) {
	f := newFake(t)
	svc := NewService(sampleIndex(t), f)

	g, err := svc.Game(context.Background(), "2026_06_DAL_GB")
	require.NoError(t, err)

	assert.Equal(t, [][3]int{{2026, 2, 6}}, f.calls, "one scoreboard request, no 18-week scan")
	assert.Equal(t, "final", g.Status)
	require.NotNil(t, g.HomeScore)
	assert.Equal(t, 23, *g.HomeScore)
	assert.Equal(t, "grass", g.Surface)
	assert.Equal(t, "Lambeau Field", g.Stadium)
	assert.Equal(t, "401873007", g.ESPNID)
}

func TestService_Game_NotFound(t *testing.T) {
	svc := NewService(sampleIndex(t), newFake(t))

	_, err := svc.Game(context.Background(), "2026_06_NYG_NYJ")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_Game_Playoffs(t *testing.T) {
	f := newFake(t)
	f.scoreboards[[3]int{2026, 3, 5}] = []espn.Event{{ID: "401873999", Week: 5, Status: espn.StatusScheduled, HomeAbbr: "SF", AwayAbbr: "KC"}}
	svc := NewService(sampleIndex(t), f)

	g, err := svc.Game(context.Background(), "2026_22_KC_SF")

	require.NoError(t, err)
	assert.Equal(t, [][3]int{{2026, 3, 5}}, f.calls)
	assert.Equal(t, "SB", g.GameType)
}

func TestService_TeamGames_IDs(t *testing.T) {
	idx := sampleIndex(t)
	svc := NewService(idx, newFake(t))

	gs, err := svc.TeamGames(context.Background(), "WAS", 2026)
	require.NoError(t, err)

	require.NotEmpty(t, gs)
	for _, id := range ids(gs) {
		assert.False(t, strings.Contains(id, "_00_"), id)
		_, ok := idx.ix.ByGameID(id)
		assert.True(t, ok, "%s must come from the index", id)
	}
}

func TestService_Weeks_PerSeason(t *testing.T) {
	svc := NewService(sampleIndex(t), newFake(t))

	weeks, err := svc.Weeks(context.Background(), 2020)
	require.NoError(t, err)
	assert.True(t, slices.Contains(weeks, 17) && slices.Contains(weeks, 18))

	none, err := svc.Weeks(context.Background(), 1990)
	require.NoError(t, err)
	assert.Empty(t, none)
}
