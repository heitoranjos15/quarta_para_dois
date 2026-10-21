package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"quarta-para-dois/backend/internal/aggregation"
	"quarta-para-dois/backend/internal/games"
	"quarta-para-dois/backend/internal/nflverse"
)

// GameService serves the game list, schedule and per-game details.
type GameService interface {
	Seasons(ctx context.Context) ([]int, error)
	Weeks(ctx context.Context, season int) ([]int, error)
	Week(ctx context.Context, season, week int) ([]games.Game, error)
	Game(ctx context.Context, gameID string) (games.Game, error)
	TeamGames(ctx context.Context, abbr string, season int) ([]games.Game, error)
}

type APIResponse[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Cached bool `json:"cached"`
	Season int  `json:"season"`
}

// GameNotes holds only what the sources actually provide; nothing is filled with placeholders.
type GameNotes struct {
	Stadium string `json:"stadium"`
	Roof    string `json:"roof"`
	Surface string `json:"surface"`
	Indoor  bool   `json:"indoor"`
}

type PlaysSource interface {
	GamePlays(ctx context.Context, gameID string) ([]nflverse.Play, error)
}

// parseGameID accepts nflverse game IDs such as 2026_06_WAS_SF.
func parseGameID(id string) (season, week int, ok bool) {
	return nflverse.ParseGameID(id)
}

func respond[T any](w http.ResponseWriter, data T, cached bool, season int, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := APIResponse[T]{
		Data: data,
		Meta: Meta{Cached: cached, Season: season},
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("encode response", "err", err)
	}
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		slog.Error("encode error response", "err", err)
	}
}

// gamePlays validates the game ID and loads its plays, writing the error response itself.
func gamePlays(w http.ResponseWriter, r *http.Request, src PlaysSource) (gameID string, season int, plays []nflverse.Play, ok bool) {
	gameID = chi.URLParam(r, "gameID")
	season, _, valid := parseGameID(gameID)
	if !valid {
		respondError(w, "invalid game id", http.StatusBadRequest)
		return "", 0, nil, false
	}

	plays, err := src.GamePlays(r.Context(), gameID)
	if errors.Is(err, nflverse.ErrGameNotFound) || (err == nil && len(plays) == 0) {
		respondError(w, "game not available yet", http.StatusNotFound)
		return "", 0, nil, false
	}
	if err != nil {
		slog.Error("fetch play data", "game_id", gameID, "err", err)
		respondError(w, "failed to fetch play data", http.StatusServiceUnavailable)
		return "", 0, nil, false
	}
	return gameID, season, plays, true
}

func SeasonsListHandler(gs GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seasons, err := gs.Seasons(r.Context())
		if err != nil {
			slog.Error("list seasons", "err", err)
			respondError(w, "failed to fetch seasons", http.StatusServiceUnavailable)
			return
		}
		latest := 0
		if len(seasons) > 0 {
			latest = seasons[len(seasons)-1]
		}
		respond(w, nonNil(seasons), false, latest, http.StatusOK)
	}
}

func SeasonsWeeksHandler(gs GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year, err := strconv.Atoi(chi.URLParam(r, "year"))
		if err != nil || year < 1999 {
			respondError(w, "invalid season", http.StatusBadRequest)
			return
		}
		weeks, err := gs.Weeks(r.Context(), year)
		if err != nil {
			slog.Error("list weeks", "season", year, "err", err)
			respondError(w, "failed to fetch weeks", http.StatusServiceUnavailable)
			return
		}
		respond(w, nonNil(weeks), false, year, http.StatusOK)
	}
}

func SeasonsGamesHandler(gs GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year, err := strconv.Atoi(chi.URLParam(r, "year"))
		if err != nil || year < 1999 {
			respondError(w, "invalid season", http.StatusBadRequest)
			return
		}
		week, err := strconv.Atoi(chi.URLParam(r, "week"))
		if err != nil || week < 1 {
			respondError(w, "invalid week", http.StatusBadRequest)
			return
		}

		list, err := gs.Week(r.Context(), year, week)
		if errors.Is(err, games.ErrNotFound) {
			respondError(w, "week not found", http.StatusNotFound)
			return
		}
		if err != nil {
			slog.Error("fetch schedule", "season", year, "week", week, "err", err)
			respondError(w, "failed to fetch schedules", http.StatusServiceUnavailable)
			return
		}
		respond(w, nonNil(list), false, year, http.StatusOK)
	}
}

// loadGame validates the game ID and loads the game, writing the error response itself.
func loadGame(w http.ResponseWriter, r *http.Request, gs GameService) (games.Game, bool) {
	gameID := chi.URLParam(r, "gameID")
	if _, _, ok := parseGameID(gameID); !ok {
		respondError(w, "invalid game id", http.StatusBadRequest)
		return games.Game{}, false
	}
	game, err := gs.Game(r.Context(), gameID)
	if errors.Is(err, games.ErrNotFound) {
		respondError(w, "game not found", http.StatusNotFound)
		return games.Game{}, false
	}
	if err != nil {
		slog.Error("fetch game", "game_id", gameID, "err", err)
		respondError(w, "failed to fetch game", http.StatusServiceUnavailable)
		return games.Game{}, false
	}
	return game, true
}

func GameDetailHandler(gs GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if game, ok := loadGame(w, r, gs); ok {
			respond(w, game, false, game.Season, http.StatusOK)
		}
	}
}

func GameNotesHandler(gs GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		game, ok := loadGame(w, r, gs)
		if !ok {
			return
		}
		respond(w, GameNotes{Stadium: game.Stadium, Roof: game.Roof, Surface: game.Surface, Indoor: game.Indoor}, false, game.Season, http.StatusOK)
	}
}

func TeamGamesHandler(gs GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamAbbr := strings.ToUpper(chi.URLParam(r, "abbr"))

		season := time.Now().Year()
		if time.Now().Month() < 9 {
			season--
		}
		if s, err := strconv.Atoi(r.URL.Query().Get("season")); err == nil {
			season = s
		}

		list, err := gs.TeamGames(r.Context(), teamAbbr, season)
		if err != nil {
			slog.Error("fetch team schedule", "team", teamAbbr, "season", season, "err", err)
			respondError(w, "failed to fetch team schedule", http.StatusServiceUnavailable)
			return
		}
		respond(w, nonNil(list), false, season, http.StatusOK)
	}
}

func nonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func GameStatsHandler(src PlaysSource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, season, plays, ok := gamePlays(w, r, src)
		if !ok {
			return
		}
		respond(w, aggregation.ComputeTeamStats(plays), true, season, http.StatusOK)
	}
}

func GamePlaysHandler(src PlaysSource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, season, gamePlays, ok := gamePlays(w, r, src)
		if !ok {
			return
		}

		// Optional query params for filtering
		if quarter, err := strconv.Atoi(r.URL.Query().Get("quarter")); err == nil {
			gamePlays = filterPlays(gamePlays, func(p nflverse.Play) bool { return p.Quarter == quarter })
		}
		if drive, err := strconv.Atoi(r.URL.Query().Get("drive")); err == nil {
			gamePlays = filterPlays(gamePlays, func(p nflverse.Play) bool { return p.FixedDrive == drive })
		}

		respond(w, aggregation.ProcessPlays(gamePlays), true, season, http.StatusOK)
	}
}

func filterPlays(plays []nflverse.Play, keep func(nflverse.Play) bool) []nflverse.Play {
	var out []nflverse.Play
	for _, p := range plays {
		if keep(p) {
			out = append(out, p)
		}
	}
	return out
}

type DriveSummaryResponse struct {
	Quarters map[int]aggregation.QuarterDriveSummary `json:"quarters"`
	Overtime *aggregation.QuarterDriveSummary        `json:"overtime,omitempty"`
}

func GameDrivesHandler(src PlaysSource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, season, plays, ok := gamePlays(w, r, src)
		if !ok {
			return
		}

		// Get home and away teams from the first play
		quarters := aggregation.ComputeDrives(plays, plays[0].HomeTeam, plays[0].AwayTeam)
		respond(w, DriveSummaryResponse{Quarters: quarters}, true, season, http.StatusOK)
	}
}
