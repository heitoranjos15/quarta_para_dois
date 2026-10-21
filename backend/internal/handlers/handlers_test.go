package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/aggregation"
	"quarta-para-dois/backend/internal/games"
	"quarta-para-dois/backend/internal/nflverse"
)

type mockPlaysSource struct {
	pbpData []nflverse.Play
	err     error
}

// GamePlays returns the plays of one game, or ErrGameNotFound when there are none.
func (m *mockPlaysSource) GamePlays(ctx context.Context, gameID string) ([]nflverse.Play, error) {
	if m.err != nil {
		return nil, m.err
	}
	var out []nflverse.Play
	for _, p := range m.pbpData {
		if p.GameID == gameID {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil, nflverse.ErrGameNotFound
	}
	return out, nil
}

type fakeGameService struct {
	games []games.Game
	err   error
}

func (f *fakeGameService) Seasons(context.Context) ([]int, error) { return []int{2023, 2024}, f.err }
func (f *fakeGameService) Weeks(context.Context, int) ([]int, error) {
	return []int{1, 2, 17, 18, 19}, f.err
}
func (f *fakeGameService) Week(context.Context, int, int) ([]games.Game, error) {
	return f.games, f.err
}
func (f *fakeGameService) TeamGames(context.Context, string, int) ([]games.Game, error) {
	return f.games, f.err
}
func (f *fakeGameService) Game(_ context.Context, id string) (games.Game, error) {
	if f.err != nil {
		return games.Game{}, f.err
	}
	for _, g := range f.games {
		if g.GameID == id {
			return g, nil
		}
	}
	return games.Game{}, games.ErrNotFound
}

func testGame() games.Game {
	return games.Game{
		GameID: "2024_01_DEN_LAC", Season: 2024, Week: 1, Status: "final",
		HomeTeam: "LAC", AwayTeam: "DEN", Stadium: "SoFi Stadium", Roof: "dome", Surface: "turf", Indoor: true,
	}
}

func setupTestRouter() *chi.Mux {
	client := &mockPlaysSource{
		pbpData: []nflverse.Play{
			{
				PlayID:               1,
				GameID:               "2024_01_DEN_LAC",
				Quarter:              1,
				Down:                 1,
				YardsToGo:            10,
				PlayType:             "pass",
				Description:          "Herbert pass complete to Williams for 15 yards",
				EPA:                  f64(1.2),
				Posteam:              "LAC",
				Defteam:              "DEN",
				Yardline100:          75,
				YardsGained:          15,
				PassAttempt:          1,
				RushAttempt:          0,
				CompletePass:         1,
				Touchdown:            0,
				PasserPlayerName:     "J.Herbert",
				ReceiverPlayerName:   "M.Williams",
				RusherPlayerName:     "",
				TacklerPlayerName:    "P.Surtain",
				ThirdDownConverted:   0,
				GameSecondsRemaining: 3600,
				Drive:                1,
				HomeTeam:             "LAC",
				AwayTeam:             "DEN",
			},
			{
				PlayID:               2,
				GameID:               "2024_01_DEN_LAC",
				Quarter:              1,
				Down:                 1,
				YardsToGo:            10,
				PlayType:             "run",
				Description:          "Ekeler rush for 5 yards",
				EPA:                  f64(0.3),
				Posteam:              "LAC",
				Defteam:              "DEN",
				Yardline100:          60,
				YardsGained:          5,
				PassAttempt:          0,
				RushAttempt:          1,
				CompletePass:         0,
				Touchdown:            0,
				PasserPlayerName:     "",
				ReceiverPlayerName:   "",
				RusherPlayerName:     "A.Ekeler",
				TacklerPlayerName:    "J.Simmons",
				ThirdDownConverted:   0,
				GameSecondsRemaining: 3560,
				Drive:                1,
				HomeTeam:             "LAC",
				AwayTeam:             "DEN",
			},
		},
	}

	gs := &fakeGameService{games: []games.Game{testGame()}}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/seasons", SeasonsListHandler(gs))
		r.Get("/seasons/{year}/weeks", SeasonsWeeksHandler(gs))
		r.Get("/seasons/{year}/weeks/{week}/games", SeasonsGamesHandler(gs))
		r.Get("/games/{gameID}", GameDetailHandler(gs))
		r.Get("/games/{gameID}/stats", GameStatsHandler(client))
		r.Get("/games/{gameID}/plays", GamePlaysHandler(client))
		r.Get("/games/{gameID}/notes", GameNotesHandler(gs))
		r.Get("/teams/{abbr}/games", TeamGamesHandler(gs))
	})
	return r
}

func TestSeasonsListHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/seasons", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]int]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
	assert.Greater(t, len(resp.Data), 0)
	assert.Contains(t, resp.Data, 2024)
}

func TestSeasonsWeeksHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/seasons/2024/weeks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]int]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, 1)
}

func TestSeasonsGamesHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/seasons/2024/weeks/1/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]games.Game]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "2024_01_DEN_LAC", resp.Data[0].GameID)
}

func TestGameDetailHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_DEN_LAC", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[games.Game]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "2024_01_DEN_LAC", resp.Data.GameID)
}

func TestGameStatsHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_DEN_LAC/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[aggregation.TeamComparisonStats]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NotNil(t, resp.Data)
	assert.Equal(t, "LAC", resp.Data.HomeTeam)
	assert.Equal(t, "DEN", resp.Data.AwayTeam)
	assert.Greater(t, resp.Data.HomeOffense.Plays, 0)
}

func TestGamePlaysHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_DEN_LAC/plays", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]aggregation.QuarterPlays]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NotNil(t, resp.Data)
	assert.Len(t, resp.Data, 1) // Only Q1 in test data
	assert.Len(t, resp.Data[0].Plays, 2)
}

func TestGameNotesHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_DEN_LAC/notes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[GameNotes]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, GameNotes{Stadium: "SoFi Stadium", Roof: "dome", Surface: "turf", Indoor: true}, resp.Data)
}

func TestGameNotesHandler_NoPlaceholders(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/games/2024_01_DEN_LAC/notes", nil))

	var raw struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &raw))
	assert.NotContains(t, raw.Data, "vegas")
	assert.NotContains(t, raw.Data, "condition")
	assert.ElementsMatch(t, []string{"stadium", "roof", "surface", "indoor"}, keys(raw.Data))
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestGameDetailHandler_NotFound(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/games/2024_01_XXX_YYY", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGameDetailHandler_InvalidID(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/games/bad", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSeasonsGamesHandler_UpstreamError(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/api/seasons/{year}/weeks/{week}/games", SeasonsGamesHandler(&fakeGameService{err: assert.AnError}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/seasons/2024/weeks/1/games", nil))

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestTeamGamesHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/teams/LAC/games?season=2024", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]games.Game]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
}

func TestGameStatsHandler_NotFound(t *testing.T) {
	client := &mockPlaysSource{
		pbpData: []nflverse.Play{}, // Empty - no plays for this game
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/games/{gameID}/stats", GameStatsHandler(client))
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_DEN_LAC/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGameStatsHandler_Error(t *testing.T) {
	client := &mockPlaysSource{
		err: assert.AnError,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/games/{gameID}/stats", GameStatsHandler(client))
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_DEN_LAC/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func f64(v float64) *float64 { return &v }

func statsRouter(src PlaysSource) http.Handler {
	r := chi.NewRouter()
	r.Get("/api/games/{gameID}/stats", GameStatsHandler(src))
	return r
}

func TestGameStatsHandler_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	statsRouter(&mockPlaysSource{}).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/games/not-an-id/stats", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGameStatsHandler_NotAvailableYet(t *testing.T) {
	src := &mockPlaysSource{err: nflverse.ErrGameNotFound}
	w := httptest.NewRecorder()
	statsRouter(src).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/games/2026_06_WAS_SF/stats", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"game not available yet"}`, w.Body.String())
}

func TestGameStatsHandler_UpstreamErrorIs503(t *testing.T) {
	src := &mockPlaysSource{err: assert.AnError}
	w := httptest.NewRecorder()
	statsRouter(src).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/games/2026_06_WAS_SF/stats", nil))

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}
