package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"quarta-para-dois/backend/internal/aggregation"
	"quarta-para-dois/backend/internal/nflverse"
)

type espnClientInterface interface {
	GetSchedule(ctx context.Context, season, week int) ([]nflverse.Game, error)
	GetGame(ctx context.Context, season, week int, gameID string) (*nflverse.Game, error)
	GetTeamSchedule(ctx context.Context, teamAbbr string, season int) ([]nflverse.Game, error)
}

type APIResponse[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Cached bool `json:"cached"`
	Season int  `json:"season"`
}

type GameNotes struct {
	Weather  *WeatherInfo `json:"weather,omitempty"`
	Injuries []InjuryInfo `json:"injuries,omitempty"`
	Vegas    *VegasInfo   `json:"vegas,omitempty"`
}

type WeatherInfo struct {
	Temperature int    `json:"temperature"`
	Wind        int    `json:"wind"`
	Condition   string `json:"condition"`
	Stadium     string `json:"stadium"`
	Roof        string `json:"roof"`
}

type InjuryInfo struct {
	Player   string `json:"player"`
	Team     string `json:"team"`
	Position string `json:"position"`
	Status   string `json:"status"`
	Detail   string `json:"detail"`
}

type VegasInfo struct {
	Spread float64 `json:"spread"`
	Total  float64 `json:"total"`
	HomeML int     `json:"home_ml"`
	AwayML int     `json:"away_ml"`
}

// PlaysSource serves one game's play-by-play.
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

func SeasonsListHandler(w http.ResponseWriter, r *http.Request) {
	currentYear := time.Now().Year()
	if time.Now().Month() < 9 {
		currentYear--
	}

	seasons := make([]int, 0, currentYear-1998)
	for y := 1999; y <= currentYear; y++ {
		seasons = append(seasons, y)
	}

	respond(w, seasons, false, currentYear, http.StatusOK)
}

func SeasonsWeeksHandler(w http.ResponseWriter, r *http.Request) {
	yearStr := chi.URLParam(r, "year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1999 {
		respondError(w, "invalid season", http.StatusBadRequest)
		return
	}

	// Regular season weeks (1-18) + playoffs
	maxWeek := 18
	if year >= 2021 {
		maxWeek = 18
	} else if year >= 1978 {
		maxWeek = 16
	}

	weeks := make([]int, maxWeek)
	for i := 1; i <= maxWeek; i++ {
		weeks[i-1] = i
	}

	respond(w, weeks, false, year, http.StatusOK)
}

func SeasonsGamesHandler(espnClient espnClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		yearStr := chi.URLParam(r, "year")
		weekStr := chi.URLParam(r, "week")

		year, err := strconv.Atoi(yearStr)
		if err != nil || year < 1999 {
			respondError(w, "invalid season", http.StatusBadRequest)
			return
		}

		week, err := strconv.Atoi(weekStr)
		if err != nil || week < 1 {
			respondError(w, "invalid week", http.StatusBadRequest)
			return
		}

		games, err := espnClient.GetSchedule(r.Context(), year, week)
		if err != nil {
			slog.Error("fetch schedule", "season", year, "week", week, "err", err)
			respondError(w, "failed to fetch schedules", http.StatusServiceUnavailable)
			return
		}

		respond(w, games, false, year, http.StatusOK)
	}
}

func GameDetailHandler(espnClient espnClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		season, _, _ := parseGameID(gameID)

		ctx := r.Context()
		game, err := espnClient.GetGame(ctx, season, 0, gameID) // week=0 means search all weeks
		if err != nil {
			slog.Error("fetch game", "game_id", gameID, "err", err)
			respondError(w, "failed to fetch game", http.StatusServiceUnavailable)
			return
		}

		if game == nil {
			respondError(w, "game not found", http.StatusNotFound)
			return
		}

		respond(w, game, false, season, http.StatusOK)
	}
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

func GameNotesHandler(espnClient espnClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		season, _, _ := parseGameID(gameID)

		ctx := r.Context()
		game, err := espnClient.GetGame(ctx, season, 0, gameID)
		if err != nil {
			slog.Error("fetch game", "game_id", gameID, "err", err)
			respondError(w, "failed to fetch game", http.StatusServiceUnavailable)
			return
		}

		if game == nil {
			respondError(w, "game not found", http.StatusNotFound)
			return
		}

		notes := GameNotes{
			Weather: &WeatherInfo{
				Temperature: game.Temp,
				Wind:        game.Wind,
				Condition:   "Clear",
				Stadium:     game.Stadium,
				Roof:        game.Roof,
			},
			Injuries: []InjuryInfo{},
			Vegas: &VegasInfo{
				Spread: 0,
				Total:  0,
			},
		}

		respond(w, notes, false, season, http.StatusOK)
	}
}

func TeamGamesHandler(espnClient espnClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamAbbr := chi.URLParam(r, "abbr")
		seasonStr := r.URL.Query().Get("season")

		season := time.Now().Year()
		if time.Now().Month() < 9 {
			season--
		}
		if seasonStr != "" {
			if s, err := strconv.Atoi(seasonStr); err == nil {
				season = s
			}
		}

		ctx := r.Context()
		games, err := espnClient.GetTeamSchedule(ctx, teamAbbr, season)
		if err != nil {
			slog.Error("fetch team schedule", "team", teamAbbr, "season", season, "err", err)
			respondError(w, "failed to fetch team schedule", http.StatusServiceUnavailable)
			return
		}

		respond(w, games, false, season, http.StatusOK)
	}
}
