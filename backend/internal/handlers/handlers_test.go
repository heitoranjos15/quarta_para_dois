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
	"quarta-para-dois/backend/internal/nflverse"
)

type mockNFLVerseClient struct {
	pbpData    []nflverse.Play
	shouldFail bool
	err        error
}

func (m *mockNFLVerseClient) GetPBP(season int) ([]nflverse.Play, error) {
	if m.shouldFail {
		return nil, m.err
	}
	return m.pbpData, nil
}

func (m *mockNFLVerseClient) GetSchedules(season int) ([]nflverse.Game, error) {
	if m.shouldFail {
		return nil, m.err
	}
	return []nflverse.Game{
		{
			GameID:   "2024_01_LAC_DEN",
			Season:   2024,
			Week:     1,
			HomeTeam: "LAC",
			AwayTeam: "DEN",
			HomeScore: 27,
			AwayScore: 20,
		},
	}, nil
}

func (m *mockNFLVerseClient) GetTeams() ([]nflverse.Team, error) {
	if m.shouldFail {
		return nil, m.err
	}
	return []nflverse.Team{
		{TeamAbbr: "LAC", TeamName: "Los Angeles Chargers"},
		{TeamAbbr: "DEN", TeamName: "Denver Broncos"},
	}, nil
}

func (m *mockNFLVerseClient) GetRosters(season int) ([]nflverse.Roster, error) {
	if m.shouldFail {
		return nil, m.err
	}
	return []nflverse.Roster{}, nil
}

func setupTestRouter() *chi.Mux {
	client := &mockNFLVerseClient{
		pbpData: []nflverse.Play{
			{
				PlayID:             1,
				GameID:             "2024_01_LAC_DEN",
				Quarter:            1,
				Down:               1,
				YardsToGo:          10,
				PlayType:           "pass",
				Description:        "Herbert pass complete to Williams for 15 yards",
				EPA:                1.2,
				Posteam:            "LAC",
				Defteam:            "DEN",
				Yardline100:        75,
				YardsGained:        15,
				PassAttempt:        1,
				RushAttempt:        0,
				CompletePass:       1,
				Touchdown:          0,
				PasserPlayerName:   "J.Herbert",
				ReceiverPlayerName: "M.Williams",
				RusherPlayerName:   "",
				TacklerPlayerName:  "P.Surtain",
				ThirdDownConverted: 0,
				GameSecondsRemaining: 3600,
				Drive:              1,
				HomeTeam:           "LAC",
				AwayTeam:           "DEN",
			},
			{
				PlayID:             2,
				GameID:             "2024_01_LAC_DEN",
				Quarter:            1,
				Down:               1,
				YardsToGo:          10,
				PlayType:           "run",
				Description:        "Ekeler rush for 5 yards",
				EPA:                0.3,
				Posteam:            "LAC",
				Defteam:            "DEN",
				Yardline100:        60,
				YardsGained:        5,
				PassAttempt:        0,
				RushAttempt:        1,
				CompletePass:       0,
				Touchdown:          0,
				PasserPlayerName:   "",
				ReceiverPlayerName: "",
				RusherPlayerName:   "A.Ekeler",
				TacklerPlayerName:  "J.Simmons",
				ThirdDownConverted: 0,
				GameSecondsRemaining: 3560,
				Drive:              1,
				HomeTeam:           "LAC",
				AwayTeam:           "DEN",
			},
		},
	}

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, nflverseClientKey, client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Route("/api", func(r chi.Router) {
		r.Get("/seasons", SeasonsListHandler)
		r.Get("/seasons/{year}/weeks", SeasonsWeeksHandler)
		r.Get("/seasons/{year}/weeks/{week}/games", SeasonsGamesHandler)
		r.Get("/games/{gameID}", GameDetailHandler(client))
		r.Get("/games/{gameID}/stats", GameStatsHandler(client))
		r.Get("/games/{gameID}/plays", GamePlaysHandler(client))
		r.Get("/games/{gameID}/notes", GameNotesHandler(client))
		r.Get("/teams/{abbr}/games", TeamGamesHandler(client))
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
	json.Unmarshal(w.Body.Bytes(), &resp)
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
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, 1)
}

func TestSeasonsGamesHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/seasons/2024/weeks/1/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]nflverse.Game]
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "2024_01_LAC_DEN", resp.Data[0].GameID)
}

func TestGameDetailHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_LAC_DEN", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[nflverse.Game]
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "2024_01_LAC_DEN", resp.Data.GameID)
}

func TestGameStatsHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_LAC_DEN/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[aggregation.TeamComparisonStats]
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.NotNil(t, resp.Data)
	assert.Equal(t, "LAC", resp.Data.HomeTeam)
	assert.Equal(t, "DEN", resp.Data.AwayTeam)
	assert.Greater(t, resp.Data.HomeOffense.Plays, 0)
}

func TestGamePlaysHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_LAC_DEN/plays", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]aggregation.QuarterPlays]
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.NotNil(t, resp.Data)
	assert.Len(t, resp.Data, 1) // Only Q1 in test data
	assert.Len(t, resp.Data[0].Plays, 2)
}

func TestGameNotesHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_LAC_DEN/notes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[GameNotes]
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
}

func TestTeamGamesHandler(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/teams/LAC/games?season=2024", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse[[]nflverse.Game]
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp.Data)
}

func TestGameStatsHandler_NotFound(t *testing.T) {
	client := &mockNFLVerseClient{
		pbpData:    []nflverse.Play{}, // Empty - no plays for this game
		shouldFail: false,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/games/{gameID}/stats", GameStatsHandler(client))
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_LAC_DEN/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGameStatsHandler_Error(t *testing.T) {
	client := &mockNFLVerseClient{
		shouldFail: true,
		err:        assert.AnError,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/games/{gameID}/stats", GameStatsHandler(client))
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games/2024_01_LAC_DEN/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}